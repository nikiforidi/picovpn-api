package main

import (
	"context"
	"log"
	"os"
	"strings"
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
		authParts := strings.Split(context.GetHeader("authorization"), " ")
		if len(authParts) != 2 {
			context.AbortWithStatusJSON(401, map[string]any{
				"message": "Unauthorized",
			})
			return
		}

		log.Println(authParts)

		authType := authParts[0]
		authData := authParts[1]

		switch authType {
		case "X-Telegram-Data":
			// Validate init data. We consider init data sign valid for 1 hour from their
			// creation moment.
			if err := initdata.Validate(authData, token, time.Hour); err != nil {
				context.AbortWithStatusJSON(401, map[string]any{
					"message": err.Error(),
				})
				return
			}

			// Parse init data. We will surely need it in the future.
			initData, err := initdata.Parse(authData)
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
}

// Middleware which shows the user init data.
func showInitDataMiddleware(context *gin.Context) {
	initData, ok := ctxInitData(context.Request.Context())
	if !ok {
		context.AbortWithStatusJSON(401, map[string]any{
			"message": "Init data not found",
		})
		return
	}

	context.JSON(200, initData)
}

func main() {
	// Your secret bot token.
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	// Set gin mode to release.
	r := gin.New()

	gin.SetMode(gin.ReleaseMode)
	// Set gin logger.

	r.Use(
		authMiddleware(token),
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: func(param gin.LogFormatterParams) string {
				return "[" + param.TimeStamp.Format(time.RFC3339) + "] " +
					param.Method + " " +
					param.Path + " " +
					param.ClientIP + " " +
					param.ErrorMessage + " " +
					param.Latency.String() + "\n"
			},
			Output: os.Stdout,
		}),
		cors.New(cors.Config{
			AllowOrigins:     []string{"https://picovpn.ru", "https://www.picovpn.ru", "https://picovpn.ru:8080", "https://www.picovpn.ru:8080"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", "X-Telegram-Data", "Authorization"},
			ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://picovpn.ru"
			},
			MaxAge: 12 * time.Hour,
		}),
		gin.Recovery())
	r.POST("/api/auth", showInitDataMiddleware)
	// r.POST("/api/users", UserAdd)

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
