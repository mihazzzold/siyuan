# JSON-структура файла `.sy` SiYuan — руководство по чтению/записи для ИИ

> Базовая версия спецификации: `2` (актуальна для всех файлов).
> Проверено на образцах: `20200825162036-4dx365o.sy` (элементы форматирования), `20200905090211-2vixtlf.sy` (типы блоков).
> Все выводы основаны на реальных образцах и исходном коде Lute / ядра SiYuan. Поля, помеченные `【inferred】`, не наблюдались в образцах напрямую — перед их генерацией сверьтесь с реальным образцом.
> Сопутствующий документ: [`WORKSPACE.ru.md`](./WORKSPACE.ru.md) описывает общую структуру рабочего пространства на диске (как организованы блокноты, родительские/дочерние документы и ресурсы); этот документ сосредоточен на **внутренней** JSON-структуре файла `.sy`.

## 0. Одним предложением

Файл `.sy` — это AST-дерево Lute, сериализованное в JSON. Корневой узел всегда `NodeDocument`; тело — это массив `Children` с рекурсивной вложенностью. Внешней схемы нет — всё состояние живёт в дереве.

---

## 0.5 Когда читать/писать `.sy` напрямую (порядок приоритетов)

SiYuan предлагает три официальных пути изменения данных: **HTTP API, MCP и CLI**. **По умолчанию предпочитайте их.** Ядро берёт на себя сериализацию AST, выделение ID блоков и синхронизацию двух индексов: индекса дерева блоков (`blocktree.db` — отображение «ID блока → путь к файлу», от которого зависят ссылки на блоки и хлебные крошки) и индекса полнотекстового поиска (`siyuan.db` + FTS5). Прямая запись файлов обходит всё это и легко приводит к рассинхронизации индексов.

**Читайте/пишите `.sy` как JSON только тогда, когда официальные пути неудобны.** Подходящие сценарии:
- Массовая офлайн-миграция (холодная инициализация рабочего пространства, импорт внешних данных; структуру рабочего пространства на диске см. в [`WORKSPACE.ru.md`](./WORKSPACE.ru.md))
- Статистика только для чтения, анализ, собственный экспорт / конвертация форматов
- Исправление низкоуровневых структурных проблем (устаревшие файлы, недопустимые узлы)
- Программная генерация каркасов / шаблонов

Разделение труда между четырьмя путями:

| Путь | Роль | Возможности изменения |
|---|---|---|
| **HTTP API** | Онлайн, во время работы | Богатейшие — полный CRUD документов/блоков (`filetree/*`, `block/*`, `transactions`) |
| **MCP** | Набор инструментов для LLM | Подмножество для ИИ-агентов, работающих с документами онлайн |
| **CLI** | Пакетная работа / эксплуатация | Импорт, экспорт, синхронизация, SQL и другие задачи командной строки |
| **Прямое чтение/запись `.sy`** | Предмет этого руководства | Офлайн, массовая, низкоуровневая структурная работа |

> ⚠️ После прямой записи файлов обычно требуется «перестроение индекса», прежде чем поиск/ссылки на блоки заработают. Если SiYuan запущен, предпочитайте HTTP API — пусть ядро само занимается сериализацией и синхронизацией индексов.

---

## 1. Структура верхнего уровня

```json
{
  "ID": "20200825162036-4dx365o",
  "Spec": "2",
  "Type": "NodeDocument",
  "Properties": {
    "icon": "1f4f0",
    "id": "20200825162036-4dx365o",
    "title": "排版元素",
    "type": "doc",
    "updated": "20260616224229"
  },
  "Children": [ ... ]
}
```

| Ключ верхнего уровня | Обязателен | Значение |
|---|---|---|
| `ID` | ✅ | ID блока документа. **Равен имени файла без `.sy`** (не генерируется случайно) |
| `Spec` | ✅ | Всегда `"2"` |
| `Type` | ✅ | Всегда `"NodeDocument"` |
| `Properties` | ✅ | IAL уровня документа — см. §8 |
| `Children` | ✅ | Массив дочерних блоков тела; должен содержать хотя бы один блок |

> ⚠️ Путь к файлу строго соответствует корневому ID: `data/<box>/<...>/<rootID>.sy`. Изменение корневого ID означает переименование файла — не меняйте его без необходимости. Полную структуру файловой системы см. в [`WORKSPACE.ru.md`](./WORKSPACE.ru.md).

---

## 2. Семантика общих полей (применимо к каждому узлу)

