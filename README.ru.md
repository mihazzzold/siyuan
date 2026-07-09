<p align="center">
<img alt="SiYuan" src="https://b3log.org/images/brand/siyuan-128.png">
<br>
<em>Рефакторинг мышления</em>
<br><br>
<a title="Build Status" target="_blank" href="https://github.com/siyuan-note/siyuan/actions/workflows/cd.yml"><img src="https://img.shields.io/github/actions/workflow/status/siyuan-note/siyuan/cd.yml?style=flat-square"></a>
<a title="Releases" target="_blank" href="https://github.com/siyuan-note/siyuan/releases"><img src="https://img.shields.io/github/release/siyuan-note/siyuan.svg?style=flat-square&color=9CF"></a>
<a title="Downloads" target="_blank" href="https://github.com/siyuan-note/siyuan/releases"><img src="https://img.shields.io/github/downloads/siyuan-note/siyuan/total.svg?style=flat-square&color=blueviolet"></a>
<br>
<a title="Docker Pulls" target="_blank" href="https://hub.docker.com/r/b3log/siyuan"><img src="https://img.shields.io/docker/pulls/b3log/siyuan.svg?style=flat-square&color=green"></a>
<a title="Docker Image Size" target="_blank" href="https://hub.docker.com/r/b3log/siyuan"><img src="https://img.shields.io/docker/image-size/b3log/siyuan.svg?style=flat-square&color=ff96b4"></a>
<a title="Hits" target="_blank" href="https://github.com/siyuan-note/siyuan"><img src="https://hits.b3log.org/siyuan-note/siyuan.svg"></a>
<br>
<a title="AGPLv3" target="_blank" href="https://www.gnu.org/licenses/agpl-3.0.txt"><img src="http://img.shields.io/badge/license-AGPLv3-orange.svg?style=flat-square"></a>
<a title="Code Size" target="_blank" href="https://github.com/siyuan-note/siyuan"><img src="https://img.shields.io/github/languages/code-size/siyuan-note/siyuan.svg?style=flat-square&color=yellow"></a>
<a title="GitHub Pull Requests" target="_blank" href="https://github.com/siyuan-note/siyuan/pulls"><img src="https://img.shields.io/github/issues-pr-closed/siyuan-note/siyuan.svg?style=flat-square&color=FF9966"></a>
<br>
<a title="GitHub Commits" target="_blank" href="https://github.com/siyuan-note/siyuan/commits/master"><img src="https://img.shields.io/github/commit-activity/m/siyuan-note/siyuan.svg?style=flat-square"></a>
<a title="Last Commit" target="_blank" href="https://github.com/siyuan-note/siyuan/commits/master"><img src="https://img.shields.io/github/last-commit/siyuan-note/siyuan.svg?style=flat-square&color=FF9900"></a>
<br><br>
<a title="Twitter" target="_blank" href="https://twitter.com/b3logos"><img alt="Twitter Follow" src="https://img.shields.io/twitter/follow/b3logos?label=Follow&style=social"></a>
<a title="Discord" target="_blank" href="https://discord.gg/dmMbCqVX7G"><img alt="Chat on Discord" src="https://img.shields.io/discord/808152298789666826?label=Discord&logo=Discord&style=social"></a>
<br><br>
<a href="https://trendshift.io/repositories/3949" target="_blank"><img src="https://trendshift.io/api/badge/repositories/3949" alt="siyuan-note%2Fsiyuan | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
</p>

<p align="center">
<a href="README.md">English</a>
| <a href="README.zh-CN.md">中文</a>
| <a href="README.ja.md">日本語</a>
| <a href="README.tr.md">Türkçe</a>
| <b>Русский</b>
</p>

---

## Содержание

