# `ticket_triage`

**Stability:** beta  
**Phase:** 2  
**Source:** `internal/handlers/ticket_triage.go`

## Purpose

Первичный triage входящей Jira issue: подбирает похожие задачи через RAG и на их основе предлагает команду-исполнителя и приоритет. Детерминированный, без LLM-генерации.

## When to use

- Новый тикет упал в инбокс, нужно быстро понять «куда отдать» и «насколько горит».
- Хочется увидеть родственные задачи в одном вызове вместе с решением по приоритету.

## When NOT to use

- Нужны только похожие задачи без решения по команде/приоритету → `similar_issues`.
- Проект не индексирован → SuggestedTeam будет пустой, Priority — по ключевым словам.
- Нужно человеко-читаемое обоснование с нюансами → этот tool возвращает короткий детерминированный rationale.

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `issue_key` | string | yes | Источник triage | `"ABC-42"` |
| `project_key` | string | yes | Проект для RAG-поиска | `"ABC"` |
| `top_k` | int | no | Количество похожих задач (1–20, default 10) | `10` |

## Output

| Field | Type | Description |
|---|---|---|
| `source` | `Issue` | Базовые поля исходной задачи |
| `suggested_team` | string | Наиболее частый assignee в похожих задачах; `""`, если assignees не найдены |
| `priority` | string | `"high"` / `"medium"` / `"low"` — по keyword-эвристике |
| `rationale` | string | Краткое детерминированное объяснение обоих решений |
| `similar_issues` | `[]Hit` | Похожие задачи (self-match отфильтрован) |

## Heuristics

### Priority (первый матч выигрывает)

| Level | Keywords (lowercase substring) |
|---|---|
| `high` | `outage`, `prod down`, `sev1`, `sev-1`, `blocker`, `critical`, `p0` |
| `medium` | `prod`, `production`, `customer`, `regression`, `p1` |
| `low` | всё остальное |

Проверка по `strings.ToLower(summary + " " + description)`.

### Suggested team

1. Собирается список assignees из excerpt'ов похожих задач (поле `Assignee:` в rendered-doc).
2. Самый частый побеждает; тай-брейк — порядок первого появления.
3. Если assignees не найдены — `""` и rationale `Team not inferred (no assignees in similar issues)`.

## Example

```json
{"name": "ticket_triage", "arguments": {"issue_key": "ABC-42", "project_key": "ABC", "top_k": 10}}
```

## Errors

- `issue_key is required` / `project_key is required`
- `top_k must be <= 20`
- Пропагируются ошибки `GetIssue` и retriever.

## Data sources & freshness

- `GetIssue` — Jira live.
- Похожие задачи — pgvector, свежесть = последняя `bin/mcp-jira-index index --project=ABC`.

## Cost & rate limits

- 1 Jira API call + 1 embedding + 1 pgvector query на вызов.

## Permissions

- Jira: read access к проекту.
- HTTP transport: `MCP_API_KEY` обязателен.

## Known limitations

- Keyword-эвристика приоритета — **только английский**, русские аналоги не покрываются.
- SuggestedTeam ограничен assignees, которые присутствуют в индексированных похожих задачах; если проект не проиндексирован или у похожих нет assignee — поле пустое.
- Без LLM-генерации: rationale — фиксированный шаблон.
- Без анализа labels / components — только assignee frequency.
