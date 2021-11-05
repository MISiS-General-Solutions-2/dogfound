import math
import sys

import torch
from torch import nn
from torchvision import models
from torchvision import transforms
from classificator.dog_type_classificator import _read_dict
import cv2

file = "./classificator/imagenet_classes.txt"

dogs_indx_list = [i for i in range(151, 276)]

device = torch.device('cuda:0' if torch.cuda.is_available() else 'cpu')
sys.path.append('yolov5/')

yolo_model = torch.load('models/detect/yolo_finetuned_v2.pt',
                        map_location=device)['model'].float().eval().autoshape()
yolo_model = yolo_model.to(device)

yolo_model.conf = 0.27
yolo_model.iou = 0.45


class MultiOutputModel(nn.Module):
    def __init__(self, n_color_classes=3, n_tail_classes=2):
        super().__init__()
        self.resnet = models.resnet34(pretrained=False)
        self.base_model = nn.Sequential(
            *(list(self.resnet.children())[:-1])
        )  # take the model without classifier

        last_channel = (
            models.resnet34().fc.in_features
        )  # size of the layer before the classifier

        # create separate classifiers for our outputs
        self.color = nn.Sequential(
            nn.Dropout(p=0.5),
            nn.Linear(in_features=last_channel, out_features=n_color_classes),
        )
        self.tail = nn.Sequential(
            nn.Dropout(p=0.5),
            nn.Linear(in_features=last_channel, out_features=n_tail_classes),
        )

    def forward(self, x):
        x = self.base_model(x)

        # reshape from [batch, channels, 1, 1] to [batch, channels] to put it into classifier
        x = torch.flatten(x, start_dim=1)
        return {
            "color": self.color(x),
            "tail": self.tail(x),
        }


classifier_model = MultiOutputModel()
classifier_model = classifier_model.to(device)
weigths = torch.load('models/classifier/resnet.pt',
                     map_location=device)['model']
classifier_model.load_state_dict(weigths)

breed_model = models.resnet34(pretrained=True)

tr_pipe = transforms.Compose([
    transforms.ToPILImage(),
    transforms.Resize(224),
    transforms.ToTensor(),
    transforms.Normalize((0.485, 0.456, 0.406),
                         (0.229, 0.224, 0.225)),
])


def eval_on_image(file_path: str):
    nd_array = cv2.imread(file_path)
    nd_array = cv2.cvtColor(nd_array, cv2.COLOR_BGR2RGB)

    eval_res = yolo_model(nd_array)
    eval_res = eval_res.pandas().xyxy[0]

    eval_res = eval_res.query(
        "name=='dog' or name=='person' or name=='cat' or name=='bird'")

    return eval_res


def run_analytics(file_path, response):
    cv_image = cv2.imread(file_path)
    cv_image = cv2.cvtColor(cv_image, cv2.COLOR_BGR2RGB)
    res_df = yolo_model(cv_image).pandas().xywh[0]
    animals = res_df.query("name=='dog' or name=='cat' or name=='bird'")
    if len(animals) > 0:
        response.is_animal_there = 1
        dogs = res_df.query("name=='dog'")
        if len(dogs) > 0:
            response.is_it_a_dog = 1
            best_dog = dogs.sort_values(
                by='confidence', ascending=False).iloc[0]
            coords = (int(best_dog['ycenter']-best_dog['height']//2),
                      int(best_dog['ycenter']+best_dog['height']//2),
                      int(best_dog['xcenter']-best_dog['width']//2),
                      int(best_dog['xcenter']+best_dog['width']//2))
            response.additional.crop = [
                coords[2], coords[0], coords[3], coords[1]]
            dog_crop = cv_image[coords[0]:coords[1], coords[2]:coords[3]]
            # Define tail and color classes
            classes = classifier_model(
                tr_pipe(dog_crop).unsqueeze(0).to(device))
            dog_breed = breed_model(tr_pipe(dog_crop).unsqueeze(0).to(device))
            color_cl = classes['color'].argmax().cpu().detach().item() + 1
            tail_cl = classes['tail'].argmax().cpu().detach().item() + 1
            response.color = color_cl
            response.tail = tail_cl
            # Define dog breed
            dog_breed_cl = ""
            label_dict = _read_dict(file)
            dog_breed_idx = dog_breed.argmax().cpu().detach().item()
            if dog_breed_idx in dogs_indx_list:
                dog_breed_cl = label_dict[dog_breed_idx]

            response.breed = dog_breed_cl
            humans = res_df.query("name=='person'")
            if len(humans) > 0:
                best_dog_coords = best_dog['xcenter'], best_dog['ycenter']
                for i, r in humans.iterrows():
                    person_coords = r['xcenter'], r['ycenter']
                    dist = math.dist(person_coords, best_dog_coords)
                    if dist <= 300:
                        response.is_the_owner_there = 1
    return response
