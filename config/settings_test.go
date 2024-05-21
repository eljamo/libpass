package config

import (
	"encoding/json"
	"testing"

	"github.com/eljamo/libpass/v7/config/option"
	"github.com/google/go-cmp/cmp"
)

func TestDefaultSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *Settings
	}{
		{
			name: "default settings",
			want: &Settings{
				CaseTransform:           option.CaseTransformRandom,
				NumPasswords:            3,
				NumWords:                3,
				PaddingCharacter:        option.PaddingCharacterRandom,
				PaddingCharactersAfter:  2,
				PaddingCharactersBefore: 2,
				PaddingDigitsAfter:      2,
				PaddingDigitsBefore:     2,
				PaddingType:             option.PaddingTypeFixed,
				Preset:                  option.PresetDefault,
				SeparatorAlphabet:       option.DefaultSpecialCharacters,
				SeparatorCharacter:      option.SeparatorCharacterRandom,
				SymbolAlphabet:          option.DefaultSpecialCharacters,
				WordLengthMax:           8,
				WordLengthMin:           4,
				WordList:                option.WordListEN,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DefaultSettings()

			if !cmp.Equal(got, tt.want) {
				t.Errorf("DefaultSettings() = %+v, want %+v\nDiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []map[string]any
		want    *Settings
		wantErr bool
	}{
		{
			name:    "Default settings when no input maps provided",
			input:   nil,
			want:    DefaultSettings(),
			wantErr: false,
		},
		{
			name: "Merge input maps and create settings",
			input: []map[string]any{
				{"num_passwords": 5},
				{"num_words": 4},
			},
			want: &Settings{
				CaseTransform:           option.CaseTransformRandom,
				NumPasswords:            5,
				NumWords:                4,
				PaddingCharacter:        option.PaddingCharacterRandom,
				PaddingCharactersAfter:  2,
				PaddingCharactersBefore: 2,
				PaddingDigitsAfter:      2,
				PaddingDigitsBefore:     2,
				PaddingType:             option.PaddingTypeFixed,
				Preset:                  option.PresetDefault,
				SeparatorAlphabet:       option.DefaultSpecialCharacters,
				SeparatorCharacter:      option.SeparatorCharacterRandom,
				SymbolAlphabet:          option.DefaultSpecialCharacters,
				WordLengthMax:           8,
				WordLengthMin:           4,
				WordList:                option.WordListEN,
			},
			wantErr: false,
		},
		{
			name: "Invalid input maps result in error",
			input: []map[string]any{
				{"num_passwords": "invalid"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := New(tt.input...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []map[string]any
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Merge multiple maps",
			input: []map[string]any{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			want:    map[string]any{"key1": "value1", "key2": "value2"},
			wantErr: false,
		},
		{
			name:    "Empty input maps",
			input:   nil,
			want:    map[string]any{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotBytes, err := mergeMaps(tt.input...)
			if (err != nil) != tt.wantErr {
				t.Errorf("mergeMaps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var got map[string]any
			if err := json.Unmarshal(gotBytes, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("mergeMaps() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestJsonToSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []byte
		want    *Settings
		wantErr bool
	}{
		{
			name: "Valid JSON",
			input: []byte(`{
				"case_transform": "upper",
				"num_passwords": 5
			}`),
			want: &Settings{
				CaseTransform:           "upper",
				NumPasswords:            5,
				NumWords:                0,
				PaddingCharacter:        "",
				PaddingCharactersAfter:  0,
				PaddingCharactersBefore: 0,
				PaddingDigitsAfter:      0,
				PaddingDigitsBefore:     0,
				PaddingType:             "",
				Preset:                  "",
				SeparatorAlphabet:       nil,
				SeparatorCharacter:      "",
				SymbolAlphabet:          nil,
				WordLengthMax:           0,
				WordLengthMin:           0,
				WordList:                "",
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   []byte(`{"num_passwords": "invalid"}`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := &Settings{}
			err := jsonToSettings(got, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonToSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !cmp.Equal(got, tt.want) {
				t.Errorf("jsonToSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapToJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   map[string]any
		want    string
		wantErr bool
	}{
		{
			name:    "Valid map to JSON",
			input:   map[string]any{"key": "value"},
			want:    `{"key":"value"}`,
			wantErr: false,
		},
		{
			name:    "Empty map to JSON",
			input:   map[string]any{},
			want:    `{}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := mapToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("mapToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("mapToJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
