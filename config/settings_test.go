package config

import (
	"testing"

	"github.com/eljamo/libpass/v7/config/option"
	"github.com/google/go-cmp/cmp"
)

func TestDefaultSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want *Settings
	}{
		{
			name: "default settings",
			want: &Settings{
				CaseTransform:           option.CaseTransformRandom,
				NumPasswords:            3,
				NumWords:                3,
				PaddingCharacter:        option.PaddingCharacterRandom,
				PaddingCharactersAfter:  2,
				PaddingCharactersBefore: 2,
				PaddingDigitsAfter:      2,
				PaddingDigitsBefore:     2,
				PaddingType:             option.PaddingTypeFixed,
				Preset:                  option.PresetDefault,
				SeparatorAlphabet:       option.DefaultSpecialCharacters,
				SeparatorCharacter:      option.SeparatorCharacterRandom,
				SymbolAlphabet:          option.DefaultSpecialCharacters,
				WordLengthMax:           8,
				WordLengthMin:           4,
				WordList:                option.WordListEN,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DefaultSettings()

			if !cmp.Equal(got, tt.want) {
				t.Errorf("DefaultSettings() = %+v, want %+v\nDiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
