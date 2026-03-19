# Helix Learning Path

Этот документ описывает рекомендуемый путь освоения HelixTrainer для
пользователя, который хочет не просто знать отдельные команды Helix, а
уверенно применять редактор в реальной разработке.

Он дополняет:

- [Tutorial Coverage Map](./TUTORIAL_COVERAGE.md)
- [Workflow Roadmap](./WORKFLOW_ROADMAP.md)
- [Helix Feature Gaps](./HELIX_FEATURE_GAPS.md)
- [Hotkey Muscle Memory Backlog](./HOTKEY_MUSCLE_MEMORY_BACKLOG.md)

## Goal

Цель HelixTrainer не в том, чтобы пользователь выучил максимальное число
команд. Цель в другом:

- начать мыслить через selections, а не через “переместился и печатаю”
- видеть в Helix инструмент для реальной работы по коду, а не tutor toy
- собирать быстрые цепочки действий из movement, selection, search,
  tree-sitter, LSP, окон и буферов

Дополнительный принцип для этого проекта:

- приоритет у motor memory для hotkey-driven workflow
- command-line mode полезен, но не должен быть центром progression
- core track должен учить минимизировать уход в `:` и чаще решать задачу
  через movement, pickers, LSP, windows, buffers, repeat и selections

Текущее состояние проекта:

- core hotkey track можно считать завершённым для `v1`
- optional command-line track остаётся как справочный и дополнительный
- дальнейшая работа должна быть направлена на refinement, а не на поиск
  крупных пробелов в core usage

## Core Mental Models

### 1. Movement Is Selection

В Helix движение почти всегда одновременно является и способом
выделения. Это главный сдвиг мышления для новичка.

Что нужно усвоить:

- `w`, `e`, `b`, `f`, `t`, `gg`, `ge` это не только навигация
- `d`, `c`, `y`, `r`, `R`, `J`, `<`, `>` работают по текущему выделению
- хорошее редактирование начинается не с “как удалить”, а с “как точно
  выделить нужную область”

Подходящие challenge-ы:

- `hello_world`
- `movement_find_fields`
- `movement_find_args`
- `edit_change_object`

### 2. Select First, Then Transform

Сильная сторона Helix раскрывается, когда пользователь сначала сужает
область редактирования, а затем применяет трансформацию.

Что нужно усвоить:

- `v`, `x`, `s`, `S`, `C`, `Alt-s` нужны для отбора целей правки
- после narrowing любая операция становится безопаснее и быстрее
- поиск и замена без сужения области часто хуже, чем несколько точных
  выделений

Подходящие challenge-ы:

- `search_select`
- `search_select_next`
- `filter_keep`
- `select_split_comma`
- `selection_cycle_duplicates`

### 3. Structure Beats Text Hacks

На коде Helix особенно силён тогда, когда пользователь опирается на
структуру AST и text objects, а не на удачный regex.

Что нужно усвоить:

- `mi`, `ma`, `mif`, `mip`, `Alt-o` полезнее случайных движений
- tree-sitter-команды делают правки устойчивыми к шуму
- structural edit лучше переносится на реальный код, чем “найти скобку”

Подходящие challenge-ы:

- `text_objects_advanced`
- `select_tree_object_function`
- `select_tree_parameter_drop`
- `select_syntax_tree`
- `select_syntax_expand_call`

### 4. Workspace Fluency Matters

Глубокое владение Helix начинается не внутри одного файла, а при работе
по проекту.

Что нужно усвоить:

- искать символ, диагностику или файл быстрее, чем скроллить вручную
- делать jump -> edit -> jump back как единый цикл
- держать несколько буферов и split-окон под задачу, а не “на всякий
  случай”

Подходящие challenge-ы:

- `picker_files`
- `picker_files_split`
- `picker_buffer`
- `jumplist_return`
- `window_close_others_focus`
- `buffer_switch_history`

### 5. LSP Is Part Of Editing

