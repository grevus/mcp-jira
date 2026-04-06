# MCP Jira Core MVP — Implementation Plan (re-decomposed)

## Context

Существующий план `docs/superpowers/plans/2026-04-06-core-mvp.md` (22 таски, ~3500 строк) корректен по сути, но многие таски бундлят по 3-5 файлов и сотни строк кода в одном коммите. Цель этого плана — раздробить ту же самую работу на ~40 таск поменьше: каждая = один файл (или одна пара prod+test), один коммит, легко ревьюить диф. Архитектура и tech stack — без изменений (см. spec `docs/superpowers/specs/2026-04-06-mcp-jira-design.md`, source of truth).

**Final destination:** после approval план будет сохранён в `docs/superpowers/plans/2026-04-07-core-mvp-small-tasks.md` (или перезапишет старый — на усмотрение пользователя).

## Принципы декомпозиции

- **Один файл на таску** там, где это возможно. Test и prod могут идти вместе только когда это TDD-цикл одной маленькой функции.
- **Каждая таска = отдельный коммит**. Без «WIP, доделаю в следующей».
- **Никакого dead code между тасками**: handler ссылается на узкий интерфейс, который описан в той же таске; реальная имплементация интерфейса (`*jira.HTTPClient`) приходит в более ранней таске или передаётся как fake в тестах.
- **Wiring отделён от логики**: `register.AddTool` для каждого tool — отдельная мини-таска после хендлера, чтобы видеть в дифе именно регистрацию.

## File Structure

Идентичен spec §4 и существующему плану — без изменений. Перечисляю только новые/измененные группировки:

- `internal/jira/client.go` дробится на: `client_base.go` (struct + конструктор + `do()` helper), `client_issues.go`, `client_sprints.go`, `client_docs.go`.
- `internal/rag/store/postgres.go` дробится на: `postgres.go` (struct + `New` + `Close`), `postgres_upsert.go`, `postgres_query.go`, `postgres_stats.go`.
- `internal/rag/embed/voyage.go` и `openai.go` остаются по одному файлу (батчинг + единый HTTP-вызов — единое целое).
- `internal/register/register.go` дробится на `adapt.go` (generic adapter) и `register.go` (`Register(srv, jc, ret)`).

## Tasks

> Каждая таска ниже разворачивается в финальном плане в стандартный TDD-цикл: failing test → run → minimal impl → run → commit. Здесь даю заголовок, файлы, и одно-двух-строчное описание содержания. Полный код шагов появится в финальном плане под `docs/superpowers/plans/`.

### Phase 0 — Bootstrap

1. **Task 1: `go.mod` deps + `.gitignore`** — добавить modules: go-sdk/mcp, echo/v4, pgx/v5, goose/v3, testify, testcontainers-go (+ postgres module). Файлы: `go.mod`, `go.sum`, `.gitignore`. Sanity-проверка: `go build ./...` (пусто, ок).
2. **Task 2: `docker-compose.yml`** — один сервис `pgvector/pgvector:pg17`, порт 5432, volume, env согласован с дефолтным `DATABASE_URL`. Файл: `docker-compose.yml`.

### Phase 1 — Config

3. **Task 3: `config.Mode` + `Config` struct** — определить тип `Mode` (`stdio|http|index`), пустой `Config`. Файл: `internal/config/config.go`.
4. **Task 4: `config.Load(mode)` — Jira/DB env (общая часть)** — тест + код для required `JIRA_BASE_URL`, `JIRA_EMAIL`, `JIRA_API_TOKEN`, `DATABASE_URL`. Файлы: `internal/config/config.go`, `internal/config/config_test.go`.
5. **Task 5: `config.Load` — embedder выбор** — `RAG_EMBEDDER` default `voyage`; require соответствующий API-ключ. Тест: 4 кейса (default, voyage явный, openai, неизвестный).
6. **Task 6: `config.Load` — http-only env** — для `http`: required `MCP_API_KEY`, optional `MCP_ADDR` (default `:8080`). Для других режимов — игнорируется. Тест: stdio без `MCP_API_KEY` ок, http без — ошибка.

### Phase 2 — Jira: helpers и DTO

7. **Task 7: JQL `quoteJQL` (TDD)** — table-driven тесты сначала: пустая, кавычка, бэкслэш, нормальный ключ. Файлы: `internal/jira/jql.go`, `internal/jira/jql_test.go`.
8. **Task 8: JQL `validateProjectKey`** — whitelist `^[A-Z][A-Z0-9_]*$`. Те же файлы.
9. **Task 9: DTO types** — `Issue`, `SprintHealth`, `IssueDoc`, `ListIssuesParams`. Никаких методов, никаких импортов. Файл: `internal/jira/types.go`.

