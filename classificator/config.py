class CFG:
    # Data path
    TRAIN_PATH = "./data/train"
    VAL_PATH = "./data/val"
    # Logging
    LOG_DIR = "./logs"
    OUTPUT_DIR = "resnet34_baseline"

    # Model setup
    # chk = "./models/resnet34_best.pt"
    chk = ""
    model_name = "resnet34"
    pretrained = True
    freeze = True

    # Main config
    GPU_ID = 0
    seed = 42
    target_size = 6

    # Train configs
    epochs = 150
    early_stopping = 10
    batch_size = 16
    size = 224
    MEAN = [0.485, 0.456, 0.406]  # ImageNet values
    STD = [0.229, 0.224, 0.225]  # ImageNet values
    num_workers = 8
    print_freq = 5

    # Optimizer config
    lr = 1e-4
    momentum = 0.9
    min_lr = 1e-2
    weight_decay = 1e-5
