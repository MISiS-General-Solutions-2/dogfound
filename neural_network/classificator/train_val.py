import time

import torch
from tqdm import tqdm

from config import CFG
from loss import get_loss
from utils.utils import AverageMeter, get_score, timeSince


def train_fn(train_loader, model, optimizer, epoch, device):
    batch_time = AverageMeter()
    data_time = AverageMeter()
    losses = AverageMeter()
    train_accuracy_color = AverageMeter()
    train_accuracy_tail = AverageMeter()
    # switch to train mode
    model.train()
    start = end = time.time()

    # Iterate over dataloader
    for i, train_batch in enumerate(tqdm(train_loader)):
        # measure data loading time
        data_time.update(time.time() - end)
        # zero the gradients
        optimizer.zero_grad()

        images = train_batch["img"]
        images = images.to(device)

        labels = train_batch["labels"]

        batch_size = images.size(0)

        output = model(images)

        loss, _ = get_loss(output, labels, device)

        # Compute gradients and do step
        loss.backward()
        optimizer.step()

        # record loss
        losses.update(loss.item(), batch_size)

        # Train accuracy for color
        color_class = output["color"].argmax(dim=1)
        acc_color = torch.mean((color_class == labels["color"]).float())
        train_accuracy_color.update(acc_color, batch_size)

        # Train accuracy for tail
        tail_class = output["tail"].argmax(dim=1)
        acc_tail = torch.mean((tail_class == labels["tail"]).float())
        train_accuracy_tail.update(acc_tail, batch_size)

        # measure elapsed time
        batch_time.update(time.time() - end)
        end = time.time()
        if i % CFG.print_freq == 0 or i == (len(train_loader) - 1):
            print(
                "Epoch: [{Epoch:d}][{Iter:d}/{Len:d}] "
                "Data {data_time.val:.3f} ({data_time.avg:.3f}) "
                "Elapsed {remain:s} "
                "Loss: {loss.val:.4f}({loss.avg:.4f}) ".format(
                    Epoch=epoch + 1,
                    Iter=i,
                    Len=len(train_loader),
                    data_time=data_time,
                    loss=losses,
                    remain=timeSince(start, float(i + 1) / len(train_loader)),
                )
            )
    return {
        "train_loss": losses.avg,
        "train_color_acc": train_accuracy_color.avg,
        "train_tail_acc": train_accuracy_tail.avg,
    }


def valid_fn(valid_loader, model, device):
    batch_time = AverageMeter()
    data_time = AverageMeter()

    losses = AverageMeter()
    val_accuracy_color = AverageMeter()
    val_accuracy_tail = AverageMeter()

    # switch to evaluation mode
    model.eval()
    start = end = time.time()
    for step, val_batch in enumerate(valid_loader):
        # measure data loading time
        data_time.update(time.time() - end)
        images = val_batch["img"]
        images = images.to(device)
        batch_size = images.size(0)

        labels = val_batch["labels"]

        # compute loss
        with torch.no_grad():
            output = model(images)
        loss, _ = get_loss(output, labels, device)
        losses.update(loss.item(), batch_size)

        # record accuracy
        preds_color = output["color"].softmax(1).cpu().data.numpy()
        preds_tail = output["tail"].softmax(1).cpu().data.numpy()

        # scoring on validation set
        val_acc_color = get_score(
            labels["color"].cpu().data.numpy(), preds_color.argmax(1), metric="accuracy"
        )
        val_acc_tail = get_score(
            labels["tail"].cpu().data.numpy(), preds_tail.argmax(1), metric="accuracy"
        )

        val_accuracy_color.update(val_acc_color, batch_size)
        val_accuracy_tail.update(val_acc_tail, batch_size)

        # measure elapsed time
        batch_time.update(time.time() - end)
        end = time.time()
        if step % CFG.print_freq == 0 or step == (len(valid_loader) - 1):
            print(
                "EVAL: [{Step:d}/{Len:d}] "
                "Data {data_time.val:.3f} ({data_time.avg:.3f}) "
                "Elapsed {remain:s} "
                "Loss: {loss.val:.4f}({loss.avg:.4f}) ".format(
                    Step=step,
                    Len=len(valid_loader),
                    data_time=data_time,
                    loss=losses,
                    remain=timeSince(start, float(step + 1) / len(valid_loader)),
                )
            )
    return {
        "val_loss": losses.avg,
        "val_color_acc": val_accuracy_color.avg,
        "val_tail_acc": val_accuracy_tail.avg,
    }
