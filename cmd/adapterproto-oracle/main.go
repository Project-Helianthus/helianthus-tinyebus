package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Project-Helianthus/helianthus-tinyebus/firmware/adapterproto"
)

func main() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(adapterproto.DefaultOracleReport()); err != nil {
		log.Fatal(err)
	}
}
