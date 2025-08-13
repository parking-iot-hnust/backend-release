/*
 * Copyright (c) 2025 LTQY. All rights reserved.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 *
 */

package main

import (
	"log"
	"net/http"
	"os"

	"urban_traffic_backend/handlers"
	"urban_traffic_backend/middleware"
	"urban_traffic_backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 初始化数据库
	models.InitDB()

	// 创建Gin路由器
	r := gin.Default()

	// CORS配置 - 允许所有源
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// API路由组
	api := r.Group("/api")
	{
		// 用户认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), handlers.Logout)
			auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
		}

		// 车辆管理路由
		vehicles := api.Group("/vehicles")
		vehicles.Use(middleware.AuthMiddleware())
		{
			vehicles.GET("", handlers.GetVehicles)
			vehicles.POST("", handlers.CreateVehicle)
			vehicles.PUT("/:id/default", handlers.SetDefaultVehicle)
			vehicles.DELETE("/:id", handlers.DeleteVehicle)
		}

		// 停车场路由
		parking := api.Group("/parking")
		{
			parking.GET("/lots/nearby", handlers.GetNearbyParkingLots)
			parking.GET("/lots/:id", handlers.GetParkingLotDetails)
			parking.PUT("/lots/:id/availability", handlers.UpdateParkingLotAvailability)
			parking.GET("/stats", handlers.GetParkingStats)
			parking.GET("/current", middleware.AuthMiddleware(), handlers.GetCurrentParkingStatus)
		}

		// 用户停车会话路由
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			userParking := user.Group("/parking")
			{
				userParking.GET("/session/current", handlers.GetCurrentParkingSession)
				userParking.GET("/session/history", handlers.GetParkingSessionHistory)
				userParking.GET("/session/:sessionId/navigation", handlers.RefreshParkingNavigation)
				userParking.POST("/session/:sessionId/pay", handlers.PayCurrentParkingFee)
				userParking.POST("/session/:sessionId/extend", handlers.ExtendParkingSession)
			}
		}

		// 交通流量路由
		traffic := api.Group("/traffic")
		{
			traffic.GET("/flow", handlers.GetTrafficFlow)
			traffic.GET("/realtime", handlers.GetRealTimeTraffic)
			traffic.GET("/inout-flow", handlers.GetInOutFlowData)
			traffic.GET("/user-stats", handlers.GetTrafficUserStats)
			traffic.GET("/crossing-rate", handlers.GetCarCrossingRate)
			traffic.GET("/flow-chart", handlers.GetTrafficFlowChart)
			traffic.GET("/heatmap", handlers.GetTrafficHeatmap)
			traffic.GET("/congestion-reports", handlers.GetCongestionReports)
		}

		// 空气质量路由
		airQuality := api.Group("/air-quality")
		{
			airQuality.GET("/current", handlers.GetCurrentAirQuality)
			airQuality.GET("/history", handlers.GetAirQualityHistory)
			airQuality.GET("/stats", handlers.GetAirQualityStats)
			airQuality.POST("/update", handlers.UpdateAirQuality)
		}

		// 停车统计路由
		parking.GET("/saturation", handlers.GetParkingSaturation)
		parking.GET("/occupancy-rate", handlers.GetParkingOccupancyRate)
		parking.GET("/total-occupancy", handlers.GetTotalOccupancyRate)
		parking.GET("/motor-congestion", handlers.GetMotorParkingCongestion)
		parking.GET("/congestion-chart", handlers.GetParkingCongestionChart)
		parking.GET("/activity-analysis", handlers.GetParkingActivityAnalysis)
		parking.GET("/activity-realtime", handlers.GetParkingActivityRealtime)

		// 设备管理路由
		devices := api.Group("/devices")
		{
			devices.GET("/expenses", handlers.GetDeviceExpenses)
			devices.GET("/maintenance", handlers.GetDeviceMaintenanceRecords)
			devices.GET("/fault-stats", handlers.GetDeviceFaultStats)
			devices.GET("/alarms", handlers.GetDeviceAlarms)
			devices.GET("/alarms/stats", handlers.GetDeviceAlarmStats)
			devices.GET("/fault-trend", handlers.GetDeviceFaultTrend)
		}

		// 收费系统路由
		toll := api.Group("/toll")
		{
			toll.GET("/records", handlers.GetTollSystemData)
		}

		// 施工统计路由
		construction := api.Group("/construction")
		{
			construction.GET("/stats", handlers.GetConstructionStats)
		}

		// 监控系统路由
		monitoring := api.Group("/monitoring")
		{
			monitoring.GET("/cameras", handlers.GetMonitoringCameras)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
