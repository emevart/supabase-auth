# Supabase Auth fork: правила для агентов

Обновлено: 2026-09-09.

Этот публичный репозиторий — downstream fork Supabase Auth для SdamEx.
Канонические инструкции Claude Code, Codex и других агентов — этот файл и
вложенные `AGENTS.md`; `CLAUDE.md` импортирует корневой файл.

## Начало задачи

- Прочитай поручение, `git status`, branch/HEAD и
  [checkpoint форка](docs/fork/README.md). Локальная рабочая база форк-патчей —
  `yandex`; не переноси изменение на `master` без отдельного решения о sync.
- Прочитай ближайший scoped `AGENTS.md`, затронутый код, `CONTRIBUTING.md`,
  релевантные targets `Makefile` и CI. Upstream README сохраняет атрибуцию и
  общие инструкции; fork contract имеет приоритет для downstream-операций.
- Не считай план, комментарий Make target или старый CI run доказательством.
  Проверяй фактические зависимости target и текущий diff.

## Границы безопасности auth

- Изменения OAuth, identity linking, verified email, PKCE/state, redirect URI,
  refresh/provider tokens, sessions и migrations имеют высокий риск. Нужны
  узкие тесты, независимое review и проверка совместимости SdamEx.
- Секреты, env values, OAuth codes, tokens, cookies, реальные пользовательские
  данные и дампы БД не попадают в git, команды, issue или отчёты.
- Не запускай live OAuth/SMTP, production/staging DB, registry publish или release
  без явного scope. Разрешение SdamEx на конкретную миграцию или rollout не даёт
  blanket-разрешения этому fork; уже согласованный точный scope переживает handoff.

## Изменения и параллельность

- Существенное изменение provider/auth contract сначала описывается в spec или
  issue: пользовательский результат, trust boundary, совместимость, rollout и rollback.
- Самостоятельные параллельные изменения веди в отдельных checkout/worktree и
  ветках. Subagents используют для ограниченного исследования или review;
  один координатор отвечает за интеграцию. Не давай двум агентам одновременно
  менять provider dispatch, configuration, migrations или workflow.
- Выбирай effort и проверки по риску. Не запускай параллельно тяжёлые Go race
  тесты и Docker/DB окружения без проверки ресурсов и изоляции.

## Находки вне scope

Побочный вопрос, баг или предложение записывается с source/evidence, разделением
`observation` и `hypothesis`, влиянием, зависимостью и owner решения. Сначала
проверь и найди дубликат; затем `accept`, связать с backlog, `defer` с причиной
и trigger либо `reject` с причиной. Не расширяй текущий diff молча и не создавай
issue на каждую мысль. Blocker текущей задачи подними сразу; независимая находка
может ждать отдельного трека. Внешние сообщения и карточки требуют поручения.

## Команды

- Безопасные статические первые шаги: `git diff --check`, `gofmt -l <files>`,
  узкий `go test` только если его окружение понятно и БД не нужна.
- `make all` выполняет `vet`, security/static tools и build, но **не tests**,
  несмотря на help comment.
- `make test` сначала собирает native и linux/arm64 binaries, затем запускает
  весь набор с coverage, `-p 1` и race detector; тестам нужна PostgreSQL.
- `make docker-test` заканчивается `docker compose down -v` и может удалить
  общие volumes. Не запускай его в shared окружении.
- `make build-strip` перезаписывает `internal/utilities/version.go`. Не используй
  как read-only verifier. `format` также изменяет файлы.
- Targets установки tools используют `@latest`; результат зависит от времени.
  Не устанавливай их автоматически и фиксируй версии перед созданием fork gate.

## CI, release и consumer

- Существующие publish/release/dogfooding workflows принадлежат upstream Supabase:
  они нацелены на Supabase registries, AWS roles, branches и endpoints. Они не
  являются готовым pipeline форка и не должны запускаться для fork release.
- Yandex provider и Docker `RELEASE_VERSION` — downstream patches. Контракт
  потребителя описан в `docs/fork/README.md`; любые изменения provider name,
  settings/config/env, identity data, callback или image version координируются
  с SdamEx до merge.
- Не публикуй image, tag, release или артефакт без отдельного проверенного
  fork pipeline и явного разрешения на точные target, version и SHA.

## Завершение и handoff

- Для code change выполни узкие тесты и применимые CI-equivalent проверки;
  явно перечисли тяжёлые/DB/live gates, которые не запускались.
- Обнови fork checkpoint при изменении downstream contract или sync state.
  Передай branch/base/HEAD, diff, проверки, holds, auth risks и следующий шаг.
- Не коммить и не push без соответствующего поручения. Не переносить сюда
  release-процесс или approvals другого репозитория.