| Поле | Тип | Наличие | Значение |
|---|---|---|---|
| `Type` | string | **обязательно у каждого узла** | Дискриминатор типа, например `"NodeParagraph"` |
| `ID` | string | только у блочных узлов | 22-символьный ID блока; строчные/маркерные узлы **его не несут** |
| `Data` | string | у некоторых | Текст / HTML / сырой markdown; **может отсутствовать** (не считайте, что оно всегда есть) |
| `Properties` | object | у большинства блоков | IAL, `map[string]string` |
| `Children` | array | контейнерные / вложимые узлы | Массив дочерних узлов |
| Поля, специфичные для типа | - | по типу | например `HeadingLevel`, `ListData`, `TextMarkType`, `AttributeViewID` |

**Ключевое правило-дискриминатор:** наличие у узла `ID` определяет, является ли он блоком. Есть `ID` ⇒ блок (его `Properties.id` должно равняться его `ID`); нет `ID` ⇒ строчный/маркерный узел.

---

## 3. Правила ID и меток времени

- **Формат ID:** `YYYYMMDDHHMMSS-xxxxxxx` = 14-значная метка времени + `-` + 7 случайных символов `[a-z0-9]`. Пример: `20210104091228-ttcj9nm`.
- **Корневой ID** берётся из имени файла и **не** генерируется заново.
- **ID дочерних блоков** генерируются по схеме выше.
- `Properties.updated` — та же 14-значная метка времени; семантика: «время последнего обновления».
- При изменении `ID` любого блока нужно **синхронно** обновить `Properties.id`. `Properties.updated` также следует обновить до текущего времени.

---

## 4. Каталог типов узлов

### Блочные узлы (имеют ID)

**Листовые блоки:** `NodeParagraph`, `NodeHeading`, `NodeThematicBreak`, `NodeHTMLBlock`, `NodeCodeBlock`, `NodeMathBlock`, `NodeTable`, `NodeBlockQueryEmbed`, `NodeAttributeView`, `NodeIFrame`, `NodeVideo`, `NodeAudio`, `NodeWidget`, `NodeCustomBlock`, `NodeGitConflict`

**Контейнерные блоки:** `NodeList`, `NodeListItem`, `NodeBlockquote`, `NodeCallout`, `NodeSuperBlock`

### Строчные / маркерные узлы (без ID)

`NodeText`, `NodeTextMark`, `NodeImage`, `NodeKramdownSpanIAL`, `NodeHeadingC8hMarker`, `NodeBlockquoteMarker`, `NodeTaskListItemMarker`, `NodeBang`, `NodeOpenBracket`, `NodeCloseBracket`, `NodeOpenParen`, `NodeCloseParen`, `NodeLinkText`, `NodeLinkDest`, `NodeCodeBlockCode`, `NodeCodeBlockFenceOpenMarker`, `NodeCodeBlockFenceInfoMarker`, `NodeCodeBlockFenceCloseMarker`, `NodeMathBlockContent`, `NodeMathBlockOpenMarker`, `NodeMathBlockCloseMarker`, `NodeSuperBlockOpenMarker`, `NodeSuperBlockLayoutMarker`, `NodeSuperBlockCloseMarker`, `NodeOpenBrace`, `NodeCloseBrace`, `NodeBlockQueryEmbedScript`, `NodeTableHead`, `NodeTableRow`, `NodeTableCell`

> Типы, отключённые в SiYuan (сноски / оглавление / YAML / LinkRef / HeadingID и т. д.), перечислены в §11 и никогда не встречаются в `.sy` — они намеренно исключены из этого каталога.

---

## 5. Типы блоков подробно (с примерами для копирования)

### 5.1 Абзац

```json
{ "Type": "NodeParagraph", "ID": "...", "Properties": { "id": "...", "updated": "..." },
  "Children": [ { "Type": "NodeText", "Data": "Это пример абзаца." } ] }
```

### 5.2 Заголовок

```json
{ "Type": "NodeHeading", "ID": "...", "HeadingLevel": 2,
  "Properties": { "id": "...", "updated": "..." },
  "Children": [ { "Type": "NodeText", "Data": "Заголовок" } ] }
```

- `HeadingLevel` в диапазоне `1`–`6`.
- `NodeHeadingC8hMarker` (`Data` вида `"## "`) **необязателен** — законно и его наличие, и отсутствие. При генерации рекомендуется **опускать** его для краткости.
- SiYuan рекомендует использовать **заголовки второго уровня в начале тела документа**, а не первого.

### 5.3 Списки (ключевое: тип различается через `ListData.Typ`)

> **★ Жёсткое структурное ограничение списков:** непосредственными детьми `NodeList` могут быть **только** `NodeListItem` (обеспечивается `CanContain` в Lute, `ast/node.go:993` `return NodeListItem == nodeType`). Абзацы, блоки кода, подсписки и любые другие блоки **нельзя** прикреплять напрямую под `NodeList` — сначала их нужно обернуть в `NodeListItem`.

