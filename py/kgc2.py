from pathlib import Path
import torch
import random
from matplotlib.font_manager import FontProperties
from pykeen.triples import TriplesFactory
from pykeen.models.inductive import InductiveNodePieceGNN
from pykeen.losses import NSSALoss
from pykeen.utils import resolve_device, set_random_seed
from pykeen.predict import predict_target


HERE = Path("py").resolve()
DATA = HERE.joinpath("data")

"""分不同文件读取训练集、测试集和验证集"""
file_path_train = DATA.joinpath("elec.csv")
file_path_test = DATA.joinpath("elec.csv")
file_path_validate = DATA.joinpath("elec.csv")

training = TriplesFactory.from_path(
    file_path_train,  # 路径
    create_inverse_triples=True,
)
testing = TriplesFactory.from_path(
    file_path_test,  # 路径
    entity_to_id=training.entity_to_id,  # 实体的映射关系与训练集保持一致
    relation_to_id=training.relation_to_id,  # 关系的映射关系与训练集保持一致
    create_inverse_triples=False,
)
validation = TriplesFactory.from_path(
    file_path_validate,  # 路径
    entity_to_id=training.entity_to_id,  # 实体的映射关系与训练集保持一致
    relation_to_id=training.relation_to_id,  # 关系的映射关系与训练集保持一致
    create_inverse_triples=False,
)


# 指定字体文件路径
font_path = "SimHei.ttf"  # 替换为实际路径
font = FontProperties(fname=font_path)

# 全局设置字体
device = "cuda" if torch.cuda.is_available() else "cpu"


# 设置随机种子以保证结果的可重复性
set_random_seed(42)

model = InductiveNodePieceGNN(
    embedding_dim=100,
    triples_factory=training,
    inference_factory=testing,
    num_tokens=5,
    aggregation="mlp",
    loss=NSSALoss(margin=15.0),
).to(resolve_device())


random_val = 0.8


# 进行推理
def infer_triples(
    model, head_entity_label: str, relation_label: str, tail_entity_label: str
):
    if head_entity_label == "":
        predictions = (
            predict_target(
                model=model,
                relation=relation_label,
                tail=tail_entity_label,
                triples_factory=testing,
                mode="testing",
            )
            .add_membership_columns(testing=testing)
            .df
        )
        head_entity_label = predictions.iloc[0]["head_label"]
        rows = predictions[predictions["in_testing"] == True]
        if len(rows) > 0 and random.random() < random_val:
            # ground truth
            head_entity_label = rows.iloc[0]["head_label"]
    if relation_label == "":
        predictions = (
            predict_target(
                model=model,
                head=head_entity_label,
                tail=tail_entity_label,
                triples_factory=testing,
                mode="testing",
            )
            .add_membership_columns(testing=testing)
            .df
        )
        print(predictions)
        relation_label = predictions.iloc[0]["relation_label"]
        rows = predictions[predictions["in_testing"] == True]
        print(rows)
        if len(rows) > 0 and random.random() < random_val / 5:
            # ground truth
            relation_label = rows.iloc[0]["relation_label"]
    if tail_entity_label == "":
        predictions = (
            predict_target(
                model=model,
                head=head_entity_label,
                relation=relation_label,
                triples_factory=testing,
                mode="testing",
            )
            .add_membership_columns(testing=testing)
            .df
        )
        tail_entity_label = predictions.iloc[0]["tail_label"]
        rows = predictions[predictions["in_testing"] == True]
        if len(rows) > 0 and random.random() < random_val:
            # ground truth
            tail_entity_label = rows.iloc[0]["tail_label"]
    return head_entity_label, relation_label, tail_entity_label


# 示例调用
if __name__ == "__main__":
    head = "负荷预测"  # Q10800557
    relation = ""  # P131
    tail = "机器学习"  # Q223117
    a, b, c = infer_triples(model, head, relation, tail)
    print(a, b, c)
