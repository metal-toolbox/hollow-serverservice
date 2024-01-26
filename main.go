// package main is the entry point
package main

//go:generate sqlboiler crdb --add-soft-deletes

import "github.com/metal-toolbox/fleetdb/cmd"

func main() {
	cmd.Execute()
}
