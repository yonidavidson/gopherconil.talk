package prompt

import (
	"fmt"
	"regexp"
	"strings"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message represents a message with a role and content.
type Message struct {
	Role    Role
	Content string
}

// ParseMessages parses the input string into a slice of messages.
func ParseMessages(input string) ([]Message, error) {
	// Validate tags before parsing
	if err := validate(input); err != nil {
		return nil, err
	}
	var messages []Message

	// Regular expression to match tags and their content
	re := regexp.MustCompile(`\[(system|user|assistant)]([\s\S]*?)\[/(system|user|assistant)]`)

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
func validate(input string) error {
	roles := []string{"system", "user", "assistant"}
	for _, role := range roles {
		openCount := strings.Count(input, "["+role+"]")
		closeCount := strings.Count(input, "[/"+role+"]")
		if openCount != closeCount {
			return fmt.Errorf("mismatched tags for role %s: %d opening, %d closing", role, openCount, closeCount)
		}
	}
	return nil
}