```
✅ Правильно                     ❌ Неправильно
NodeList                         NodeList
└─ NodeListItem                   ├─ NodeParagraph        ← недопустимо
   └─ NodeParagraph               └─ NodeCodeBlock         ← недопустимо
```

**Вложенные списки** записываются оборачиванием ещё одного `NodeList` (`NodeListItem` попадает в ветку `CanContain` по умолчанию и не может напрямую содержать другой `NodeListItem`):

```
✅ Правильно                     ❌ Неправильно
NodeList                         NodeList
└─ NodeListItem                   └─ NodeListItem
   ├─ NodeParagraph                  ├─ NodeParagraph
   └─ NodeList  ← подсписок          └─ NodeListItem  ← недопустимо
      └─ NodeListItem
         └─ NodeParagraph
```

**Маркированный список** (`Typ` опущен):

```json
{ "Type": "NodeList", "ID": "...", "ListData": {},
  "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeListItem", "ID": "...",
      "ListData": { "BulletChar": 42, "Marker": "Kg==" },
      "Properties": { "id": "..." },
      "Children": [
        { "Type": "NodeParagraph", "ID": "...", "Properties": { "id": "..." },
          "Children": [ { "Type": "NodeText", "Data": "Пункт один" } ] }
      ] }
  ] }
```

**Нумерованный список** (`Typ: 1`):

```json
{ "Type": "NodeListItem", "ID": "...", "Data": "1",
  "ListData": { "Typ": 1, "Tight": true, "Start": 1, "Delimiter": 46, "Padding": 3, "Marker": "MQ==", "Num": 1 },
  "Properties": { "id": "..." }, "Children": [ ... ] }
```

**Список задач** (`Typ: 3`); первый дочерний узел — `NodeTaskListItemMarker`:

```json
{ "Type": "NodeListItem", "ID": "...",
  "ListData": { "Typ": 3, "Tight": true, "BulletChar": 45, "Padding": 2, "Checked": true, "Marker": "LQ==", "Num": -1 },
  "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeTaskListItemMarker", "TaskListItemChecked": true },
    { "Type": "NodeParagraph", "ID": "...", "Properties": { "id": "..." },
      "Children": [ { "Type": "NodeText", "Data": "Задача один" } ] }
  ] }
```

### 5.4 Поля `ListData` полностью (★ здесь ошибаются чаще всего)

| Поле | Тип (в коде) | Форма в JSON | Значение |
|---|---|---|---|
| `Typ` | int | число | **Дискриминатор типа списка:** опущено = маркированный, `1` = нумерованный, `3` = задачи |
| `Tight` | bool | boolean | Плотный (без пустых строк); необязательно |
| `BulletChar` | byte | число | **ASCII-код** маркера для маркированных списков/списков задач (`42` = `*`, `45` = `-`) |
| `Delimiter` | byte | число | ASCII-код разделителя нумерованного списка (`46` = `.`) |
| `Start` | int | число | Начальный номер нумерованного списка |
| `Num` | int | число | Номер данного пункта; для маркированных/задач обычно опущен или `-1` |
| `Padding` | int | число | Отступ; необязательно |
| `Checked` | bool | boolean | Отмечен ли пункт задачи (агрегируется на уровне списка) |
| `Marker` | []byte | **строка base64** | Текст маркера, **закодированный в base64**; может включать разделитель (`"MS4="` = `1.`) или нет (`"MQ=="` = `1`) |

> Ключевое различие: **`BulletChar`/`Delimiter` в коде имеют тип `byte` и появляются в JSON как целочисленные коды**; **`Marker` в коде — `[]byte` и появляется в JSON как строка base64**. `Marker`/`BulletChar`/`Delimiter` все имеют `omitempty` и могут быть опущены.

### 5.5 Маркер задачи

```json
// Отмечено
{ "Type": "NodeTaskListItemMarker", "TaskListItemChecked": true }
// Не отмечено
{ "Type": "NodeTaskListItemMarker" }
```

> В реальных файлах `.sy` `NodeTaskListItemMarker` **обычно не имеет поля `Data`** (SiYuan восстанавливает маркер из DOM-атрибута `data-task` и не хранит исходный `[X]`/`[ ]`): отмеченные пункты несут только `Type`+`TaskListItemChecked`, неотмеченные — только `Type`. Если всё же записывать `Data` вручную, отмеченная форма должна быть с заглавной `[X]` (не `[x]`) — но это используется только при нормализованном экспорте в Markdown.

### 5.6 Цитата

```json
{ "Type": "NodeBlockquote", "ID": "...", "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeBlockquoteMarker", "Data": "> " },
    { "Type": "NodeParagraph", "ID": "...", "Properties": { "id": "..." },
      "Children": [ { "Type": "NodeText", "Data": "Цитируемое содержимое" } ] }
  ] }
```

