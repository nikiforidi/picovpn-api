package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var TOKEN string

func init() {
	// Get token from environment variable.
	TOKEN = os.Getenv("TELEGRAM_BOT_TOKEN")
}

func Validate(context *gin.Context) {
	// Init data in raw format.
	initDataRaw := context.GetHeader("X-Telegram-Data")

	// Define how long since init data generation date init data is valid.
	expIn := 24 * time.Hour

	// Will return error in case, init data is invalid.
	err := initdata.Validate(initDataRaw, TOKEN, expIn)
	if err != nil {
		context.AbortWithStatusJSON(401, map[string]any{
			"message": "Unauthorized",
		})
		return
	}
	initData, err := initdata.Parse(initDataRaw)
	if err != nil {
		context.AbortWithStatusJSON(401, map[string]any{
			"message": "Unauthorized",
		})
		return
	}
	// If init data is valid, you can use it.
	user, err := UserGetByTelegramID(initData.User.ID)
	if err != nil {
		context.String(500, err.Error())
		return
	}

	if user != nil {
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
			return
		}
		context.JSON(200, user)
		return
	}
}

func main() {
	// Your secret bot tgoken.

	r := gin.New()

	r.Use(cors.Default())
	// r.GET("/", showInitDataMiddleware)
	r.POST("/auth", Validate)
	// r.POST("/api/users", UserAdd)

	// err := r.Run(":8080")

	err := r.RunTLS(":8080", "/etc/letsencrypt/live/picovpn.ru/fullchain.pem", "/etc/letsencrypt/live/picovpn.ru/privkey.pem")
	if err != nil {
		panic(err)
	}
}
