package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"api/internal/checker"
	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

var ErrTaskNotFound = errors.New("задание не найдено")

type CheckHandler struct {
	client    *supabase.Client
	pythonURL string
}

func NewCheckHandler(client *supabase.Client) *CheckHandler {
	pythonURL := os.Getenv("PYTHON_URL")
	if pythonURL == "" {
		pythonURL = "http://localhost:5080"
	}
	return &CheckHandler{
		client:    client,
		pythonURL: pythonURL,
	}
}

type CheckRequest struct {
	TaskID string `json:"task_id" binding:"required"`
	Answer string `json:"answer" binding:"required"`
}

type CheckResponse struct {
	Correct       bool   `json:"correct"`
	CorrectAnswer string `json:"correct_answer"`
	Explanation   string `json:"explanation"`
}

// Check — POST /api/v1/check
func (h *CheckHandler) Check(c *gin.Context) {
	var req CheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "нужны task_id и answer"})
		return
	}

	if !isValidUUID(req.TaskID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный task_id"})
		return
	}

	task, err := h.getTask(req.TaskID)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "задание не найдено"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
		}
		return
	}

	taskType := getStringField(task, "task_type")
	if taskType == "" {
		answer := getStringField(task, "answer")
		if answer == "" || answer == "-" {
			taskType = "code"
		} else {
			taskType = "choice"
		}
	}

	result := checker.Check(taskType, getStringField(task, "answer"), req.Answer)

	if result.NeedsPython {
		pythonResult, err := h.checkViaPython(task, req.Answer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ошибка проверки через Python: " + err.Error(),
			})
			return
		}
		result.Correct = pythonResult.Correct
	}

	explanation := getStringField(task, "solution")

	shortID := req.TaskID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	if result.Correct {
		fmt.Printf("  \033[1;32m  ✓ task=%s correct\033[0m\n", shortID)
	} else {
		fmt.Printf("  \033[1;31m  ✗ task=%s wrong → %s\033[0m\n", shortID, result.CorrectAnswer)
	}

	c.JSON(http.StatusOK, CheckResponse{
		Correct:       result.Correct,
		CorrectAnswer: result.CorrectAnswer,
		Explanation:   explanation,
	})
}

func (h *CheckHandler) getTask(taskID string) (map[string]interface{}, error) {
	var tasks []map[string]interface{}
	endpoint := fmt.Sprintf("tasks?select=*&id=eq.%s&limit=1", taskID)
	err := h.client.Query(endpoint, false, &tasks)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, ErrTaskNotFound
	}
	return tasks[0], nil
}

func (h *CheckHandler) checkViaPython(task map[string]interface{}, userAnswer string) (*checker.Result, error) {
	payload := map[string]string{
		"task_id":     getStringField(task, "id"),
		"task_type":   getStringField(task, "task_type"),
		"content":     getStringField(task, "content"),
		"answer":      getStringField(task, "answer"),
		"user_answer": userAnswer,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(h.pythonURL+"/ai/v1/check", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("чтение ответа Python: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Python вернул %d: %s", resp.StatusCode, string(respBody))
	}

	var pyResult struct {
		Correct bool `json:"correct"`
	}
	if err := json.Unmarshal(respBody, &pyResult); err != nil {
		return nil, fmt.Errorf("парсинг ответа Python: %w", err)
	}

	return &checker.Result{
		Correct:     pyResult.Correct,
		NeedsPython: false,
	}, nil
}

func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