> `NodeBlockquoteMarker.Data` может быть `">"` или `"> "` — оба варианта допустимы.

### 5.7 Выноска (GFM Alert)

```json
{ "Type": "NodeCallout", "ID": "...",
  "CalloutType": "NOTE", "CalloutTitle": "Note", "CalloutIcon": "✏️",
  "Properties": { "id": "...", "updated": "..." },
  "Children": [ { "Type": "NodeParagraph", "ID": "...", "Properties": { "id": "..." },
    "Children": [ { "Type": "NodeText", "Data": "Содержимое выноски" } ] } ] }
```

| `CalloutType` | `CalloutTitle` | `CalloutIcon` |
|---|---|---|
| `NOTE` | `Note` | `✏️` |
| `TIP` | `Tip` | `💡` |
| `IMPORTANT` | `Important` | `❗` |
| `WARNING` | `Warning` | `⚠️` |
| `CAUTION` | `Caution` | `🚨` |

> `CalloutIcon` — **буквальный символ эмодзи**, не base64 и не кодовая точка.

### 5.8 Суперблок (вкладываемый; трёхчастная структура)

```json
{ "Type": "NodeSuperBlock", "ID": "...", "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeSuperBlockOpenMarker" },
    { "Type": "NodeSuperBlockLayoutMarker", "Data": "col" },
    { "Type": "NodeSuperBlock", "ID": "...", "Properties": { "id": "..." }, "Children": [ ... вложенный суперблок, Data "row" ... ] },
    { "Type": "NodeSuperBlockCloseMarker" }
  ] }
```

> `NodeSuperBlockLayoutMarker.Data` может быть только `"row"` (горизонтально) или `"col"` (вертикально). Суперблоки могут вкладываться; это единственный контейнер, способный содержать любой блок (включая самого себя).

### 5.9 Блок встраивания (пятичастная структура `{{ ... }}`)

```json
{ "Type": "NodeBlockQueryEmbed", "ID": "...", "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeOpenBrace" },
    { "Type": "NodeOpenBrace" },
    { "Type": "NodeBlockQueryEmbedScript", "Data": "select * from blocks where id='20210428212840-8rqwn5o'" },
    { "Type": "NodeCloseBrace" },
    { "Type": "NodeCloseBrace" }
  ] }
```

### 5.10 Блок кода (четырёхчастная структура; только fenced)

```json
{ "Type": "NodeCodeBlock", "ID": "...", "IsFencedCodeBlock": true,
  "CodeBlockFenceChar": 96, "CodeBlockFenceLen": 3,
  "CodeBlockOpenFence": "YGBg", "CodeBlockInfo": "Z28=", "CodeBlockCloseFence": "YGBg",
  "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeCodeBlockFenceOpenMarker", "Data": "```", "CodeBlockFenceLen": 3 },
    { "Type": "NodeCodeBlockFenceInfoMarker", "CodeBlockInfo": "Z28=" },
    { "Type": "NodeCodeBlockCode", "Data": "package main\n...\n" },
    { "Type": "NodeCodeBlockFenceCloseMarker", "Data": "```", "CodeBlockFenceLen": 3 }
  ] }
```

Примечания:
- `NodeCodeBlockCode` несёт содержимое кода (в `Data`, сырой текст с экранированными `\n`); это строчный дочерний узел `NodeCodeBlock`.
- Окружающие маркеры ограждения (Open/Info/Close) — тоже строчные дочерние узлы.
- `CodeBlockInfo` — **язык, закодированный в base64** (`"Z28="` = `go`). Семь полей родителя (`IsFencedCodeBlock`/`CodeBlockFenceChar`/`CodeBlockFenceLen`/`CodeBlockOpenFence`/`CodeBlockInfo`/`CodeBlockCloseFence`) все имеют `omitempty` и могут опускаться по необходимости — новые файлы `.sy` часто пишут только `"IsFencedCodeBlock": true`.
- SiYuan **не поддерживает блоки кода с отступом** (`SetIndentCodeBlock(false)`); все блоки кода — fenced.

### 5.11 Математический блок (трёхчастная структура)

```json
{ "Type": "NodeMathBlock", "ID": "...", "Properties": { "id": "..." },
  "Children": [
    { "Type": "NodeMathBlockOpenMarker" },
    { "Type": "NodeMathBlockContent", "Data": "a^2 + b^2 = c^2" },
    { "Type": "NodeMathBlockCloseMarker" }
  ] }
