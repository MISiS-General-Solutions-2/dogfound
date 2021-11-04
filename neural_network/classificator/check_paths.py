import os

from config import CFG


def check_paths(path):
    head, tail = os.path.split(path)
    file_path = f"{CFG.DATA_PATH}/{head}/crops/dog/{tail}"
    return os.path.exists(file_path)
