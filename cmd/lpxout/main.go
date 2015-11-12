package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/apg/lpxgen"
)

var (
	count      = flag.Int("count", 100, "Number of lines to emit.")
	uniqTokens = flag.Int("tokens", 10, "Number of tokens to utilize")
	logdist    = flag.String("dist", "default", "Distribution of log types. <type>:0.9,<type>:0.1")
)

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] URL\n", os.Args[0])
		flag.PrintDefaults()
	}

	if *uniqTokens > 0 {
		lpxgen.UniqTokens = *uniqTokens
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: tokens must be greater than 0\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *count <= 0 {
		fmt.Fprintf(os.Stderr, "ERROR: tokens must be greater than 0\n\n")
		flag.Usage()
		os.Exit(1)
	}

	log := lpxgen.ProbLogFromString(*logdist)

	for i := 0; i < *count; i++ {
		print(log.String())
	}
}
