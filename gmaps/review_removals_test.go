package gmaps

import "testing"

func TestParseReviewRemovals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		minCount int
		maxCount int
	}{
		{
			name:     "german notice",
			input:    "201 bis 250 Bewertungen aufgrund von Beschwerden wegen Diffamierung entfernt.",
			minCount: 201,
			maxCount: 250,
		},
		{
			name:     "english notice",
			input:    "201 to 250 reviews removed due to defamation complaints.",
			minCount: 201,
			maxCount: 250,
		},
		{
			name:     "dash separator",
			input:    "201-250 reviews removed due to complaints.",
			minCount: 201,
			maxCount: 250,
		},
		{
			name:     "invalid text",
			input:    "No removals listed for this business.",
			minCount: 0,
			maxCount: 0,
		},
		{
			name:     "price range should not match",
			input:    "Kebabimbiss · 1-10 €",
			minCount: 0,
			maxCount: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			minCount, maxCount := parseReviewRemovals(tt.input)
			if minCount != tt.minCount || maxCount != tt.maxCount {
				t.Fatalf("unexpected parsed range for %q: got %d-%d, want %d-%d", tt.input, minCount, maxCount, tt.minCount, tt.maxCount)
			}
		})
	}
}

func TestExtractReviewRemovalNoticeFromAny(t *testing.T) {
	t.Parallel()

	data := []any{
		"foo",
		[]any{
			"bar",
			map[string]any{
				"x": "201 bis 250 Bewertungen aufgrund von Beschwerden wegen Diffamierung entfernt.",
			},
		},
	}

	notice := extractReviewRemovalNoticeFromAny(data)
	if notice == "" {
		t.Fatal("expected notice to be extracted")
	}

	minCount, maxCount := parseReviewRemovals(notice)
	if minCount != 201 || maxCount != 250 {
		t.Fatalf("unexpected parsed range: got %d-%d, want 201-250", minCount, maxCount)
	}
}
