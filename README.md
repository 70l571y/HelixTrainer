# HelixTrainer

Тренируйте навыки работы с редактором Helix на код-челленджах.

## Установка

### Для пользователей (простая установка)
Вы можете установить напрямую из GitHub без клонирования:

```bash
go install github.com/70l571y/HelixTrainer/cmd/hxtrainer@latest
```

### Для разработчиков
Если вы хотите внести вклад или изменить код:

1.  **Клонируйте репозиторий:**
    ```bash
    git clone https://github.com/70l571y/HelixTrainer.git
    cd HelixTrainer
    ```

2.  **Установите зависимости:**
    ```bash
    go mod download
    ```

3.  **Соберите и установите:**
    ```bash
    go install ./cmd/hxtrainer
    ```

## Расположение данных

HelixTrainer хранит прогресс, базу данных и встроенные челленджи в стандартном каталоге конфигурации вашей системы:

*   **Linux**: `~/.config/hxtrainer/` (или `$XDG_CONFIG_HOME/hxtrainer`)
*   **macOS**: `~/Library/Application Support/hxtrainer/`
*   **Windows**: `%APPDATA%\hxtrainer\`

База данных хранится в файле `hxtrainer.db`, а каталог `challenges_data` будет автоматически создан при первом запуске установленного бинарника. Для сброса прогресса используйте `hxtrainer stats reset`.

## Шпаргалка Команд

```bash
# запуск и выбор challenge
hxtrainer play
hxtrainer play hello_world
hxtrainer play --track core --strategy weak-skills
hxtrainer play --difficulty medium --tag lsp_reference

# просмотр challenge-ов
hxtrainer list
hxtrainer list --json
hxtrainer list --track core
hxtrainer list --difficulty easy --tag movement_basic

# статистика и история
hxtrainer stats
hxtrainer stats --json
hxtrainer stats --track optional
hxtrainer stats --difficulty hard --tag lsp_reference
hxtrainer history
hxtrainer history hello_world
hxtrainer history hello_world --json

# очередь практики
hxtrainer queue
hxtrainer queue --strategy progression
hxtrainer queue --strategy weak-skills --track core --limit 10
hxtrainer queue --json

# перенос прогресса
hxtrainer stats export attempts.json
hxtrainer stats import attempts.json
hxtrainer stats import attempts.json --replace

# обслуживание
hxtrainer stats reset
hxtrainer stats reset --yes
hxtrainer upgrade
hxtrainer doctor
hxtrainer doctor --json

# shell completion
hxtrainer completion bash
hxtrainer completion zsh
hxtrainer completion fish
hxtrainer completion powershell

# разработка и packaging
go test ./...
make test
make build
make install
make release-snapshot
```

## Использование

### 1. Игра

Начните сеанс прохождения челленджей. Система интеллектуально подберёт челлендж на основе вашего прогресса.

```bash
hxtrainer play
```

Helix откроется, и вы получите задачу и подсказки.
Когда закончите, просто выйдите с помощью `:wq`, и ваше решение будет проверено.

```
# hello_world
# Task: Исправьте оператор print, чтобы вывести 'Hello, World!'

fmt.Println("Helo, Wolrd!")


# Базовое редактирование:
# 1. Перемещайтесь с 'h', 'j', 'k', 'l'.
# 2. Перейдите к опечатке.
# 3. Замените символ: 'r' -> правильный символ.
# 4. Или измените текст: 'c' -> ввод -> Esc.
```

Вы также можете запустить конкретный челлендж по его ID:

```bash
hxtrainer play <challenge_id>
```
*Пример: `hxtrainer play hello_world`*

Для управляемого progression доступны фильтры и стратегия выбора:

```bash
hxtrainer play --track core --strategy weak-skills
hxtrainer play --difficulty medium --tag lsp_reference
```

Поддерживаемые стратегии:

* `smart` — текущий адаптивный выбор по истории.
* `progression` — первый нерешённый challenge по учебному порядку.
* `weak-skills` — приоритет навыкам, которые ещё слабо покрыты практикой.


### 2. Список челленджей

Посмотреть все доступные челленджи и статус их выполнения.

```bash
hxtrainer list
```

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┓
┃ ID                          ┃ Сложность  ┃ Язык     ┃ Метки                                                               ┃ Статус    ┃
┡━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━━╇━━━━━━━━━━╇━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━┩
│ edit_delete_block           │ Easy       │ go       │ select_object, edit_delete                                          │ Выполнено │
│ edit_join_lines             │ Medium     │ go       │ edit_join, select_line                                              │ Выполнено │
│ hello_world                 │ Easy       │ go       │ edit_insert, edit_delete, movement_basic                            │ Выполнено │
│ movement_long_jump          │ Easy       │ go       │ movement_goto, edit_change                                          │ Выполнено │
│ surround_add_parens         │ Medium     │ go       │ surround_add, select_basic                                          │ Не решено │
└─────────────────────────────┴────────────┴──────────┴─────────────────────────────────────────────────────────────────────┴───────────┘
```

