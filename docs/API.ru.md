[English](API.md)
| [中文](API.zh-CN.md)
| [日本語](API.ja.md)
| **Русский**

* [Спецификация](#спецификация)
    * [Параметры и возвращаемые значения](#параметры-и-возвращаемые-значения)
    * [Аутентификация](#аутентификация)
* [Блокноты](#блокноты)
    * [Список блокнотов](#список-блокнотов)
    * [Открыть блокнот](#открыть-блокнот)
    * [Закрыть блокнот](#закрыть-блокнот)
    * [Переименовать блокнот](#переименовать-блокнот)
    * [Создать блокнот](#создать-блокнот)
    * [Удалить блокнот](#удалить-блокнот)
    * [Получить конфигурацию блокнота](#получить-конфигурацию-блокнота)
    * [Сохранить конфигурацию блокнота](#сохранить-конфигурацию-блокнота)
* [Документы](#документы)
    * [Создать документ из Markdown](#создать-документ-из-markdown)
    * [Переименовать документ](#переименовать-документ)
    * [Удалить документ](#удалить-документ)
    * [Переместить документы](#переместить-документы)
    * [Получить человекочитаемый путь по пути](#получить-человекочитаемый-путь-по-пути)
    * [Получить человекочитаемый путь по ID](#получить-человекочитаемый-путь-по-id)
    * [Получить путь хранения по ID](#получить-путь-хранения-по-id)
    * [Получить ID по человекочитаемому пути](#получить-id-по-человекочитаемому-пути)
* [Ресурсы](#ресурсы)
    * [Загрузить ресурсы](#загрузить-ресурсы)
* [Блоки](#блоки)
    * [Вставить блоки](#вставить-блоки)
    * [Вставить блоки в начало](#вставить-блоки-в-начало)
    * [Добавить блоки в конец](#добавить-блоки-в-конец)
    * [Обновить блок](#обновить-блок)
    * [Удалить блок](#удалить-блок)
    * [Переместить блок](#переместить-блок)
    * [Свернуть блок](#свернуть-блок)
    * [Развернуть блок](#развернуть-блок)
    * [Получить kramdown блока](#получить-kramdown-блока)
    * [Получить дочерние блоки](#получить-дочерние-блоки)
    * [Перенести ссылки на блок](#перенести-ссылки-на-блок)
* [Атрибуты](#атрибуты)
    * [Установить атрибуты блока](#установить-атрибуты-блока)
    * [Получить атрибуты блока](#получить-атрибуты-блока)
* [База данных](#база-данных)
    * [Отрисовка](#отрисовка)
    * [Получение](#получение)
    * [Получить значения первичного ключа](#получить-значения-первичного-ключа)
    * [Поиск](#поиск)
    * [Установить значение ячейки](#установить-значение-ячейки)
    * [Добавить элементы](#добавить-элементы)
    * [Удалить элементы](#удалить-элементы)
    * [Сменить макет](#сменить-макет)
    * [Настроить группировку](#настроить-группировку)
    * [Получить фильтры и сортировку](#получить-фильтры-и-сортировку)
    * [Установить фильтры](#установить-фильтры)
    * [Установить сортировку](#установить-сортировку)
    * [Добавить поле](#добавить-поле)
    * [Удалить поле](#удалить-поле)
    * [Задать глобальный порядок полей](#задать-глобальный-порядок-полей)
    * [Задать порядок полей в представлении](#задать-порядок-полей-в-представлении)
* [SQL](#sql)
    * [Выполнить SQL-запрос](#выполнить-sql-запрос)
    * [Сбросить транзакцию](#сбросить-транзакцию)
* [Шаблоны](#шаблоны)
    * [Отрисовать шаблон](#отрисовать-шаблон)
    * [Отрисовать Sprig](#отрисовать-sprig)
* [Файлы](#файлы)
    * [Получить файл](#получить-файл)
    * [Записать файл](#записать-файл)
    * [Удалить файл](#удалить-файл)
    * [Переименовать файл](#переименовать-файл)
    * [Список файлов](#список-файлов)
* [Экспорт](#экспорт)
    * [Экспорт Markdown](#экспорт-markdown)
    * [Экспорт файлов и папок](#экспорт-файлов-и-папок)
* [Конвертация](#конвертация)
    * [Pandoc](#pandoc)
* [Уведомления](#уведомления)
    * [Отправить сообщение](#отправить-сообщение)
    * [Отправить сообщение об ошибке](#отправить-сообщение-об-ошибке)
* [Сеть](#сеть)
    * [Форвард-прокси](#форвард-прокси)
* [Система](#система)
    * [Получить прогресс загрузки](#получить-прогресс-загрузки)
    * [Получить версию системы](#получить-версию-системы)
    * [Получить текущее время системы](#получить-текущее-время-системы)

---

## Спецификация

### Параметры и возвращаемые значения

* Конечная точка: `http://127.0.0.1:6806`
* Все методы — POST
* Для интерфейсов с параметрами параметр — это JSON-строка, помещаемая в тело запроса, заголовок
  Content-Type — `application/json`
* Возвращаемое значение

   ````json
   {
     "code": 0,
     "msg": "",
     "data": {}
   }
   ````

    * `code`: ненулевое значение при исключениях
    * `msg`: в нормальной ситуации пустая строка, при ошибке возвращается текст ошибки
    * `data`: может быть `{}`, `[]` или `NULL`, в зависимости от интерфейса

### Аутентификация

API-токен можно посмотреть в <kbd>Настройки - О программе</kbd>, заголовок запроса: `Authorization: Token xxx`

## Блокноты

### Список блокнотов

* `/api/notebook/lsNotebooks`
* Без параметров
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "notebooks": [
        {
          "id": "20210817205410-2kvfpfn", 
          "name": "Test Notebook",
          "icon": "1f41b",
          "sort": 0,
          "closed": false
        },
        {
          "id": "20210808180117-czj9bvb",
          "name": "SiYuan User Guide",
          "icon": "1f4d4",
          "sort": 1,
          "closed": false
        }
      ]
    }
  }
  ```

### Открыть блокнот

* `/api/notebook/openNotebook`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0"
  }
  ```

    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Закрыть блокнот

* `/api/notebook/closeNotebook`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0"
  }
  ```

    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Переименовать блокнот

* `/api/notebook/renameNotebook`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0",
    "name": "New name for notebook"
  }
  ```

    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Создать блокнот

* `/api/notebook/createNotebook`
* Параметры

  ```json
  {
    "name": "Notebook name"
  }
  ```
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "notebook": {
        "id": "20220126215949-r1wvoch",
        "name": "Notebook name",
        "icon": "",
        "sort": 0,
        "closed": false
      }
    }
  }
  ```

### Удалить блокнот

* `/api/notebook/removeNotebook`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0"
  }
  ```

    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Получить конфигурацию блокнота

* `/api/notebook/getNotebookConf`
* Параметры

  ```json
  {
    "notebook": "20210817205410-2kvfpfn"
  }
  ```

    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "box": "20210817205410-2kvfpfn",
      "conf": {
        "name": "Test Notebook",
        "closed": false,
        "refCreateSavePath": "",
        "createDocNameTemplate": "",
        "dailyNoteSavePath": "/daily note/{{now | date \"2006/01\"}}/{{now | date \"2006-01-02\"}}",
        "dailyNoteTemplatePath": ""
      },
      "name": "Test Notebook"
    }
  }
  ```

### Сохранить конфигурацию блокнота

* `/api/notebook/setNotebookConf`
* Параметры

  ```json
  {
    "notebook": "20210817205410-2kvfpfn",
    "conf": {
        "name": "Test Notebook",
        "closed": false,
        "refCreateSavePath": "",
        "createDocNameTemplate": "",
        "dailyNoteSavePath": "/daily note/{{now | date \"2006/01\"}}/{{now | date \"2006-01-02\"}}",
        "dailyNoteTemplatePath": ""
      }
  }
  ```

    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "name": "Test Notebook",
      "closed": false,
      "refCreateSavePath": "",
      "createDocNameTemplate": "",
      "dailyNoteSavePath": "/daily note/{{now | date \"2006/01\"}}/{{now | date \"2006-01-02\"}}",
      "dailyNoteTemplatePath": ""
    }
  }
  ```

## Документы

### Создать документ из Markdown

* `/api/filetree/createDocWithMd`
* Параметры

  ```json
  {
    "notebook": "20210817205410-2kvfpfn",
    "path": "/foo/bar",
    "markdown": ""
  }
  ```

    * `notebook`: ID блокнота
    * `path`: путь документа; должен начинаться с / и разделять уровни через / (path здесь соответствует
      полю hpath базы данных)
    * `markdown`: содержимое в формате GFM Markdown
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": "20210914223645-oj2vnx2"
  }
  ```

    * `data`: ID созданного документа
    * При повторном вызове этого интерфейса с тем же `path` существующий документ не будет перезаписан

### Переименовать документ

* `/api/filetree/renameDoc`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0",
    "path": "/20210902210113-0avi12f.sy",
    "title": "New document title"
  }
  ```

    * `notebook`: ID блокнота
    * `path`: путь документа
    * `title`: новый заголовок документа
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

Переименование документа по `id`:

* `/api/filetree/renameDocByID`
* Параметры

  ```json
  {
    "id": "20210902210113-0avi12f",
    "title": "New document title"
  }
  ```

    * `id`: ID документа
    * `title`: новый заголовок документа
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Удалить документ

* `/api/filetree/removeDoc`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0",
    "path": "/20210902210113-0avi12f.sy"
  }
  ```

    * `notebook`: ID блокнота
    * `path`: путь документа
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

Удаление документа по `id`:

* `/api/filetree/removeDocByID`
* Параметры

  ```json
  {
    "id": "20210902210113-0avi12f"
  }
  ```

    * `id`: ID документа
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Переместить документы

* `/api/filetree/moveDocs`
* Параметры

  ```json
  {
    "fromPaths": ["/20210917220056-yxtyl7i.sy"],
    "toNotebook": "20210817205410-2kvfpfn",
    "toPath": "/"
  }
  ```

    * `fromPaths`: исходные пути
    * `toNotebook`: ID целевого блокнота
    * `toPath`: целевой путь
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

Перемещение документов по `id`:

* `/api/filetree/moveDocsByID`
* Параметры

  ```json
  {
    "fromIDs": ["20210917220056-yxtyl7i"],
    "toID": "20210817205410-2kvfpfn"
  }
  ```

    * `fromIDs`: ID исходных документов
    * `toID`: ID целевого родительского документа или ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Получить человекочитаемый путь по пути

* `/api/filetree/getHPathByPath`
* Параметры

  ```json
  {
    "notebook": "20210831090520-7dvbdv0",
    "path": "/20210917220500-sz588nq/20210917220056-yxtyl7i.sy"
  }
  ```

    * `notebook`: ID блокнота
    * `path`: путь документа
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": "/foo/bar"
  }
  ```

### Получить человекочитаемый путь по ID

* `/api/filetree/getHPathByID`
* Параметры

  ```json
  {
    "id": "20210917220056-yxtyl7i"
  }
  ```

    * `id`: ID блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": "/foo/bar"
  }
  ```

### Получить путь хранения по ID

* `/api/filetree/getPathByID`
* Параметры

  ```json
  {
    "id": "20210808180320-fqgskfj"
  }
  ```

    * `id`: ID блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
    "notebook": "20210808180117-czj9bvb",
    "path": "/20200812220555-lj3enxa/20210808180320-fqgskfj.sy"
    }
  }
  ```

### Получить ID по человекочитаемому пути

* `/api/filetree/getIDsByHPath`
* Параметры

  ```json
  {
    "path": "/foo/bar",
    "notebook": "20210808180117-czj9bvb"
  }
  ```

    * `path`: человекочитаемый путь
    * `notebook`: ID блокнота
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
        "20200813004931-q4cu8na"
    ]
  }
  ```

## Ресурсы

### Загрузить ресурсы

* `/api/asset/upload`
* Параметр — HTTP Multipart-форма

    * `assetsDirPath`: путь к папке хранения ресурсов, корнем считается папка data, например:
        * `"/assets/"`: папка workspace/data/assets/
        * `"/assets/sub/"`: папка workspace/data/assets/sub/

      В обычной ситуации рекомендуется первый вариант — хранение в папке assets рабочего пространства;
      размещение в подкаталоге имеет побочные эффекты, см. главу про ресурсы в руководстве пользователя.
    * `file[]`: список загружаемых файлов
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "errFiles": [""],
      "succMap": {
        "foo.png": "assets/foo-20210719092549-9j5y79r.png"
      }
    }
  }
  ```

    * `errFiles`: список имён файлов, при обработке которых произошла ошибка
    * `succMap`: для успешно обработанных файлов ключ — имя файла при загрузке, значение —
      assets/foo-id.png; используется для замены адресов ссылок на ресурсы в существующем Markdown-содержимом
      на загруженные адреса

## Блоки

### Вставить блоки

* `/api/block/insertBlock`
* Параметры

  ```json
  {
    "dataType": "markdown",
    "data": "foo**bar**{: style=\"color: var(--b3-font-color8);\"}baz",
    "nextID": "",
    "previousID": "20211229114650-vrek5x6",
    "parentID": ""
  }
  ```

    * `dataType`: тип вставляемых данных, значение может быть `markdown` или `dom`
    * `data`: вставляемые данные
    * `nextID`: ID следующего блока, используется для привязки позиции вставки
    * `previousID`: ID предыдущего блока, используется для привязки позиции вставки
    * `parentID`: ID родительского блока, используется для привязки позиции вставки

  Хотя бы одно из `nextID`, `previousID`, `parentID` должно иметь значение; приоритет: `nextID` > `previousID` >
  `parentID`
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "doOperations": [
          {
            "action": "insert",
            "data": "<div data-node-id=\"20211230115020-g02dfx0\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<strong style=\"color: var(--b3-font-color8);\">bar</strong>baz</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>",
            "id": "20211230115020-g02dfx0",
            "parentID": "",
            "previousID": "20211229114650-vrek5x6",
            "retData": null
          }
        ],
        "undoOperations": null
      }
    ]
  }
  ```

    * `action.data`: DOM, сгенерированный для вставленного блока
    * `action.id`: ID вставленного блока

### Вставить блоки в начало

* `/api/block/prependBlock`
* Параметры

  ```json
  {
    "data": "foo**bar**{: style=\"color: var(--b3-font-color8);\"}baz",
    "dataType": "markdown",
    "parentID": "20220107173950-7f9m1nb"
  }
  ```

    * `dataType`: тип вставляемых данных, значение может быть `markdown` или `dom`
    * `data`: вставляемые данные
    * `parentID`: ID родительского блока, используется для привязки позиции вставки
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "doOperations": [
          {
            "action": "insert",
            "data": "<div data-node-id=\"20220108003710-hm0x9sc\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<strong style=\"color: var(--b3-font-color8);\">bar</strong>baz</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>",
            "id": "20220108003710-hm0x9sc",
            "parentID": "20220107173950-7f9m1nb",
            "previousID": "",
            "retData": null
          }
        ],
        "undoOperations": null
      }
    ]
  }
  ```

    * `action.data`: DOM, сгенерированный для вставленного блока
    * `action.id`: ID вставленного блока

### Добавить блоки в конец

* `/api/block/appendBlock`
* Параметры

  ```json
  {
    "data": "foo**bar**{: style=\"color: var(--b3-font-color8);\"}baz",
    "dataType": "markdown",
    "parentID": "20220107173950-7f9m1nb"
  }
  ```

    * `dataType`: тип вставляемых данных, значение может быть `markdown` или `dom`
    * `data`: вставляемые данные
    * `parentID`: ID родительского блока, используется для привязки позиции вставки
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "doOperations": [
          {
            "action": "insert",
            "data": "<div data-node-id=\"20220108003642-y2wmpcv\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<strong style=\"color: var(--b3-font-color8);\">bar</strong>baz</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>",
            "id": "20220108003642-y2wmpcv",
            "parentID": "20220107173950-7f9m1nb",
            "previousID": "20220108003615-7rk41t1",
            "retData": null
          }
        ],
        "undoOperations": null
      }
    ]
  }
  ```

    * `action.data`: DOM, сгенерированный для вставленного блока
    * `action.id`: ID вставленного блока

### Обновить блок

* `/api/block/updateBlock`
* Параметры

  ```json
  {
    "dataType": "markdown",
    "data": "foobarbaz",
    "id": "20211230161520-querkps"
  }
  ```

    * `dataType`: тип обновляемых данных, значение может быть `markdown` или `dom`
    * `data`: обновляемые данные
    * `id`: ID обновляемого блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "doOperations": [
          {
            "action": "update",
            "data": "<div data-node-id=\"20211230161520-querkps\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<strong>bar</strong>baz</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>",
            "id": "20211230161520-querkps",
            "parentID": "",
            "previousID": "",
            "retData": null
            }
          ],
        "undoOperations": null
      }
    ]
  }
  ```

    * `action.data`: DOM, сгенерированный для обновлённого блока

### Удалить блок

* `/api/block/deleteBlock`
* Параметры

  ```json
  {
    "id": "20211230161520-querkps"
  }
  ```

    * `id`: ID удаляемого блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "doOperations": [
          {
            "action": "delete",
            "data": null,
            "id": "20211230162439-vtm09qo",
            "parentID": "",
            "previousID": "",
            "retData": null
          }
        ],
       "undoOperations": null
      }
    ]
  }
  ```

### Переместить блок

* `/api/block/moveBlock`
* Параметры

  ```json
  {
    "id": "20230406180530-3o1rqkc",
    "previousID": "20230406152734-if5kyx6",
    "parentID": "20230404183855-woe52ko"
  }
  ```

    * `id`: ID перемещаемого блока
    * `previousID`: ID предыдущего блока, используется для привязки позиции вставки
    * `parentID`: ID родительского блока, используется для привязки позиции вставки; `previousID` и `parentID`
      не могут быть пустыми одновременно, если заданы оба, приоритет у `previousID`
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
        {
            "doOperations": [
                {
                    "action": "move",
                    "data": null,
                    "id": "20230406180530-3o1rqkc",
                    "parentID": "20230404183855-woe52ko",
                    "previousID": "20230406152734-if5kyx6",
                    "nextID": "",
                    "retData": null,
                    "srcIDs": null,
                    "name": "",
                    "type": ""
                }
            ],
            "undoOperations": null
        }
    ]
  }
  ```

### Свернуть блок

* `/api/block/foldBlock`
* Параметры

  ```json
  {
    "id": "20231224160424-2f5680o"
  }
  ```

    * `id`: ID сворачиваемого блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Развернуть блок

* `/api/block/unfoldBlock`
* Параметры

  ```json
  {
    "id": "20231224160424-2f5680o"
  }
  ```

    * `id`: ID разворачиваемого блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Получить kramdown блока

* `/api/block/getBlockKramdown`
* Параметры

  ```json
  {
    "id": "20201225220954-dlgzk1o"
  }
  ```

    * `id`: ID запрашиваемого блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "id": "20201225220954-dlgzk1o",
      "kramdown": "* {: id=\"20201225220954-e913snx\"}Create a new notebook, create a new document under the notebook\n  {: id=\"20210131161940-kfs31q6\"}\n* {: id=\"20201225220954-ygz217h\"}Enter <kbd>/</kbd> in the editor to trigger the function menu\n  {: id=\"20210131161940-eo0riwq\"}\n* {: id=\"20201225220954-875yybt\"}((20200924101200-gss5vee \"Navigate in the content block\")) and ((20200924100906-0u4zfq3 \"Window and tab\"))\n  {: id=\"20210131161940-b5uow2h\"}"
    }
  }
  ```

### Получить дочерние блоки

* `/api/block/getChildBlocks`
* Параметры

  ```json
  {
    "id": "20230506212712-vt9ajwj"
  }
  ```

    * `id`: ID родительского блока
    * Блоки под заголовком также считаются дочерними блоками
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "id": "20230512083858-mjdwkbn",
        "type": "h",
        "subType": "h1"
      },
      {
        "id": "20230513213727-thswvfd",
        "type": "s"
      },
      {
        "id": "20230513213633-9lsj4ew",
        "type": "l",
        "subType": "u"
      }
    ]
  }
  ```

### Перенести ссылки на блок

* `/api/block/transferBlockRef`
* Параметры

  ```json
  {
    "fromID": "20230612160235-mv6rrh1",
    "toID": "20230613093045-uwcomng",
    "refIDs": ["20230613092230-cpyimmd"]
  }
  ```

    * `fromID`: ID блока-определения
    * `toID`: ID целевого блока
    * `refIDs`: ID ссылающихся блоков, указывающих на блок-определение; необязательно, если не указано —
      будут перенесены все ссылающиеся блоки
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

## Атрибуты

### Установить атрибуты блока

* `/api/attr/setBlockAttrs`
* Параметры

  ```json
  {
    "id": "20210912214605-uhi5gco",
    "attrs": {
      "custom-attr1": "line1\nline2"
    }
  }
  ```

    * `id`: ID блока
    * `attrs`: атрибуты блока, пользовательские атрибуты должны иметь префикс `custom-`
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Получить атрибуты блока

* `/api/attr/getBlockAttrs`
* Параметры

  ```json
  {
    "id": "20210912214605-uhi5gco"
  }
  ```

    * `id`: ID блока
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "custom-attr1": "line1\nline2",
      "id": "20210912214605-uhi5gco",
      "title": "PDF Annotation Demo",
      "type": "doc",
      "updated": "20210916120715"
    }
  }
  ```

## SQL

### Выполнить SQL-запрос

* `/api/query/sql`
* Параметры

  ```json
  {
    "stmt": "SELECT * FROM blocks WHERE content LIKE'%content%' LIMIT 7"
  }
  ```

    * `stmt`: SQL-выражение
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      { "col": "val" }
    ]
  }
  ```

Примечание: в целях безопасности данных доступ к этому интерфейсу запрещён в режиме публикации (Publish Mode).

### Сбросить транзакцию

* `/api/sqlite/flushTransaction`
* Без параметров
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

## Шаблоны

### Отрисовать шаблон

* `/api/template/render`
* Параметры

  ```json
  {
    "id": "20220724223548-j6g0o87",
    "path": "F:\\SiYuan\\data\\templates\\foo.md"
  }
  ```

    * `id`: ID документа, из которого вызывается отрисовка
    * `path`: абсолютный путь к файлу шаблона
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "content": "<div data-node-id=\"20220729234848-dlgsah7\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220729234840\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>",
      "path": "F:\\SiYuan\\data\\templates\\foo.md"
    }
  }
  ```

### Отрисовать Sprig

* `/api/template/renderSprig`
* Параметры

  ```json
  {
    "template": "/daily note/{{now | date \"2006/01\"}}/{{now | date \"2006-01-02\"}}"
  }
  ```
    * `template`: содержимое шаблона
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": "/daily note/2023/03/2023-03-24"
  }
  ```

## Файлы

### Получить файл

* `/api/file/getFile`
* Параметры

  ``json {
  "path": "/data/20210808180117-6v0mkxr/20200923234011-ieuun1p.sy"
  }
  ``
    * `path`: путь к файлу относительно пути рабочего пространства
* Возвращаемое значение

    * Код состояния ответа `200`: содержимое файла
    * Код состояния ответа `202`: информация об исключении

      ```json
      {
        "code": 404,
        "msg": "",
        "data": null
      }
      ```

        * `code`: ненулевое значение при исключениях

            * `-1`: ошибка разбора параметров
            * `403`: доступ запрещён (файл не находится в рабочем пространстве)
            * `404`: не найдено (файл не существует)
            * `405`: метод не разрешён (это каталог)
            * `500`: ошибка сервера (не удалось получить сведения о файле / прочитать файл)
        * `msg`: текст с описанием ошибки

### Записать файл

* `/api/file/putFile`
* Параметр — HTTP Multipart-форма

    * `path`: путь к файлу относительно пути рабочего пространства
    * `isDir`: создавать ли папку; при `true` создаётся только папка, `file` игнорируется
    * `modTime`: время последнего доступа и изменения, Unix-время
    * `file`: загружаемый файл
* Возвращаемое значение

   ```json
   {
     "code": 0,
     "msg": "",
     "data": null
   }
   ```

### Удалить файл

* `/api/file/removeFile`
* Параметры

  ```json
  {
    "path": "/data/20210808180117-6v0mkxr/20200923234011-ieuun1p.sy"
  }
  ```
    * `path`: путь к файлу относительно пути рабочего пространства
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Переименовать файл

* `/api/file/renameFile`
* Параметры

  ```json
  {
    "path": "/data/assets/image-20230523085812-k3o9t32.png",
    "newPath": "/data/assets/test-20230523085812-k3o9t32.png"
  }
  ```
    * `path`: путь к файлу относительно пути рабочего пространства
    * `newPath`: новый путь к файлу относительно пути рабочего пространства
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Список файлов

* `/api/file/readDir`
* Параметры

  ```json
  {
    "path": "/data/20210808180117-6v0mkxr/20200923234011-ieuun1p"
  }
  ```
    * `path`: путь к каталогу относительно пути рабочего пространства
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": [
      {
        "isDir": true,
        "isSymlink": false,
        "name": "20210808180303-6yi0dv5",
        "updated": 1691467624
      },
      {
        "isDir": false,
        "isSymlink": false,
        "name": "20210808180303-6yi0dv5.sy",
        "updated": 1663298365
      }
    ]
  }
  ```

## Экспорт

### Экспорт Markdown

* `/api/export/exportMdContent`
* Параметры

  ```json
  {
    "id": ""
  }
  ```

    * `id`: ID блока документа для экспорта
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "hPath": "/Please Start Here",
      "content": "## 🍫 Content Block\n\nIn SiYuan, the only important core concept is..."
    }
  }
  ```

    * `hPath`: человекочитаемый путь
    * `content`: содержимое Markdown

### Экспорт файлов и папок

* `/api/export/exportResources`
* Параметры

  ```json
  {
    "paths": [
      "/conf/appearance/boot",
      "/conf/appearance/langs",
      "/conf/appearance/emojis/conf.json",
      "/conf/appearance/icons/index.html"
    ],
    "name": "zip-file-name"
  }
  ```

    * `paths`: список путей к экспортируемым файлам или папкам; одинаковые имена файлов/папок будут перезаписаны
    * `name`: (необязательно) имя экспортируемого файла; по умолчанию `export-YYYY-MM-DD_hh-mm-ss.zip`
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "path": "temp/export/zip-file-name.zip"
    }
  }
  ```

    * `path`: путь к созданному файлу `*.zip`
        * Структура каталогов в `zip-file-name.zip` выглядит так:
            * `zip-file-name`
                * `boot`
                * `langs`
                * `conf.json`
                * `index.html`

## Конвертация

### Pandoc

* `/api/convert/pandoc`
* Рабочий каталог
    * При выполнении команды pandoc рабочим каталогом будет `workspace/temp/convert/pandoc/${dir}`
    * Сначала можно записать конвертируемый файл в этот каталог через API [`Записать файл`](#записать-файл)
    * Затем вызвать API для конвертации — конвертированный файл также будет записан в этот каталог
    * Наконец, вызвать API [`Получить файл`](#получить-файл), чтобы получить конвертированный файл
        * Или вызвать API [Создать документ из Markdown](#создать-документ-из-markdown)
        * Или вызвать внутренний API `importStdMd`, чтобы импортировать конвертированную папку напрямую
* Параметры

  ```json
  {
    "dir": "test",
    "args": [
      "--to", "markdown_strict-raw_html",
      "foo.epub",
      "-o", "foo.md"
   ]
  }
  ```

    * `args`: параметры командной строки Pandoc
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
       "path": "/temp/convert/pandoc/test"
    }
  }
  ```
    * `path`: путь относительно рабочего пространства

## Уведомления

### Отправить сообщение

* `/api/notification/pushMsg`
* Параметры

  ```json
  {
    "msg": "test",
    "timeout": 7000
  }
  ```
    * `timeout`: длительность отображения сообщения в миллисекундах. Поле можно опустить, по умолчанию 7000
      миллисекунд
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
        "id": "62jtmqi"
    }
  }
  ```
    * `id`: ID сообщения

### Отправить сообщение об ошибке

* `/api/notification/pushErrMsg`
* Параметры

  ```json
  {
    "msg": "test",
    "timeout": 7000
  }
  ```
    * `timeout`: длительность отображения сообщения в миллисекундах. Поле можно опустить, по умолчанию 7000
      миллисекунд
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
        "id": "qc9znut"
    }
  }
  ```
    * `id`: ID сообщения

## Сеть

### Форвард-прокси

* `/api/network/forwardProxy`
* Параметры

  ```json
  {
    "url": "https://b3log.org/siyuan/",
    "method": "GET",
    "timeout": 7000,
    "contentType": "text/html",
    "headers": [
        {
            "Cookie": ""
        }
    ],
    "payload": {},
    "payloadEncoding": "text",
    "responseEncoding": "text"
  }
  ```

    * `url`: URL для переадресации
    * `method`: HTTP-метод, по умолчанию `POST`
    * `timeout`: тайм-аут в миллисекундах, по умолчанию `7000`
    * `contentType`: Content-Type, по умолчанию `application/json`
    * `headers`: HTTP-заголовки
    * `payload`: тело HTTP-запроса, объект или строка
    * `payloadEncoding`: схема кодирования, используемая `payload`, по умолчанию `text`; допустимые значения:

        * `text`
        * `base64` | `base64-std`
        * `base64-url`
        * `base32` | `base32-std`
        * `base32-hex`
        * `hex`
    * `responseEncoding`: схема кодирования поля `body` в данных ответа, по умолчанию `text`; допустимые
      значения:

        * `text`
        * `base64` | `base64-std`
        * `base64-url`
        * `base32` | `base32-std`
        * `base32-hex`
        * `hex`
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "body": "",
      "bodyEncoding": "text",
      "contentType": "text/html",
      "elapsed": 1976,
      "headers": {
      },
      "status": 200,
      "url": "https://b3log.org/siyuan"
    }
  }
  ```

    * `bodyEncoding`: схема кодирования поля `body`, соответствует полю `responseEncoding` в запросе,
      по умолчанию `text`; допустимые значения:

        * `text`
        * `base64` | `base64-std`
        * `base64-url`
        * `base32` | `base32-std`
        * `base32-hex`
        * `hex`

## Система

### Получить прогресс загрузки

* `/api/system/bootProgress`
* Без параметров
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "details": "Finishing boot...",
      "progress": 100
    }
  }
  ```

### Получить версию системы

* `/api/system/version`
* Без параметров
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": "1.3.5"
  }
  ```

### Получить текущее время системы

* `/api/system/currentTime`
* Без параметров
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": 1631850968131
  }
  ```

    * `data`: точность в миллисекундах

## База данных

База данных (внутренне — «attribute view», представление атрибутов) хранит структурированные данные в виде полей (столбцов) и элементов (строк). Каждая база данных идентифицируется по `avID` и может встраиваться в документ через один или несколько блоков базы данных (`blockID`). Одна база данных может содержать несколько представлений (`viewID`) с разными типами макета: `table` (таблица), `gallery` (галерея) и `kanban` (канбан).

Типы полей (`keyType`):

| Значение     | Описание                        |
|--------------|---------------------------------|
| `block`      | Первичный ключ (привязанный блок) |
| `text`       | Текст                           |
| `number`     | Число                           |
| `date`       | Дата                            |
| `select`     | Одиночный выбор                 |
| `mSelect`    | Множественный выбор             |
| `url`        | URL                             |
| `email`      | Электронная почта               |
| `phone`      | Телефон                         |
| `mAsset`     | Ресурс                          |
| `template`   | Шаблон                          |
| `created`    | Время создания                  |
| `updated`    | Время обновления                |
| `checkbox`   | Флажок                          |
| `relation`   | Связь                           |
| `rollup`     | Свёртка (rollup)                |
| `lineNumber` | Номер строки                    |

### Отрисовка

* `/api/av/renderAttributeView`
* Параметры

  ```json
  {
    "id": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t",
    "viewID": "",
    "page": 1,
    "pageSize": 50,
    "query": "",
    "groupPaging": {},
    "createIfNotExist": true
  }
  ```

    * `id`: ID базы данных
    * `blockID`: блок базы данных, встраивающий эту базу. Используется для определения активного представления и доступа при публикации. Опустите при отрисовке отвязанной базы данных
    * `viewID`: представление для отрисовки. Если опущено, используется текущее представление (поле `viewID`)
    * `page`: номер страницы, начиная с 1. По умолчанию `1`
    * `pageSize`: элементов на страницу. `-1` или отсутствие — использовать значение представления по умолчанию (`50`)
    * `query`: необязательный полнотекстовый фильтр по значениям первичного ключа
    * `groupPaging`: необязательная конфигурация постраничности для сгруппированных (канбан) представлений
    * `createIfNotExist`: при `true` (по умолчанию) создаёт представление по умолчанию, если у базы данных его нет
* Возвращаемое значение (реальный ответ, табличный макет, показана одна строка):

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "name": "API 测试",
      "id": "20240118120204-kwyzf77",
      "viewType": "table",
      "viewID": "20240118120204-7rnmyc1",
      "isMirror": false,
      "views": [
        {
          "id": "20240118120204-7rnmyc1",
          "icon": "",
          "name": "表格",
          "desc": "",
          "hideAttrViewName": false,
          "type": "table",
          "pageSize": 50
        }
      ],
      "view": {
        "id": "20240118120204-7rnmyc1",
        "icon": "",
        "name": "表格",
        "desc": "",
        "hideAttrViewName": false,
        "filters": [],
        "sorts": [],
        "group": null,
        "pageSize": 50,
        "showIcon": true,
        "wrapField": false,
        "groupFolded": false,
        "groupHidden": 0,
        "columns": [
          {
            "id": "20240118120204-w6cggab",
            "name": "主键",
            "type": "block",
            "icon": "",
            "wrap": false,
            "hidden": false,
            "desc": "",
            "calc": null,
            "numberFormat": "",
            "template": "",
            "pin": false,
            "width": ""
          }
        ],
        "rows": [
          {
            "id": "20240118203831-fkfvvtx",
            "cells": [
              {
                "id": "20240118203911-xrg9obl",
                "value": {
                  "id": "20240118203911-xrg9obl",
                  "keyID": "20240118120204-w6cggab",
                  "blockID": "20240118203831-fkfvvtx",
                  "type": "block",
                  "createdAt": 1706843791000,
                  "updatedAt": 1706843791000,
                  "block": {
                    "id": "20240118203831-fkfvvtx",
                    "content": "3",
                    "created": 1706843791000,
                    "updated": 1706843791000
                  }
                },
                "valueType": "block",
                "color": "",
                "bgColor": ""
              }
            ]
          }
        ],
        "rowCount": 5
      }
    }
  }
  ```

    * `data.view`: отрисованный экземпляр представления. Форма зависит от `viewType` — `table` возвращает `columns`/`rows`/`rowCount`, `gallery` возвращает `columns`/`rows`, `kanban` возвращает `columns`/`groups` (каждая группа сама является экземпляром представления с `groupKey`/`groupValue`). `view` также несёт `filters`, `sorts`, `group`, `showIcon`, `wrapField`, `groupFolded`, `groupHidden`. Внимание: активные фильтры/группировка могут сделать `rows` пустым даже при `rowCount` > 0
    * `data.view.columns[]`: каждый содержит `id`, `name`, `type`, `icon`, `wrap`, `hidden`, `desc`, `calc`, `numberFormat`, `template`, `pin`, `width`; столбцы `select`/`mSelect` дополнительно несут `options`
    * `data.view.rows[].id`: **ID строки** (ID элемента). Для привязанной строки совпадает с ID привязанного блока; для отвязанной строки это сгенерированный ID элемента, не совпадающий ни с одним блоком
    * `data.view.rows[].cells[].value`: объект `Value` — все формы значений см. в разделе [Установить значение ячейки](#установить-значение-ячейки). `createdAt`/`updatedAt` — int64-метки времени в миллисекундах
    * `data.views`: метаданные каждого представления (без строк)
    * `data.isMirror`: `true`, если блок базы данных является зеркалом (копией только для чтения) базы данных

### Получение

* `/api/av/getAttributeView`
* Параметры

  ```json
  {
    "id": "20240118120204-kwyzf77"
  }
  ```

    * `id`: ID базы данных
* Возвращаемое значение (реальный ответ, сокращён — массивы `keyValues`/`views` усечены):

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "av": {
        "spec": 4,
        "id": "20240118120204-kwyzf77",
        "name": "API 测试",
        "keyValues": [
          {
            "key": {
              "id": "20240118120204-w6cggab",
              "name": "主键",
              "type": "block",
              "icon": "",
              "desc": "",
              "numberFormat": "",
              "template": ""
            },
            "values": [
              {
                "id": "20240118203911-xrg9obl",
                "keyID": "20240118120204-w6cggab",
                "blockID": "20240118203831-fkfvvtx",
                "type": "block",
                "createdAt": 1706843791000,
                "updatedAt": 1706843791000,
                "block": {
                  "id": "20240118203831-fkfvvtx",
                  "content": "3",
                  "created": 1706843791000,
                  "updated": 1706843791000
                }
              }
            ]
          }
        ],
        "keyIDs": null,
        "viewID": "20240118120204-7rnmyc1",
        "views": [
          {
            "id": "20240118120204-7rnmyc1",
            "icon": "",
            "name": "表格",
            "hideAttrViewName": false,
            "desc": "",
            "pageSize": 50,
            "type": "table",
            "table": {
              "spec": 0,
              "id": "20240118120204-grokgmm",
              "showIcon": true,
              "wrapField": false,
              "columns": [
                {
                  "id": "20240118120204-w6cggab",
                  "wrap": false,
                  "hidden": false,
                  "pin": false,
                  "width": ""
                }
              ],
              "rowIds": null
            },
            "itemIds": ["20240118203818-ct041hj", "20240118203855-sqzbja0", "20240118203831-fkfvvtx", "20240118203842-kc31ovy", "20240531235026-uiap07y"],
            "groupCreated": 0,
            "groupItemIds": null,
            "groupFolded": false,
            "groupHidden": 0,
            "groupSort": 0
          }
        ]
      }
    }
  }
  ```

    * `data.av`: полное определение `AttributeView` — поля (`keyValues`), порядок полей (`keyIDs`, может быть `null`), текущее представление (`viewID`) и все представления с их сырой конфигурацией макета (`table`/`gallery`/`kanban`) и порядком элементов (`itemIds`). Возвращается сырое определение (без отрисованных строк и постраничности); для вычисленных строк используйте [Отрисовку](#отрисовка)

### Получить значения первичного ключа

* `/api/av/getAttributeViewPrimaryKeyValues`
* Параметры

  ```json
  {
    "id": "20240118120204-kwyzf77",
    "keyword": "",
    "page": 1,
    "pageSize": 16
  }
  ```

    * `id`: ID базы данных
    * `keyword`: необязательный фильтр-подстрока по тексту первичного ключа (без учёта регистра)
    * `page`: номер страницы, начиная с 1. По умолчанию `1`
    * `pageSize`: элементов на страницу. `-1` или отсутствие — `16`. Значения сортируются по `block.updated` по убыванию
* Возвращаемое значение (реальный ответ, показано одно значение):

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "name": "API 测试",
      "blockIDs": ["20240118120201-kldj15t"],
      "rows": {
        "key": {
          "id": "20240118120204-w6cggab",
          "name": "主键",
          "type": "block",
          "icon": "",
          "desc": "",
          "numberFormat": "",
          "template": ""
        },
        "values": [
          {
            "id": "20240118203911-xrg9obl",
            "keyID": "20240118120204-w6cggab",
            "blockID": "20240118203831-fkfvvtx",
            "type": "block",
            "createdAt": 1706843791000,
            "updatedAt": 1706843791000,
            "block": {
              "id": "20240118203831-fkfvvtx",
              "content": "3",
              "created": 1706843791000,
              "updated": 1706843791000
            }
          }
        ]
      }
    }
  }
  ```

    * `data.rows`: объект `KeyValues`, содержащий поле первичного ключа (`block`) и его постраничные значения
    * `data.blockIDs`: ID всех блоков базы данных (зеркал), ссылающихся на эту базу

### Поиск

* `/api/av/searchAttributeView`
* Параметры

  ```json
  {
    "keyword": "API",
    "excludes": []
  }
  ```

    * `keyword`: ключевое слово поиска (сопоставляется с именем базы данных)
    * `excludes`: необязательный список ID баз данных, исключаемых из результатов
* Возвращаемое значение (реальный ответ):

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "results": [
        {
          "avID": "20240118120204-kwyzf77",
          "avName": "API 测试",
          "viewName": "",
          "viewID": "",
          "viewLayout": "",
          "blockID": "20240118120201-kldj15t",
          "hPath": "正在跟进的问题/数据库/API",
          "children": [
            {
              "avID": "20240118120204-kwyzf77",
              "avName": "API 测试",
              "viewName": "表格",
              "viewID": "20240118120204-7rnmyc1",
              "viewLayout": "table",
              "blockID": "20240118120201-kldj15t",
              "hPath": "正在跟进的问题/数据库/API"
            }
          ]
        }
      ]
    }
  }
  ```

    * `data.results[]`: каждый результат верхнего уровня группирует базу данных по `avID`; его `children[]` перечисляют отдельные представления (`viewName`/`viewID`/`viewLayout`)

### Установить значение ячейки

Обновляет одну ячейку (одно поле одной строки). Это основная конечная точка для записи значений ячеек. `value` в запросе — частичный объект `Value`, форма которого зависит от `keyType` поля. Наиболее распространённые формы значений:

| `keyType`  | Форма `value`                                                                                                        |
|------------|----------------------------------------------------------------------------------------------------------------------|
| `block`    | `{"block": {"content": "First row", "id": "<boundBlockID>"}, "isDetached": false}`                                  |
| `text`     | `{"text": {"content": "Some text"}}`                                                                                 |
| `number`   | `{"number": {"content": 42, "isNotEmpty": true}}` (очистка — `{"isNotEmpty": false}`)                               |
| `date`     | `{"date": {"content": 1676042451000, "isNotEmpty": true}}` (метка времени в миллисекундах)                          |
| `select`   | `{"mSelect": [{"content": "Done", "color": "1"}]}` (не более одной опции)                                           |
| `mSelect`  | `{"mSelect": [{"content": "A", "color": "1"}, {"content": "B", "color": "2"}]}`                                      |
| `url`      | `{"url": {"content": "https://siyuan.com"}}`                                                                         |
| `email`    | `{"email": {"content": "a@b.com"}}`                                                                                  |
| `phone`    | `{"phone": {"content": "1234567890"}}`                                                                               |
| `checkbox` | `{"checkbox": {"checked": true}}`                                                                                    |

> ⚠️ `itemID` — это **ID строки** (`rows[].id` из [Отрисовки](#отрисовка)). Для привязанной строки ID строки равен ID привязанного блока; для отвязанной строки это сгенерированный ID элемента. Передача неправильного ID сохранит значение как «сироту», которое не появится в отрисованной ячейке.

* `/api/av/setAttributeViewBlockAttr`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "keyID": "20240531232156-ahsyx8l",
    "itemID": "20240118203831-fkfvvtx",
    "value": {
      "type": "number",
      "number": {
        "content": 42,
        "isNotEmpty": true
      }
    }
  }
  ```

    * `avID`: ID базы данных
    * `keyID`: ID поля (обновляемого столбца)
    * `itemID`: **ID строки** (`rows[].id` из [Отрисовки](#отрисовка)). Устаревший параметр `rowID` будет удалён после 2026-12-01; используйте `itemID`
    * `value`: частичный объект `Value` (см. таблицу выше). Неизвестные или неподдерживаемые ключи игнорируются
* Возвращаемое значение (реальный ответ, числовое значение):

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "value": {
        "id": "20240531235048-4zisj1p",
        "keyID": "20240531232156-ahsyx8l",
        "blockID": "20240118203831-fkfvvtx",
        "type": "number",
        "createdAt": 1717170648596,
        "updatedAt": 1781610266432,
        "number": {
          "content": 42,
          "isNotEmpty": true,
          "format": "",
          "formattedContent": "42"
        }
      }
    }
  }
  ```

    * `data.value`: полностью нормализованное значение после обновления (с вычисленными полями, такими как `number.formattedContent`). Используйте его для обновления UI, а не повторную отправку тела запроса

### Добавить элементы

Добавляет один или несколько элементов (строк). Каждый источник может либо привязать существующий блок (`isDetached: false`), либо создать отвязанную строку, живущую только внутри представления (`isDetached: true`).

* `/api/av/addAttributeViewBlocks`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t",
    "viewID": "",
    "groupID": "",
    "previousID": "",
    "srcs": [
      {
        "id": "20240118120201-kldj15t",
        "isDetached": false,
        "content": "New row"
      }
    ],
    "ignoreDefaultFill": false
  }
  ```

    * `avID`: ID базы данных
    * `blockID`: блок базы данных, владеющий этой базой (определяет целевое представление/группу)
    * `viewID`: целевое представление. Если опущено, используется текущее
    * `groupID`: ID целевой группы для канбан-представлений. Опустите для таблицы/галереи
    * `previousID`: вставить после элемента с этим ID. Пустое значение — добавить в конец
    * `srcs[].id`: для привязанных блоков (`isDetached: false`) — ID привязываемого блока. Должен соответствовать формату ID узла
    * `srcs[].isDetached`: `true` — создать отвязанную строку; `false` — привязать существующий блок
    * `srcs[].content`: отображаемый текст первичного ключа (используется при `isDetached: true` или для переопределения содержимого привязанного блока)
    * `srcs[].itemID`: необязательный явный ID элемента. Генерируется автоматически, если опущен
    * `ignoreDefaultFill`: при `true` пропускает автозаполнение значений по умолчанию в поля фильтров/группировки
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

    * Конечная точка возвращает `null`; после успешного вызова вызовите [Отрисовку](#отрисовка), чтобы получить обновлённые строки (включая новые ID строк, необходимые для обновления ячеек)

### Удалить элементы

Удаляет один или несколько элементов (строк). Отвязанные строки удаляются; привязанные блоки отвязываются (сам блок документа не удаляется).

* `/api/av/removeAttributeViewBlocks`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "srcIDs": ["20240118203831-fkfvvtx"]
  }
  ```

    * `avID`: ID базы данных
    * `srcIDs`: ID строк (`rows[].id` из [Отрисовки](#отрисовка)) для удаления
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Сменить макет

Переключает тип макета текущего представления между `table`, `gallery` и `kanban`. При успехе сервер заново отрисовывает представление и возвращает его (та же форма, что у [Отрисовки](#отрисовка)).

* `/api/av/changeAttrViewLayout`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t",
    "layoutType": "kanban"
  }
  ```

    * `avID`: ID базы данных
    * `blockID`: блок базы данных, владеющий представлением
    * `layoutType`: целевой макет — `table`, `gallery` или `kanban`
* Возвращаемое значение: та же форма, что у [Отрисовки](#отрисовка). При переключении на `kanban` с настроенной группировкой `data.view` несёт массив `groups[]`; каждая группа — экземпляр представления с `groupKey`, `groupValue` и специфичными для канбана полями (`coverFrom`, `cardAspectRatio`, `cardSize`, `fitImage`, `displayFieldName`, `fillColBackgroundColor`, `fields`)

### Настроить группировку

Устанавливает или снимает правило группировки канбан-представления. Если `group.field` пуст, группировка удаляется. При успехе сервер заново отрисовывает представление и возвращает его.

* `/api/av/setAttrViewGroup`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t",
    "group": {
      "field": "20240118203822-io6ofxb",
      "method": 0,
      "order": 0,
      "hideEmpty": false
    }
  }
  ```

    * `avID`: ID базы данных
    * `blockID`: блок базы данных, владеющий представлением
    * `group`: правило группировки
    * `group.field`: ID поля (столбца), по которому группировать. Пустая строка снимает группировку
    * `group.method`: метод группировки — `0` по значению, `1` по числовому диапазону, `2` по относительной дате, `3` по дню, `4` по неделе, `5` по месяцу, `6` по году
    * `group.range`: необязательно. Требуется при `method` равном `1` (числовой диапазон): `{ "numStart": 0, "numEnd": 100, "numStep": 10 }`
    * `group.order`: порядок групп — `0` по возрастанию, `1` по убыванию, `2` вручную, `3` в порядке опций выбора
    * `group.hideEmpty`: скрывать ли пустые группы
* Возвращаемое значение: та же форма, что у [Отрисовки](#отрисовка)

### Получить фильтры и сортировку

Возвращает текущие правила фильтрации и сортировки представления, привязанного к блоку базы данных.

* `/api/av/getAttributeViewFilterSort`
* Параметры

  ```json
  {
    "id": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t"
  }
  ```

    * `id`: ID базы данных
    * `blockID`: блок базы данных, владеющий представлением
* Возвращаемое значение (реальный ответ, фильтры/сортировки не настроены):

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "filters": [],
      "sorts": []
    }
  }
  ```

  Когда они настроены (реальный перехваченный ответ), фильтр и сортировка выглядят так:

  ```json
  {
    "code": 0,
    "msg": "",
    "data": {
      "filters": [
        {
          "column": "20240118203822-io6ofxb",
          "operator": "=",
          "value": {
            "type": "select",
            "mSelect": [
              { "content": "Done", "color": "1" }
            ]
          }
        }
      ],
      "sorts": [
        {
          "column": "20240118120204-w6cggab",
          "order": "DESC"
        }
      ]
    }
  }
  ```

    * `data.filters`: массив `ViewFilter`. Верхний уровень содержит единственный корневой групповой узел `{ "combination": "and"|"or", "filters": [...] }`; элементы массива — либо листовые фильтры, либо вложенные групповые узлы, что позволяет рекурсивно комбинировать AND/OR.
    * `data.filters[].column`: ID поля (столбца), к которому применяется фильтр (только для листового узла)
    * `data.filters[].operator`: оператор фильтра (см. таблицу операторов ниже; только для листового узла)
    * `data.filters[].value`: значение фильтра, объект `Value` (формы значений см. в разделе [Установить значение ячейки](#установить-значение-ячейки); только для листового узла)
    * `data.filters[].relativeDate`: необязательный дескриптор относительной даты для фильтров по дате (`{ "count": 7, "unit": 0, "direction": -1 }`; `unit`: `0` день, `1` неделя, `2` месяц, `3` год; `direction`: `-1` до, `0` в этот период, `1` после; только для листового узла)
    * `data.filters[].combination`: комбинатор группы, `"and"` или `"or"` (только для группового узла)
    * `data.filters[].filters`: дочерние узлы фильтров, рекурсивно `ViewFilter` (только для группового узла)
    * `data.sorts`: массив `ViewSort`
    * `data.sorts[].column`: ID поля (столбца), к которому применяется сортировка
    * `data.sorts[].order`: `ASC` или `DESC`

  Операторы фильтров:

  | Значение             | Описание             |
  |----------------------|----------------------|
  | `=`                  | Равно                |
  | `!=`                 | Не равно             |
  | `>`                  | Больше               |
  | `>=`                 | Больше или равно     |
  | `<`                  | Меньше               |
  | `<=`                 | Меньше или равно     |
  | `Contains`           | Содержит             |
  | `Does not contains`  | Не содержит          |
  | `Is empty`           | Пусто                |
  | `Is not empty`       | Не пусто             |
  | `Starts with`        | Начинается с         |
  | `Ends with`          | Заканчивается на     |
  | `Is between`         | Между                |
  | `Is true`            | Истина (флажок)      |
  | `Is false`           | Ложь (флажок)        |

### Установить фильтры

* `/api/av/setAttrViewFilters`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t",
    "data": [
      {
        "column": "20240118203822-io6ofxb",
        "operator": "=",
        "value": {
          "type": "select",
          "mSelect": [
            { "content": "Done", "color": "1" }
          ]
        }
      }
    ]
  }
  ```

    * `avID`: ID базы данных
    * `blockID`: блок базы данных, владеющий представлением
    * `data`: полный новый массив объектов `ViewFilter`, который **целиком заменяет** существующие фильтры представления (см. [Получить фильтры и сортировку](#получить-фильтры-и-сортировку)). Передайте `[]`, чтобы очистить все фильтры. Верхний уровень содержит единственный корневой групповой узел `{ "combination": "and"|"or", "filters": [...] }`; элементы массива — либо листовые фильтры, либо вложенные групповые узлы, что позволяет рекурсивно комбинировать AND/OR
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Установить сортировку

* `/api/av/setAttrViewSorts`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "blockID": "20240118120201-kldj15t",
    "data": [
      {
        "column": "20240118120204-w6cggab",
        "order": "DESC"
      }
    ]
  }
  ```

    * `avID`: ID базы данных
    * `blockID`: блок базы данных, владеющий представлением
    * `data`: полный новый массив объектов `ViewSort`, который **целиком заменяет** существующие сортировки представления (см. [Получить фильтры и сортировку](#получить-фильтры-и-сортировку)). Передайте `[]`, чтобы очистить все сортировки
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

    * После успешного вызова проверьте сохранение через [Получить фильтры и сортировку](#получить-фильтры-и-сортировку)

### Добавить поле

Добавляет новое поле (столбец). Поле добавляется в каждое представление (таблица/галерея/канбан) в позицию после `previousKeyID` (или в позицию по умолчанию, если он пуст).

* `/api/av/addAttributeViewKey`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "keyID": "20240118120204-7k9wzbp",
    "keyName": "状态",
    "keyType": "select",
    "keyIcon": "",
    "previousKeyID": "20240118120204-w6cggab"
  }
  ```

    * `avID`: ID базы данных
    * `keyID`: ID нового поля. Должен быть валидным ID узла, сгенерированным `Lute.NewNodeID()` (14-значная метка времени + `-` + 7 случайных буквенно-цифровых символов, например `20240118120204-abc1234`)
    * `keyName`: отображаемое имя поля
    * `keyType`: тип поля — одно из `text`, `number`, `date`, `select`, `mSelect`, `url`, `email`, `phone`, `mAsset`, `template`, `created`, `updated`, `checkbox`, `relation`, `rollup`, `lineNumber`. `block` (первичный ключ) через эту конечную точку добавить нельзя
    * `keyIcon`: необязательный значок поля (эмодзи или пустая строка)
    * `previousKeyID`: вставить новый столбец после поля с этим ID. Пустая строка — позиция по умолчанию для макета (первый столбец для таблицы, последний для галереи/канбана)
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Удалить поле

Удаляет поле (столбец) и все его значения. Возвращает `code: -1` с `msg: "key not found"`, если `keyID` не существует.

* `/api/av/removeAttributeViewKey`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "keyID": "20240118120204-7k9wzbp",
    "removeRelationDest": false
  }
  ```

    * `avID`: ID базы данных
    * `keyID`: ID удаляемого поля
    * `removeRelationDest`: при `true`, если поле является связью, также удаляет соответствующее поле обратной связи из целевой базы данных. По умолчанию `false`
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Задать глобальный порядок полей

Переупорядочивает поле (столбец) глобально — перемещает `keyID` в позицию после `previousKeyID` в порядке полей, затрагивая каждое представление.

* `/api/av/sortAttributeViewKey`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "keyID": "20240118203822-io6ofxb",
    "previousKeyID": "20240118120204-w6cggab"
  }
  ```

    * `avID`: ID базы данных
    * `keyID`: ID перемещаемого поля
    * `previousKeyID`: ID поля, после которого следует разместить `keyID`. Пустая строка перемещает его на первую позицию
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```

### Задать порядок полей в представлении

Переупорядочивает столбец в рамках макета одного представления (например, порядок столбцов таблицы), не меняя глобальный порядок полей.

* `/api/av/sortAttributeViewViewKey`
* Параметры

  ```json
  {
    "avID": "20240118120204-kwyzf77",
    "viewID": "20240118120204-7rnmyc1",
    "keyID": "20240118203822-io6ofxb",
    "previousKeyID": "20240118120204-w6cggab"
  }
  ```

    * `avID`: ID базы данных
    * `viewID`: целевое представление. Если пусто, используется текущее
    * `keyID`: ID перемещаемого поля
    * `previousKeyID`: ID поля, после которого следует разместить `keyID`. Пустая строка перемещает его на первую позицию
* Возвращаемое значение

  ```json
  {
    "code": 0,
    "msg": "",
    "data": null
  }
  ```
