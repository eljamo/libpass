package merger

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		maps []map[string]any
		want map[string]any
	}{
		{
			name: "No maps",
			maps: []map[string]any{},
			want: map[string]any{},
		},
		{
			name: "Single map",
			maps: []map[string]any{
				{"key1": "value1", "key2": "value2"},
			},
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "Multiple maps with unique keys",
			maps: []map[string]any{
				{"key1": "value1"},
				{"key2": "value2"},
				{"key3": "value3"},
			},
			want: map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			name: "Multiple maps with overlapping keys",
			maps: []map[string]any{
				{"key1": "value1"},
				{"key1": "value2"},
				{"key1": "value3"},
			},
			want: map[string]any{
				"key1": "value3",
			},
		},
		{
			name: "Mixed keys",
			maps: []map[string]any{
				{"key1": "value1", "key2": "value2"},
				{"key2": "newValue2", "key3": "value3"},
			},
			want: map[string]any{
				"key1": "value1",
				"key2": "newValue2",
				"key3": "value3",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Map(tt.maps...)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
