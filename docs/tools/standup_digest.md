# `standup_digest`

**Stability:** beta  
**Phase:** 1  
**Source:** `internal/handlers/standup_digest.go`

## Purpose

Собирает async standup дайджест за временной диапазон: группировка задач по done / in-progress / blocked / other.

## When to use

- Ежедневный async standup бота.
- Быстрый "что изменилось за вчера" на команду.

## When NOT to use

- Нужна только одна категория (например, blocked) → `list_issues` с фильтром.
- Свежие изменения нужны в реальном времени с webhook'ами (не поддерживается).

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `team_key` | string | yes | Jira project key (или команды, если проект = команда) | `"ABC"` |
| `from` | string | yes | Нижняя граница `updated`, `YYYY-MM-DD` или `YYYY-MM-DD HH:MM` | `"2026-04-07"` |
| `to` | string | yes | Верхняя граница | `"2026-04-08"` |
| `limit` | int | no | Максимум задач (default 25, max 100) | `50` |

## Output

| Field | Type | Description |
|---|---|---|
| `yesterday_summary` | string | Список done-задач, по строке на задачу, или `"(none)"` |
| `today_focus` | string | Список in-progress |
| `blockers` | `[]Issue` | Задачи с "blocked" в статусе |
| `notable_changes` | []string | Остальные движения (To Do, In Review и т.п.) |

## Example

```json
{
  "name": "standup_digest",
  "arguments": {"team_key": "ABC", "from": "2026-04-07", "to": "2026-04-08"}
}
```

## Errors

- `team_key is required`, `from and to are required`
- `jira: invalid date ...` при неправильном формате.
- Jira API errors.

## Data sources & freshness

- Jira live (JQL `project = X AND updated >= ... AND updated <= ...`).

## Cost & rate limits

- 1 Jira `/search/jql` вызов на запрос.

## Permissions

- Jira: read access к проекту.
- HTTP transport: `MCP_API_KEY`.

## Known limitations

- Нет фильтра по assignee / team field — `team_key` = project key.
- Maximum 100 задач (Jira `maxResults` cap).
- Категоризация "done/in progress/blocked" — по имени статуса, не по `statusCategory` (консистентность с MVP).
