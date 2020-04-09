package random

import (
	"crypto/rand"
	"encoding/base64"
)

// Is needed for the testing
// When not isRandom always returns one string("1")
type Generator struct {
	isRandom bool
}

func NewGenerator(random bool) *Generator {
	return &Generator{
		isRandom: random,
	}
}

func (g *Generator) GetString() string {
	if g.isRandom {
		rnd := make([]byte, 16)
		rand.Read(rnd)
		res := base64.URLEncoding.EncodeToString(rnd)
		return res
	} else {
		return "1"
	}
}
