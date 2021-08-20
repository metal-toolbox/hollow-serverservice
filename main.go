package main

//go:generate sqlboiler crdb

import "go.hollow.sh/dcim/cmd"

func main() {
	cmd.Execute()
}
