package register

import (
	"context"

	"github.com/grevus/mcp-jira/internal/handlers"
	"github.com/grevus/mcp-jira/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// JiraClient объединяет узкие интерфейсы всех обработчиков, работающих с Jira.
// *jira.HTTPClient автоматически удовлетворяет этому интерфейсу.
type JiraClient interface {
	ListIssues(ctx context.Context, p jira.ListIssuesParams) ([]jira.Issue, error)
	GetIssue(ctx context.Context, key string) (jira.Issue, string, error)
	GetSprintHealth(ctx context.Context, boardID int) (jira.SprintHealth, error)
	GetSprintReport(ctx context.Context, boardID, sprintID int) (jira.SprintReport, error)
	GetSprintScopeChanges(ctx context.Context, sprintID int) (added, removed []string, err error)
	GetIssueComments(ctx context.Context, issueKey string) ([]string, error)
}

// Register регистрирует все MCP-инструменты в srv.
// jc — клиент Jira; ret — retriever, реализующий handlers.KnowledgeRetriever.
func Register(srv *mcp.Server, jc JiraClient, ret handlers.KnowledgeRetriever) {
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list_issues",
		Description: "Search Jira issues using JQL filters (project, status, assignee).",
	}, adapt(handlers.ListIssues(jc)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get_sprint_health",
		Description: "Return health metrics for the active sprint of a Jira Software board.",
	}, adapt(handlers.SprintHealth(jc)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "search_jira_knowledge",
		Description: "Semantic search over indexed Jira issues for a given project. Use for free-text questions when you don't have a specific issue key. Not a substitute for live Jira filters — data is as fresh as the last indexer run.",
	}, adapt(handlers.SearchKnowledge(ret)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "similar_issues",
		Description: "Find Jira issues semantically similar to a given one. Use for duplicate detection, incident correlation, known-fix discovery. Relies on the RAG index — issue must belong to an indexed project.",
	}, adapt(handlers.SimilarIssues(jc, ret)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "sprint_health_report",
		Description: "Return a deterministic health report for a Jira Software sprint: aggregates, blocked issues, risk level, action items, and scope_added/scope_removed issue keys (when sprint_id is provided, derived from changelog). If sprint_id is 0, the active sprint of the board is used.",
	}, adapt(handlers.SprintHealthReport(jc, jc)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "standup_digest",
		Description: "Produce an async standup digest for a team over a time window: grouped by done / in-progress / blocked / other. team_key is a Jira project key; from/to are dates (YYYY-MM-DD[ HH:MM]).",
	}, adapt(handlers.StandupDigest(jc)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "engineering_qa",
		Description: "Answer an engineering question by returning relevant Jira issues from the RAG index as citations for the LLM client. Jira-only until Phase 3 adds Confluence.",
	}, adapt(handlers.EngineeringQA(ret)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "incident_context",
		Description: "Gather incident context for a Jira issue: basic fields, similar past incidents from RAG, plus deterministically extracted suspected causes and recommended checks from the description and comments.",
	}, adapt(handlers.IncidentContext(jc, jc, ret)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "ticket_triage",
		Description: "Primary triage of an incoming Jira issue: finds similar issues via RAG and infers a suggested owning team (frequency of assignees) and priority (keyword heuristic). Deterministic, no LLM generation.",
	}, adapt(handlers.TicketTriage(jc, ret)))

	mcp.AddTool(srv, &mcp.Tool{
		Name:        "release_risk_check",
		Description: "Assess the risk of an upcoming release by fix_version: collects open and blocked Jira issues and semantically-similar postmortems from RAG, returning a deterministic risk level and short summary.",
	}, adapt(handlers.ReleaseRiskCheck(jc, ret)))
}
