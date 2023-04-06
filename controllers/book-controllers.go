package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Digisata/dts-hactiv8-golang-chap2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller interface {
	CreateBook(c *gin.Context)
	GetBook(c *gin.Context)
	GetBookById(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
}

type bookController struct {
	DB *gorm.DB
}

func NewBookController(db *gorm.DB) Controller {
	return &bookController{
		DB: db,
	}
}

func (controller bookController) CreateBook(c *gin.Context) {
	bookRequest := models.Book{}

	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := controller.DB.Create(&bookRequest).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, bookRequest)
}

func (controller bookController) GetBook(c *gin.Context) {
	books := []models.Book{}

	err := controller.DB.Find(&books).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, books)
}

func (controller bookController) GetBookById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	book := models.Book{}
	err = controller.DB.First(&book, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
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
	bookRequest := models.Book{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := c.ShouldBindJSON(&bookRequest); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	book := models.Book{}
	err = controller.DB.Model(&book).Where("id = ?", id).Updates(models.Book{Title: bookRequest.Title, Author: bookRequest.Author, Description: bookRequest.Description}).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (controller bookController) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	book := models.Book{}
	err = controller.DB.Where("id = ?", id).Delete(&book).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successfully",
	})
}
