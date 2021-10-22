import albumentations as A
import cv2
import torchvision
from albumentations.pytorch.transforms import ToTensorV2
from torch.utils.data import DataLoader
from torchvision.datasets.folder import IMG_EXTENSIONS


def _albumentations_loader(path):
    """Load an image."""
    image = cv2.imread(path)
    image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    return image


class AlbumentationsImageFolder(torchvision.datasets.DatasetFolder):
    """Helper class to apply augmentations and form the dataset."""

    def __init__(
        self,
        root,
        transform=None,
        target_transform=None,
        loader=_albumentations_loader,
        is_valid_file=None,
    ):
        super(AlbumentationsImageFolder, self).__init__(
            root,
            loader,
            IMG_EXTENSIONS if is_valid_file is None else None,
            transform=transform,
            target_transform=target_transform,
            is_valid_file=is_valid_file,
        )
        self.imgs = self.samples

    def __getitem__(self, index):
        """
        Args:
            index (int): Index
        Returns:
            tuple: (sample, target) where target is class_index of the target class.
        """

        path, target = self.samples[index]
        sample = self.loader(path)
        if self.transform is not None:
            augmented = self.transform(image=sample)
            sample = augmented["image"]
        if self.target_transform is not None:
            target = self.target_transform(target)
        return sample, target


def train_val_dataloaders(cfg):
    """Form the dataloaders for training and validation and store them in the dictionary.
    :param train_path: path to images for trainin
    :param val_path: path to images for validation
    :param batch_size: size of the batch
    :return: the dictionary with dataloaders
    """
    train_transform = A.Compose(
        [
            A.Resize(cfg.size, cfg.size),
            A.Normalize(),
            ToTensorV2(),
        ]
    )

    val_transforms = A.Compose(
        [A.Resize(cfg.size, cfg.size), A.Normalize(), A.pytorch.transforms.ToTensorV2()]
    )

    train_dataset = AlbumentationsImageFolder(cfg.TRAIN_PATH, train_transform)
    val_dataset = AlbumentationsImageFolder(cfg.VAL_PATH, val_transforms)

    dataloader = dict()

    dataloader["train"] = DataLoader(
        dataset=train_dataset,
        batch_size=cfg.batch_size,
        shuffle=True,
        num_workers=4,
        drop_last=True,
    )

    dataloader["val"] = DataLoader(
        dataset=val_dataset,
        batch_size=cfg.batch_size,
        shuffle=True,
        num_workers=4,
        drop_last=True,
    )

    return dataloader
