import argparse
from fastapi import FastAPI
from pydantic import BaseModel
import hanlp
import uvicorn
import json
from typing import Tuple


# 加载模型（只加载一次）
HanLP = hanlp.load(
    hanlp.pretrained.mtl.CLOSE_TOK_POS_NER_SRL_DEP_SDP_CON_ELECTRA_BASE_ZH
)
ner = HanLP["ner/msra"]
ner.dict_tags = {
    # 电力设备类 (EQUIP)
    ("发电机",): ("S-EQUIP",),
    ("变压器",): ("S-EQUIP",),
    ("断路器",): ("S-EQUIP",),
    ("逆变器",): ("S-EQUIP",),
    ("整流器",): ("S-EQUIP",),
    ("电容器",): ("S-EQUIP",),
    ("智能电表",): ("S-EQUIP",),
    ("冷却塔",): ("S-EQUIP",),
    ("锅炉",): ("S-EQUIP",),
    ("涡轮机",): ("S-EQUIP",),
    # 设施类 (FAC)
    ("发电厂",): ("S-FAC",),
    ("变电站",): ("S-FAC",),
    ("输电线路",): ("B-FAC", "E-FAC"),
    ("配电线路",): ("B-FAC", "E-FAC"),
    ("充电桩",): ("S-FAC",),
    ("水电站",): ("S-FAC",),
    ("核电厂",): ("S-FAC",),
    # 技术术语 (TECH)
    ("欧姆定律",): ("S-TECH",),
    ("电磁感应",): ("S-TECH",),
    ("光伏效应",): ("S-TECH",),
    ("无功补偿",): ("S-TECH",),
    ("功率因数",): ("S-TECH",),
    ("特高压输电",): ("B-TECH", "E-TECH"),
    ("柔性输电",): ("B-TECH", "E-TECH"),
    ("电力电子",): ("B-TECH", "E-TECH"),
    ("继电保护",): ("B-TECH", "E-TECH"),
    # 参数指标 (PARAM)
    ("电压",): ("S-PARAM",),
    ("电流",): ("S-PARAM",),
    ("电阻",): ("S-PARAM",),
    ("频率",): ("S-PARAM",),
    ("功率",): ("S-PARAM",),
    ("谐波",): ("S-PARAM",),
    # 材料类 (MAT)
    ("铜",): ("S-MAT",),
    ("绝缘油",): ("S-MAT",),
    ("六氟化硫",): ("S-MAT",),
    ("超导材料",): ("S-MAT",),
    # 能源类型 (ENER)
    ("电能",): ("S-ENER",),
    ("风能",): ("S-ENER",),
    ("水能",): ("S-ENER",),
    ("太阳能",): ("S-ENER",),
    ("生物质能",): ("S-ENER",),
    # 组织机构 (ORG)
    ("SCADA",): ("S-ORG",),
    ("EMS",): ("S-ORG",),
    ("DMS",): ("S-ORG",),
    ("WAMS",): ("S-ORG",),
    # 安全类 (SAFE)
    ("接地保护",): ("B-SAFE", "E-SAFE"),
    ("漏电保护",): ("B-SAFE", "E-SAFE"),
    ("耐压试验",): ("B-SAFE", "E-SAFE"),
    ("防污闪",): ("B-SAFE", "E-SAFE"),
    # 操作类 (OP)
    ("负荷预测",): ("B-OP", "E-OP"),
    ("电力调度",): ("B-OP", "E-OP"),
    ("削峰填谷",): ("B-OP", "E-OP"),
    ("状态监测",): ("B-OP", "E-OP"),
    # 故障类 (FAULT)
    ("短路",): ("S-FAULT",),
    ("开路",): ("S-FAULT",),
    ("接地故障",): ("B-FAULT", "E-FAULT"),
    ("局部放电",): ("B-FAULT", "E-FAULT"),
    # 单位类 (UNIT)
    ("欧姆",): ("S-UNIT",),
    ("千伏",): ("S-UNIT",),
    ("kW",): ("S-UNIT",),
    ("kWh",): ("S-UNIT",),
}


app = FastAPI()


class BaseRsp(BaseModel):
    code: int
    msg: str


class NERReq(BaseModel):
    text: str


class NERRsp(BaseRsp):
    text: str


@app.post("/ner", response_model=NERRsp)
async def ner(req: NERReq):
    result = HanLP(req.text, tasks="ner").to_pretty()
    # result_json = json.dumps(result, ensure_ascii=False)
    return NERRsp(code=0, msg="success", text=result)


class KGCReq(BaseModel):
    head: str
    relation: str
    tail: str


class KGCRsp(BaseRsp):
    head: str
    relation: str
    tail: str


from kgc import model, infer_triples, benchmark


# KGC
def KGCModel(head: str, relation: str, tail: str) -> Tuple[str, str, str]:
    return infer_triples(model, head, relation, tail)


@app.post("/kgc", response_model=KGCRsp)
async def kgc(req: KGCReq):
    head, relation, tail = KGCModel(req.head, req.relation, req.tail)
    return KGCRsp(code=0, msg="success", head=head, relation=relation, tail=tail)


class KGCBenchmarkRsp(BaseRsp):
    hits_at_1: str
    mrr: str


@app.get("/kgc_benchmark", response_model=KGCBenchmarkRsp)
async def kgc_benchmark():
    hits_at_1, mrr = benchmark()
    return KGCBenchmarkRsp(
        code=0, msg="success", hits_at_1=str(hits_at_1), mrr=str(mrr)
    )


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", type=str, default="127.0.0.1", help="Host address")
    parser.add_argument("--port", type=int, default=8081, help="Port number")

    args = parser.parse_args()

    print(f"Starting NER service at http://{args.host}:{args.port}...")
    uvicorn.run(app, host=args.host, port=args.port)
