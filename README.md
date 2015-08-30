# whoapi [![GoDoc](https://godoc.org/github.com/tg/whoapi?status.svg)](https://godoc.org/github.com/tg/whoapi)
A simple Go package and command line tool for accessing WhoAPI services.

## Installation
You need a go environment set up. See: http://golang.org/doc/install

Grab package and install command line tool:
```
go get github.com/tg/whoapi/...
```

## Usage
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tg/whoapi"
)

func main() {
	// Create client and fetch the data.
	// API key can be also provided in WHOAPI_KEY environmental variable.
	client := &whoapi.Client{Key: "your-api-key"}

	data, err := client.Get("whois", "github.com")
	if err != nil {
		log.Fatal(err)
		// Might be still worth checking whether we got a whois data back.
		// See whoapi.ErrWhoisNotYetSupported for more details.
	}
	
	// Extract whatever you're interested in
	var whois struct {
		Created string `json:"date_created"`
	}
	err = json.Unmarshal(data, &whois)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("github was created on", whois.Created)
}
```

## CLI
After installing the command line tool shipped with this package, you can query WhoAPI directly from your terminal.
Assuming you have your key in WHOAPI_KEY environmental variable, you can simply try the following:
```
$ whoapi taken github.com
{
	"status": "0",
	"taken": "1"
}
```
Sorry, but it's already taken.

## Caveats
WhoAPI has a bad habit of mixing JSON values for the same key.
For example, status code can be either `status: 0` or `status: "0"` (sometimes a number, sometimes a string).
This can be pain in statically checked languages like Go when you want to parse the value into a specific type.
Taking `whoapi taken` example from CLI section, if the domain is free, the following response is returned:
```
$ whoapi taken github666.com
{
	"status": "0",
	"taken": 0
}
```
This time `taken` is a number, while before it was a string. There is a `whoapi.Int64` type especially for this occasion,
which can be unmarshaled from both forms. Use it in your JSON structure when you expect this might happen.

## TODO?
- Full response structure for every request type (eliminating the need for manual JSON parsing)
- Automatic throttle limit (limits can be obtained through `myaccount` request)
