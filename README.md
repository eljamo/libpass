# libpass

A Go library specifically designed for generating secure and memorable passwords. It serves as the backbone for [mempass](https://https://github.com/eljamo/mempass)


## Install

```
go get github.com/eljamo/libpass/v5
```

## Basic Usage

```
package main

import (
	"fmt"

	"github.com/eljamo/libpass/v5/config"
	"github.com/eljamo/libpass/v5/service"
)

func main() {
	specialCharacters := []string{
		"!", "@", "$", "%", "^", "&", "*", "-", "+", "=", ":", "|", "~", "?", "/", ".", ";",
	}
	config := &config.Settings{
		WordList:                "EN",
		NumPasswords:            3,
		NumWords:                3,
		WordLengthMin:           4,
		WordLengthMax:           8,
		CaseTransform:           "RANDOM",
		SeparatorCharacter:      "RANDOM",
		SeparatorAlphabet:       specialCharacters,
		PaddingDigitsBefore:     2,
		PaddingDigitsAfter:      2,
		PaddingType:             "FIXED",
		PaddingCharacter:        "RANDOM",
		SymbolAlphabet:          specialCharacters,
		PaddingCharactersBefore: 2,
		PaddingCharactersAfter:  2,
	}

	svc, err := service.NewPasswordGeneratorService(config)
	if err != nil {
		fmt.Println(err)
	}

	pws, err := svc.Generate()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pws)
}
```