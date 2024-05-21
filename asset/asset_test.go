package asset

import (
	"embed"
	"errors"
	"testing"

	"github.com/eljamo/libpass/v7/config/option"
	"github.com/google/go-cmp/cmp"
)

// Mock embedded files for testing
//
//go:embed test_data/*
var testFiles embed.FS

func TestKeyToFile(t *testing.T) {
	tests := []struct {
		key      string
		fileType string
		want     string
		wantOk   bool
	}{
		{key: option.PresetAppleID, fileType: option.ConfigKeyPreset, want: "appleid.json", wantOk: true},
		{key: option.WordListEN, fileType: option.ConfigKeyWordList, want: "en.txt", wantOk: true},
		{key: "invalid", fileType: option.ConfigKeyPreset, want: "", wantOk: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.key, func(t *testing.T) {
			t.Parallel()
			got, ok := keyToFile(tt.key, tt.fileType)
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("keyToFile(%q, %q) = %v, %v; want %v, %v", tt.key, tt.fileType, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestLoadJSONFileData(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		readerFunc  func(string) ([]byte, error)
		expected    map[string]any
		expectError bool
	}{
		{
			name:     "Successful Read",
			filePath: "test_data/sample.json",
			readerFunc: func(string) ([]byte, error) {
				return []byte(`{"key": "value"}`), nil
			},
			expected:    map[string]any{"key": "value"},
			expectError: false,
		},
		{
			name:     "Read Error",
			filePath: "test_data/error.json",
			readerFunc: func(string) ([]byte, error) {
				return nil, errors.New("read error")
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:     "Invalid JSON",
			filePath: "test_data/invalid.json",
			readerFunc: func(string) ([]byte, error) {
				return []byte(`{"key": "value"`), nil // Invalid JSON
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := loadJSONFileData(tt.filePath, tt.readerFunc)
			if (err != nil) != tt.expectError {
				t.Errorf("loadJSONFileData() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("loadJSONFileData() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestReadAndFilterWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filePath string
		minLen   int
		maxLen   int
		want     []string
		wantErr  bool
	}{
		{
			name:     "Valid words with a length of 6",
			filePath: "test_data/words.txt",
			minLen:   6,
			maxLen:   6,
			want:     []string{"banana", "cherry"},
			wantErr:  false,
		},
		{
			name:     "No words with length between 10 and 12",
			filePath: "test_data/words.txt",
			minLen:   10,
			maxLen:   12,
			want:     nil,
			wantErr:  false,
		},
		{
			name:     "File does not exist",
			filePath: "test_data/nonexistent.txt",
			minLen:   3,
			maxLen:   6,
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := readAndFilterWords(tt.filePath, tt.minLen, tt.maxLen, testFiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAndFilterWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("readAndFilterWords() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadJSONFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		content     []byte
		expected    map[string]any
		expectError bool
	}{
		{
			name:        "Valid JSON File",
			filePath:    "test_data/valid.json",
			content:     []byte(`{"name": "test"}`),
			expected:    map[string]any{"name": "test"},
			expectError: false,
		},
		{
			name:        "Invalid JSON File",
			filePath:    "test_data/invalid.json",
			content:     []byte(`{"name": "test"`), // Malformed JSON
			expected:    nil,
			expectError: true,
		},
		{
			name:        "File Not Found",
			filePath:    "test_data/not_found.json",
			content:     nil,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Mock os.ReadFile to return predefined content
			readerFunc := func(string) ([]byte, error) {
				if tt.content == nil {
					return nil, errors.New("file not found")
				}
				return tt.content, nil
			}

			got, err := loadJSONFileData(tt.filePath, readerFunc)
			if (err != nil) != tt.expectError {
				t.Errorf("LoadJSONFile() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if diff := cmp.Diff(got, tt.expected); diff != "" {
				t.Errorf("LoadJSONFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetJSONPreset(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Valid JSON Preset",
			key:  option.PresetAppleID,
			want: map[string]any{
				"word_list":                 "EN",
				"num_passwords":             float64(3),
				"num_words":                 float64(3),
				"word_length_min":           float64(5),
				"word_length_max":           float64(7),
				"case_transform":            "RANDOM",
				"separator_character":       "RANDOM",
				"separator_alphabet":        []any{"-", ":", ".", ","},
				"padding_digits_before":     float64(2),
				"padding_digits_after":      float64(2),
				"padding_type":              "FIXED",
				"padding_character":         "RANDOM",
				"symbol_alphabet":           []any{"!", "?", "@", "&"},
				"padding_characters_before": float64(1),
				"padding_characters_after":  float64(1),
			},
			wantErr: false,
		},
		{
			name:    "Invalid Key",
			key:     "invalid_key",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := GetJSONPreset(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJSONPreset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("GetJSONPreset() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
