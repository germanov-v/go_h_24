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
	emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[a-z]{2,4}$`) // !!!!ПОДСМОТРЕЛ!!!! -
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	//u, err := getUsers(r)
	//if err != nil {
	//	return nil, fmt.Errorf("get users error: %w", err)
	//}
	//return countDomains(u, domain)

	//file, err := os.Open(domain)
	//if err != nil {
	//	return nil, fmt.Errorf("couldn't open the file - alas and ah", err)
	//}
	//defer file.Close()

	stats := make(DomainStat)

	fScanner := bufio.NewScanner(r)
	fScanner.Split(bufio.ScanLines)

	for fScanner.Scan() {
		var user User
		//text := fScanner.Text()
		//
		//err = json.NewDecoder(strings.NewReader(text)).Decode(&user)

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
	isEmail := emailRegex.MatchString(u.Email)
	if isEmail {
		parts := strings.Split(u.Email, "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid email address")
		}
		domain := strings.ToLower(parts[1])
		num := (*stats)[domain]
		num++
		(*stats)[domain] = num
	}

	return nil
}

func countDomains(u users, domain string) (DomainStat, error) {
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
