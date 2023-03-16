package routes

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sofianhadi1983/fiber-cucumber-demo/database"
	"github.com/sofianhadi1983/fiber-cucumber-demo/models"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/books", AddBook)
	app.Get("/books", Book)
	app.Put("/books/:id", Update)
	app.Delete("/books", Delete)
}

//AddBook
func AddBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.DB.Db.Create(&book)

	resp := make([]models.Book, 0)
	resp = append(resp, *book)
	return c.Status(http.StatusCreated).JSON(resp)
}

//Book
func Book(c *fiber.Ctx) error {
	books := []models.Book{}
	title := c.Query("title")

	if title != "" {
		database.DB.Db.Where("title = ?", title).Find(&books)
	} else {
		database.DB.Db.Find(&books)
	}

	return c.Status(200).JSON(books)
}

//Update
func Update(c *fiber.Ctx) error {
	book := new(models.Book)

	id := c.Params("id")
	ui64, _ := strconv.ParseUint(id, 10, 64)
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	book.ID = uint(ui64)
	database.DB.Db.Updates(&book)

	return c.Status(200).JSON(book)
}

//Delete
func Delete(c *fiber.Ctx) error {
	book := []models.Book{}
	title := new(models.Book)
	if err := c.BodyParser(title); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	database.DB.Db.Where("title = ?", title.Title).Delete(&book)

	return c.Status(200).JSON("deleted")
}
