# `similar_issues`

**Stability:** beta  
**Phase:** 1  
**Source:** `internal/handlers/similar_issues.go`

## Purpose

Находит Jira issues, семантически похожие на заданную задачу — для обнаружения дубликатов, корреляции инцидентов, поиска known-fix.

## When to use

- Нужно найти "похожие баги" для triage.
- Коррелируем текущий инцидент с прошлыми.
- Ищем, не решалась ли эта проблема раньше.

## When NOT to use

- Есть ключевой поисковой запрос без исходной задачи → используй `search_jira_knowledge`.
- Проект не индексирован в pgvector → tool вернёт пустой список.

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `project_key` | string | yes | Ключ проекта Jira, в рамках которого искать | `"ABC"` |
| `issue_key` | string | yes | Источник: issue, к которой ищем похожие | `"ABC-42"` |
| `top_k` | int | no | Количество результатов (1–20, default 5) | `5` |

## Output

| Field | Type | Description |
|---|---|---|
| `source` | `Issue` | Базовые поля исходной задачи |
| `similar_issues` | `[]Hit` | Список похожих с `issue_key`, `summary`, `score`, `excerpt`. Исходная задача отфильтрована. |

## Example

```json
{"name": "similar_issues", "arguments": {"project_key": "ABC", "issue_key": "ABC-42", "top_k": 5}}
```

## Errors

- `issue_key is required` / `project_key is required`
- `top_k must be <= 20`
- Пропагируются ошибки Jira (`GET /rest/api/3/issue/...`) и retriever.

## Data sources & freshness

- `GetIssue` — Jira live.
- `Search` — pgvector, свежесть = последняя `bin/mcp-jira-index index --project=ABC`.

## Cost & rate limits

- 1 Jira API call + 1 embedding request + 1 pgvector query на вызов.

## Permissions

- Jira: read access к проекту.
- HTTP transport: `MCP_API_KEY` обязателен.

## Known limitations

- Работает только с уже индексированными проектами.
- Сходство — bi-encoder embedding; hybrid search / rerank не поддерживаются.
