package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
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
	initData := initdata.InitData{}
	err = json.Unmarshal(b, &initData)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": err,
		})
		return
	}
	plan := Plan{ExpiresAt: time.Now().AddDate(0, 1, 0)}
	result := DB.Create(&plan)
	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": result.Error,
		})
		return
	}
	user := User{
		PlanID:     plan.ID,
		Plan:       plan,
		TelegramID: initData.User.ID,
	}
	result = DB.Create(&user)
	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
			"message": result.Error,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, user)
}
