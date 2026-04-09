# `sprint_health_report`

**Stability:** beta  
**Phase:** 1  
**Source:** `internal/handlers/sprint_report.go`

## Purpose

Детерминированный health-отчёт по спринту: агрегаты, список blocked задач, уровень риска и action items.

## When to use

- EM/PM хочет быстрый срез активного спринта.
- Перед standup — нужно видеть blocked и риск.

## When NOT to use

- Нужны просто базовые counters → `get_sprint_health`.
- Нужны scope changes (added/removed) — пока не реализовано (phase 2, анализ changelog).

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `board_id` | int | yes | ID Jira Software board | `42` |
| `sprint_id` | int | no | Конкретный sprint. `0`/omit → активный | `123` |

## Output

| Field | Type | Description |
|---|---|---|
| `report.health` | `SprintHealth` | Totals, done, in_progress, blocked, velocity |
| `report.blocked_issues` | `[]Issue` | Разобранные blocked задачи |
| `report.scope_added` / `scope_removed` | `[]Issue` | **пустые** (phase 2) |
| `summary` | string | Человекочитаемое описание |
| `risk_level` | `"low"` / `"medium"` / `"high"` | >20% blocked → high; >10% → medium |
| `action_items` | []string | Текст "Unblock KEY: summary" на каждую blocked |

## Example

```json
{"name": "sprint_health_report", "arguments": {"board_id": 42}}
```

## Errors

- `board_id is required`
- Jira Agile API errors (нет активного спринта, 404, 5xx).

## Data sources & freshness

- Jira Agile API v1.0 live: `/board/{id}/sprint?state=active`, `/sprint/{id}/issue`.

## Cost & rate limits

- 1–2 Jira Agile API вызова на запрос.

## Permissions

- Jira: Agile scope + read access к board.
- HTTP transport: `MCP_API_KEY`.

## Known limitations

- Scope changes (added/removed) не считаются — требует `expand=changelog` и анализа истории полей.
- `story points` берутся из `customfield_10016` (дефолт Jira Cloud); для других инстансов нужно параметризовать.
