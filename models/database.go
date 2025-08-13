/*
 * Copyright (c) 2025 LTQY. All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package models

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// 从环境变量获取数据库配置
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/urban_traffic?charset=utf8mb4&parseTime=True&loc=Local"
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Connected to MySQL database successfully")

	// 自动迁移数据库表
	err = DB.AutoMigrate(
		&User{}, &Vehicle{}, &ParkingRecord{}, &ParkingLot{}, &SpecialSpot{}, &ParkingSession{},
		// 交通相关表
		&TrafficFlow{}, &TrafficUserStats{}, &TrafficHeatmap{}, &CongestionReport{},
		&InOutFlowData{}, &CarCrossingRate{},
		// 设备相关表
		&Device{}, &DeviceExpense{}, &DeviceMaintenanceRecord{}, &DeviceFaultStats{},
		&DeviceAlarm{}, &DeviceAlarmStats{},
		// 统计相关表
		&ParkingSaturation{}, &ParkingOccupancyRate{}, &TotalOccupancy{},
		&MotorParkingCongestion{}, &TollRecord{}, &ConstructionStats{}, &MonitoringCamera{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建默认用户和测试数据
	createDefaultUsers()
	createTestVehicles()
	createTestParkingSessions()

	// 创建模拟数据
	CreateSimulationData()
}

func createDefaultUsers() {
	var count int64
	DB.Model(&User{}).Count(&count)

	if count == 0 {
		defaultUsers := []User{
			{
				Username: "admin",
				Password: "admin",
				UserType: "admin",
				Email:    "admin@example.com",
			},
			{
				Username: "user",
				Password: "user",
				UserType: "user",
				Email:    "user@example.com",
			},
		}

		for _, user := range defaultUsers {
			user.HashPassword()
			DB.Create(&user)
		}
		log.Println("Default users created")
	}
}

func createTestVehicles() {
	var count int64
	DB.Model(&Vehicle{}).Count(&count)

	if count == 0 {
		// 获取普通用户ID - 所有测试车辆都分配给普通用户
		var normalUser User
		DB.Where("username = ?", "user").First(&normalUser)

		testVehicles := []Vehicle{
			{
				UserID:      normalUser.ID,
				PlateNumber: "浙A·12345",
				Brand:       "特斯拉",
				Model:       "Model 3",
				Color:       "白色",
				Type:        "新能源车",
				RegDate:     "2022-05-15",
				IsDefault:   true,
			},
			{
				UserID:      normalUser.ID,
				PlateNumber: "浙B·67890",
				Brand:       "本田",
				Model:       "CR-V",
				Color:       "灰色",
				Type:        "小型汽车",
				RegDate:     "2021-10-22",
				IsDefault:   false,
			},
			{
				UserID:      normalUser.ID,
				PlateNumber: "浙C·11111",
				Brand:       "丰田",
				Model:       "凯美瑞",
				Color:       "黑色",
				Type:        "小型汽车",
				RegDate:     "2023-01-10",
				IsDefault:   false,
			},
			{
				UserID:      normalUser.ID,
				PlateNumber: "浙D·22222",
				Brand:       "比亚迪",
				Model:       "秦PLUS DM-i",
				Color:       "蓝色",
				Type:        "新能源车",
				RegDate:     "2023-06-08",
				IsDefault:   false,
			},
			{
				UserID:      normalUser.ID,
				PlateNumber: "浙E·33333",
				Brand:       "小鹏",
				Model:       "P7",
				Color:       "银色",
				Type:        "新能源车",
				RegDate:     "2023-09-15",
				IsDefault:   false,
			},
		}

		for _, vehicle := range testVehicles {
			DB.Create(&vehicle)
		}

		log.Println("Test vehicles created")

		// 创建停车记录（在车辆创建之后）
		createTestParkingRecords()
	}
}

func createTestParkingRecords() {
	var vehicles []Vehicle
	DB.Find(&vehicles)

	// 首先创建停车场数据（如果不存在）
	createTestParkingLots()

	// 获取停车场ID
	var parkingLots []ParkingLot
	DB.Find(&parkingLots)

	for i, vehicle := range vehicles {
		// 为每辆车创建不同的停车记录
		lotIndex := i % len(parkingLots)
		parkingLot := parkingLots[lotIndex]

		records := []ParkingRecord{
			{
				VehicleID:    vehicle.ID,
				ParkingLotID: &parkingLot.ID,
				Location:     parkingLot.Name,
				StartTime:    "2023-08-16 09:30:00",
				EndTime:      "2023-08-16 11:30:00",
				Fee:          15.50,
				Duration:     2.0,
				SpotType:     "normal",
			},
			{
				VehicleID:    vehicle.ID,
				ParkingLotID: &parkingLot.ID,
				Location:     parkingLot.Name,
				StartTime:    "2023-08-10 14:00:00",
				EndTime:      "2023-08-10 17:00:00",
				Fee:          22.00,
				Duration:     3.0,
				SpotType:     "charging",
			},
		}

		for _, record := range records {
			DB.Create(&record)
		}
	}
}

func createTestParkingLots() {
	var count int64
	DB.Model(&ParkingLot{}).Count(&count)

	if count == 0 {
		testParkingLots := []ParkingLot{
			{
				Name:           "西湖停车场",
				Address:        "杭州市西湖区西湖景区内",
				Latitude:       30.2394,
				Longitude:      120.1509,
				TotalSpots:     120,
				AvailableSpots: 25,
				HourlyRate:     10.00,
				IsActive:       true,
				OperatingHours: "24小时",
				PaymentMethods: "微信,支付宝,现金",
			},
			{
				Name:           "中央公园停车场",
				Address:        "杭州市上城区中央公园旁",
				Latitude:       30.2635,
				Longitude:      120.1709,
				TotalSpots:     200,
				AvailableSpots: 42,
				HourlyRate:     10.00,
				IsActive:       true,
				OperatingHours: "24小时",
				PaymentMethods: "微信,支付宝,现金",
			},
			{
				Name:           "东方广场停车场",
				Address:        "杭州市江干区东方广场地下",
				Latitude:       30.2793,
				Longitude:      120.2194,
				TotalSpots:     80,
				AvailableSpots: 5,
				HourlyRate:     10.00,
				IsActive:       true,
				OperatingHours: "24小时",
				PaymentMethods: "微信,支付宝,现金",
			},
		}

		for _, lot := range testParkingLots {
			DB.Create(&lot)
		}

		// 创建特殊车位
		createTestSpecialSpots()
		log.Println("Test parking lots created")
	}
}

func createTestSpecialSpots() {
	var parkingLots []ParkingLot
	DB.Find(&parkingLots)

	for _, lot := range parkingLots {
		specialSpots := []SpecialSpot{
			{
				ParkingLotID:   lot.ID,
				SpotType:       "charging",
				TotalCount:     5,
				AvailableCount: 3,
				AdditionalFee:  5.00,
			},
			{
				ParkingLotID:   lot.ID,
				SpotType:       "disabled",
				TotalCount:     3,
				AvailableCount: 2,
				AdditionalFee:  0.00,
			},
			{
				ParkingLotID:   lot.ID,
				SpotType:       "vip",
				TotalCount:     2,
				AvailableCount: 1,
				AdditionalFee:  15.00,
			},
		}

		for _, spot := range specialSpots {
			DB.Create(&spot)
		}
	}
}

func createTestParkingSessions() {
	var count int64
	DB.Model(&ParkingSession{}).Count(&count)

	if count == 0 {
		// 获取用户、车辆和停车场数据
		var user User
		DB.Where("username = ?", "user").First(&user)

		var vehicle Vehicle
		DB.Where("user_id = ? AND is_default = ?", user.ID, true).First(&vehicle)

		var parkingLot ParkingLot
		DB.Where("name = ?", "西湖停车场").First(&parkingLot)

		// 创建一个活跃的停车会话
		now := time.Now()
		startTime := now.Add(-2 * time.Hour)         // 2小时前开始停车
		nextBillingTime := now.Add(58 * time.Minute) // 下次计费时间
		nextFeeAmount := 5.0

		activeSession := ParkingSession{
			UserID:                user.ID,
			VehicleID:             vehicle.ID,
			ParkingLotID:          parkingLot.ID,
			SpotCode:              "A-101",
			SpotType:              "normal",
			StartTime:             startTime,
			EndTime:               nil, // 未结束
			Status:                "active",
			FeeRate:               10.0,
			FeeCurrent:            10.0, // 已经停了2小时，费用10元
			NextBillingTime:       &nextBillingTime,
			NextFeeAmount:         &nextFeeAmount,
			CurrentBillingCycle:   1,
			PricingRule:           "首小时10元，后续每小时5元",
			NavigationStatus:      "parked",
			RemainingDistanceM:    0,
			EstimatedMinutes:      0,
			DestinationLat:        parkingLot.Latitude,
			DestinationLon:        parkingLot.Longitude,
			UserPositionLat:       &parkingLot.Latitude,
			UserPositionLon:       &parkingLot.Longitude,
			ProgressToDestination: 100,
		}

		DB.Create(&activeSession)
		log.Println("Test parking session created")
	}
}

// CreateSimulationData 创建模拟数据
func CreateSimulationData() {
	createTrafficData()
	createDeviceData()
	createStatisticsData()
	log.Println("Simulation data created successfully")
}

func createTrafficData() {
	// 创建交通流量数据
	var count int64
	DB.Model(&TrafficFlow{}).Count(&count)
	if count == 0 {
		flows := []TrafficFlow{
			{Location: "一层北", FlowCount: 1250, Speed: 45.5, Direction: "inbound", VehicleType: "小型汽车", Timestamp: time.Now().Add(-1 * time.Hour)},
			{Location: "一层东", FlowCount: 880, Speed: 35.2, Direction: "outbound", VehicleType: "小型汽车", Timestamp: time.Now().Add(-1 * time.Hour)},
			{Location: "二层南", FlowCount: 950, Speed: 42.1, Direction: "inbound", VehicleType: "小型汽车", Timestamp: time.Now().Add(-1 * time.Hour)},
		}
		for _, flow := range flows {
			DB.Create(&flow)
		}
	}

	// 创建交通用户统计数据
	DB.Model(&TrafficUserStats{}).Count(&count)
	if count == 0 {
		stats := []TrafficUserStats{
			{MotorCount: 65, NonMotorCount: 22, PedestrianCount: 18, Timestamp: time.Now().Add(-30 * time.Minute), Location: "市中心区域"},
			{MotorCount: 80, NonMotorCount: 30, PedestrianCount: 25, Timestamp: time.Now().Add(-25 * time.Minute), Location: "市中心区域"},
			{MotorCount: 75, NonMotorCount: 40, PedestrianCount: 20, Timestamp: time.Now().Add(-20 * time.Minute), Location: "市中心区域"},
		}
		for _, stat := range stats {
			DB.Create(&stat)
		}
	}

	// 创建交通热力图数据
	DB.Model(&TrafficHeatmap{}).Count(&count)
	if count == 0 {
		roads := []string{"中山北路", "南京路", "延安高架", "内环路", "外环路", "浦东大道", "世纪大道", "陆家嘴环路", "北京路", "四川路"}
		for i, road := range roads {
			heatmap := TrafficHeatmap{
				RoadName:        road,
				CongestionLevel: 0.3 + float64(i%5)*0.15, // 0.3-0.9之间变化
				GridX:           i / 5,
				GridY:           i % 5,
				TimeFilter:      "now",
				Timestamp:       time.Now(),
			}
			DB.Create(&heatmap)
		}
	}

	// 创建拥堵播报数据
	DB.Model(&CongestionReport{}).Count(&count)
	if count == 0 {
		reports := []CongestionReport{
			{Location: "延安高架路中段", Severity: "严重", Description: "因施工作业导致车道缩减，车辆通行缓慢", Duration: 120, ReportTime: time.Now().Add(-2 * time.Hour), Status: "active"},
			{Location: "南京路与淮海路交叉口", Severity: "轻微", Description: "红绿灯故障，交警现场指挥交通", Duration: 45, ReportTime: time.Now().Add(-1 * time.Hour), Status: "active"},
		}
		for _, report := range reports {
			DB.Create(&report)
		}
	}
}

func createDeviceData() {
	// 创建设备数据
	var count int64
	DB.Model(&Device{}).Count(&count)
	if count == 0 {
		devices := []Device{
			{DeviceName: "摄像头-001", DeviceType: "监控设备", Location: "一层北入口", Status: "normal", SerialNumber: "CAM001", Manufacturer: "海康威视", InstallDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local)},
			{DeviceName: "感应器-001", DeviceType: "车位检测", Location: "二层A区", Status: "normal", SerialNumber: "SEN001", Manufacturer: "博世", InstallDate: time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local)},
			{DeviceName: "充电桩-001", DeviceType: "充电设备", Location: "地下一层", Status: "normal", SerialNumber: "CHG001", Manufacturer: "特来电", InstallDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.Local)},
		}
		for _, device := range devices {
			DB.Create(&device)
		}

		// 创建设备支出数据
		expenses := []DeviceExpense{
			{DeviceID: 1, ExpenseType: "维护费用", Amount: 500.00, Description: "摄像头清洁和调试", ExpenseDate: time.Date(2025, 8, 1, 0, 0, 0, 0, time.Local), Period: "month"},
			{DeviceID: 2, ExpenseType: "电费", Amount: 120.00, Description: "车位检测器电费", ExpenseDate: time.Date(2025, 8, 5, 0, 0, 0, 0, time.Local), Period: "month"},
			{DeviceID: 3, ExpenseType: "电费", Amount: 2500.00, Description: "充电桩电费", ExpenseDate: time.Date(2025, 8, 1, 0, 0, 0, 0, time.Local), Period: "month"},
		}
		for _, expense := range expenses {
			DB.Create(&expense)
		}

		// 创建设备故障统计
		faultStats := []DeviceFaultStats{
			{DeviceType: "监控设备", FaultCount: 2, RepairCount: 2, FaultRate: 12.50, StatDate: time.Now(), Period: "daily", TrendDirection: "down", PercentageChange: -5.2},
			{DeviceType: "车位检测", FaultCount: 3, RepairCount: 2, FaultRate: 25.00, StatDate: time.Now(), Period: "daily", TrendDirection: "up", PercentageChange: 8.3},
			{DeviceType: "充电设备", FaultCount: 1, RepairCount: 1, FaultRate: 25.00, StatDate: time.Now(), Period: "daily", TrendDirection: "down", PercentageChange: -10.5},
		}
		for _, stat := range faultStats {
			DB.Create(&stat)
		}
	}
}

func createStatisticsData() {
	// 创建停车饱和度数据
	var count int64
	DB.Model(&ParkingSaturation{}).Count(&count)
	if count == 0 {
		// 获取停车场数据
		var parkingLots []ParkingLot
		DB.Find(&parkingLots)

		for _, lot := range parkingLots {
			saturation := ParkingSaturation{
				ParkingLotID:   lot.ID,
				SaturationRate: float64(lot.TotalSpots-lot.AvailableSpots) / float64(lot.TotalSpots) * 100,
				OccupiedSpots:  lot.TotalSpots - lot.AvailableSpots,
				TotalSpots:     lot.TotalSpots,
				Timestamp:      time.Now().Add(-1 * time.Hour),
				Hour:           8,
			}
			DB.Create(&saturation)
		}
	}

	// 创建总占用率数据
	DB.Model(&TotalOccupancy{}).Count(&count)
	if count == 0 {
		totalOccupancy := TotalOccupancy{
			TotalSpots:     700,
			OccupiedSpots:  548,
			OccupancyRate:  78.29,
			AvailableSpots: 152,
			Timestamp:      time.Now().Add(-1 * time.Hour),
			Period:         "hourly",
		}
		DB.Create(&totalOccupancy)
	}

	// 创建监控摄像头数据
	DB.Model(&MonitoringCamera{}).Count(&count)
	if count == 0 {
		cameras := []MonitoringCamera{
			{CameraName: "入口监控-001", Location: "主入口A", Latitude: 30.2594, Longitude: 120.1644, Status: "online", StreamUrl: "rtmp://192.168.1.101/live/cam001", Type: "球机", Resolution: "1080P", ViewAngle: 360, NightVision: true, InstallDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.Local)},
			{CameraName: "出口监控-002", Location: "主出口B", Latitude: 30.2595, Longitude: 120.1645, Status: "online", StreamUrl: "rtmp://192.168.1.102/live/cam002", Type: "枪机", Resolution: "4K", ViewAngle: 90, NightVision: true, InstallDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.Local)},
		}
		for _, camera := range cameras {
			DB.Create(&camera)
		}
	}
}
