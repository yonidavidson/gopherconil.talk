package prompt

import (
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
