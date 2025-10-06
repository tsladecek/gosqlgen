package gosqlgen

import "math/rand"

var r = rand.New(rand.NewSource(1618))

func randomInt(max int) int {
	return r.Intn(max)
}

func randomString(length int, alphabet []rune) string {
	out := make([]rune, length)
	for i := range length {
		out[i] = alphabet[randomInt(len(alphabet))]
	}

	return string(out)
}
