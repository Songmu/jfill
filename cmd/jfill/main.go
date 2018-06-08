package main

import (
	"os"

	"github.com/Songmu/jfill"
)

func main() {
	os.Exit(jfill.Run(os.Args[1:]))
}
