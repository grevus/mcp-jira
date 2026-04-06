package config

// Mode определяет режим запуска сервера.
type Mode string

const (
	ModeStdio Mode = "stdio"
	ModeHTTP  Mode = "http"
	ModeIndex Mode = "index"
)

// Config содержит всю конфигурацию приложения.
// Конструктор Load появится в Task 4.
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
