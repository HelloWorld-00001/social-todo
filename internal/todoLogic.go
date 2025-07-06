package internal

import (
	"fmt"
	. "github.com/coderconquerer/go-login-app/internal/TodoItem/models"
	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateTodo(c *gin.Context, db *gorm.DB) {
	var todo Todo
	if err := c.ShouldBind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var todo Todo

	if err := db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	var input Todo
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&todo).Updates(input)
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func GetAllTodos(c *gin.Context, db *gorm.DB) {
	var todos []Todo

	// default pagination
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	fmt.Printf("Limit: %d and Offset: %d\n", limit, offset)
	result := db.Limit(limit).Offset(offset).Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  todos,
		"total": result.RowsAffected,
	})
}
