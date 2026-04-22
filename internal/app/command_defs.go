package app

import "github.com/spf13/cobra"

var playCmd = &cobra.Command{
	Use:   "play [challenge_id]",
	Short: "Запустить челлендж",
	Long:  "Запускает сеанок прохождения челленджа. Если указан ID - запускает конкретный челлендж, иначе выбирает умно.",
	Run:   runPlay,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список всех челленджей",
	Long:  "Показывает все доступные челленджи с их статусом выполнения.",
	Run:   runList,
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Показать статистику прогресса",
	Long:  "Показывает детальную статистику: последние попытки, лучшие времена, вехи.",
	Run:   runStats,
}

var statsResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Сбросить статистику",
	Long:  "Удаляет все попытки и рекорды из статистики после подтверждения. Подтвердить можно интерактивно или через --yes без интерактивного подтверждения.",
	Args:  cobra.NoArgs,
	RunE:  runStatsReset,
}

var statsExportCmd = &cobra.Command{
	Use:   "export <path>",
	Short: "Экспортировать статистику",
	Args:  cobra.ExactArgs(1),
	RunE:  runStatsExport,
}

var statsImportCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "Импортировать статистику",
	Args:  cobra.ExactArgs(1),
	RunE:  runStatsImport,
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Обновить HelixTrainer",
	Long:  "Проверяет и устанавливает последнюю версию из GitHub.",
	Run:   runUpgrade,
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Проверить окружение",
	Long:  "Проверяет наличие внешних зависимостей и выводит ключевые пути данных HelixTrainer.",
	Run:   runDoctor,
}

var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Показать рекомендуемую очередь челленджей",
	Long:  "Строит очередь рекомендуемых челленджей с учётом фильтров и стратегии выбора.",
	Run:   runQueue,
}

var historyCmd = &cobra.Command{
	Use:   "history [challenge_id]",
	Short: "Показать историю попыток",
	Long:  "Показывает последние попытки по всем челленджам или по одному challenge ID.",
	Run:   runHistory,
}
