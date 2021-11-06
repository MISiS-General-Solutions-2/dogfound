import os

import torch
import torch.nn as nn
from torchvision import models

from config import CFG


def save_model(model, epoch, trainloss, valloss, val_color_acc, val_tail_acc, name):
    """Saves PyTorch model."""
    torch.save(
        {
            "model": model.state_dict(),
            "epoch": epoch,
            "train_loss": trainloss,
            "val_loss": valloss,
            "val_color_accuracy": val_color_acc,
            "val_color_tail": val_tail_acc,
        },
        os.path.join(os.path.join(CFG.LOG_DIR, CFG.OUTPUT_DIR, "weights"), name),
    )


class MultiOutputModel(nn.Module):
    def __init__(self, n_color_classes=3, n_tail_classes=2):
        super().__init__()
        self.resnet = models.resnet34(pretrained=True)
        self.base_model = nn.Sequential(
            *(list(self.resnet.children())[:-1])
        )  # take the model without classifier

        last_channel = (
            models.resnet34().fc.in_features
        )  # size of the layer before the classifier

        # create separate classifiers for our outputs
        self.color = nn.Sequential(
            nn.Dropout(p=0.5),
            nn.Linear(in_features=last_channel, out_features=n_color_classes),
        )
        self.tail = nn.Sequential(
            nn.Dropout(p=0.5),
            nn.Linear(in_features=last_channel, out_features=n_tail_classes),
        )

    def forward(self, x):
        x = self.base_model(x)

        # reshape from [batch, channels, 1, 1] to [batch, channels] to put it into classifier
        x = torch.flatten(x, start_dim=1)
        return {
            "color": self.color(x),
            "tail": self.tail(x),
        }
