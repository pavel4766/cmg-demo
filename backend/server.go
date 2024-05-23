package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// Construct the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, dbport, dbname)

	// Print the connection string
	fmt.Println("Connection string:", connStr)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	// Initialize the datebase
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS todos (
        id SERIAL PRIMARY KEY,
        item TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)

	if err != nil {
		log.Fatalf("Error creating todos table: %v", err)
	}

	fmt.Println("Successfully created todos table!")

	app := fiber.New()

	app.Get("/api", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/api", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/api/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/api/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // React's local dev server
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	// Fetch todo items from the database
	rows, err := db.Query("SELECT id, item FROM todos")
	if err != nil {
		log.Println("Error querying todos:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Create a slice to store todo items
	var todos []fiber.Map

	// Iterate over the rows and append todo items to the slice
	for rows.Next() {
		var id int
		var item string
		if err := rows.Scan(&id, &item); err != nil {
			log.Println("Error scanning todo:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		}
		todos = append(todos, fiber.Map{
			"id":   id,
			"item": item,
		})
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over todos:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	// Return the todo items as JSON in the response
	return c.JSON(fiber.Map{"todos": todos})
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	var todo struct {
		Item string `json:"item"`
	}

	if err := c.BodyParser(&todo); err != nil {
		return err
	}

	_, err := db.Exec("INSERT INTO todos (item) VALUES ($1)", todo.Item)
	if err != nil {
		return err
	}

	return c.JSON(todo)
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	var todo struct {
		ID   int    `json:"id"`
		Item string `json:"item"`
	}

	if err := c.BodyParser(&todo); err != nil {
		return err
	}

	_, err := db.Exec("UPDATE todos SET item = $1 WHERE id = $2", todo.Item, todo.ID)
	if err != nil {
		return err
	}

	return c.JSON(todo)
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	var todo struct {
		ID int `json:"id"`
	}

	if err := c.BodyParser(&todo); err != nil {
		return err
	}

	_, err := db.Exec("DELETE FROM todos WHERE id = $1", todo.ID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
