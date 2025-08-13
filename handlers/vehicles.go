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
	"strconv"
	"time"

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
)

type CreateVehicleRequest struct {
	PlateNumber string `json:"plate_number" binding:"required"`
	Brand       string `json:"brand" binding:"required"`
	Model       string `json:"model" binding:"required"`
	Color       string `json:"color"`
	Type        string `json:"type"`
	RegDate     string `json:"reg_date"`
}

// GetVehicles 获取用户的车辆列表
func GetVehicles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	var vehicles []models.Vehicle
	if err := models.DB.Where("user_id = ?", userID).Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取车辆列表失败"})
		return
	}

	// 为每辆车添加统计信息
	var vehiclesWithStats []models.VehicleWithStats
	for _, vehicle := range vehicles {
		vehicleWithStats := models.VehicleWithStats{Vehicle: vehicle}

		// 获取停车统计
		var stats struct {
			ParkingCount int     `json:"parking_count"`
			TotalFee     float64 `json:"total_fee"`
			TotalHours   float64 `json:"total_hours"`
		}

		models.DB.Model(&models.ParkingRecord{}).
			Where("vehicle_id = ? AND end_time != ''", vehicle.ID).
			Select("COUNT(*) as parking_count, COALESCE(SUM(fee), 0) as total_fee, COALESCE(SUM(duration), 0) as total_hours").
			Scan(&stats)

		vehicleWithStats.Stats = stats

		// 获取最后一次停车记录
		var lastRecord models.ParkingRecord
		if err := models.DB.Where("vehicle_id = ?", vehicle.ID).
			Order("created_at DESC").First(&lastRecord).Error; err == nil {
			vehicleWithStats.LastParking = &struct {
				Location string `json:"location"`
				Time     string `json:"time"`
			}{
				Location: lastRecord.Location,
				Time:     lastRecord.StartTime,
			}
		}

		vehiclesWithStats = append(vehiclesWithStats, vehicleWithStats)
	}

	c.JSON(http.StatusOK, gin.H{
		"vehicles": vehiclesWithStats,
		"count":    len(vehiclesWithStats),
	})
}

// CreateVehicle 添加新车辆
func CreateVehicle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	var req CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查车牌号是否已存在
	var existingVehicle models.Vehicle
	if err := models.DB.Where("plate_number = ?", req.PlateNumber).First(&existingVehicle).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该车牌号已存在"})
		return
	}

	// 检查是否是用户的第一辆车
	var count int64
	models.DB.Model(&models.Vehicle{}).Where("user_id = ?", userID).Count(&count)
	isDefault := count == 0

	// 创建新车辆
	vehicle := models.Vehicle{
		UserID:      userID.(uint),
		PlateNumber: req.PlateNumber,
		Brand:       req.Brand,
		Model:       req.Model,
		Color:       req.Color,
		Type:        req.Type,
		RegDate:     req.RegDate,
		IsDefault:   isDefault,
	}

	if vehicle.Type == "" {
		vehicle.Type = "小型汽车"
	}
	if vehicle.RegDate == "" {
		vehicle.RegDate = time.Now().Format("2006-01-02")
	}

	if err := models.DB.Create(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加车辆失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "车辆添加成功",
		"vehicle": vehicle,
	})
}

// SetDefaultVehicle 设置默认车辆
func SetDefaultVehicle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	vehicleID := c.Param("id")
	id, err := strconv.ParseUint(vehicleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的车辆ID"})
		return
	}

	// 验证车辆是否属于当前用户
	var vehicle models.Vehicle
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "车辆不存在"})
		return
	}

	// 取消其他车辆的默认状态
	models.DB.Model(&models.Vehicle{}).Where("user_id = ?", userID).Update("is_default", false)

	// 设置当前车辆为默认
	if err := models.DB.Model(&vehicle).Update("is_default", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置默认车辆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "默认车辆设置成功"})
}

// DeleteVehicle 删除车辆
func DeleteVehicle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	vehicleID := c.Param("id")
	id, err := strconv.ParseUint(vehicleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的车辆ID"})
		return
	}

	// 验证车辆是否属于当前用户
	var vehicle models.Vehicle
	if err := models.DB.Where("id = ? AND user_id = ?", id, userID).First(&vehicle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "车辆不存在"})
		return
	}

	// 删除相关的停车记录
	models.DB.Where("vehicle_id = ?", id).Delete(&models.ParkingRecord{})

	// 删除车辆
	if err := models.DB.Delete(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除车辆失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "车辆删除成功"})
}