### Phase 3 — Jira HTTP client (по endpoint-у на таску)

10. **Task 10: `HTTPClient` baseline + `do()` helper** — конструктор `NewHTTPClient(baseURL, email, token, *http.Client)`, приватный `do(ctx, method, path, body) (*http.Response, error)` с basic auth и `Accept: application/json`. Тест на basic-auth-заголовок через `httptest.Server`. Файлы: `internal/jira/client_base.go`, `internal/jira/client_base_test.go`.
11. **Task 11: `ListIssues` happy path** — построение JQL из `ListIssuesParams`, GET `/rest/api/3/search/jql?fields=...&maxResults=...`, парсинг ответа. Тест на одном фикстурном JSON. Файлы: `internal/jira/client_issues.go`, `internal/jira/client_issues_test.go`.
12. **Task 12: `ListIssues` фильтры (status/assignee/limit)** — расширить тест (table-driven) и реализацию: добавление кляуз `AND status = ...`, `AND assignee = ...`, дефолт limit=25, max 100.
13. **Task 13: `ListIssues` ошибки** — non-2xx → `fmt.Errorf("jira: GET %s -> %d", path, status)`. Тест на 401 и 500.
14. **Task 14: `GetSprintHealth` — fetch active sprint** — GET `/rest/agile/1.0/board/{id}/sprint?state=active`, парсинг id+name. Тест с httptest. Файлы: `internal/jira/client_sprints.go`, `internal/jira/client_sprints_test.go`.
15. **Task 15: `GetSprintHealth` — aggregation** — GET `/rest/agile/1.0/sprint/{sprintId}/issue?fields=status`, агрегация `Total/Done/InProgress/Blocked`, `Velocity` (story points done). Тест с фикстурой на 5 issue.
16. **Task 16: `IterateIssueDocs` — pagination loop** — каркас канала, одна страница, без comments. Файлы: `internal/jira/client_docs.go`, `internal/jira/client_docs_test.go`.
17. **Task 17: `IterateIssueDocs` — comments fetch per issue** — для каждой задачи догрузить `/rest/api/3/issue/{key}/comment`. Расширить тест.
18. **Task 18: `IterateIssueDocs` — changelog → status history** — `expand=changelog`, парсинг переходов статусов, форматирование `"YYYY-MM-DD: From → To"`. Расширить тест фикстурой с 2-3 переходами.
19. **Task 19: `IterateIssueDocs` — linked issues + ошибки** — парсинг `issuelinks`, error-channel на http-фейлы, корректное закрытие каналов.

### Phase 4 — Auth

20. **Task 20: `auth.Middleware` (constant-time)** — TDD: 3 теста (valid → next, invalid → 401, missing → 401). Реализация — stdlib `func(http.Handler) http.Handler` с `subtle.ConstantTimeCompare`. Файлы: `internal/auth/apikey.go`, `internal/auth/apikey_test.go`.

### Phase 5 — Handlers (без MCP)

21. **Task 21: `Handler[In, Out]` тип** — один файл с одним типом, без зависимостей кроме `context`. Файл: `internal/handlers/handler.go`.
22. **Task 22: `handlers.ListIssues` — Input/Output + узкий интерфейс** — `IssueLister` интерфейс, `ListIssuesInput`/`ListIssuesOutput`, конструктор `ListIssues(l IssueLister) Handler[...]`. Тест на ручной fake. Файлы: `internal/handlers/issues.go`, `internal/handlers/issues_test.go`.
23. **Task 23: `handlers.SprintHealth`** — `SprintReader`, `SprintHealthInput`/`Output`, конструктор + тест с fake. Файлы: `internal/handlers/sprints.go`, `internal/handlers/sprints_test.go`.
24. **Task 24: `handlers.SearchKnowledge` — Input/Output** — типы и узкий интерфейс `KnowledgeRetriever`. (Реализация retriever-а — позже; в этой таске handler работает с fake.) Файлы: `internal/handlers/knowledge.go`, `internal/handlers/knowledge_test.go`.
25. **Task 25: `handlers.SearchKnowledge` — валидация top_k** — `top_k` default 5, max 20, error при >20. Расширить тест.

### Phase 6 — RAG: embed

