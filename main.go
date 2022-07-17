package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	database "github.com/snykk/simple_go_crud/database"
	"github.com/snykk/simple_go_crud/models"
	"gorm.io/gorm"
)

type Task struct {
	Task     string `json:"task"`
	Priority string `json:"priority"`
	Is_done  bool   `json:"is_done"`
}

type Repository struct {
	DB *gorm.DB
}

func greeting(context *fiber.Ctx) error {
	return context.SendString("Welcome to amazing api")
}

func (r *Repository) CreateTask(context *fiber.Ctx) error {
	task := Task{}
	task.Is_done = false

	err := context.BodyParser(&task)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "request failed",
			})
		log.Fatal("Error: ", err)
		return err
	}

	err = r.DB.Create(&task).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "could not create task",
			})
		fmt.Println("Error: ", err)
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"ok":      true,
		"message": "task added successfully",
	})
	return nil
}

func (r *Repository) DeleteTask(context *fiber.Ctx) error {
	taskModel := models.Task{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(taskModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "could not delete task",
		})
		log.Fatal("Error: ", err.Error)
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"ok":      true,
		"message": "task delete successfully",
	})
	return nil
}

func (r *Repository) GetTasks(context *fiber.Ctx) error {
	taskModel := &[]models.Task{}

	err := r.DB.Find(taskModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "could not get Tasks",
			})
		log.Fatal("Error: ", err)
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"ok":      false,
		"message": "Tasks fetched successfully",
		"data":    taskModel,
	})
	return nil
}

func (r *Repository) GetTaskByID(context *fiber.Ctx) error {

	id := context.Params("id")
	taskModel := &models.Task{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the id is", id)

	err := r.DB.Where("id = ?", id).First(taskModel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "could not get the task",
			})
		log.Fatal("Error: ", err)
		return nil
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"ok":      true,
		"message": "task id fetched successfully",
		"data":    taskModel,
	})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/task", r.CreateTask)
	api.Delete("task/:id", r.DeleteTask)
	api.Get("/task/:id", r.GetTaskByID)
	api.Get("/tasks", r.GetTasks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateTasks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	app.Get("/", greeting)
	r.SetupRoutes(app)
	app.Listen(":8080")
}
