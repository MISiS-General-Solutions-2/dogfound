import argparse
import os
import time

import torch
import torch.nn as nn
from config import CFG
from model import get_model, save_model
from torch.optim import Adam
from torch.utils.tensorboard import SummaryWriter
from torch_lr_finder import LRFinder
from train_val import train_fn, valid_fn
from utils.data_processing import train_val_dataloaders
from utils.utils import init_logger, save_batch, seed_torch


def main():
    parser = argparse.ArgumentParser(
        description="Define whether to save train batch figs or find optimal learning rate"
    )
    parser.add_argument(
        "--save_batch_fig",
        action="store_true",
        help="Whether to save a sample of a batch or not",
    )
    parser.add_argument(
        "--find_lr",
        action="store_true",
        help="Whether to find optimal learning rate or not",
    )

    args = parser.parse_args()
    save_single_batch = args.save_batch_fig
    find_lr = args.find_lr

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

    # Classes
    CLASS_NAMES = [
        "colorful_long",
        "colorful_short",
        "dark_long",
        "dark_short",
        "light_long",
        "light_short",
    ]

    # ====================================================
    # Form dataloaders
    # ====================================================

    device = torch.device(f"cuda:{CFG.GPU_ID}")

    dataloaders = train_val_dataloaders(CFG)

    train_loader = dataloaders["train"]
    val_loader = dataloaders["val"]

    # Show batch to see the effect of augmentations
    if save_single_batch:
        LOGGER.info("Creating dir to save samples of a batch...")
        path_to_figs = os.path.join(logger_path, "batch_figs")
        os.makedirs(path_to_figs)
        LOGGER.info("Saving figures of a single batch...")
        save_batch(train_loader, CLASS_NAMES, path_to_figs, CFG)
        LOGGER.info("Figures have been saved!")

    # ====================================================
    # model & optimizer
    # ====================================================
    model = get_model(CFG)
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

    """
    scheduler = torch.optim.lr_scheduler.CyclicLR(
        optimizer,
        base_lr=CFG.min_lr,
        max_lr=CFG.lr,
        mode="triangular2",
        step_size_up=2652,
    )
    """

    scheduler = None

    # ====================================================
    # loop
    # ====================================================
    criterion = nn.CrossEntropyLoss()
    LOGGER.info("Select CrossEntropyLoss criterion")

    if find_lr:
        print("Finding oprimal learning rate...")
        # Add this line before running `LRFinder`
        lr_finder = LRFinder(model, optimizer, criterion, device="cuda")
        lr_finder.range_test(train_loader, end_lr=100, num_iter=100)
        lr_finder.plot()  # to inspect the loss-learning rate graph
        lr_finder.reset()  # to reset the model and optimizer to their initial state
        print("Optimal learning rate has been found!")

    best_epoch = 0
    best_acc_score = 0.0
    best_f1_score = 0.0

    count_bad_epochs = 0  # Count epochs that don't improve the score

    for epoch in range(CFG.epochs):
        start_time = time.time()

        if epoch >= 50:
            for name, child in model.named_children():
                for param in child.parameters():
                    param.requires_grad = True

        # train
        avg_train_loss, train_acc = train_fn(
            train_loader,
            model,
            criterion,
            optimizer,
            epoch,
            device,
            scheduler=scheduler,
        )

        # eval
        avg_val_loss, val_acc_score, val_f1_score = valid_fn(
            val_loader, model, criterion, device
        )

        cur_lr = optimizer.param_groups[0]["lr"]
        LOGGER.info(f"Current learning rate: {cur_lr}")

        tb.add_scalar("Learning rate", cur_lr, epoch + 1)
        tb.add_scalar("Train Loss", avg_train_loss, epoch + 1)
        tb.add_scalar("Train accuracy", train_acc, epoch + 1)
        tb.add_scalar("Val Loss", avg_val_loss, epoch + 1)
        tb.add_scalar("Val accuracy score", val_acc_score, epoch + 1)
        tb.add_scalar("Val f1 score", val_f1_score, epoch + 1)

        elapsed = time.time() - start_time

        LOGGER.info(
            f"Epoch {epoch+1} - avg_train_loss: {avg_train_loss:.4f} \
            avg_val_loss: {avg_val_loss:.4f}  time: {elapsed:.0f}s"
        )
        LOGGER.info(
            f"Epoch {epoch+1} - Accuracy: {val_acc_score} - F1-score {val_f1_score}"
        )

        best_acc_bool = False
        best_f1_bool = False

        # Update best score
        if val_acc_score >= best_acc_score:
            best_acc_score = val_acc_score
            best_acc_bool = True

        if val_f1_score >= best_f1_score:
            best_f1_score = val_f1_score
            best_f1_bool = True

        if best_acc_bool and best_f1_bool:
            LOGGER.info(
                f"Epoch {epoch+1} - Save Best Accuracy: {best_acc_score:.4f} - \
                Save Best F1-score: {best_f1_score:.4f} Model"
            )
            save_model(
                model, epoch + 1, avg_train_loss, avg_val_loss, val_f1_score, "best.pt"
            )
            best_epoch = epoch + 1
            count_bad_epochs = 0
        else:
            count_bad_epochs += 1
        print(count_bad_epochs)
        LOGGER.info(f"Number of bad epochs {count_bad_epochs}")
        # Early stopping
        if count_bad_epochs > CFG.early_stopping:
            LOGGER.info(
                f"Stop the training, since the score has not improved for {CFG.early_stopping} epochs!"
            )
            save_model(
                model,
                epoch + 1,
                avg_train_loss,
                avg_val_loss,
                val_f1_score,
                f"{CFG.model_name}_epoch{epoch+1}_last.pth",
            )
            break
        elif epoch + 1 == CFG.epochs:
            LOGGER.info(f"Reached the final {epoch+1} epoch!")
            save_model(
                model,
                epoch + 1,
                avg_train_loss,
                avg_val_loss,
                val_f1_score,
                f"{CFG.model_name}_epoch{epoch+1}_final.pth",
            )

    LOGGER.info(
        f"AFTER TRAINING: Epoch {best_epoch}: Best Accuracy: {best_acc_score:.4f} - \
                Best F1-score: {best_f1_score:.4f}"
    )
    tb.close()


if __name__ == "__main__":
    main()
