package main

import (
	"fmt"
	"strings"
)
const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvwxyz-=[];',.!@#$%^&*()_+{}:|<>? "
const charsetLen = len(charset)
func XOR(input, key string) (output string) {
	for i := 0; i < len(input); i++ {
		var inputChar = strings.Index(charset, string(input[i]))
		m := string(key[i % len(key)])
		keyChar := strings.Index(charset, m)
		var xorIndex = inputChar ^ keyChar
		if xorIndex >= charsetLen {
			xorIndex = inputChar
		}
		output += string(charset[xorIndex % charsetLen])
	}
	return output
}
func main() {

	const secretKey = "thatwastooeasy"// solution from the first challenge
	fmt.Println(XOR("i04 CIR17QMJ G3: C8NU3IVIIIO3H", secretKey))
}