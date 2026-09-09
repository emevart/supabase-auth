# SdamEx fork: checkpoint и контракт потребителя

Обновлено: 2026-09-09.

Этот репозиторий сохраняет upstream [Supabase Auth](../../README.md) и содержит
downstream patches, нужные SdamEx. Документ описывает проверенное состояние
исходников; он не подтверждает текущее состояние live auth или registry.

## Checkpoint

| Поле | Проверенное значение |
| --- | --- |
| Fork remote | `https://github.com/emevart/supabase-auth.git` |
| База fork-разработки | `yandex` |
| Проверенный HEAD | `0a39c986aebf1b4f931d78842356a4d67e015a58` |
| Merge base с `origin/master` | `6a8fc26f6042ba92cda51bc62b70f1b61c2cd12e` |
| Downstream commits после merge base | 2 |
| Upstream-only commits в локальных refs | 85 |

GitHub default branch — master (проверено через gh repo view). Для новой задачи
над downstream-кодом явно выбирай базу yandex: локальный origin/HEAD не меняет
настройку default branch на GitHub. Изменение default branch — отдельное решение.

В локальных refs `origin/HEAD` указывает на `origin/yandex`, хотя upstream-style
workflows и CONTRIBUTING ориентированы на `master`. Интегрировать эту docs-правку
следует в `yandex`. Sync с `origin/master` — отдельная высокорисковая работа:
нельзя rebase/merge вслепую из-за большого расхождения и auth/migration изменений.

Downstream diff против merge base состоит из пяти файлов:

- новый `internal/api/provider/yandex.go`;
- wiring `yandex` в external dispatch, configuration и settings;
- Docker build argument `RELEASE_VERSION`, передаваемый в `make build`.

## Контракт SdamEx

Read-only проверка consumer repo подтверждает:

- SdamEx собирает и разворачивает отдельный auth image из ветки `yandex`;
- версия образа содержит downstream suffix вида `v2.184.0-yandex.1` и должна
  собираться с соответствующим Docker `--build-arg RELEASE_VERSION=...`;
- конфигурация использует `GOTRUE_EXTERNAL_YANDEX_ENABLED`, `CLIENT_ID`, `SECRET`
  и `REDIRECT_URI`; значения остаются вне этого репозитория;
- callback должен точно совпадать с окружением, а доступность интерфейса SdamEx
  отдельно управляется consumer feature flag `yandex_oauth_enabled`.

Канонический rollout/image path остаётся в consumer runbook
`docs/deployment/YANDEX-OAUTH-SETUP.md` репозитория SdamEx. Здесь не дублируются
registry path, credentials и live status. Изменения provider identifier, env/config,
settings response, subject/identity metadata, email trust, callback, migrations,
version output или image layout требуют consumer review до интеграции.

## Auth trust boundaries

- Yandex — OAuth 2.0 provider, текущий код возвращает `RequiresPKCE() == false`.
  Изменение требует проверки общего state/PKCE flow и реального provider contract.
- Код доверяет Yandex profile и помечает возвращённые email verified. Это влияет
  на account linking; предположение должно быть подтверждено актуальной provider
  документацией и тестами перед изменением или upstream sync.
- `u.ID` становится identity subject. Нельзя заменять его login/email или менять
  namespace provider без плана совместимости существующих identities.
- Redirect/referrer allow-list и callback различаются по окружениям. Неверный URI
  может сломать вход или отправить auth material не тому получателю.
- Provider access/refresh tokens могут участвовать в callback response. Их нельзя
  логировать, сохранять в evidence или передавать через issue.
- Изменения sessions/refresh-token rotation способны массово разлогинить пользователей;
  они требуют отдельного rollout, rollback и consumer acceptance.

## Команды и их реальный эффект

| Команда | Реальный эффект | Статус для агента |
| --- | --- | --- |
| `gofmt -l <files>` | Только перечисляет неформатированные Go files | безопасная static-проверка |
| `git diff --check` | Whitespace errors текущего diff | обязательна для docs/change |
| `make all` | `vet`, `gosec`, `staticcheck`, `exhaustive`, build | не запускает tests; может установить `@latest` tools |
| `make build` | Скачивает/verifies modules, строит native и linux/arm64 binaries | пишет build artifacts, сеть и ресурсы |
| `make test` | Сначала `make build`, затем coverage + весь suite с `-race` | тяжёлый, требует DB/env |
| `make docker-test` | Поднимает DB/auth, мигрирует, тестирует, затем `down -v` | опасен для shared volumes |
| `make build-strip` | Переписывает `internal/utilities/version.go`, строит artifact | не использовать как verifier |
| `make format` | `gofmt -s -w .` | изменяет весь Go tree |

`check-gosec`, `check-staticcheck`, `check-exhaustive` и `check-oapi-codegen`
устанавливают инструменты с `@latest`; это источник недетерминированности. До
создания fork CI версии надо зафиксировать отдельной scoped правкой.

## CI и release

Test workflow создаёт PostgreSQL, применяет migrations и запускает race suite.
Publish/release/dogfooding workflows ориентированы на инфраструктуру Supabase:
Supabase container registries, AWS roles/buckets, upstream release refs и live
Supabase endpoints. В fork они **не являются готовым release pipeline**.

Не запускать publish/release и не заявлять fork release-ready. Нужен отдельный
дизайн: target registry, immutable tag и SHA, credentials/OIDC, multi-arch build,
scan, staging consumer update, OAuth acceptance, rollback и production gate.

## Находки и следующий шаг

Побочная находка проходит цикл: source/evidence → observation/hypothesis →
question/bug/proposal → impact/dependency/owner → verify → `accepted`, `backlog`,
`deferred` с причиной/trigger или `rejected` с причиной. Сначала ищется дубликат;
issue не создаётся на каждую мысль, текущий scope молча не расширяется.

Следующий разумный этап — independent review этих docs, затем PR в `yandex`.
Отдельными карточками, после triage, могут идти upstream sync и fork-specific CI;
эта docs-правка их не реализует.
