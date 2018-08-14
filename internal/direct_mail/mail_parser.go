package direct_mail

// メールログddjkdをパースして、FROMアドレスを取得する

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

type MailParser struct {
}

func NewMailParser() *MailParser {
	return &MailParser{}
}

func (p *MailParser) Parser(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("[(MailParser) Parser] Open file error: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Scan first line to get return path for extracting mail address
	scanner.Scan()
	returnPath := scanner.Text()

	// Mail log not invalid
	if !strings.HasPrefix(returnPath, "Return-Path:") {
		return "", fmt.Errorf("[(MailParser) Parser] Mail log not invalid: path = %v", path)
	}

	mailAddress := strings.TrimLeft(returnPath, "Return-Path: <")
	mailAddress = strings.TrimRight(mailAddress, ">")

	return mailAddress, nil
}