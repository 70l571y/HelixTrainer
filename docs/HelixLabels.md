# Метки Helix

Этот файл содержит стандартизированные метки, используемые в конфигурациях челленджей для идентификации тестируемых функций Helix.

## Перемещение (Movement)
*   `movement_basic`: Базовое перемещение курсора (h, j, k, l).
*   `movement_word`: Перемещение по словам (w, b, e, W, B, E).
*   `movement_count`: Числовые префиксы перед движениями (2w, 3e).
*   `movement_find`: Поиск символа (f, t, F, T).
*   `movement_goto`: Команды перехода (gg, ge, gh, gl и т.д.).
*   `movement_jump_label`: Переход по двухсимвольным меткам (gw).
*   `movement_match`: Переход к парной скобке (mm).
*   `jumplist`: Переходы по jumplist (Ctrl-s, Ctrl-o, Ctrl-i).

## Редактирование (Editing)
*   `edit_insert`: Вход в режим вставки (i, a, I, A, o, O).
*   `edit_delete`: Удаление текста (d).
*   `edit_change`: Изменение текста (c).
*   `edit_comment`: Переключение комментариев (Ctrl-c).
*   `edit_replace`: Замена символов (r, R).
*   `edit_case`: Смена регистра (~, `).
*   `edit_yank_paste`: Копирование и вставка (y, p, P).
*   `edit_replace_yanked`: Замена выделения содержимым yank/register (R).
*   `edit_join`: Объединение строк (J).
*   `edit_indent`: Отступы (<, >).
*   `edit_open_line`: Открытие новой строки для вставки (o, O).
*   `format_buffer`: Форматирование буфера или выделения (=).
*   `format_selection`: Форматирование только выбранного фрагмента (= на выделении).
*   `number_increment`: Изменение чисел под курсором (Ctrl-a, Ctrl-x).
*   `undo_redo`: Отмена и повтор изменений (u, U).
*   `repeat_insert`: Повтор последней вставки (.).
*   `repeat_find`: Повтор последнего выделения через f/t (Alt-.).

## Выделение (Selection)
*   `select_basic`: Базовый режим выделения (v).
*   `select_line`: Выделение строки (x).
*   `select_regex`: Regex-выделение в текущем выделении (s).
*   `select_split_regex`: Разделение выделения по regex (S).
*   `select_split_lines`: Разделение выделений по строкам (Alt-s).
*   `select_cursor`: Добавление курсоров (C).
*   `select_object`: Выделение текстовых объектов (ma, mi).
*   `select_tree_object`: Tree-sitter textobjects (`mif`, `mit`, `mia` и т.д.).
*   `select_tree_parameter`: Tree-sitter textobjects для параметров (`mip`, `map`).
*   `select_syntax`: Расширение выделения по синтаксическому дереву (Alt-o).
*   `select_all`: Выделение всего файла (%).
*   `selection_reset`: Свернуть или развернуть выделение (;, Alt-;).
*   `selection_align`: Выровнять несколько выделений (&).
*   `selection_cycle`: Переключить основное выделение ((, ), Alt-,).

## Поиск (Search)
*   `search_basic`: Базовый поиск (/, ?).
*   `search_global`: Глобальный поиск по workspace (`<space>/`).
*   `search_selection`: Поиск с использованием текущего выделения (*).
*   `search_next`: Навигация по совпадениям (n, N).
*   `search_select_next`: Добавление выделений через v + n/N.
*   `search_replace_selections`: Замена через select-all/regex/change (`%`, `s`, `c`).
*   `search_subselect`: Повторный `s` внутри уже найденных совпадений для точечной замены.

## LSP и диагностика
*   `lsp_definition`: Переход к определению символа (gd).
*   `lsp_declaration`: Переход к объявлению символа (gD).
*   `lsp_references`: Переход или выбор ссылок на символ (gr, `<space>h`).
*   `lsp_select_references`: Выбор всех references символа как selection (`<space>h`).
*   `lsp_rename`: Переименование символа через LSP (`<space>r`).
*   `lsp_code_action`: Применение code action (`<space>a`).
*   `lsp_type_definition`: Переход к type definition символа (`gy`).
*   `lsp_hover`: Просмотр hover-документации по символу (`<space>k`).
*   `diagnostics_nav`: Навигация по диагностике (`[d`, `]d`).

## Окружение (Surround)
*   `surround_add`: Добавление окружения (ms).
*   `surround_replace`: Замена окружения (mr).
*   `surround_delete`: Удаление окружения (md).

## Продвинутое (Advanced)
*   `multicursor`: Использование нескольких курсоров.
*   `macro`: Запись и воспроизведение макросов (q, Q).
*   `register`: Работа с именованными регистрами (").
*   `register_blackhole`: Удаление без сохранения в регистр (`"_`).
*   `clipboard_system`: Работа с системным clipboard (Space+y, Space+p).
*   `filter`: Фильтрация курсоров (K, Alt-K).
*   `shell_pipe`: Конвейер shell для выделения (|).
*   `shell_transform`: Трансформация выделения через shell-команду.
*   `selection_rotate_contents`: Циклический сдвиг содержимого выделений (Alt-(, Alt-)).

## Окна и буферы (Windows & Buffers)
*   `command_mode`: Командный режим (:, команды редактора).
*   `file_save`: Сохранение и выход через командный режим (:w, :wq, :q).
*   `window_split`: Разделение окон (:hs, :vs).
*   `window_split_current`: Разделение текущего буфера (Ctrl-w v, Ctrl-w s).
*   `window_move`: Переключение между окнами (<space>w h/j/k/l).
*   `window_layout`: Перестановка и трансформация окон (Ctrl-w HJKL, Ctrl-w t).
*   `window_close`: Закрытие текущего или лишних окон (Ctrl-w q, Ctrl-w o).
*   `buffer_open`: Открытие буфера (:o).
*   `buffer_switch`: Переключение буферов (:b, :previous-buffer).

## Picker и навигация
*   `picker_symbols`: Навигация по символам (<space>s).
*   `picker_files`: Поиск файлов (<space>f).
*   `picker_changed_files`: Навигация по изменённым git-файлам (<space>g).
*   `picker_buffer`: Выбор буфера (<space>b).
*   `picker_workspace_symbols`: Поиск символов по workspace (<space>S).
*   `picker_diagnostics`: Переход по диагностике через picker (<space>d, <space>D).
*   `command_palette`: Поиск и запуск команд через palette (<space>?).

## Meta Track
*   `track_core_hotkey`: Задача относится к основному hotkey-first progression track.
*   `track_optional_command_line`: Задача относится к optional command-line track и сознательно тренирует `:`-workflow.
