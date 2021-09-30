package main

import "fmt"

// Version information.
var (
	BuildTS   = "None"
	GitHash   = "None"
	GitBranch = "None"
	Version   = "None"
)

//go:generate go run -ldflags '-X "main.BuildTS=`date`" -X "main.GitHash=`git rev-parse HEAD`" -X "main.GitBranch=`git rev-parse --abbrev-ref HEAD`"' main.go
func main() {
	fmt.Println("Version:          ", Version)
	fmt.Println("Git Branch:       ", GitBranch)
	fmt.Println("Git Commit:       ", GitHash)
	fmt.Println("Build Time (UTC): ", BuildTS)
}
