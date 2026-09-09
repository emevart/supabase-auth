# Правила для auth migrations

- Считай применённые migration files неизменяемой историей. Новое изменение БД
  получает новый упорядоченный файл и совместимо с текущим и предыдущим auth image.
- Перед merge оцени lock/backfill/runtime, повторный запуск, rollback или forward
  fix, а также влияние на users, identities, sessions и refresh tokens.
- Не запускай migration против staging/production или общей локальной БД без
  явного target и разрешения. Согласование миграций в consumer repo не является
  blanket-разрешением для произвольной миграции fork.
- Для теста используй изолированную одноразовую БД. `make docker-test` удаляет
  compose volumes через `down -v`, поэтому не запускай его в shared project.
