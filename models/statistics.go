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

	"gorm.io/gorm"
)

// ParkingSaturation 停车饱和度数据
type ParkingSaturation struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ParkingLotID   uint      `json:"parking_lot_id"`                           // 停车场ID
	SaturationRate float64   `gorm:"type:decimal(5,2)" json:"saturation_rate"` // 饱和度
	OccupiedSpots  int       `gorm:"not null" json:"occupied_spots"`           // 已占用车位
	TotalSpots     int       `gorm:"not null" json:"total_spots"`              // 总车位数
	Timestamp      time.Time `json:"timestamp"`                                // 记录时间
	Hour           int       `gorm:"not null" json:"hour"`                     // 小时(0-23)

	// 关联
	ParkingLot ParkingLot `gorm:"foreignKey:ParkingLotID" json:"-"`
}

// ParkingOccupancyRate 停车占用率数据
type ParkingOccupancyRate struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ParkingLotID  uint      `json:"parking_lot_id"`                          // 停车场ID
	OccupancyRate float64   `gorm:"type:decimal(5,2)" json:"occupancy_rate"` // 占用率
	Period        string    `gorm:"size:20" json:"period"`                   // 时间段
	Timestamp     time.Time `json:"timestamp"`                               // 记录时间

	// 关联
	ParkingLot ParkingLot `gorm:"foreignKey:ParkingLotID" json:"-"`
}

// TotalOccupancy 总占用率统计
type TotalOccupancy struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	TotalSpots     int       `gorm:"not null" json:"total_spots"`             // 总车位数
	OccupiedSpots  int       `gorm:"not null" json:"occupied_spots"`          // 已占用车位
	OccupancyRate  float64   `gorm:"type:decimal(5,2)" json:"occupancy_rate"` // 占用率
	AvailableSpots int       `json:"available_spots"`                         // 可用车位
	Timestamp      time.Time `json:"timestamp"`                               // 记录时间
	Period         string    `gorm:"size:20" json:"period"`                   // 统计周期
}

// MotorParkingCongestion 非机动车停车拥堵数据
type MotorParkingCongestion struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Location        string    `gorm:"size:100;not null" json:"location"`        // 位置
	CongestionLevel string    `gorm:"size:20;not null" json:"congestion_level"` // 拥堵等级
	MotorCount      int       `gorm:"not null" json:"motor_count"`              // 非机动车数量
	Capacity        int       `gorm:"not null" json:"capacity"`                 // 容量
	CongestionRate  float64   `gorm:"type:decimal(5,2)" json:"congestion_rate"` // 拥堵率
	Timestamp       time.Time `json:"timestamp"`                                // 记录时间
	ReportedBy      string    `gorm:"size:100" json:"reported_by"`              // 上报人
}

// TollRecord 收费记录
type TollRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VehicleID     uint       `json:"vehicle_id"`                                // 车辆ID
	PlateNumber   string     `gorm:"size:30;not null" json:"plate_number"`      // 车牌号
	TollStation   string     `gorm:"size:100;not null" json:"toll_station"`     // 收费站
	EntryTime     time.Time  `json:"entry_time"`                                // 入口时间
	ExitTime      *time.Time `json:"exit_time"`                                 // 出口时间
	Amount        float64    `gorm:"type:decimal(10,2)" json:"amount"`          // 收费金额
	PaymentMethod string     `gorm:"size:50" json:"payment_method"`             // 支付方式
	Status        string     `gorm:"size:20;default:'completed'" json:"status"` // 状态
	Distance      float64    `gorm:"type:decimal(10,2)" json:"distance"`        // 距离(公里)

	// 关联
	Vehicle Vehicle `gorm:"foreignKey:VehicleID" json:"-"`
}

// ConstructionStats 施工统计
type ConstructionStats struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ProjectName   string     `gorm:"size:100;not null" json:"project_name"` // 项目名称
	Location      string     `gorm:"size:100;not null" json:"location"`     // 施工位置
	StartDate     time.Time  `json:"start_date"`                            // 开始日期
	EndDate       *time.Time `json:"end_date"`                              // 结束日期
	Status        string     `gorm:"size:20;not null" json:"status"`        // 状态
	Progress      float64    `gorm:"type:decimal(5,2)" json:"progress"`     // 进度百分比
	ImpactLevel   string     `gorm:"size:20" json:"impact_level"`           // 影响等级
	TrafficImpact string     `gorm:"size:500" json:"traffic_impact"`        // 交通影响
	Budget        float64    `gorm:"type:decimal(12,2)" json:"budget"`      // 预算
	ActualCost    float64    `gorm:"type:decimal(12,2)" json:"actual_cost"` // 实际成本
}

// MonitoringCamera 监控摄像头
type MonitoringCamera struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	CameraName  string    `gorm:"size:100;not null" json:"camera_name"`   // 摄像头名称
	Location    string    `gorm:"size:100;not null" json:"location"`      // 位置
	Latitude    float64   `gorm:"type:decimal(10,8)" json:"latitude"`     // 纬度
	Longitude   float64   `gorm:"type:decimal(11,8)" json:"longitude"`    // 经度
	Status      string    `gorm:"size:20;default:'online'" json:"status"` // 状态
	StreamUrl   string    `gorm:"size:255" json:"stream_url"`             // 视频流地址
	Type        string    `gorm:"size:50" json:"type"`                    // 摄像头类型
	Resolution  string    `gorm:"size:20" json:"resolution"`              // 分辨率
	ViewAngle   int       `json:"view_angle"`                             // 视角角度
	NightVision bool      `gorm:"default:false" json:"night_vision"`      // 夜视功能
	InstallDate time.Time `json:"install_date"`                           // 安装日期
}
