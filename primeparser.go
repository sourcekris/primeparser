// Package primeparser parses the list of large prime numbers from primes.utm.edu.
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	primesURI = "https://primes.utm.edu/primes/lists/all.txt"
	primeLine = regexp.MustCompile(`^\s+(\d+)[a-z]?\s+(\S+)\s+(\d+)\s+(\S+)\s+(.*)`)
)

const (
	mersenne = iota
	lucas
)

type prime struct {
	expression  string
	digits      int
	description string
}

func (p *prime) String() string {
	return fmt.Sprintf("Prime: %s, Num. Digits: %d, Desc: %s", p.expression, p.digits, p.description)
}

// getPrimes gets a http.Response object containing the prime database and is abstracted out
// so it can be replaced in unit tests.
func getPrimes(c *http.Client, url string) (*http.Response, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func main() {
	hc := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	r, err := getPrimes(hc, primesURI)
	if err != nil {
		log.Fatalf("failed downloading primes: %v", err)
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		log.Fatalf("failed download, received unexpected status code: %d", r.StatusCode)
	}

	// db, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Fatalf("failed reading database from stream: %v", err)
	// }
	var primes []*prime
	sc := bufio.NewScanner(r.Body)
	for sc.Scan() {
		if primeLine.MatchString(sc.Text()) {
			r := primeLine.FindStringSubmatch(sc.Text())
			if len(r) < 5 {
				fmt.Printf("regexerr: %s\n", sc.Text())
				continue
			}

			digits, _ := strconv.Atoi(r[3])
			primes = append(primes, &prime{
				expression:  r[2],
				digits:      digits,
				description: r[5],
			})
		}
	}

	fmt.Printf("got %d prime expressions\n", len(primes))

	for _, p := range primes {
		fmt.Println(p)
	}

	// TODO(sewid): Scan file line by line, parse prime components with regexp.
	// calculate the prime using the various number sequence algorithms.
}
