package validator

import (
	"testing"
)

func TestHasElementWithLengthGreaterThanOne(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    []string
		want bool
	}{
		{
			name: "Empty slice",
			s:    []string{},
			want: false,
		},
		{
			name: "All elements have length 1",
			s:    []string{"a", "b", "c"},
			want: false,
		},
		{
			name: "One element with length greater than 1",
			s:    []string{"a", "bc", "d"},
			want: true,
		},
		{
			name: "Multiple elements with length greater than 1",
			s:    []string{"abc", "de", "f"},
			want: true,
		},
		{
			name: "All elements have length greater than 1",
			s:    []string{"abc", "def", "ghi"},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := HasElementWithLengthGreaterThanOne(tt.s)
			if got != tt.want {
				t.Errorf("HasElementWithLengthGreaterThanOne(%v) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestIsElementInSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		s    []int
		e    int
		want bool
	}{
		{
			name: "Empty slice",
			s:    []int{},
			e:    1,
			want: false,
		},
		{
			name: "Element in slice",
			s:    []int{1, 2, 3},
			e:    2,
			want: true,
		},
		{
			name: "Element not in slice",
			s:    []int{1, 2, 3},
			e:    4,
			want: false,
		},
		{
			name: "Multiple occurrences of element",
			s:    []int{1, 2, 3, 2},
			e:    2,
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsElementInSlice(tt.s, tt.e)
			if got != tt.want {
				t.Errorf("IsElementInSlice(%v, %v) = %v, want %v", tt.s, tt.e, got, tt.want)
			}
		})
	}

	t.Run("Test with string slice", func(t *testing.T) {
		t.Parallel()

		strTests := []struct {
			name string
			s    []string
			e    string
			want bool
		}{
			{
				name: "Element in slice",
				s:    []string{"a", "b", "c"},
				e:    "b",
				want: true,
			},
			{
				name: "Element not in slice",
				s:    []string{"a", "b", "c"},
				e:    "d",
				want: false,
			},
		}

		for _, st := range strTests {
			st := st

			t.Run(st.name, func(t *testing.T) {
				t.Parallel()

				got := IsElementInSlice(st.s, st.e)
				if got != st.want {
					t.Errorf("IsElementInSlice(%v, %v) = %v, want %v", st.s, st.e, got, st.want)
				}
			})
		}
	})
}
