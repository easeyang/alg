package main

import (
	"structure_algorithm/huffuman"
)

func main() {
	str := "2233366666688888888999999999"
	bts := []byte(str)
	huffuman.Encode(bts)
}
