from pykeen.triples import TriplesFactory
from pathlib import Path

HERE = Path("py").resolve()
DATA = HERE.joinpath("data")

"""分不同文件读取训练集、测试集和验证集"""
file_path_train = DATA.joinpath("elec.csv")
file_path_test = DATA.joinpath("elec.csv")
file_path_validate = DATA.joinpath("elec.csv")

training = TriplesFactory.from_path(
    file_path_train,  # 路径
    create_inverse_triples=False,
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
import matplotlib.pyplot as plt
from matplotlib.font_manager import FontProperties

# 指定字体文件路径
font_path = "SimHei.ttf"  # 替换为实际路径
font = FontProperties(fname=font_path)

# 全局设置字体
plt.rcParams["font.sans-serif"] = ["SimHei"]  # 使用 SimHei 字体
plt.rcParams["axes.unicode_minus"] = False  # 解决负号 '-' 显示问题
device = "cuda" if torch.cuda.is_available() else "cpu"
from pykeen.datasets import get_dataset
from pykeen.pipeline import pipeline
from pykeen.trackers import PythonResultTracker

result = pipeline(
    training=training,
    testing=testing,  # 必传
    validation=validation,  # 可不传
    model="RGCN",
    model_kwargs=dict(
        embedding_dim=512,
        num_layers=3,
    ),
    training_kwargs=dict(
        num_epochs=200,
        batch_size=1024,
        use_tqdm_batch=False,
        callbacks="evaluation-loss",
        callbacks_kwargs=dict(triples_factory=validation, prefix="validation"),
    ),
    optimizer_kwargs=dict(
        lr=0.01,  # 学习率可以适当调整
    ),
    evaluation_kwargs=dict(use_tqdm=False),
    random_seed=42,
    device=device,  # 指定使用GPU
    result_tracker=result_tracker,
)


# 使用模型预测
from pykeen import predict

predict.predict_target(
    model=result.model, head="高温环境", relation="运行", triples_factory=testing
).add_membership_columns(testing=testing).df
