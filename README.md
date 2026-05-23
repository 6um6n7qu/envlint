# envlint

Validates `.env` files against a schema definition to catch missing or malformed variables before deploy.

---

## Installation

```bash
go install github.com/yourname/envlint@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envlint.git && cd envlint && go build ./...
```

---

## Usage

Define a schema file (`.env.schema`):

```ini
DATABASE_URL=required,url
PORT=required,number
DEBUG=optional,bool
APP_NAME=required
```

Then validate your `.env` file against it:

```bash
envlint --schema .env.schema --env .env
```

**Example output:**

```
✗ DATABASE_URL: missing required variable
✗ PORT: value "abc" is not a valid number
✓ DEBUG: ok
✓ APP_NAME: ok

2 error(s) found.
```

Exit code is non-zero when validation fails, making it easy to integrate into CI pipelines:

```yaml
- name: Validate env
  run: envlint --schema .env.schema --env .env
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--schema` | `.env.schema` | Path to the schema file |
| `--env` | `.env` | Path to the env file to validate |
| `--quiet` | `false` | Suppress output, use exit code only |

---

## License

MIT © yourname