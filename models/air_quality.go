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

// AirQuality 空气质量数据
type AirQuality struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Location  string    `gorm:"size:100;not null" json:"location"`      // 监测位置
	AQI       int       `gorm:"not null" json:"aqi"`                    // 空气质量指数
	Level     string    `gorm:"size:20;not null" json:"level"`          // 空气质量等级
	PM25      float64   `gorm:"type:decimal(5,2);not null" json:"pm25"` // PM2.5浓度
	PM10      float64   `gorm:"type:decimal(5,2);not null" json:"pm10"` // PM10浓度
	O3        float64   `gorm:"type:decimal(5,2);not null" json:"o3"`   // 臭氧浓度
	NO2       float64   `gorm:"type:decimal(5,2);not null" json:"no2"`  // 二氧化氮浓度
	SO2       float64   `gorm:"type:decimal(5,2);not null" json:"so2"`  // 二氧化硫浓度
	CO        float64   `gorm:"type:decimal(5,2);not null" json:"co"`   // 一氧化碳浓度
	Timestamp time.Time `json:"timestamp"`                              // 数据时间戳
}

// AirQualityStats 空气质量统计
type AirQualityStats struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Date             time.Time `gorm:"type:date;not null" json:"date"`            // 统计日期
	AvgAQI           float64   `gorm:"type:decimal(5,2);not null" json:"avg_aqi"` // 平均AQI
	MaxAQI           int       `gorm:"not null" json:"max_aqi"`                   // 最高AQI
	MinAQI           int       `gorm:"not null" json:"min_aqi"`                   // 最低AQI
	PrimaryPollutant string    `gorm:"size:20" json:"primary_pollutant"`          // 主要污染物
	GoodHours        int       `gorm:"default:0" json:"good_hours"`               // 优良小时数
	PollutedHours    int       `gorm:"default:0" json:"polluted_hours"`           // 污染小时数
}
