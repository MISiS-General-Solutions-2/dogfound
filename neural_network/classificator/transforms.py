import albumentations as A
from albumentations.pytorch import ToTensorV2

from config import CFG

# ====================================================
# Transforms
# ====================================================

simple_train_transform = A.Compose(
    [
        A.Resize(CFG.size, CFG.size),
        A.HorizontalFlip(),
        A.HueSaturationValue(),
        A.RandomBrightnessContrast(),
        A.Normalize(
            mean=CFG.mean,
            std=CFG.std,
        ),
        ToTensorV2(),
    ]
)

complex_train_transform = A.Compose(
    [
        A.Resize(CFG.size, CFG.size, p=1.0),
        A.OneOf(
            [
                A.RandomRotate90(),
                A.Flip(),
                A.Transpose(),
            ],
            p=0.3,
        ),
        A.OneOf(
            [
                A.IAAAdditiveGaussianNoise(),
                A.GaussNoise(),
                A.ISONoise(),
            ],
            p=0.6,
        ),
        A.OneOf(
            [
                A.ImageCompression(),
                A.JpegCompression(),
            ],
            p=0.7,
        ),
        A.OneOf(
            [
                A.MotionBlur(p=0.2),
                A.MedianBlur(blur_limit=3, p=0.1),
                A.Blur(blur_limit=3, p=0.1),
            ],
            p=0.4,
        ),
        A.OneOf(
            [
                A.GridDistortion(p=0.1),
                A.IAAPiecewiseAffine(p=0.3),
            ],
            p=0.4,
        ),
        A.OneOf(
            [
                A.CLAHE(clip_limit=2),
                A.IAASharpen(),
                A.IAAEmboss(),
                A.RandomBrightnessContrast(),
            ],
            p=0.3,
        ),
        A.HueSaturationValue(p=0.3),
        A.Normalize(
            mean=CFG.mean,
            std=CFG.std,
        ),
        ToTensorV2(),
    ]
)


valid_transform = A.Compose(
    [
        A.Resize(CFG.size, CFG.size),
        A.Normalize(
            mean=CFG.mean,
            std=CFG.std,
        ),
        ToTensorV2(),
    ]
)
