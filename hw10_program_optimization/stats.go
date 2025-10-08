package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	result = make(users, 0, 100000)

	scanner := bufio.NewScanner(r)
	var user User

	for scanner.Scan() {
		if err = user.UnmarshalJSON(scanner.Bytes()); err != nil {
			return result, fmt.Errorf("unmarshal error: %w", err)
		}
		result = append(result, user)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	targetDomain := strings.ToLower(domain)

	for _, user := range u {
		email := strings.ToLower(user.Email)

		index := strings.LastIndex(email, "@")
		if index == -1 {
			continue
		}

		fullDomain := email[index+1:]
		if strings.HasSuffix(fullDomain, "."+targetDomain) {
			result[fullDomain]++
		}
	}

	return result, nil
}
