package main

import (
	"Avito/internal/db"
	"Avito/internal/handler_manager"
	"Avito/internal/middleware"
	"Avito/internal/repository"
	"Avito/internal/service"
	"Avito/internal/utils"
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	httpPort := os.Getenv("SERVER_PORT")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDb(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.GetPool(ctx).Close()

	repo := repository.NewRepository(database)
	svc := service.NewService(repo)
	jwtGen := &utils.JWTGen{}
	hm := handler_manager.NewHandlerManager(svc, jwtGen)

	http.Handle("/api/info", middleware.AuthMiddleware(http.HandlerFunc(hm.Info)))
	http.Handle("/api/buy/", middleware.AuthMiddleware(http.HandlerFunc(hm.Buy)))
	http.Handle("/api/sendCoin", middleware.AuthMiddleware(http.HandlerFunc(hm.SendCoin)))
	http.HandleFunc("/api/auth", hm.Auth)

	log.Println("http serer start listening on port:", httpPort)
	err = http.ListenAndServe(":"+httpPort, nil)

	if err != nil {
		log.Fatal(err)
	}
}
