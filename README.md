# libpass

A Go library specifically designed for generating secure and memorable passwords. It serves as the backbone for [mempass](https://github.com/eljamo/mempass)


## Install

```
go get github.com/eljamo/libpass/v6
```

## Basic Usage

```
package main

import (
	"fmt"

	"github.com/eljamo/libpass/v6/config"
	"github.com/eljamo/libpass/v6/service"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
	}

	svc, err := service.NewPasswordGeneratorService(cfg)
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