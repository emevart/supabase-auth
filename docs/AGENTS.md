# Правила для fork docs

- `docs/fork/README.md` хранит текущий downstream checkpoint и consumer contract;
  upstream README и CONTRIBUTING сохраняют происхождение и общие инструкции.
- Разделяй verified repository facts, consumer-doc evidence и непроверенное live
  состояние. Не копируй secrets, registry credentials, приватные endpoints или env values.
- После изменения проверяй относительные ссылки, `git diff --check` и отсутствие
  ложных заявлений о CI, DB, OAuth, publish или release.