```

### 5.12 Блоки HTML / IFrame / Video / Audio (листовые; содержимое в `Data` верхнего уровня)

```json
{ "Type": "NodeHTMLBlock", "ID": "...", "Data": "<div>\n<ruby>你<rt>nǐ</rt>...</div>", "Properties": { "id": "..." } }
{ "Type": "NodeIFrame", "ID": "...", "Data": "<iframe src=\"...\"></iframe>", "Properties": { "id": "..." } }
{ "Type": "NodeVideo", "ID": "...", "Data": "<video controls src=\"assets/x.mp4\"></video>", "Properties": { "id": "..." } }
{ "Type": "NodeAudio", "ID": "...", "Data": "<audio controls src=\"assets/x.wav\"></audio>", "Properties": { "id": "..." } }
```

> Эти четыре типа **не имеют `Children`**; HTML-содержимое (JSON-экранированное) находится непосредственно в `Data` верхнего уровня.

### 5.13 Таблица

```json
{ "Type": "NodeTable", "ID": "...", "TableAligns": [0, 0, 0],
  "Properties": { "id": "...", "colgroup": "||" },
  "Children": [
    { "Type": "NodeTableHead", "Data": "thead", "Children": [
      { "Type": "NodeTableRow", "Data": "tr", "Children": [
        { "Type": "NodeTableCell", "Data": "th", "Children": [ { "Type": "NodeText", "Data": "Заголовок" } ] }
      ] }
    ] },
    { "Type": "NodeTableRow", "Data": "tr", "Children": [
      { "Type": "NodeTableCell", "Data": "td", "Children": [ { "Type": "NodeText", "Data": "Ячейка" } ] }
    ] }
  ] }
```

- Вложенность фиксирована: `NodeTable > NodeTableHead/NodeTableRow > NodeTableCell > строчные узлы`.
- `TableAligns`: массив int с выравниванием по столбцам, `0` = по умолчанию/слева.
- `Data` (`thead`/`tr`/`th`/`td`) в компактных файлах может быть **опущено**.
- Ширины столбцов хранятся в `Properties.colgroup` (через `|`).

### 5.14 Блок AttributeView (база данных; листовой)

```json
{ "Type": "NodeAttributeView", "ID": "...",
  "Properties": { "custom-sy-av-view": "20251230141609-lcme2fh", "id": "...", "updated": "..." },
  "AttributeViewID": "20251230141609-2kvghrg",
  "AttributeViewType": "table" }
```

- **Не имеет `Children`.**
- `AttributeViewID` указывает на данные AV-таблицы (хранятся в отдельном `.json` — **не** выдумывайте этот ID).
- `AttributeViewType`: `table` / `kanban` / `gallery` и т. д.
- `custom-sy-av-view` фиксирует ID текущего представления.

> ИИ рекомендуется **не создавать новые блоки AttributeView**, поскольку данные таблицы находятся не в `.sy` — им нужны сопутствующие файлы.

### 5.15 Горизонтальная линия

```json
{ "Type": "NodeThematicBreak", "ID": "...", "Properties": { "id": "..." } }
```

---

## 6. Строчные узлы подробно

### 6.1 `NodeText` (обычный текст)

```json
{ "Type": "NodeText", "Data": "обычный текст" }
```

`Data` **может отсутствовать** (часто как заполнитель с пробелом нулевой ширины: `{ "Type": "NodeText" }`).

### 6.2 `NodeTextMark` (единый носитель современного строчного форматирования)

В файлах `.sy` жирный/курсив/ссылка/строчный код/ссылка на блок и т. д. — **почти всегда** `NodeTextMark`, а **не** `NodeStrong`/`NodeEmphasis`/`NodeLink`. Разновидность определяется полем `TextMarkType`.

| `TextMarkType` | Значение | Обязательные поля |
|---|---|---|
| `text` | обычный текст | `TextMarkTextContent` |
| `strong` | жирный | `TextMarkTextContent` |
| `em` | курсив | `TextMarkTextContent` |
| `u` | подчёркивание | `TextMarkTextContent` |
| `s` | зачёркивание (двойная тильда `~~`) | `TextMarkTextContent` |
| `mark` | выделение | `TextMarkTextContent` |
| `sup` / `sub` | верхний/нижний индекс | `TextMarkTextContent` |
| `kbd` | клавиша клавиатуры | `TextMarkTextContent` |
| `code` | строчный код | `TextMarkTextContent` |
| `tag` | тег `#тег#` | `TextMarkTextContent` |
| `a` | гиперссылка | `TextMarkAHref`, `TextMarkTextContent` (необязательно `TextMarkATitle`) |
| `block-ref` | ссылка на блок | `TextMarkBlockRefID`, `TextMarkBlockRefSubtype`, `TextMarkTextContent` |
| `inline-math` | строчная формула | `TextMarkInlineMathContent` (**без** `TextMarkTextContent`) |
| `inline-memo` | строчная заметка | `TextMarkInlineMemoContent`, `TextMarkTextContent` |
| `file-annotation-ref` | ссылка на аннотацию файла | `TextMarkFileAnnotationRefID`, `TextMarkTextContent` |

