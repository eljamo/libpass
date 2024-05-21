package service

import (
	"errors"
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

var errmockRNGService = errors.New("mock RNG Service Error")

type mockErrRNGService struct{}

func (s *mockErrRNGService) GenerateWithMax(max int) (int, error) {
	return 0, errmockRNGService
}

func (s *mockErrRNGService) Generate() (int, error) {
	return 0, errmockRNGService
}

func (s *mockErrRNGService) GenerateDigit() (int, error) {
	return 0, errmockRNGService
}

func (s *mockErrRNGService) GenerateSlice(length int) ([]int, error) {
	return nil, errmockRNGService
}

func (s *mockErrRNGService) GenerateSliceWithMax(length, max int) ([]int, error) {
	return nil, errmockRNGService
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
		t.Run(tc.name, func(t *testing.T) {
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
		})
	}
}

func TestRNGGenerate(t *testing.T) {
	t.Parallel()

	rngSvc := NewRNGService()

	t.Run("Generate", func(t *testing.T) {
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

	t.Run("valid slice", func(t *testing.T) {
		length, max := 5, 10
		slice, err := rng.GenerateSliceWithMax(length, max)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(slice) != length {
			t.Errorf("expected slice length %d, got %d", length, len(slice))
		}
		for _, num := range slice {
			if num < 0 || num >= max {
				t.Errorf("number %d is out of bounds [0, %d)", num, max)
			}
		}
	})

	t.Run("negative length", func(t *testing.T) {
		_, err := rng.GenerateSliceWithMax(-1, 10)
		if err == nil {
			t.Error("expected error for negative length, got nil")
		}
	})

	t.Run("max less than 1", func(t *testing.T) {
		_, err := rng.GenerateSliceWithMax(5, 0)
		if err == nil {
			t.Error("expected error for max < 1, got nil")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		slice, err := rng.GenerateSliceWithMax(0, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(slice) != 0 {
			t.Errorf("expected empty slice, got %d elements", len(slice))
		}
	})
}

func TestGenerateSlice(t *testing.T) {
	t.Parallel()
	rng := NewRNGService()

	t.Run("valid slice", func(t *testing.T) {
		length := 5
		slice, err := rng.GenerateSlice(length)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(slice) != length {
			t.Errorf("expected slice length %d, got %d", length, len(slice))
		}
	})

	t.Run("negative length", func(t *testing.T) {
		_, err := rng.GenerateSlice(-1)
		if err == nil {
			t.Error("expected error for negative length, got nil")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		slice, err := rng.GenerateSlice(0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(slice) != 0 {
			t.Errorf("expected empty slice, got %d elements", len(slice))
		}
	})
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
