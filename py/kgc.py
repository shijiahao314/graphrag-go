from pathlib import Path
import torch
import random
import numpy as np
from pykeen.triples import TriplesFactory
from pykeen.models import RGCN
from pykeen.utils import resolve_device, set_random_seed
from pykeen.predict import predict_target


# 设置随机种子，确保各模块结果一致
SEED = 42
set_random_seed(SEED)
random.seed(SEED)
np.random.seed(SEED)
torch.manual_seed(SEED)
if torch.cuda.is_available():
    torch.cuda.manual_seed_all(SEED)
    # 设置 cudnn 为确定性模式，可能会影响性能
    torch.backends.cudnn.deterministic = True
    torch.backends.cudnn.benchmark = False

HERE = Path("py").resolve()
DATA = HERE.joinpath("data")

# 分别读取训练集、测试集和验证集
file_path_train = DATA.joinpath("elec.csv")
file_path_test = DATA.joinpath("elec.csv")
file_path_validate = DATA.joinpath("elec.csv")

training = TriplesFactory.from_path(
    file_path_train,
    create_inverse_triples=False,
)

# 统一使用训练集的映射关系加载测试集和验证集
testing = TriplesFactory.from_path(
    file_path_test,
    entity_to_id=training.entity_to_id,
    relation_to_id=training.relation_to_id,
    create_inverse_triples=False,
)

validation = TriplesFactory.from_path(
    file_path_validate,
    entity_to_id=training.entity_to_id,
    relation_to_id=training.relation_to_id,
    create_inverse_triples=False,
)

device = resolve_device()

# 创建RGCN模型并加载训练时的权重
model = RGCN(
    triples_factory=training,
    embedding_dim=512,
    num_layers=3,
).to(device)

model_path = DATA.joinpath("synthesized_rgcn.pt")
model.load_state_dict(torch.load(model_path, map_location=device))


# 简化推理函数，直接调用 predict_target 完成推理
def infer_triples(
    model, head_entity_label: str, relation_label: str, tail_entity_label: str
):
    predictions = predict_target(
        model=model,
        head=head_entity_label or None,
        relation=relation_label or None,
        tail=tail_entity_label or None,
        triples_factory=training,
    ).df

    # 若输入为空，则使用预测结果中的第一个结果作为默认值
    if not head_entity_label:
        head_entity_label = predictions.iloc[0]["head_label"]
    if not relation_label:
        relation_label = predictions.iloc[0]["relation_label"]
    if not tail_entity_label:
        tail_entity_label = predictions.iloc[0]["tail_label"]

    return head_entity_label, relation_label, tail_entity_label


# 通用评估函数，计算 Hits@1 和 MRR
def evaluate(model, testing, target_type):
    correct = 0
    total_mrr = 0
    total = 0
    for head_id, relation_id, tail_id in testing.mapped_triples:
        head_label = testing.entity_id_to_label[head_id.item()]
        relation_label = testing.relation_id_to_label[relation_id.item()]
        tail_label = testing.entity_id_to_label[tail_id.item()]

        if target_type == "head":
            predictions = predict_target(
                model=model,
                relation=relation_label,
                tail=tail_label,
                triples_factory=testing,
            ).df
            target_label = head_label
            target_column = "head_label"
        elif target_type == "relation":
            predictions = predict_target(
                model=model,
                head=head_label,
                tail=tail_label,
                triples_factory=testing,
            ).df
            target_label = relation_label
            target_column = "relation_label"
        elif target_type == "tail":
            predictions = predict_target(
                model=model,
                head=head_label,
                relation=relation_label,
                triples_factory=testing,
            ).df
            target_label = tail_label
            target_column = "tail_label"

        # 使用argmax()来定位排名
        rank = (predictions[target_column] == target_label).argmax() + 1

        # 如果预测的最上面结果和目标标签匹配，则为正确
        top_predicted = predictions.iloc[0][target_column]
        if top_predicted == target_label:
            correct += 1

        total_mrr += 1 / rank
        total += 1

    # 计算Hits@1和MRR
    hits_at_1 = correct / total if total > 0 else 0
    mrr = total_mrr / total if total > 0 else 0
    return hits_at_1, mrr


# 返回尾实体预测的 Hits@1 和 MRR
def benchmark():
    # return Hits@1, MRR
    tail_hits_at_1, tail_mrr = evaluate(model, testing, "tail")
    return tail_hits_at_1, tail_mrr


# 示例调用
if __name__ == "__main__":
    head = "输电线路"
    relation = "连接"
    tail = ""  # 若尾实体为空，则预测结果中第一个结果作为默认值
    a, b, c = infer_triples(model, head, relation, tail)
    print(f"推理结果: {a} → {b} → {c}")

    head_hits_at_1, head_mrr = evaluate(model, testing, "head")
    relation_hits_at_1, relation_mrr = evaluate(model, testing, "relation")
    tail_hits_at_1, tail_mrr = evaluate(model, testing, "tail")

    print(f"头实体 Hits@1: {head_hits_at_1:.3f}")
    print(f"关系 Hits@1: {relation_hits_at_1:.3f}")
    print(f"尾实体 Hits@1: {tail_hits_at_1:.3f}")
    print(f"头实体 MRR: {head_mrr:.3f}")
    print(f"关系 MRR: {relation_mrr:.3f}")
    print(f"尾实体 MRR: {tail_mrr:.3f}")
