# Go EET Client

Implementation of EET client in Go

## Installation

```sh
$ go get github/prochac/eet
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/prochac/eet"
)

func main(){
    d, _ := eet.NewDispatcher(eet.PlaygroundService,
		"EET_CA1_Playground-CZ00000019.p12",
		"eet",
    )
    r := eet.Receipt{
		UuidZpravy: "49ee3022-de4e-447c-b07f-a550b2378410",
		DicPopl:    "CZ00000019",
		IdProvoz:   273,
		IdPokl:     "/5546/RO24",
		PoradCis:   "0/6460/ZQ42",
		DatTrzby:   time.Now(),
		CelkTrzba:  0,
		Rezim:      eet.RegularRegime,
    }
    response, err := d.SendPayment(r)
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Println("Fik: ", response.Fik)
}
```

## Thanks

Thanks for help and inspiration
 - https://github.com/ondrejnov/eet
 - https://github.com/v154c1/pyEET
