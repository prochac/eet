# goEET

Go EET Client

## Installation

```sh
$ go get github/prochac/goEET
```

## Example

```go
package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/prochac/goEET"
)

func main(){
    d, _ := goEET.NewDispatcher(goEET.PlaygroundService,
		"cert_test/EET_CA1_Playground-CZ00000019.key",
		"cert_test/EET_CA1_Playground-CZ00000019.pem",
		"")
    r := goEET.Receipt{
		UuidZpravy: "49ee3022-de4e-447c-b07f-a550b2378410",
		DicPopl:    "CZ00000019",
		IdProvoz:   273,
		IdPokl:     "/5546/RO24",
		PoradCis:   "0/6460/ZQ42",
		DatTrzby:   time.Now(),
		CelkTrzba:  0,
		Rezim:      goEET.RegularRegime,
    }
    response, err := d.SendPayment(r)
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Println("Fik: ", response.Fik)
}
```