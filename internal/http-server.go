package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/daos"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/gateways"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/handlers"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/middlewares"
	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/usecases"
	webhttp "github.com/gsaaraujo/ecommerce-api-scenario-1/internal/web-http"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HttpServer struct {
	echo   *echo.Echo
	logger *slog.Logger
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		echo:   echo.New(),
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (h *HttpServer) Ready() {
	h.logger.Info("http server getting ready")

	h.echo.HidePort = true
	h.echo.HideBanner = true
	h.echo.Use(middleware.RequestID())
	h.echo.Use(middlewares.NewEchoRequestLoggerMiddleware(h.logger))
	h.echo.Use(middlewares.NewEchoRecoverMiddleware(h.logger))

	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		h.logger.Error(err.Error())
		os.Exit(1)
	}

	secretsClient := secretsmanager.NewFromConfig(defaultConfig)

	awsSecretsGateway := gateways.NewAwsSecretsGateway(secretsClient)

	postgresUrl, err := awsSecretsGateway.Get("POSTGRES_URL")
	if err != nil {
		h.logger.Error(err.Error())
		os.Exit(1)
	}

	pgxPool, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		h.logger.Error(err.Error())
		os.Exit(1)
	}

	accessTokenSigningKey, err := awsSecretsGateway.Get("ACCESS_TOKEN_SIGNING_KEY")
	if err != nil {
		h.logger.Error(err.Error())
		os.Exit(1)
	}

	// rabbitmqUrl, err := awsSecretsGateway.Get("RABBITMQ_URL")
	// if err != nil {
	// 	h.logger.Error(err.Error())
	// 	os.Exit(1)
	// }

	// rabbitmqConn, err := amqp091.Dial(rabbitmqUrl)
	// if err != nil {
	// 	h.logger.Error(err.Error())
	// 	os.Exit(1)
	// }

	// redisUrl, err := awsSecretsGateway.Get("REDIS_URL")
	// if err != nil {
	// 	h.logger.Error(err.Error())
	// 	os.Exit(1)
	// }

	// redisParsedUrl, err := redis.ParseURL(redisUrl)
	// if err != nil {
	// 	h.logger.Error(err.Error())
	// 	os.Exit(1)
	// }

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr:     redisParsedUrl.Addr,
	// 	Password: redisParsedUrl.Password,
	// })

	jsonBodyValidator, err := webhttp.NewJSONBodyValidator()
	if err != nil {
		h.logger.Error(err.Error())
		os.Exit(1)
	}

	customerDAO := daos.NewCustomerDAO(pgxPool)
	inventoryDAO := daos.NewInventoryDAO(pgxPool)
	cartDAO := daos.NewCartDAO(pgxPool)
	cartItemDAO := daos.NewCartItemDAO(pgxPool)
	productDAO := daos.NewProductDAO(pgxPool)
	addressDAO := daos.NewAddressDAO(pgxPool)

	httpZipCodeGateway := gateways.NewHttpZipCodeGateway(awsSecretsGateway)

	loginUsecase := usecases.NewLoginUsecase(customerDAO, awsSecretsGateway)
	signUpUsecase := usecases.NewSignUpUsecase(pgxPool, customerDAO)
	addProductUsecase := usecases.NewAddProductUsecase(pgxPool)
	addStockUsecase := usecases.NewAddStockUsecase(pgxPool, inventoryDAO)
	publishProductUsecase := usecases.NewPublishProductUsecase(pgxPool, productDAO)
	addProductToCartUsecase := usecases.NewAddProductToCartUsecase(pgxPool, cartDAO, cartItemDAO, productDAO, inventoryDAO)
	removeProductFromCartUsecase := usecases.NewRemoveProductFromCartUsecase(pgxPool, cartDAO, cartItemDAO)
	increaseProductQuantityInCartUsecase := usecases.NewIncreaseProductQuantityInCartUsecase(pgxPool, cartDAO, cartItemDAO, inventoryDAO)
	decreaseProductQuantityInCartUsecase := usecases.NewDecreaseProductQuantityInCartUsecase(pgxPool, cartDAO, cartItemDAO)
	addAddressUsecase := usecases.NewAddAddressUsecase(addressDAO, httpZipCodeGateway)

	loginHandler := handlers.NewLoginHandler(jsonBodyValidator, loginUsecase)
	signUpHandler := handlers.NewSignUpHandler(jsonBodyValidator, signUpUsecase)
	addProductHandler := handlers.NewAddProductHandler(jsonBodyValidator, addProductUsecase)
	addStockHandler := handlers.NewAddStockHandler(jsonBodyValidator, addStockUsecase)
	publishProductHandler := handlers.NewPublishProductHandler(jsonBodyValidator, publishProductUsecase)
	addProductToCartHandler := handlers.NewAddProductToCartHandler(jsonBodyValidator, addProductToCartUsecase)
	removeProductFromCartHandler := handlers.NewRemoveProductFromCartHandler(jsonBodyValidator, removeProductFromCartUsecase)
	increaseProductQuantityInCartHandler := handlers.NewIncreaseProductQuantityInCartHandler(jsonBodyValidator, increaseProductQuantityInCartUsecase)
	decreaseProductQuantityInCartHandler := handlers.NewDecreaseProductQuantityInCartHandler(jsonBodyValidator, decreaseProductQuantityInCartUsecase)
	getCartHandler := handlers.NewGetCartHandler(pgxPool, cartDAO)
	addAddressHandler := handlers.NewAddAddressHandler(jsonBodyValidator, addAddressUsecase)

	h.echo.GET("/health", func(c echo.Context) error {
		return c.NoContent(204)
	})

	v1 := h.echo.Group("/v1")

	v1.POST("/login", loginHandler.Handle)
	v1.POST("/sign-up", signUpHandler.Handle)

	echoJWTMiddleware := middlewares.NewEchoJWTMiddleware(accessTokenSigningKey)
	v1.POST("/admin/add-product", addProductHandler.Handle, echoJWTMiddleware)
	v1.POST("/admin/add-stock", addStockHandler.Handle, echoJWTMiddleware)
	v1.POST("/admin/publish-product", publishProductHandler.Handle, echoJWTMiddleware)

	v1.POST("/add-product-to-cart", addProductToCartHandler.Handle, echoJWTMiddleware)
	v1.POST("/remove-product-from-cart", removeProductFromCartHandler.Handle, echoJWTMiddleware)
	v1.POST("/increase-product-quantity-in-cart", increaseProductQuantityInCartHandler.Handle, echoJWTMiddleware)
	v1.POST("/decrease-product-quantity-in-cart", decreaseProductQuantityInCartHandler.Handle, echoJWTMiddleware)
	v1.POST("/add-address", addAddressHandler.Handle, echoJWTMiddleware)

	v1.GET("/cart", getCartHandler.Handle, echoJWTMiddleware)

	h.logger.Info("http server is now ready")
}

func (h *HttpServer) Start() {
	h.Ready()
	h.logger.Info("http server successfully started")
	err := h.echo.Start(":3333")

	if err != nil {
		h.logger.Error(err.Error())
		os.Exit(1)
	}
}

func (h *HttpServer) Echo() *echo.Echo {
	return h.echo
}
