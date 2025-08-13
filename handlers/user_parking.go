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
	"strconv"
	"time"

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
)

// ParkingSessionResponse 停车会话响应结构
type ParkingSessionResponse struct {
	ID                            uint           `json:"id"`
	VehiclePlate                  string         `json:"vehicle_plate"`
	ParkingLot                    ParkingLotInfo `json:"parking_lot"`
	SpotCode                      string         `json:"spot_code"`
	SpotType                      string         `json:"spot_type"`
	StartTime                     string         `json:"start_time"`
	DurationMinutes               int            `json:"duration_minutes"`
	FeeCurrent                    float64        `json:"fee_current"`
	NextBillingTime               *string        `json:"next_billing_time"`
	NextFeeAmount                 *float64       `json:"next_fee_amount"`
	BillingProgressPercent        int            `json:"billing_progress_percent"`
	RemainingMinutesToNextBilling int            `json:"remaining_minutes_to_next_billing"`
	CurrentBillingCycle           int            `json:"current_billing_cycle"`
	PricingRule                   string         `json:"pricing_rule"`
	Navigation                    NavigationInfo `json:"navigation"`
	Status                        string         `json:"status"`
}

type ParkingLotInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Floor string `json:"floor,omitempty"`
	Area  string `json:"area,omitempty"`
}

type NavigationInfo struct {
	Status                       string    `json:"status"`
	RemainingDistanceM           int       `json:"remaining_distance_m"`
	EstimatedMinutes             int       `json:"estimated_minutes"`
	Destination                  Position  `json:"destination"`
	UserPosition                 *Position `json:"user_position,omitempty"`
	ProgressToDestinationPercent int       `json:"progress_to_destination_percent,omitempty"`
}

