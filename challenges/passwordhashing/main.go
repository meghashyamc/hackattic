package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log/slog"
	"os"

	"github.com/meghashyamc/hackattic/pkg/auth"
	"github.com/meghashyamc/hackattic/pkg/problem"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

// Problem: https://hackattic.com/challenges/password_hashing

type Problem struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
	PBKDF2   PBKDF2 `json:"pbkdf2"`
	Scrypt   Scrypt `json:"scrypt"`
}

type PBKDF2 struct {
	Rounds int    `json:"rounds"`
	Hash   string `json:"hash"`
}

type Scrypt struct {
	N       int    `json:"N"`
	R       int    `json:"r"`
	P       int    `json:"p"`
	BufLen  int    `json:"buflen"`
	Control string `json:"_control"`
}

type Solution struct {
	SHA256 string `json:"sha256"`
	HMAC   string `json:"hmac"`
	PBKDF2 string `json:"pbkdf2"`
	Scrypt string `json:"scrypt"`
}

func main() {
	accessToken, err := auth.GetAccessToken()
	if err != nil {
		os.Exit(1)
	}
	problemName := "password_hashing"
	var problemDetails Problem
	if err := problem.Get(problemName, accessToken, &problemDetails); err != nil {
		os.Exit(1)
	}

	solution, err := hashPassword(&problemDetails)
	if err != nil {
		os.Exit(1)
	}

	if err := problem.Submit(problemName, accessToken, *solution); err != nil {
		os.Exit(1)
	}

}

func hashPassword(problemDetails *Problem) (*Solution, error) {

	passwordBytes := []byte(problemDetails.Password)
	saltBytes, err := base64.StdEncoding.DecodeString(problemDetails.Salt)
	if err != nil {
		slog.Error("got an unexpected error when decoding salt", "err", err)
		return nil, err
	}

	return &Solution{
		SHA256: calculateSHA256(passwordBytes),
		HMAC:   calculateHMACSHA256(passwordBytes, saltBytes),
		PBKDF2: calculatePBKDF2(passwordBytes, saltBytes, problemDetails.PBKDF2.Rounds),
		Scrypt: calculateScrypt(passwordBytes, saltBytes, problemDetails.Scrypt.N, problemDetails.Scrypt.R, problemDetails.Scrypt.P, problemDetails.Scrypt.BufLen),
	}, nil
}

func calculateSHA256(input []byte) string {
	hasher := sha256.New()
	hasher.Write(input)
	sha256Hash := hasher.Sum(nil)
	return hex.EncodeToString(sha256Hash)
}

func calculateHMACSHA256(input []byte, salt []byte) string {
	mac := hmac.New(sha256.New, salt)
	mac.Write(input)
	hmacResult := mac.Sum(nil)
	return hex.EncodeToString(hmacResult)
}

func calculatePBKDF2(input []byte, salt []byte, iterations int) string {
	pbkdf2Result := pbkdf2.Key(input, salt, iterations, 32, sha256.New)

	return hex.EncodeToString(pbkdf2Result)
}

func calculateScrypt(input []byte, salt []byte, N int, r int, p int, buflen int) string {
	scryptResult, err := scrypt.Key(input, salt, N, r, p, buflen)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(scryptResult)
}
