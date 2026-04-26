package gmaps

import (
	"regexp"
	"strconv"
	"strings"
)

var reviewRemovalCountRegex = regexp.MustCompile(`(?i)\b(\d+)\s*(?:bis|to|-|–)\s*(\d+)\b`)
var reviewRemovalOverCountRegex = regexp.MustCompile(`(?i)(?:über|ueber|mehr als|over|more than)\s*(\d+)`)
var reviewRemovalContextRegex = regexp.MustCompile(`(?i)(?:(?:diffam|defamation|beschwerde(?:n)?|complaints?).*(?:entfernt|removed))|(?:(?:entfernt|removed).*(?:diffam|defamation|beschwerde(?:n)?|complaints?))`)
var reviewRemovalInlineRegex = regexp.MustCompile(`(?i)((?:\d+\s*(?:bis|to|-|–)\s*\d+|(?:über|ueber|mehr als|over|more than)\s*\d+)[\s\S]{0,200}?(?:bewertungen|reviews)[\s\S]{0,250}?(?:beschwerde(?:n)?|complaints?|diffam(?:ierung|ation)?)[\s\S]{0,150}?(?:entfernt|removed))`)

// parseReviewRemovals extracts a min/max range from the defamation-removal notice text.
func parseReviewRemovals(raw string) (int, int) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return 0, 0
	}
	if !reviewRemovalContextRegex.MatchString(normalized) {
		return 0, 0
	}

	matches := reviewRemovalCountRegex.FindStringSubmatch(normalized)
	if len(matches) == 3 {
		minCount, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, 0
		}

		maxCount, err := strconv.Atoi(matches[2])
		if err != nil {
			return 0, 0
		}

		if minCount > maxCount {
			return maxCount, minCount
		}

		return minCount, maxCount
	}

	overMatches := reviewRemovalOverCountRegex.FindStringSubmatch(normalized)
	if len(overMatches) != 2 {
		return 0, 0
	}

	threshold, err := strconv.Atoi(overMatches[1])
	if err != nil {
		return 0, 0
	}

	return threshold + 1, 0
}

func extractReviewRemovalNoticeFromAny(v any) string {
	switch val := v.(type) {
	case string:
		candidate := strings.TrimSpace(val)
		if candidate == "" {
			return ""
		}

		if !reviewRemovalContextRegex.MatchString(candidate) {
			if match := reviewRemovalInlineRegex.FindString(candidate); match != "" {
				return strings.TrimSpace(match)
			}

			return ""
		}

		if reviewRemovalCountRegex.MatchString(candidate) || reviewRemovalOverCountRegex.MatchString(candidate) {
			return candidate
		}

		if match := reviewRemovalInlineRegex.FindString(candidate); match != "" {
			return strings.TrimSpace(match)
		}

		return ""
	case []any:
		for _, item := range val {
			if notice := extractReviewRemovalNoticeFromAny(item); notice != "" {
				return notice
			}
		}
	case map[string]any:
		for _, item := range val {
			if notice := extractReviewRemovalNoticeFromAny(item); notice != "" {
				return notice
			}
		}
	}

	return ""
}
