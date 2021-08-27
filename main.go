package main

//go:generate sqlboiler crdb

import "go.hollow.sh/serverservice/cmd"

func main() {
	cmd.Execute()
}
