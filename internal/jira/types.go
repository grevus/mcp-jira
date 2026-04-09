package jira

import "time"

type ListIssuesParams struct {
	ProjectKey  string
	Status      string
	Assignee    string
	FixVersion  string // опционально; подставляется в JQL `fixVersion = "..."`
	UpdatedFrom string // опционально; формат "YYYY-MM-DD" или "YYYY-MM-DD HH:MM"
	UpdatedTo   string // опционально; тот же формат
	Limit       int    // default 25, max 100 (применяется в client, не здесь)
}

type Issue struct {
	Key      string `json:"key"`
	Summary  string `json:"summary"`
	Status   string `json:"status"`
	Assignee string `json:"assignee,omitempty"`
}

type SprintHealth struct {
	BoardID    int     `json:"board_id"`
	SprintName string  `json:"sprint_name"`
	Total      int     `json:"total"`
	Done       int     `json:"done"`
	InProgress int     `json:"in_progress"`
	Blocked    int     `json:"blocked"`
	Velocity   float64 `json:"velocity"`
}

// SprintReport — расширенный health-отчёт для sprint_health_report tool.
// Дополняет SprintHealth списками blocked/scope и детерминированным summary.
type SprintReport struct {
	Health         SprintHealth `json:"health"`
	BlockedIssues  []Issue      `json:"blocked_issues"`
	ScopeAdded     []Issue      `json:"scope_added"`    // TODO(phase2): заполнять через changelog
	ScopeRemoved   []Issue      `json:"scope_removed"`  // TODO(phase2): заполнять через changelog
}

type IssueDoc struct {
	ProjectKey    string
	Key           string
	Summary       string
	Status        string
	Assignee      string
	Description   string
	Comments      []string  // плоский список текстов
	StatusHistory []string  // "2026-01-01: To Do → In Progress"
	LinkedIssues  []string  // ключи связанных
	UpdatedAt     time.Time
}
