package service

import (
	"testing"
	"unicode/utf8"

	"github.com/eljamo/libpass/v7/asset"
	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
)

const presetContractIterations = 25

// presetLengthBounds pins the documented length guarantees of the embedded
// presets, measured in runes. A min of 0 means no minimum; presets without an
// entry only guarantee non-empty passwords.
var presetLengthBounds = map[string]struct{ min, max int }{
	option.PresetNTLM:          {14, 14},
	option.PresetWeb16:         {0, 16},
	option.PresetWeb16XKPasswd: {0, 16},
	option.PresetWeb32:         {0, 32},
	option.PresetWiFi:          {63, 63},
}

// TestPresetContracts generates passwords from every embedded preset and
// asserts the length guarantees documented in option.PresetDescriptionMap,
// so a preset edit which breaks a documented promise fails CI.
func TestPresetContracts(t *testing.T) {
	t.Parallel()

	for _, preset := range option.Presets {
		t.Run(preset, func(t *testing.T) {
			t.Parallel()

			pm, err := asset.GetJSONPreset(preset)
			if err != nil {
				t.Fatalf("GetJSONPreset(%q) error = %v", preset, err)
			}

			cfg, err := config.New(pm)
			if err != nil {
				t.Fatalf("config.New(%q) error = %v", preset, err)
			}

			svc, err := NewPasswordGeneratorService(cfg)
			if err != nil {
				t.Fatalf("NewPasswordGeneratorService(%q) error = %v", preset, err)
			}

			bounds, hasBounds := presetLengthBounds[preset]
			for range presetContractIterations {
				pws, err := svc.Generate()
				if err != nil {
					t.Fatalf("Generate(%q) error = %v", preset, err)
				}

				for _, pw := range pws {
					pwLen := utf8.RuneCountInString(pw)
					if pwLen == 0 {
						t.Fatalf("Generate(%q) returned an empty password", preset)
					}

					if !hasBounds {
						continue
					}

					if pwLen < bounds.min || pwLen > bounds.max {
						t.Fatalf(
							"Generate(%q) password %q is %d runes, want between %d and %d",
							preset, pw, pwLen, bounds.min, bounds.max,
						)
					}
				}
			}
		})
	}
}