type Position struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type ParkingSessionHistoryItem struct {
	ID              uint    `json:"id"`
	ParkingLotName  string  `json:"parking_lot_name"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	TotalFee        float64 `json:"total_fee"`
	DurationMinutes int     `json:"duration_minutes"`
}

// GetCurrentParkingSession 获取当前用户的活跃停车会话
func GetCurrentParkingSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var session models.ParkingSession
	result := models.DB.Preload("Vehicle").Preload("ParkingLot").Where("user_id = ? AND status = ?", userID, "active").First(&session)

	if result.Error != nil {
		// 没有找到活跃会话
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"data":        nil,
			"server_time": time.Now().Format("2006-01-02 15:04:05"),
		})
		return
	}

	// 计算停车时长（分钟）
	durationMinutes := int(time.Since(session.StartTime).Minutes())

	// 计算计费进度
	var billingProgressPercent int
	var remainingMinutesToNextBilling int

	if session.NextBillingTime != nil {
		totalBillingDuration := 60 // 一小时计费周期
		remainingMinutes := int(session.NextBillingTime.Sub(time.Now()).Minutes())
		if remainingMinutes < 0 {
			remainingMinutes = 0
		}
		remainingMinutesToNextBilling = remainingMinutes
		billingProgressPercent = int((float64(totalBillingDuration-remainingMinutes) / float64(totalBillingDuration)) * 100)
	}

	// 构建响应
	response := ParkingSessionResponse{
		ID:           session.ID,
		VehiclePlate: session.Vehicle.PlateNumber,
		ParkingLot: ParkingLotInfo{
			ID:    session.ParkingLot.ID,
			Name:  session.ParkingLot.Name,
			Floor: "B1", // 可以从数据库字段获取
			Area:  "A区", // 可以从数据库字段获取
		},
		SpotCode:                      session.SpotCode,
		SpotType:                      session.SpotType,
		StartTime:                     session.StartTime.Format("2006-01-02 15:04:05"),
		DurationMinutes:               durationMinutes,
		FeeCurrent:                    session.FeeCurrent,
		NextBillingTime:               formatTimePtr(session.NextBillingTime),
		NextFeeAmount:                 session.NextFeeAmount,
		BillingProgressPercent:        billingProgressPercent,
		RemainingMinutesToNextBilling: remainingMinutesToNextBilling,
		CurrentBillingCycle:           session.CurrentBillingCycle,
		PricingRule:                   session.PricingRule,
		Navigation: NavigationInfo{
			Status:             session.NavigationStatus,
			RemainingDistanceM: session.RemainingDistanceM,
			EstimatedMinutes:   session.EstimatedMinutes,
			Destination: Position{
				Lat: session.DestinationLat,
				Lon: session.DestinationLon,
			},
			UserPosition:                 formatPositionPtr(session.UserPositionLat, session.UserPositionLon),
			ProgressToDestinationPercent: session.ProgressToDestination,
		},
		Status: session.Status,
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"data":        response,
		"server_time": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// GetParkingSessionHistory 获取停车会话历史记录
func GetParkingSessionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 5
	}

	offset := (page - 1) * pageSize

	var sessions []models.ParkingSession
	var total int64

	// 计算总数
	models.DB.Model(&models.ParkingSession{}).Where("user_id = ? AND status = ?", userID, "ended").Count(&total)

	// 获取历史记录
	models.DB.Preload("ParkingLot").Where("user_id = ? AND status = ?", userID, "ended").
		Order("start_time DESC").Limit(pageSize).Offset(offset).Find(&sessions)

	var historyItems []ParkingSessionHistoryItem
	for _, session := range sessions {
		item := ParkingSessionHistoryItem{
			ID:             session.ID,
			ParkingLotName: session.ParkingLot.Name,
			StartTime:      session.StartTime.Format("2006-01-02 15:04:05"),
			TotalFee:       session.FeeCurrent,
		}

		if session.EndTime != nil {
			item.EndTime = session.EndTime.Format("2006-01-02 15:04:05")
			item.DurationMinutes = int(session.EndTime.Sub(session.StartTime).Minutes())
		}

		historyItems = append(historyItems, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      historyItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// RefreshParkingNavigation 刷新导航信息
func RefreshParkingNavigation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	var session models.ParkingSession
	result := models.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "停车会话不存在"})
		return
	}

	// 模拟导航数据更新（实际应用中可能需要调用地图API）
	if session.NavigationStatus == "en_route" {
		// 模拟到达过程
		session.RemainingDistanceM = int(math.Max(0, float64(session.RemainingDistanceM)-50)) // 每次刷新减少50米
		session.EstimatedMinutes = int(math.Max(0, float64(session.EstimatedMinutes)-1))      // 时间减少1分钟

		if session.RemainingDistanceM <= 10 {
			session.NavigationStatus = "in_garage"
		}

		// 更新进度
		totalDistance := 2000.0 // 假设总距离2公里
		session.ProgressToDestination = int((totalDistance - float64(session.RemainingDistanceM)) / totalDistance * 100)
	}

	// 保存更新
	models.DB.Save(&session)

	// 返回更新后的导航信息
	navigation := NavigationInfo{
		Status:             session.NavigationStatus,
		RemainingDistanceM: session.RemainingDistanceM,
		EstimatedMinutes:   session.EstimatedMinutes,
		Destination: Position{
			Lat: session.DestinationLat,
			Lon: session.DestinationLon,
		},
		UserPosition:                 formatPositionPtr(session.UserPositionLat, session.UserPositionLon),
		ProgressToDestinationPercent: session.ProgressToDestination,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    navigation,
	})
}

// PayCurrentParkingFee 支付当前停车费用
func PayCurrentParkingFee(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	var session models.ParkingSession
	result := models.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "停车会话不存在"})
		return
	}

	// 模拟支付过程
	now := time.Now()
	session.EndTime = &now
	session.Status = "ended"

	// 创建停车记录
	record := models.ParkingRecord{
		VehicleID:    session.VehicleID,
		ParkingLotID: &session.ParkingLotID,
		Location:     "", // 从ParkingLot获取
		StartTime:    session.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:      now.Format("2006-01-02 15:04:05"),
		Fee:          session.FeeCurrent,
		Duration:     time.Since(session.StartTime).Hours(),
		SpotType:     session.SpotType,
	}

	// 在事务中更新
	tx := models.DB.Begin()
	if err := tx.Save(&session).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支付失败"})
		return
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支付失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "支付成功",
	})
}

// ExtendParkingSession 延长停车时间
func ExtendParkingSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的会话ID"})
		return
	}

	var request struct {
		Minutes int `json:"minutes"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	var session models.ParkingSession
	result := models.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "停车会话不存在"})
		return
	}

	// 延长下次计费时间
	if session.NextBillingTime != nil {
		extendedTime := session.NextBillingTime.Add(time.Duration(request.Minutes) * time.Minute)
		session.NextBillingTime = &extendedTime
	}

	models.DB.Save(&session)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "延长成功",
	})
}

// 辅助函数
func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	formatted := t.Format("2006-01-02 15:04:05")
	return &formatted
}

func formatPositionPtr(lat, lon *float64) *Position {
	if lat == nil || lon == nil {
		return nil
	}
	return &Position{
		Lat: *lat,
		Lon: *lon,
	}
}
