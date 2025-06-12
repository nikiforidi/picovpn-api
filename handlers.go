package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/anatolio-deb/picovpnd/picovpnd"
	"github.com/gin-gonic/gin"
)

// func try(context *gin.Context) {
// 	plan := Plan{
// 		ExpiresAt: time.Now().AddDate(0, 1, 0),
// 	}
// 	result := DB.Create(&plan)
// 	if result.Error != nil {
// 		context.AbortWithStatusJSON(500, map[string]any{
// 			"message": result.Error,
// 		})
// 		return
// 	}

// 	user := &User{TelegramID: initData.User.ID, PlanID: plan.ID, Plan: plan}
// 	result = DB.Create(&user)
// 	if result.Error != nil {
// 		context.AbortWithStatusJSON(500, map[string]any{
// 			"message": result.Error,
// 		})
// 		return
// 	}
// }

func userGet(context *gin.Context) {
	tgid := context.Param("tgid")
	i, err := strconv.Atoi(tgid)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": err,
		})
		return
	}
	user, err := UserGetByTelegramID(int64(i))
	if err != nil {
		context.AbortWithStatusJSON(http.StatusNotFound, map[string]any{
			"message": err,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, user)
}

func userAdd(context *gin.Context) {
	b, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": err,
		})
		return
	}
	password := Password{}
	err = json.Unmarshal(b, &password)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": err,
		})
		return
	}
	if !password.IsValid() {
		context.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{
			"message": "Password is not valid",
		})
		return
	}
	initData, ok := ctxInitData(context.Request.Context())
	if !ok {
		context.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
			"message": "Init data not found",
		})
		return
	}
	// Check if user already exists
	user, err := UserGetByTelegramID(initData.User.ID)
	if err != nil {
		// If user does not exist, create a new one
		plan := Plan{ExpiresAt: time.Now().AddDate(0, 1, 0)}
		result := DB.Create(&plan)
		if result.Error != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": result.Error,
			})
			return
		}
		user = &User{
			TelegramUsername: initData.User.Username,
			TelegramID:       initData.User.ID,
			PlanID:           plan.ID,
			Plan:             plan,
			// ChatID:           initData.User.ChatID,
			// TelegramUsername: initData.User.Username,
			// Account:          initData.User.Account,
			// Wallet:           initData.User.Wallet,
		}

		result = DB.Create(&user)
		if result.Error != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": result.Error,
			})
			return
		}
		context.IndentedJSON(http.StatusOK, user)
		// TODO: Get the daemon address from the config or environment variable
		// Create a new DaemonClient and add the user
		resp, err := NewDaemonClient("").UserAdd(context, &picovpnd.UserAddRequest{
			Username: user.TelegramUsername,
			Password: password.Password,
		})
		if err != nil {
			errors := []string{err.Error(), resp.Error}
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": errors,
			})
		} else {
			// If user already exists, return the existing user
			context.IndentedJSON(http.StatusOK, user)
			return
		}
	}
}
