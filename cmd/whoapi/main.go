package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/tg/whoapi"
)

func main() {
	// Flags
	var (
		key string // api key (if not set $WHOAPI_KEY will be used by whoapi package)
		raw bool   // print raw response
	)
	flag.StringVar(&key, "key", "", "API key (if not set WHOAPI_KEY environment variable will be used)")
	flag.BoolVar(&raw, "raw", false, "Print raw response from server")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: whoapi [options] request [domain]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `			
  See https://whoapi.com/api-functions.html for possible request types.  
			
Examples:
		
  whoapi whois whoapi.com
  whoapi cert google.com
  whoapi myaccount
`)
	}

	flag.Parse()

	// Process arguments
	var req, domain string
	switch args := flag.Args(); len(args) {
	case 2:
		domain = args[1]
		fallthrough
	case 1:
		req = args[0]
	case 0:
		flag.Usage()
		os.Exit(2)
	default:
		fmt.Println("Too many arguments")
		os.Exit(2)
	}

	// Fetch data
	client := &whoapi.Client{Key: key}
	data, err := client.GetRaw(req, domain)

	if err != nil {
		fmt.Println(err)
		os.Exit(3) // failed to get data from server
	}

	if raw {
		os.Stdout.Write(data)
	} else {
		var out bytes.Buffer
		err = json.Indent(&out, data, "", "\t")
		if err != nil {
			fmt.Println(err)
			os.Exit(4) // failed to decode json
		}

		out.WriteTo(os.Stdout)
	}
	fmt.Println()
}
