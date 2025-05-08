import torch
import numpy as np
import matplotlib.pyplot as plt
from pykeen.pipeline import pipeline
from pykeen.datasets import Nations
from sklearn.decomposition import PCA
from pykeen.trackers import PythonResultTracker

result_tracker = PythonResultTracker()
import pandas
import seaborn

from pykeen.datasets import get_dataset
from pykeen.pipeline import pipeline
from pykeen.trackers import PythonResultTracker

# 检查GPU可用性
device = "cuda:0" if torch.cuda.is_available() else "cpu"
print(f"Using device: {device}")


# 1. 加载数据集
dataset = Nations()
print("训练三元组数量:", dataset.training.num_triples)
print("实体数量:", dataset.training.num_entities)
print("关系数量:", dataset.training.num_relations)

# 2. 创建并训练RGCN模型（使用GPU）
result = pipeline(
    model="RGCN",
    dataset=dataset,
    model_kwargs=dict(
        embedding_dim=512,
        num_layers=3,
    ),
    training_kwargs=dict(
        num_epochs=200,
        batch_size=1024,
        use_tqdm_batch=False,
        callbacks="evaluation-loss",
        callbacks_kwargs=dict(triples_factory=dataset.validation, prefix="validation"),
    ),
    optimizer_kwargs=dict(
        lr=0.01,  # 学习率可以适当调整
    ),
    evaluation_kwargs=dict(use_tqdm=False),
    random_seed=42,
    device=device,  # 指定使用GPU
    result_tracker=result_tracker,
)

# 3. 保存模型
result.save_to_directory("rgcn_nations_model")

# 4. 评估结果
metric_results = result.metric_results.to_df()
print("\n评估结果:")
print(metric_results)

# 5. 可视化训练损失（注意将数据移回CPU）
result.plot(er_kwargs=dict(plot_relations=True))

plt.show()

grid = seaborn.relplot(
    data=pandas.DataFrame(
        data=[
            [step, step_metrics.get("loss"), step_metrics.get("validation.loss")]
            for step, step_metrics in result_tracker.metrics.items()
        ],
        columns=["step", "training", "validation"],
    )
    .set_index("step")
    .rolling(window=5)
    .agg(["min", "mean", "max"]),
    kind="line",
)
grid.fig.show()

# 6. 实体嵌入可视化（将嵌入数据移回CPU）
entity_embeddings = result.model.entity_representations[0]().detach().cpu().numpy()
pca = PCA(n_components=2)
embeddings_2d = pca.fit_transform(entity_embeddings)

plt.figure(figsize=(12, 8))
plt.scatter(embeddings_2d[:, 0], embeddings_2d[:, 1], alpha=0.7)

# 添加标签（只显示部分避免拥挤）
entity_labels = list(dataset.training.entity_to_id.keys())
for i, entity in enumerate(entity_labels):
    if i % 3 == 0:  # 只标注1/3的实体
        plt.annotate(entity, (embeddings_2d[i, 0], embeddings_2d[i, 1]))

plt.title("Entity Embeddings (PCA)")
plt.xlabel("PCA Component 1")
plt.ylabel("PCA Component 2")
plt.grid()
plt.savefig("entity_embeddings.png")
plt.show()

# 7. 知识图谱补全（推理示例，进行预测）
model = result.model
triples_factory = result.training


def predict_tails(head, relation, k=5):
    """预测最有可能的k个尾实体"""
    head_id = triples_factory.entity_to_id[head]
    relation_id = triples_factory.relation_to_id[relation]

    # 在GPU上准备所有可能的尾实体
    tail_ids = torch.arange(triples_factory.num_entities, device=model.device)
    head_ids = torch.full_like(tail_ids, fill_value=head_id)
    relation_ids = torch.full_like(tail_ids, fill_value=relation_id)
    triples = torch.stack([head_ids, relation_ids, tail_ids], dim=1)

    # 批量预测（保持在GPU上）
    scores = model.predict(triples).cpu().numpy()

    # 获取排序结果
    ranked_indices = np.argsort(-scores)
    print(f"\n预测 '{head} - {relation} - ?' 的Top-{k}结果:")
    for i in range(k):
        tail_id = ranked_indices[i]
        tail = triples_factory.entity_id_to_label[tail_id]
        print(f"{i + 1}. {tail} (分数: {scores[tail_id]:.4f})")


# 8. 关系嵌入可视化（将数据移回CPU）
relation_embeddings = result.model.relation_representations[0]().detach().cpu().numpy()
pca = PCA(n_components=2)
rel_embeddings_2d = pca.fit_transform(relation_embeddings)

plt.figure(figsize=(10, 8))
plt.scatter(rel_embeddings_2d[:, 0], rel_embeddings_2d[:, 1], alpha=0.7)

# 添加关系标签
for i, relation in enumerate(dataset.training.relation_to_id.keys()):
    plt.annotate(relation, (rel_embeddings_2d[i, 0], rel_embeddings_2d[i, 1]))

plt.title("Relation Embeddings (PCA)")
plt.xlabel("PCA Component 1")
plt.ylabel("PCA Component 2")
plt.grid()
plt.savefig("relation_embeddings.png")
plt.show()

# 9. 知识图谱可视化（使用networkx）
try:
    import networkx as nx

    # 创建子图可视化
    G = nx.DiGraph()

    # 添加部分三元组（避免可视化太密集）
    sample_triples = dataset.training.triples[:20]

    for h, r, t in sample_triples:
        # h_label = dataset.training.entity_id_to_label[h]
        # r_label = dataset.training.relation_id_to_label[r]
        # t_label = dataset.training.entity_id_to_label[t]
        G.add_edge(h, t, relation=r)

    # 绘制网络图
    plt.figure(figsize=(15, 12))
    pos = nx.spring_layout(G)
    nx.draw(G, pos, with_labels=True, node_size=2000, node_color="skyblue")
    edge_labels = nx.get_edge_attributes(G, "label")
    for (u, v), label in edge_labels.items():
        x = (pos[u][0] + pos[v][0]) / 2
        y = (pos[u][1] + pos[v][1]) / 2
        plt.text(x, y, label, bbox=dict(facecolor="white", alpha=0.5))
    plt.title("Knowledge Graph Subset")
    plt.savefig("kg_subgraph.png")
    plt.show()
except ImportError:
    print("\n要使用网络图可视化，请先安装networkx: pip install networkx")

result.plot_er(
    relations={"negativebehavior"},
    apply_limits=False,
    plot_relations=True,
)

# 所有指标
result.metric_results.to_df()


# 给定三元组的头实体与关系使用模型预测
from pykeen import predict

predict.predict_target(
    model=model, head="brazil", relation="intergovorgs", triples_factory=tf
).add_membership_columns(testing=dataset.testing).df

# 自动过滤掉非归纳的预测

predict.predict_target(
    model=model, head="brazil", relation="intergovorgs", triples_factory=tf
).filter_triples(dataset.testing).df

# 所有的预测
predict.predict_all(model=model).process(
    factory=result.training
).add_membership_columns(**dataset.factory_dict).df
