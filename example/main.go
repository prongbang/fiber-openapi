package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/inteniquetic/fiberopenapi"
	"github.com/prongbang/oapigen"
)

type User struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	NullString *string `json:"nullString"`
	NullInt    *int64  `json:"nullInt"`
	NullFloat  *string `json:"nullFloat"`
	NullBool   *bool   `json:"nullBool"`
}

func main() {
	app := fiber.New()

	// OpenAPI middleware
	app.Use(fiberopenapi.New(oapigen.Config{
		Title:     "My API",
		Version:   "1.0.0",
		ServerURL: "http://localhost:3000",
		JSONPath:  "/openapi.json",
		DocsPath:  "/docs",
		SpecFile:  "./docs/openapi.json",
		Observe:   oapigen.ObsEnable,
	}))

	// --- GET: /healthz
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// --- Users: GET list / POST create
	app.Get("/users", func(c *fiber.Ctx) error {
		q := c.Query("q")
		users := []User{{ID: 1, Name: "Ada"}, {ID: 2, Name: "Bob"}}
		if q != "" {
			out := []User{}
			for _, u := range users {
				if q == u.Name {
					out = append(out, u)
				}
			}
			return c.JSON(out)
		}
		return c.JSON(users)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		var body struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}
		// mock create
		u := User{ID: 3, Name: body.Name}
		return c.Status(fiber.StatusCreated).JSON(u)
	})

	// --- Users: GET/PUT/PATCH/DELETE by id (path param)
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return c.JSON(User{ID: mustAtoi(c.Params("id")), Name: "Test"})
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		id := mustAtoi(c.Params("id"))
		var body struct {
			Name string `json:"name"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}
		return c.JSON(User{ID: id, Name: body.Name})
	})

	app.Patch("/users/:id", func(c *fiber.Ctx) error {
		id := mustAtoi(c.Params("id"))
		var body map[string]any
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
		}
		// mock partial update
		name := "Test"
		if v, ok := body["name"].(string); ok && v != "" {
			name = v
		}
		return c.JSON(User{ID: id, Name: name})
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	// --- OPTIONS: Allow header /users
	app.Options("/users", func(c *fiber.Ctx) error {
		c.Set("Allow", "GET,POST,OPTIONS")
		return c.SendStatus(fiber.StatusNoContent)
	})

	// --- HEAD: header no body
	app.Head("/ping", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		c.Set("X-Ping", "pong")
		return c.SendStatus(fiber.StatusOK)
	})

	// --- ALL:echo
	app.All("/echo", func(c *fiber.Ctx) error {
		payload := fiber.Map{
			"method":      c.Method(),
			"path":        c.Path(),
			"query":       c.Queries(),
			"contentType": string(c.Request().Header.ContentType()),
			"requestHeaders": func() map[string]string {
				h := map[string]string{}
				c.Request().Header.VisitAll(func(k, v []byte) {
					h[string(k)] = string(v)
				})
				return h
			}(),
			"body": string(c.Body()),
		}
		return c.JSON(payload)
	})

	fmt.Println("listening on :3000")
	_ = app.Listen(":3000")
}

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
