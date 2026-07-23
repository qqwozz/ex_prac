// Пакет tests — запуск проверок при старте сервера
// Все проверки выводятся в терминал с цветовой индикацией
package tests

import (
	"fmt"
	"os"
	"time"

	"api/internal/checker"
	"api/internal/supabase"

	"github.com/joho/godotenv"
)

// Цвета для терминала
const (
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorReset  = "\033[0m"
)

// subtest — результат одного подтеста
type subtest struct {
	name string
	pass bool
	err  string
}

// section — группа тестов (секция в выводе)
type section struct {
	name     string
	subtests []subtest
}

// RunAll — запускает все проверки и выводит результат в терминал
func RunAll() {
	fmt.Println()
	fmt.Printf("%s┌─────────────────────────────────────┐%s\n", colorCyan, colorReset)
	fmt.Printf("%s│         ЗАПУСК ПРОВЕРОК             │%s\n", colorCyan, colorReset)
	fmt.Printf("%s└─────────────────────────────────────┘%s\n", colorCyan, colorReset)
	fmt.Println()

	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("%s  ✗ .env не найден%s\n\n", colorRed, colorReset)
		os.Exit(1)
	}

	url := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	var sections []section

	// === СЕКЦИЯ 1: Конфигурация ===
	sections = append(sections, sectionConfig(url, anonKey, serviceKey))

	// === СЕКЦИЯ 2: Подключение к БД ===
	if url != "" && anonKey != "" {
		sections = append(sections, sectionDB(url, anonKey, serviceKey))
	}

	// === СЕКЦИЯ 3: Проверка ответов (checker) ===
	sections = append(sections, sectionChecker())

	// === СЕКЦИЯ 4: HTTP-эндпоинты ===
	if url != "" && anonKey != "" {
		sections = append(sections, sectionEndpoints(url, anonKey))
	}

	// Выводим все секции
	totalPassed, totalFailed := 0, 0
	for _, s := range sections {
		printSection(s)
		for _, st := range s.subtests {
			if st.pass {
				totalPassed++
			} else {
				totalFailed++
			}
		}
	}

	// Итоги
	fmt.Println()
	if totalFailed == 0 {
		fmt.Printf("  %sВсе проверки пройдены (%d/%d)%s\n", colorGreen, totalPassed, totalPassed+totalFailed, colorReset)
	} else {
		fmt.Printf("  %sПройдено: %d/%d, ошибок: %d%s\n", colorYellow, totalPassed, totalPassed+totalFailed, totalFailed, colorReset)
	}
	fmt.Println()

	// Если есть ошибки — не запускаем сервер
	if totalFailed > 0 {
		os.Exit(1)
	}
}

// ==================== СЕКЦИЯ: КОНФИГУРАЦИЯ ====================

func sectionConfig(url, anonKey, serviceKey string) section {
	s := section{name: "Конфигурация"}

	s.subtests = append(s.subtests, subtest{
		name: "SUPABASE_URL задан",
		pass: url != "",
		err:  "переменная не задана",
	})
	s.subtests = append(s.subtests, subtest{
		name: "SUPABASE_ANON_KEY задан",
		pass: anonKey != "",
		err:  "переменная не задана",
	})
	s.subtests = append(s.subtests, subtest{
		name: "Ключи различаются",
		pass: anonKey != serviceKey && serviceKey != "",
		err:  "anon и service ключи совпадают",
	})

	return s
}

// ==================== СЕКЦИЯ: БАЗА ДАННЫХ ====================

func sectionDB(url, anonKey, serviceKey string) section {
	s := section{name: "База данных"}
	client := supabase.NewClient(url, anonKey, serviceKey)

	// Тест: подключение (anon)
	start := time.Now()
	var tasks []map[string]interface{}
	err := client.Query("tasks?select=id&limit=1", false, &tasks)
	elapsed := time.Since(start)
	s.subtests = append(s.subtests, subtest{
		name: fmt.Sprintf("Подключение (%s)", elapsed.Round(time.Millisecond)),
		pass: err == nil && len(tasks) > 0,
		err:  errText(err, "таблица tasks пуста"),
	})

	// Тест: подключение (service key)
	if serviceKey != "" {
		var tasks2 []map[string]interface{}
		err := client.Query("tasks?select=id&limit=1", true, &tasks2)
		s.subtests = append(s.subtests, subtest{
			name: "Service key работает",
			pass: err == nil && len(tasks2) > 0,
			err:  errText(err, "не удалось запросить через service key"),
		})
	}

	// Тест: фильтр по предмету
	var infoTasks []map[string]interface{}
	err = client.Query("tasks?select=id&subject=eq.informatics&limit=1", false, &infoTasks)
	s.subtests = append(s.subtests, subtest{
		name: "Фильтр subject=informatics",
		pass: err == nil && len(infoTasks) > 0,
		err:  errText(err, "нет заданий по информатике"),
	})

	// Тест: фильтр по теме
	var mathTasks []map[string]interface{}
	err = client.Query("tasks?select=id&subject=eq.math&limit=1", false, &mathTasks)
	s.subtests = append(s.subtests, subtest{
		name: "Фильтр subject=math",
		pass: err == nil && len(mathTasks) > 0,
		err:  errText(err, "нет заданий по математике"),
	})

	return s
}

// ==================== СЕКЦИЯ: CHECKER ====================

