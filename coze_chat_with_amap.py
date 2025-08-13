#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Copyright (c) 2025 LTQY. All rights reserved.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.

This example is about how to use the streaming interface to start a chat request
and handle chat events with Amap location support
"""

import os
import requests
import json
# Our official coze sdk for Python [cozepy](https://github.com/coze-dev/coze-py)
from cozepy import COZE_CN_BASE_URL

# Get an access_token through personal access token or oauth.
coze_api_token = 'pat_XXXXXX'  # 请替换为正确的PAT Token (以pat_开头)
# The default access is api.coze.com, but if you need to access api.coze.cn,
# please use base_url to configure the api endpoint to access
coze_api_base = COZE_CN_BASE_URL

from cozepy import Coze, TokenAuth, Message, ChatStatus, MessageContentType, ChatEventType  # noqa

# Init the Coze client through the access_token.
coze = Coze(auth=TokenAuth(token=coze_api_token), base_url=coze_api_base)

# Create a bot instance in Coze, copy the last number from the web link as the bot's ID.
bot_id = 'bot-id'

# 调试模式开关：设置为True时会显示发送给AI的完整消息
DEBUG_MODE = False

# 高德地图API配置
AMAP_API_KEY = 'api-key'  # 您的高德地图API Key

# 位置获取方式配置
LOCATION_METHOD = 'amap'  # 'ip' 使用IP定位, 'amap' 使用高德地图定位

# The user id identifies the identity of a user. Developers can use a custom business ID
# or a random string.
user_id = '123'

def get_location_by_ip():
    """通过IP地址获取位置信息"""
    try:
        # 获取IP地址
        ip_response = requests.get('https://httpbin.org/ip', timeout=5)
        current_ip = ip_response.json()['origin']
        
        # 获取位置信息
        location_response = requests.get(f'http://ip-api.com/json/{current_ip}', timeout=5)
        location_data = location_response.json()
        
        if location_data.get('status') == 'success':
            lat = location_data.get('lat', 0)
            lon = location_data.get('lon', 0)
            city = location_data.get('city', '未知')
            region = location_data.get('regionName', '未知')
            return lat, lon, city, region, 'IP定位'
        
    except Exception:
        pass
    
    return None

def get_location_by_amap():
    """通过高德地图API获取位置信息"""
    try:
        # 首先获取IP地址
        ip_response = requests.get('https://httpbin.org/ip', timeout=5)
        current_ip = ip_response.json()['origin']
        
        # 使用高德地图IP定位API
        amap_url = f'https://restapi.amap.com/v3/ip?ip={current_ip}&key={AMAP_API_KEY}'
        amap_response = requests.get(amap_url, timeout=10)
        amap_data = amap_response.json()
        
        if amap_data.get('status') == '1' and amap_data.get('info') == 'OK':
            # 解析高德返回的位置信息
            province = amap_data.get('province', '未知')
            city = amap_data.get('city', '未知')
            
            # 如果有矩形区域信息，计算中心点
            rectangle = amap_data.get('rectangle', '')
            if rectangle:
                coords = rectangle.split(';')
                if len(coords) == 2:
                    # 左下角和右上角坐标
                    bottom_left = coords[0].split(',')
                    top_right = coords[1].split(',')
                    
                    if len(bottom_left) == 2 and len(top_right) == 2:
                        # 计算中心点
                        lon = (float(bottom_left[0]) + float(top_right[0])) / 2
                        lat = (float(bottom_left[1]) + float(top_right[1])) / 2
                        return lat, lon, city, province, '高德地图IP定位'
            
            # 如果没有矩形信息且有城市信息，使用城市中心坐标
            if city and city != '[]' and city != '未知':
                # 获取城市的地理编码
                geo_url = f'https://restapi.amap.com/v3/geocode/geo?address={city}&key={AMAP_API_KEY}'
                geo_response = requests.get(geo_url, timeout=10)
                geo_data = geo_response.json()
                
                if geo_data.get('status') == '1' and geo_data.get('geocodes'):
                    location = geo_data['geocodes'][0]['location']
                    lon, lat = map(float, location.split(','))
                    return lat, lon, city, province, '高德地图地理编码'
        
    except Exception as e:
        if DEBUG_MODE:
            print(f"高德地图定位失败: {e}")
    
    return None

def get_current_location():
    """获取当前位置的经纬度信息（不在界面显示）"""
    location_result = None
    
    # 根据配置选择定位方式
    if LOCATION_METHOD == 'amap':
        location_result = get_location_by_amap()
        if not location_result:
            # 高德地图失败，降级到IP定位
            location_result = get_location_by_ip()
    elif LOCATION_METHOD == 'ip':
        location_result = get_location_by_ip()
    
    # 如果所有方法都失败，使用默认位置（湖南科技大学）
    if not location_result:
        return 27.9087, 112.6112, "湘潭", "湖南", '默认位置(湖南科技大学)'
    
    return location_result

def simple_chat(message):
    """简单对话函数"""
    print(f"用户输入: {message}")
    
    # 获取当前位置信息（静默获取，不显示给用户）
    lat, lon, city, region, method = get_current_location()
    
    # 构建包含位置信息的消息
    location_info = f"[位置信息] 用户当前位置: {region} {city}, 经纬度: {lat}, {lon}"
    enhanced_message = f"{message}\n\n{location_info}"
    
    # 调试模式：显示发送给AI的完整消息
    if DEBUG_MODE:
        print(f"[调试] 定位方式: {method}")
        print(f"[调试] 发送给AI的完整消息:\n{enhanced_message}\n")
    
    print("AI回复: ", end="", flush=True)
    
    # Call the coze.chat.stream method to create a chat. The create method is a streaming
    # chat and will return a Chat Iterator. Developers should iterate the iterator to get
    # chat event and handle them.
    for event in coze.chat.stream(
        bot_id=bot_id,
        user_id=user_id,
        additional_messages=[
            Message.build_user_question_text(enhanced_message),
        ],
    ):
        if event.event == ChatEventType.CONVERSATION_MESSAGE_DELTA:
            print(event.message.content, end="", flush=True)

        if event.event == ChatEventType.CONVERSATION_CHAT_COMPLETED:
            print()
            print(f"token usage: {event.chat.usage.token_count}")
            print("-" * 50)

def interactive_chat():
    """交互式对话"""
    print("=== Coze 智能体对话（高德地图定位）===")
    print("💡 提示: 每次对话都会自动附加您的当前位置信息")
    print(f"🗺️  定位方式: {LOCATION_METHOD.upper()}")
    print("🔧 调试: 将 DEBUG_MODE 设为 True 可查看发送给AI的完整消息")
    print("输入 'quit' 或 'exit' 退出对话")
    print("-" * 60)
    
    # 显示当前位置（仅在启动时显示一次）
    lat, lon, city, region, method = get_current_location()
    print(f"📍 检测到您的位置: {region} {city} ({lat}, {lon})")
    print(f"🔍 定位方式: {method}")
    print("-" * 60)
    
    while True:
        try:
            user_input = input("\n请输入您的问题: ").strip()
            
            if user_input.lower() in ['quit', 'exit', '退出']:
                print("再见!")
                break
                
            if not user_input:
                print("请输入有效的问题")
                continue
                
            simple_chat(user_input)
            
        except KeyboardInterrupt:
            print("\n\n程序被中断，再见!")
            break
        except Exception as e:
            print(f"\n发生错误: {e}")
            print("请重试...")

if __name__ == "__main__":
    # 可以选择直接发送单个消息或进入交互模式
    
    # 选项1: 发送单个消息（测试用）
    # simple_chat("告诉我附近有什么好吃的餐厅")
    
    # 选项2: 交互式对话（推荐）
    interactive_chat()
    
    # 使用说明：
    # - 设置 LOCATION_METHOD = 'amap' 使用高德地图定位（推荐，更准确）
    # - 设置 LOCATION_METHOD = 'ip' 使用IP定位（备选方案）
    # - 设置 DEBUG_MODE = True 查看完整的发送消息和定位过程
