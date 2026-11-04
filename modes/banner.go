package modes

import (
	"fmt"
	"os"
)

const banner = `
   ____  ___  _  _   ____   ___  _
  / ___|/ _ \| || | / ___| / _ \| |
 | |  _| | | | || |_\___ \| | | | |
 | |_| | |_| |__   _|___) | |_| | |___
  \____|\___/   |_| |____/ \__\_\_____|

`

// printBanner writes to stderr because stdout carries query results, which
// callers may pipe or diff (e.g. `go4sql -stream < queries.sql > out`).
func printBanner() {
	fmt.Fprint(os.Stderr, banner)
}
