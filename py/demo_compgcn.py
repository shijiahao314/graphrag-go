"""The challenge's datasets."""

from pathlib import Path

from pykeen.datasets.inductive.base import DisjointInductivePathDataset
from typing_extensions import Literal

__all__ = [
    "InductiveLPDataset",
    "Size",
]

HERE = Path(__file__).parent.resolve()
DATA = HERE.joinpath("data")

Size = Literal["small", "large"]


class InductiveLPDataset(DisjointInductivePathDataset):
    """An inductive link prediction dataset for the ILPC 2022 Challenge."""

    def __init__(self, size: Size = "small", **kwargs):
        """Initialize the inductive link prediction dataset.

        :param size: "small" or "large"
        :param kwargs: keyword arguments to forward to the base dataset class, cf. DisjointInductivePathDataset
        """
        super().__init__(
            transductive_training_path=DATA.joinpath(size, "train.txt"),
            inductive_inference_path=DATA.joinpath(size, "inference.txt"),
            inductive_validation_path=DATA.joinpath(size, "inference_validation.txt"),
            inductive_testing_path=DATA.joinpath(size, "inference_test.txt"),
            create_inverse_triples=True,
            eager=True,
            **kwargs,
        )


import logging
import os

os.chdir(r"/Users/shijiahao/proj/graphrag-go/py")
from pathlib import Path
import torch
from pykeen.evaluation import RankBasedEvaluator
from pykeen.losses import NSSALoss
from pykeen.models.inductive import InductiveNodePiece, InductiveNodePieceGNN
from pykeen.trackers import ConsoleResultTracker, WANDBResultTracker
from pykeen.training import SLCWATrainingLoop
from pykeen.typing import TESTING, TRAINING, VALIDATION
from pykeen.utils import resolve_device, set_random_seed
from torch.optim import Adam


import datetime

now = datetime.datetime.now()
date_hour_str = now.strftime("%Y-%m-%d %H")

# 设置随机种子以保证结果的可重复性
set_random_seed(42)

# 参数设置
dataset = "small"  # 或者 "large"
embedding_dim = 100
tokens = 5
learning_rate = 0.0001
margin = 15.0
num_negatives = 4
batch_size = 256
epochs = 100
wandb = True  # 如果使用Weights & Biases，则设为True
save = True  # 是否保存模型
# gnn = False  # 使用带有GCN层的Inductive NodePiece模型
log_level = "INFO"

logging.basicConfig(level=log_level)
HERE = Path().resolve()
DATA = HERE.joinpath("data")
print(dataset)
# 加载数据集
dataset = InductiveLPDataset(size=dataset)
loss = NSSALoss(margin=margin)

model_cls = InductiveNodePieceGNN
model = model_cls(
    embedding_dim=embedding_dim,
    triples_factory=dataset.transductive_training,
    inference_factory=dataset.inductive_inference,
    num_tokens=tokens,
    aggregation="mlp",
    loss=loss,
).to(resolve_device())
optimizer = Adam(params=model.parameters(), lr=learning_rate)

if wandb:
    tracker = WANDBResultTracker(project="inductive_ilp", config=dict())
    tracker.start_run(run_name="gnn_compgcn" + date_hour_str)
else:
    tracker = ConsoleResultTracker()

training_loop = SLCWATrainingLoop(
    triples_factory=dataset.transductive_training,
    model=model,
    optimizer=optimizer,
    result_tracker=tracker,
    negative_sampler_kwargs=dict(num_negs_per_pos=num_negatives),
    mode=TRAINING,
)

valid_evaluator = RankBasedEvaluator(
    mode=VALIDATION,
    metrics=["hits_at_k"] * 5,
    metrics_kwargs=[dict(k=k) for k in (1, 3, 5, 10, 100)],
    add_defaults=True,
)
test_evaluator = RankBasedEvaluator(
    mode=TESTING,
    metrics=["hits_at_k"] * 5,
    metrics_kwargs=[dict(k=k) for k in (1, 3, 5, 10, 100)],
    add_defaults=True,
)

training_loop.train(
    triples_factory=dataset.transductive_training,
    num_epochs=epochs,
    batch_size=batch_size,
    callbacks="evaluation",
    callbacks_kwargs=dict(
        evaluator=valid_evaluator,
        evaluation_triples=dataset.inductive_validation.mapped_triples,
        prefix="validation",
        frequency=1,
        additional_filter_triples=dataset.inductive_inference.mapped_triples,
        batch_size=batch_size,
    ),
)

result = test_evaluator.evaluate(
    model=model,
    mapped_triples=dataset.inductive_testing.mapped_triples,
    additional_filter_triples=[
        dataset.inductive_inference.mapped_triples,
        dataset.inductive_validation.mapped_triples,
    ],
    batch_size=batch_size,
)

for metric, metric_label in [
    ("inverse_harmonic_mean_rank", "MRR"),
    *((f"hits_at_{k}", f"Hits@{k}") for k in (100, 10, 5, 3, 1)),
    ("adjusted_arithmetic_mean_rank_index", "AMRI"),
]:
    logging.info(f"Test {metric_label:10}: {result.get_metric(name=metric):.5f}")

torch.save(model, DATA.joinpath("small_model_gnn.pth" + date_hour_str))


from pykeen.evaluation import RankBasedEvaluator, MacroRankBasedEvaluator

test_evaluator = MacroRankBasedEvaluator(
    mode=TESTING,
    metrics=["hits_at_k"] * 5,
    metrics_kwargs=[dict(k=k) for k in (1, 3, 5, 10, 100)],
    add_defaults=True,
)
result = test_evaluator.evaluate(
    model=model,
    mapped_triples=dataset.inductive_testing.mapped_triples,
    additional_filter_triples=[
        dataset.inductive_inference.mapped_triples,
        dataset.inductive_validation.mapped_triples,
    ],
    batch_size=batch_size,
)

for metric, metric_label in [
    ("inverse_harmonic_mean_rank", "MRR"),
    *((f"hits_at_{k}", f"Hits@{k}") for k in (100, 10, 5, 3, 1)),
    ("adjusted_arithmetic_mean_rank_index", "AMRI"),
]:
    logging.info(f"Test {metric_label:10}: {result.get_metric(name=metric):.5f}")

from pykeen.predict import predict_target

# 示例：预测缺失的尾实体
# 假设我们有一个三元组 ("head_entity_label", "relation_label", ?)
head_entity_label = "Q999726"  # 替换为实际的头实体标签
relation_label = "P101"  # 替换为实际的关系标签

# 使用模型进行预测
predictions = predict_target(
    model=model,
    head=head_entity_label,
    relation=relation_label,
    triples_factory=dataset.transductive_training,  # 使用训练集的 triples_factory
    mode=TRAINING,  # 指定模式为训练模式
)

# 将预测结果转换为 DataFrame 并显示前 5 个预测
df_predictions = predictions.df
print("Top 5 predicted tail entities:")
print(df_predictions.head())
# 所有指标
result.to_df()
