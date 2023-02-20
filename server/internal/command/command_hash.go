package command

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

func Hash(_ *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter string > ")
	str, _ := reader.ReadString('\n')

	if str == "" {
		return errors.New("string cannot be empty")
	}

	str = strings.TrimSpace(str)

	enc, err := bcrypt.GenerateFromPassword([]byte(str), 5)

	if err != nil {
		return err
	}

	log.Println("Result: " + string(enc))

	return nil
}
