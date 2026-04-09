# MCP Jira Tools — каталог

Канонический источник контрактов MCP tools. Описания в коде (`internal/register/register.go`) и здесь должны совпадать — при рассинхронизации правится и то и другое.

| Name | Phase | Stability | Transport | Data source | Doc |
|---|---|---|---|---|---|
| `list_issues` | MVP | stable | stdio/http | Jira live | [list_issues.md](list_issues.md) |
| `get_sprint_health` | MVP | stable | stdio/http | Jira Agile live | [get_sprint_health.md](get_sprint_health.md) |
| `search_jira_knowledge` | MVP | stable | stdio/http | pgvector RAG | [search_jira_knowledge.md](search_jira_knowledge.md) |
| `similar_issues` | 1 | beta | stdio/http | Jira live + RAG | [similar_issues.md](similar_issues.md) |
| `sprint_health_report` | 1 | beta | stdio/http | Jira Agile live | [sprint_health_report.md](sprint_health_report.md) |
| `standup_digest` | 1 | beta | stdio/http | Jira live | [standup_digest.md](standup_digest.md) |
| `incident_context` | 2 | beta | stdio/http | Jira live + RAG | [incident_context.md](incident_context.md) |
| `engineering_qa` | 2 | beta | stdio/http | pgvector RAG | [engineering_qa.md](engineering_qa.md) |
| `ticket_triage` | 2 | beta | stdio/http | Jira live + RAG | [ticket_triage.md](ticket_triage.md) |
| `release_risk_check` | 2 | beta | stdio/http | Jira live + RAG | [release_risk_check.md](release_risk_check.md) |
| `runbook_for_signal` | 3 | planned | — | Confluence | — |
| `onboarding_path` | 3 | planned | — | Confluence + directory | — |
| `policy_guardrail_check` | 3 | planned | — | Policies RAG | — |

Phase 1 и Phase 2 tools полностью реализованы в этой ветке (Jira-only RAG). Phase 3 требует Confluence-коннектора и отдельного плана.

Шаблон для новых tools — [`_template.md`](_template.md).
