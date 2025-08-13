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
	"os"
	"time"

	"urban_traffic_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	UserType string `json:"userType" binding:"required"`
}

type LoginResponse struct {
	Token   string      `json:"token"`
	User    models.User `json:"user"`
	Message string      `json:"message"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查找用户
	var user models.User
	if err := models.DB.Where("username = ? AND user_type = ?", req.Username, req.UserType).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "urban-traffic-secret-key-2024"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"username":  user.Username,
		"user_type": user.UserType,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   tokenString,
		User:    user,
		Message: "登录成功",
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "退出登录成功"})
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到用户信息"})
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}
