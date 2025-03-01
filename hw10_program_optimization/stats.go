package hw10programoptimization

import (
	"bufio"
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

//type users [100_000]User // Array

//var (
//	emailRegex = regexp.MustCompile(`(?i)\.[a-z]{2,4}$`)
//	//regexp.MustCompile(`^[a-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[a-z]{2,4}$`)
//)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {

	stats := make(DomainStat)

	fScanner := bufio.NewScanner(r)
	fScanner.Split(bufio.ScanLines)
	domainRegex := //regexp.MustCompile("(?i)^(?:[^@]+@)(.+\\." + regexp.QuoteMeta(domain) + ")$")
		//regexp.MustCompile("(?i)@(.+?)\\." + regexp.QuoteMeta(domain) + "$")
		//regexp.MustCompile("(?i)@(.+\\." + regexp.QuoteMeta(domain) + ")$") // !!!
		//regexp.MustCompile("\\." + regexp.QuoteMeta(domain))
		regexp.MustCompile(`(?i)"Email"\s*:\s*"[^@]+@(.+\.` + regexp.QuoteMeta(domain) + `)"`)
	for fScanner.Scan() {
		err := calcDomainsNonJson(fScanner.Bytes(), &stats, domainRegex)
		if err != nil {
			return nil, err
		}
		//var user User
		//
		//err := json.Unmarshal(fScanner.Bytes(), &user)
		//if err != nil {
		//	continue
		//}
		//err = calcDomains(user, &stats, domainRegex)
		//if err != nil {
		//	//	return nil, fmt.Errorf("failed calc domain", err)
		//	return nil, err
		//}
	}
	return stats, nil
}

func calcDomainsNonJson(b []byte, stats *DomainStat, regexp *regexp.Regexp) error {

	matches := regexp.FindSubmatch(b)
	if len(matches) <= 1 {
		return nil
	}

	//match := regexp.Match(b)
	//if !match {
	//	return nil
	//}

	//_, domain, found := strings.Cut(u.Email, "@")
	//if !found {
	//	return fmt.Errorf("invalid email address")
	//}
	(*stats)[(strings.ToLower(string(matches[1])))]++
	//(*stats)[(domain)]++
	//(*stats)[customToLowe(matches[1])]++
	return nil
}

//
//func calcDomains(u User, stats *DomainStat, regexp *regexp.Regexp) error {
//
//	//matches := regexp.FindStringSubmatch(u.Email)
//	//if len(matches) <= 1 {
//	//	return nil
//	//}
//
//	match := regexp.MatchString(u.Email)
//	if !match {
//		return nil
//	}
//
//	_, domain, found := strings.Cut(u.Email, "@")
//	if !found {
//		return fmt.Errorf("invalid email address")
//	}
//	(*stats)[(strings.ToLower(domain))]++
//	//(*stats)[(domain)]++
//	//(*stats)[customToLowe(matches[1])]++
//	return nil
//}

//func customToLower(s string) string {
//	for i := 0; i < len(s); i++ {
//		c := s[i]
//		//
//		if c >= 'A' && c <= 'Z' {
//
//			b := make([]byte, len(s))
//			copy(b, s[:i]) // здесь уже в нижнем регистре
//			for j := i; j < len(s); j++ {
//				c = s[j]
//				if c >= 'A' && c <= 'Z' {
//					b[j] = c + 32 // для ASCII сдвигаемся: A-65 => a-97
//				} else {
//					b[j] = c
//				}
//			}
//			return string(b)
//		}
//	}
//	return s
//}
