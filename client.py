#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Copyright (c) 2025 LTQY. All rights reserved.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
"""

import argparse
from datetime import datetime

import cv2
import numpy as np
import os
import sys
import tritonclient.http as httpclient
import math
import requests
import json
import time

# 存储字符到数字的映射关系
text_map = {}
# 存储车牌到停车时间的映射关系
mp1 = {}
# 存储车牌离开时间的映射关系
mp2 = {}
# 停车时间间隔，单位为秒
in_seconds = 300


def getLabelMap():
    """从文件中读取字符到数字的映射关系"""
    with open('rec_word_dict.txt', 'r') as file:
        text_data = file.readlines()

    for i, line in enumerate(text_data, start=1):
        text_map[i] = line.strip()


def decode(text_index, text_prob=None, is_remove_duplicate=False):
    """将文本索引转换为文本标签"""
    result_list = []
    ignored_tokens = [0]
    batch_size = len(text_index)
    for batch_idx in range(batch_size):
        selection = np.ones(len(text_index[batch_idx]), dtype=bool)
        if is_remove_duplicate:
            selection[1:] = text_index[batch_idx][1:] != text_index[batch_idx][:-1]
        for ignored_token in ignored_tokens:
            selection &= text_index[batch_idx] != ignored_token
        if text_prob is not None:
            selection &= text_prob[batch_idx] >= 0

        char_list = [
            text_map[text_id]
            for text_id in text_index[batch_idx][selection]
        ]
        if text_prob is not None:
            conf_list = text_prob[batch_idx][selection]
        else:
            conf_list = [1] * len(selection)
        if len(conf_list) == 0:
            conf_list = [0]

        text = ''.join(char_list)
        result_list.append((text, np.mean(conf_list).tolist()))
    return result_list


def post_process_image(results):
    """处理推理结果"""
    output1 = results.as_numpy('softmax_2.tmp_0')
    preds_idx = output1.argmax(axis=2)
    preds_prob = output1.max(axis=2)
    res = decode(preds_idx, preds_prob, is_remove_duplicate=True)
    return res


def resize_norm_img(img):
    """调整图像大小并进行归一化处理"""
    imgC, imgH, imgW = (3, 48, 64)
    h, w = img.shape[:2]
    wh_ratio = w * 1.0 / h

    max_wh_ratio = imgW / imgH
    max_wh_ratio = max(max_wh_ratio, wh_ratio)

    assert imgC == img.shape[2]
    imgW = int((imgH * max_wh_ratio))

    ratio = w / float(h)

    if math.ceil(imgH * ratio) > imgW:
        resized_w = imgW
    else:
        resized_w = int(math.ceil(imgH * ratio))
    resized_image = cv2.resize(img, (resized_w, imgH))
    resized_image = resized_image.astype('float32')
    resized_image = resized_image.transpose((2, 0, 1)) / 255
    resized_image -= 0.5
    resized_image /= 0.5
    padding_im = np.zeros((imgC, imgH, imgW), dtype=np.float32)
    padding_im[:, :, 0:resized_w] = resized_image
    return padding_im


def pre_process_image(images):
    """对图像列表进行预处理"""
    images = resize_norm_img(images)
    return images


def deal_plate(plate, url, parking_id):
    """处理车牌信息"""
    current_time = datetime.now()
    if plate not in mp1:
        mp1[plate] = 0
    elif mp1[plate] == 0:
        mp1[plate] = current_time
    elif (current_time - mp1[plate]).total_seconds() >= in_seconds:
        mp2[plate] = mp1[plate]
        send_message(plate, url, parking_id, 0, mp1[plate])
        mp1[plate] = -1
    for a in list(mp2.keys()):
        b = mp2[a]  # 获取车牌a的停车时间
        time = datetime.now() # 获取当前时间
        if (time - b).total_seconds() >= in_seconds:  # 如果离开时间大于等于in_seconds
            del mp2[a]  # 删除车牌a的记录
            send_message(a, url, parking_id, 1, time)
            if a in mp1:
                mp1[a] = 0  # 将时间间隔记录重置为0


def send_message(plate, url, parking_id, type, parking_time):
    """发送HTTP POST请求到服务器"""
    data = {
        "plate": plate,
        "parking_id": parking_id,
        "parking_time": parking_time.isoformat(),
        "type": type
    }
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, json=data, headers=headers)

    if response.ok:
        print("消息发送成功")
    else:
        print(f"发送失败，状态码：{response.status_code}")


def image_batch(image_list, batch_size):
    """生成图像批次，用于批量推理"""
    length = len(image_list)
    idx = 0
    while idx < length:
        images = image_list[idx:idx + batch_size]
        idx = idx + batch_size
        yield images


def main(args):
    """程序的主入口"""
    triton_client = httpclient.InferenceServerClient(url=args.server_url)
    cap = cv2.VideoCapture(args.video_url)
    print(cap.isOpened())
    while True:
        ret, frame = cap.read()

        if not ret:
            break
        image = pre_process_image(frame)
        request_images = np.expand_dims(image, axis=0)

        inputs = []
        inputs.append(
            httpclient.InferInput('x', request_images.shape, "FP32")
        )
        inputs[0].set_data_from_numpy(request_images, binary_data=False)

        outputs = []
        outputs.append(
            httpclient.InferRequestedOutput('softmax_2.tmp_0', binary_data=False)
        )
        results = triton_client.infer(
            'm-official-39', inputs=inputs, outputs=outputs)

        results = post_process_image(results)

        for text, confidence in results:
            if len(text) >= 8:
                deal_plate(text, args.post_url, args.parking_id)
                print("successfully send message")
                print(f"plate: {text}, confidence: {confidence}")
            cv2.putText(frame, text, (50, 50), cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 255, 0), 2)
        cv2.imshow('Frame', frame)
        if cv2.waitKey(1) & 0xFF == ord('q'):
            break
    cap.release()
    cv2.destroyAllWindows()


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(description='Triton HTTP Client')
    parser.add_argument("--server_url", help="Server http url", required=False,
                         type=str)
    parser.add_argument("--video_url", help="Video", required=False,
                        default="/dev/video", type=str)
    parser.add_argument("--parking_id", help="parking_id", required=False,
                         type=str)
    parser.add_argument("--post_url", help="post_url", required=False,
                         type=str)
    parser.add_argument("--max_batch", help="Max batch size",required=False,
                        default=1, type=int)
    args = parser.parse_args()
    return args

if __name__ == '__main__':
    getLabelMap()
    args = parse_args()
    sys.exit(main(args))