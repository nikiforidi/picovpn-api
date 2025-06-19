package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/anatolio-deb/picovpnd/grpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

		// daemons, err := DaemonsGetAll()
		// if err != nil {
		// 	context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
		// 		"message": err,
		// 	})
		// 	return
		// }
		// if len(daemons) == 0 {
		// 	context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
		// 		"message": "No daemons found",
		// 	})
		// 	return
		// }
		// Propogate new user to ocserve server instances through the daemons
		conn, err := grpc.NewClient("daemon:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("could not connect to daemon: %v", err)
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": err,
			})
			return
		}
		defer conn.Close()
		c := pb.NewOpenConnectServiceClient(conn)
		r, err := c.UserAdd(context.Request.Context(), &pb.UserAddRequest{
			Username: initData.User.Username,
			Password: password.Password,
		})
		if err != nil {
			log.Printf("could not add user: %v", err)
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": err,
			})
			return
		}
		if r.Error != "" {
			log.Printf("error adding user: %s", r.Error)
			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
				"message": r.Error,
			})
			return
		}
		log.Printf("user %s added successfully", initData.User.Username)
		context.IndentedJSON(http.StatusOK, map[string]string{
			"message":  "User added successfully",
			"username": initData.User.Username,
		})
		return
	}
	context.IndentedJSON(http.StatusOK, map[string]string{
		"message":  "User already exists",
		"username": user.TelegramUsername,
	})
}

// func registerDaemon(context *gin.Context) {
// 	b, err := io.ReadAll(context.Request.Body)
// 	if err != nil {
// 		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
// 			"message": err,
// 		})
// 		return
// 	}
// 	daemon := Daemon{}
// 	err = json.Unmarshal(b, &daemon)
// 	if err != nil {
// 		context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
// 			"message": err,
// 		})
// 		return
// 	}
// 	daemonRec, err := DaemonGetByAddress(daemon.Address)
// 	if err != nil {
// 		result := DB.Create(&daemon)
// 		if result.Error != nil {
// 			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
// 				"message": result.Error,
// 			})
// 			return
// 		}
// 		context.IndentedJSON(http.StatusOK, daemon)
// 	} else if daemonRec.Address == daemon.Address {
// 		daemonRec.Port = daemon.Port
// 		daemonRec.CertPEM = daemon.CertPEM
// 		result := DB.Save(&daemonRec)
// 		if result.Error != nil {
// 			context.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{
// 				"message": result.Error,
// 			})
// 			return
// 		}
// 		context.IndentedJSON(http.StatusOK, daemonRec)
// 		return
// 	}
// }
