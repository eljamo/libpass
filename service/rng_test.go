package service

import (
	"errors"
	"fmt"
	"testing"
)

type mockRNGService struct{}

func (s *mockRNGService) GenerateWithMax(max int) (int, error) {
	return 1, nil
}

func (s *mockRNGService) Generate() (int, error) {
	return 1, nil
}

func (s *mockRNGService) GenerateDigit() (int, error) {
	return 1, nil
}

func (s *mockRNGService) GenerateSlice(length int) ([]int, error) {
	return s.GenerateSliceWithMax(length, 2)
}

func (s *mockRNGService) GenerateSliceWithMax(length, max int) ([]int, error) {
	slice := make([]int, length)
	for i := 0; i < length; i++ {
		slice[i] = 2
	}

	return slice, nil
}

type mockEvenRNGService struct{}

func (s *mockEvenRNGService) GenerateWithMax(max int) (int, error) {
	return 2, nil
}

func (s *mockEvenRNGService) Generate() (int, error) {
	return 2, nil
}

func (s *mockEvenRNGService) GenerateDigit() (int, error) {
	return 2, nil
}

func (s *mockEvenRNGService) GenerateSlice(length int) ([]int, error) {
	return s.GenerateSliceWithMax(length, 2)
}

func (s *mockEvenRNGService) GenerateSliceWithMax(length, max int) ([]int, error) {
	slice := make([]int, length)
	for i := 0; i < length; i++ {
		slice[i] = 2
	}

	return slice, nil
}

var errMockRNGService = errors.New("mock RNG Service Error")

type mockErrRNGService struct{}

func (s *mockErrRNGService) GenerateWithMax(max int) (int, error) {
	return 0, errMockRNGService
}

func (s *mockErrRNGService) Generate() (int, error) {
	return 0, errMockRNGService
}

func (s *mockErrRNGService) GenerateDigit() (int, error) {
	return 0, errMockRNGService
}

func (s *mockErrRNGService) GenerateSlice(length int) ([]int, error) {
	return nil, errMockRNGService
}

func (s *mockErrRNGService) GenerateSliceWithMax(length, max int) ([]int, error) {
	return nil, errMockRNGService
}

func TestRNGGenerateWithMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		max       int
		expectErr bool
	}{
		{"ValidMax", 100, false},
		{"NegativeMax", -1, true},
	}

	rngSvc := NewRNGService()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, createTestFunc(tc, rngSvc))
	}
}

func createTestFunc(tc struct {
	name      string
	max       int
	expectErr bool
}, rngSvc *DefaultRNGService,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		generated, err := rngSvc.GenerateWithMax(tc.max)

		if tc.expectErr {
			if err == nil {
				t.Errorf("Expected an error for max = %v, but got none", tc.max)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for max = %v: %v", tc.max, err)
			}
			if generated < 0 || generated >= tc.max {
				t.Errorf("Generated number is out of bounds for max = %v: got %v", tc.max, generated)
			}
		}
	}
}

func TestRNGGenerate(t *testing.T) {
	t.Parallel()

	rngSvc := NewRNGService()

	t.Run("Generate", func(t *testing.T) {
		t.Parallel()
		generated, err := rngSvc.Generate()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if generated < 0 {
			t.Errorf("Generated number is negative: got %v", generated)
		}
	})
}

func TestGenerateSliceWithMax(t *testing.T) {
	t.Parallel()
	rng := NewRNGService()

	tests := []struct {
		name     string
		length   int
		max      int
		wantErr  bool
		wantLen  int
		checkVal func([]int) error
	}{
		{
			name:     "valid slice",
			length:   5,
			max:      10,
			wantErr:  false,
			wantLen:  5,
			checkVal: checkSliceValues,
		},
		{
			name:    "negative length",
			length:  -1,
			max:     10,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "max less than 1",
			length:  5,
			max:     0,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "zero length",
			length:  0,
			max:     10,
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			slice, err := rng.GenerateSliceWithMax(tt.length, tt.max)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateSliceWithMax() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(slice) != tt.wantLen {
				t.Errorf("expected slice length %d, got %d", tt.wantLen, len(slice))
			}

			if tt.checkVal != nil {
				if err := tt.checkVal(slice); err != nil {
					t.Errorf("value check failed: %v", err)
				}
			}
		})
	}
}

func checkSliceValues(slice []int) error {
	for _, num := range slice {
		if num < 0 || num >= 10 {
			return fmt.Errorf("number %d is out of bounds [0, %d)", num, 10)
		}
	}
	return nil
}

func TestGenerateSlice(t *testing.T) {
	t.Parallel()
	rng := NewRNGService()

	tests := []struct {
		name    string
		length  int
		wantErr bool
		wantLen int
	}{
		{
			name:    "valid slice",
			length:  5,
			wantErr: false,
			wantLen: 5,
		},
		{
			name:    "negative length",
			length:  -1,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "zero length",
			length:  0,
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			slice, err := rng.GenerateSlice(tt.length)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateSlice() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(slice) != tt.wantLen {
				t.Errorf("expected slice length %d, got %d", tt.wantLen, len(slice))
			}
		})
	}
}

func TestGenerateDigit(t *testing.T) {
	t.Parallel()

	rngSvc := NewRNGService()

	for i := 0; i < 100; i++ {
		digit, err := rngSvc.GenerateDigit()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if digit < 0 || digit > 9 {
			t.Errorf("generated digit is out of range: got %d, want 0-9", digit)
		}
	}
}
