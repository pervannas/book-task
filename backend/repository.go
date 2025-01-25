package main

import (
	"database/sql"
	"fmt"
	"log"
)

func connectToDatabase() *sql.DB {
	dsn := "user:password@tcp(127.0.0.1:3306)/db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection is alive
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database!")
	return db
}

func setupDatabase() *sql.DB {
	db := connectToDatabase()

	books, err := getAllBooks(db)
	if err != nil {
		log.Fatal("Couldn't find books", err)
	}

	for i := 0; i < len(books); i++ {
		deleteBookById(db, books[i].Id)
	}
	// TODO: These should be id 1-3 always
	insertBook(db, "The Fellowship of the Ring")
	insertBook(db, "The Two Towers")
	insertBook(db, "The Return of the King")

	return db
}

func getAllBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query("SELECT id, title FROM book")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.Id, &book.Title); err != nil {
			return []Book{}, fmt.Errorf("error with row %d, %v", book.Id, err)
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return []Book{}, fmt.Errorf("error in rows %v", err)
	}

	return books, nil
}

func insertBook(db *sql.DB, title string) error {
	_, err := db.Exec("insert into book (title) values (?);", title)
	if err != nil {
		log.Printf("Error while inserting value into table\n")
	}
	return err
}

func createMyTable(db *sql.DB) {
	sql := "CREATE TABLE IF NOT EXISTS book ( id INT AUTO_INCREMENT PRIMARY KEY, title NVARCHAR(100) NOT NULL)"
	_, err := db.Exec(sql)
	if err != nil {
		log.Printf("Error while creating table\n")
	}
}

func getBookById(db *sql.DB, id int) (Book, error) {
	row := db.QueryRow("SELECT * FROM book where id=?", id)
	var book Book
	if err := row.Scan(&book.Id, &book.Title); err != nil {
		if err == sql.ErrNoRows {
			return Book{}, fmt.Errorf("cannot find book with id %d", id)
		}
		return Book{}, fmt.Errorf("error with row %d, %v", id, err)
	}
	return book, nil
}

func updateBookById(db *sql.DB, title string, id int) error {
	// TODO: Make sure that it does not add a book if there is no book at id
	_, err := db.Exec("UPDATE book SET title=? WHERE id=?;", title, id)
	if err != nil {
		log.Printf("Error while updating value in table\n")
	}
	return err
}

func deleteBookById(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM book WHERE id=?;", id)
	if err != nil {
		log.Printf("Could not delete book by id %d, err: %v", id, err)
	}
	return err
}
