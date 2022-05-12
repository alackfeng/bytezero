package utils

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"time"
)


var randomLetters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var randomLettersLen = len(randomLetters)

// RandomBytes -
func RandomBytes(len int, r []byte) []byte {
    if r == nil {
        r = make([]byte, len)
    }
    t := time.Now().UnixNano()
    for i:=0; i<len; i++ {
        n, _ := rand.Int(rand.Reader, big.NewInt(t))
        r[i] = randomLetters[n.Int64() % int64(randomLettersLen)]
    }
    return r
}

// RandomBytesQ -
func RandomBytesQ(len int, r []byte) []byte {
    if r == nil {
        r = make([]byte, len)
    }
    t := time.Now().UnixNano()
    mrand.Seed(t)
    for i:=0; i<len; i++ {
        n := mrand.Int63n(t)
        r[i] = randomLetters[n % int64(randomLettersLen)]
    }
    return r
}