26. **Task 26: `embed.Embedder` interface** — `Embed`, `Dimension`, `Name`. Файл: `internal/rag/embed/embedder.go`.
27. **Task 27: `VoyageEmbedder` — single batch happy path** — POST к `api.voyageai.com`, парсинг ответа. Тест на httptest с одним батчем. Файлы: `internal/rag/embed/voyage.go`, `internal/rag/embed/voyage_test.go`.
28. **Task 28: `VoyageEmbedder` — batching (128) и порядок** — расщепить вход на батчи по 128, склеить в правильном порядке. Тест на 200 текстов.
29. **Task 29: `VoyageEmbedder` — ошибки и dimension** — non-2xx → ошибка; `Dimension()` = 1024.
30. **Task 30: `OpenAIEmbedder` — happy path** — модель `text-embedding-3-small`, `dimensions=1024`. Тест с httptest. Файлы: `internal/rag/embed/openai.go`, `internal/rag/embed/openai_test.go`.
31. **Task 31: `OpenAIEmbedder` — batching (100) и ошибки** — аналогично Voyage.

### Phase 7 — RAG: store

32. **Task 32: `store` types** — `Document`, `Filter`, `Hit`, интерфейс `Store`. Файл: `internal/rag/store/store.go`.
33. **Task 33: миграции — `001_init.sql` + `embed.FS`** — `CREATE EXTENSION vector`, `CREATE TABLE issues_index(...)`, HNSW индекс. Файлы: `internal/rag/store/migrations/001_init.sql`, `internal/rag/store/migrations.go` (`embed.FS` + `Migrate(ctx, db)` через goose).
34. **Task 34: `PgvectorStore` baseline** — `New(ctx, dsn)` + `Close()`. Без CRUD. Файл: `internal/rag/store/postgres.go`.
35. **Task 35: testcontainers helper** — общий `setupPgvector(t)` под build tag `+integration`. Файл: `internal/rag/store/testhelp_test.go` (под `//go:build integration`).
36. **Task 36: `PgvectorStore.Upsert`** — TDD: интеграционный тест (insert + повторный insert обновляет). Файл: `internal/rag/store/postgres_upsert.go` + тест в `postgres_upsert_test.go`.
37. **Task 37: `PgvectorStore.Query`** — SQL с `<=>`, фильтр по project_key, excerpt = `substring(content,1,300)`. Интеграционный тест: вставить 3 doc, поискать.
38. **Task 38: `PgvectorStore.Stats`** — `SELECT count(*) WHERE project_key=$1`. Тест.

### Phase 8 — RAG: indexer и retriever

39. **Task 39: `Indexer.RenderDoc`** — pure function: `IssueDoc → string` по шаблону из spec §5.5. Чистый unit-тест на golden string. Файл: `internal/rag/index/render.go` + `render_test.go`.
40. **Task 40: `Indexer` struct + `Reindex` happy path** — `Reader+Embedder+Store`, последовательно: drain reader → embed → upsert. Без транзакции. Тест на fake-ах. Файлы: `internal/rag/index/indexer.go`, `internal/rag/index/indexer_test.go`.
41. **Task 41: `Indexer.Reindex` — батчинг + транзакционный TRUNCATE/INSERT** — добавить интерфейс `TxStore` (или метод `Store.WithinTx`), TRUNCATE по project_key + insert в одной транзакции. Расширить fake store с историей вызовов; добавить интеграционный тест против реального pgvector.
42. **Task 42: `Retriever.Search`** — `Embedder.Embed([query])` → `Store.Query`. Файлы: `internal/rag/retriever/retriever.go`, `internal/rag/retriever/retriever_test.go`.

### Phase 9 — Register (мост к MCP SDK)

43. **Task 43: `adapt[In, Out]` — success path** — generic функция, JSON marshal Out → TextContent + structured Out. Unit-тест без реального сервера. Файлы: `internal/register/adapt.go`, `internal/register/adapt_test.go`.
44. **Task 44: `adapt` — error path** — handler возвращает ошибку → `IsError: true`, `Content: TextContent{err.Error()}`, `error == nil`. Расширить тест.
45. **Task 45: `Register(srv, jc, ret)` — list_issues** — `mcp.AddTool(srv, &mcp.Tool{Name: "list_issues", ...}, adapt(handlers.ListIssues(jc)))`. Файл: `internal/register/register.go` + `register_test.go` (smoke: tool появился в `srv.Tools()` если SDK даёт API).
46. **Task 46: `Register` — get_sprint_health** — одна строка + расширение теста.
47. **Task 47: `Register` — search_jira_knowledge** — одна строка + расширение теста.

### Phase 10 — Бинари

