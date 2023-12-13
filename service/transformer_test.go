package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/eljamo/libpass/v5/config"
	"github.com/eljamo/libpass/v5/config/option"
)

func TestNewTransformerService(t *testing.T) {
	t.Parallel()

	mockRNGService := &MockEvenRNGService{}

	validTransformType := option.Upper
	invalidTransformType := "invalid"

	tests := []struct {
		name          string
		caseTransform string
		wantErr       bool
	}{
		{
			name:          "Valid configuration",
			caseTransform: validTransformType,
			wantErr:       false,
		},
		{
			name:          "Invalid configuration",
			caseTransform: invalidTransformType,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := &config.Settings{CaseTransform: tt.caseTransform}
			_, err := NewTransformerService(cfg, mockRNGService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransformerService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultTransformerService_Transform(t *testing.T) {
	t.Parallel()

	rngs := &MockRNGService{}
	erngs := &MockEvenRNGService{}

	tests := []struct {
		name      string
		cfg       *config.Settings
		rngSvc    RNGService
		input     []string
		expected  []string
		expectErr bool
	}{
		{
			name:     "Alternate",
			cfg:      &config.Settings{CaseTransform: option.Alternate},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hello", "WORLD"},
		},
		{
			name:     "Alternate Lettercase",
			cfg:      &config.Settings{CaseTransform: option.AlternateLettercase},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hElLo", "wOrLd"},
		},
		{
			name:     "Capitalisation",
			cfg:      &config.Settings{CaseTransform: option.Capitalise},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"Hello", "World"},
		},
		{
			name:     "Capitalisation Inversed",
			cfg:      &config.Settings{CaseTransform: option.CapitaliseInvert},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hELLO", "wORLD"},
		},
		{
			name:     "Inversion",
			cfg:      &config.Settings{CaseTransform: option.Invert},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hELLO", "wORLD"},
		},
		{
			name:     "Lower",
			cfg:      &config.Settings{CaseTransform: option.Lower},
			rngSvc:   rngs,
			input:    []string{"HELLO", "WORLD"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Lower Vowel Upper Consonant",
			cfg:      &config.Settings{CaseTransform: option.LowerVowelUpperConsonant},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"HeLLo", "WoRLD"},
		},
		{
			name:     "Sentence",
			cfg:      &config.Settings{CaseTransform: option.Sentence},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"Hello", "world"},
		},
		{
			name:     "Upper",
			cfg:      &config.Settings{CaseTransform: option.Upper},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"HELLO", "WORLD"},
		},
		{
			name:     "Random",
			cfg:      &config.Settings{CaseTransform: option.Random},
			rngSvc:   rngs,
			input:    []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Random with even RNG",
			cfg:      &config.Settings{CaseTransform: option.Random},
			rngSvc:   erngs,
			input:    []string{"hello", "world"},
			expected: []string{"HELLO", "WORLD"},
		},
		{
			name:     "Empty slice",
			cfg:      &config.Settings{CaseTransform: option.Random},
			rngSvc:   rngs,
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "Special characters slice",
			cfg:      &config.Settings{CaseTransform: option.Random},
			rngSvc:   rngs,
			input:    []string{"-", "&"},
			expected: []string{"-", "&"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, err := NewTransformerService(tt.cfg, tt.rngSvc)
			if err != nil {
				t.Errorf("unexpected error with service init: %s", err)
			}

			got, err := service.Transform(tt.input)

			if (err != nil) != tt.expectErr {
				t.Errorf("DefaultTransformerService.Transform() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DefaultTransformerService.Transform() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDefaultTransformerService_Validate(t *testing.T) {
	t.Parallel()

	validCaseTransforms := []string{
		option.Alternate,
		option.AlternateLettercase,
		option.Capitalise,
		option.CapitaliseInvert,
		option.Invert,
		option.Lower,
		option.LowerVowelUpperConsonant,
		option.Random,
		option.Upper,
		option.None,
	}

	for _, validTransform := range validCaseTransforms {
		validTransform := validTransform
		t.Run(fmt.Sprintf("Valid case transform: %s", validTransform), func(t *testing.T) {
			t.Parallel()

			cfg := DefaultTransformerService{
				cfg: &config.Settings{CaseTransform: validTransform},
			}
			if err := cfg.validate(); err != nil {
				t.Errorf("validate() with valid case transform %s returned error: %v", validTransform, err)
			}
		})
	}

	t.Run("Invalid case transform", func(t *testing.T) {
		t.Parallel()

		cfg := DefaultTransformerService{
			cfg: &config.Settings{CaseTransform: "invalid_case_transform"},
		}
		if err := cfg.validate(); err == nil {
			t.Error("validate() with invalid case transform did not return an error")
		}
	})
}
