package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Book struct {
	Id    int    `json:"id" xml:"id"`
	Title string `json:"title" xml:"title"`
}

func startupAndShutdown(e *echo.Echo, ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down server")
		}
	}()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func getBooks(e *echo.Echo, db *sql.DB) {
	e.GET("/book", func(c echo.Context) error {
		books, err := getAllBooks(db)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("No books found, err: %v", err))
		}

		return c.JSON(http.StatusOK, books)
	})
}

func getSingleBook(e *echo.Echo, db *sql.DB) {
	e.GET("/book/:id", func(c echo.Context) error {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.Error(err)
		}
		book, err := getBookById(db, idInt)
		if err != nil {
			return c.String(http.StatusNotFound, "No books found\n")
		}

		return c.JSON(http.StatusOK, book)
	})
}

func createBook(e *echo.Echo, db *sql.DB) {
	e.POST("/book", func(c echo.Context) error {
		var body Book

		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid input, err: %v", err))
		}
		err := insertBook(db, body.Title)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid input, err: %v", err))
		}
		return c.String(http.StatusOK, "Book created")
	})
}

func updateBook(e *echo.Echo, db *sql.DB) {
	e.PUT("/book", func(c echo.Context) error {
		var body Book

		if err := c.Bind(&body); err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid input, err: %v", err))
		}
		if body.Title == "" {
			return c.String(http.StatusBadRequest, "Title and/or id cannot be zero")
		}
		err := updateBookById(db, body.Title, body.Id)
		if err != nil {
			return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid input, err: %v", err))
		}
		return c.String(http.StatusOK, "Book updated")
	})
}

func deleteBook(e *echo.Echo, db *sql.DB) {
	e.DELETE("/book/:id", func(c echo.Context) error {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Need to input an index. id: %s", id))
		}
		err = deleteBookById(db, idInt)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Got error while trying to delete item at id %d, err: %v", idInt, err))
		}
		return c.String(http.StatusOK, fmt.Sprintf("Book with id %d deleted", idInt))
	})
}

func main() {
	db := setupDatabase()

	defer db.Close()

	e := echo.New()
	// CORS middleware
	e.Use(echo.MiddlewareFunc(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*") // Allow all origins
			c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "Content-Type, Authorization")
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	}))
	getBooks(e, db)

	getSingleBook(e, db)

	createBook(e, db)

	updateBook(e, db)

	deleteBook(e, db)

	startupAndShutdown(e, context.Background())
}
