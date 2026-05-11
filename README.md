> 🌐 **English** | **[한국어](README.ko.md)**

# LOADSTAR CLI

A Go-based LOADSTAR metadata management tool. Manage project work units (WayPoints), indexes (Maps), TODOs, and logs shared between AI agents and humans — all from the command line.

> 📌 New to LOADSTAR? Start with the [openLoadstar overview](https://github.com/openLoadstar/openLoadstar) and the [methodology spec](https://github.com/openLoadstar/spec).

---

## 🛠️ Installation

### Prerequisites

- Go 1.21 or later

### Build

```bash
git clone https://github.com/openLoadstar/cli.git
cd cli
go build -o bin/loadstar.exe .
```

Add the built binary (`bin/loadstar.exe`) to your PATH or invoke it by absolute path.

### Quick verification

```bash
loadstar --help
```

---

## 📋 Commands

| Command | Purpose |
|:---|:---|
| `loadstar init` | Initialize the `.loadstar/` directory structure |
| `loadstar show [FILTER] [--recent]` | List WayPoints (address · STATUS · LAST_MODIFIED) — keyword filter + sort by most recently modified |
| `loadstar validate` | Verify referential integrity across all elements, report broken links |
| `loadstar log [TIME_RANGE] [FILTER]` | Search change log — time range like `7d`, `3h` + keyword filter |
| `loadstar log add <ADDR> <KIND> "<MSG>"` | Directly add a log entry |
| `loadstar todo sync` | Auto-sync TODOs based on WayPoint STATUS |
| `loadstar todo list` | Current PENDING / ACTIVE / BLOCKED task list |
| `loadstar todo history [MAP_ADDR]` | History of completed TECH_SPEC items |
| `loadstar question [FILTER] [--with-resolved]` | Query unresolved OPEN_QUESTIONS |
| `loadstar question done <ADDR> <QID>` | Transition a RESOLVED question to DONE |
| `loadstar question close <ADDR> <QID> [reason]` | Close a question directly without creating a Decision file |
| `loadstar question stats` | Aggregate OPEN / DEFERRED / RESOLVED / DONE counts |

> Run `--help` on any command for detailed options.

---

## 🚀 Quick Start

### Introducing LOADSTAR into a new project

```bash
cd /my/project

# 1. Initialize the metadata directory
loadstar init

# 2. Create the first WayPoint & Map via AI, or add directly to .loadstar/WAYPOINT/
#    (see "AI Session Entry Prompt" in the openLoadstar README for step-by-step instructions)

# 3. Check current state
loadstar show           # all WayPoints
loadstar show --recent  # sorted by most recently modified
loadstar show frontend  # keyword filter: "frontend"

# 4. Validate references
loadstar validate
```

### Day-to-day metadata operations

```bash
# Current task list
loadstar todo list

# Sync TODOs after changing WayPoint STATUS
loadstar todo sync

# Browse completion history
loadstar todo history
loadstar todo history M://root/cli

# Search change log
loadstar log 7d                    # last 7 days
loadstar log cmd_show              # keyword filter
loadstar log 2d ISSUE              # time range + KIND filter

# Check unresolved questions requiring a human decision
loadstar question
loadstar question --with-resolved  # include resolved items
```

### Recording a meta event directly

```bash
loadstar log add W://root/cli/cmd_show MODIFIED "added --recent flag to show command"
```

---

## 🧭 Address Convention

```
M://root/cli            →  .loadstar/MAP/root.cli.md
W://root/cli/cmd_show   →  .loadstar/WAYPOINT/root.cli.cmd_show.md
```

- **M (Map)**: An index for grouping WayPoints — no STATUS, represents hierarchy only. Carries a **GOAL** slot for the high-level intent of the sub-tree.
- **W (WayPoint)**: The smallest work unit — composed of IDENTITY / CONNECTIONS / CODE_MAP / TODO / ISSUE. Carries a **GOAL** slot for the intent this WayPoint fulfils.
- **D (Data WayPoint / dwp)**: Metadata record for a data artifact (config, reference data) that changes independently of code — stored in `.loadstar/DATA_WAYPOINT/`.

---

## 📂 Directory Structure

```
.loadstar/
├── MAP/              M:// elements (Markdown)
├── WAYPOINT/         W:// elements (Markdown)
├── DATA_WAYPOINT/    D:// elements — data artifact metadata (Markdown)
├── DECISIONS/        OPEN_QUESTIONS decision records (ADR)
├── COMMON/           Project settings
└── .clionly/         ⚠️ CLI-only — do not edit directly (AI or human)
    ├── LOG/          Change history log
    └── TODO/         TODO_LIST · WP_SNAPSHOT (managed by sync)
```

> Directly editing `.clionly/` permanently breaks consistency between LOG and the actual metadata state.

---

## 🤖 AI Collaboration Workflow

1. **Session start** — AI loads `LOADSTAR_INIT.md` and the SPEC, then runs `loadstar show` / `loadstar todo list` / `loadstar question` to understand the current state.
2. **Before modifying code** — Register `- [ ] task description` in the target WayPoint's TODO section.
3. **After modification** — Check it off as `- [x] YYYY-MM-DD task description`.
4. **WayPoint fully complete** — Change STATUS from `S_PRG → S_STB`.
5. **Sync TODOs** — Run `loadstar todo sync` to auto-update TODO_LIST from WP STATUS.
6. **Validate** — Run `loadstar validate` before finishing to confirm no broken references (also detects ORPHAN/DANGLING Data WayPoints).

> In Claude Code, a PostToolUse Hook can be configured to automatically output a TODO registration/update reminder whenever source files are edited.

---

## 🧩 Dependencies

- [cobra](https://github.com/spf13/cobra) — CLI framework

Minimal external dependencies policy (beyond the standard library).

---

## 🔗 Related Projects

- 🌐 **[openLoadstar](https://github.com/openLoadstar/openLoadstar)** — Full ecosystem overview
- 📖 **[spec](https://github.com/openLoadstar/spec)** — LOADSTAR methodology specification
- 🖥️ **[ui](https://github.com/openLoadstar/ui)** — Spring Boot + React Explorer UI
- 🔌 **[mcp](https://github.com/openLoadstar/mcp)** — Python MCP server (for external AI clients: Claude Desktop, Cursor, etc.)

---

## 📮 Contributing / Security

- 🤝 **Contributing**: [openLoadstar/CONTRIBUTING.md](https://github.com/openLoadstar/openLoadstar/blob/main/CONTRIBUTING.md)
- 🔒 **Security**: [openLoadstar/SECURITY.md](https://github.com/openLoadstar/openLoadstar/blob/main/SECURITY.md) — Please use GitHub Security Advisories.
- 💬 **Questions & Ideas**: [GitHub Discussions](https://github.com/openLoadstar/openLoadstar/discussions)

---

## 📄 License

[Apache License 2.0](./LICENSE)
