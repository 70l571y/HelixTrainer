// Package judge предоставляет функции для проверки решений челленджей.
package judge

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// CheckSolution проверяет решение пользователя против эталонного.
// Поддерживаемые режимы: exact, ignore_whitespace, ast, go_ast
func CheckSolution(userText, goalText, judgeMode string) bool {
	switch judgeMode {
	case "ast", "go_ast":
		return checkSolutionGoAST(userText, goalText)
	case "ignore_whitespace":
		return checkSolutionIgnoreWhitespace(userText, goalText)
	default: // "exact"
		return checkSolutionExact(userText, goalText)
	}
}

// checkSolutionExact выполняет точное посимвольное сравнение.
func checkSolutionExact(userText, goalText string) bool {
	return strings.TrimSpace(userText) == strings.TrimSpace(goalText)
}

// checkSolutionIgnoreWhitespace сравнивает без учёта пробелов и пустых строк.
func checkSolutionIgnoreWhitespace(userText, goalText string) bool {
	return normalizeCode(userText) == normalizeCode(goalText)
}

// normalizeCode нормализует код, удаляя лишние пробелы и пустые строки.
func normalizeCode(code string) string {
	lines := strings.Split(code, "\n")
	var normalized []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}

	return strings.Join(normalized, "\n")
}

// checkSolutionGoAST сравнивает код через AST Go.
func checkSolutionGoAST(userText, goalText string) bool {
	userNormalized, err := normalizeGoCode(userText)
	if err != nil {
		return false
	}

	goalNormalized, err := normalizeGoCode(goalText)
	if err != nil {
		return false
	}

	return userNormalized == goalNormalized
}

// normalizeGoCode нормализует Go код через форматирование.
func normalizeGoCode(code string) (string, error) {
	// Парсим код
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		return "", err
	}

	// Форматируем
	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, file)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateDiff генерирует unified diff между двумя текстами.
func GenerateDiff(userText, goalText string) string {
	userLines := strings.Split(userText, "\n")
	goalLines := strings.Split(goalText, "\n")

	var sb strings.Builder
	sb.WriteString("--- Your Solution\n")
	sb.WriteString("+++ Goal Solution\n")

	// Простая реализация diff
	maxLen := len(userLines)
	if len(goalLines) > maxLen {
		maxLen = len(goalLines)
	}

	for i := 0; i < maxLen; i++ {
		var userLine, goalLine string
		if i < len(userLines) {
			userLine = userLines[i]
		}
		if i < len(goalLines) {
			goalLine = goalLines[i]
		}

		if userLine != goalLine {
			if i < len(userLines) {
				sb.WriteString(fmt.Sprintf("-%s\n", userLine))
			}
			if i < len(goalLines) {
				sb.WriteString(fmt.Sprintf("+%s\n", goalLine))
			}
		} else {
			sb.WriteString(fmt.Sprintf(" %s\n", userLine))
		}
	}

	return sb.String()
}
