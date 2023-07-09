package common

import (
	"math/rand"
	"time"
)

func RandomString(n1 int, n2 int) string {
	var letters = []byte("asdfghjklzcxvbnmqwertyuiopadjbhksndjkdsqwefcbuknASDBIUEBCUIQNAOQIWDN")
	var numbers = []byte("1234567890")

	result1 := make([]byte, n1)
	result2 := make([]byte, n2)

	rand.Seed(time.Now().Unix())

	for i := range result1 {
		result1[i] = letters[rand.Intn(len(letters))]
	}

	rand.Seed(time.Now().Unix())
	for i := range result2 {
		result2[i] = numbers[rand.Intn(len(numbers))]
	}

	result := append(result1, result2...)

	return string(result)
}
