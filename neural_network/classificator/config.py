class CFG:
    # Data path
    DATA_PATH = "./data/new_detect"
    CLEAN_DATA_PATH = "./data/clean_dogs"
    PATH_CSV = "./data/train_remarked.csv"
    PATH_OWNER = "./data/owner+dog.csv"
    CLEAN_CSV = "./data/clean_dogs.csv"
    NEW_CSV = "./data/new_images.csv"
    # Logging
    LOG_DIR = "./logs"
    OUTPUT_DIR = "resnet34_baseline+all_data_freeze-true"
    TEST_PATH = ""

    # Model setup
    model_name = "resnet34"
    pretrained = True
    freeze = True

    # Main config
    GPU_ID = 1
    seed = 42
    target_size = 5

    # Train configs
    epochs = 75
    early_stopping = 10
    batch_size = 16
    size = 224
    mean = [0.485, 0.456, 0.406]  # ImageNet values
    std = [0.229, 0.224, 0.225]  # ImageNet values
    num_workers = 8
    print_freq = 5

    # Optimizer config
    lr = 1e-4
    momentum = 0.9
    min_lr = 1e-2
    weight_decay = 1e-5