В практической работе Helix раскрывается через LSP и diagnostics так же
сильно, как через базовую modal editing механику.

Что нужно усвоить:

- `gd`, `gr`, `<space>r`, `<space>a`, `<space>d`, `<space>D`
- navigation, rename, code action и diagnostics это один workflow-класс
- multi-file refactor должен ощущаться как редакторный сценарий, а не
  как переключение в IDE

Подходящие challenge-ы:

- `lsp_definition_helper`
- `lsp_references_local`
- `lsp_references_crossfile`
- `lsp_rename_service`
- `lsp_code_action_import`
- `lsp_code_action_remove_import`
- `diagnostics_nav_fix`
- `picker_diagnostics_workspace`

## Tracks

Для HelixTrainer полезно различать два разных трека.

### Core Hotkey Track

Это главный трек проекта. Он должен строить мышечную память вокруг
быстрых сочетаний клавиш, а не вокруг ввода typable commands.

В core входят:

- movement and goto
- search, repeat and jumplist
- selections, multicursor, split and filter
- text objects and tree-sitter selections
- pickers
- LSP hotkeys
- windows and buffers, если workflow решается хоткеями
- registers and clipboard
- shell and formatting
- macros

### Optional Command-Line Track

Этот трек полезен как дополнительный, но не должен быть основным
маршрутом пользователя, если цель проекта - скорость работы в редакторе
при разработке.

Сюда относятся сценарии, где главная ценность challenge-а в наборе
команды через `:`, а не в hotkey workflow.

## Current Command-Line Track

Ниже задачи, которые сейчас опираются на command-line mode как на
основной приём и поэтому лучше считать optional-треком:

- `command_open_file`
- `buffer_open_command_cycle`
- `buffer_copy_paste`
- `search_replace_regex`
- `open_quoted_path`
- `write_to_new_path`
- `update_only_if_modified`

Пограничные, но всё же core-задачи:

- `lsp_code_action_import`
- `lsp_code_action_remove_import`
- `buffer_previous_next`
- `buffer_switch_history`
- `picker_buffer`
- `windows_split`
- `matrix_split`

У code action задач основной путь решения проходит через `<space>a`, а
не через `:`.

У buffer/window задач выше hotkey-first формулировки уже переписаны и
их стоит считать core-задачами, а не command-line-зависимыми.

## Remaining Optional Rewrites

На текущем этапе high-priority rewrites уже выполнены. Ниже остаются
только задачи, которые сознательно можно оставить optional либо
переписать позже, если нужен ещё более чистый core-track.

### Optional-Track Candidates

- `command_open_file`
  Понизить до optional reference challenge или заменить на picker-based
  open flow.

- `buffer_open_command_cycle`
  Заменить на file picker + previous buffer cycle.

- `buffer_copy_paste`
  Переписать под picker/files workflow без `:o`.

- `search_replace_regex`
  Оставить как optional command-line challenge, а в core делать упор на
  selection-driven replace и search workflows.

### Keep Optional

- `open_quoted_path`
- `write_to_new_path`
- `update_only_if_modified`

Эти задачи полезны как знание о возможностях Helix, но слабо совпадают
с основной целью проекта: развить быструю моторику разработки без
частого ухода в `:`.

## Recommended Track

Ниже рекомендуемый порядок прохождения для пользователя, который хочет
получить именно рабочее владение Helix.

### Stage 1: Basic Modal Muscle Memory

Цель: перестать тянуться к стрелкам и вставке “по месту”.

Challenge-ы:

1. `hello_world`
2. `open_line_answers`
3. `movement_long_jump`
4. `undo_redo_fixup`
5. `goto_line`

Критерий перехода дальше:

- пользователь уверенно двигается без стрелок
- понимает разницу между normal и insert
- не боится отменять и повторять правки

### Stage 2: Selection-Driven Editing

Цель: научиться сначала выделять, потом менять.

Challenge-ы:

