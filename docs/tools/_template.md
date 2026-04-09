# `<tool_name>`

**Stability:** experimental | beta | stable  
**Phase:** N  
**Source:** `internal/handlers/<file>.go`

## Purpose

Одно предложение: что делает tool.

## When to use

- …

## When NOT to use

- …

## Input

| Field | Type | Required | Description | Example |
|---|---|---|---|---|
| `field` | string | yes | … | `"ABC"` |

## Output

| Field | Type | Description |
|---|---|---|

## Example

```json
// request
{"name": "<tool_name>", "arguments": {}}
```

```json
// response
{}
```

## Errors

- `<message>` — когда возникает.

## Data sources & freshness

- …

## Cost & rate limits

- …

## Permissions

- Jira scopes: …
- HTTP transport: требуется `MCP_API_KEY`.

## Known limitations

- …
