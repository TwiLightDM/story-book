package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"story-book/internal/config"
	"story-book/internal/middlewares"
	"story-book/internal/services/bookservice"
	"story-book/internal/services/cardservice"
	"story-book/internal/services/cartservice"
	"story-book/internal/services/favouriteservice"
	"story-book/internal/services/genreservice"
	"story-book/internal/services/ratingservice"
	"story-book/internal/services/shopservice"
	"story-book/internal/services/userservice"
	"story-book/package/databases/postgres"
	"story-book/package/services/encryptservice"
	"story-book/package/services/jwtservice"
	"story-book/package/services/validateservice"
	"syscall"
	"time"

	_ "story-book/internal/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

func Run(cfg *config.Config) error {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())

	db, err := postgres.InitDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)
	if err != nil {
		return err
	}

	jwtService := jwtservice.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessDuration, cfg.JWT.RefreshDuration)
	encryptService := encryptservice.NewEncryptionService(cfg.SaltLength)
	validateService := validateservice.NewValidationService(cfg.MinPasswordSize)

	authMiddleware := middlewares.AuthMiddleware(jwtService)

	userRepository := userservice.NewUserRepository(db)
	userService := userservice.NewUserService(userRepository, jwtService, encryptService, validateService)
	userHandler := userservice.NewUserHandler(userService)

	bookRepository := bookservice.NewBookRepository(db)
	bookService := bookservice.NewBookService(bookRepository)
	bookHandler := bookservice.NewBookHandler(bookService)

	genreRepository := genreservice.NewGenreRepository(db)
	genreService := genreservice.NewGenreService(genreRepository)
	genreHandler := genreservice.NewGenreHandler(genreService)

	favouriteRepository := favouriteservice.NewFavouriteRepository(db)
	favouriteService := favouriteservice.NewFavouriteService(favouriteRepository)
	favouriteHandler := favouriteservice.NewFavouriteHandler(favouriteService)

	ratingRepository := ratingservice.NewRatingRepository(db)
	ratingService := ratingservice.NewRatingService(ratingRepository)
	ratingHandler := ratingservice.NewRatingHandler(ratingService)

	shopRepository := shopservice.NewShopRepository(db)
	shopService := shopservice.NewShopService(shopRepository)
	shopHandler := shopservice.NewShopHandler(shopService)

	cardRepository := cardservice.NewCardRepository(db)
	cardService := cardservice.NewCardService(cardRepository)
	cardHandler := cardservice.NewCardHandler(cardService)

	cartRepository := cartservice.NewCartRepository(db)
	cartService := cartservice.NewCartService(cartRepository)
	cartHandler := cartservice.NewCartHandler(cartService)

	registerRoutes(e, authMiddleware, userHandler, bookHandler, genreHandler, favouriteHandler, ratingHandler, shopHandler, cardHandler, cartHandler)

	server := &http.Server{
		Addr:    ":" + cfg.BackendPort,
		Handler: e,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func(db *gorm.DB) {
		log.Printf("Backend started on :%s", cfg.BackendPort)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}

		if errors.Is(err, http.ErrServerClosed) {
			err = postgres.CloseDB(db)
			if err != nil {
				log.Fatalf("failed to close database connection: %v", err)
			}
		}
	}(db)

	<-ctx.Done()
	log.Println("Shutting down backend server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("Backend server stopped gracefully")
	return nil
}

func registerRoutes(e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	userHandler *userservice.UserHandler,
	bookHandler *bookservice.BookHandler,
	genreHandler *genreservice.GenreHandler,
	favouriteHandler *favouriteservice.FavouriteHandler,
	ratingHandler *ratingservice.RatingHandler,
	shopHandler *shopservice.ShopHandler,
	cardHandler *cardservice.CardHandler,
	cartHandler *cartservice.CartHandler,
) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	public := e.Group("/auth")
	public.POST("/login", userHandler.Login)
	public.POST("/signup", userHandler.SignUp)
	public.POST("/admin", userHandler.SignUpAdmin, authMiddleware)
	public.POST("/refresh", userHandler.Refresh, authMiddleware)
	public.POST("/reset-password", userHandler.ResetPassword, authMiddleware)

	users := e.Group("/users", authMiddleware)
	users.GET("/me", userHandler.ReadSelf)
	users.GET("/:id", userHandler.ReadUser)
	users.PUT("/me", userHandler.UpdateUser)
	users.PATCH("/me/password", userHandler.ChangePassword)
	users.DELETE("/me", userHandler.DeleteUser)

	books := e.Group("/books")
	books.POST("", bookHandler.CreateBook, authMiddleware)
	books.GET("", bookHandler.ReadBooks)
	books.GET("/:id", bookHandler.ReadBook)
	books.PUT("/:id", bookHandler.UpdateBook, authMiddleware)
	books.DELETE("/:id", bookHandler.DeleteBook, authMiddleware)

	genres := e.Group("/genres", authMiddleware)
	genres.POST("", genreHandler.CreateGenre)
	genres.DELETE("/:book_id", genreHandler.DeleteGenre)

	favourites := e.Group("/favourites", authMiddleware)
	favourites.POST("", favouriteHandler.CreateFavourite)
	favourites.GET("", favouriteHandler.ReadFavourites)
	favourites.DELETE("/:book_id", favouriteHandler.DeleteFavourite)

	ratings := e.Group("/ratings", authMiddleware)
	ratings.POST("", ratingHandler.CreateRating)
	ratings.DELETE("/:book_id", ratingHandler.DeleteRating)

	shops := e.Group("/shops")
	shops.POST("", shopHandler.CreateShop, authMiddleware)
	shops.GET("", shopHandler.ReadShops)
	shops.GET("/:id", shopHandler.ReadShop)
	shops.PUT("/:id", shopHandler.UpdateShop, authMiddleware)
	shops.DELETE("/:id", shopHandler.DeleteShop, authMiddleware)

	cards := e.Group("/cards", authMiddleware)
	cards.POST("", cardHandler.CreateCard)
	cards.GET("", cardHandler.ReadCards)
	cards.DELETE("", cardHandler.DeleteCard)

	carts := e.Group("/carts", authMiddleware)
	carts.POST("", cartHandler.CreateCart)
	carts.GET("", cartHandler.ReadCarts)
	carts.PATCH("/:book_id", cartHandler.UpdateCart)
	carts.DELETE("/:book_id", cartHandler.DeleteCart)
}
