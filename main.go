package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type contextKey string

const (
	_initDataKey contextKey = "init-data"
)

// Returns new context with specified init data.
func withInitData(ctx context.Context, initData initdata.InitData) context.Context {
	return context.WithValue(ctx, _initDataKey, initData)
}

// Returns the init data from the specified context.
func ctxInitData(ctx context.Context) (initdata.InitData, bool) {
	initData, ok := ctx.Value(_initDataKey).(initdata.InitData)
	return initData, ok
}

// Middleware which authorizes the external client.
func authMiddleware(token string) gin.HandlerFunc {
	return func(context *gin.Context) {
		// We expect passing init data in the Authorization header in the following format:
		// <auth-type> <auth-data>
		// <auth-type> must be "tma", and <auth-data> is Telegram Mini Apps init data.
		body, err := io.ReadAll(context.Request.Body)
		if err != nil {
			context.AbortWithStatusJSON(400, map[string]any{
				"message": "Bad request: " + err.Error(),
			})
			return
		}
		tma := AuthBody{}
		err = json.Unmarshal(body, &tma)
		if err != nil {
			context.AbortWithStatusJSON(400, map[string]any{
				"message": "Bad request: " + err.Error(),
			})
			return
		}
		log.Println(string(body))

		if err := initdata.Validate(tma.TMA, token, time.Hour); err != nil {
			context.AbortWithStatusJSON(401, map[string]any{
				"message": err.Error(),
			})
			return
		}

		// Parse init data. We will surely need it in the future.
		initData, err := initdata.Parse(tma.TMA)
		if err != nil {
			context.AbortWithStatusJSON(500, map[string]any{
				"message": err.Error(),
			})
			return
		}

		context.Request = context.Request.WithContext(
			withInitData(context.Request.Context(), initData),
		)
	}
}

// Middleware which shows the user init data.
func authenticateOrCreateNewUser(context *gin.Context) {
	initData, ok := ctxInitData(context.Request.Context())
	if !ok {
		context.AbortWithStatusJSON(401, map[string]any{
			"message": "Init data not found",
		})
		return
	}

	user, err := UserGetByTelegramID(initData.User.ID)
	if err != nil {
		context.String(500, err.Error())
	} else if user != nil {
		context.JSON(200, user)
	} else {
		user := User{
			// PlanID:     plan.ID,
			// Plan:       plan,
			TelegramID:       initData.User.ID,
			ChatID:           initData.Chat.ID,
			TelegramUsername: initData.User.Username,
		}
		result := DB.Create(&user)
		if result.Error != nil {
			context.String(500, result.Error.Error())
		} else {
			context.JSON(200, user)
		}
	}
}

// func UserAdd(context *gin.Context) {
// 	b, err := io.ReadAll(context.Request.Body)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	user := userAdd{}
// 	err = json.Unmarshal(b, &user)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	context.JSON(200, userID)

// 	// fmt.Println(user.Password)
// }

func main() {
	// Your secret bot tgoken.
	token := os.Getenv("TG_BOT_TOKEN")

	r := gin.New()

	r.Use(authMiddleware(token), cors.Default())
	r.POST("/auth", authenticateOrCreateNewUser)
	// r.POST("/api/users", UserAdd)

	// err := r.Run(":8080")

	err := r.RunTLS(":8080", "/etc/letsencrypt/live/picovpn.ru/fullchain.pem", "/etc/letsencrypt/live/picovpn.ru/privkey.pem")
	if err != nil {
		panic(err)
	}
}
