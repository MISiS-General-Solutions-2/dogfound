import argparse
import ast

import torch
from PIL import Image
from torchvision import models, transforms

file = "./imagenet_classes.txt"


def _read_dict(file):
    with open(file, "r") as data:
        dictionary = ast.literal_eval(data.read())
    return dictionary


def main():
    parser = argparse.ArgumentParser(description="Path to the image")
    parser.add_argument(
        "--path",
        type=str,
        help="Path to the image",
    )
    args = parser.parse_args()
    path_image = args.path
    # Load model
    resnet = models.resnet34(pretrained=True)
    resnet.eval()
    # Load image
    img = Image.open(path_image)
    # Transform image
    transform = transforms.Compose(
        [
            transforms.Resize(254),
            transforms.ToTensor(),
            transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
        ]
    )
    img_t = transform(img)
    batch_t = torch.unsqueeze(img_t, 0)
    # Inference
    out = resnet(batch_t)
    print(out.shape)

    label_dict = _read_dict(file)

    _, index = torch.max(out, 1)
    percentage = torch.nn.functional.softmax(out, dim=1)[0] * 100

    print(label_dict[index[0].item()], percentage[index[0].item()].item())

    _, indices = torch.sort(out, descending=True)
    print(
        [
            (label_dict[idx.item()], percentage[idx.item()].item())
            for idx in indices[0][:5]
        ]
    )


if __name__ == "__main__":
    main()
