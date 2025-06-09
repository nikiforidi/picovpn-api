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

type PlanType int

func (s PlanType) String() string {
	switch s {
	case Monthly:
		return "monthly"
	case Yearly:
		return "yearly"
	}
	return "unknown"
}

const (
	Monthly PlanType = iota
	Yearly
)

type UserPlan struct {
	gorm.Model
	UserID    uint
	User      User
	Type      PlanType
	IsActive  bool
	ExpiresAt time.Time
}

type AuthBody struct {
	TMA string `json:"tma"`
}
