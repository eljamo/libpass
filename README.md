# libpass

A Go package for generating secure and memorable passwords. It serves as the backbone for [mempass](https://github.com/eljamo/mempass) and [mempass-api](https://github.com/eljamo/mempass-api)

## Install

```
go get github.com/eljamo/libpass/v8
```

## Basic Usage

```
package main

import (
	"fmt"

	"github.com/eljamo/libpass/v8/config"
	"github.com/eljamo/libpass/v8/service"
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

```bash
$ go run ./example/example.go

[^^77%TWIKI%ardently%STORM%58^^ ^^57.HIGH.DOLL.GRAY.67^^ ::90:passive:FEELS:WASTING:40::]
```

### Run the tests

```bash
go test --race --shuffle on ./...
```