48. **Task 48: `cmd/server/main.go` — флаги и `config.Load`** — флаг `--transport=stdio|http`, `config.Load(mode)`, конструирование `*jira.HTTPClient`, `embed.*Embedder`, `*store.PgvectorStore`, `*retriever.Retriever`, `mcp.NewServer + register.Register`. Без транспорта. Файл: `cmd/server/main.go`.
49. **Task 49: `cmd/server` — stdio транспорт** — `srv.Run(ctx, &mcp.StdioTransport{})`. Явный комментарий: «никакого stdout, только log (stderr)». `go build` smoke.
50. **Task 50: `cmd/server` — Streamable HTTP + Echo + auth** — Echo, `auth.Middleware` через `echo.WrapMiddleware`, `mcp.NewStreamableHTTPHandler` через `echo.WrapHandler` на `/mcp` (`ANY`). `go build` smoke.
51. **Task 51: `cmd/index/main.go` — каркас + `migrate` subcommand** — простой ручной парсинг (без cobra), `migrate` → `store.Migrate(ctx, db)`. Файл: `cmd/index/main.go`.
52. **Task 52: `cmd/index` — `index --project=KEY` subcommand** — собрать `Indexer`, вызвать `Reindex`, залогировать `indexed N docs in T`.

### Phase 11 — Документация

53. **Task 53: README — Quickstart + env матрица** — `docker compose up`, `go build`, env по mode, два примера запуска, Claude Desktop конфиг. Файл: `README.md`.
54. **Task 54: README — «Adding a new Jira tool»** — пошаговый рецепт из spec §10.
55. **Task 55: CLAUDE.md финализация** — обновить ссылки на новый план; убедиться, что secret-rules в актуальном состоянии. Файл: `CLAUDE.md`.

## Critical files to be modified

- `internal/handlers/*` — без импортов `mcp`/`echo`. Это инвариант, проверяется глазами на каждой PR-таске 21–25 и 45–47.
- `internal/register/{adapt,register}.go` — единственный мост к `go-sdk/mcp`. При апдейте SDK правки только здесь.
- `cmd/server/main.go` — единственное место с Echo. В stdio-ветке Echo вообще не создаётся.
- `internal/rag/store/migrations/001_init.sql` — `vector(1024)` фиксирован; смена размерности embedder-а ломает эту схему.

## Reuse / existing utilities

В репо ещё нет кода — переиспользовать нечего. Но для каждой таски, начиная с 10, проверяй существующий код перед написанием нового (например, тест-helper для `httptest.Server` из таски 10 переиспользуется в 11–19).

## Verification (end-to-end)

После таски 52 (всё собрано):

1. `docker compose up -d`
2. `go build -o bin/mcp-jira ./cmd/server && go build -o bin/mcp-jira-index ./cmd/index`
3. Экспортировать env (см. README).
4. `bin/mcp-jira-index migrate` → ноль ошибок, в БД появилась таблица `issues_index`.
5. `bin/mcp-jira-index index --project=<реальный ключ>` → лог `indexed N docs`.
6. `bin/mcp-jira --transport=stdio` под Claude Desktop → tools `list_issues`, `get_sprint_health`, `search_jira_knowledge` появляются и отвечают.
7. `MCP_API_KEY=secret bin/mcp-jira --transport=http` → `curl -X POST -H "X-API-Key: secret" localhost:8080/mcp ...` отвечает; без ключа → 401.
8. `go test ./...` зелёный.
9. `go test -tags=integration ./...` зелёный (требует Docker).

## Self-review checklist (для писателя финального плана)

- Каждая таска 21–25, 45–47 не импортирует `mcp`/`echo` (handlers) либо импортирует только в `register`/`cmd/server`.
- Имена методов на интерфейсах handler-ов совпадают с реальными методами `*jira.HTTPClient` (`ListIssues`, `GetSprintHealth`, `IterateIssueDocs`).
- Размерность 1024 встречается в: миграции, Voyage `Dimension()`, OpenAI `Dimension()` + `dimensions=1024` в HTTP-запросе, проверке в `cmd/index`.
- В stdio-ветке `cmd/server` нет ничего, что пишет в stdout.
- Все 22 таски старого плана покрыты (старая → новая):
  - Old 1 → New 1–2; Old 2 → New 3–6; Old 3 → New 7–8; Old 4 → New 9; Old 5 → New 10–13; Old 6 → New 14–15; Old 7 → New 16–19; Old 8 → New 20; Old 9 → New 21; Old 10 → New 22; Old 11 → New 23; Old 12 → New 26–29; Old 13 → New 30–31; Old 14 → New 32–33; Old 15 → New 34–38; Old 16 → New 39–41; Old 17 → New 42; Old 18 → New 24–25; Old 19 → New 43–47; Old 20 → New 48–50; Old 21 → New 51–52; Old 22 → New 53–55.
