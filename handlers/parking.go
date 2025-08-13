/*
 * Copyright (c) 2025 LTQY. All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package handlers

import (
	"math"
	"net/http"
	"sort"
	"strconv"

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
)

// 计算两点间距离
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // 地球半径（公里）

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// GetNearbyParkingLots 获取附近停车场
func GetNearbyParkingLots(c *gin.Context) {
	// 获取用户当前位置（从查询参数）
	latStr := c.DefaultQuery("lat", "30.2594")
	lonStr := c.DefaultQuery("lon", "120.1644")
	sortBy := c.DefaultQuery("sort", "distance")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	// 获取所有活跃的停车场
	var parkingLots []models.ParkingLot
	if err := models.DB.Preload("SpecialSpots").Where("is_active = ?", true).Find(&parkingLots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parking lots"})
		return
	}

	// 计算距离并构建响应数据
	var parkingLotsWithDistance []models.ParkingLotWithDistance
	for _, lot := range parkingLots {
		distance := calculateDistance(lat, lon, lot.Latitude, lot.Longitude)
		distanceStr := ""
		if distance < 1 {
			distanceStr = strconv.FormatFloat(distance*1000, 'f', 0, 64) + "m"
		} else {
			distanceStr = strconv.FormatFloat(distance, 'f', 1, 64) + "km"
		}

		// 构建特殊车位映射
		specialSpotsMap := make(map[string]models.SpecialSpot)
		for _, spot := range lot.SpecialSpots {
			specialSpotsMap[spot.SpotType] = spot
		}

		parkingLotsWithDistance = append(parkingLotsWithDistance, models.ParkingLotWithDistance{
			ParkingLot:      lot,
			Distance:        distanceStr,
			SpecialSpotsMap: specialSpotsMap,
		})
	}

	// 根据排序条件排序
	switch sortBy {
	case "available":
		sort.Slice(parkingLotsWithDistance, func(i, j int) bool {
			return parkingLotsWithDistance[i].AvailableSpots > parkingLotsWithDistance[j].AvailableSpots
		})
	case "rate":
		sort.Slice(parkingLotsWithDistance, func(i, j int) bool {
			rateI := float64(parkingLotsWithDistance[i].AvailableSpots) / float64(parkingLotsWithDistance[i].TotalSpots)
			rateJ := float64(parkingLotsWithDistance[j].AvailableSpots) / float64(parkingLotsWithDistance[j].TotalSpots)
			return rateI > rateJ
		})
	default: // distance
		sort.Slice(parkingLotsWithDistance, func(i, j int) bool {
			distI := calculateDistance(lat, lon, parkingLotsWithDistance[i].Latitude, parkingLotsWithDistance[i].Longitude)
			distJ := calculateDistance(lat, lon, parkingLotsWithDistance[j].Latitude, parkingLotsWithDistance[j].Longitude)
			return distI < distJ
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        parkingLotsWithDistance,
		"update_time": "",
		"user_location": gin.H{
			"latitude":  lat,
			"longitude": lon,
		},
	})
}

// UpdateParkingLotAvailability 更新停车场可用车位数
func UpdateParkingLotAvailability(c *gin.Context) {
	lotID := c.Param("id")

	var lot models.ParkingLot
	if err := models.DB.First(&lot, lotID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking lot not found"})
		return
	}

	var updateData struct {
		AvailableSpots int `json:"available_spots"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// 验证可用车位数不能超过总车位数
	if updateData.AvailableSpots > lot.TotalSpots || updateData.AvailableSpots < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid available spots count"})
		return
	}

	// 更新可用车位数
	if err := models.DB.Model(&lot).Update("available_spots", updateData.AvailableSpots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parking lot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Parking lot availability updated successfully",
	})
}

// GetParkingLotDetails 获取停车场详细信息
func GetParkingLotDetails(c *gin.Context) {
	lotID := c.Param("id")

	var lot models.ParkingLot
	if err := models.DB.Preload("SpecialSpots").First(&lot, lotID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parking lot not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    lot,
	})
}
