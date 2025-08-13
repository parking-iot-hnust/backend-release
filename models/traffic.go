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

// TrafficFlow 交通流量数据
type TrafficFlow struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Location    string    `gorm:"size:100;not null" json:"location"` // 路段位置
	FlowCount   int       `gorm:"not null" json:"flow_count"`        // 车流量
	Speed       float64   `gorm:"type:decimal(5,2)" json:"speed"`    // 平均速度 km/h
	Direction   string    `gorm:"size:20" json:"direction"`          // 方向 inbound/outbound
	VehicleType string    `gorm:"size:30" json:"vehicle_type"`       // 车辆类型
	Timestamp   time.Time `json:"timestamp"`                         // 数据时间戳
}

// TrafficUserStats 交通用户统计
type TrafficUserStats struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MotorCount      int       `gorm:"not null" json:"motor_count"`      // 机动车数量
	NonMotorCount   int       `gorm:"not null" json:"non_motor_count"`  // 非机动车数量
	PedestrianCount int       `gorm:"not null" json:"pedestrian_count"` // 行人数量
	Timestamp       time.Time `json:"timestamp"`                        // 统计时间
	Location        string    `gorm:"size:100" json:"location"`         // 统计位置
}

// TrafficHeatmap 交通热力图数据
type TrafficHeatmap struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RoadName        string    `gorm:"size:100;not null" json:"road_name"`        // 道路名称
	CongestionLevel float64   `gorm:"type:decimal(3,2)" json:"congestion_level"` // 拥堵程度 0-1
	GridX           int       `gorm:"not null" json:"grid_x"`                    // 网格X坐标
	GridY           int       `gorm:"not null" json:"grid_y"`                    // 网格Y坐标
	TimeFilter      string    `gorm:"size:20" json:"time_filter"`                // 时间过滤器
	Timestamp       time.Time `json:"timestamp"`                                 // 数据时间
}

// CongestionReport 拥堵播报
type CongestionReport struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Location    string    `gorm:"size:100;not null" json:"location"`      // 拥堵位置
	Severity    string    `gorm:"size:20;not null" json:"severity"`       // 严重程度
	Description string    `gorm:"size:500" json:"description"`            // 描述
	Duration    int       `gorm:"not null" json:"duration"`               // 持续时间(分钟)
	ReportTime  time.Time `json:"report_time"`                            // 播报时间
	Status      string    `gorm:"size:20;default:'active'" json:"status"` // 状态
}

// InOutFlowData 出入流量监控数据
type InOutFlowData struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Location     string    `gorm:"size:100;not null" json:"location"` // 监控位置
	InboundFlow  int       `gorm:"not null" json:"inbound_flow"`      // 入流量
	OutboundFlow int       `gorm:"not null" json:"outbound_flow"`     // 出流量
	NetFlow      int       `json:"net_flow"`                          // 净流量
	Timestamp    time.Time `json:"timestamp"`                         // 记录时间
	Hour         int       `gorm:"not null" json:"hour"`              // 小时(0-23)
}

// CarCrossingRate 车辆通过率数据
type CarCrossingRate struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Location     string    `gorm:"size:100;not null" json:"location"`      // 路口位置
	TotalCount   int       `gorm:"not null" json:"total_count"`            // 总车辆数
	PassedCount  int       `gorm:"not null" json:"passed_count"`           // 通过车辆数
	CrossingRate float64   `gorm:"type:decimal(5,2)" json:"crossing_rate"` // 通过率
	Timestamp    time.Time `json:"timestamp"`                              // 记录时间
	Period       string    `gorm:"size:20" json:"period"`                  // 时间段
}
