package config

import (
	"fmt"
	"os"
)

// Mode определяет режим запуска сервера.
type Mode string

const (
	ModeStdio Mode = "stdio"
	ModeHTTP  Mode = "http"
	ModeIndex Mode = "index"
)

// Config содержит всю конфигурацию приложения.
type Config struct {
	Mode         Mode
	JiraBaseURL  string
	JiraEmail    string
	JiraAPIToken string
	DatabaseURL  string
	RAGEmbedder  string // "voyage" или "openai"
	VoyageAPIKey string
	OpenAIAPIKey string
	MCPAPIKey    string // только для http
	MCPAddr      string // только для http, default ":8080"
}

// Load читает переменные окружения и возвращает Config для указанного mode.
// Общие обязательные переменные: JIRA_BASE_URL, JIRA_EMAIL, JIRA_API_TOKEN, DATABASE_URL.
func Load(mode Mode) (*Config, error) {
	required := []struct {
		env   string
		value *string
	}{
		{"JIRA_BASE_URL", nil},
		{"JIRA_EMAIL", nil},
		{"JIRA_API_TOKEN", nil},
		{"DATABASE_URL", nil},
	}

	values := make([]string, len(required))
	for i, r := range required {
		v := os.Getenv(r.env)
		if v == "" {
			return nil, fmt.Errorf("config: %s is required", r.env)
		}
		values[i] = v
	}

	return &Config{
		Mode:         mode,
		JiraBaseURL:  values[0],
		JiraEmail:    values[1],
		JiraAPIToken: values[2],
		DatabaseURL:  values[3],
	}, nil
}
