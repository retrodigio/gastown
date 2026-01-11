package cmd

import (
	"strings"
	"testing"
)

func TestGetSubscribersFromDescription(t *testing.T) {
	tests := []struct {
		name     string
		desc     string
		expected []string
	}{
		{
			name:     "new format single subscriber",
			desc:     "Convoy tracking 2 issues\nSubscribers: mayor/",
			expected: []string{"mayor/"},
		},
		{
			name:     "new format multiple subscribers",
			desc:     "Convoy tracking 2 issues\nSubscribers: mayor/, deacon/, human@email.com",
			expected: []string{"mayor/", "deacon/", "human@email.com"},
		},
		{
			name:     "legacy notify format",
			desc:     "Convoy tracking 2 issues\nNotify: mayor/",
			expected: []string{"mayor/"},
		},
		{
			name:     "no subscribers",
			desc:     "Convoy tracking 2 issues\nMolecule: mol-123",
			expected: nil,
		},
		{
			name:     "empty subscribers line",
			desc:     "Convoy tracking 2 issues\nSubscribers: ",
			expected: nil,
		},
		{
			name:     "subscribers with extra whitespace",
			desc:     "Convoy tracking 2 issues\nSubscribers:  mayor/ ,  deacon/  ",
			expected: []string{"mayor/", "deacon/"},
		},
		{
			name:     "subscribers with molecule after",
			desc:     "Convoy tracking 2 issues\nSubscribers: mayor/\nMolecule: mol-123",
			expected: []string{"mayor/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSubscribersFromDescription(tt.desc)
			if len(got) != len(tt.expected) {
				t.Errorf("getSubscribersFromDescription() = %v, want %v", got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("getSubscribersFromDescription()[%d] = %q, want %q", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestUpdateSubscribersInDescription(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		subscribers []string
		wantContain string
	}{
		{
			name:        "add to empty description",
			desc:        "Convoy tracking 2 issues",
			subscribers: []string{"mayor/"},
			wantContain: "Subscribers: mayor/",
		},
		{
			name:        "add multiple subscribers",
			desc:        "Convoy tracking 2 issues",
			subscribers: []string{"mayor/", "deacon/"},
			wantContain: "Subscribers: mayor/, deacon/",
		},
		{
			name:        "replace existing subscribers",
			desc:        "Convoy tracking 2 issues\nSubscribers: old@example.com",
			subscribers: []string{"new@example.com"},
			wantContain: "Subscribers: new@example.com",
		},
		{
			name:        "migrate from legacy notify",
			desc:        "Convoy tracking 2 issues\nNotify: mayor/",
			subscribers: []string{"mayor/", "deacon/"},
			wantContain: "Subscribers: mayor/, deacon/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := updateSubscribersInDescription(tt.desc, tt.subscribers)
			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("updateSubscribersInDescription() = %q, want to contain %q", got, tt.wantContain)
			}
		})
	}
}

func TestUpdateSubscribersInDescription_NoLegacy(t *testing.T) {
	// After migration, Notify: should not appear
	desc := "Convoy tracking 2 issues\nNotify: mayor/"
	subscribers := []string{"mayor/", "deacon/"}

	got := updateSubscribersInDescription(desc, subscribers)

	if strings.Contains(got, "Notify:") {
		t.Errorf("updateSubscribersInDescription() should replace Notify: with Subscribers:, got %q", got)
	}
}