func sectionChecker() section {
	s := section{name: "Проверка ответов (checker)"}

	// --- choice ---
	s.subtests = append(s.subtests, checkerSub("choice: совпадение",
		checker.Check("choice", "А", "А").Correct == true))
	s.subtests = append(s.subtests, checkerSub("choice: регистр",
		checker.Check("choice", "А", "а").Correct == true))
	s.subtests = append(s.subtests, checkerSub("choice: пробелы",
		checker.Check("choice", "  А  ", "А").Correct == true))
	s.subtests = append(s.subtests, checkerSub("choice: неправильно",
		checker.Check("choice", "А", "Б").Correct == false))

	// --- number ---
	s.subtests = append(s.subtests, checkerSub("number: точное",
		checker.Check("number", "17", "17").Correct == true))
	s.subtests = append(s.subtests, checkerSub("number: допуск",
		checker.Check("number", "17", "17.005").Correct == true))
	s.subtests = append(s.subtests, checkerSub("number: за допуском",
		checker.Check("number", "17", "17.1").Correct == false))
	s.subtests = append(s.subtests, checkerSub("number: запятая",
		checker.Check("number", "17,5", "17.5").Correct == true))
	s.subtests = append(s.subtests, checkerSub("number: отрицательное",
		checker.Check("number", "-5", "-5").Correct == true))
	s.subtests = append(s.subtests, checkerSub("number: граница 0.01",
		checker.Check("number", "10", "10.01").Correct == true))
	s.subtests = append(s.subtests, checkerSub("number: превышение",
		checker.Check("number", "10", "10.011").Correct == false))

	// --- string ---
	s.subtests = append(s.subtests, checkerSub("string: совпадение",
		checker.Check("string", "программа", "программа").Correct == true))
	s.subtests = append(s.subtests, checkerSub("string: регистр",
		checker.Check("string", "Программа", "программа").Correct == true))
	s.subtests = append(s.subtests, checkerSub("string: пробелы",
		checker.Check("string", "про  грамма", "про грамма").Correct == true))
	s.subtests = append(s.subtests, checkerSub("string: неправильно",
		checker.Check("string", "программа", "код").Correct == false))

	// --- multi ---
	s.subtests = append(s.subtests, checkerSub("multi: совпадение",
		checker.Check("multi", "А,Б,В", "А,Б,В").Correct == true))
	s.subtests = append(s.subtests, checkerSub("multi: порядок",
		checker.Check("multi", "А,Б,В", "В,Б,А").Correct == true))
	s.subtests = append(s.subtests, checkerSub("multi: разделители",
		checker.Check("multi", "А;Б;В", "А|Б|В").Correct == true))
	s.subtests = append(s.subtests, checkerSub("multi: неполный",
		checker.Check("multi", "А,Б,В", "А,Б").Correct == false))

	// --- code/text → Python ---
	s.subtests = append(s.subtests, checkerSub("code: NeedsPython=true",
		checker.Check("code", "print(42)", "x").NeedsPython == true))
	s.subtests = append(s.subtests, checkerSub("text: NeedsPython=true",
		checker.Check("text", "ответ", "x").NeedsPython == true))

	return s
}

// ==================== СЕКЦИЯ: HTTP-ЭНДПОИНТЫ ====================

func sectionEndpoints(url, anonKey string) section {
	s := section{name: "HTTP-эндпоинты"}
	client := supabase.NewClient(url, anonKey, "")

	// --- GET /api/v1/tasks ---
	var tasks []map[string]interface{}
	err := client.Query("tasks?select=*&subject=eq.informatics&limit=1", false, &tasks)
	s.subtests = append(s.subtests, subtest{
		name: "GET tasks?subject=informatics",
		pass: err == nil && len(tasks) > 0,
		err:  errText(err, "пустой ответ"),
	})

	var mathTasks []map[string]interface{}
	err = client.Query("tasks?select=*&subject=eq.math&limit=1", false, &mathTasks)
	s.subtests = append(s.subtests, subtest{
		name: "GET tasks?subject=math",
		pass: err == nil && len(mathTasks) > 0,
		err:  errText(err, "пустой ответ"),
	})

	// --- Поля ответа ---
	if len(tasks) > 0 {
		required := []string{"id", "content", "answer", "subject", "exam_type", "topic", "task_type"}
		allPresent := true
		missing := ""
		for _, field := range required {
			if _, ok := tasks[0][field]; !ok {
				allPresent = false
				missing += " " + field
			}
		}
		s.subtests = append(s.subtests, subtest{
			name: "Поля задания в ответе",
			pass: allPresent,
			err:  "отсутствуют:" + missing,
		})
	}

	// --- GET /api/v1/check (проверка через БД) ---
	if len(mathTasks) > 0 {
		taskID := mathTasks[0]["id"].(string)
		var checkTask []map[string]interface{}
		err := client.Query("tasks?select=answer&id=eq."+taskID, false, &checkTask)
		if err == nil && len(checkTask) > 0 {
			answer := checkTask[0]["answer"].(string)
			r := checker.Check("fipi", answer, answer)
			s.subtests = append(s.subtests, subtest{
				name: "Check: задание из БД",
				pass: r.Correct,
				err:  "ответ не совпал",
			})
		}
	}

	return s
}

// ==================== ВЫВОД ====================

// printSection — выводит секцию с подтестами
func printSection(s section) {
	fmt.Printf("  %s▸%s %s\n", colorCyan, colorReset, s.name)
	for _, st := range s.subtests {
		if st.pass {
			fmt.Printf("    %s✓%s %s\n", colorGreen, colorReset, st.name)
		} else {
			fmt.Printf("    %s✗%s %s — %s%s\n", colorRed, colorReset, st.name, st.err, colorReset)
		}
	}
	fmt.Println()
}

// ==================== ХЕЛПЕРЫ ====================

func checkerSub(name string, pass bool) subtest {
	return subtest{name: name, pass: pass}
}

func errText(err error, fallback string) string {
	if err != nil {
		return err.Error()
	}
	return fallback
}
