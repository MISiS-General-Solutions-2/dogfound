import torch.nn.functional as F


def get_loss(net_output, ground_truth, device):
    ground_truth["color"] = ground_truth["color"].to(device)
    ground_truth["tail"] = ground_truth["tail"].to(device)
    color_loss = F.cross_entropy(net_output["color"], ground_truth["color"])
    tail_loss = F.cross_entropy(net_output["tail"], ground_truth["tail"])
    loss = 0.6*color_loss + 0.4*tail_loss
    return loss, {"color": color_loss, "tail": tail_loss}
