package main

import (
	"fmt"

	"github.com/reconmap/cli/internal/api"
)

func main() {
	const usage = `Reconmap pentest automation tool.

Usage: rmap [OPTIONS] COMMAND

Commands
 - get clients|projects|tasks|vulnerabilities
 - create clients|projects|tasks|vulnerabilities
 - import
 - run
 - upload

Find out more information at https://reconmap.org/.
`
	fmt.Println(usage)

	fmt.Println(api.RetrieveData())
}
