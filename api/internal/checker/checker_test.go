// Тесты логики проверки ответов
// Каждый тип задания — отдельная функция, каждый кейс — подтест
package checker

import "testing"

// ==================== CHOICE ====================

func TestCheckChoice(t *testing.T) {
	t.Run("совпадение", func(t *testing.T) {
		r := Check("choice", "А", "А")
		assertCorrect(t, r, true)
	})
	t.Run("регистр_независимый", func(t *testing.T) {
		r := Check("choice", "А", "а")
		assertCorrect(t, r, true)
	})
	t.Run("пробелы_обрезаются", func(t *testing.T) {
		r := Check("choice", "  А  ", "А")
		assertCorrect(t, r, true)
	})
	t.Run("разные_буквы", func(t *testing.T) {
		r := Check("choice", "А", "Б")
		assertCorrect(t, r, false)
	})
	t.Run("пустой_ответ_ученика", func(t *testing.T) {
		r := Check("choice", "А", "")
		assertCorrect(t, r, false)
	})
	t.Run("оба_пустые", func(t *testing.T) {
		r := Check("choice", "", "")
		assertCorrect(t, r, true)
	})
	t.Run("цифры_как_выбор", func(t *testing.T) {
		r := Check("choice", "1", "1")
		assertCorrect(t, r, true)
	})
	t.Run("длинный_ответ", func(t *testing.T) {
		r := Check("choice", "Правильный ответ", "Правильный ответ")
		assertCorrect(t, r, true)
	})
	t.Run("не_нужен_python", func(t *testing.T) {
		r := Check("choice", "А", "А")
		assertNeedsPython(t, r, false)
	})
	t.Run("correct_answer_сохраняется", func(t *testing.T) {
		r := Check("choice", "Вариант 3", "Вариант 3")
		assertCorrectAnswer(t, r, "Вариант 3")
	})
}

// ==================== NUMBER ====================

func TestCheckNumber(t *testing.T) {
	t.Run("точное_совпадение", func(t *testing.T) {
		r := Check("number", "17", "17")
		assertCorrect(t, r, true)
	})
	t.Run("в_допуске", func(t *testing.T) {
		r := Check("number", "17", "17.005")
		assertCorrect(t, r, true)
	})
	t.Run("за_допуском", func(t *testing.T) {
		r := Check("number", "17", "17.1")
		assertCorrect(t, r, false)
	})
	t.Run("запятая_вместо_точки", func(t *testing.T) {
		r := Check("number", "17,5", "17.5")
		assertCorrect(t, r, true)
	})
	t.Run("дробное_число", func(t *testing.T) {
		r := Check("number", "3.14", "3.14")
		assertCorrect(t, r, true)
	})
	t.Run("отрицательное", func(t *testing.T) {
		r := Check("number", "-5", "-5")
		assertCorrect(t, r, true)
	})
	t.Run("ноль", func(t *testing.T) {
		r := Check("number", "0", "0")
		assertCorrect(t, r, true)
	})
	t.Run("пробелы_вокруг", func(t *testing.T) {
		r := Check("number", "  42  ", "42")
		assertCorrect(t, r, true)
	})
	t.Run("граница_допуска_0.01", func(t *testing.T) {
		r := Check("number", "10", "10.01")
		assertCorrect(t, r, true)
	})
	t.Run("превышение_допуска_0.011", func(t *testing.T) {
		r := Check("number", "10", "10.011")
		assertCorrect(t, r, false)
	})
	t.Run("не_нужен_python", func(t *testing.T) {
		r := Check("number", "5", "5")
		assertNeedsPython(t, r, false)
	})
	t.Run("нечисловой_ответ_сравнивается_как_строка", func(t *testing.T) {
		r := Check("number", "abc", "abc")
		assertCorrect(t, r, true)
	})
	t.Run("нечисловой_и_числовой_не_совпадают", func(t *testing.T) {
		r := Check("number", "abc", "123")
		assertCorrect(t, r, false)
	})
}

// ==================== STRING ====================

