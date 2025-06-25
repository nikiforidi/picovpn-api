package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Your secret bot token.
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	// Set gin mode to release.
	r := gin.New()

	gin.SetMode(gin.ReleaseMode)
	// Set gin logger.

	r.Use(
		authMiddleware(token),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return strings.Contains(origin, "picovpn.ru")
			},
			MaxAge: 24 * time.Hour,
		}),
		gin.Recovery(),
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
	)
	r.POST("/api/auth", showInitDataMiddleware)
	r.GET("/api/users/", userGet)
	r.POST("/api/users", userAdd)
	r.POST("/api/daemons", registerDaemon)
	r.POST("/api/password-reset", passwordReset)

	// Run the server on port 8080 with TLS.
	// Make sure to replace the paths to your SSL certificate and key files.
	// You can use Let's Encrypt or any other certificate authority.
	// Make sure to have the certificate and key files in the specified paths.
	// If you are using Let's Encrypt, you can use the following command to generate the certificate and key files:
	// sudo certbot certonly --standalone -d picovpn.ru -d www.picovpn.ru
	// Make sure to have the certificate and key files in the specified paths.
	// If you are using a self-signed certificate, you can use the following command to generate the certificate and key files:
	// openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
	// err := r.RunTLS(":8080", "/etc/letsencrypt/live/picovpn.ru/fullchain.pem", "/etc/letsencrypt/live/picovpn.ru/privkey.pem")
	err := r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
