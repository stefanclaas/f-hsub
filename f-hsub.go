package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
)

type hsub struct {
	key     string
	subject string
}

func (h *hsub) hsubtest() bool {
	sublen := len(h.subject)
	if sublen <= 32 || sublen > 96 {
		fmt.Println("Error: not a valid hsub")
		return false
	}
	iv, err := hex.DecodeString(h.subject[:16])
	if err != nil {
		return false
	}
	digest := sha256.New()
	digest.Write(iv)
	digest.Write([]byte(h.key))
	newhsub := hex.EncodeToString(append(iv, digest.Sum(nil)...))[:sublen]
	return newhsub == h.subject
}

func findValidSubjectsInFile(filename string, key string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var outputFile *os.File
	var headers []string

	scanner := bufio.NewScanner(file)
Loop:
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ".") {
			if outputFile != nil {
				outputFile.Close()
				outputFile = nil
			}
			headers = nil
			continue Loop
		}
		if strings.Contains(line, "Subject:") {
			parts := strings.Split(line, "Subject:")
			if len(parts) > 1 {
				h := &hsub{key: key, subject: strings.TrimSpace(parts[1])}
				if h.hsubtest() {
					if outputFile != nil {
						outputFile.Close()
					}
					outputFileName := fmt.Sprintf("valid_hsub_%s.txt", h.subject)
					outputFile, err = os.Create(outputFileName)
					if err != nil {
						fmt.Println("Error creating output file:", err)
						os.Exit(1)
					}
					fmt.Println("Valid Subject:", h.subject)
					for _, header := range headers {
						fmt.Fprintln(outputFile, header)
					}
					headers = nil
				} else {
					continue
				}
			}
		}

		if outputFile != nil {
			fmt.Fprintln(outputFile, line)
		} else {
			headers = append(headers, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	if outputFile != nil {
		outputFile.Close()
	}
}

func main() {
	flag.Parse()
	cmdargs := flag.Args()
	switch len(cmdargs) {
	case 2:
		findValidSubjectsInFile(cmdargs[0], cmdargs[1])
	default:
		fmt.Println("Usage: f-hsub <filename> <key>")
		os.Exit(2)
	}
}
