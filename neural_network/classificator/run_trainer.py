import os
import random
import time

import numpy as np
import pandas as pd
import torch
from check_paths import check_paths
from config import CFG
from dataset import TrainValDataset
from model import MultiOutputModel, save_model
from sklearn.model_selection import train_test_split
from torch.optim import Adam
from torch.utils.data import DataLoader
from torch.utils.tensorboard import SummaryWriter
from train_val import train_fn, valid_fn
from transforms import get_transforms
from utils.utils import init_logger, seed_torch


def main():
    # Path to log
    logger_path = os.path.join(CFG.LOG_DIR, CFG.OUTPUT_DIR)

    # Create dir for saving logs and weights
    print(f"Creating dir {CFG.OUTPUT_DIR} for saving logs")
    os.makedirs(os.path.join(logger_path, "weights"))
    print(f"Dir {CFG.OUTPUT_DIR} has been created!")

    # Define logger to save train logs
    LOGGER = init_logger(os.path.join(logger_path, "train.log"))
    # Write to tensorboard
    tb = SummaryWriter(logger_path)

    # Set seed
    seed_torch(seed=CFG.seed)

    LOGGER.info("Reading data...")
    train_df = pd.read_csv(CFG.PATH_CSV)
    mask_train = train_df["image"].apply(check_paths)

    train_df = train_df[mask_train]

    train_df["color"] = train_df["color"].apply(lambda x: x - 1)
    train_df["tail"] = train_df["tail"].apply(lambda x: x - 1)

    print(train_df.shape)

    train_df = train_df[["image", "color", "tail"]]

    print(train_df.shape)

    owner_df = pd.read_csv(CFG.PATH_OWNER)
    mask_owner = owner_df["image"].apply(check_paths)

    owner_df = owner_df[mask_owner]

    owner_df["color"] = owner_df["color"].apply(lambda x: x - 1)
    owner_df["tail"] = owner_df["tail"].apply(lambda x: x - 1)

    print(owner_df.shape)

    final_df = pd.concat([train_df, owner_df], axis=0)

    print(final_df.shape)

    print(final_df.head())

    # Classes
    COLORS = [
        "dark",
        "light",
        "colorful",
    ]

    TAIL = ["short", "long"]

    LOGGER.info("Splitting data for training and validation...")
    X_train, X_val, y_train, y_val = train_test_split(
        final_df["image"],
        final_df[["color", "tail"]],
        test_size=0.2,
        random_state=CFG.seed,
        shuffle=True,
    )

    train_fold = pd.concat([X_train, y_train], axis=1)
    LOGGER.info("train shape: ")
    LOGGER.info(train_fold.shape)
    valid_fold = pd.concat([X_val, y_val], axis=1)
    LOGGER.info("valid shape: ")
    LOGGER.info(valid_fold.shape)

    # ====================================================
    # Form dataloaders
    # ====================================================

    device = torch.device(f"cuda:{CFG.GPU_ID}")

    train_dataset = TrainValDataset(
        train_fold, transform=get_transforms(cfg=CFG, data="train")
    )
    valid_dataset = TrainValDataset(
        valid_fold, transform=get_transforms(cfg=CFG, data="valid")
    )

    def seed_worker(worker_id):
        worker_seed = torch.initial_seed() % 2 ** 32
        np.random.seed(worker_seed)
        random.seed(worker_seed)

    g = torch.Generator()
    g.manual_seed(0)

    train_loader = DataLoader(
        train_dataset,
        batch_size=CFG.batch_size,
        shuffle=True,
        num_workers=CFG.num_workers,
        worker_init_fn=seed_worker,
        generator=g,
        drop_last=True,
    )
    val_loader = DataLoader(
        valid_dataset,
        batch_size=CFG.batch_size,
        shuffle=False,
        num_workers=CFG.num_workers,
        worker_init_fn=seed_worker,
        generator=g,
        drop_last=False,
    )

    # ====================================================
    # model & optimizer
    # ====================================================
    model = MultiOutputModel()
    model.to(device)

    if CFG.freeze:
        for name, child in model.named_children():
            for param in child.parameters():
                param.requires_grad = False
            if name == "fc":
                for param in child.parameters():
                    param.requires_grad = True

    LOGGER.info(f"Batch size {CFG.batch_size}")
    LOGGER.info(f"Input size {CFG.size}")

    optimizer = Adam(model.parameters(), lr=CFG.lr, weight_decay=CFG.weight_decay)

    # ====================================================
    # loop
    best_acc_score = 0.0

    for epoch in range(CFG.epochs):
        start_time = time.time()

        # train
        train_dict = train_fn(
            train_loader,
            model,
            optimizer,
            epoch,
            device,
        )

        avg_train_loss = train_dict["train_loss"]
        train_color_acc = train_dict["train_color_acc"]
        train_tail_acc = train_dict["train_tail_acc"]

        # eval
        val_dict = valid_fn(val_loader, model, device)

        avg_val_loss = val_dict["val_loss"]
        val_color_acc = val_dict["val_color_acc"]
        val_tail_acc = val_dict["val_tail_acc"]

        cur_lr = optimizer.param_groups[0]["lr"]
        LOGGER.info(f"Current learning rate: {cur_lr}")

        tb.add_scalar("Learning rate", cur_lr, epoch + 1)
        tb.add_scalar("Train Loss", avg_train_loss, epoch + 1)
        tb.add_scalar("Train color accuracy", train_color_acc, epoch + 1)
        tb.add_scalar("Train tail accuracy", train_tail_acc, epoch + 1)
        tb.add_scalar("Val Loss", avg_val_loss, epoch + 1)
        tb.add_scalar("Val color accuracy", val_color_acc, epoch + 1)
        tb.add_scalar("Val tail accuracy", val_tail_acc, epoch + 1)

        elapsed = time.time() - start_time

        LOGGER.info(
            f"Epoch {epoch+1} - avg_train_loss: {avg_train_loss:.4f} \
            avg_val_loss: {avg_val_loss:.4f}  time: {elapsed:.0f}s"
        )
        LOGGER.info(
            f"Epoch {epoch+1} - Color Accuracy: {val_color_acc} - Tail Accuracy {val_tail_acc}"
        )

        mean_acc = (val_color_acc + val_tail_acc) / 2
        save_model(
            model,
            epoch + 1,
            avg_train_loss,
            avg_val_loss,
            val_color_acc,
            val_tail_acc,
            f"{epoch}.pt",
        )
        # Update best score
        if mean_acc >= best_acc_score:
            LOGGER.info(f"Epoch {epoch+1} - Best Accuracy: {mean_acc:.4f}")
    tb.close()


if __name__ == "__main__":
    main()
