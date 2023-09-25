/*
Copyright 2023 Codenotary Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"database/sql"
	"strconv"

	_ "github.com/codenotary/immudb/pkg/stdlib"

	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	connStr := "immudb://immudb:immudb@127.0.0.1:3322/todos?sslmode=disable"

	db, err := sql.Open("immudb", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.ExecContext(
		context.Background(),
		"CREATE TABLE IF NOT EXISTS todos(id INTEGER AUTO_INCREMENT, description VARCHAR(256), PRIMARY KEY id)",
	)
	if err != nil {
		log.Fatal(err)
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		return indexHandler(c, db)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return postHandler(c, db)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		return putHandler(c, db)
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		return deleteHandler(c, db)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

type TODO struct {
	ID          int
	Description string
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var id int
	var description string

	var todos []TODO

	rows, err := db.Query("SELECT id, description FROM todos")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &description)
		if err != nil {
			return err
		}

		todos = append(todos, TODO{ID: id, Description: description})
	}

	return c.Render("index", fiber.Map{
		"Todos": todos,
	})
}

func postHandler(c *fiber.Ctx, db *sql.DB) error {
	newTODO := TODO{}

	err := c.BodyParser(&newTODO)
	if err != nil {
		return err
	}

	if newTODO.Description != "" {
		_, err := db.Exec("INSERT INTO todos(description) VALUES ($1)", newTODO.Description)
		if err != nil {
			return err
		}
	}

	return c.Redirect("/")
}

func putHandler(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return err
	}

	description := c.Query("description")

	_, err = db.Exec("UPDATE todos SET description=$1 WHERE id=$2", description, id)
	if err != nil {
		return err
	}

	return c.SendString("updated")
}

func deleteHandler(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM todos WHERE id=$1", id)
	if err != nil {
		return err
	}

	return c.SendString("deleted")
}
