# Правила для GitHub Actions

- Существующие workflows унаследованы от Supabase и не являются fork-ready:
  publish/release обращаются к Supabase registries, AWS roles, release branches
  и endpoints. Не запускай и не включай их для форка без отдельного дизайна.
- Любой fork pipeline использует отдельные credentials, least privilege,
  immutable version + source SHA, сборку обеих нужных архитектур, scan и
  consumer-controlled rollout. Секреты передаются только через GitHub secrets/OIDC.
- CI ставит `staticcheck`, `exhaustive` и `gosec` через `@latest`; сначала
  зафиксируй версии и проверь совместимость, не меняй это попутно с provider fix.
- Изменение workflow требует independent review. Не считай зелёный upstream CI
  подтверждением публикации или совместимости SdamEx image.
