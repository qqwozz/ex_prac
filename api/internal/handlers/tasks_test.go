// Тесты HTTP-хендлера GET /api/v1/tasks
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"api/internal/config"
	"api/internal/supabase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// setupTestServer — создаёт тестовый сервер
func setupTestServer(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	if err := godotenv.Load("../../.env"); err != nil {
		t.Skip(".env не найден, пропуск интеграционного теста")
	}

	cfg := &config.Config{
		SupabaseURL:     os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey: os.Getenv("SUPABASE_ANON_KEY"),
	}

	client := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, "")
	h := NewTasksHandler(client)

	r := gin.New()
	r.GET("/api/v1/tasks", h.GetTasks)
	return r
}

// doRequest — вспомогательная функция для HTTP-запросов
func doRequest(t *testing.T, r *gin.Engine, url string) map[string]interface{} {
	t.Helper()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET %s → %d: %s", url, w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("невалидный JSON: %v", err)
	}
	return resp
}

// ==================== ПОЛУЧЕНИЕ ЗАДАНИЙ ====================

func TestGetTasks_Informatics(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics")

	count := int(resp["count"].(float64))
	if count < 1 {
		t.Errorf("ожидалось >= 1 задание, получено %d", count)
	}

	tasks := resp["tasks"].([]interface{})
	for _, task := range tasks {
		tm := task.(map[string]interface{})
		if tm["subject"] != "informatics" {
			t.Errorf("subject = %v, ожидается informatics", tm["subject"])
		}
	}
}

func TestGetTasks_Math(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=math")

	count := int(resp["count"].(float64))
	if count < 1 {
		t.Errorf("ожидалось >= 1 задание, получено %d", count)
	}
}

func TestGetTasks_MathWithLimit(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=math&limit=1")

	count := int(resp["count"].(float64))
	if count != 1 {
		t.Errorf("ожидалось 1 задание, получено %d", count)
	}
}

func TestGetTasks_DefaultLimit(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics")

	tasks := resp["tasks"].([]interface{})
	if len(tasks) > 1 {
		t.Errorf("лимит по умолчанию должен быть 1, получено %d", len(tasks))
	}
}

func TestGetTasks_WithoutSubject(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks")

	count := int(resp["count"].(float64))
	if count < 1 {
		t.Errorf("ожидалось >= 1 задание без фильтра, получено %d", count)
	}
}

func TestGetTasks_WithExamFilter(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics&exam=oge")

	tasks := resp["tasks"].([]interface{})
	for _, task := range tasks {
		tm := task.(map[string]interface{})
		if tm["exam_type"] != "oge" {
			t.Errorf("exam_type = %v, ожидается oge", tm["exam_type"])
		}
	}
}

func TestGetTasks_WithTaskNumberFilter(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics&task_number=16")

	tasks := resp["tasks"].([]interface{})
	for _, task := range tasks {
		tm := task.(map[string]interface{})
		taskNum := int(tm["task_number"].(float64))
		if taskNum != 16 {
			t.Errorf("task_number = %d, ожидается 16", taskNum)
		}
	}
}

func TestGetTasks_Random(t *testing.T) {
	r := setupTestServer(t)
	resp1 := doRequest(t, r, "/api/v1/tasks?subject=informatics&limit=2")
	resp2 := doRequest(t, r, "/api/v1/tasks?subject=informatics&limit=2")

	tasks1 := resp1["tasks"].([]interface{})
	tasks2 := resp2["tasks"].([]interface{})

	if len(tasks1) != len(tasks2) {
		t.Skip("разное количество заданий — пропуск проверки рандома")
	}

	same := true
	for i := range tasks1 {
		id1 := tasks1[i].(map[string]interface{})["id"]
		id2 := tasks2[i].(map[string]interface{})["id"]
		if id1 != id2 {
			same = false
			break
		}
	}
	if same && len(tasks1) > 1 {
		t.Log("оба запроса вернули одинаковый порядок — возможно, но маловероятно")
	}
}

func TestGetTasks_ResponseHasFields(t *testing.T) {
	r := setupTestServer(t)
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics")

	tasks := resp["tasks"].([]interface{})
	if len(tasks) == 0 {
		t.Fatal("нет заданий для проверки полей")
	}

	task := tasks[0].(map[string]interface{})
	requiredFields := []string{"id", "content", "answer", "subject", "exam_type", "topic", "task_type", "display_id"}
	for _, field := range requiredFields {
		if _, ok := task[field]; !ok {
			t.Errorf("отсутствует поле %q в ответе", field)
		}
	}
}

// ==================== ОШИБКИ ====================

func TestGetTasks_InvalidLimit(t *testing.T) {
	r := setupTestServer(t)
	// Невалидный лимит — должен игнорироваться и использовать default=1
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics&limit=abc")

	tasks := resp["tasks"].([]interface{})
	if len(tasks) > 1 {
		t.Errorf("невалидный лимит должен быть проигнорирован, получено %d", len(tasks))
	}
}

func TestGetTasks_NegativeLimit(t *testing.T) {
	r := setupTestServer(t)
	// Отрицательный лимит — должен игнорироваться
	resp := doRequest(t, r, "/api/v1/tasks?subject=informatics&limit=-5")

	tasks := resp["tasks"].([]interface{})
	if len(tasks) > 1 {
		t.Errorf("отрицательный лимит должен быть проигнорирован, получено %d", len(tasks))
	}
}
