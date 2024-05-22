package service

import (
	"slices"
	"testing"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
)

func TestNewSeparatorService(t *testing.T) {
	t.Parallel()

	mockRNGService := &mockRNGService{}

	tests := []struct {
		name    string
		cfg     *config.Settings
		wantErr bool
	}{
		{
			name:    "Valid configuration",
			cfg:     &config.Settings{SeparatorCharacter: "*", SeparatorAlphabet: []string{"!", "@", "#", "$", "%"}},
			wantErr: false,
		},
		{
			name:    "Invalid configuration - invalid separator character",
			cfg:     &config.Settings{SeparatorCharacter: "invalid"},
			wantErr: true,
		},
		{
			name:    "Valid configuration - separator alphabet",
			cfg:     &config.Settings{SeparatorCharacter: option.SeparatorCharacterRandom, SeparatorAlphabet: []string{""}},
			wantErr: false,
		},

		{
			name:    "Valid configuration - separator alphabet",
			cfg:     &config.Settings{SeparatorCharacter: option.SeparatorCharacterRandom, SeparatorAlphabet: []string{"a"}},
			wantErr: false,
		},
		{
			name:    "Invalid configuration - empty separator alphabet",
			cfg:     &config.Settings{SeparatorCharacter: option.SeparatorCharacterRandom, SeparatorAlphabet: []string{"aaa"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSeparatorService(tt.cfg, mockRNGService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSeparatorService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeparatorServiceSeparate(t *testing.T) {
	t.Parallel()

	rngs := &mockRNGService{}
	erngs := &mockEvenRNGService{}

	// Define test cases
	tests := []struct {
		name      string
		cfg       *config.Settings
		rngSvc    RNGService
		input     []string
		expected  []string
		expectErr error
	}{
		{
			name:     "With fixed separator",
			cfg:      &config.Settings{SeparatorCharacter: "-"},
			rngSvc:   rngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"-", "a", "-", "b", "-", "c", "-"},
		},
		{
			name:     "With empty slice",
			cfg:      &config.Settings{SeparatorCharacter: "-"},
			rngSvc:   rngs,
			input:    []string{},
			expected: []string{"-"},
		},
		{
			name:     "With random separator",
			cfg:      &config.Settings{SeparatorCharacter: option.SeparatorCharacterRandom, SeparatorAlphabet: []string{"!", "-", "="}},
			rngSvc:   rngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"-", "a", "-", "b", "-", "c", "-"},
		},
		{
			name:     "With random separator with RNG returning a even number",
			cfg:      &config.Settings{SeparatorCharacter: option.SeparatorCharacterRandom, SeparatorAlphabet: []string{"!", "-", "="}},
			rngSvc:   erngs,
			input:    []string{"a", "b", "c"},
			expected: []string{"=", "a", "=", "b", "=", "c", "="},
		},
	}

	// Run test cases
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			runSeparatorServiceSeparateTest(t, tt.cfg, tt.rngSvc, tt.input, tt.expected, tt.expectErr)
		})
	}
}

// Helper function to run the test cases
func runSeparatorServiceSeparateTest(t *testing.T, cfg *config.Settings, rngSvc RNGService, input []string, expected []string, expectErr error) {
	t.Helper()

	service, err := NewSeparatorService(cfg, rngSvc)
	if err != nil {
		t.Errorf("unexpected error with service init: %s", err)
		return
	}

	got, err := service.Separate(input)

	if expectErr != nil {
		if err == nil {
			t.Errorf("expected error but got none")
		} else if err.Error() != expectErr.Error() {
			t.Errorf("expected error %q but got %q", expectErr, err)
		}
		return
	}

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !slices.Equal(got, expected) {
		t.Errorf("expected %v, but got %v", expected, got)
	}
}
