package utils

import (
	"math/rand/v2"
	"strings"
)

var characters = "abcdefghijlkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

const id_len = 12

func IdGen() string {
	return IdGenOfLength(id_len)
}

func IdGenOfLength(length int) string {
	var id strings.Builder
	for i := 0; i < length; i++ {
		id.WriteByte(characters[rand.IntN(len(characters))])
	}

	return id.String()
}
