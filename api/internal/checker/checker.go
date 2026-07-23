package checker

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type Result struct {
	Correct       bool   `json:"correct"`
	CorrectAnswer string `json:"correct_answer"`
	NeedsPython   bool   `json:"needs_python"`
}

// Check — проверяет ответ по типу задания
func Check(taskType, correctAnswer, userAnswer string) Result {
	switch taskType {
	case "choice":
		return checkChoice(correctAnswer, userAnswer)
	case "number":
		return checkNumber(correctAnswer, userAnswer)
	case "string":
		return checkString(correctAnswer, userAnswer)
	case "multi":
		return checkMulti(correctAnswer, userAnswer)
	case "code", "text":
		return Result{
			Correct:       false,
			CorrectAnswer: correctAnswer,
			NeedsPython:   true,
		}
	default:
		return checkString(correctAnswer, userAnswer)
	}
}

func checkChoice(correct, user string) Result {
	correct = strings.TrimSpace(correct)
	user = strings.TrimSpace(user)
	return Result{
		Correct:       strings.EqualFold(correct, user),
		CorrectAnswer: correct,
		NeedsPython:   false,
	}
}

func checkNumber(correct, user string) Result {
	correct = strings.TrimSpace(correct)
	user = strings.TrimSpace(user)

	correctVal, err1 := parseNumber(correct)
	userVal, err2 := parseNumber(user)

	if err1 != nil || err2 != nil {
		return Result{
			Correct:       strings.EqualFold(correct, user),
			CorrectAnswer: correct,
			NeedsPython:   false,
		}
	}

	tolerance := 0.01
	return Result{
		Correct:       math.Abs(correctVal-userVal) <= tolerance,
		CorrectAnswer: correct,
		NeedsPython:   false,
	}
}

func checkString(correct, user string) Result {
	correct = normalize(correct)
	user = normalize(user)
	return Result{
		Correct:       correct == user,
		CorrectAnswer: strings.TrimSpace(correct),
		NeedsPython:   false,
	}
}

func checkMulti(correct, user string) Result {
	correctSet := parseSet(correct)
	userSet := parseSet(user)

	if len(correctSet) != len(userSet) {
		return Result{
			Correct:       false,
			CorrectAnswer: correct,
			NeedsPython:   false,
		}
	}

	for i := range correctSet {
		if correctSet[i] != userSet[i] {
			return Result{
				Correct:       false,
				CorrectAnswer: correct,
				NeedsPython:   false,
			}
		}
	}

	return Result{
		Correct:       true,
		CorrectAnswer: correct,
		NeedsPython:   false,
	}
}

func parseNumber(s string) (float64, error) {
	s = strings.Replace(s, ",", ".", -1)
	s = strings.TrimSpace(s)
	var v float64
	_, err := fmt.Sscanf(s, "%f", &v)
	return v, err
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	for strings.Contains(s, "  ") {
		s = strings.Replace(s, "  ", " ", -1)
	}
	return s
}

// parseSet — парсит строку в отсортированный список (a,b,c | a;b;c | a|b|c)
func parseSet(s string) []string {
	var parts []string
	if strings.Contains(s, ",") {
		parts = strings.Split(s, ",")
	} else if strings.Contains(s, ";") {
		parts = strings.Split(s, ";")
	} else if strings.Contains(s, "|") {
		parts = strings.Split(s, "|")
	} else {
		parts = strings.Fields(s)
	}

	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			result = append(result, p)
		}
	}
	sort.Strings(result)
	return result
}
