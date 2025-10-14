// @title ROMA System API
// @version 0.1
// @description API para gestión de ejercicios, programas, sesiones y modo maestro.
// @BasePath /
// @schemes http
//
// @contact.name ROMA Dev Team
// @contact.url https://github.com/vicepalma/roma-system
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Esquema: "Bearer <token>"

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vicepalma/roma-system/backend/internal/middleware"
	"github.com/vicepalma/roma-system/backend/internal/repository"
	"github.com/vicepalma/roma-system/backend/internal/security"
	"github.com/vicepalma/roma-system/backend/internal/service"
	httpHandlers "github.com/vicepalma/roma-system/backend/internal/transport/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	sr "github.com/vicepalma/roma-system/backend/internal/repository"
	ss "github.com/vicepalma/roma-system/backend/internal/service"
	sh "github.com/vicepalma/roma-system/backend/internal/transport/http"
)

// loadEnv loads .env (if present) and reads required vars.
func loadEnv() (port, dbURL, env string) {
	_ = godotenv.Load() // no falla si no existe .env
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL = os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL no definido (configura backend/.env)")
	}
	env = os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	return
}

// openDB opens a GORM connection and does a quick ping.
func openDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error abriendo DB: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error obteniendo sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// ping inicial
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("DB no responde al ping: %v", err)
	}
	return db
}

func main() {
	defTZ := os.Getenv("DEFAULT_TZ")
	if defTZ == "" {
		defTZ = "America/Santiago"
	}

	port, dbURL, env := loadEnv()
	db := openDB(dbURL)
	sqlDB, _ := db.DB()

	// router
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	// CORS global, ANTES de las rutas y de cualquier middleware de auth
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173", // Vite dev
			"http://127.0.0.1:5173", // por si el browser usa 127.0.0.1
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // si planeas cookies; con Bearer no es necesario, pero no molesta
		MaxAge:           12 * time.Hour,
	}))
	// después de crear r:
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.RequestID(), middleware.AccessLog())
	// health endpoints
	r.GET("/health/db", func(c *gin.Context) {
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// http server with graceful shutdown
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Repos
	userRepo := repository.NewUserRepository(db)
	exRepo := repository.NewExerciseRepository(db)
	progRepo := repository.NewProgramRepository(db)
	progSvc := service.NewProgramService(progRepo)
	progH := httpHandlers.NewProgramHandler(progSvc)

	histRepo := repository.NewHistoryRepository(db)
	histSvc := service.NewHistoryService(histRepo)
	histH := httpHandlers.NewHistoryHandler(histSvc)

	coachRepo := repository.NewCoachRepository(db)
	coachSvc := service.NewCoachService(coachRepo, histSvc)
	coachH := httpHandlers.NewCoachHandler(coachSvc, histSvc, userRepo)

	sessRepo := sr.NewSessionRepository(db)
	sessSvc := ss.NewSessionService(sessRepo, coachSvc)
	sessH := sh.NewSessionHandler(sessSvc)

	healthH := httpHandlers.NewHealthHandler(db)

	inviteRepo := repository.NewInviteRepository(db)
	inviteSvc := service.NewInviteService(inviteRepo, coachSvc, "")
	inviteH := httpHandlers.NewInviteHandler(inviteSvc)

	// Handlers
	authH := httpHandlers.NewAuthHandler(userRepo, db)

	// Servicio ejercicios
	svc := service.NewExerciseService(exRepo)
	exH := httpHandlers.NewExerciseHandler(svc)

	// Rutas públicas
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })

	// Grupo público con limiter para /auth
	pubAuth := r.Group("/", middleware.NewLimiter(3, 5, 5*time.Minute).Gin())
	authH.Register(pubAuth) // /auth/* y /me (me se auto-protege dentro del handler)                                              // /auth/* y /me (ojo: /me ya internamente requiere AuthRequired)

	// Health (elige una)
	healthH.Register(r) // público
	// healthH.Register(api) // protegido

	// Rutas protegidas
	api := r.Group("/api", security.AuthRequired())
	exH.Register(api)
	progH.Register(api)
	sessH.Register(api)
	histH.Register(api)
	coachH.Register(api)
	inviteH.Register(api)

	// start async
	go func() {
		log.Printf("API escuchando en http://localhost:%s (env=%s)", port, env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown forzado: %v", err)
	}
	_ = sqlDB.Close()
	log.Println("servidor detenido correctamente")
}
