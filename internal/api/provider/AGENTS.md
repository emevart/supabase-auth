# Правила для OAuth providers

- Provider boundary считается недоверенной: проверяй token exchange, endpoint,
  profile decoding, stable subject, email presence/verification и error mapping.
- Не меняй provider identifier, default scopes, callback contract или metadata
  без тестов dispatch/config/settings и проверки consumer contract.
- Для Yandex отдельно проверяй осознанное `RequiresPKCE`, OAuth state общего
  потока, точный redirect URI, доверие к `default_email`/`emails`, стабильность
  `ID` как subject и обращение с provider access/refresh tokens.
- Live OAuth требует отдельного scope и тестового приложения. Unit-тесты используют
  локальный mock HTTP server и не должны содержать реальные credentials или PII.
