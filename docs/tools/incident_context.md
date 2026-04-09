# `incident_context`

**Stability:** beta  
**Phase:** 2  
**Source:** `internal/handlers/incident_context.go`

## Purpose

Собирает контекст инцидента по Jira-задаче: базовые поля, похожие прошлые инциденты из RAG, а также детерминированно извлечённые из описания и комментариев гипотезы причин и рекомендуемые проверки.

## When to use

- Начало triage по инциденту — быстрый срез "что известно + что было похожее".
- On-call хочет понять, не было ли такого раньше.
- Нужна отправная точка для runbook, без LLM-саммаризации.

## When NOT to use

- Нужна подробная документация/runbook — используй `search_jira_knowledge` (и `runbook_for_signal` в Phase 3).
- Проект не проиндексирован — `related_incidents` будет пустым.

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `issue_key` | string | yes | Исходная задача инцидента | `"ABC-42"` |
| `project_key` | string | yes | Проект для RAG-поиска похожих | `"ABC"` |
| `top_k` | int | no | Количество похожих (1–20, default 5) | `5` |

## Output

| Field | Type | Description |
|---|---|---|
| `source` | `Issue` | Базовые поля исходной задачи |
| `related_incidents` | `[]Hit` | Похожие задачи из того же проекта, исходная отфильтрована |
| `suspected_causes` | `[]string` | Предложения из описания/комментариев с маркерами "caused by", "root cause", "due to", "because of" (cap 5) |
| `recommended_checks` | `[]string` | Предложения с маркерами "check", "verify", "rollback", "restart", "monitor" (cap 5) |
| `docs_links` | `[]string` | Всегда пустой в Phase 2 — зарезервировано под Confluence-коннектор |

## Example

```json
{"name": "incident_context", "arguments": {"issue_key": "ABC-42", "project_key": "ABC", "top_k": 5}}
```

## Errors

- `issue_key is required` / `project_key is required`
- `top_k must be <= 20`
- Пропагируются ошибки `GetIssue`, `GetIssueComments`, retriever.

## Data sources & freshness

- `GetIssue`, `GetIssueComments` — Jira live.
- `Search` — pgvector, свежесть = последняя `bin/mcp-jira-index index --project=ABC`.

## Cost & rate limits

- 2 Jira API calls + 1 embedding request + 1 pgvector query на вызов.

## Permissions

- Jira: read access к проекту и комментариям.
- HTTP transport: `MCP_API_KEY` обязателен.

## Known limitations

- **Jira-only.** Пока нет Confluence-коннектора (Phase 3), `docs_links` всегда пуст.
- Извлечение причин/проверок — строго детерминированные keyword-правила, без LLM. Шум возможен, но воспроизводимо.
- Split предложений простой: по `. ! ? \n`. Сложные форматы (маркдаун-таблицы, код) могут давать неидеальные фрагменты.
