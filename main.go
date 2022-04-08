package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "args %s\n", os.Args[1:])
	_, _ = fmt.Fprintf(os.Stdout, "::set-output name=time::%s\n", time.Now().Format(time.RFC3339))
}
