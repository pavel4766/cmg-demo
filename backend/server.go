package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
	"math/big"
	"math/rand"
	"os"
	"time"
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

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Method() == "OPTIONS" {
			// Return early for preflight requests
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

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

	// New route for factorial calculation
	app.Get("/api/factorial", func(c *fiber.Ctx) error {
		return factorialHandler(c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

// Factorial calculation function using big integers
func factorial(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

// Handler for factorial calculation
func factorialHandler(c *fiber.Ctx) error {
  number := rand.Int63n(10000) // Generate a random number up to 100

	start := time.Now()

	var result *big.Int
	for i := 0; i < 100; i++ {
		result = factorial(number)
	}

	elapsed := time.Since(start)

	// Log the result and the duration of the computation
	log.Printf("Computed factorial of %d in %s\n", number, elapsed)

	return c.JSON(fiber.Map{
		"number":    number,
		"factorial": result.String(),
		"duration":  elapsed.String(),
	})
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
