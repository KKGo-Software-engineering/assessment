package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nata-non/assessment/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=2022 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fmt.Println("Please use server.go for main file")
	db.AutoMigrate()
	fmt.Println("start at port:", os.Getenv("PORT"))
	list := []model.User{}
	db.Find(&list)
	fmt.Println(list)
	app := fiber.New()
	app.Get("/expenses", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(model.User{
			ID:     0,
			Title:  "Shopee",
			Amount: 690,
			Note:   "Pay later",
			Tags:   []string{"Dog", "Cat"},
		})
	})
	app.Post("/expenses", func(ctx *fiber.Ctx) error {
		//a := new(model.User)
		p := struct {
			Title  string   `json:"title"`
			Amount int      `json:"amount"`
			Note   string   `json:"note"`
			Tags   []string `json:"tags"`
		}{}
		if err := ctx.Status(http.StatusBadRequest).BodyParser(&p); err != nil {
			return err
		}
		a := model.User{
			Title:  p.Title,
			Amount: p.Amount,
			Note:   p.Note,
			Tags:   p.Tags,
		}
		db.Create(&a)
		return ctx.Status(http.StatusOK).JSON(a)
	})
	app.Get("/expenses/:id", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(model.User{
			ID:     0,
			Title:  "Shopee",
			Amount: 690,
			Note:   "Pay later",
			Tags:   []string{"Dog", "Cat"},
		})
	})
	app.Put("/expenses/:id", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(model.User{
			ID:     0,
			Title:  "Shopee",
			Amount: 690,
			Note:   "Pay later",
			Tags:   []string{"Dog", "Cat"},
		})
	})
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
