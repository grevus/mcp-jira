# `engineering_qa`

**Stability:** beta  
**Phase:** 2  
**Source:** `internal/handlers/engineering_qa.go`

## Purpose

Отвечает на инженерный вопрос, возвращая релевантные Jira issues из RAG-индекса как цитаты для LLM-клиента.

## When to use

- Нужно быстро найти исторический контекст по техническому вопросу.
- Ищем прецеденты решения похожей инженерной задачи.
- LLM-клиент собирает ответ с опорой на конкретные issue.

## When NOT to use

- Нужен актуальный источник истины из живой документации — этот tool не заменяет Confluence/runbooks.
- Confluence-коннектор ещё не подключён (Phase 3), поэтому ответы ограничены только Jira issue.
- Проект не индексирован в pgvector → вернётся пустой список.

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `question` | string | yes | Инженерный вопрос в свободной форме | `"How do we rotate DB creds?"` |
| `project_key` | string | yes | Ключ проекта Jira, в котором искать | `"ABC"` |
| `context_hint` | string | no | Доп. ключевые слова, конкатенируются к запросу | `"postgres vault"` |
| `top_k` | int | no | Количество цитат (1–20, default 5) | `5` |

## Output

| Field | Type | Description |
|---|---|---|
| `citations` | `[]Hit` | Релевантные issues с `issue_key`, `summary`, `status`, `score`, `excerpt`. |

## Example

```json
{"name": "engineering_qa", "arguments": {"question": "How do we rotate DB creds?", "project_key": "ABC", "context_hint": "postgres vault", "top_k": 5}}
```

## Errors

- `question is required`
- `project_key is required`
- `top_k must be <= 20`
- Пропагируются ошибки retriever (embedder / pgvector).

## Data sources & freshness

- `Search` — pgvector RAG-индекс, только Jira issues. Свежесть = последняя `bin/mcp-jira-index index --project=ABC`.

## Cost & rate limits

- 1 embedding request + 1 pgvector query на вызов.

## Permissions

- Jira: read access к проекту (на этапе индексации).
- HTTP transport: `MCP_API_KEY` обязателен.

## Known limitations

- Работает только с уже индексированными проектами.
- `docs_links` пустой до Phase 3 (Confluence-коннектор). Поле `doc_hits` добавится под флагом `RAG_DOCS_ENABLED`.
- Детерминированный bi-encoder retrieval; hybrid search / rerank не поддерживаются.
- Handler не генерирует текст ответа — LLM-клиент формулирует его сам поверх цитат.
