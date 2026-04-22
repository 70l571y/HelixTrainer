package app

type doctorCheck struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Details string `json:"details"`
}

type listJSONEntry struct {
	ID         string   `json:"id"`
	Difficulty string   `json:"difficulty"`
	Language   string   `json:"language"`
	Tags       []string `json:"tags"`
	Status     string   `json:"status"`
	Completed  bool     `json:"completed"`
}

type statsRecentAttemptJSON struct {
	Timestamp       string  `json:"timestamp"`
	ChallengeID     string  `json:"challenge_id"`
	Result          string  `json:"result"`
	DurationSeconds float64 `json:"duration_seconds"`
}

type statsChallengeJSON struct {
	ID              string   `json:"id"`
	Status          string   `json:"status"`
	BestTimeSeconds *float64 `json:"best_time_seconds,omitempty"`
	Milestone       string   `json:"milestone,omitempty"`
	Attempts        int      `json:"attempts"`
}

type statsJSONPayload struct {
	RecentAttempts []statsRecentAttemptJSON `json:"recent_attempts"`
	Challenges     []statsChallengeJSON     `json:"challenges"`
}

type queueJSONEntry struct {
	ID         string   `json:"id"`
	Difficulty string   `json:"difficulty"`
	Tags       []string `json:"tags"`
	Reason     string   `json:"reason"`
}

type historyJSONEntry struct {
	ChallengeID     string  `json:"challenge_id"`
	Timestamp       string  `json:"timestamp"`
	IsCorrect       bool    `json:"is_correct"`
	DurationSeconds float64 `json:"duration_seconds"`
}
