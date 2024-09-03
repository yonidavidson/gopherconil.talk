package prompt

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type (
	// Message represents a message with a role and content.
	Message struct {
		Role    Role
		Content string
	}
	// Role represents the role of a message.
	Role string
)

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// ParseMessages transforms the prompt into a slice of messages.
func ParseMessages(input string, data any) ([]Message, error) {
	pt, err := parse(input, data)
	if err != nil {
		return nil, err
	}
	input = string(pt)
	// Validate tags before parsing
	if err := validate(pt); err != nil {
		return nil, err
	}
	var messages []Message

	// Regular expression to match tags and their content
	re := regexp.MustCompile(`<(system|user|assistant)>([\s\S]*?)</(system|user|assistant)>`)

	// Find all matches in the input string
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		role := Role(match[1])
		content := strings.TrimSpace(match[2])

		message := Message{
			Role:    role,
			Content: content,
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// validate checks if the input string has matching opening and closing tags for each role.
func validate(input []byte) error {
	roles := []string{"system", "user", "assistant"}
	for _, role := range roles {
		openCount := strings.Count(string(input), "<"+role+">")
		closeCount := strings.Count(string(input), "</"+role+">")
		if openCount != closeCount {
			return fmt.Errorf("mismatched tags for role %s: %d opening, %d closing", role, openCount, closeCount)
		}
	}
	return nil
}

func parse(promptTemplate string, data any) ([]byte, error) {
	tmpl, err := template.New("talk").Funcs(template.FuncMap{
		"limitTokens": limitTokens,
		"multiply": func(a, b float64) float64 {
			return a * b
		},
	}).Parse(promptTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}
	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}
	return []byte(result.String()), nil
}

func limitTokens(s string, maxTokens float64) string {
	const avgTokenLength = 4 // Average token length heuristic
	maxChars := int(maxTokens * avgTokenLength)

	if len(s) <= maxChars {
		return s
	}
	return s[:maxChars]
}