1. `movement_find_fields`
2. `edit_change_object`
3. `search_select`
4. `search_select_next`
5. `select_split_comma`
6. `selection_reset_sentences`

Критерий перехода дальше:

- пользователь регулярно использует `v`, `x`, `s`, `S`
- умеет narrowing области вместо глобальной замены

### Stage 3: Productive Text Transformations

Цель: собрать быстрые локальные правки из repeat, registers и multicursor.

Challenge-ы:

1. `repeat_insert_suffix`
2. `replace_with_yank`
3. `registers_named`
4. `register_blackhole_swap`
5. `select_cursor_prefix`
6. `filter_keep`
7. `shell_transform_prefix`

Критерий перехода дальше:

- пользователь не теряет нужный yank в длинной правке
- использует multicursor и filters осознанно, а не случайно

### Stage 4: Structural Code Editing

Цель: перейти от текстовых хаков к syntax-aware редактированию.

Challenge-ы:

1. `text_objects_advanced`
2. `select_tree_object_function`
3. `select_tree_parameter_drop`
4. `select_syntax_tree`
5. `select_syntax_expand_call`
6. `surround_replace_brackets`

Критерий перехода дальше:

- пользователь предпочитает text objects и syntax expansion
- умеет делать правки по форме кода, а не по совпадению текста

### Stage 5: Workspace Navigation

Цель: уверенно перемещаться по проекту и держать контекст задачи.

Challenge-ы:

1. `picker_files`
2. `picker_buffer`
3. `picker_files_split`
4. `jumplist_return`
5. `buffer_switch_history`
6. `window_split_current`
7. `window_close_others_focus`

Критерий перехода дальше:

- пользователь быстро открывает нужный файл или буфер
- умеет возвращаться к исходной точке без потери контекста

### Stage 6: LSP And Diagnostics Workflows

Цель: сделать Helix полноценным инструментом работы по кодовой базе.

Challenge-ы:

1. `lsp_definition_helper`
2. `lsp_references_local`
3. `lsp_references_crossfile`
4. `lsp_rename_service`
5. `lsp_code_action_import`
6. `lsp_code_action_remove_import`
7. `diagnostics_nav_fix`
8. `picker_diagnostics_workspace`
9. `picker_workspace_symbols_jump`

Критерий перехода дальше:

- пользователь воспринимает LSP как часть normal workflow
- умеет делать multi-file fix без ручного перебора файлов

### Stage 7: Orchestration And Real Refactors

Цель: соединять несколько классов команд в один рабочий цикл.

Challenge-ы:

1. `matrix_split`
2. `search_global_todo`
3. `picker_diagnostics_triple`
4. `buffer_copy_paste`
5. `function_to_class`

Критерий освоения:

- пользователь уверенно сочетает search, split, buffer history,
  structure-aware edit и formatting
- прохождение challenge-а сводится к плану действий, а не к поиску одной
  “магической команды”

## How To Read Coverage

Полезно различать три вопроса:

1. Есть ли в проекте challenge на эту команду?
2. Есть ли challenge, где эта команда является основным приёмом?
3. Есть ли challenge, где эта команда используется внутри реального
   multi-step workflow?

[Tutorial Coverage Map](./TUTORIAL_COVERAGE.md)
в основном отвечает на первый вопрос. Этот документ отвечает на второй и
третий.

## What To Add Next

Если цель именно углубить понимание Helix, а не просто расширить покрытие,
следующие challenge-ы дадут лучший эффект:

1. package-level `lsp_rename` на 3-4 файла
2. diagnostics triage с `jump back` после исправления
3. multi-cursor + filter chain на полуреальных данных
4. global search -> split -> compare -> local structural edit
5. tree-sitter refactor с сигнатурой и несколькими call sites

## Authoring Implication

Новый challenge стоит считать учебно ценным только если он даёт ответ на
вопрос: какой реальный рабочий навык Helix станет у пользователя лучше
после прохождения?

Если ответ звучит как “узнает ещё одну кнопку”, этого недостаточно.
