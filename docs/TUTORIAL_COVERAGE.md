# Tutorial Coverage Map

Этот документ нужен не для полного инвентаря challenge-ов, а для
понимания качества покрытия по классам навыков.

## Dedicated Strong Coverage

- базовое редактирование и `edit_change`
- буферный workflow
- diagnostics navigation
- LSP references / definition
- split/window navigation
- multicursor editing

## Dedicated But Thin Coverage

- `movement_count`
- `movement_jump_label`
- `movement_match`
- `edit_comment`
- `edit_open_line`
- `number_increment`
- `undo_redo`
- `select_cursor`
- `selection_reset`
- `window_split_current`

## Incidental-Only Or Missing

На текущем этапе из ранее известных official workflows явных дыр не
осталось. Дальше задача не в том, чтобы “добавить хоть что-то”, а в том,
чтобы углублять уже покрытые workflow в более реалистичные multi-step
сценарии.

## Coverage Rule

Функцию можно считать покрытой только если есть хотя бы один challenge,
где она является intended primary workflow, а не побочным шагом.
