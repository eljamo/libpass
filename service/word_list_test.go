package service

import (
	"testing"

	"github.com/eljamo/libpass/v7/config"
)

func TestNewWordListService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *config.Settings
		wantErr bool
	}{
		{
			name:    "Valid Config",
			cfg:     &config.Settings{NumWords: 5, WordList: "EN_SMALL", WordLengthMin: 2, WordLengthMax: 10},
			wantErr: false,
		},
		{
			name:    "Invalid Config - WordLength",
			cfg:     &config.Settings{NumWords: 5, WordList: "EN_SMALL", WordLengthMin: 15, WordLengthMax: 20},
			wantErr: true,
		},
		{
			name:    "Invalid Config - WordLengthMax",
			cfg:     &config.Settings{NumWords: 5, WordList: "EN_SMALL", WordLengthMin: 10, WordLengthMax: 2},
			wantErr: true,
		},
		{
			name:    "Invalid Config - NumWords",
			cfg:     &config.Settings{NumWords: 1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewWordListService(tt.cfg, &MockRNGService{})
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: NewWordListService() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestDefaultWordListServiceGetWords(t *testing.T) {
	t.Parallel()

	cfg := &config.Settings{NumWords: 5, WordList: "EN_SMALL", WordLengthMin: 2, WordLengthMax: 10}
	service, err := NewWordListService(cfg, &MockRNGService{})
	if err != nil {
		t.Fatalf("Failed to create WordListService: %v", err)
	}

	t.Run("GetWords", func(t *testing.T) {
		words, err := service.GetWords()
		if err != nil {
			t.Errorf("GetWords() error = %v", err)
		}
		if len(words) != cfg.NumWords {
			t.Errorf("GetWords() returned %d words, want %d", len(words), cfg.NumWords)
		}
	})
}
