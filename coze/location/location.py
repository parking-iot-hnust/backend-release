from runtime import Args
from typings.location.location import Input, Output
import requests
import json
import string
from runtime import Args
from typing import TypedDict, List
from dataclasses import dataclass
import requests
import mysql.connector

db = mysql.connector.connect(
    host="yourhost",
    port=3306,
    user="root",
    password="password",
    database="parking_lots"
)


# Output
class Output(TypedDict, total=False):
    """传递附近是否有停车场"""
    message: str
    """停车场的信息"""
    parking_lots: List[dict]


# Input
@dataclass
class Input:
    """目标地点"""
    target: str


def handler(args: Args[Input]) -> Output:
    input = args.input
    logger = args.logger
    url = "https://restapi.amap.com/v3/place/around"
    params = {
        "key": "your_api_key",
        "location": get_coordinates(input.target),
        "radius": 500,
        "types": "停车场",
        "offset": 10,
        "extensions": "all"
    }
    response = requests.get(url, params=params)
    data = response.json()
    if "pois" not in data:
        logger.warning("No parking lots found.")
        db.close()
        return {"message": "No parking lots found."}
    parking_lots = data["pois"]
    output = []
    for lot in parking_lots:
        total, now_parking = search(lot, logger)
        output.append({
            "name": lot["name"],
            "location": lot["location"],
            "tatol": total,
            "now_parking": now_parking
        })
    db.close
    return {"message": "Parking lots found.", "parking_lots": output}


def get_coordinates(place_name: str) -> str:
    """
    将地点名称转换为经纬度坐标

    参数:
    place_name (str): 地点名称

    返回:
    str: 经纬度坐标，格式为 "latitude,longitude"
    """
    url = f"https://restapi.amap.com/v3/geocode/geo?key=keyxxxxxxxxx&address={place_name}"
    response = requests.get(url)
    data = response.json()
    if "geocodes" not in data:
        raise ValueError(f"无法获取 {place_name} 的经纬度坐标")
    geocode = data["geocodes"][0]
    location = geocode["location"]
    return location


def search(lot, logger):
    res = manage_parking_data(lot["name"], lot["adcode"], lot["entr_location"], logger)
    return res


def get_next_parking_id(area_code, logger):
    cursor = db.cursor()
    # 查询当前地区编码下的最大parking_id的后半部分数字
    sql = f"""
        SELECT MAX(CAST(SUBSTRING_INDEX(parking_id, '-', -1) AS UNSIGNED)) AS max_id
        FROM parking_lots_data
        WHERE parking_id LIKE '{area_code}-%'
        """
    cursor.execute(sql)
    result = cursor.fetchone()

    cursor.close()
    if result:
        # 如果存在，返回当前最大值加1
        return area_code + '-' + str(result[0] + 1)
    else:

        # 如果不存在，从1开始
        return area_code + '-1'


def parse_location(entr_location):
    latitude, longitude = map(float, entr_location.split(','))
    return latitude, longitude


def manage_parking_data(parking_name, area_code, entr_location, logger):
    cursor = db.cursor()
    total = 0
    now_parking = 0
    # 解析经纬度
    latitude, longitude = parse_location(entr_location)
    # 查询数据库
    cursor.execute(
        """
        SELECT total, now_parking
        FROM parking_lots_data
        WHERE parking_name = %s AND SUBSTRING(parking_id, 1, 6) = %s
        """, (parking_name, area_code))
    result = cursor.fetchone()

    if result:
        # 如果存在结果，获取total和now_parking
        total, now_parking = result
    else:
        # 如果不存在结果，生成新的parking_id并插入数据
        new_parking_id = get_next_parking_id(area_code, logger)
        sql = """
        INSERT INTO parking_lots_data (parking_id, parking_name, latitude, longitude, total, now_parking)
        VALUES (%s, %s, %s, %s, %s, %s)
        """
        cursor.execute(sql, (new_parking_id, parking_name, latitude, longitude, 0, 0))
        db.commit()
    return total, now_parking
