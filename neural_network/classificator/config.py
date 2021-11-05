class CFG:
    # Data path
    DATA_PATH = "./data/detect"
    PATH_CSV = "./data/all.csv"
    PATH_OWNER = "./data/colorstails.csv"
    # Logging
    LOG_DIR = "./logs"
    OUTPUT_DIR = "test"
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
    epochs = 50
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