Если нужен машинно-читаемый вывод:

```bash
hxtrainer list --json
```

Также доступны фильтры:

```bash
hxtrainer list --track core
hxtrainer list --difficulty easy --tag movement_basic
```

### 3. Статистика

Проверьте детальную статистику прогресса.

```bash
hxtrainer stats
```

Для интеграций и скриптов доступен JSON-режим:

```bash
hxtrainer stats --json
```

Фильтры работают и для статистики:

```bash
hxtrainer stats --track optional
hxtrainer stats --difficulty hard --tag lsp_reference
```

Экспорт и импорт прогресса:

```bash
hxtrainer stats export attempts.json
hxtrainer stats import attempts.json
hxtrainer stats import attempts.json --replace
```

Это меню покажет ваши лучшие времена и вехи. Каждый уровень имеет трофеи, которые вы получаете за быстрое прохождение: бронза, серебро, золото и автор.

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━┓
┃ Челлендж                    ┃ Статус    ┃ Лучшее    ┃ Веха      ┃ Попыток  ┃
┡━━━━━━━━━━━━━━━━━━━━━━━━━━━━━╇━━━━━━━━━━━╇━━━━━━━━━━━╇━━━━━━━━━━━╇━━━━━━━━━━┩
│ edit_delete_block           │ Выполнено │ 5.67с     │ 🟢 Author │ 2        │
│ edit_join_lines             │ Выполнено │ 2.84с     │ 🟢 Author │ 2        │
│ hello_world                 │ Выполнено │ 12.14с    │ 🟢 Author │ 3        │
│ movement_long_jump          │ Выполнено │ 6.11с     │ 🟢 Author │ 6        │
│ surround_add_parens         │ Не решено │ -         │ -         │ 0        │
└─────────────────────────────┴───────────┴───────────┴───────────┴──────────┘
```

### 4. Очередь практики

Чтобы посмотреть, что тренировать дальше, не открывая Helix сразу:

```bash
hxtrainer queue
hxtrainer queue --strategy progression
hxtrainer queue --strategy weak-skills --track core --limit 10
hxtrainer queue --json
```

### 5. История попыток

Историю можно посмотреть по одному challenge или по всем сразу:

```bash
hxtrainer history
hxtrainer history hello_world
hxtrainer history hello_world --json
```

### 6. Сброс статистики

Если нужно начать заново, сбросьте статистику отдельной командой. Она удаляет прогресс попыток и рекорды, но не удаляет челленджи.

```bash
hxtrainer stats reset
```

Команда попросит подтверждение перед удалением.

Если вы хотите пропустить подтверждение, используйте явный флаг:

```bash
hxtrainer stats reset --yes
```

### 7. Обновление

Поддерживайте HelixTrainer в актуальном состоянии, чтобы получать новые челленджи.
Команда `hxtrainer upgrade` проверяет наличие новой версии через GitHub Releases и подсказывает команду обновления.

```bash
go install github.com/70l571y/HelixTrainer/cmd/hxtrainer@latest
```

### 8. Диагностика окружения

Если `hxtrainer play` не запускается как ожидается, проверьте окружение:

```bash
hxtrainer doctor
```

Команда покажет наличие `hx` и `git`, а также ключевые пути данных HelixTrainer.
Для автоматизации есть JSON-режим:

```bash
hxtrainer doctor --json
```

### 9. Shell completion

Можно сгенерировать completion script для shell:

```bash
hxtrainer completion bash
hxtrainer completion zsh
hxtrainer completion fish
hxtrainer completion powershell
```

## Packaging

Для локальной сборки и проверки:

```bash
make test
make build
make install
```

Для snapshot-релиза через GoReleaser:

```bash
make release-snapshot
```

## Разработка

Для запуска тестов:

```bash
go test ./...
```
