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

// GetCurrentAirQuality 获取当前空气质量数据
func GetCurrentAirQuality(c *gin.Context) {
	location := c.DefaultQuery("location", "停车场周边")

	var airQuality models.AirQuality
	result := models.DB.Where("location = ?", location).
		Order("timestamp desc").
		First(&airQuality)

	if result.Error != nil {
		// 如果没有数据，返回默认值
		airQuality = models.AirQuality{
			Location:  location,
			AQI:       75,
			Level:     "良",
			PM25:      28.0,
			PM10:      55.0,
			O3:        90.0,
			NO2:       30.0,
			SO2:       10.0,
			CO:        0.8,
			Timestamp: time.Now(),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    airQuality,
		"message": "获取空气质量数据成功",
	})
}

// GetAirQualityHistory 获取空气质量历史数据
func GetAirQualityHistory(c *gin.Context) {
	location := c.DefaultQuery("location", "停车场周边")
	hours := c.DefaultQuery("hours", "24")

	var airQualityData []models.AirQuality

	// 计算时间范围
	var timeRange time.Time
	switch hours {
	case "1":
		timeRange = time.Now().Add(-1 * time.Hour)
	case "6":
		timeRange = time.Now().Add(-6 * time.Hour)
	case "12":
		timeRange = time.Now().Add(-12 * time.Hour)
	case "24":
		timeRange = time.Now().Add(-24 * time.Hour)
	case "48":
		timeRange = time.Now().Add(-48 * time.Hour)
	default:
		timeRange = time.Now().Add(-24 * time.Hour)
	}

	result := models.DB.Where("location = ? AND timestamp >= ?", location, timeRange).
		Order("timestamp desc").
		Limit(50).
		Find(&airQualityData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取空气质量历史数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    airQualityData,
		"message": "获取空气质量历史数据成功",
	})
}

// GetAirQualityStats 获取空气质量统计数据
func GetAirQualityStats(c *gin.Context) {
	days := c.DefaultQuery("days", "7")

	var stats []models.AirQualityStats

	// 计算日期范围
	var dateRange time.Time
	switch days {
	case "3":
		dateRange = time.Now().AddDate(0, 0, -3)
	case "7":
		dateRange = time.Now().AddDate(0, 0, -7)
	case "15":
		dateRange = time.Now().AddDate(0, 0, -15)
	case "30":
		dateRange = time.Now().AddDate(0, 0, -30)
	default:
		dateRange = time.Now().AddDate(0, 0, -7)
	}

	result := models.DB.Where("date >= ?", dateRange).
		Order("date desc").
		Find(&stats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取空气质量统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "获取空气质量统计数据成功",
	})
}

// UpdateAirQuality 更新空气质量数据（供数据采集使用）
func UpdateAirQuality(c *gin.Context) {
	var airQuality models.AirQuality
	if err := c.ShouldBindJSON(&airQuality); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 设置时间戳
	airQuality.Timestamp = time.Now()

	// 根据AQI计算等级
	if airQuality.AQI <= 50 {
		airQuality.Level = "优"
	} else if airQuality.AQI <= 100 {
		airQuality.Level = "良"
	} else if airQuality.AQI <= 150 {
		airQuality.Level = "轻度污染"
	} else if airQuality.AQI <= 200 {
		airQuality.Level = "中度污染"
	} else if airQuality.AQI <= 300 {
		airQuality.Level = "重度污染"
	} else {
		airQuality.Level = "严重污染"
	}

	if err := models.DB.Create(&airQuality).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新空气质量数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    airQuality,
		"message": "空气质量数据更新成功",
	})
}
