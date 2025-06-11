package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramUsername string `json:"telegramUsername"`
	TelegramID       int64  `json:"telegramId"`
	PlanID           uint
	Plan             Plan

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

type Password struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirmation"`
}

func (p Password) IsValid() bool {
	return p.Password != "" && p.Password == p.PasswordConfirm
}

type OCServer struct {
	gorm.Model
	DaemonAddress  string `json:"daemonAddress"`
	DaemonPort     int    `json:"daemonPort"`
	DaemonToken    string `json:"daemonToken"`
	DaemonUser     string `json:"daemonUser"`
	DaemonPassword string `json:"daemonPassword"`
	DaemonSSL      bool   `json:"daemonSSL"`
	DaemonCert     string `json:"daemonCert"`
}

func (o OCServer) UserAdd(username, password string) {

}

func (o OCServer) UserLock(username string) map[string]string {
	return map[string]string{
		"username": username,
	}
}

func (o OCServer) UserUnlock(username string) map[string]string {
	return map[string]string{
		"username": username,
	}
}

func (o OCServer) UserDelete(username string) map[string]string {
	return map[string]string{
		"username": username,
	}
}

func (o OCServer) UserGet(username string) map[string]string {
	return map[string]string{
		"username": username,
	}
}

func (o OCServer) UserList() map[string]string {
	return map[string]string{}
}
