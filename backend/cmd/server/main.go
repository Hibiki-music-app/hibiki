package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Hibiki-music-app/hibiki/backend/internal"
	"github.com/Hibiki-music-app/hibiki/backend/internal/handlers"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := internal.Load()

	// connect to db once instance allowed
	internal.InitDB(cfg.Host, cfg.User, cfg.Passwd, cfg.Host, cfg.PortInt)

	// config & init webauthn
	wconfig := &webauthn.Config{
		RPDisplayName: "HibikiGo Webauthn",               // Display Name for your site
		RPID:          cfg.Host,                          // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:5173"}, // The origin URLs allowed for WebAuthn (cause the cookies is created on front and RPO reject him if doesnt match) -> to change into variable
	}
	w, err := webauthn.New(wconfig)
	if err != nil {
		panic(err)
	}

	//start echo instance, cors, routes, middleware
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api/passkey")
	api.POST("/registerStart", func(c echo.Context) error {
		return handlers.BeginRegistration(w, c)
	})
	api.POST("/registerFinish", handlers.FinishRegistration)
	api.POST("/loginStart", handlers.BeginLogin)
	api.POST("/loginFinish", handlers.FinishLogin)
	api.GET("/private", handlers.PrivatePage, internal.LoggedInMiddleware())

	e.Logger.Fatal(e.Start(cfg.Port))
	origin := fmt.Sprintf("%s://%s%s", cfg.Proto, cfg.Host, cfg.Port)
	log.Printf("[INFO] server start on %s", origin)

}
