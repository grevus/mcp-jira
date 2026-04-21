# `release_risk_check`

**Stability:** beta
**Phase:** 2
**Source:** `internal/handlers/release_risk_check.go`

## Purpose

Оценивает риск предстоящего релиза по `fixVersion`: собирает открытые и заблокированные задачи из Jira и семантически близкие постмортемы из RAG, выдавая детерминированный `risk_level` и краткое summary.

## When to use

- Перед cut'ом релиза: понять, что ещё висит и есть ли похожие инциденты в истории.
- Go/no-go meeting — быстрая сводка по версии.

## When NOT to use

- Нет зафиксированной `fixVersion` на задачах — результат будет пустым.
- Нужна детальная картина спринта → `sprint_health_report`.

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `project_key` | string | yes | Ключ проекта Jira | `"ABC"` |
| `fix_version` | string | yes | Значение fixVersion | `"1.42.0"` |
| `services_involved` | []string | no | Сервисы для усиления семантического запроса | `["payments","auth"]` |
| `top_k` | int | no | Сколько постмортемов вернуть (1–20, default 5) | `5` |

## Output

| Field | Type | Description |
|---|---|---|
| `fix_version` | string | Эхо входа |
| `open_issues` | `[]Issue` | Не Done и не Blocked |
| `blocked_issues` | `[]Issue` | Status содержит `block` (case-insensitive) |
| `related_postmortems` | `[]Hit` | RAG-хиты по запросу `postmortem incident <fix_version> <services>` |
| `risk_level` | string | `low` / `medium` / `high` (см. ниже) |
| `missing_runbooks` | []string | Всегда пустой массив в Phase 2 |
| `summary` | string | Одна строка с агрегатами |

## Risk rules (deterministic)

- `high`: `len(blocked) >= 3` **или** `len(related_postmortems) >= 3`.
- `medium`: `len(blocked) >= 1` **или** `len(open) >= 10`.
- иначе `low`.

## Example

```json
{"name": "release_risk_check", "arguments": {"project_key": "ABC", "fix_version": "1.42.0", "services_involved": ["payments"], "top_k": 5}}
```

## Errors

- `fix_version is required` / `project_key is required`
- `top_k must be <= 20`
- Пропагируются ошибки Jira `ListIssues` и RAG `Search`.

## Data sources & freshness

- Jira `/rest/api/3/search/jql` — live.
- pgvector — свежесть равна последней `bin/mcp-issues-index index --project=...`.

## Cost & rate limits

- 1 Jira search + 1 embedding + 1 pgvector query.

## Permissions

- Jira: read access к проекту.
- HTTP transport: `MCP_API_KEY` обязателен.

## Known limitations

- `missing_runbooks` всегда пуст — заполнится в Phase 3 после появления Confluence-коннектора.
- Без changelog-анализа (scope changes), без rerank, без hybrid search.
- Лимит выгрузки — 100 issues на fixVersion.
