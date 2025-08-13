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

// Device 设备信息
type Device struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DeviceName   string    `gorm:"size:100;not null" json:"device_name"`      // 设备名称
	DeviceType   string    `gorm:"size:50;not null" json:"device_type"`       // 设备类型
	Location     string    `gorm:"size:100" json:"location"`                  // 设备位置
	Status       string    `gorm:"size:20;default:'normal'" json:"status"`    // 设备状态
	SerialNumber string    `gorm:"size:100;uniqueIndex" json:"serial_number"` // 序列号
	Manufacturer string    `gorm:"size:100" json:"manufacturer"`              // 制造商
	InstallDate  time.Time `json:"install_date"`                              // 安装日期
}

// DeviceExpense 设备支出
type DeviceExpense struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DeviceID    uint      `gorm:"not null" json:"device_id"`            // 设备ID
	ExpenseType string    `gorm:"size:50;not null" json:"expense_type"` // 支出类型
	Amount      float64   `gorm:"type:decimal(10,2)" json:"amount"`     // 金额
	Description string    `gorm:"size:255" json:"description"`          // 描述
	ExpenseDate time.Time `json:"expense_date"`                         // 支出日期
	Period      string    `gorm:"size:20" json:"period"`                // 期间(month/quarter/year)

	// 关联
	Device Device `gorm:"foreignKey:DeviceID" json:"-"`
}

// DeviceMaintenanceRecord 设备维修记录
type DeviceMaintenanceRecord struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DeviceID        uint       `gorm:"not null" json:"device_id"`                // 设备ID
	MaintenanceType string     `gorm:"size:50;not null" json:"maintenance_type"` // 维修类型
	Description     string     `gorm:"size:500" json:"description"`              // 描述
	Status          string     `gorm:"size:20;default:'pending'" json:"status"`  // 状态
	Technician      string     `gorm:"size:100" json:"technician"`               // 技术员
	Cost            float64    `gorm:"type:decimal(10,2)" json:"cost"`           // 费用
	StartTime       time.Time  `json:"start_time"`                               // 开始时间
	EndTime         *time.Time `json:"end_time"`                                 // 结束时间
	Priority        string     `gorm:"size:20;default:'medium'" json:"priority"` // 优先级

	// 关联
	Device Device `gorm:"foreignKey:DeviceID" json:"-"`
}

// DeviceFaultStats 设备故障统计
type DeviceFaultStats struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DeviceType       string    `gorm:"size:50;not null" json:"device_type"`        // 设备类型
	FaultCount       int       `gorm:"not null" json:"fault_count"`                // 故障数量
	RepairCount      int       `gorm:"not null" json:"repair_count"`               // 修复数量
	FaultRate        float64   `gorm:"type:decimal(5,2)" json:"fault_rate"`        // 故障率
	StatDate         time.Time `json:"stat_date"`                                  // 统计日期
	Period           string    `gorm:"size:20" json:"period"`                      // 统计周期
	TrendDirection   string    `gorm:"size:10" json:"trend_direction"`             // 趋势方向
	PercentageChange float64   `gorm:"type:decimal(5,2)" json:"percentage_change"` // 变化百分比
}

// DeviceAlarm 设备报警
type DeviceAlarm struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	DeviceID    uint       `gorm:"not null" json:"device_id"`              // 设备ID
	AlarmType   string     `gorm:"size:50;not null" json:"alarm_type"`     // 报警类型
	Severity    string     `gorm:"size:20;not null" json:"severity"`       // 严重程度
	Message     string     `gorm:"size:500" json:"message"`                // 报警消息
	Status      string     `gorm:"size:20;default:'active'" json:"status"` // 状态
	AlarmTime   time.Time  `json:"alarm_time"`                             // 报警时间
	AckTime     *time.Time `json:"ack_time"`                               // 确认时间
	ResolveTime *time.Time `json:"resolve_time"`                           // 解决时间
	Location    string     `gorm:"size:100" json:"location"`               // 位置

	// 关联
	Device Device `gorm:"foreignKey:DeviceID" json:"-"`
}

// DeviceAlarmStats 设备报警统计
type DeviceAlarmStats struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	AlarmType string    `gorm:"size:50;not null" json:"alarm_type"` // 报警类型
	Count     int       `gorm:"not null" json:"count"`              // 数量
	Hour      int       `gorm:"not null" json:"hour"`               // 小时(0-23)
	StatDate  time.Time `json:"stat_date"`                          // 统计日期
	Severity  string    `gorm:"size:20" json:"severity"`            // 严重程度
}
