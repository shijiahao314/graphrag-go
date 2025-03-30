from pykeen.predict import predict_target, predict_triples
from pykeen.models.inductive import InductiveNodePieceGNN
from pykeen.datasets.inductive.base import DisjointInductivePathDataset
from typing_extensions import Literal
from pathlib import Path
from pykeen.losses import NSSALoss
from pykeen.utils import resolve_device, set_random_seed
import logging

# 设置随机种子以保证结果的可重复性
set_random_seed(42)

HERE = Path("py").resolve()
DATA = HERE.joinpath("data")


# 禁用 PyKEEN 日志
logging.getLogger("pykeen").setLevel(logging.ERROR)

Size = Literal["small", "large"]


class InductiveLPDataset(DisjointInductivePathDataset):
    """An inductive link prediction dataset for the ILPC 2022 Challenge."""

    def __init__(self, size: Size = "small", **kwargs):
        super().__init__(
            transductive_training_path=DATA.joinpath(size, "train.txt"),
            inductive_inference_path=DATA.joinpath(size, "inference.txt"),
            inductive_validation_path=DATA.joinpath(size, "inference_validation.txt"),
            inductive_testing_path=DATA.joinpath(size, "inference_test.txt"),
            create_inverse_triples=True,
            eager=True,
            **kwargs,
        )


# 预加载模型，提高推理效率
def load_model(dataset_size: Size = "small"):
    dataset = InductiveLPDataset(size=dataset_size)
    model = InductiveNodePieceGNN(
        embedding_dim=100,
        triples_factory=dataset.transductive_training,
        inference_factory=dataset.inductive_inference,
        num_tokens=5,
        aggregation="mlp",
        loss=NSSALoss(margin=15.0),
    ).to(resolve_device())
    return model, dataset


# 进行推理
def infer_triples(
    model, dataset, head_entity_label: str, relation_label: str, tail_entity_label: str
):
    if head_entity_label == "":
        predictions = predict_target(
            model=model,
            relation=relation_label,
            tail=tail_entity_label,
            triples_factory=dataset.transductive_training,
            mode="training",
        )
        return predictions.df.iloc[0]["head_label"], relation_label, tail_entity_label
    if relation_label == "":
        predictions = predict_target(
            model=model,
            head=head_entity_label,
            tail=tail_entity_label,
            triples_factory=dataset.transductive_training,
            mode="training",
        )
        return (
            head_entity_label,
            predictions.df.iloc[0]["relation_label"],
            tail_entity_label,
        )
    if tail_entity_label == "":
        predictions = predict_target(
            model=model,
            head=head_entity_label,
            relation=relation_label,
            triples_factory=dataset.transductive_training,
            mode="training",
        )
        return head_entity_label, relation_label, predictions.df.iloc[0]["tail_label"]
    return head_entity_label, relation_label, tail_entity_label


# 初始化模型和数据集（只需要执行一次）
model, dataset = load_model("small")

# 示例调用
if __name__ == "__main__":
    head = "Q999726"  # Q10800557
    relation = "P101"  # P131
    tail = "Q223117"  # Q223117
    a, b, c = infer_triples(model, dataset, head, relation, tail)
    print(a, b, c)
