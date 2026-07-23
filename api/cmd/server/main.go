package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"api/internal/config"
	"api/internal/handlers"
	"api/internal/supabase"
	"api/internal/tests"

	"github.com/gin-gonic/gin"
)

func main() {
	tests.RunAll()
	cfg := config.Load()
	client := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseServiceKey)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())

	tasks := handlers.NewTasksHandler(client)
	check := handlers.NewCheckHandler(client)
	r.GET("/api/v1/tasks", tasks.GetTasks)
	r.POST("/api/v1/check", check.Check)

	printBanner(cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func requestLogger() gin.HandlerFunc {
	const (
		reset  = "\033[0m"
		red    = "\033[1;31m"
		green  = "\033[1;32m"
		yellow = "\033[1;33m"
		cyan   = "\033[1;36m"
		gray   = "\033[90m"
		bold   = "\033[1m"
	)

	return func(c *gin.Context) {
		start := time.Now()
		query := c.Request.URL.RawQuery
		clientIP := c.ClientIP()

		c.Next()

		elapsed := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		statusColor := green
		if status >= 500 {
			statusColor = red
		} else if status >= 400 {
			statusColor = yellow
		}

		t := time.Now().Format("15:04:05")

		// Строка запроса
		queryPart := ""
		if query != "" {
			queryPart = fmt.Sprintf("%s?%s", path, query)
		} else {
			queryPart = path
		}

		// Время — цветом по скорости
		timeColor := green
		timeVal := elapsed.Round(time.Millisecond)
		if elapsed > 1*time.Second {
			timeColor = red
		} else if elapsed > 200*time.Millisecond {
			timeColor = yellow
		}

		logLine := fmt.Sprintf("  %s[%s]%s %s%-7s%s %s%s%s %s%d%s %s%s%s",
			gray, t, reset,
			cyan, method, reset,
			bold, queryPart, reset,
			statusColor, status, reset,
			timeColor, timeVal, reset,
		)

		// Доп. строка для POST /check — результат проверки
		if method == "POST" && path == "/api/v1/check" {
			logLine += fmt.Sprintf("  %sfrom=%s%s", gray, clientIP, reset)
		}

		// Ошибки — жирная красная строка
		if status >= 400 {
			logLine += fmt.Sprintf("\n  %s  ✗ %s%s", red, http.StatusText(status), reset)
		}

		// Supabase slow query — предупреждение
		if elapsed > 1*time.Second {
			logLine += fmt.Sprintf("\n  %s  ⚠ slow query > 1s%s", yellow, reset)
		}

		fmt.Println(logLine)
	}
}

func printBanner(port string) {
	fmt.Println()
	fmt.Printf("  \033[32m┌─────────────────────────────────────┐\033[0m\n")
	fmt.Printf("  \033[32m│  Сервер запущен                     │\033[0m\n")
	fmt.Printf("  \033[32m│  http://localhost:%s               │\033[0m\n", port)
	fmt.Printf("  \033[32m│                                     │\033[0m\n")
	fmt.Printf("  \033[32m│  GET  /api/v1/tasks                 │\033[0m\n")
	fmt.Printf("  \033[32m│  POST /api/v1/check                 │\033[0m\n")
	fmt.Printf("  \033[32m└─────────────────────────────────────┘\033[0m\n")
	fmt.Println()
}
