from runtime import Args
from typings.get_time.get_time import Input, Output
import requests
import json
import string
from dataclasses import dataclass
import mysql.connector
from typing import TypedDict
from typing import List

db = mysql.connector.connect(
    host="your_host",
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


# @dataclass
# class Output:
#     parking_time_seconds: int

def handler(args: Args[Input]) -> Output:
    input = args.input
    logger = args.logger

    sql = "SELECT parking_time FROM parking_info WHERE numberplate = %s order by record_id  DESC LIMIT 1"

    cursor.execute(sql, (input.car_number,))  # 参数化查询
    results = cursor.fetchall()

    if (results):
        # 将查询结果转换为整数列表
        time_int_list = [int(row[0]) for row in results]  # 将每个元组的第一个元素转换为整数

        # 如果只需要一个时间值，可以进一步处理 time_int_list
        # 例如，如果停车时间是按顺序排列的，你可以取第一个或最后一个时间
        first_time = time_int_list[0]  # if time_int_list else None
        # return {"time":time_int_list}
        # 根据需要返回数据，这里返回第一个时间值作为示例
        return {"time": first_time}
    db.close
    return {"time": 0}


