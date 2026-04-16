# Hotkey Muscle Memory Backlog

Этот backlog prioritizes challenge-ы, которые укрепляют именно
hotkey-first usage, а не знание командного режима как такового.

## Implemented In Current Cycle

- `change_without_yank_preserve_register`
  Difficulty: Medium
  Focus: `Alt-c`, `register`, `edit_replace_yanked`
  Why: modern documented no-yank workflow почти не покрыт.

- `clipboard_replace_selection`
  Difficulty: Easy
  Focus: `clipboard_system`, `edit_replace_yanked`
  Why: `<space>R` documented, но в progression почти отсутствует.

- `picker_global_search_register_query`
  Difficulty: Medium
  Focus: `search_global`, `register`, `picker_files`
  Why: picker/prompt register insertion через `Ctrl-r` documented и
  полезен для реальной работы.

## Implemented In Follow-Up Batch

- `window_split_current_compare_move`
  Difficulty: Medium
  Focus: `window_split_current`, `window_move`, `edit_yank_paste`
  Why: усиливает single-buffer split workflow на реальном reference-copy.

- `match_bracket_cleanup_with_noise`
  Difficulty: Medium
  Focus: `movement_match`, `edit_insert`
  Why: даёт более прикладной anti-cheese сценарий на `mm`.

- `rename_symbol_two_hop`
  Difficulty: Hard
  Focus: `lsp_rename`, `buffer_switch`, `edit_change`
  Why: делает rename более устойчивым workspace loop, а не одиночным приёмом.

- `picker_buffer_dual_fix`
  Difficulty: Medium
  Focus: `picker_buffer`, `buffer_open`, `buffer_switch`
  Why: делает buffer picker primary workflow, а не incidental шагом.

- `buffer_picker_history_reopen`
  Difficulty: Medium
  Focus: `picker_buffer`, `buffer_switch`, `edit_change`
  Why: соединяет buffer picker с history-return в одном коротком loop.

- `changed_files_picker_fix`
  Difficulty: Medium
  Focus: `picker_changed_files`, `buffer_switch`, `edit_change`
  Why: закрывает отдельный git-aware picker workflow через runtime dirty workspace.

- `command_palette_reopen_picker`
  Difficulty: Medium
  Focus: `command_palette`, `picker_files`, `buffer_switch`
  Why: даёт dedicated путь на `<space>?`, а не incidental доступ к командам.

- `popup_docs_then_fix`
  Difficulty: Medium
  Focus: `lsp_hover`, `lsp_definition`, `edit_change`
  Why: добавляет popup-docs workflow через `<space>k`.

- `picker_symbols_local_roundtrip`
  Difficulty: Medium
  Focus: `picker_symbols`, `edit_change`
  Why: усиливает local symbols picker отдельным scenario в одном буфере.

- `picker_preview_split_compare`
  Difficulty: Medium
  Focus: `picker_files`, `window_split`, `window_move`
  Why: закрывает dedicated workflow на picker open-in-split.

- `picker_background_open_then_jump`
  Difficulty: Medium
  Focus: `picker_files`, `picker_buffer`, `buffer_switch`
  Why: закрывает dedicated workflow на background-open из picker.

- `select_cursor_applied_refactor`
  Difficulty: Medium
  Focus: `select_cursor`, `multicursor`, `edit_change`
  Why: переводит навык из isolated-demo в applied refactor.

- `selection_reset_edit_cycle`
  Difficulty: Medium
  Focus: `selection_reset`, `select_split_regex`, `edit_case`
  Why: добавляет второй realistic loop на reset/expand/collapse selections.

- `undo_redo_branching_fix`
  Difficulty: Medium
  Focus: `undo_redo`, `edit_change`
  Why: делает undo/redo частью нормального edit flow, а не единичной демонстрацией.

- `movement_count_applied_cleanup`
  Difficulty: Medium
  Focus: `movement_count`, `movement_goto`, `edit_change`
  Why: добавляет applied scenario с реальным смещением по файлу.

- `movement_jump_label_followup`
  Difficulty: Medium
  Focus: `movement_jump_label`, `edit_change`
  Why: усиливает `gw` вторым сценарием после базового target jump.

- `edit_comment_toggle_block`
  Difficulty: Easy
  Focus: `edit_comment`, `select_line`
  Why: добавляет противоположный direction toggle: comment instead of uncomment.

- `edit_open_line_steps`
  Difficulty: Easy
  Focus: `edit_open_line`, `edit_insert`
  Why: даёт applied сценарий на обрамление существующей строки сверху и снизу.

- `number_increment_batch_tune`
  Difficulty: Easy
  Focus: `number_increment`
  Why: добавляет серию числовых правок вместо одиночного inc/dec.

## P2 Remaining

На текущем этапе явных P2 gaps для official feature coverage не осталось.

## P3

Следующий backlog уже должен быть про усложнение цепочек действий, а не
про закрытие явных одиночных пробелов покрытия.
