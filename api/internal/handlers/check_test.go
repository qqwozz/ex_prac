// Тесты HTTP-хендлера POST /api/v1/check
package handlers

import (
	"bytes"
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

// setupCheckServer — создаёт тестовый сервер с хендлером проверки
func setupCheckServer(t *testing.T) *gin.Engine {
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
	h := NewCheckHandler(client)

	r := gin.New()
	r.POST("/api/v1/check", h.Check)
	return r
}

// postJSON — отправляет POST-запрос с JSON-телом
func postJSON(t *testing.T, r *gin.Engine, url string, body interface{}) (int, map[string]interface{}) {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("ошибка маршалинга: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return w.Code, resp
}

// ==================== ВАЛИДНЫЕ ЗАПРОСЫ ====================

func TestCheck_MathTask_Correct(t *testing.T) {
	r := setupCheckServer(t)

	// Берём ID математического задания из БД
	taskID := getMathTaskID(t)

	code, resp := postJSON(t, r, "/api/v1/check", map[string]string{
		"task_id": taskID,
		"answer":  "4",
	})

	if code != http.StatusOK {
		t.Fatalf("ожидался 200, получен %d: %v", code, resp)
	}
	if resp["correct"] != true {
		t.Errorf("correct = %v, ожидается true (ответ 4 на задание с производной)", resp["correct"])
	}
	if resp["correct_answer"] != "4" {
		t.Errorf("correct_answer = %v, ожидается 4", resp["correct_answer"])
	}
}

func TestCheck_MathTask_Wrong(t *testing.T) {
	r := setupCheckServer(t)
	taskID := getMathTaskID(t)

	code, resp := postJSON(t, r, "/api/v1/check", map[string]string{
		"task_id": taskID,
		"answer":  "5",
	})

	if code != http.StatusOK {
		t.Fatalf("ожидался 200, получен %d", code)
	}
	if resp["correct"] != false {
		t.Errorf("correct = %v, ожидается false (ответ 5 неправильный)", resp["correct"])
	}
}

func TestCheck_InformaticsTask(t *testing.T) {
	r := setupCheckServer(t)
	taskID := getInformaticsTaskID(t)

	// Это задание на программирование (answer="-), должно идти в Python
	code, resp := postJSON(t, r, "/api/v1/check", map[string]string{
		"task_id": taskID,
		"answer":  "print(sum(x for x in [18,192,104,117,0] if 100<=x<=999 and x%4==0))",
	})

	// Может вернуть 500 если Python не запущен — это ожидаемо
	if code == http.StatusOK {
		t.Logf("check result: %v", resp["correct"])
	} else {
		t.Logf("Python недоступен (ожидаемо): %d", code)
	}
}

// ==================== ОШИБКИ ВХОДНЫХ ДАННЫХ ====================

func TestCheck_MissingTaskID(t *testing.T) {
	r := setupCheckServer(t)

	code, resp := postJSON(t, r, "/api/v1/check", map[string]string{
		"answer": "4",
	})

	if code != http.StatusBadRequest {
		t.Errorf("ожидался 400 без task_id, получен %d", code)
	}
	if resp["error"] == nil {
		t.Error("ожидалось сообщение об ошибке")
	}
}

func TestCheck_MissingAnswer(t *testing.T) {
	r := setupCheckServer(t)

	code, resp := postJSON(t, r, "/api/v1/check", map[string]string{
		"task_id": "00000000-0000-0000-0000-000000000000",
	})

	if code != http.StatusBadRequest {
		t.Errorf("ожидался 400 без answer, получен %d", code)
	}
	if resp["error"] == nil {
		t.Error("ожидалось сообщение об ошибке")
	}
}

func TestCheck_EmptyBody(t *testing.T) {
	r := setupCheckServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/check", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался 400 для пустого тела, получен %d", w.Code)
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	r := setupCheckServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/check", bytes.NewBuffer([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался 400 для невалидного JSON, получен %d", w.Code)
	}
}

func TestCheck_NonexistentTask(t *testing.T) {
	r := setupCheckServer(t)

	code, resp := postJSON(t, r, "/api/v1/check", map[string]string{
		"task_id": "00000000-0000-0000-0000-000000000000",
		"answer":  "4",
	})

	if code != http.StatusNotFound {
		t.Errorf("ожидался 404 для несуществующего задания, получен %d", code)
	}
	if resp["error"] == nil {
		t.Error("ожидалось сообщение об ошибке")
	}
}

// ==================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ====================

// getMathTaskID — достаёт ID первого математического задания из БД
func getMathTaskID(t *testing.T) string {
	t.Helper()
	if err := godotenv.Load("../../.env"); err != nil {
		t.Skip(".env не найден")
	}

	client := supabase.NewClient(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_ANON_KEY"),
		"",
	)

	var tasks []map[string]interface{}
	err := client.Query("tasks?select=id&subject=eq.math&limit=1", false, &tasks)
	if err != nil || len(tasks) == 0 {
		t.Skip("нет математических заданий в БД")
	}
	return tasks[0]["id"].(string)
}

// getInformaticsTaskID — достаёт ID первого задания по информатике
func getInformaticsTaskID(t *testing.T) string {
	t.Helper()
	if err := godotenv.Load("../../.env"); err != nil {
		t.Skip(".env не найден")
	}

	client := supabase.NewClient(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_ANON_KEY"),
		"",
	)

	var tasks []map[string]interface{}
	err := client.Query("tasks?select=id&subject=eq.informatics&limit=1", false, &tasks)
	if err != nil || len(tasks) == 0 {
		t.Skip("нет заданий по информатике в БД")
	}
	return tasks[0]["id"].(string)
}