- [💡 Введение](#-введение)
- [🔮 Возможности](#-возможности)
- [🏗️ Архитектура и экосистема](#️-архитектура-и-экосистема)
- [🌟 История звёзд](#-история-звёзд)
- [🗺️ Дорожная карта](#️-дорожная-карта)
- [🚀 Загрузка и установка](#-загрузка-и-установка)
  - [Магазины приложений](#магазины-приложений)
  - [Установочный пакет](#установочный-пакет)
  - [Пакетный менеджер](#пакетный-менеджер)
  - [Размещение в Docker](#размещение-в-docker)
  - [Размещение в Unraid](#размещение-в-unraid)
  - [Размещение в TrueNAS](#размещение-в-truenas)
  - [Предварительные сборки (Insider Preview)](#предварительные-сборки-insider-preview)
- [⌨️ Интерфейс командной строки](#️-интерфейс-командной-строки)
- [🏘️ Сообщество](#️-сообщество)
- [🛠️ Руководство по разработке](#️-руководство-по-разработке)
- [❓ Часто задаваемые вопросы](#-часто-задаваемые-вопросы)
  - [Как SiYuan хранит данные?](#как-siyuan-хранит-данные)
  - [Поддерживается ли синхронизация данных через сторонние облачные диски?](#поддерживается-ли-синхронизация-данных-через-сторонние-облачные-диски)
  - [SiYuan — открытое ПО?](#siyuan--открытое-по)
  - [Как обновиться до новой версии?](#как-обновиться-до-новой-версии)
  - [Что делать, если у некоторых блоков (например, абзацев внутри элементов списка) нет значка блока?](#что-делать-если-у-некоторых-блоков-например-абзацев-внутри-элементов-списка-нет-значка-блока)
  - [Что делать, если потерян ключ репозитория данных?](#что-делать-если-потерян-ключ-репозитория-данных)
  - [Нужно ли платить за использование?](#нужно-ли-платить-за-использование)
- [🙏 Благодарности](#-благодарности)
  - [Участники](#участники)

---

## 💡 Введение

SiYuan — это система управления личными знаниями, ориентированная на приватность, с поддержкой детальных ссылок на уровне блоков и Markdown WYSIWYG.

Добро пожаловать на [форум SiYuan](https://liuyun.io), где можно узнать больше.

Онлайн-руководство пользователя: [English](https://siyuan-en.b3log.org/)

![feature0.png](https://b3logfile.com/file/2025/11/feature0-GfbhEqf.png)

![feature51.png](https://b3logfile.com/file/2025/11/feature5-1-7DJSfEP.png)

## 🔮 Возможности

Большинство функций бесплатны, в том числе для коммерческого использования.

- Блоки контента
  - Ссылки на уровне блоков и двунаправленные связи
  - Пользовательские атрибуты
  - Встраивание SQL-запросов
  - Протокол `siyuan://`
- Редактор
  - Блочный стиль
  - Markdown WYSIWYG
  - Структура списков (outline)
  - Фокусировка на блоке (zoom-in)
  - Редактирование больших документов на миллион слов
  - Математические формулы, диаграммы, блок-схемы, диаграммы Ганта, временные диаграммы, нотные станы и т. д.
  - Веб-клиппер
  - Связь с аннотациями PDF
- Экспорт
  - Ссылки на блоки и встраивание
  - Стандартный Markdown вместе с ресурсами
  - PDF, Word и HTML
  - Копирование в WeChat MP, Zhihu и Yuque
- База данных
  - Табличное представление
- Флеш-карточки с интервальным повторением
- Написание текстов с ИИ и чат вопрос-ответ через OpenAI API
- Tesseract OCR
- Многовкладочный интерфейс, разделение экрана перетаскиванием
- Шаблоны-сниппеты
- Сниппеты JavaScript/CSS
- Приложения для Android/iOS/HarmonyOS
- Развёртывание в Docker
- [API](https://github.com/siyuan-note/siyuan/blob/master/docs/API.md)
- Маркетплейс сообщества

Некоторые функции доступны только платным подписчикам, подробнее см. [Цены](https://b3log.org/siyuan/en/pricing.html).

## 🏗️ Архитектура и экосистема

![SiYuan Arch](https://b3logfile.com/file/2023/05/SiYuan_Arch-Sgu8vXT.png "SiYuan Arch")

| Проект                                                   | Описание                  | Форки                                                                           | Звёзды                                                                               | 
|----------------------------------------------------------|---------------------------|---------------------------------------------------------------------------------|--------------------------------------------------------------------------------------|
| [lute](https://github.com/88250/lute)                    | Движок редактора          | ![GitHub forks](https://img.shields.io/github/forks/88250/lute)                 | ![GitHub Repo stars](https://img.shields.io/github/stars/88250/lute)                 |
| [chrome](https://github.com/siyuan-note/siyuan-chrome)   | Расширение Chrome/Edge    | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/siyuan-chrome)  | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/siyuan-chrome)  |
| [bazaar](https://github.com/siyuan-note/bazaar)          | Маркетплейс сообщества    | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/bazaar)         | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/bazaar)         |
| [dejavu](https://github.com/siyuan-note/dejavu)          | Репозиторий данных        | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/dejavu)         | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/dejavu)         |
| [petal](https://github.com/siyuan-note/petal)            | API плагинов              | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/petal)          | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/petal)          |
| [android](https://github.com/siyuan-note/siyuan-android) | Приложение Android        | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/siyuan-android) | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/siyuan-android) |
| [ios](https://github.com/siyuan-note/siyuan-ios)         | Приложение iOS            | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/siyuan-ios)     | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/siyuan-ios)     |
| [harmony](https://github.com/siyuan-note/siyuan-harmony) | Приложение HarmonyOS      | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/siyuan-harmony) | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/siyuan-harmony) |
| [riff](https://github.com/siyuan-note/riff)              | Интервальное повторение   | ![GitHub forks](https://img.shields.io/github/forks/siyuan-note/riff)           | ![GitHub Repo stars](https://img.shields.io/github/stars/siyuan-note/riff)           |

## 🌟 История звёзд

<a href="https://star-history.com/#siyuan-note/siyuan&Date">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=siyuan-note/siyuan&type=Date&theme=dark" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=siyuan-note/siyuan&type=Date" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=siyuan-note/siyuan&type=Date" />
 </picture>
</a>

## 🗺️ Дорожная карта

- [План развития и прогресс SiYuan](https://github.com/orgs/siyuan-note/projects/1)
- [Журнал изменений SiYuan](CHANGELOG.md)

## 🚀 Загрузка и установка

На настольных и мобильных устройствах рекомендуется в первую очередь устанавливать приложение через магазин приложений — так в будущем можно будет обновлять версию в один клик.

### Магазины приложений

Мобильные:

- [App Store](https://apps.apple.com/us/app/siyuan/id1583226508)
- [Google Play](https://play.google.com/store/apps/details?id=org.b3log.siyuan)
- [F-Droid](https://f-droid.org/packages/org.b3log.siyuan)

Настольные:

- [Microsoft Store](https://apps.microsoft.com/detail/9p7hpmxp73k4)

### Установочный пакет

- [B3log](https://b3log.org/siyuan/en/download.html)
- [GitHub](https://github.com/siyuan-note/siyuan/releases)

### Пакетный менеджер

#### `siyuan`

[![Packaging status](https://repology.org/badge/vertical-allrepos/siyuan.svg)](https://repology.org/project/siyuan/versions)

#### `siyuan-note`

[![Packaging status](https://repology.org/badge/vertical-allrepos/siyuan-note.svg)](https://repology.org/project/siyuan-note/versions)

### Размещение в Docker

<details>
<summary>Развёртывание в Docker</summary>

#### Обзор

Самый простой способ развернуть SiYuan на сервере — через Docker.

- Имя образа: `b3log/siyuan`
- [Страница образа](https://hub.docker.com/r/b3log/siyuan)

#### Структура файлов

Программа целиком расположена в `/opt/siyuan/` — по сути это структура папки resources установочного пакета Electron:

- appearance: значки, темы, языки
- guide: документы руководства пользователя
- stage: интерфейс и статические ресурсы
- kernel: программа ядра

#### Точка входа

Точка входа задаётся при сборке Docker-образа: `ENTRYPOINT ["/opt/siyuan/entrypoint.sh"]`. Этот скрипт позволяет изменить `PUID` и `PGID` пользователя, от имени которого будет работать процесс внутри контейнера. Это особенно важно для решения проблем с правами доступа при монтировании каталогов с хоста. `PUID` (ID пользователя) и `PGID` (ID группы) можно передать как переменные окружения, что упрощает обеспечение корректных прав при доступе к смонтированным каталогам хоста.

Используйте следующие параметры при запуске контейнера командой `docker run b3log/siyuan`:

> **Примечание:** Начиная с v3.7.0 подкоманду `serve` нужно передавать явно (например, `docker run b3log/siyuan serve --workspace=...`). Выполните `docker run --rm b3log/siyuan serve --help`, чтобы увидеть все параметры запуска.

- `--workspace`: путь к папке рабочего пространства, монтируется в контейнер через `-v` на хосте
- `--accessAuthCode`: пароль экрана блокировки

Дополнительные параметры можно узнать с помощью `--help`. Вот пример команды запуска с новыми переменными окружения:

```bash
docker run -d \
  -v workspace_dir_host:workspace_dir_container \
  -p 6806:6806 \
  -e PUID=1001 -e PGID=1002 \
  b3log/siyuan \
  serve \
  --workspace=workspace_dir_container \
  --accessAuthCode=xxx
```

- `PUID`: пользовательский ID (необязательно, по умолчанию `1000`)
- `PGID`: ID группы (необязательно, по умолчанию `1000`)
- `workspace_dir_host`: путь к папке рабочего пространства на хосте
- `workspace_dir_container`: путь к папке рабочего пространства в контейнере, указанный в `--workspace`
  - Альтернативно путь можно задать через переменную окружения `SIYUAN_WORKSPACE_PATH`. Если заданы обе, приоритет всегда у командной строки
- `accessAuthCode`: пароль экрана блокировки (обязательно **измените его**, иначе любой сможет получить доступ к вашим данным)
  - Альтернативно пароль экрана блокировки можно задать через переменную окружения `SIYUAN_ACCESS_AUTH_CODE`. Если заданы обе, приоритет всегда у командной строки
  - Чтобы отключить пароль экрана блокировки, установите переменную окружения `SIYUAN_ACCESS_AUTH_CODE_BYPASS=true`
- `SIYUAN_LANG`: язык интерфейса (необязательно, в Docker по умолчанию `en`, если не задан). Принимает теги BCP 47, такие как `zh-CN`/`zh-TW`/`en`/`ja`/`pt-BR`; устаревшие значения с подчёркиванием вроде `zh_CN`/`en_US` также принимаются для обратной совместимости. Не задавайте эту переменную, если хотите, чтобы язык, выбранный в **Настройках**, сохранялся между перезапусками; если она задана, то применяется при каждом запуске и переопределяет сохранённую настройку
  - Альтернативно используйте параметр командной строки `--lang`. Если заданы оба, приоритет у командной строки

Для простоты рекомендуется настроить одинаковый путь к папке рабочего пространства на хосте и в контейнере, например задать и `workspace_dir_host`, и `workspace_dir_container` как `/siyuan/workspace`. Соответствующая команда запуска:

```bash
docker run -d \
  -v /siyuan/workspace:/siyuan/workspace \
  -p 6806:6806 \
  -e PUID=1001 -e PGID=1002 \
  b3log/siyuan \
  serve \
  --workspace=/siyuan/workspace/ \
  --accessAuthCode=xxx
```

#### Docker Compose

Пользователи, запускающие SiYuan через Docker Compose, могут передавать переменные окружения `PUID` и `PGID` для настройки ID пользователя и группы. Пример конфигурации Docker Compose:

```yaml
version: "3.9"
services:
  main:
    image: b3log/siyuan
    command: ['serve', '--workspace=/siyuan/workspace/', '--accessAuthCode=${AuthCode}']
    ports:
      - 6806:6806
    volumes:
      - /siyuan/workspace:/siyuan/workspace
    restart: unless-stopped
    environment:
      - TZ=${YOUR_TIME_ZONE}    # Список идентификаторов часовых поясов: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
      - PUID=${YOUR_USER_PUID}  # Пользовательский ID
      - PGID=${YOUR_USER_PGID}  # ID группы
```

В этой конфигурации:

- `PUID` и `PGID` задаются динамически и передаются в контейнер
- Если эти переменные не указаны, используется значение по умолчанию `1000`

Указывая `PUID` и `PGID` в переменных окружения, вы избавляетесь от необходимости явно задавать директиву `user` (`user: '1000:1000'`) в compose-файле. Контейнер динамически настроит пользователя и группу на основе этих переменных при запуске.

#### Права пользователя

В образе скрипт `entrypoint.sh` обеспечивает создание пользователя и группы `siyuan` с указанными `PUID` и `PGID`. Поэтому при создании папки рабочего пространства на хосте обратите внимание на то, чтобы владелец и группа папки соответствовали `PUID` и `PGID`, которые вы планируете использовать. Например:

```bash
chown -R 1001:1002 /siyuan/workspace
```

Если вы используете собственные значения `PUID` и `PGID`, скрипт точки входа обеспечит создание правильного пользователя и группы внутри контейнера, а владелец смонтированных томов будет соответствующим образом скорректирован. Нет необходимости вручную передавать `-u` в `docker run` или `docker-compose` — переменные окружения выполнят настройку.

#### Скрытие порта

Используйте обратный прокси NGINX, чтобы скрыть порт 6806. Обратите внимание:

- Настройте обратное проксирование WebSocket для `/ws`

#### Примечания

- Обязательно проверьте корректность смонтированного тома, иначе данные будут потеряны после удаления контейнера
- Не используйте переписывание URL (URL rewriting) для перенаправления, иначе могут возникнуть проблемы с аутентификацией; рекомендуется настроить обратный прокси
- При проблемах с правами доступа убедитесь, что переменные окружения `PUID` и `PGID` соответствуют владельцу смонтированных каталогов на хосте

#### Ограничения

- Не поддерживается подключение настольных и мобильных приложений, работа возможна только через браузер
- Не поддерживается экспорт в форматы PDF, HTML и Word
- Не поддерживается импорт Markdown-файлов

</details>

### Размещение в Unraid

<details>
<summary>Развёртывание в Unraid</summary>

Примечание: сначала выполните в терминале `chown -R 1000:1000 /mnt/user/appdata/siyuan`

Шаблон для справки:

```
Web UI: 6806
Container Port: 6806
Container Path: /home/siyuan
Host path: /mnt/user/appdata/siyuan
PUID: 1000
PGID: 1000
Publish parameters: serve --accessAuthCode=******(пароль экрана блокировки)
```

</details>

### Размещение в TrueNAS

<details>
<summary>Развёртывание в TrueNAS</summary>

Примечание: сначала выполните приведённые ниже команды в TrueNAS Shell. Замените `Pool_1/Apps_Data/siyuan` на путь к вашему датасету.

```shell
zfs create Pool_1/Apps_Data/siyuan
chown -R 1001:1002 /mnt/Pool_1/Apps_Data/siyuan
chmod 755 /mnt/Pool_1/Apps_Data/siyuan
```

Перейдите в Apps - DiscoverApps - More Options (сверху справа, рядом с Custom App) - Install via YAML

Шаблон для справки:

```yaml
services:
  siyuan:
    image: b3log/siyuan
    container_name: siyuan
    command: ['serve', '--workspace=/siyuan/workspace/', '--accessAuthCode=2222']
    ports:
      - 6806:6806
    volumes:
      - /mnt/Pool_1/Apps_Data/siyuan:/siyuan/workspace  # Замените на путь к вашему датасету
    restart: unless-stopped
    environment:
      - TZ=America/New_York  # При необходимости замените на свой часовой пояс
      - PUID=1001
      - PGID=1002
```

</details>

### Предварительные сборки (Insider Preview)

Перед крупными обновлениями мы выпускаем предварительные сборки, см. [https://github.com/siyuan-note/insider](https://github.com/siyuan-note/insider).

## ⌨️ Интерфейс командной строки

Встроенный CLI предоставляет прямой доступ к данным рабочего пространства — запущенный сервер не требуется.

### Быстрый старт

```bash
# Список всех блокнотов
siyuan notebook list -w ~/SiYuan

# Полнотекстовый поиск с выводом в JSON
siyuan search "keyword" -w ~/SiYuan -f json

# Экспорт документа в Markdown
siyuan export md --id <block-id> -w ~/SiYuan
```

### Доступные команды

| Категория | Команды |
|----------|----------|
| Блокноты и документы | `notebook`, `document`, `dailynote` — CRUD и ежедневные заметки |
| Контент | `block`, `attr`, `outline` — чтение/запись блоков, атрибуты, структура |
| Метаданные | `tag`, `bookmark`, `template` — теги, закладки, шаблоны-сниппеты |
| Запросы | `search`, `sql` — полнотекстовый поиск и SQL-запросы |
| Ссылки | `ref` — обратные ссылки и упоминания |
| Импорт/экспорт | `export`, `import`, `inbox` — Markdown, HTML, предпросмотр, Word, .sy.zip, Data, облачный инбокс |
| Управление данными | `repo`, `history`, `sync` — снимки, версии, облачная синхронизация |
| Утилиты | `asset`, `file` — ресурсы и файловая система |
| База данных | `database` — управление attribute view |
| Сервер | `serve` — запуск HTTP-сервера ядра |
| Рабочее пространство и система | `workspace`, `system` — список, просмотр, информация о системе |

Выполните `siyuan --help`, чтобы увидеть полное дерево команд. Используйте `-f json` (по умолчанию `-f table`) для вывода, удобного для скриптов. Большинство изменяющих команд также поддерживают `--dry-run` — предпросмотр изменений без их применения.

### Настройка

Бинарный файл CLI — `SiYuan-Kernel` в `<install>/resources/kernel`.
Установщик для Windows добавляет его в PATH автоматически.
На macOS/Linux создайте символическую ссылку вручную:

```bash
# macOS
ln -s /Applications/SiYuan.app/Contents/Resources/kernel/SiYuan-Kernel /usr/local/bin/siyuan

# Linux
ln -s /path/to/SiYuan/resources/kernel/SiYuan-Kernel /usr/local/bin/siyuan
```

## 🏘️ Сообщество

- [Форум обсуждений (англ.)](https://liuyun.io)
- [Сводка сообществ пользователей](https://liuyun.io/article/1687779743723)
- [Awesome SiYuan](https://github.com/siyuan-note/awesome)

## 🛠️ Руководство по разработке

См. [Руководство по разработке](https://github.com/siyuan-note/siyuan/blob/master/.github/CONTRIBUTING.md).

## ❓ Часто задаваемые вопросы

### Как SiYuan хранит данные?

Данные сохраняются в папке data рабочего пространства:

- `assets` — все вставленные ресурсы
- `emojis` — изображения эмодзи
- `snippets` — сниппеты кода
- `storage` — условия запросов, раскладки, флеш-карточки и т. д.
- `templates` — шаблоны-сниппеты
- `widgets` — виджеты
- `plugins` — плагины
- `public` — публичные данные
- Остальные папки — это блокноты, созданные пользователем; файлы с расширением `.sy` в папке блокнота хранят данные документов, формат данных — JSON

### Поддерживается ли синхронизация данных через сторонние облачные диски?

Синхронизация данных через сторонние облачные диски не поддерживается, иначе данные могут быть повреждены.

Хотя сторонние диски синхронизации не поддерживаются, поддерживается подключение к стороннему облачному хранилищу (привилегия подписчиков).

Кроме того, можно рассмотреть ручной экспорт и импорт данных для синхронизации:

- Настольная версия: <kbd>Настройки</kbd> - <kbd>Экспорт</kbd> - <kbd>Экспорт данных</kbd> / <kbd>Импорт данных</kbd>
- Мобильная версия: <kbd>Правая панель</kbd> - <kbd>О программе</kbd> - <kbd>Экспорт данных</kbd> / <kbd>Импорт данных</kbd>

### SiYuan — открытое ПО?

SiYuan полностью открыт, вклад приветствуется:

- [Пользовательский интерфейс и ядро](https://github.com/siyuan-note/siyuan)
- [Android](https://github.com/siyuan-note/siyuan-android)
- [iOS](https://github.com/siyuan-note/siyuan-ios)
- [HarmonyOS](https://github.com/siyuan-note/siyuan-harmony)
- [Расширение-клиппер для Chrome](https://github.com/siyuan-note/siyuan-chrome)

Подробнее см. [Руководство по разработке](https://github.com/siyuan-note/siyuan/blob/master/.github/CONTRIBUTING.md).

### Как обновиться до новой версии?

- Если приложение установлено через магазин приложений, обновляйте его через магазин
- Если оно установлено через установочный пакет на настольном компьютере, можно включить опцию <kbd>Настройки</kbd> - <kbd>О программе</kbd> - <kbd>Автоматически загружать установочный пакет обновления</kbd> — тогда SiYuan будет автоматически скачивать последнюю версию установочного пакета и предлагать установку
- Если оно установлено вручную из установочного пакета, скачайте установочный пакет заново и установите

Проверить обновления можно в <kbd>Настройки</kbd> - <kbd>О программе</kbd> - <kbd>Текущая версия</kbd> - <kbd>Проверить обновление</kbd>, либо следите за [официальной страницей загрузки](https://b3log.org/siyuan/en/download.html) или [GitHub Releases](https://github.com/siyuan-note/siyuan/releases), чтобы получать новые версии.

### Что делать, если у некоторых блоков (например, абзацев внутри элементов списка) нет значка блока?

Значок блока не отображается у первого дочернего блока внутри элемента списка. Можно поместить курсор в этот блок и вызвать его меню сочетанием <kbd>Ctrl+/</kbd>.

### Что делать, если потерян ключ репозитория данных?

- Если ключ репозитория данных ранее был корректно инициализирован на нескольких устройствах, он одинаков на всех устройствах, и его можно получить в <kbd>Настройки</kbd> - <kbd>О программе</kbd> - <kbd>Ключ репозитория данных</kbd> - <kbd>Скопировать строку ключа</kbd>
- Если он не был настроен корректно (например, ключи на разных устройствах различаются) или все устройства недоступны и строку ключа получить невозможно, ключ можно сбросить следующим образом:

  1. Вручную сделайте резервную копию данных: используйте <kbd>Экспорт данных</kbd> или просто скопируйте папку <kbd>workspace/data/</kbd> в файловой системе
  2. <kbd>Настройки</kbd> - <kbd>О программе</kbd> - <kbd>Ключ репозитория данных</kbd> - <kbd>Сбросить репозиторий данных</kbd>
  3. Заново инициализируйте ключ репозитория данных. После инициализации ключа на одном устройстве импортируйте его на остальных
  4. Облако использует новый каталог синхронизации, старый каталог синхронизации больше недоступен и может быть удалён
  5. Существующие облачные снимки больше недоступны и могут быть удалены

### Нужно ли платить за использование?

Большинство функций бесплатны, в том числе для коммерческого использования.

Привилегии подписчика доступны только после оплаты, подробнее см. [Цены](https://b3log.org/siyuan/en/pricing.html).

## 🙏 Благодарности

Появление SiYuan было бы невозможно без множества проектов с открытым исходным кодом и их участников. См. исходный код проекта: kernel/go.mod, app/package.json и домашние страницы проектов.

Рост SiYuan неотделим от отзывов и поддержки пользователей. Спасибо всем за помощь SiYuan ❤️

### Участники

Присоединяйтесь к нам и вносите свой вклад в код SiYuan вместе с нами.

<a href="https://github.com/siyuan-note/siyuan/graphs/contributors">
   <img src="https://contrib.rocks/image?repo=siyuan-note/siyuan" />
</a>
