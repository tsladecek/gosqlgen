package gosqlgen

import "math/rand"

var r = rand.New(rand.NewSource(1618))

func RandomInt(max int) int {
	return r.Intn(max)
}

func RandomString(length int, alphabet []rune) string {
	out := make([]rune, length)
	for i := range length {
		out[i] = alphabet[RandomInt(len(alphabet))]
	}

	return string(out)
}
