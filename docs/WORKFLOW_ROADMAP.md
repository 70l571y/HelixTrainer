# Workflow Roadmap

Roadmap ниже группирует развитие HelixTrainer по типам реальной работы,
а не просто по списку команд.

## 1. Single-Buffer Editing Fluency

Цель:
научить пользователя быстро редактировать код внутри одного файла через
selection-first мышление.

Примеры:

- movement + change
- textobject cleanup
- surround workflows
- selection split/filter/transform

## 2. Repetition And State Preservation

Цель:
научить пользователя повторять операции и не терять нужное содержимое.

Примеры:

- macros
- repeat insert
- repeat find
- registers
- blackhole / no-yank workflows

## 3. Project Navigation

Цель:
сделать workspace navigation естественной частью editing loop.

Примеры:

- file and buffer pickers
- workspace symbol picker
- diagnostics picker
- jumplist return
- split comparison

## 4. LSP-Driven Refactor

Цель:
закрепить мысль, что Helix editing и LSP navigation — это один и тот же
рабочий цикл.

Примеры:

- definition / declaration / type definition
- references
- rename
- code actions

## 5. Optional Command-Line Knowledge

Цель:
сохранить полезные `:`-scenarios как reference-track, не превращая их в
центр progression.

Примеры:

- `:open`
- `:write path`
- `:update`
- regex replace через command-line