func TestCheckString(t *testing.T) {
	t.Run("точное_совпадение", func(t *testing.T) {
		r := Check("string", "программа", "программа")
		assertCorrect(t, r, true)
	})
	t.Run("регистр_независимый", func(t *testing.T) {
		r := Check("string", "Программа", "программа")
		assertCorrect(t, r, true)
	})
	t.Run("пробелы_обрезаются", func(t *testing.T) {
		r := Check("string", "  программа  ", "программа")
		assertCorrect(t, r, true)
	})
	t.Run("множественные_пробелы", func(t *testing.T) {
		r := Check("string", "про  грамма", "про грамма")
		assertCorrect(t, r, true)
	})
	t.Run("разные_слова", func(t *testing.T) {
		r := Check("string", "программа", "код")
		assertCorrect(t, r, false)
	})
	t.Run("пустая_строка", func(t *testing.T) {
		r := Check("string", "", "")
		assertCorrect(t, r, true)
	})
	t.Run("не_нужен_python", func(t *testing.T) {
		r := Check("string", "ответ", "ответ")
		assertNeedsPython(t, r, false)
	})
	t.Run("кириллица_и_латиница", func(t *testing.T) {
		r := Check("string", "ответ", "otvet")
		assertCorrect(t, r, false)
	})
}

// ==================== MULTI ====================

func TestCheckMulti(t *testing.T) {
	t.Run("точное_совпадение", func(t *testing.T) {
		r := Check("multi", "А,Б,В", "А,Б,В")
		assertCorrect(t, r, true)
	})
	t.Run("разный_порядок", func(t *testing.T) {
		r := Check("multi", "А,Б,В", "В,Б,А")
		assertCorrect(t, r, true)
	})
	t.Run("точка_с_запятой", func(t *testing.T) {
		r := Check("multi", "А;Б;В", "А;Б;В")
		assertCorrect(t, r, true)
	})
	t.Run("разные_разделители", func(t *testing.T) {
		r := Check("multi", "А;Б;В", "А|Б|В")
		assertCorrect(t, r, true)
	})
	t.Run("не_полный_набор", func(t *testing.T) {
		r := Check("multi", "А,Б,В", "А,Б")
		assertCorrect(t, r, false)
	})
	t.Run("лишний_ответ", func(t *testing.T) {
		r := Check("multi", "А,Б", "А,Б,В")
		assertCorrect(t, r, false)
	})
	t.Run("один_элемент", func(t *testing.T) {
		r := Check("multi", "А", "А")
		assertCorrect(t, r, true)
	})
	t.Run("пробелы_вокруг", func(t *testing.T) {
		r := Check("multi", " А , Б ", "а,б")
		assertCorrect(t, r, true)
	})
	t.Run("не_нужен_python", func(t *testing.T) {
		r := Check("multi", "1,2", "1,2")
		assertNeedsPython(t, r, false)
	})
}

// ==================== CODE / TEXT ====================

func TestCheckCodeNeedsPython(t *testing.T) {
	t.Run("code_требует_python", func(t *testing.T) {
		r := Check("code", "print(42)", "print(42)")
		assertNeedsPython(t, r, true)
	})
	t.Run("text_требует_python", func(t *testing.T) {
		r := Check("text", "сложное задание", "ответ ученика")
		assertNeedsPython(t, r, true)
	})
	t.Run("code_correct_answer_сохраняется", func(t *testing.T) {
		r := Check("code", "print(42)", "print(0)")
		assertCorrectAnswer(t, r, "print(42)")
	})
	t.Run("code_correct_не_проверяется", func(t *testing.T) {
		r := Check("code", "print(42)", "print(42)")
		// До отправки в Python correct = false
		assertCorrect(t, r, false)
	})
}

// ==================== UNKNOWN TYPE ====================

func TestCheckUnknownType(t *testing.T) {
	t.Run("неизвестный_тип_как_string", func(t *testing.T) {
		r := Check("unknown_type", "ответ", "ответ")
		assertCorrect(t, r, true)
	})
	t.Run("неизвестный_тип_неправильный", func(t *testing.T) {
		r := Check("unknown_type", "ответ", "другое")
		assertCorrect(t, r, false)
	})
	t.Run("не_нужен_python", func(t *testing.T) {
		r := Check("unknown_type", "а", "а")
		assertNeedsPython(t, r, false)
	})
}

// ==================== Хелперы ====================

func assertCorrect(t *testing.T, r Result, want bool) {
	t.Helper()
	if r.Correct != want {
		t.Errorf("Correct = %v, want %v", r.Correct, want)
	}
}

func assertNeedsPython(t *testing.T, r Result, want bool) {
	t.Helper()
	if r.NeedsPython != want {
		t.Errorf("NeedsPython = %v, want %v", r.NeedsPython, want)
	}
}

func assertCorrectAnswer(t *testing.T, r Result, want string) {
	t.Helper()
	if r.CorrectAnswer != want {
		t.Errorf("CorrectAnswer = %q, want %q", r.CorrectAnswer, want)
	}
}