Примеры:

```json
{ "Type": "NodeTextMark", "TextMarkType": "a", "TextMarkAHref": "https://ld246.com", "TextMarkTextContent": "гиперссылка" }
{ "Type": "NodeTextMark", "TextMarkType": "block-ref", "TextMarkBlockRefID": "20200812220555-lj3enxa", "TextMarkBlockRefSubtype": "s", "TextMarkTextContent": "ссылка на блок" }
{ "Type": "NodeTextMark", "TextMarkType": "inline-math", "TextMarkInlineMathContent": "a^2 + b^2 = c^2" }
{ "Type": "NodeTextMark", "TextMarkType": "inline-memo", "TextMarkInlineMemoContent": "строчная заметка", "TextMarkTextContent": "заметка" }
```

- `TextMarkBlockRefSubtype`: `"s"` = статический текст якоря, `"d"` = динамический текст якоря (текст якоря следует за содержимым целевого блока; обратите внимание, что «блок встраивания» — это отдельный узел `NodeBlockQueryEmbed`, не связанный с этим).
- `TextMarkType` может объединять несколько меток через пробел, например `"strong em"`.
- `TextMarkTextContent` присутствует не у каждого типа (у `inline-math` его нет).
- Зачёркивание **поддерживает только двойную тильду `~~x~~`**, не одинарную `~x~` (`SetGFMStrikethrough1(false)`).
- Экранирование обратной косой чертой — **не** подтип `NodeTextMark`: оно отображается в отдельный узел `NodeBackslash` и никогда не появляется как значение `TextMarkType`.

### 6.3 Стилизованный строчный текст (★ обязательно парный)

`NodeTextMark`, несущий цвет/эффекты (с `Properties.style`), **должен сразу сопровождаться** узлом `NodeKramdownSpanIAL`, и оба должны содержать в точности одинаковый текст стиля:

```json
{ "Type": "NodeTextMark", "Properties": { "style": "color: var(--b3-font-color1); background-color: var(--b3-font-background1);" },
  "TextMarkType": "strong", "TextMarkTextContent": "цвет 1" },
{ "Type": "NodeKramdownSpanIAL", "Data": "{: style=\"color: var(--b3-font-color1); background-color: var(--b3-font-background1);\"}" }
```

> При генерации стилизованного строчного текста эти два узла должны идти парой, иначе круговое преобразование kramdown потеряет стиль.

### 6.4 `NodeImage` (семичастная структура)

```json
{ "Type": "NodeImage", "Data": "span", "Children": [
  { "Type": "NodeBang" },
  { "Type": "NodeOpenBracket" },
  { "Type": "NodeLinkText", "Data": "альтернативный текст" },
  { "Type": "NodeCloseBracket" },
  { "Type": "NodeOpenParen" },
  { "Type": "NodeLinkDest", "Data": "assets/image-2021.png" },
  { "Type": "NodeCloseParen" }
] }
```

- Сам узел изображения имеет `Data` = `"span"`.
- Маркеры `NodeBang`/`NodeOpenBracket`/`NodeCloseBracket`/`NodeOpenParen`/`NodeCloseParen` **могут опускать** `Data`.
- Только `NodeLinkText` и `NodeLinkDest` несут `Data`.

---

## 7. Соглашение о кодировании base64 (★ обязательно к прочтению)

