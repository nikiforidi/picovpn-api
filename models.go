package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID int64 `json:"telegramId"`
	PlanID     uint
	Plan       Plan
	// ChatID           int64  `json:"chatId"`
	// TelegramUsername string `json:"telegramUsername"`
	// Account          string `json:"account"`
	// Wallet           string `json:"wallet"`
}

type Plan struct {
	gorm.Model
	// UserID    uint      `json:"userId"`
	// User      User      `json:"user"`
	ExpiresAt time.Time `json:"expiresAt"`
}
