package main

import (
	"context"
	"database/sql"
	"github.com/benciks/flow-backend/internal/database/db"
	graph2 "github.com/benciks/flow-backend/internal/graph"
	"github.com/benciks/flow-backend/internal/middleware"
	"github.com/labstack/echo/v4"
	em "github.com/labstack/echo/v4/middleware"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conn, err := sql.Open("sqlite3", "./flow.db")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Migration
	schemaFile, err := os.ReadFile("./internal/database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if _, err := conn.ExecContext(ctx, string(schemaFile)); err != nil {
		log.Fatal(err)
	}

	// Prepare data folders
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		if err := os.Mkdir("./data", 0755); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("./data/authorized_keys"); os.IsNotExist(err) {
		if err := os.Mkdir("./data/authorized_keys", 0755); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("./data/timewarrior"); os.IsNotExist(err) {
		if err := os.Mkdir("./data/timewarrior", 0755); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("./data/taskwarrior"); os.IsNotExist(err) {
		if err := os.Mkdir("./data/taskwarrior", 0755); err != nil {
			log.Fatal(err)
		}
	}

	// Load environment variables
	err = godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)

	e := echo.New()
	e.Use(middleware.Auth(queries))
	e.Use(em.CORS())

	srv := handler.NewDefaultServer(graph2.NewExecutableSchema(graph2.Config{Resolvers: &graph2.Resolver{
		DB: queries,
	}}))

	playgroundHandler := playground.Handler("GraphQL playground", "/query")
	e.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.POST("/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.Logger.Fatal(e.Start(":" + port))
}
