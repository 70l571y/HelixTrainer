package challenges

import (
	"testing"
	"time"

	"github.com/70l571y/HelixTrainer/internal/database"
)

func TestFilterChallenges(t *testing.T) {
	challengeList := []Challenge{
		{ID: "core_easy", Difficulty: "Easy", Tags: []string{"movement_basic", "track_core_hotkey"}},
		{ID: "optional_easy", Difficulty: "Easy", Tags: []string{"command_mode", "track_optional_command_line"}},
		{ID: "core_medium", Difficulty: "Medium", Tags: []string{"movement_basic", "track_core_hotkey"}},
	}

	filtered := FilterChallenges(challengeList, ChallengeFilter{
		Difficulty: "easy",
		Track:      "core",
		Tags:       []string{"movement_basic"},
	})

	if len(filtered) != 1 {
		t.Fatalf("FilterChallenges() len = %d, want 1", len(filtered))
	}
	if filtered[0].ID != "core_easy" {
		t.Fatalf("FilterChallenges()[0].ID = %q, want %q", filtered[0].ID, "core_easy")
	}
}

func TestSelectWeakestChallengePrefersWeakTag(t *testing.T) {
	now := time.Now()
	challengeList := []Challenge{
		{ID: "movement_one", Difficulty: "Easy", Tags: []string{"movement_basic", "track_core_hotkey"}},
		{ID: "movement_two", Difficulty: "Medium", Tags: []string{"movement_basic", "track_core_hotkey"}},
		{ID: "lsp_one", Difficulty: "Easy", Tags: []string{"lsp_reference", "track_core_hotkey"}},
	}

	attempts := []database.Attempt{
		{ChallengeID: "movement_one", IsCorrect: true, Timestamp: now.Add(-3 * time.Minute), Duration: 3},
		{ChallengeID: "movement_two", IsCorrect: false, Timestamp: now.Add(-2 * time.Minute), Duration: 10},
	}

	got := SelectWeakestChallenge(challengeList, attempts)
	if got.ID != "lsp_one" {
		t.Fatalf("SelectWeakestChallenge() = %q, want %q", got.ID, "lsp_one")
	}
}

func TestSelectProgressionChallengePrefersFirstUnsolved(t *testing.T) {
	challengeList := []Challenge{
		{ID: "easy_done", Difficulty: "Easy"},
		{ID: "easy_next", Difficulty: "Easy"},
		{ID: "medium_later", Difficulty: "Medium"},
	}

	attempts := []database.Attempt{
		{ChallengeID: "easy_done", IsCorrect: true},
	}

	got := SelectProgressionChallenge(challengeList, attempts)
	if got.ID != "easy_next" {
		t.Fatalf("SelectProgressionChallenge() = %q, want %q", got.ID, "easy_next")
	}
}
