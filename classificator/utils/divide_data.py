import os
import random
import shutil

DATA_PATH = "../data"
INPUT_PATH = os.path.join(DATA_PATH, 'detect')


def main():
    list_dirs = os.listdir(INPUT_PATH)

    for dir_name in list_dirs:
        if dir_name == "owner+dog":
            continue
        train_path = os.path.join(DATA_PATH, "train", dir_name)
        val_path = os.path.join(DATA_PATH, "val", dir_name)
        # Create dirs to copy the files
        os.makedirs(train_path, exist_ok=True)
        os.makedirs(val_path, exist_ok=True)
        dog_path = os.path.join(INPUT_PATH, dir_name, "crops", "dog")
        for image_file in os.listdir(dog_path):
            if random.uniform(0, 1) > 0.25:
                shutil.copy(os.path.join(dog_path, image_file), train_path)
            else:
                shutil.copy(os.path.join(dog_path, image_file), val_path)


if __name__ == "__main__":
    main()
