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
	// If init data is valid, you can parse it.
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
	// If user is not found, you can create a new one.
	if user == nil {
		user = &User{
			// PlanID:     plan.ID,
			// Plan:       plan,
			TelegramID:       initData.User.ID,
			ChatID:           initData.Chat.ID,
			TelegramUsername: initData.User.Username,
		}
		// Save user to database.
		result := DB.Create(&user)
		// If there is an error while saving user to database, return error.
		if result.Error != nil {
			context.AbortWithStatusJSON(500, map[string]any{
				"message": "Internal Server Error",
			})
			return
		}
	}
	// If you want to update user data, you can do it here.
	context.JSON(200, user)
}

func main() {
	// Set gin mode to release.
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	// Set gin logger.
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return "[" + param.TimeStamp.Format(time.RFC3339) + "] " +
				param.Method + " " +
				param.Path + " " +
				param.ClientIP + " " +
				param.ErrorMessage + " " +
				param.Latency.String() + "\n"
		},
		Output: os.Stdout,
	}))
	// Set gin recovery.
	r.Use(gin.Recovery())
	// Set CORS middleware.
	// You can customize CORS settings here.
	// For example, you can allow only specific origins, methods, headers, etc
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://picovpn.ru", "https://www.picovpn.ru", "https://picovpn.ru:8080", "https://www.picovpn.ru:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", "X-Telegram-Data", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://picovpn.ru"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.POST("/auth", Validate)
	// r.POST("/api/users", UserAdd)

	// err := r.Run(":8080")
	// Run the server on port 8080 with TLS.
	// Make sure to replace the paths to your SSL certificate and key files.
	// You can use Let's Encrypt or any other certificate authority.
	// Make sure to have the certificate and key files in the specified paths.
	// If you are using Let's Encrypt, you can use the following command to generate the certificate and key files:
	// sudo certbot certonly --standalone -d picovpn.ru -d www.picovpn.ru
	// Make sure to have the certificate and key files in the specified paths.
	// If you are using a self-signed certificate, you can use the following command to generate the certificate and key files:
	// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
	err := r.RunTLS(":8080", "/etc/letsencrypt/live/picovpn.ru/fullchain.pem", "/etc/letsencrypt/live/picovpn.ru/privkey.pem")
	if err != nil {
		panic(err)
	}
}
