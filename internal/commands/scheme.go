package commands

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
}

type StreamDelta struct {
	Content string `json:"content,omitempty"` // Use omitempty if content might be absent
	Role    string `json:"role,omitempty"`    // Role might appear in the first chunk
}

type StreamChoice struct {
	Delta        StreamDelta `json:"delta"`
	Index        int         `json:"index"`
	FinishReason *string     `json:"finish_reason"` // Pointer allows null/omitted
}

type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
	// Usage   UsageStats     `json:"usage,omitempty"` // Usage might appear at the end
}

type MessageEntry struct {
	Time    string  `json:"time"`
	Message Message `json:"message"`
}

func NewSystemMessage(content string) Message {
	return Message{Role: "system", Content: content}
}
func NewUserMessage(content string) Message {
	return Message{Role: "user", Content: content}
}
func NewAssistantMessage(content string) Message {
	return Message{Role: "assistant", Content: content}
}

func NewMessageEntry(time string, message Message) MessageEntry {
	return MessageEntry{
		Time:    time,
		Message: message,
	}
}
