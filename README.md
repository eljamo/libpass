# libpass

A Go package for generating secure and memorable passwords. It serves as the backbone for [mempass](https://github.com/eljamo/mempass)

## Install

```
go get github.com/eljamo/libpass/v7
```

## Basic Usage

```
package main

import (
	"fmt"

	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/service"
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

## Contributing

### Clone the repo

```bash
git clone https://github.com/eljamo/libpass
cd libpass
```

### Run the tests

```bash
go test ./...
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
