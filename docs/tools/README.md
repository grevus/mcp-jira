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
| `incident_context` | 2 | planned | — | Jira + RAG + Confluence | — |
| `engineering_qa` | 2 | planned | — | RAG + Confluence | — |
| `ticket_triage` | 2 | planned | — | Jira + RAG | — |
| `release_risk_check` | 2 | planned | — | Jira + RAG | — |
| `runbook_for_signal` | 3 | planned | — | Confluence | — |
| `onboarding_path` | 3 | planned | — | Confluence + directory | — |
| `policy_guardrail_check` | 3 | planned | — | Policies RAG | — |

Phase 1 tools полностью реализованы в этой ветке. Phase 2/3 требуют новых источников данных и отдельного плана.

Шаблон для новых tools — [`_template.md`](_template.md).
