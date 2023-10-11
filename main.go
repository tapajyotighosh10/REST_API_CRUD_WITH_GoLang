//

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tapajyotighosh10/newRestApi/models"
	"github.com/tapajyotighosh10/newRestApi/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	log.Println(err, "''''''''''''''''")

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBookByID(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book delete successfully",
	})
	return nil
}
func (r *Repository) UpdateBookByID(context *fiber.Ctx) error {
	log.Println("gghggggggggggggggg")
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "ID cannot be empty",
		})
		return nil
	}
	log.Println("hdjhdjjhdasjbdjasbjdbsjkn")
	// Parse the request body to get the updated book data.
	updatedBook := &models.Books{}
	if err := context.BodyParser(&updatedBook); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Invalid request body",
		})
		return err
	}
	log.Println("ttttttttttttttttttttttttt")

	existingBook := &models.Books{}
	fmt.Println("the ID is", id)
	log.Println("aaaaaaaaaaaaaaaaaaaaaa")
	err := r.DB.Where("id = ?", id).First(existingBook).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	existingBook.Author = updatedBook.Author
	existingBook.Title = updatedBook.Title
	existingBook.Publisher = updatedBook.Publisher

	if err := r.DB.Save(existingBook).Error; err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "Could not update the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Book updated successfully",
		"data":    existingBook,
	})

	return nil

}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully",
		"data":    bookModel,
	})
	return nil
}
func (r *Repository) FetchBooks(context *fiber.Ctx) error {
	// Define a slice to hold book records
	var bookModels []models.Books

	// Fetch book data from the database
	err := r.DB.Find(&bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "Could not get books"})
		return err
	}

	// Create a new Excel workbook
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "ID")
	f.SetCellValue("Sheet1", "B1", "Author")
	f.SetCellValue("Sheet1", "C1", "Title")
	f.SetCellValue("Sheet1", "D1", "Publisher")

	// Populate the Excel sheet with book data
	for i, value := range bookModels {
		id := value.ID
		author := value.Author
		// t := author
		// log.Println(t, *t, &t, "daat")
		title := value.Title
		publisher := value.Publisher
		// log.Println(&author, "hhhhhhhhhhhhh")

		a := strconv.Itoa(i + 2)
		f.SetCellValue("Sheet1", "A"+a, id)
		f.SetCellValue("Sheet1", "B"+a, *author)
		f.SetCellValue("Sheet1", "C"+a, *title)
		f.SetCellValue("Sheet1", "D"+a, *publisher)
	}

	// Save the Excel file
	err = f.SaveAs("MyBooks.xlsx")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Excel File created and data inserted")
	}

	// Return a response indicating success
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Books fetched successfully",
		"data":    bookModels,
	})

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBookByID)
	api.Put("/update_book/:id", r.UpdateBookByID)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
	api.Get("/fetch", r.FetchBooks)

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.Conn(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
