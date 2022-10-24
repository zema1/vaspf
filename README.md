# Package vaspf

Changes:
- add context parameter to all net calls
- fix some bugs in redirect cases
- add go module support

[![Documentation](https://godoc.org/github.com/zema1/vaspf?status.svg)](http://godoc.org/github.com/asggo/spf)

Package spf parses an SPF record and determines if a given IP address
is allowed to send email based on that record. SPF handles all of the
mechanisms defined at http://www.open-spf.org/SPF_Record_Syntax/.

## Example

```
go get -u "github.com/zema1/vaspf"
```

```Go
package main

import (
	"context"
	"github.com/zema1/vaspf"
)

func main() {

	SMTPClientIP := "1.1.1.1"
	envelopeFrom := "info@example.com"
	ctx := context.Background()

	result, err := spf.SPFTest(ctx, SMTPClientIP, envelopeFrom)
	if err != nil {
		panic(err)
	}

	switch result {
	case spf.Pass:
		// allow action
	case spf.Fail:
		// deny action
	}
	//...
}

```
