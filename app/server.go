package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kritsana135/assessment/db"
	"github.com/Kritsana135/assessment/expense/delivery/http_"
	"github.com/Kritsana135/assessment/expense/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Connect()

	expUCase := usecase.NewExpUsecase()

	r := gin.Default()
	r.GET("/health", GetHealthCheck)

	http_.NewExpenseHandler(&r.RouterGroup, expUCase)

	r.Run(":" + os.Getenv("PORT"))
}

type ResponseHealthCheck struct {
	Message string    `json:"message"`
	Uptime  time.Time `json:"uptime"`
}

func GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseHealthCheck{
		Message: "OK",
		Uptime:  time.Now(),
	})
}
