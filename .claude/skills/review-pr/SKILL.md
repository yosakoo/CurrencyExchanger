---
# Скилл запускается вручную: /review-pr 42
# $ARGUMENTS — номер PR в GitHub
name: review-pr
description: Провести code review PR по стандартам команды. Проверяет стиль, тесты, безопасность и оставляет комментарии в GitHub.
disable-model-invocation: true
allowed-tools:
  - Read
  - Bash
  - Grep
  - mcp__github__get_pull_request
  - mcp__github__get_pull_request_files
  - mcp__github__create_review
  - mcp__github__create_review_comment
---

Проведи code review для PR #$ARGUMENTS:

1. Получи метаданные PR через GitHub MCP — title, description, base branch
2. Получи список изменённых файлов и diff
3. Для каждого изменённого файла проверь:
   - Соответствие соглашениям из CLAUDE.md
   - Корректность обработки ошибок
   - Наличие тестов для новой функциональности
   - Потенциальные проблемы с безопасностью (SQL injection, утечки данных, незащищённые endpoints)
   - Читаемость и именование
4. Оставь inline-комментарии к конкретным строкам через GitHub MCP
5. Итоговый вердикт:
   - APPROVE — если всё в порядке
   - REQUEST_CHANGES — с перечнем обязательных правок
   - COMMENT — если есть вопросы но блокеров нет
