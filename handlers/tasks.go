package handlers

import (
	"net/http"
	"strconv"
	"time"
	"verba/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Task Представляет собой структуру задачи
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTask Отвечает за создание записи о задаче в базе данных.
func CreateTask(c *gin.Context) {
	var task Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO tasks (title, description, due_date) 
		VALUES ($1, $2, $3)
		RETURNING id, title, description, due_date, created_at, updated_at
	`
	err := config.DB.QueryRow(c, query, task.Title, task.Description, task.DueDate).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере."})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetAllTasks Возвращает список весь задач
func GetAllTasks(c *gin.Context) {
	query := `
		SELECT id, title, description, due_date, created_at, updated_at 
		FROM tasks
	`
	rows, err := config.DB.Query(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере."})
		return
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере."})
			return
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере."})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask Возвращает задачу по id
func GetTask(c *gin.Context) {
	var task Task

	id := c.Param("id")
	query := `SELECT id, title, description, due_date, created_at, updated_at FROM tasks WHERE ID=$1`
	err := config.DB.QueryRow(c, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере."})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask Обновляет задачу по id
func UpdateTask(c *gin.Context) {
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат данных"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный формат данных"})
		return
	}
	task.ID = id

	query := `
		UPDATE tasks 
		SET title = $1, description = $2, due_date = $3, updated_at = NOW() 
		WHERE id = $4
		RETURNING id, title, description, due_date, created_at, updated_at
	`

	err = config.DB.QueryRow(c, query, task.Title, task.Description, task.DueDate, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask Удаляет задачу по id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM tasks WHERE ID=$1`

	result, err := config.DB.Exec(c, query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Проблема на сервере"})
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.Status(http.StatusNoContent)
}
