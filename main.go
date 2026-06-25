package main

import (
	"bufio"
	"cmp"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"
)

type Log struct {
	Ip         string
	HttpStatus string
	Endpoint   string
}

type kv struct {
	Key   string
	Value int
}

func newLog(ip, httpstatus, endpoint string) Log {
	return Log{
		Ip:         ip,
		HttpStatus: httpstatus,
		Endpoint:   endpoint,
	}
}

func logenerateLogFile() error {
	neededBytes := 8 * 1024 * 1024
	bytes := 0

	file, err := os.Create("server.log")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	method := []string{"GET", "POST", "DELETE"}
	path := []string{"/api/v1/auth/login", "/api/v2/products/1093", "/api/v1/user/profile"}
	protocol := []string{"HTTP/1.1", "HTTP/2", "HTTP/3"}
	status := []string{"200", "301", "400", "401", "404"}

	buf := bufio.NewWriter(file)
	defer buf.Flush()
	for bytes < neededBytes {
		curTime := time.Now()
		s := fmt.Sprintf("%d.%d.%d.%d - - [%s] \"%s %s %s\" %s\n", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), curTime, method[rand.Intn(3)], path[rand.Intn(3)], protocol[rand.Intn(3)], status[rand.Intn(3)])
		n, err := buf.WriteString(s)
		if err != nil {
			return fmt.Errorf("failed writing string to buffer: %w", err)
		}
		bytes += n
	}
	return nil
}

func readFile(m map[string]int) error {
	file, err := os.Open("server.log")
	if err != nil {
		return fmt.Errorf("could not open file %w", err)
	}

	defer file.Close()
	fmt.Println("----- FIle: -----")

	var logs []Log

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		parts := strings.Fields(line)

		if len(parts) < 3 {
			continue
		}

		ip := parts[0]
		endpoint := parts[9]
		status := parts[11]
		logEntry := newLog(ip, status, endpoint)
		(m)[ip]++
		fmt.Println(logEntry.Ip, logEntry.Endpoint, logEntry.HttpStatus)
		logs = append(logs, logEntry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	fmt.Println("------ End of file -----")
	return nil
}

func main() {

	err1 := generateLogFile()
	if err1 != nil {
		fmt.Println("Error while creating file")
	}

	mostIp := map[string]int{}
	err2 := readFile(mostIp)
	if err2 != nil {
		fmt.Println("Error while reading file")
	}

	var sortedList []kv
	for k, v := range mostIp {
		sortedList = append(sortedList, kv{k, v})
	}

	slices.SortFunc(sortedList, func(i, j kv) int {
		return cmp.Compare(j.Value, i.Value)
	})

	fmt.Println("Most seen ip addresses: ")
	for i := 0; i < 5; i++ {
		fmt.Println(i+1, ":", sortedList[i].Key, ":", sortedList[i].Value)
	}
}
