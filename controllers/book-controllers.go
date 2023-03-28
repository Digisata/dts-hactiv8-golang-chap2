package controllers

import (
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

var bookData = []Book{}

func CreateBook(c *gin.Context) {
	newBook := Book{}

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	newBook.ID = len(bookData) + 1
	bookData = append(bookData, newBook)

	c.String(http.StatusCreated, "Created")
}

func GetBook(c *gin.Context) {
	c.JSON(http.StatusOK, bookData)
}

func GetBookById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for _, book := range bookData {
		if book.ID == id {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "NOT FOUND",
		"message": fmt.Sprintf("Book with ID %d not found", id),
	})
}

func UpdateBook(c *gin.Context) {
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

	for i, book := range bookData {
		if book.ID == id {
			bookRequest.ID = id
			bookData[i] = bookRequest
			c.String(http.StatusOK, "Updated")
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "NOT FOUND",
		"message": fmt.Sprintf("Book with ID %d not found", id),
	})
}

func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for i, book := range bookData {
		if book.ID == id {
			bookData = append(bookData[:i], bookData[i+1:]...)
			c.String(http.StatusOK, "Deleted")
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"status":  "NOT FOUND",
		"message": fmt.Sprintf("Book with ID %d not found", id),
	})
}
