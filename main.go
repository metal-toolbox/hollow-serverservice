package main

//go:generate sqlboiler crdb --add-soft-deletes

import "go.hollow.sh/serverservice/cmd"

// @title           serverservice API
// @version         1.0
// @description     serverservice API holds hardware and asset information

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1

func main() {
	cmd.Execute()
}
