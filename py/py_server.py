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
    # 电力设备类
    ("检查", "变压器"): ("O", "S-EQUIP"),  # B:设备开始 M:设备中间 E:设备结束 S:单字设备
    ("更换", "断路器"): ("O", "S-EQUIP"),
    ("绝缘", "子"): ("B-COMPONENT", "E-COMPONENT"),
    # 电力技术术语
    ("过", "电压", "保护"): ("B-TECH", "M-TECH", "E-TECH"),
    ("谐波", "抑制"): ("B-TECH", "E-TECH"),
    ("短路", "电流"): ("B-TECH", "E-TECH"),
    # 组织机构
    ("国家", "电网"): ("B-ORG", "E-ORG"),  # 国家电网有限公司
    ("南方", "电网"): ("B-ORG", "E-ORG"),
    ("华能", "集团"): ("B-ORG", "E-ORG"),
    # 电力设施
    ("1000", "千伏", "变电站"): ("B-FAC", "I-FAC", "E-FAC"),
    ("输电", "线路"): ("B-FAC", "E-FAC"),
    ("配电", "网络"): ("B-FAC", "E-FAC"),
    # 安全操作
    ("接地", "线"): ("B-SAFETY", "E-SAFETY"),
    ("绝缘", "手套"): ("B-SAFETY", "E-SAFETY"),
    ("安全", "规程"): ("B-SAFETY", "E-SAFETY"),
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


from kgc import model, dataset, infer_tail_entity


# KGC
def KGCModel(head: str, relation: str, tail: str) -> Tuple[str, str, str]:
    return infer_tail_entity(model, dataset, head, relation, tail)


@app.post("/kgc", response_model=KGCRsp)
async def kgc(req: KGCReq):
    head, relation, tail = KGCModel(req.head, req.relation, req.tail)
    return KGCRsp(code=0, msg="success", head=head, relation=relation, tail=tail)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", type=str, default="127.0.0.1", help="Host address")
    parser.add_argument("--port", type=int, default=8081, help="Port number")

    args = parser.parse_args()

    print(f"Starting NER service at http://{args.host}:{args.port}...")
    uvicorn.run(app, host=args.host, port=args.port)
