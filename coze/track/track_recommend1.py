from runtime import Args
from typings.track_recommend.track_recommend import Input, Output
import requests
from dataclasses import dataclass


@dataclass
class Input:
    start_lat: float
    start_lng: float
    end_lat: float
    end_lng: float


def handler(args: Args[Input]) -> Output:
    '''
    处理请求，获取导航路线

    参数：
    args (Args[Input])：包含输入参数的对象：两对经纬度

    返回值：
    Output：包含导航路线信息的对象
    '''
    input = args.input
    logger = args.logger
    url = f"https://restapi.amap.com/v3/direction/driving"
    params = {
        "strategy": 1,
        "key": "your_api_key",
        "origin": f"{input.start_lng},{input.start_lat}",
        "destination": f"{input.end_lng},{input.end_lat}",
    }

    response = requests.get(url, params=params)
    data = response.json()
    logger.warning(data)

    if data["status"] == "0":
        return {"message": "无法获取路线信息"}
    routes = data["route"]["paths"]
    if not routes:
        return {"message": "无法获取路线信息"}

    # 输出完整的路径推荐
    output = []
    for lot in routes:
        strategy = lot['strategy']
        track = ""
        for step in lot["steps"]:
            track += "; " + step['instruction']
        output.append({
            "strategy": strategy,
            "track": track
        })

    # 返回
    return {
        "message": output
    }

