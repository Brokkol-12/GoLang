package main

import (
	"fmt"
	"golang/configs"
	"golang/internal/auth"
	"golang/internal/link"
	"golang/internal/stat"
	"golang/internal/user"
	"golang/pkg/db"
	"golang/pkg/event"
	"golang/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	//Repositries
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick()

	//Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	return stack(router)

}

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
