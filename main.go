package main

import (
	"flag"
	"fmt"
	"os"
)

var tokenFlag string

func main() {
	_, _ = fmt.Fprintf(os.Stdout, "the token is set? %v\n", tokenFlag != "please-set-me")
	result := "great success"
	_, _ = fmt.Fprintf(os.Stdout, "::set-output name=result::%s\n", result)
}

func init() {
	flag.StringVar(&tokenFlag, "github-token", "", "GitHub Access Token for accessing the wiki repo")
	flag.Parse()
}
