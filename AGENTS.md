# Project instructions

This repository keeps its Codex workflow skills in `.agents/skills/`.

Also read `.agents/AGENTS.md` as the project's architecture, technology, and
coding-standard supplement before making implementation changes.

Before starting a task, inspect the relevant `SKILL.md` file(s) and follow their
instructions when the user's request matches one of these triggers:

| Trigger | Skill |
| --- | --- |
| `/archive`, `dọn dẹp bộ nhớ` | `.agents/skills/archive-memory/SKILL.md` |
| `/code-api`, `viết api` | `.agents/skills/code-api/SKILL.md` |
| `/code-db`, `viết schema` | `.agents/skills/code-db/SKILL.md` |
| `/code-ui`, `code giao diện` | `.agents/skills/code-ui/SKILL.md` |
| Docker/containerization work | `.agents/skills/docker-expert/SKILL.md` |
| `/integrate`, `nối api` | `.agents/skills/integrate-api/SKILL.md` |
| `/save`, `lưu ngữ cảnh` | `.agents/skills/save-context/SKILL.md` |
| `/plan-backend`, `quy hoạch hệ thống` | `.agents/skills/system-planner/SKILL.md` |

These are project-local skills. Do not modify or relocate `.agents/skills` unless
the user explicitly asks for migration. When multiple skills match, load all
matching `SKILL.md` files and apply them in the order relevant to the task.

Slash command dispatch details and usage examples are maintained in
`.agents/skills/AGENTS.md`.
