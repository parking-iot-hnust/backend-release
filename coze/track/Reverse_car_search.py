from runtime import Args
from typings.Reverse_car_search.Reverse_car_search import Input, Output
import requests
import json
import string
from dataclasses import dataclass
import mysql.connector
from typing import TypedDict
from typing import List

db = mysql.connector.connect(
    host="yourhost",
    port=3306,
    user="root",
    password="password",
    connect_timeout=30,
    database="parking_lots"
)

cursor = db.cursor()


@dataclass
class Input:
    car_number: str
    now_location: str


def handler(args: Args[Input]) -> Output:
    '''
    根据车牌号查找对应的摄像头名称
    '''
    input = args.input
    logger = args.logger
    sql = "SELECT camera_location FROM parking_car WHERE numberplate = %s"
    cursor.execute(sql, (input.car_number,))  # 注意这里使用参数化查询
    results = cursor.fetchall()

    # 将查询结果转换为字符串列表
    camera_locations = [', '.join(row) for row in results]

    # 由于fetchall()返回的是列表的列表，我们需要将其转换为扁平的列表
    # 如果每个结果只有一列，可以直接使用列表推导式
    camera_locations_str = [item for item in camera_locations]

    db.close
    # 将列表转换为JSON格式的字符串
    return {"message": camera_locations_str}
