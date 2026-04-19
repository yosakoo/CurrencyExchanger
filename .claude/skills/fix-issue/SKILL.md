---
# Скилл запускается вручную: /fix-issue 247
# $ARGUMENTS — номер задачи в GitHub Issues
name: fix-issue
description: Исправить задачу из GitHub Issues по номеру. Читает issue, находит код, вносит правку, пишет тест, открывает PR.
disable-model-invocation: true
allowed-tools:
  - Read
  - Edit
  - Write
  - Bash
  - Grep
  - Glob
  - mcp__github__get_issue
  - mcp__github__create_pull_request
  - mcp__github__create_issue_comment
---

Исправь задачу GitHub #$ARGUMENTS:

1. Прочитай задачу через GitHub MCP — получи title, body, labels
2. Изучи кодовую базу: найди файлы связанные с проблемой через Grep и Glob
3. Прочитай найденные файлы и смежный код для понимания контекста
4. Предложи план исправления в одном абзаце и жди подтверждения
5. После подтверждения — реализуй исправление согласно правилам из CLAUDE.md
6. Напиши unit-тест для исправленного кода
7. Запусти тесты: `make test` — убедись что всё зелёное
8. Запусти линтер: `make lint`
9. Создай коммит: `fix: <краткое описание> (#$ARGUMENTS)`
10. Открой PR через GitHub MCP с описанием что было сделано и почему
