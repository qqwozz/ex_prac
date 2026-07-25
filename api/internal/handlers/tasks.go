package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func isValidUUID(s string) bool {
	return uuidRe.MatchString(s)
}

type TasksHandler struct {
	client *supabase.Client
}

func NewTasksHandler(client *supabase.Client) *TasksHandler {
	return &TasksHandler{client: client}
}

// GetTasks — GET /api/v1/tasks
func (h *TasksHandler) GetTasks(c *gin.Context) {
	filters := []string{}

	if subject := c.Query("subject"); subject != "" {
		filters = append(filters, fmt.Sprintf("subject=eq.%s", url.QueryEscape(subject)))
	}
	if exam := c.Query("exam"); exam != "" {
		filters = append(filters, fmt.Sprintf("exam_type=eq.%s", url.QueryEscape(exam)))
	}
	if taskType := c.Query("type"); taskType != "" {
		filters = append(filters, fmt.Sprintf("task_type=eq.%s", url.QueryEscape(taskType)))
	}
	if topic := c.Query("topic"); topic != "" {
		filters = append(filters, fmt.Sprintf("topic=eq.%s", url.QueryEscape(topic)))
	}
	if difficulty := c.Query("difficulty"); difficulty != "" {
		filters = append(filters, fmt.Sprintf("difficulty=lte.%s", url.QueryEscape(difficulty)))
	}
	if taskNumber := c.Query("task_number"); taskNumber != "" {
		filters = append(filters, fmt.Sprintf("task_number=eq.%s", url.QueryEscape(taskNumber)))
	}

	limit := 1
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	const maxLimit = 100
	if limit > maxLimit {
		limit = maxLimit
	}

	endpoint := "tasks?select=*"
	if len(filters) > 0 {
		endpoint += "&" + strings.Join(filters, "&")
	}

	var tasks []map[string]interface{}
	if err := h.client.Query(endpoint, false, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rand.Shuffle(len(tasks), func(i, j int) {
		tasks[i], tasks[j] = tasks[j], tasks[i]
	})

	if limit > len(tasks) {
		limit = len(tasks)
	}
	tasks = tasks[:limit]

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"count": len(tasks),
	})
}

// GetTaskByID — GET /api/v1/tasks/:id
func (h *TasksHandler) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	if !isValidUUID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный UUID"})
		return
	}

	var tasks []map[string]interface{}
	endpoint := fmt.Sprintf("tasks?select=*&id=eq.%s&limit=1", id)
	if err := h.client.Query(endpoint, false, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "задание не найдено"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": tasks[0]})
}
