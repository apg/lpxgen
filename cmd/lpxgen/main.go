package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/apg/lpxgen"
)

var (
	count      = flag.Int("count", 1000, "Number of batches to send")
	uniqTokens = flag.Int("tokens", 10, "Number of tokens to utilize")
	minbatch   = flag.Int("min", 1, "Minimum number of messages in a batch")
	maxbatch   = flag.Int("max", 100, "Maximum number of messages in a batch")
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

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: No URL given\n\n")
		flag.Usage()
		os.Exit(1)
	}

	url := flag.Arg(0)

	if *minbatch > *maxbatch {
		tmp := *minbatch
		*minbatch = *maxbatch
		*maxbatch = tmp
	} else if *minbatch == *maxbatch {
		*maxbatch++
	}

	gen := lpxgen.NewGenerator(*minbatch, *maxbatch, lpxgen.ProbLogFromString(*logdist))

	client := &http.Client{}
	for i := 0; i < *count; i++ {
		req := gen.Generate(url)

		if resp, err := client.Do(req); err != nil {
			fmt.Fprintf(os.Stderr, "Error while performing request: %q\n", err)
		} else if resp.Status[0] != '2' {
			fmt.Fprintf(os.Stderr, "Non 2xx response recieved: %s\n", resp.Status)
		}
	}
}
