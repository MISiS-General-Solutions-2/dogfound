import torch
import cv2
import pandas as pd

torch.device('cuda:0' if torch.cuda.is_available() else 'cpu')
yolo_model = torch.hub.load(
    'ultralytics/yolov5', 'custom', path='models/detect/yolo_finetuned_v2.pt')

yolo_model.conf = 0.27
yolo_model.iou = 0.45


def eval_on_image(file_path: str):
    nd_array = cv2.imread(file_path)
    nd_array = cv2.cvtColor(nd_array, cv2.COLOR_BGR2RGB)

    eval_res = yolo_model(nd_array)
    eval_res = eval_res.pandas().xyxy[0]

    eval_res = eval_res.query(
        "name=='dog' or name=='person' or name=='cat' or name=='bird'")

    return eval_res


def run_analytics(result_df, response):

    return response
