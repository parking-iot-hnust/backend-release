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

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
)

// GetDeviceExpenses 获取设备支出数据
func GetDeviceExpenses(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	var expenses []models.DeviceExpense
	result := models.DB.Where("period = ?", period).
		Order("expense_date desc").
		Find(&expenses)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 按支出类型分组统计
	expenseMap := make(map[string]float64)
	for _, expense := range expenses {
		expenseMap[expense.ExpenseType] += expense.Amount
	}

	// 转换为前端需要的格式
	var data []gin.H
	for expenseType, amount := range expenseMap {
		data = append(data, gin.H{
			"name":  expenseType,
			"value": amount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"message": "获取设备支出数据成功",
	})
}

// GetDeviceMaintenanceRecords 获取设备维修记录
func GetDeviceMaintenanceRecords(c *gin.Context) {
	var records []models.DeviceMaintenanceRecord
	result := models.DB.Preload("Device").
		Order("start_time desc").
		Limit(100).
		Find(&records)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    records,
		"message": "获取设备维修记录成功",
	})
}

// GetDeviceFaultStats 获取设备故障统计
func GetDeviceFaultStats(c *gin.Context) {
	var stats []models.DeviceFaultStats
	result := models.DB.Order("stat_date desc").
		Limit(10).
		Find(&stats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 按设备类型分组
	statsMap := make(map[string]*gin.H)
	for _, stat := range stats {
		if _, exists := statsMap[stat.DeviceType]; !exists {
			statsMap[stat.DeviceType] = &gin.H{
				"device_type":       stat.DeviceType,
				"fault_count":       stat.FaultCount,
				"repair_count":      stat.RepairCount,
				"fault_rate":        stat.FaultRate,
				"trend_direction":   stat.TrendDirection,
				"percentage_change": stat.PercentageChange,
			}
		}
	}

	var data []gin.H
	for _, statData := range statsMap {
		data = append(data, *statData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"message": "获取设备故障统计成功",
	})
}

// GetDeviceAlarms 获取设备报警
func GetDeviceAlarms(c *gin.Context) {
	var alarms []models.DeviceAlarm
	result := models.DB.Preload("Device").
		Where("status = ?", "active").
		Order("alarm_time desc").
		Limit(50).
		Find(&alarms)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    alarms,
		"message": "获取设备报警数据成功",
	})
}

// GetDeviceAlarmStats 获取设备报警统计
func GetDeviceAlarmStats(c *gin.Context) {
	var alarmStats []models.DeviceAlarmStats
	result := models.DB.Order("stat_date desc, hour").
		Limit(24). // 最近24小时
		Find(&alarmStats)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 按报警类型分组统计
	typeStats := make(map[string]int)
	hourlyStats := make(map[int]int)

	for _, stat := range alarmStats {
		typeStats[stat.AlarmType] += stat.Count
		hourlyStats[stat.Hour] += stat.Count
	}

	response := gin.H{
		"type_distribution":   typeStats,
		"hourly_distribution": hourlyStats,
		"total_alarms": func() int {
			total := 0
			for _, count := range typeStats {
				total += count
			}
			return total
		}(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
		"message": "获取设备报警统计成功",
	})
}

// GetDeviceFaultTrend 获取设备故障趋势
func GetDeviceFaultTrend(c *gin.Context) {
	var trends []models.DeviceFaultStats
	result := models.DB.Where("period = ?", "daily").
		Order("stat_date desc").
		Limit(30). // 最近30天
		Find(&trends)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	// 按日期和设备类型组织数据
	trendData := make(map[string][]gin.H)

	for _, trend := range trends {
		if _, exists := trendData[trend.DeviceType]; !exists {
			trendData[trend.DeviceType] = []gin.H{}
		}

		trendData[trend.DeviceType] = append(trendData[trend.DeviceType], gin.H{
			"date":              trend.StatDate.Format("2006-01-02"),
			"fault_rate":        trend.FaultRate,
			"fault_count":       trend.FaultCount,
			"percentage_change": trend.PercentageChange,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    trendData,
		"message": "获取设备故障趋势成功",
	})
}
