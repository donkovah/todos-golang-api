package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	Email     *string
	Password  *string
	UpdatedAt time.Time
}

type Task struct {
	ID          uint
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Status      string
	UserID      int
	User        User
	CompletedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Todo struct {
	ID          uint
	Name        string
	Status      string
	CompletedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Comment struct {
	ID        uint
	Body      *string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&User{}, &Task{}, &Comment{})

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "API health check",
		})
	})

	r.GET("/tasks", func(c *gin.Context) {
		var tasks []Task
		result := db.Find(&tasks)

		if result.Error != nil {
			errorMessage := "Failed to retrieve tasks from the database"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"error":  errorMessage,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "Tasks retrieved successfully",
			"data":   tasks,
		})
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		var task Task
		idParam := c.Param("id")
		taskID, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		result := db.First(&task, taskID)

		if result.Error != nil {
			errorMessage := "Task not found"
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": 500,
				"error":  errorMessage,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "Task retrieved successfully",
			"data":   task,
		})
	})

	r.POST("/tasks", func(c *gin.Context) {
		var task Task
		err := c.ShouldBindJSON(&task)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.Create(&task)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create tasks",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "Task Created",
			"data":   task,
		})
	})

	r.POST("/tasks/:id/comment", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Comment on a task",
		})
	})

	r.PATCH("/tasks/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Update a Task",
		})
	})

	r.DELETE("/tasks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Delete a task",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
