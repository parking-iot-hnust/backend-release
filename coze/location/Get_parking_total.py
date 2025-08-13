from runtime import Args
from typings.Get_parking_total.Get_parking_total import Input, Output
import requests
import json
import string
from runtime import Args
from typing import TypedDict, List
from dataclasses import dataclass
import requests
import mysql.connector

db = mysql.connector.connect(
    host="your_host",
    port=3306,
    user="root",
    password="password",
    connect_timeout=30,
    database="parking_lots"
)

def handler(args: Args[Input]) -> Output:
    input = args.input
    logger = args.logger
    total = 0

    cursor = db.cursor()  # 创建游标
    sql = "SELECT total FROM parking_lots_data WHERE parking_name = %s"
    cursor.execute(sql, (input.parking_name,))
    results = cursor.fetchall()

    if results:
        total_pass = [int(row[0]) for row in results]
        total = sum(total_pass)  # 累加所有结果中的总数
    else:
        total = 0
    db.close
    return {"total": total}
