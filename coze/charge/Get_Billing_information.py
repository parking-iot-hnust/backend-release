from runtime import Args
from typings.Get_Billing_information.Get_Billing_information import Input, Output
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
    connect_timeout  = 30,
    database="parking_lots"
)

cursor = db.cursor()

@dataclass
class Input:
    numberplate: str


def handler(args: Args[Input]) -> Output:
    input = args.input
    logger = args.logger

    # 多表查询——1
    sql = "SELECT parking_id FROM parking_info WHERE numberplate = %s"
    cursor.execute(sql, (input.numberplate,))  # 参数化查询
    results = cursor.fetchall()
    if(results) :
        parking_id1 = [str(row[0]) for row in results]  # 假设结果是一个元组列表，我们只取第一个元素
        parking_id = str(parking_id1[0])
        parking_id = parking_id[:8]

        # 多表查询——2
        sql2 = "SELECT * FROM parking_charge WHERE parking_id = %s"
        cursor.execute(sql2, (parking_id,))  # 参数化查询
        result = cursor.fetchall()

        # 分别将第一个子列表的元素赋值给其他变量
        first_sublist = result[0]
        id_number = first_sublist[0]  # "430302-1"
        reg1 = first_sublist[1]      # 1
        reg2 = first_sublist[2]      # 2
        reg3 = first_sublist[3]      # 1
        unit = first_sublist[4]      # 1
        free_time = first_sublist[5]      # 1
        #return {"message": value1}
        db.close
        return {
            "reg1" : reg1,
            "reg2" : reg2,
            "reg3" : reg3,
            "unit" : unit,
            "free_time" : free_time
        }
    else :
        db.close
        return {
            "reg1" : 0,
            "reg2" : 0,
            "reg3" : 0,
            "unit" : 0,
            "free_time" : 0
        }