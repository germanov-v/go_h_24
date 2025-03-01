package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
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

type users [100_000]User // Array

var (
	emailRegex = regexp.MustCompile(`(?i)\.[a-z]{2,4}$`)
	//regexp.MustCompile(`^[a-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[a-z]{2,4}$`)
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {

	stats := make(DomainStat)

	fScanner := bufio.NewScanner(r)
	fScanner.Split(bufio.ScanLines)

	for fScanner.Scan() {
		var user User

		err := json.Unmarshal(fScanner.Bytes(), &user)
		if err != nil {
			continue
		}
		err = calcDomains(user, &stats)
		if err != nil {
			//	return nil, fmt.Errorf("failed calc domain", err)
			return nil, err
		}
	}
	return stats, nil
}

func calcDomains(u User, stats *DomainStat) error {

	if u.Email == "" || !emailRegex.MatchString(u.Email) {
		return nil
	}

	_, domain, found := strings.Cut(u.Email, "@")
	if !found {
		return fmt.Errorf("invalid email address")
	}
	(*stats)[strings.ToLower(domain)]++
	return nil
}
