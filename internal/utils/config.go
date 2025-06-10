// go-backend\internal\utils\config.go
package utils

type PromptUpdateBody struct {
	Key    string `json:"key"`    // ex: ai_prompt.line
	Model  string `json:"model"`  // ex: gpt-4o
	Prompt string `json:"prompt"` // ex: actual prompt text
}
