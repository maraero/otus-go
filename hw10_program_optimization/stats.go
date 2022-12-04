package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	parser := jsoniter.ConfigCompatibleWithStandardLibrary
	result := make(DomainStat)
	ending := "." + domain

	bScanner := bufio.NewScanner(r)

	var user User
	for bScanner.Scan() {
		err := parser.Unmarshal(bScanner.Bytes(), &user)
		if err != nil {
			return nil, err
		}
		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, ending) {
			fDomain := strings.SplitN(email, "@", 2)[1]
			result[fDomain]++
		}
	}
	return result, nil
}
