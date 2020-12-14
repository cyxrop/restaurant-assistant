package app

import (
	"context"
	baseLog "log"

	"go.uber.org/zap"

	"restaurant-assistant/pkg/app/server"
	"restaurant-assistant/pkg/app/service"
	"restaurant-assistant/pkg/config"
	"restaurant-assistant/pkg/handler"
	"restaurant-assistant/pkg/logger"
	"restaurant-assistant/pkg/storage/cache"
	"restaurant-assistant/pkg/storage/database"
	"restaurant-assistant/pkg/storage/database/repository"
)

func Run() {
	// Init config
	cfg := config.Init()
	ctx := context.Background()

	// Init logger
	log, err := logger.New()
	if err != nil {
		panic("cannot create logger")
	}

	defer func(*zap.Logger) {
		if err := log.Sync(); err != nil {
			baseLog.Fatal("sync logger", err.Error())
		}
	}(log)

	// Get DB Connection
	db, err := database.GetConnection(cfg.DB)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	// Init repository
	userRepository := repository.NewUserRepository(db, ctx)
	productRepository := repository.NewProductRepository(db, ctx)
	orderRepository := repository.NewOrderRepository(db, ctx)

	// Get Redis connection
	rdb := cache.GetConnection(cfg.Cache)

	// Init cache provider
	cacheProvider := cache.NewProvider(ctx, rdb)

	// Init services
	authenticationService, err := service.NewAuthenticationService(
		log,
		cacheProvider,
		userRepository,
	)
	if err != nil {
		log.Fatal("cannot create authentication service", zap.Error(err))
	}

	userService, err := service.NewUserService(
		log,
		userRepository,
	)
	if err != nil {
		log.Fatal("cannot create user service", zap.Error(err))
	}

	productService, err := service.NewProductService(
		log,
		productRepository,
		userRepository,
	)
	if err != nil {
		log.Fatal("cannot create user service", zap.Error(err))
	}

	orderService, err := service.NewOrderService(
		log,
		userRepository,
		orderRepository,
		productRepository,
	)
	if err != nil {
		log.Fatal("cannot create user service", zap.Error(err))
	}

	// Init server
	ras, err := server.NewRestaurantAssistantServer(
		log,
		authenticationService,
		userService,
		productService,
		orderService,
	)
	if err != nil {
		log.Fatal("cannot create ra ser{ver", zap.Error(err))
	}

	// Init handler
	handler.Init(ras)
}
