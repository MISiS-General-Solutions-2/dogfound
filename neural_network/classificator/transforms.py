import albumentations as A
from albumentations.pytorch import ToTensorV2


# ====================================================
# Transforms
# ====================================================
def get_transforms(*, cfg, data):

    if data == "train":
        return A.Compose(
            [
                A.Resize(cfg.size, cfg.size, p=1.0),
                A.RandomRotate90(),
                A.Flip(),
                A.Transpose(),
                A.OneOf(
                    [
                        A.IAAAdditiveGaussianNoise(),
                        A.GaussNoise(),
                        A.ISONoise(),
                    ],
                    p=0.2,
                ),
                A.OneOf(
                    [
                        A.ImageCompression(),
                        A.JpegCompression(),
                    ],
                    p=0.2,
                ),
                A.OneOf(
                    [
                        A.MotionBlur(p=0.2),
                        A.MedianBlur(blur_limit=3, p=0.1),
                        A.Blur(blur_limit=3, p=0.1),
                    ],
                    p=0.2,
                ),
                A.OneOf(
                    [
                        A.GridDistortion(p=0.1),
                        A.IAAPiecewiseAffine(p=0.3),
                    ],
                    p=0.2,
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
                    mean=cfg.mean,
                    std=cfg.std,
                ),
                ToTensorV2(),
            ]
        )

    elif data == "valid":
        return A.Compose(
            [
                A.Resize(cfg.size, cfg.size),
                A.Normalize(
                    mean=cfg.mean,
                    std=cfg.std,
                ),
                ToTensorV2(),
            ]
        )
