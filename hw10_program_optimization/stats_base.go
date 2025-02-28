package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func GetDomainStat1(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers1(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers1(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r) // stream
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains1(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email)) // вынести из цикла, что нибудь с предкомпиляцией поискать на го
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
