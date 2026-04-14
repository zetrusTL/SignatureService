package domain

import "strings"

type Algorithm string

const (
	AlgorithmRSA Algorithm = "RSA"
	AlgorithmECC Algorithm = "ECC"
)

func ParseAlgorithm(s string) (Algorithm, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "RSA":
		return AlgorithmRSA, nil
	case "ECC":
		return AlgorithmECC, nil
	default:
		return "", ErrUnknownAlgorithm
	}
}
