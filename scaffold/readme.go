package scaffold

import (
	"log"
	"os"
)

func GenerateReadmeFile(projectPath string) {
	f, err := os.Create(projectPath + "/goProto.md")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	content := `# GoProto

A personal CLI scaffolding tool for Go projects. GoProto automates the tedious parts of starting a new Go project — folder structure, framework wiring, database setup, migrations, and query generation — so you can get straight to building.

---

## What It Does

- Scaffolds a Go project structure inside your **current directory** (no forced subdirectory)
- Wraps [Go Blueprint](https://github.com/Melkeydev/go-blueprint) to generate the base project
- Wires up [sqlc](https://sqlc.dev/) for type-safe SQL query generation
- Sets up [Goose](https://github.com/pressly/goose) for database migrations
- Supports adding pre-built code templates (auth, etc.) to your project on demand

---

## Project Types

| Type | Description |
|------|-------------|
| ` + "`api`" + ` | REST API project with framework + database setup |
| ` + "`normal`" + ` | General Go project with standard structure (no database) — *coming soon* |

---

## Supported Options

**Frameworks**

- Gin
- Chi
- ` + "`net/http`" + ` (standard library)

**Databases**

- PostgreSQL (via pgx)
- SQLite

---

## Prerequisites

Make sure the following are installed before using GoProto:

- [Go](https://go.dev/dl/) 1.21+
- [Go Blueprint](https://github.com/Melkeydev/go-blueprint) — ` + "`go install github.com/melkeydev/go-blueprint@latest`" + `
- [Goose](https://github.com/pressly/goose) — ` + "`brew install goose`" + ` or ` + "`go install github.com/pressly/goose/v3/cmd/goose@latest`" + `
- [sqlc](https://sqlc.dev/) — ` + "`brew install sqlc`" + ` or ` + "`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`" + `

---

## Installation

` + "```bash" + `
go install github.com/yourusername/goproto@latest
` + "```" + `

> Replace ` + "`yourusername`" + ` with your actual GitHub username.

---

## Usage

1. Create your project folder and navigate into it:

` + "```bash" + `
mkdir my-project && cd my-project
` + "```" + `

2. Run GoProto:

` + "```bash" + `
goproto
` + "```" + `

3. Answer the prompts:
   - Project type (` + "`api`" + ` or ` + "`normal`" + `)
   - Framework (` + "`gin`" + `, ` + "`chi`" + `, ` + "`net/http`" + `)
   - Database (` + "`postgres`" + `, ` + "`sqlite`" + `)

GoProto will scaffold the full project structure, configure sqlc, and set up your migrations directory.

---

## Post-Generation Steps

After scaffolding, you'll still need to run a few commands manually:

` + "```bash" + `
# Generate type-safe query code from your SQL
sqlc generate

# Create your first migration
goose create init_schema sql -dir ./migrations

# Run migrations (see your Makefile for shortcuts)
make migrate
` + "```" + `

---

## Templates *(WIP)*

GoProto will support adding pre-built code snippets to an existing project:

` + "```bash" + `
goproto add auth
` + "```" + `

This will copy the relevant template files directly into your project. GoProto ships with built-in templates, with plans to support custom user-defined templates in the future.

---

## Roadmap

- [x] In-place project scaffolding (no forced subdirectory)
- [x] Go Blueprint wrapper with framework + database selection
- [x] sqlc configuration generation
- [ ] Goose migration setup
- [ ] Makefile with common commands
- [ ] ` + "`normal`" + ` project type
- [ ] Template system (` + "`goproto add <template>`" + `)
- [ ] Custom user templates

---

## Notes

GoProto is a personal tool built for my own workflow. It's opinionated by design — it only supports the tools and patterns I actually use. If it happens to work for you too, great.
`

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
}
