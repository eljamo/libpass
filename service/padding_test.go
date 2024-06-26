package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
)

func TestNewPaddingService(t *testing.T) {
	t.Parallel()

	mockRNGService := &mockRNGService{}

	tests := []struct {
		name    string
		cfg     *config.Settings
		wantErr bool
	}{
		{
			name: "Valid configuration",
			cfg: &config.Settings{
				PaddingDigitsBefore: 2, PaddingDigitsAfter: 2, PaddingCharacter: "*",
				SymbolAlphabet: []string{"!", "@", "#", "$", "%"}, PaddingType: option.PaddingTypeFixed,
			},
			wantErr: false,
		},
		{
			name:    "Invalid configuration - negative padding digits before and after",
			cfg:     &config.Settings{PaddingDigitsBefore: -1, PaddingDigitsAfter: -1},
			wantErr: true,
		},
		{
			name:    "Invalid configuration - negative padding character before and after",
			cfg:     &config.Settings{PaddingType: option.PaddingTypeFixed, PaddingCharactersBefore: -1, PaddingCharactersAfter: -1},
			wantErr: true,
		},
		{
			name:    "Invalid configuration - invalid padding type",
			cfg:     &config.Settings{PaddingCharacter: "invalid", PaddingType: option.PaddingTypeFixed},
			wantErr: true,
		},
		{
			name:    "Invalid configuration - empty symbol alphabet",
			cfg:     &config.Settings{PaddingCharacter: option.PaddingCharacterRandom, SymbolAlphabet: []string{}},
			wantErr: true,
		},
		{
			name:    "Valid configuration - symbol alphabet",
			cfg:     &config.Settings{PaddingCharacter: option.PaddingCharacterRandom, SymbolAlphabet: []string{""}},
			wantErr: false,
		},
		{
			name:    "Valid configuration - symbol alphabet",
			cfg:     &config.Settings{PaddingCharacter: option.PaddingCharacterRandom, SymbolAlphabet: []string{"a"}},
			wantErr: false,
		},

		{
			name:    "Invalid configuration - too large symbol alphabet element",
			cfg:     &config.Settings{PaddingCharacter: option.PaddingCharacterRandom, SymbolAlphabet: []string{"aaa"}},
			wantErr: true,
		},
		{
			name:    "Invalid configuration - invalid padding to length",
			cfg:     &config.Settings{PadToLength: -1, PaddingType: option.PaddingTypeAdaptive},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewPaddingService(tt.cfg, mockRNGService)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPaddingService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPad(t *testing.T) {
	t.Parallel()

	cfg := &config.Settings{
		PaddingDigitsBefore:     2,
		PaddingDigitsAfter:      2,
		SeparatorCharacter:      "-",
		PaddingType:             option.PaddingTypeFixed,
		PaddingCharactersBefore: 2,
		PaddingCharactersAfter:  2,
		PadToLength:             20,
		PaddingCharacter:        "*",
		SymbolAlphabet:          []string{"!"},
	}
	rngs := &mockRNGService{}
	s, err := NewPaddingService(cfg, rngs)
	if err != nil {
		t.Errorf("service init error: %v", err)
	}

	t.Run("FixedPaddingWithConfig", func(t *testing.T) {
		t.Parallel()
		input := []string{"-test-"}
		expected := "**11-test-11**"

		result, err := s.Pad(input)
		if err != nil {
			t.Errorf("Pad() error = %v, expectErr %v", err, false)
		}
		if result != expected {
			t.Errorf("Pad() = %v, expected %v", result, expected)
		}
	})
}

func TestDigits(t *testing.T) {
	t.Parallel()

	rngs := &mockRNGService{}
	tests := []struct {
		name     string
		input    []string
		before   int
		after    int
		expected []string
	}{
		{
			name:     "valid padding",
			input:    []string{"a", "b", "c"},
			expected: []string{"1", "1", "a", "b", "c", "1", "1"},
			before:   2,
			after:    2,
		},
		{
			name:     "no padding",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
			before:   0,
			after:    0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Settings{
				PaddingDigitsBefore: tt.before,
				PaddingDigitsAfter:  tt.after,
			}

			s, err := NewPaddingService(cfg, rngs)
			if err != nil {
				t.Errorf("service init error: %v", err)
			}
			got, err := s.digits(tt.input)
			if err != nil {
				t.Errorf("digits() error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("digits() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestGenerateRandomDigits(t *testing.T) {
	t.Parallel()

	cfg := &config.Settings{}
	tests := []struct {
		name     string
		count    int
		number   int
		expected []string
	}{
		{
			name:     "generate 3 digits",
			count:    3,
			number:   1,
			expected: []string{"1", "1", "1"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rngs := &mockRNGService{}
			s, err := NewPaddingService(cfg, rngs)
			if err != nil {
				t.Errorf("service init error: %v", err)
			}

			got, err := s.generateRandomDigits(tt.count)
			if err := checkError(tt.count, err); err != nil {
				t.Errorf("generateRandomDigits() error = %v, expectErr nil", err)
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("generateRandomDigits() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func checkError(count int, err error) error {
	if count < 0 {
		if err == nil {
			return fmt.Errorf("generateRandomDigits() expected error for count %d, got nil", count)
		}
	}
	return nil
}

func TestRemoveEdgeSeparatorCharacter(t *testing.T) {
	t.Parallel()

	cfg := &config.Settings{SeparatorCharacter: "-"}

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no separator at edges",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at start",
			input:    []string{"-", "a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at end",
			input:    []string{"a", "b", "c", "-"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at both ends",
			input:    []string{"-", "a", "b", "c", "-"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "input with only separators",
			input:    []string{"-", "-"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &DefaultPaddingService{cfg: cfg}
			got := s.removeEdgeSeparatorCharacter(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("removeEdgeSeparatorCharacter() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestRemoveRandomEdgeSeparatorCharacter(t *testing.T) {
	t.Parallel()

	cfg := &config.Settings{SeparatorCharacter: "RANDOM", SeparatorAlphabet: []string{"!", "-", "="}}

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no separator at edges",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at start",
			input:    []string{"-", "a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at end",
			input:    []string{"a", "b", "c", "-"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator at both ends",
			input:    []string{"-", "a", "b", "c", "-"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "input with only separators",
			input:    []string{"-", "-"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &DefaultPaddingService{cfg: cfg}
			got := s.removeEdgeSeparatorCharacter(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("removeEdgeSeparatorCharacter() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestSymbols(t *testing.T) {
	t.Parallel()

	rngs := &mockRNGService{}

	tests := []struct {
		name       string
		cfg        *config.Settings
		pw         string
		want       string
		expectErr  bool
		errMessage string
	}{
		{
			name: "fixed padding with specific character",
			cfg: &config.Settings{
				PaddingType:             option.PaddingTypeFixed,
				PaddingCharacter:        "*",
				PaddingCharactersBefore: 2,
				PaddingCharactersAfter:  2,
				PadToLength:             10,
			},
			pw:   "pass",
			want: "**pass**",
		},
		{
			name: "adaptive padding to specific length",
			cfg: &config.Settings{
				PaddingType:      option.PaddingTypeAdaptive,
				PaddingCharacter: "*",
				PadToLength:      10,
			},
			pw:   "pass",
			want: "pass******",
		},
		{
			name: "no padding",
			cfg: &config.Settings{
				PaddingType: option.PaddingTypeNone,
			},
			pw:   "password",
			want: "password",
		},
		{
			name: "random padding character",
			cfg: &config.Settings{
				PaddingType:             option.PaddingTypeFixed,
				PaddingCharacter:        option.PaddingCharacterRandom,
				PaddingCharactersBefore: 2,
				PaddingCharactersAfter:  2,
				PadToLength:             10,
				SymbolAlphabet:          []string{"!", "-", "_"},
			},
			pw:   "pass",
			want: "--pass--",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s, err := NewPaddingService(tt.cfg, rngs)
			if err != nil {
				t.Errorf("service init error: %v", err)
			}

			got, err := s.symbols(tt.pw)
			if (err != nil) != tt.expectErr {
				t.Errorf("symbols() error = %v, expectErr %v", err, tt.expectErr)
			}
			if tt.expectErr && err.Error() != tt.errMessage {
				t.Errorf("symbols() error = %v, expectErr message %v", err, tt.errMessage)
			}
			if !tt.expectErr && got != tt.want {
				t.Errorf("symbols() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFixed(t *testing.T) {
	t.Parallel()

	testNormalPadding(t)
	testNoPadding(t)
	testPaddingBeforeOnly(t)
	testPaddingAfterOnly(t)
}

func testNormalPadding(t *testing.T) {
	t.Helper()

	cfg := &config.Settings{
		PaddingCharactersBefore: 2,
		PaddingCharactersAfter:  3,
	}
	svc := &DefaultPaddingService{cfg: cfg, rngSvc: &mockEvenRNGService{}}

	got, err := svc.fixed("password", "*")
	if (err != nil) != false {
		t.Fatalf("fixed() error = %v, expectErr %v", err, false)
	}
	if got != "**password***" {
		t.Errorf("fixed() got = %v, want %v", got, "**password***")
	}
}

func testNoPadding(t *testing.T) {
	t.Helper()

	cfg := &config.Settings{
		PaddingCharactersBefore: 0,
		PaddingCharactersAfter:  0,
	}
	svc := &DefaultPaddingService{cfg: cfg, rngSvc: &mockEvenRNGService{}}

	got, err := svc.fixed("password", "*")
	if (err != nil) != false {
		t.Fatalf("fixed() error = %v, expectErr %v", err, false)
	}
	if got != "password" {
		t.Errorf("fixed() got = %v, want %v", got, "password")
	}
}

func testPaddingBeforeOnly(t *testing.T) {
	t.Helper()

	cfg := &config.Settings{
		PaddingCharactersBefore: 4,
		PaddingCharactersAfter:  0,
	}
	svc := &DefaultPaddingService{cfg: cfg, rngSvc: &mockEvenRNGService{}}

	got, err := svc.fixed("password", "#")
	if (err != nil) != false {
		t.Fatalf("fixed() error = %v, expectErr %v", err, false)
	}
	if got != "####password" {
		t.Errorf("fixed() got = %v, want %v", got, "####password")
	}
}

func testPaddingAfterOnly(t *testing.T) {
	t.Helper()

	cfg := &config.Settings{
		PaddingCharactersBefore: 0,
		PaddingCharactersAfter:  2,
	}
	svc := &DefaultPaddingService{cfg: cfg, rngSvc: &mockEvenRNGService{}}

	got, err := svc.fixed("password", "+")
	if (err != nil) != false {
		t.Fatalf("fixed() error = %v, expectErr %v", err, false)
	}
	if got != "password++" {
		t.Errorf("fixed() got = %v, want %v", got, "password++")
	}
}

func TestAdaptive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		pw     string
		char   string
		padLen int
		want   string
	}{
		{
			name:   "no padding needed",
			pw:     "1234567890",
			char:   "*",
			padLen: 10,
			want:   "1234567890",
		},
		{
			name:   "padding needed",
			pw:     "12345",
			char:   "*",
			padLen: 10,
			want:   "12345*****",
		},
		{
			name:   "empty password",
			pw:     "",
			char:   "*",
			padLen: 10,
			want:   "**********",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := &config.Settings{
				PadToLength: tt.padLen,
			}

			s, err := NewPaddingService(cfg, &mockEvenRNGService{})
			if err != nil {
				t.Errorf("service init error: %v", err)
			}

			got := s.adaptive(tt.pw, tt.char)

			if got != tt.want {
				t.Errorf("adaptive() got = %v, want %v", got, tt.want)
			}
		})
	}
}
