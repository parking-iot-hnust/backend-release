#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Copyright (c) 2025 LTQY. All rights reserved.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.
"""

from runtime import Args
from typings.GetDeviceID.GetDeviceID import Input, Output
import mysql.connector

db=mysql.connector.connect(
    host="your_host",
    port=3306,
    user="root",
    password="password",
    connect_timeout=30,
    database="parking_lots"
)

def handler(args: Args[Input])->Output:
    input = args.input
    target =  input.target
    results = select_mysql(target)
    parking_id = solve(results)
    db.close
    return {"parking_id": parking_id}

def select_mysql(target:str)->list:
    cursor = db.cursor()
    sql = "select parking_id from parking_lots_data where parking_name = %s"%target
    cursor.execute(sql)
    results = cursor.fetchall()
    # db.close()
    return results

def solve(results:list)->str:
    parking_id = ""
    for result in results:
        parking_id = result
    parking_id = parking_id[0]
    return parking_id
