import time

import torch
from config import CFG
from tqdm import tqdm
from utils.utils import AverageMeter, get_score, timeSince


def train_fn(train_loader, model, criterion, optimizer, epoch, device, scheduler=None):
    batch_time = AverageMeter()
    data_time = AverageMeter()
    losses = AverageMeter()
    accuracy = AverageMeter()
    # switch to train mode
    model.train()
    start = end = time.time()

    # Iterate over dataloader
    for i, (images, labels) in enumerate(tqdm(train_loader)):
        # measure data loading time
        data_time.update(time.time() - end)
        # zero the gradients
        optimizer.zero_grad()

        images = images.to(device)
        labels = labels.to(device)

        batch_size = labels.size(0)

        y_preds = model(images)
        loss = criterion(y_preds, labels)
        # Compute gradients and do step
        loss.backward()
        optimizer.step()

        if scheduler is not None:
            scheduler.step()

        # record loss
        losses.update(loss.item(), batch_size)
        classes = y_preds.argmax(dim=1)
        acc = torch.mean((classes == labels).float())
        accuracy.update(acc, batch_size)

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
    return losses.avg, accuracy.avg


def valid_fn(valid_loader, model, criterion, device):
    batch_time = AverageMeter()
    data_time = AverageMeter()
    losses = AverageMeter()
    val_accuracy = AverageMeter()
    val_f1 = AverageMeter()

    # switch to evaluation mode
    model.eval()
    start = end = time.time()
    for step, (images, labels) in enumerate(valid_loader):
        # measure data loading time
        data_time.update(time.time() - end)
        images = images.to(device)
        labels = labels.to(device)
        batch_size = labels.size(0)

        # compute loss
        with torch.no_grad():
            y_preds = model(images)
        loss = criterion(y_preds, labels)
        losses.update(loss.item(), batch_size)

        # record accuracy
        preds = y_preds.softmax(1).cpu().data.numpy()
        labels = labels.cpu().data.numpy()

        # scoring on validation set
        val_acc_score = get_score(labels, preds.argmax(1), metric="accuracy")
        val_f1_score = get_score(labels, preds.argmax(1), metric="f1_score")

        val_accuracy.update(val_acc_score, batch_size)
        val_f1.update(val_f1_score, batch_size)

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
    return losses.avg, val_accuracy.avg, val_f1.avg
