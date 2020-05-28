package main

import (
	"fmt"
	"log"

	"github.com/new-tel/newtel_api.go"
)

func main() {
	cfg := newtel_api.NewTelConfig{
		ApiKey:  "xxxxxxxxxxxxxxxxxxxx",
		ApiSign: "yyyyyyyyyyyyyyyyyyyyyyyyy",
	}
	client := newtel_api.NewTelClient(cfg)
	resp, err := client.CallPassword("+79081234567", "1234")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", resp)
}
