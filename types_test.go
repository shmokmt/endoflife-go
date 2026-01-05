package endoflife

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "valid date",
			input:    `"2024-01-15"`,
			expected: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "null",
			input:    `null`,
			expected: time.Time{},
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: time.Time{},
			wantErr:  false,
		},
		{
			name:     "invalid date format",
			input:    `"15-01-2024"`,
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Date
			err := json.Unmarshal([]byte(tt.input), &d)

			if (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !d.Time.Equal(tt.expected) {
				t.Errorf("Date.UnmarshalJSON() = %v, want %v", d.Time, tt.expected)
			}
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		date     Date
		expected string
	}{
		{
			name:     "valid date",
			date:     Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
			expected: `"2024-01-15"`,
		},
		{
			name:     "zero date",
			date:     Date{},
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.date)
			if err != nil {
				t.Errorf("Date.MarshalJSON() error = %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("Date.MarshalJSON() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		name     string
		date     Date
		expected string
	}{
		{
			name:     "valid date",
			date:     Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
			expected: "2024-01-15",
		},
		{
			name:     "zero date",
			date:     Date{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.date.String()
			if result != tt.expected {
				t.Errorf("Date.String() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestProductRelease_UnmarshalJSON(t *testing.T) {
	input := `{
		"name": "3.12",
		"codename": null,
		"label": "3.12",
		"releaseDate": "2023-10-02",
		"isLts": false,
		"ltsFrom": null,
		"isEoas": false,
		"eoasFrom": null,
		"isEol": false,
		"eolFrom": "2028-10-01",
		"isDiscontinued": false,
		"discontinuedFrom": null,
		"isEoes": false,
		"eoesFrom": null,
		"isMaintained": true,
		"latest": {
			"name": "3.12.8",
			"date": "2024-12-03",
			"link": "https://python.org/downloads/release/python-3128/"
		}
	}`

	var release ProductRelease
	err := json.Unmarshal([]byte(input), &release)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if release.Name != "3.12" {
		t.Errorf("expected name '3.12', got %s", release.Name)
	}
	if release.IsLTS {
		t.Error("expected IsLTS to be false")
	}
	if !release.IsMaintained {
		t.Error("expected IsMaintained to be true")
	}
	if release.Latest == nil {
		t.Fatal("expected Latest to be non-nil")
	}
	if release.Latest.Name != "3.12.8" {
		t.Errorf("expected latest name '3.12.8', got %s", release.Latest.Name)
	}
	if release.EOLFrom == nil {
		t.Fatal("expected EOLFrom to be non-nil")
	}
	if release.EOLFrom.Format("2006-01-02") != "2028-10-01" {
		t.Errorf("expected EOLFrom '2028-10-01', got %s", release.EOLFrom.Format("2006-01-02"))
	}
}
