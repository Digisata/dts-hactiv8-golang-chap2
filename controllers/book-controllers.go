package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"desc"`
}

type Controller interface {
	CreateBook(c *gin.Context)
	GetBook(c *gin.Context)
	GetBookById(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type bookController struct {
	DB *sql.DB
}

func NewBookController(db *sql.DB) Controller {
	return &bookController{
		DB: db,
	}
}

func (controller bookController) CreateBook(c *gin.Context) {
	newBook := Book{}

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sqlStatement := `
	INSERT INTO books (title, author, description)
	VALUES ($1, $2, $3)
	`

	err := controller.DB.QueryRowContext(c, sqlStatement, newBook.Title, newBook.Author, newBook.Description).Err()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusCreated, "Created")
}

func (controller bookController) GetBook(c *gin.Context) {
	result := []Book{}

	sqlStatement := `SELECT id, title, author, description FROM books`

	rows, err := controller.DB.QueryContext(c, sqlStatement)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		book := Book{}

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		result = append(result, book)
	}

	c.JSON(http.StatusOK, result)
}

func (controller bookController) GetBookById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sqlStatement := `SELECT id, title, author, description FROM books WHERE id = $1`

	row := controller.DB.QueryRowContext(c, sqlStatement, id)
	book := Book{}

	err = row.Scan(&book.ID, &book.Title, &book.Author, &book.Description)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "NOT FOUND",
			"message": fmt.Sprintf("Book with ID %d not found", id),
		})
	} else if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, book)
	}
}

func (controller bookController) UpdateBook(c *gin.Context) {
	bookRequest := Book{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sqlStatement := `
	UPDATE books
	SET title = $2, author = $3, description = $4
	WHERE id = $1
	`

	res, err := controller.DB.ExecContext(c, sqlStatement, id, bookRequest.Title, bookRequest.Author, bookRequest.Description)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if count, _ := res.RowsAffected(); count != 0 {
		c.String(http.StatusOK, "Updated")
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "NOT FOUND",
		"message": fmt.Sprintf("Book with ID %d not found", id),
	})
}

func (controller bookController) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	sqlStatement := `
	DELETE FROM books
	WHERE id = $1
	`

	res, err := controller.DB.ExecContext(c, sqlStatement, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if count, _ := res.RowsAffected(); count != 0 {
		c.String(http.StatusOK, "Deleted")
		return
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "NOT FOUND",
		"message": fmt.Sprintf("Book with ID %d not found", id),
	})
}
