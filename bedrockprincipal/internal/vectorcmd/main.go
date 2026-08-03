package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"go.minekube.com/connect/bedrockprincipal/internal/vectorgen"
)

func main() {
	check := flag.String("check", "", "verify that the checked-in vector file regenerates exactly")
	out := flag.String("out", "", "write regenerated vectors to this explicit path")
	flag.Parse()
	if (*check == "") == (*out == "") {
		log.Fatal("exactly one of -check or -out is required")
	}
	path := *check
	if path == "" {
		path = *out
	}
	checked, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	regenerated, err := vectorgen.Regenerate(checked)
	if err != nil {
		log.Fatal(err)
	}
	if *check != "" {
		if !bytes.Equal(checked, regenerated) {
			log.Fatal("checked-in vectors differ from deterministic regeneration")
		}
		fmt.Println("core vectors regenerate exactly")
		return
	}
	if err := os.WriteFile(path, regenerated, 0o644); err != nil {
		log.Fatal(err)
	}
}
