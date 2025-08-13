/*
 * Copyright (c) 2025 LTQY. All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`
	UserType string `gorm:"size:20;not null;default:'user'" json:"user_type"` // admin, user
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// HashPassword 哈希密码
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type Vehicle struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID      uint   `gorm:"not null" json:"user_id"`
	PlateNumber string `gorm:"size:30;not null;uniqueIndex" json:"plate_number"`
	Brand       string `gorm:"size:50;not null" json:"brand"`
	Model       string `gorm:"size:50;not null" json:"model"`
	Color       string `gorm:"size:30" json:"color"`
	Type        string `gorm:"size:30;not null;default:'小型汽车'" json:"type"`
	RegDate     string `gorm:"size:15" json:"reg_date"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`

	// 关联字段
	User           User            `gorm:"foreignKey:UserID" json:"-"`
	ParkingRecords []ParkingRecord `gorm:"foreignKey:VehicleID" json:"-"`
}

type ParkingRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID    uint    `gorm:"not null" json:"vehicle_id"`
	ParkingLotID *uint   `json:"parking_lot_id"`
	Location     string  `gorm:"size:100;not null" json:"location"`
	StartTime    string  `gorm:"size:25" json:"start_time"`
	EndTime      string  `gorm:"size:25" json:"end_time"`
	Fee          float64 `json:"fee"`
	Duration     float64 `json:"duration"`                                  // 停车时长(小时)
	SpotType     string  `gorm:"size:20;default:'normal'" json:"spot_type"` // normal, charging, disabled, vip

	// 关联字段
	Vehicle    Vehicle    `gorm:"foreignKey:VehicleID" json:"-"`
	ParkingLot ParkingLot `gorm:"foreignKey:ParkingLotID" json:"-"`
}

type ParkingLot struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name           string  `gorm:"size:100;not null" json:"name"`
	Address        string  `gorm:"size:255" json:"address"`
	Latitude       float64 `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude      float64 `gorm:"type:decimal(11,8)" json:"longitude"`
	TotalSpots     int     `gorm:"not null;default:0" json:"total_spots"`
	AvailableSpots int     `gorm:"not null;default:0" json:"available_spots"`
	HourlyRate     float64 `gorm:"type:decimal(10,2);default:10.00" json:"hourly_rate"`
	IsActive       bool    `gorm:"default:true" json:"is_active"`
	OperatingHours string  `gorm:"size:50;default:'24小时'" json:"operating_hours"`
	PaymentMethods string  `gorm:"size:100;default:'微信,支付宝,现金'" json:"payment_methods"`

	// 关联字段
	SpecialSpots    []SpecialSpot    `gorm:"foreignKey:ParkingLotID" json:"special_spots"`
	ParkingRecords  []ParkingRecord  `gorm:"foreignKey:ParkingLotID" json:"-"`
	ParkingSessions []ParkingSession `gorm:"foreignKey:ParkingLotID" json:"-"`
}

// ParkingSession 停车会话表（实时停车状态）
type ParkingSession struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID       uint       `gorm:"not null;index" json:"user_id"`
	VehicleID    uint       `gorm:"not null" json:"vehicle_id"`
	ParkingLotID uint       `gorm:"not null" json:"parking_lot_id"`
	SpotCode     string     `gorm:"size:20;not null" json:"spot_code"`         // 车位编号
	SpotType     string     `gorm:"size:20;default:'normal'" json:"spot_type"` // normal, charging, disabled, vip
	StartTime    time.Time  `gorm:"not null" json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Status       string     `gorm:"size:20;not null;default:'active'" json:"status"` // active, ended, paid

	// 费用相关
	FeeRate             float64    `gorm:"type:decimal(10,2);not null" json:"fee_rate"`           // 每小时费率
	FeeCurrent          float64    `gorm:"type:decimal(10,2);default:0" json:"fee_current"`       // 当前费用
	NextBillingTime     *time.Time `json:"next_billing_time"`                                     // 下次计费时间
	NextFeeAmount       *float64   `gorm:"type:decimal(10,2)" json:"next_fee_amount"`             // 下次计费金额
	CurrentBillingCycle int        `gorm:"default:0" json:"current_billing_cycle"`                // 当前计费周期
	PricingRule         string     `gorm:"size:100;default:'首小时10元，后续每小时5元'" json:"pricing_rule"` // 计费规则描述

	// 导航相关
	NavigationStatus      string   `gorm:"size:20;default:'en_route'" json:"navigation_status"` // en_route, in_garage, parked
	RemainingDistanceM    int      `gorm:"default:0" json:"remaining_distance_m"`               // 剩余距离（米）
	EstimatedMinutes      int      `gorm:"default:0" json:"estimated_minutes"`                  // 预计到达时间（分钟）
	DestinationLat        float64  `gorm:"type:decimal(10,8)" json:"destination_lat"`           // 目的地纬度
	DestinationLon        float64  `gorm:"type:decimal(11,8)" json:"destination_lon"`           // 目的地经度
	UserPositionLat       *float64 `gorm:"type:decimal(10,8)" json:"user_position_lat"`         // 用户当前纬度
	UserPositionLon       *float64 `gorm:"type:decimal(11,8)" json:"user_position_lon"`         // 用户当前经度
	ProgressToDestination int      `gorm:"default:0" json:"progress_to_destination_percent"`    // 到达目的地进度百分比

	// 关联字段
	User       User       `gorm:"foreignKey:UserID" json:"user"`
	Vehicle    Vehicle    `gorm:"foreignKey:VehicleID" json:"vehicle"`
	ParkingLot ParkingLot `gorm:"foreignKey:ParkingLotID" json:"parking_lot"`
}

type SpecialSpot struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ParkingLotID   uint    `gorm:"not null" json:"parking_lot_id"`
	SpotType       string  `gorm:"size:20;not null" json:"spot_type"` // charging, disabled, vip
	TotalCount     int     `gorm:"not null;default:0" json:"total_count"`
	AvailableCount int     `gorm:"not null;default:0" json:"available_count"`
	AdditionalFee  float64 `gorm:"type:decimal(10,2);default:0.00" json:"additional_fee"`

	// 关联字段
	ParkingLot ParkingLot `gorm:"foreignKey:ParkingLotID" json:"-"`
}

// ParkingLotWithDistance 包含距离信息的停车场结构
type ParkingLotWithDistance struct {
	ParkingLot
	Distance        string                 `json:"distance"`
	SpecialSpotsMap map[string]SpecialSpot `json:"special_spots_map"`
}

// VehicleWithStats 包含统计信息的车辆结构
type VehicleWithStats struct {
	Vehicle
	LastParking *struct {
		Location string `json:"location"`
		Time     string `json:"time"`
	} `json:"last_parking"`
	Stats struct {
		ParkingCount int     `json:"parking_count"`
		TotalFee     float64 `json:"total_fee"`
		TotalHours   float64 `json:"total_hours"`
	} `json:"stats"`
}
