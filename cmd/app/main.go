package main

import (
	"classting/config"
	"classting/internal/news"
	"classting/internal/school"
	"classting/internal/subscription"
	"classting/internal/user"
	"classting/pkg/db"
	"classting/pkg/router"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// infrastructure
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := db.NewSql(cfg)
	if err != nil {
		log.Fatal(err)
	}
	router := router.NewServeRouter(cfg)

	// domain
	userRepository := user.NewUserRepository(db)
	schoolRepository := school.NewSchoolRepository(db)
	newsRepository := news.NewNewsRepository(db)
	subscriptionRepository := subscription.NewSubscriptionRepository(db)

	// service
	userService := user.NewUserService(userRepository, cfg)
	schoolService := school.NewSchoolService(schoolRepository, userRepository)
	newsService := news.NewNewsService(newsRepository, schoolRepository)
	subscriptionService := subscription.NewSubscriptionService(newsRepository, schoolRepository, subscriptionRepository)

	// controller
	userController := user.NewUserController(userService)
	schoolController := school.NewSchoolController(schoolService, cfg)
	newsController := news.NewNewsController(newsService)
	subscriptionController := subscription.NewSubscriptionController(subscriptionService)

	// routes
	user.RegisterRoutes(router, userController)
	school.RegisterRoutes(router, schoolController, cfg)
	news.RegisterRoutes(router, newsController, cfg)
	subscription.RegisterRoutes(router, subscriptionController, cfg)

	// http server
	srv := &http.Server{Addr: cfg.HTTP.Port, Handler: router}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 1 seconds.")
	}
	log.Println("Server exiting")
}
