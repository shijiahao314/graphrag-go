import sys
import hanlp


def ner(text):
    # 加载 HanLP 模型
    HanLP = hanlp.load(
        hanlp.pretrained.mtl.CLOSE_TOK_POS_NER_SRL_DEP_SDP_CON_ELECTRA_BASE_ZH
    )
    # 执行命名实体识别任务
    result = HanLP(text, tasks="ner")
    # 输出结果
    print(result)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print('Usage: python ner.py "<text>"')
        sys.exit(1)
    input_text = sys.argv[1]
    ner(input_text)
