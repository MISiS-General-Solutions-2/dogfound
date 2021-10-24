import os

import cv2
import torch
from torch.utils.data import Dataset

from config import CFG


# ====================================================
# Dataset creation
# ====================================================
class TrainValDataset(Dataset):
    def __init__(self, df, transform=None):
        self.df = df
        self.file_names = df["image"].values
        self.colors = df["color"].values
        self.tails = df["tail"].values
        self.transform = transform

    def __len__(self):
        return len(self.df)

    def __getitem__(self, idx):
        file_name = self.file_names[idx]
        head, tail = os.path.split(file_name)
        file_path = f"{CFG.DATA_PATH}/{head}/crops/dog/{tail}"
        image = cv2.imread(file_path)
        image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
        if self.transform:
            augmented = self.transform(image=image)
            image = augmented["image"]
        color = torch.tensor(self.colors[idx]).long()
        tail = torch.tensor(self.tails[idx]).long()
        # return the image and all the associated labels
        dict_data = {
            "img": image,
            "labels": {
                "color": color,
                "tail": tail,
            },
        }
        return dict_data
