package asset

import (
	"slices"
	"testing"

	"github.com/eljamo/libpass/v8/config/option"
)

// These tests pin the three places a word list or preset must be registered
// (asset fileMap, the option slice, and the option description map) to each
// other, so adding one without the others fails CI.

func TestWordListOptionsMatchAssetRegistry(t *testing.T) {
	t.Parallel()

	registry := fileMap[option.ConfigKeyWordList]
	if got, want := len(option.WordLists), len(registry); got != want {
		t.Errorf("option.WordLists has %d entries, asset registry has %d", got, want)
	}

	for _, wl := range option.WordLists {
		if _, ok := registry[wl]; !ok {
			t.Errorf("option.WordLists entry %q has no file in the asset registry", wl)
		}
		if _, ok := option.WordListDescriptionMap[wl]; !ok {
			t.Errorf("option.WordLists entry %q has no entry in WordListDescriptionMap", wl)
		}
	}

	for key := range registry {
		if !slices.Contains(option.WordLists, key) {
			t.Errorf("asset registry word list %q is missing from option.WordLists", key)
		}
	}
}

func TestPresetOptionsMatchAssetRegistry(t *testing.T) {
	t.Parallel()

	registry := fileMap[option.ConfigKeyPreset]
	if got, want := len(option.Presets), len(registry); got != want {
		t.Errorf("option.Presets has %d entries, asset registry has %d", got, want)
	}

	for _, p := range option.Presets {
		if _, ok := registry[p]; !ok {
			t.Errorf("option.Presets entry %q has no file in the asset registry", p)
		}
		if _, ok := option.PresetDescriptionMap[p]; !ok {
			t.Errorf("option.Presets entry %q has no entry in PresetDescriptionMap", p)
		}
	}

	for key := range registry {
		if !slices.Contains(option.Presets, key) {
			t.Errorf("asset registry preset %q is missing from option.Presets", key)
		}
	}
}
