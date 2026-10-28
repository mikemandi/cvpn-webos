# Подача CVPN webOS в каталог Homebrew Channel

**Статус:** готовые файлы для заявки на CVPN webOS 2.2.5. Сам факт появления этого документа **не означает**, что приложение добавлено в официальный каталог.

Официальные правила: https://github.com/webosbrew/apps-repo

## Подать через сайт GitHub (рекомендуется)

1. Открой https://github.com/webosbrew/apps-repo/fork и создай **Fork** в аккаунте **mikemandi**.
2. Открой свой Fork `mikemandi/apps-repo`. В нём перейди в папку **packages** и выбери **Add file → Create new file**.
3. Назови файл **packages/su.cvpn.webos.yml**. Вставь содержимое [готовой карточки](su.cvpn.webos.yml) без дополнительных обёрток, Markdown или обратных кавычек. Эта карточка использует реальные английские скриншоты без подписок и отдельные оригинальные иконки CVPN.
4. Сохрани изменение через **Commit changes** (удобно в ветке `add-cvpn-webos`). Не меняй название package-файла или поле `pool: non-free`.
5. Нажми **Contribute → Open pull request** на своём Fork, проверь направление **base: webosbrew/apps-repo:main**, **compare: mikemandi/apps-repo:add-cvpn-webos**.
6. В качестве заголовка укажи **Add CVPN webOS — a free VPN client for rooted LG TVs**, а в описание скопируй [готовый текст заявки](HOMEBREW_PR_BODY.md). Он содержит сведения о тестах, ограничениях и честное раскрытие использования ИИ.
7. Создай PR. Дождись автоматических **Package Lint / Compatibility Check**, посмотри их отчёт и ответь сопровождающим. Только принятие PR сопровождающими Homebrew означает публикацию приложения в каталоге.

Если интерфейс GitHub показывает только ветку `main` в твоём Fork, можно сохранить файл туда и выбрать её как compare-ветку.

## До отправки

- Проверь и лично подтверди текст заявки, версию релиза и тесты на настоящем LG.
- Не обещай работу на телевизорах без root, полный kill switch, IPv6-туннель или зашифрованное HTTP-управление.
- Собственный код CVPN использует [проприетарную лицензию с бесплатным использованием](../LICENSE.ru.md); **non-free — допустимая категория Homebrew, не цена приложения**. Уведомления Xray MPL-2.0 и QR MIT/Apache-2.0 должны сохраняться.
- Публикуй только официальные [GitHub Release v2.2.5](https://github.com/mikemandi/cvpn-webos/releases/tag/v2.2.5) и обезличенные скриншоты.
- Разработку с помощью ИИ нужно раскрыть; заявитель отвечает за проверку кода и безопасность телевизора.

## Ссылки

- [Карточка Homebrew YAML](su.cvpn.webos.yml)
- [Текст PR на английском](HOMEBREW_PR_BODY.md)
- [English README проекта](../README.md)
- [Русский README](../README.ru.md)
- [Лицензия](../LICENSE.ru.md) · [Лицензирование](../LICENSE-PUBLISHING-REVIEW.ru.md)
- [Безопасность](../docs/SECURITY.ru.md)
