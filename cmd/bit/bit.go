package main

import (
	"flag"
	"log"
	"os"

	"github.com/blackchip-org/bit"
)

var errlog = log.New(os.Stderr, "", 0)
var prog = "bit"

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		errlog.Fatalf("invalid number of arguments\n")
	}
	b, err := bit.NewBuilder(flag.Arg(0))
	if err != nil {
		errlog.Fatalln(prog + ": error: " + err.Error())
	}

	err = b.Run()
	if err != nil {
		errlog.Fatalln(err)
	}
}
