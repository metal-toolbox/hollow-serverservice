package main

//go:generate sqlboiler crdb

import "go.metalkube.net/hollow/cmd"

func main() {
	cmd.Execute()
}
