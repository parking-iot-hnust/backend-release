#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Copyright (c) 2025 LTQY. All rights reserved.

This source code is licensed under the MIT license found in the
LICENSE file in the root directory of this source tree.

快速获取当前位置的IP地址和位置信息
"""

import requests
import json

def get_my_location():
    """快速获取我的位置信息"""
    print("🔍 正在定位您的位置...")
    print("=" * 40)
    
    try:
        # 方法1: 获取IP地址
        print("📍 获取您的IP地址...")
        ip_response = requests.get('https://httpbin.org/ip', timeout=5)
        my_ip = ip_response.json()['origin']
        print(f"   您的公网IP: {my_ip}")
        
        # 方法2: 获取详细位置信息
        print("\n🌍 获取位置详情...")
        location_response = requests.get(f'http://ip-api.com/json/{my_ip}?lang=zh-CN', timeout=5)
        location_data = location_response.json()
        
        if location_data.get('status') == 'success':
            print(f"   🏳️  国家: {location_data.get('country', '未知')}")
            print(f"   🏛️  省份: {location_data.get('regionName', '未知')}")
            print(f"   🏙️  城市: {location_data.get('city', '未知')}")
            print(f"   📮 邮编: {location_data.get('zip', '未知')}")
            print(f"   🌐 运营商: {location_data.get('isp', '未知')}")
            print(f"   📡 组织: {location_data.get('org', '未知')}")
            
            # 显示坐标
            if location_data.get('lat') and location_data.get('lon'):
                lat, lon = location_data.get('lat'), location_data.get('lon')
                print(f"   📍 坐标: {lat}, {lon}")
                print(f"   🗺️  地图链接:")
                print(f"      Google: https://www.google.com/maps?q={lat},{lon}")
                print(f"      百度: https://api.map.baidu.com/marker?location={lat},{lon}&title=您的位置")
                print(f"      高德: https://uri.amap.com/marker?position={lon},{lat}&name=您的位置")
        else:
            print("   ❌ 位置信息获取失败")
            
        # 方法3: 尝试获取更精确的信息
        print("\n🎯 尝试获取更精确的位置...")
        try:
            precise_response = requests.get(f'https://ipapi.co/{my_ip}/json/', timeout=5)
            precise_data = precise_response.json()
            
            if not precise_data.get('error'):
                print(f"   🏢 地区: {precise_data.get('region', '未知')}")
                print(f"   🕐 时区: {precise_data.get('timezone', '未知')}")
                print(f"   🏦 ASN: {precise_data.get('asn', '未知')}")
                print(f"   🌐 网络: {precise_data.get('network', '未知')}")
            else:
                print(f"   ⚠️  精确位置查询受限: {precise_data.get('reason', '未知')}")
        except:
            print("   ⚠️  精确位置服务暂不可用")
            
    except Exception as e:
        print(f"❌ 获取位置信息时发生错误: {e}")
    
    print("\n" + "=" * 40)
    print("✅ 位置查询完成!")

def get_local_ip():
    """获取本地网络IP地址"""
    import socket
    
    print("\n💻 本地网络信息:")
    print("-" * 30)
    
    try:
        # 获取本机名
        hostname = socket.gethostname()
        print(f"   🖥️  主机名: {hostname}")
        
        # 获取本地IP
        local_ip = socket.gethostbyname(hostname)
        print(f"   🏠 本地IP: {local_ip}")
        
        # 尝试获取更准确的本地IP（通过连接外部服务器）
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(("8.8.8.8", 80))
        real_local_ip = s.getsockname()[0]
        s.close()
        print(f"   📶 实际内网IP: {real_local_ip}")
        
    except Exception as e:
        print(f"   ❌ 获取本地IP失败: {e}")

if __name__ == "__main__":
    print("🎯 IP地址和位置快速查询工具")
    print("Author: GitHub Copilot")
    print()
    
    # 获取公网位置
    get_my_location()
    
    # 获取本地网络信息
    get_local_ip()
    
    print("\n" + "🔄 需要重新查询请重新运行程序")
