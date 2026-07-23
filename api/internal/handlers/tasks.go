package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

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
