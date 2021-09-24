package main

//go:generate sqlboiler crdb --add-soft-deletes

import "go.hollow.sh/serverservice/cmd"

func main() {
	cmd.Execute()
}