| Поле | Кодировка | Пример |
|---|---|---|
| `ListData.Marker` | base64 | `Kg==` = `*`, `MS4=` = `1.`, `MQ==` = `1` |
| `CodeBlockInfo` | base64 | `Z28=` = `go`, `amF2YQ==` = `java` |
| `CodeBlockOpenFence`/`CloseFence` | base64 | `YGBg` = ` ``` ` |
| `ListData.BulletChar`/`Delimiter` | **целочисленный ASCII-код** (**не** base64) | `42` = `*`, `46` = `.` |
| `Data` (текст абзаца, содержимое кода, ссылка, SQL и т. д.) | **сырой** (не кодируется) | `"package main\n..."` |

> Практическое правило: **поля-маркеры** вроде `Marker`/`Fence`/`Info` — base64; **поля содержимого** вроде `Data`, `TextMarkTextContent`, `TextMarkInlineMathContent` — сырые; `BulletChar`/`Delimiter` — целочисленные коды.

---

## 8. Properties (IAL) полностью

Плоский `map[string]string`.

**Уровень документа (обязательно):** `id`, `title`, `type` (всегда `"doc"`), `updated`. Необязательно: `icon` (hex кодовой точки эмодзи, например `"1f4f0"`), `title-img` (CSS).

**Уровень блока (обязательно):** `id` (= `ID` узла), `updated`. Необязательно: `style` (строчный CSS), `fold: "1"` (свёрнут), `colgroup` (ширины столбцов таблицы), произвольные пользовательские атрибуты `custom-*`.

> Авторитетный ключ — `id` в **нижнем регистре**. В некоторых старых импортированных файлах остался и `ID` в верхнем регистре — предпочитайте нижний регистр.

---

## 9. Шпаргалка по вложенности контейнеров

| Контейнер | Может содержать | Не может содержать |
|---|---|---|
| `NodeList` | **только** `NodeListItem` | любой другой блок (абзацы/блоки кода/подсписки сначала оборачиваются в `NodeListItem`) |
| `NodeListItem` | любой блок, кроме `NodeListItem` (абзац/блок кода/вложенный `NodeList`/суперблок…) | `NodeListItem` (вложенность требует ещё одного `NodeList`) |
| `NodeBlockquote` | любой блок, кроме `NodeListItem`, + один `NodeBlockquoteMarker` | `NodeListItem` |
| `NodeCallout` | любой блок, кроме `NodeListItem` | `NodeListItem` |
| `NodeSuperBlock` | **любой блок** (включая вложенные суперблоки), обёрнутый тремя маркерами | ничего запрещённого (самый разрешающий) |
| `NodeDocument` | любой блок, кроме `NodeListItem` | `NodeListItem` |

> Эти правила обеспечиваются функцией `CanContain` в Lute (`ast/node.go:988`). Нарушения вызывают аномалии разбора/отрисовки — ИИ обязан соблюдать их при генерации.

---

## 10. Соглашение о пробеле нулевой ширины

Файлы `.sy` активно используют `​` (U+200B) как разделитель-заполнитель. Вокруг изображений, строчного кода, тегов, kbd и т. п. **обычно с каждой стороны** есть `NodeText` с `Data` = `​`, чтобы отрисовка и поведение курсора оставались корректными. ИИ должен следовать этому соглашению при генерации такого содержимого.

---

## 11. Синтаксис Markdown, отключённый в SiYuan (не должен появляться в `.sy`)

SiYuan отключает следующий синтаксис через `SetXxx(false)` в `NewLute()` (`kernel/util/lute.go`). Соответствующие типы узлов **никогда** не появляются в файлах `.sy` — ИИ не должен их генерировать:

| Отключённый элемент | Соответствующие типы узлов | Примечание |
|---|---|---|
| `SetFootnotes(false)` | `NodeFootnotesDefBlock`/`NodeFootnotesDef`/`NodeFootnotesRef` | сноски, полностью отключены |
| `SetToC(false)` | `NodeToC` | оглавление `[toc]` |
| `SetIndentCodeBlock(false)` | блоки кода с отступом | поддерживаются только fenced-блоки кода |
| `SetHeadingID(false)` | `NodeHeadingID` | пользовательский ID заголовка `{#id}` |
| `SetSetext(false)` | Setext-заголовки (форма с подчёркиванием `===`/`---`) | поддерживается только ATX-стиль `#` |
| `SetYamlFrontMatter(false)` | `NodeYamlFrontMatter` | YAML front matter |
| `SetLinkRef(false)` | `NodeLinkRefDef`/`NodeLinkRefDefBlock` | определения ссылок-сносок |
| `SetGFMStrikethrough1(false)` | зачёркивание одинарной тильдой `~x~` | поддерживается только двойная тильда `~~x~~` |

> Примечание: `NewLute()` также устанавливает `SetAutoSpace(false)`, `SetCodeSyntaxHighlight(false)` и `SetExportNormalizeTaskListMarker(false)` — это несинтаксические переключатели, влияющие только на отрисовку/экспорт и не убирающие никакие типы узлов, поэтому они опущены в таблице выше.

---

## 12. Чек-лист записи для ИИ

При генерации `.sy`, который SiYuan загрузит без проблем, проверьте по пунктам:

1. ☐ Корневой `Type` = `"NodeDocument"`, `Spec` = `"2"`; корневой `ID` = имени файла (без `.sy`) и равен `Properties.id`
2. ☐ Корневой `Properties` содержит `id`/`title`/`type:"doc"`/`updated`
3. ☐ Каждый блок имеет 22-символьный `ID`; `Properties.id` = `ID`; `Properties.updated` — валидная 14-значная метка времени
4. ☐ Строчные/маркерные узлы **не несут** `ID`
5. ☐ Списки различаются через `ListData.Typ` (опущено = маркированный / `1` = нумерованный / `3` = задачи)
6. ☐ Непосредственные дети `NodeList` — **только** `NodeListItem`; вложенные списки оборачивают ещё один `NodeList`
7. ☐ `BulletChar`/`Delimiter` — byte (в JSON целочисленные коды); `Marker` — base64
8. ☐ Маркеры задач обычно без `Data` (отмеченные используют `TaskListItemChecked:true`, неотмеченные — только `Type`)
9. ☐ Блок кода — четыре части, математический блок — три, встраивание — пять, суперблок — три — структура цела
10. ☐ `NodeCodeBlockCode`/`NodeMathBlockContent` — строчные дети, обычно только `Type`+`Data`
11. ☐ base64-поля закодированы; поля содержимого остаются сырыми
12. ☐ Для строчного форматирования предпочитайте `NodeTextMark`, а не устаревшие `NodeStrong`/`NodeEmphasis`/`NodeLink`
13. ☐ За стилизованным `TextMark` должен следовать парный `NodeKramdownSpanIAL`
14. ☐ HTML/IFrame/Video/Audio/AttributeView — листья без `Children` (содержимое в `Data` или полях типа)
15. ☐ Не выдумывайте `AttributeViewID`/целевые ID `block-ref` (они должны указывать на реальные блоки/AV)
16. ☐ Не генерируйте отключённые типы (сноски/оглавление/YAML/LinkRef/HeadingID и т. д. — см. §11)

---

## 13. Подводные камни и типичные ошибки

| ❌ Неправильно | ✅ Правильно |
|---|---|
| Считать, что у каждого узла есть `Data` | `Data` может отсутствовать; у маркерных узлов его часто нет |
| Строчный узел несёт `ID` | Строчные/маркерные узлы не имеют `ID` |
| Использовать устаревшие узлы вроде `NodeStrong`/`NodeLink` | Используйте `NodeTextMark` + `TextMarkType` |
| `ListData.Typ` принимает только `1` | опущено = маркированный, `1` = нумерованный, `3` = задачи |
| Трактовать `BulletChar` как base64 | Это `byte`, в JSON — целочисленный код (`42` = `*`) |
| Принудительно писать `"Data":"[X]"` в маркер задачи | В реальных `.sy` обычно нет `Data`; отмеченные несут только `TaskListItemChecked:true` |
| Стилизованный `TextMark` без IAL | Обязательна пара с `NodeKramdownSpanIAL` |
| Добавлять `Children` в `NodeAttributeView` | Это лист — используйте `AttributeViewID`/`AttributeViewType` |
| Менять `ID` без синхронизации `Properties.id` | Оба должны совпадать |
| `inline-math` несёт `TextMarkTextContent` | У него только `TextMarkInlineMathContent` |
| Выдумывать целевые ID block-ref / AV | Цели должны реально существовать |
| Вешать абзац напрямую под `NodeList` | `NodeList` может содержать только `NodeListItem` — сначала оберните |
| Генерировать сноски/оглавление/YAML и т. д. | SiYuan отключает их; они никогда не появляются в `.sy` |

---

## 14. Минимальный шаблон записываемого документа

```json
{
  "ID": "20260628120000-abc1234",
  "Spec": "2",
  "Type": "NodeDocument",
  "Properties": {
    "id": "20260628120000-abc1234",
    "title": "Новый документ",
    "type": "doc",
    "updated": "20260628120000"
  },
  "Children": [
    {
      "Type": "NodeHeading", "ID": "20260628120001-def5678", "HeadingLevel": 2,
      "Properties": { "id": "20260628120001-def5678", "updated": "20260628120001" },
      "Children": [ { "Type": "NodeText", "Data": "Заголовок" } ]
    },
    {
      "Type": "NodeParagraph", "ID": "20260628120002-ghi9012",
      "Properties": { "id": "20260628120002-ghi9012", "updated": "20260628120002" },
      "Children": [
        { "Type": "NodeText", "Data": "Текст с " },
        { "Type": "NodeTextMark", "TextMarkType": "strong", "TextMarkTextContent": "жирным" },
        { "Type": "NodeText", "Data": "." }
      ]
    }
  ]
}
```

---

## Приложение: источники проверки

- Образец 1: `app/guide/.../20200825162036-4dx365o.sy` (элементы форматирования — покрывает почти все типы блоков)
- Образец 2: `app/guide/.../20200905090211-2vixtlf.sy` (типы блоков — включая компактные списки, AttributeView)
- Константы типов узлов и логика сериализации: `lute/ast/node.go`, `lute/render/json_renderer.go`, `dataparser/sy.go`
- Проверка вложенности списков: `lute/ast/node.go:988` (`CanContain`)
- Конфигурация отключённого синтаксиса: `kernel/util/lute.go:51` (`NewLute`)
- Поля с пометкой `【inferred】` (например, подполя `file-annotation-ref`) перед генерацией следует перепроверить на реальном образце.
