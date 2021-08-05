package main

import (
	"log"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("nvidia-smi", "--format=csv,noheader", "--query-gpu=index,memory.used")
	// cmd := exec.Command("ls", "-lah")
	data, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to call Output(): %v", err)
	}
	log.Printf("output: %s", data)
	array := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, a := range array {
		res := strings.Split(a, ",")
		log.Println(res[0], strings.Split(strings.TrimSpace(res[1]), " ")[0])
	}
}
