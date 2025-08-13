/*
 * Copyright (c) 2025 LTQY. All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
)

// TrafficFlowData 交通流量数据结构
type TrafficFlowData struct {
	ID       int    `json:"id"`
	Location string `json:"location"`
	Flow     int    `json:"flow"`
	Unit     string `json:"unit"`
	Trend    string `json:"trend"`
}

// RealTimeTrafficData 实时交通数据结构
type RealTimeTrafficData struct {
	ID     int    `json:"id"`
	Road   string `json:"road"`
	Status string `json:"status"`
	Speed  int    `json:"speed"`
	Change int    `json:"change"`
}

// InOutFlowData 进出流量数据结构
type InOutFlowData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

// ParkingStatsData 停车统计数据结构
type ParkingStatsData struct {
	TotalSpots     int     `json:"total_spots"`
	OccupiedSpots  int     `json:"occupied_spots"`
	AvailableSpots int     `json:"available_spots"`
	OccupancyRate  float64 `json:"occupancy_rate"`
	LastUpdate     string  `json:"last_update"`
}

// GetTrafficFlow 获取交通流量数据
func GetTrafficFlow(c *gin.Context) {
	// 从数据库获取交通流量数据
	var trafficData []models.TrafficFlow
	result := models.DB.Where("timestamp >= ?", time.Now().Add(-1*time.Hour)).
		Order("location").
		Find(&trafficData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取交通流量数据失败"})
		return
	}

	// 转换为前端需要的格式
	var flowData []TrafficFlowData
	for _, data := range trafficData {
		// 根据流量确定趋势
		trend := "stable"
		if data.FlowCount > 1000 {
			trend = "up"
		} else if data.FlowCount < 500 {
			trend = "down"
		}

		flowData = append(flowData, TrafficFlowData{
			ID:       int(data.ID),
			Location: data.Location,
			Flow:     data.FlowCount,
			Unit:     "辆/小时",
			Trend:    trend,
		})
	}

	// 如果没有数据，提供默认响应
	if len(flowData) == 0 {
		flowData = []TrafficFlowData{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    flowData,
		"message": "获取交通流量数据成功",
	})
}

// GetRealTimeTraffic 获取实时交通数据
func GetRealTimeTraffic(c *gin.Context) {
	// 从数据库获取实时交通数据
	var trafficData []models.TrafficFlow
	result := models.DB.Where("timestamp >= ?", time.Now().Add(-15*time.Minute)).
		Order("timestamp desc, location").
		Find(&trafficData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取实时交通数据失败"})
		return
	}

	// 转换为前端需要的格式
	var realTimeData []RealTimeTrafficData
	for _, data := range trafficData {
		// 根据速度确定状态
		status := "畅通"
		if data.Speed < 20 {
			status = "拥堵"
		} else if data.Speed < 40 {
			status = "缓行"
		}

		// 计算变化（这里简化处理，实际应该对比历史数据）
		change := 0
		if data.Speed < 20 {
			change = -5
		} else if data.Speed > 50 {
			change = 3
		}

		realTimeData = append(realTimeData, RealTimeTrafficData{
			ID:     int(data.ID),
			Road:   data.Location,
			Status: status,
			Speed:  int(data.Speed),
			Change: change,
		})
	}

	// 如果没有数据，提供默认响应
	if len(realTimeData) == 0 {
		realTimeData = []RealTimeTrafficData{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    realTimeData,
		"message": "获取实时交通数据成功",
	})
}

// GetInOutFlowData 获取进出流量数据
func GetInOutFlowData(c *gin.Context) {
	// 从数据库获取进出流量数据
	var trafficData []models.TrafficFlow
	result := models.DB.Select("location, direction, SUM(flow_count) as total_flow").
		Where("timestamp >= ? AND direction IN ('inbound', 'outbound')",
			time.Now().Add(-1*time.Hour)).
		Group("location, direction").
		Find(&trafficData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取进出流量数据失败"})
		return
	}

	// 构建不同颜色的调色板
	colors := []string{"#4CAF50", "#2196F3", "#FFC107", "#FF5722", "#00BCD4", "#9C27B0", "#E91E63", "#795548"}

	// 转换为前端需要的格式
	var inOutData []InOutFlowData
	totalFlow := 0

	// 计算总流量
	for _, data := range trafficData {
		totalFlow += data.FlowCount
	}

	// 计算百分比并构建数据
	for i, data := range trafficData {
		percentage := 0
		if totalFlow > 0 {
			percentage = (data.FlowCount * 100) / totalFlow
		}

		colorIndex := i % len(colors)

		inOutData = append(inOutData, InOutFlowData{
			Name:  data.Location + "(" + data.Direction + ")",
			Value: percentage,
			Color: colors[colorIndex],
		})
	}

	// 如果没有数据，提供默认响应
	if len(inOutData) == 0 {
		inOutData = []InOutFlowData{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    inOutData,
		"message": "获取进出流量数据成功",
	})
}

// GetParkingStats 获取停车统计数据
func GetParkingStats(c *gin.Context) {
	// 从数据库计算停车统计数据
	var totalSpots int64
	var occupiedSpots int64

	// 获取总停车位数
	result := models.DB.Model(&models.ParkingLot{}).
		Select("SUM(total_spots)").
		Scan(&totalSpots)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取停车位数据失败"})
		return
	}

	// 获取当前占用的停车位数（进行中的停车会话）
	result = models.DB.Model(&models.ParkingSession{}).
		Where("end_time IS NULL OR end_time = ''").
		Count(&occupiedSpots)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取停车占用数据失败"})
		return
	}

	availableSpots := totalSpots - occupiedSpots
	occupancyRate := 0.0
	if totalSpots > 0 {
		occupancyRate = float64(occupiedSpots) / float64(totalSpots) * 100
	}

	stats := ParkingStatsData{
		TotalSpots:     int(totalSpots),
		OccupiedSpots:  int(occupiedSpots),
		AvailableSpots: int(availableSpots),
		OccupancyRate:  occupancyRate,
		LastUpdate:     time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "获取停车统计数据成功",
	})
}

// GetCurrentParkingStatus 获取当前停车状态
func GetCurrentParkingStatus(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	// 查找用户当前的停车会话
	var parkingSession models.ParkingSession
	result := models.DB.Where("user_id = ? AND (end_time IS NULL OR end_time = '')", userID).
		Order("start_time desc").
		First(&parkingSession)

	if result.Error != nil {
		// 没有找到当前停车会话
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"is_parking": false,
			},
			"message": "用户当前未停车",
		})
		return
	}

	// 计算停车时长
	duration := time.Since(parkingSession.StartTime)

	// 计算当前费用（这里使用简单的计费逻辑）
	hours := duration.Hours()
	currentFee := hours * 5.0 // 假设每小时5元

	// 计算下次计费时间
	nextBillingHour := int(hours) + 1
	nextBilling := parkingSession.StartTime.Add(time.Duration(nextBillingHour) * time.Hour)
	remainingMinutes := int(time.Until(nextBilling).Minutes())

	// 计算当前计费周期进度
	billingProgress := int((duration.Minutes() - float64(int(hours)*60)) / 60.0 * 100)

	// 获取停车场信息以获取location
	var parkingLot models.ParkingLot
	models.DB.First(&parkingLot, parkingSession.ParkingLotID)

	parkingStatus := gin.H{
		"is_parking":        true,
		"location":          parkingLot.Name + " " + parkingSession.SpotCode,
		"start_time":        parkingSession.StartTime.Format("2006-01-02 15:04:05"),
		"duration":          fmt.Sprintf("%.0f小时%.0f分钟", duration.Hours(), duration.Minutes()-duration.Hours()*60),
		"current_fee":       currentFee,
		"billing_cycle":     int(hours) + 1,
		"next_billing":      nextBilling.Format("2006-01-02 15:04:05"),
		"next_fee":          5.00,
		"billing_progress":  billingProgress,
		"remaining_minutes": remainingMinutes,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    parkingStatus,
		"message": "获取当前停车状态成功",
	})
}

// GetTrafficUserStats 获取交通用户统计数据
func GetTrafficUserStats(c *gin.Context) {
	// 获取最近的交通用户统计数据
	var stats []models.TrafficUserStats
	result := models.DB.Where("timestamp >= ?", time.Now().Add(-2*time.Hour)).
		Order("timestamp desc").
		Limit(20).
		Find(&stats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 构建时间标签
	var labels []string
	var motorData []int
	var nonMotorData []int
	var pedestrianData []int

	if len(stats) > 0 {
		// 按时间间隔整理数据
		for i := len(stats) - 1; i >= 0; i-- { // 反向遍历以获得时间顺序
			stat := stats[i]
			labels = append(labels, stat.Timestamp.Format("15:04"))
			motorData = append(motorData, stat.MotorCount)
			nonMotorData = append(nonMotorData, stat.NonMotorCount)
			pedestrianData = append(pedestrianData, stat.PedestrianCount)
		}
	} else {
		labels = []string{}
		motorData = []int{}
		nonMotorData = []int{}
		pedestrianData = []int{}
	}

	// 格式化返回数据
	chartData := gin.H{
		"labels": labels,
		"datasets": []gin.H{
			{
				"label":       "机动车",
				"data":        motorData,
				"borderColor": "#3b82f6",
			},
			{
				"label":       "非机动车",
				"data":        nonMotorData,
				"borderColor": "#06b6d4",
			},
			{
				"label":       "行人",
				"data":        pedestrianData,
				"borderColor": "#10b981",
			},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    chartData,
		"message": "获取交通用户统计数据成功",
	})
}

// GetCarCrossingRate 获取车辆通过率数据
func GetCarCrossingRate(c *gin.Context) {
	var crossingData []models.CarCrossingRate
	result := models.DB.Order("timestamp desc").Limit(24).Find(&crossingData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    crossingData,
		"message": "获取车辆通过率数据成功",
	})
}

// GetTrafficFlowChart 获取交通流量图表数据
func GetTrafficFlowChart(c *gin.Context) {
	var flowData []models.TrafficFlow
	result := models.DB.Order("timestamp desc").Limit(50).Find(&flowData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 按小时分组统计
	hourlyFlow := make(map[int]int)
	for _, flow := range flowData {
		hour := flow.Timestamp.Hour()
		hourlyFlow[hour] += flow.FlowCount
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    hourlyFlow,
		"message": "获取交通流量图表数据成功",
	})
}

// GetTrafficHeatmap 获取交通热力图数据
func GetTrafficHeatmap(c *gin.Context) {
	timeFilter := c.DefaultQuery("time", "now")

	var heatmapData []models.TrafficHeatmap
	query := models.DB.Where("time_filter = ?", timeFilter)

	// 如果是实时数据，获取最近的数据
	if timeFilter == "now" {
		query = query.Where("timestamp >= ?", time.Now().Add(-10*time.Minute))
	}

	result := query.Order("grid_x, grid_y").Find(&heatmapData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 转换为前端需要的格式
	gridData := make([][]*gin.H, 5)
	for i := range gridData {
		gridData[i] = make([]*gin.H, 10)
	}

	for _, data := range heatmapData {
		if data.GridX < 5 && data.GridY < 10 {
			gridData[data.GridX][data.GridY] = &gin.H{
				"value": data.CongestionLevel,
				"label": data.RoadName,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    gridData,
		"message": "获取交通热力图数据成功",
	})
}

// GetCongestionReports 获取拥堵播报
func GetCongestionReports(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "5")
	limit, _ := strconv.Atoi(limitStr)

	var reports []models.CongestionReport
	result := models.DB.Where("status = ?", "active").
		Order("report_time desc").
		Limit(limit).
		Find(&reports)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    reports,
		"message": "获取拥堵播报数据成功",
	})
}
