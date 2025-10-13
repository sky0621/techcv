package uuidv7

import (
	"regexp"
	"testing"
)

func TestNewString(t *testing.T) {
	pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	ids := make(map[string]struct{})
	for i := 0; i < 10; i++ {
		id, err := NewString()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !pattern.MatchString(id) {
			t.Fatalf("id does not match v7 pattern: %s", id)
		}

		if _, exists := ids[id]; exists {
			t.Fatalf("duplicate id generated: %s", id)
		}
		ids[id] = struct{}{}
	}
}
