package config_test

import (
	"testing"

	"github.com/grevus/mcp-jira/internal/config"
	"github.com/stretchr/testify/require"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("JIRA_BASE_URL", "https://example.atlassian.net")
	t.Setenv("JIRA_EMAIL", "user@example.com")
	t.Setenv("JIRA_API_TOKEN", "secret-token")
	t.Setenv("DATABASE_URL", "postgres://localhost/testdb")
}

func TestLoad_HappyPath(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := config.Load(config.ModeStdio)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, config.ModeStdio, cfg.Mode)
	require.Equal(t, "https://example.atlassian.net", cfg.JiraBaseURL)
	require.Equal(t, "user@example.com", cfg.JiraEmail)
	require.Equal(t, "secret-token", cfg.JiraAPIToken)
	require.Equal(t, "postgres://localhost/testdb", cfg.DatabaseURL)
}

func TestLoad_MissingJiraBaseURL(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("JIRA_BASE_URL", "")

	cfg, err := config.Load(config.ModeStdio)
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "JIRA_BASE_URL")
}

func TestLoad_MissingJiraEmail(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("JIRA_EMAIL", "")

	cfg, err := config.Load(config.ModeStdio)
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "JIRA_EMAIL")
}

func TestLoad_MissingJiraAPIToken(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("JIRA_API_TOKEN", "")

	cfg, err := config.Load(config.ModeStdio)
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "JIRA_API_TOKEN")
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DATABASE_URL", "")

	cfg, err := config.Load(config.ModeStdio)
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "DATABASE_URL")
}
