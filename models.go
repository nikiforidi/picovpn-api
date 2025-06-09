package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID int64 `json:"telegramId"`
	// PlanID     uint
	// Plan       UserPlan
	// ChatID           int64  `json:"chatId"`
	// TelegramUsername string `json:"telegramUsername"`
	// Account          string `json:"account"`
	// Wallet           string `json:"wallet"`
}

type UserPlan struct {
	gorm.Model
	UserID    uint
	User      User
	ExpiresAt time.Time
}
