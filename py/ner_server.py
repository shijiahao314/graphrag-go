import argparse
from fastapi import FastAPI
from pydantic import BaseModel
import hanlp
import uvicorn
import json


# 加载模型（只加载一次）
HanLP = hanlp.load(
    hanlp.pretrained.mtl.CLOSE_TOK_POS_NER_SRL_DEP_SDP_CON_ELECTRA_BASE_ZH
)

app = FastAPI()


class NERReq(BaseModel):
    text: str


class BaseRsp(BaseModel):
    code: int
    msg: str


class NERRsp(BaseRsp):
    text: str


@app.post("/ner", response_model=NERRsp)
async def ner(req: NERReq):
    result = HanLP(req.text, tasks="ner").to_pretty()
    result_json = json.dumps(result, ensure_ascii=False)
    return NERRsp(code=0, msg="success", text=result_json)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", type=str, default="127.0.0.1", help="Host address")
    parser.add_argument("--port", type=int, default=8081, help="Port number")

    args = parser.parse_args()

    print(f"Starting NER service at http://{args.host}:{args.port}...")
    uvicorn.run(app, host=args.host, port=args.port)
