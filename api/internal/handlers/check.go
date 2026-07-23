package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"api/internal/checker"
	"api/internal/supabase"

	"github.com/gin-gonic/gin"
)

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

	task, err := h.getTask(req.TaskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "задание не найдено"})
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

	if result.Correct {
		fmt.Printf("  \033[1;32m  ✓ task=%s correct\033[0m\n", req.TaskID[:8])
	} else {
		fmt.Printf("  \033[1;31m  ✗ task=%s wrong → %s\033[0m\n", req.TaskID[:8], result.CorrectAnswer)
	}

	c.JSON(http.StatusOK, CheckResponse{
		Correct:       result.Correct,
		CorrectAnswer: result.CorrectAnswer,
		Explanation:   explanation,
	})
}

func (h *CheckHandler) getTask(taskID string) (map[string]interface{}, error) {
	var tasks []map[string]interface{}
	endpoint := "tasks?select=*&id=eq." + taskID + "&limit=1"
	err := h.client.Query(endpoint, false, &tasks)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, io.ErrUnexpectedEOF
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

	resp, err := http.Post(h.pythonURL+"/ai/v1/check", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pyResult struct {
		Correct bool `json:"correct"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pyResult); err != nil {
		return nil, err
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
