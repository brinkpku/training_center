package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	str := "tvbluyR0MwqUeiPKM1bxFA=="
	bytes, _ := base64.StdEncoding.DecodeString(str)
	fmt.Println(hex.EncodeToString(bytes))
}
