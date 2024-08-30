package prompt_test

import (
	"reflect"
	"testing"

	"github.com/yonidavidson/gophercon-israel-2024/prompt"
)

func TestParseMessages(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []prompt.Message
	}{
		{
			name: "System and User messages",
			input: `[system]You are a helpful assistant[/system]
                    [user]Hello, how are you?[/user]`,
			expected: []prompt.Message{
				{Role: prompt.RoleSystem, Content: "You are a helpful assistant"},
				{Role: prompt.RoleUser, Content: "Hello, how are you?"},
			},
		},
		{
			name: "System, User, and Assistant messages",
			input: `[system]You are a helpful assistant[/system]
                    [user]What's the weather like?[/user]
                    [assistant]I'm sorry, I don't have real-time weather information. Could you please specify a location and I can provide general climate information?[/assistant]
                    [user]How about in New York?[/user]`,
			expected: []prompt.Message{
				{Role: prompt.RoleSystem, Content: "You are a helpful assistant"},
				{Role: prompt.RoleUser, Content: "What's the weather like?"},
				{Role: prompt.RoleAssistant, Content: "I'm sorry, I don't have real-time weather information. Could you please specify a location and I can provide general climate information?"},
				{Role: prompt.RoleUser, Content: "How about in New York?"},
			},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := prompt.ParseMessages(tt.input)
			if err != nil {
				t.Errorf("ParseMessages() error = %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseMessages() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseMessagesError(t *testing.T) {
	input := `<
[system]Incomplete system message
                    [user]User message without closing tag
                    [assistant]Assistant message[/assistant]`

	_, err := prompt.ParseMessages(input)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
