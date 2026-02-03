package main

import (
	"context"
	"workout-tracker/config"
	"workout-tracker/handlers"
	"workout-tracker/middleware"
	"workout-tracker/repositories"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"*"},
	}
	r.Use(cors.New(corsConfig))
	err := loadConfig()
	if err != nil {
		panic(err)
	}
	connection, err := connectToDb()
	if err != nil {
		panic(err)
	}
	authRepository := repositories.NewAuthRepository(connection)
	exerciseRepository := repositories.NewExerciseRepository(connection)
	workoutRepository := repositories.NewWorkoutRepository(connection)
	authHandler := handlers.NewAuthHandler(authRepository)
	exerciseHandler := handlers.NewExerciseHandler(exerciseRepository)
	workoutHandler := handlers.NewWorkoutHandler(workoutRepository)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)
	})

	//users handlers
	r.POST("/user/signUp", authHandler.SignUp)
	r.POST("/user/signIn", authHandler.SignIn)
	//workout handlers
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/exercise", exerciseHandler.CreateExercise)
		auth.POST("/workout", workoutHandler.CreateWorkout)
		auth.GET("/user/workout", workoutHandler.GetMyWorkouts)
		auth.GET("/user/exercises", workoutHandler.GetMyExercises)
	}
	r.Run(config.Config.AppHost)
}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)

	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func loadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	if err := viper.BindEnv("APP_HOST"); err != nil {
		viper.SetDefault("APP_HOST", ":8070")
	}
	if err := viper.BindEnv("DB_CONNECTION_STRING"); err != nil {
		viper.SetDefault("DB_CONNECTION_STRING", "postgres://postgres:ansar2007+A@localhost:5430/task-manager-system?sslmode=disable")
	}
	if err := viper.BindEnv("JWT_SECRET_KEY"); err != nil {
		viper.SetDefault("JWT_SECRET_KEY", "supersecretkey")
	}
	if err := viper.BindEnv("JWT_EXPIRES_IN"); err != nil {
		viper.SetDefault("JWT_EXPIRES_IN", "24h")
	}
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}

	config.Config = &mapConfig

	return nil
}
