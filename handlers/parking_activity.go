/*
 * Copyright (c) 2025 LTQY. All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package handlers

import (
	"net/http"
	"time"

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
)

// ParkingActivityData 停车活动数据结构
type ParkingActivityData struct {
	Month string         `json:"month"`
	Data  map[string]int `json:"data"`
}

// ParkingActivityResponse 停车活动分析响应
type ParkingActivityResponse struct {
	Success    bool                  `json:"success"`
	Data       []ParkingActivityData `json:"data"`
	UpdateTime string                `json:"update_time"`
}

// GetParkingActivityAnalysis 获取停车活动类型分析数据
func GetParkingActivityAnalysis(c *gin.Context) {
	currentMonth := time.Now().Format("2006-01")
	previousMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

	// 从数据库获取数据
	currentData, err := calculateParkingActivityFromDB(currentMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取当前月份数据失败: " + err.Error(),
		})
		return
	}

	previousData, err := calculateParkingActivityFromDB(previousMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取上月数据失败: " + err.Error(),
		})
		return
	}

	response := ParkingActivityResponse{
		Success: true,
		Data: []ParkingActivityData{
			{
				Month: "current",
				Data:  currentData,
			},
			{
				Month: "previous",
				Data:  previousData,
			},
		},
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// calculateParkingActivityFromDB 从数据库查询停车活动数据
func calculateParkingActivityFromDB(month string) (map[string]int, error) {
	activityData := map[string]int{
		"临时停车":   0,
		"月租车辆":   0,
		"新能源车":   0,
		"访客车辆":   0,
		"高峰入场":   0,
		"夜间停车":   0,
		"特殊需求":   0,
		"短时快进快出": 0,
	}

	// 查询指定月份的所有停车会话
	var sessions []models.ParkingSession
	err := models.DB.Preload("Vehicle").Where(
		"DATE_FORMAT(start_time, '%Y-%m') = ? AND status = ?",
		month, "ended",
	).Find(&sessions).Error

	if err != nil {
		return nil, err
	}

	// 统计各类停车活动
	for _, session := range sessions {
		duration := session.EndTime.Sub(session.StartTime).Hours()
		hour := session.StartTime.Hour()

		// 根据车位类型分类
		switch session.SpotType {
		case "charging":
			activityData["新能源车"]++
		case "disabled", "vip":
			activityData["特殊需求"]++
		default:
			// 根据停车时长和时间段分类
			if duration <= 2 { // 2小时以内算短时停车
				activityData["短时快进快出"]++
			} else if duration >= 8 { // 8小时以上算夜间停车或长期停车
				if hour >= 18 || hour <= 8 {
					activityData["夜间停车"]++
				} else {
					// 可能是月租车辆（长期停车）
					activityData["月租车辆"]++
				}
			} else {
				// 中等时长停车
				if hour >= 7 && hour <= 9 || hour >= 17 && hour <= 19 {
					activityData["高峰入场"]++
				} else if session.Vehicle.PlateNumber != "" {
					// 判断是否为访客车辆（这里简化处理，实际可根据车辆注册信息判断）
					if duration <= 4 {
						activityData["访客车辆"]++
					} else {
						activityData["临时停车"]++
					}
				} else {
					activityData["临时停车"]++
				}
			}
		}
	}

	// 将统计数据转换为活跃度指数（0-100）
	maxCount := 1
	for _, count := range activityData {
		if count > maxCount {
			maxCount = count
		}
	}

	// 如果没有数据，使用默认值避免除零错误
	if maxCount == 0 {
		maxCount = 1
	}

	// 转换为百分比指数，并确保有一定的基础值
	for key, count := range activityData {
		index := int(float64(count)/float64(maxCount)*60) + 20 // 20-80的范围
		if index > 100 {
			index = 100
		}
		activityData[key] = index
	}

	return activityData, nil
}

// GetParkingActivityRealtime 获取实时停车活动数据
func GetParkingActivityRealtime(c *gin.Context) {
	// 获取实时更新的停车活动数据
	currentData, err := calculateRealtimeParkingActivityFromDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "获取实时数据失败: " + err.Error(),
		})
		return
	}

	response := gin.H{
		"success": true,
		"data": map[string]interface{}{
			"month": "current",
			"data":  currentData,
		},
		"update_time": time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// calculateRealtimeParkingActivityFromDB 从数据库查询实时停车活动数据
func calculateRealtimeParkingActivityFromDB() (map[string]int, error) {
	activityData := map[string]int{
		"临时停车":   0,
		"月租车辆":   0,
		"新能源车":   0,
		"访客车辆":   0,
		"高峰入场":   0,
		"夜间停车":   0,
		"特殊需求":   0,
		"短时快进快出": 0,
	}

	currentTime := time.Now()
	hour := currentTime.Hour()

	// 查询当前活跃的停车会话（正在进行中的）
	var activeSessions []models.ParkingSession
	err := models.DB.Preload("Vehicle").Where("status = ?", "active").Find(&activeSessions).Error
	if err != nil {
		return nil, err
	}

	// 查询最近24小时内结束的停车会话（用于分析当前趋势）
	yesterday := currentTime.Add(-24 * time.Hour)
	var recentSessions []models.ParkingSession
	err = models.DB.Preload("Vehicle").Where(
		"status = ? AND end_time >= ?",
		"ended", yesterday,
	).Find(&recentSessions).Error
	if err != nil {
		return nil, err
	}

	// 合并活跃会话和最近会话进行分析
	allSessions := append(activeSessions, recentSessions...)

	// 统计各类停车活动
	for _, session := range allSessions {
		var duration float64
		if session.EndTime != nil {
			duration = session.EndTime.Sub(session.StartTime).Hours()
		} else {
			// 对于活跃会话，计算当前已停车时长
			duration = currentTime.Sub(session.StartTime).Hours()
		}

		sessionHour := session.StartTime.Hour()

		// 根据车位类型和停车特征分类
		switch session.SpotType {
		case "charging":
			activityData["新能源车"]++
		case "disabled", "vip":
			activityData["特殊需求"]++
		default:
			// 基于时间和停车模式的智能分类
			if duration <= 1.5 { // 1.5小时以内算快进快出
				activityData["短时快进快出"]++
			} else if duration >= 10 { // 10小时以上算夜间或月租
				if sessionHour >= 18 || sessionHour <= 8 {
					activityData["夜间停车"]++
				} else {
					activityData["月租车辆"]++
				}
			} else {
				// 中等时长的停车，根据时间段判断
				if sessionHour >= 7 && sessionHour <= 9 || sessionHour >= 17 && sessionHour <= 19 {
					activityData["高峰入场"]++
				} else if duration <= 4 {
					// 短期访问
					activityData["访客车辆"]++
				} else {
					// 常规临时停车
					activityData["临时停车"]++
				}
			}
		}
	}

	// 根据当前时间段调整权重，体现实时特征
	timeFactors := calculateTimeFactors(hour)
	for key, factor := range timeFactors {
		if count, exists := activityData[key]; exists {
			activityData[key] = int(float64(count) * factor)
		}
	}

	// 转换为活跃度指数
	maxCount := 1
	for _, count := range activityData {
		if count > maxCount {
			maxCount = count
		}
	}

	// 确保有基础活跃度，避免数据过于稀疏
	if maxCount == 0 {
		maxCount = 1
	}

	// 转换为指数（30-95范围，确保有一定的活跃度基础）
	for key, count := range activityData {
		index := int(float64(count)/float64(maxCount)*65) + 30
		if index > 100 {
			index = 100
		}
		activityData[key] = index
	}

	return activityData, nil
}

// calculateTimeFactors 根据当前时间计算各活动类型的权重因子
func calculateTimeFactors(hour int) map[string]float64 {
	factors := map[string]float64{
		"临时停车":   1.0,
		"月租车辆":   1.0,
		"新能源车":   1.0,
		"访客车辆":   1.0,
		"高峰入场":   1.0,
		"夜间停车":   1.0,
		"特殊需求":   1.0,
		"短时快进快出": 1.0,
	}

	if hour >= 7 && hour <= 9 || hour >= 17 && hour <= 19 {
		// 高峰时段权重调整
		factors["高峰入场"] = 1.5
		factors["短时快进快出"] = 1.3
		factors["临时停车"] = 1.2
		factors["夜间停车"] = 0.3
	} else if hour >= 22 || hour <= 6 {
		// 夜间时段权重调整
		factors["夜间停车"] = 1.8
		factors["月租车辆"] = 1.2
		factors["高峰入场"] = 0.2
		factors["短时快进快出"] = 0.4
		factors["访客车辆"] = 0.5
	} else if hour >= 10 && hour <= 16 {
		// 平峰时段权重调整
		factors["访客车辆"] = 1.4
		factors["临时停车"] = 1.2
		factors["新能源车"] = 1.1
		factors["高峰入场"] = 0.6
	}

	return factors
}
