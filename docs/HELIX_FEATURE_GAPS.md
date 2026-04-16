# Helix Feature Gaps

Этот документ сводит вместе:

- официально задокументированные возможности Helix
- текущее покрытие HelixTrainer
- практические пробелы, которые ещё стоит закрыть

Источники истины:

- [Keymap](https://docs.helix-editor.com/master/keymap.html)
- [Commands](https://docs.helix-editor.com/master/commands.html)
- [Pickers](https://docs.helix-editor.com/master/pickers.html)
- [Surround](https://docs.helix-editor.com/master/surround.html)
- [Helix Labels](./HelixLabels.md)

## Remaining High Priority Gaps

На текущем этапе явных high-priority official gaps не осталось.

## Recently Closed In This Cycle

- `Alt-c` / no-yank change:
  закрыто challenge-ом `change_without_yank_preserve_register`.
- `<space>R` / replace selections with clipboard:
  закрыто challenge-ом `clipboard_replace_selection`.
- picker prompt register insertion (`Ctrl-r` в picker/prompt):
  закрыто challenge-ом `picker_global_search_register_query`.
- второй practical scenario для `window_split_current`:
  закрыто challenge-ом `window_split_current_compare_move`.
- второй practical scenario для `movement_match`:
  закрыто challenge-ом `match_bracket_cleanup_with_noise`.
- дополнительный workspace-aware challenge для `lsp_rename`:
  закрыто challenge-ом `rename_symbol_two_hop`.
- дополнительный dedicated batch для `picker_buffer`:
  закрыто challenge-ами `picker_buffer_dual_fix` и
  `buffer_picker_history_reopen`.
- dedicated challenge для `changed_file_picker`:
  закрыто runtime support-ом `git_dirty_files` и challenge-ом
  `changed_files_picker_fix`.
- dedicated challenge для `command_palette`:
  закрыто challenge-ом `command_palette_reopen_picker`.
- dedicated hover/popup workflow:
  закрыто challenge-ом `popup_docs_then_fix`.
- второй dedicated challenge для `picker_symbols`:
  закрыто challenge-ом `picker_symbols_local_roundtrip`.
- dedicated challenge для picker open in split:
  закрыто challenge-ом `picker_preview_split_compare`.
- dedicated challenge для picker background open:
  закрыто challenge-ом `picker_background_open_then_jump`.
- applied challenge для `select_cursor`:
  закрыто challenge-ом `select_cursor_applied_refactor`.
- applied challenge для `selection_reset`:
  закрыто challenge-ом `selection_reset_edit_cycle`.
- applied challenge для `undo_redo`:
  закрыто challenge-ом `undo_redo_branching_fix`.
- additional applied challenge для `movement_count`:
  закрыто challenge-ом `movement_count_applied_cleanup`.
- additional applied challenge для `movement_jump_label`:
  закрыто challenge-ом `movement_jump_label_followup`.
- additional applied challenge для `edit_comment`:
  закрыто challenge-ом `edit_comment_toggle_block`.
- additional applied challenge для `edit_open_line`:
  закрыто challenge-ом `edit_open_line_steps`.
- additional applied challenge для `number_increment`:
  закрыто challenge-ом `number_increment_batch_tune`.

## Thin Coverage Areas

Текущий набор больше не имеет очевидных thin-coverage функций из
официального backlog этого цикла. Дальнейшее расширение теперь скорее
про variation depth, чем про закрытие дыр.

## Structural Project Gaps

- Реестр тегов должен быть синхронизирован с challenge-config-ами.
- Authoring docs должны описывать `main_file_name`, иначе multi-file
  сценарии проектируются хуже, чем умеет runtime.
- Для command-line track и hotkey-core track нужны явно задокументированные
  meta-tags, чтобы backlog и learning-path оставались согласованными.
- Некоторые будущие challenge-ы всё ещё могут потребовать расширения
  runtime, если понадобятся более сложные project-aware сценарии
  поверх git/LSP/picker state.

## Recommended Next Batch

Если нужен следующий пакет работ после текущего цикла, разумный порядок:

1. deepen existing covered areas through harder multi-step workflows
2. добавить variation depth для shell/filter/macro chains
3. усилить project-scale LSP orchestration scenarios
4. переработать progression ordering на основе новой полноты покрытия
5. при необходимости добавить новые labels только под реально новые official features
