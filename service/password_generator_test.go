package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/eljamo/libpass/v7/config"
)

type mockTransformerService struct{}

func (m *mockTransformerService) Transform(words []string) ([]string, error) {
	return words, nil
}

type mockSeparatorService struct{}

func (m *mockSeparatorService) Separate(words []string) ([]string, error) {
	char := "-"
	separatedSlice := make([]string, 0, len(words))
	for _, element := range words {
		separatedSlice = append(separatedSlice, char, element)
	}
	separatedSlice = append(separatedSlice, char)

	return separatedSlice, nil
}

type mockPaddingService struct{}

func (m *mockPaddingService) Pad(password []string) (string, error) {
	return "!05" + strings.Join(password, "") + "67!", nil
}

type mockWordListService struct{}

func (m *mockWordListService) GetWords() ([]string, error) {
	return []string{"word1", "word2"}, nil
}

type mockTransformerErrService struct{}

func (m *mockTransformerErrService) Transform(words []string) ([]string, error) {
	return nil, errors.New("transformer error")
}

type mockSeparatorErrService struct{}

func (m *mockSeparatorErrService) Separate(words []string) ([]string, error) {
	return nil, errors.New("separator error")
}

type mockPaddingErrService struct{}

func (m *mockPaddingErrService) Pad(password []string) (string, error) {
	return "", errors.New("padding error")
}

type mockWordListErrService struct{}

func (m *mockWordListErrService) GetWords() ([]string, error) {
	return nil, errors.New("word list error")
}

func TestNewCustomPasswordGeneratorService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *config.Settings
		wantErr bool
	}{
		{
			name:    "Valid Configuration",
			config:  &config.Settings{NumPasswords: 1},
			wantErr: false,
		},
		{
			name:    "Invalid Configuration - NumPasswords is 0",
			config:  &config.Settings{NumPasswords: 0},
			wantErr: true,
		},
		{
			name:    "Invalid Configuration - NumPasswords is 11",
			config:  &config.Settings{NumPasswords: 11},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewCustomPasswordGeneratorService(tt.config, &mockTransformerService{}, &mockSeparatorService{}, &mockPaddingService{}, &mockWordListService{})
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCustomPasswordGeneratorService() with config %+v, error = %v, wantErr %v", tt.config, err, tt.wantErr)
			}
		})
	}
}

func TestPasswordGenerate(t *testing.T) {
	t.Parallel()

	setupService := func(
		transformer TransformerService,
		separator SeparatorService,
		padding PaddingService,
		wordList WordListService,
	) *DefaultPasswordGeneratorService {
		cfg := &config.Settings{NumPasswords: 2, NumWords: 2}
		return &DefaultPasswordGeneratorService{
			cfg,
			transformer,
			separator,
			padding,
			wordList,
		}
	}

	tests := []struct {
		name    string
		service *DefaultPasswordGeneratorService
		wantErr bool
	}{
		{"Valid Service", setupService(&mockTransformerService{}, &mockSeparatorService{}, &mockPaddingService{}, &mockWordListService{}), false},
		{"Word List Error", setupService(&mockTransformerService{}, &mockSeparatorService{}, &mockPaddingService{}, &mockWordListErrService{}), true},
		{"Transformer Error", setupService(&mockTransformerErrService{}, &mockSeparatorService{}, &mockPaddingService{}, &mockWordListService{}), true},
		{"Separator Error", setupService(&mockTransformerService{}, &mockSeparatorErrService{}, &mockPaddingService{}, &mockWordListService{}), true},
		{"Padding Error", setupService(&mockTransformerService{}, &mockSeparatorService{}, &mockPaddingErrService{}, &mockWordListService{}), true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			passwords, err := tt.service.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() with %v: error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}

			if !tt.wantErr && len(passwords) != 2 {
				t.Errorf("Generate() with %v: got = %v, want 2 passwords", tt.name, len(passwords))
			}
		})
	}
}
