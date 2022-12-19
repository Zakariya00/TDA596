package main

import (
	"crypto/sha1"
	"math/big"
)

/* Hashing and Related Functions */

const keySize = sha1.Size * 8

var two = big.NewInt(2)
var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)

// hash get a sha1 hash value
func hash(elt string) *big.Int {
	hasher := sha1.New()
	hasher.Write([]byte(elt))
	return new(big.Int).SetBytes(hasher.Sum(nil))
}

// jump  computes the address of a position across the ring that should be pointed to by the given finger table entry
func jump(address string, fingerentry int) *big.Int {
	n := hash(address)
	fingerentryminus1 := big.NewInt(int64(fingerentry) - 1)
	jump := new(big.Int).Exp(two, fingerentryminus1, nil)
	sum := new(big.Int).Add(n, jump)

	return new(big.Int).Mod(sum, hashMod)
}

func between(start, elt, end *big.Int, inclusive bool) bool {
	if end.Cmp(start) > 0 {
		return (start.Cmp(elt) < 0 && elt.Cmp(end) < 0) || (inclusive && elt.Cmp(end) == 0)
	} else {
		return start.Cmp(elt) < 0 || elt.Cmp(end) < 0 || (inclusive && elt.Cmp(end) == 0)
	}
}

// between1 returns true if elt is between start and end on the ring, accounting for the boundary where the ring
// loops back on itself. If inclusive is true, it tests if elt is in (start,end],
// otherwise it tests for (start,end).
func between1(start, elt, end string, inclusive bool) bool {
	bStart, _ := new(big.Int).SetString(start, 10)
	bElt, _ := new(big.Int).SetString(elt, 10)
	bEnd, _ := new(big.Int).SetString(end, 10)

	return between(bStart, bElt, bEnd, inclusive)
}
