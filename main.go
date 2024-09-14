package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/blakehulett7/goToYourMenu"
)

var startMenuOptions = []goToYourMenu.MenuOption{
	{
		Name:    "Login",
		Command: func() { fmt.Println("Login") },
	},
	{
		Name:    "Create Account",
		Command: CreateAccount,
	},
	{
		Name:    "Exit",
		Command: func() { os.Exit(0) },
	},
}

const hostUrl = "http://localhost:8080"
const hostVersion = "v1"

func main() {
	for {
		Run("clear")
		fmt.Println("Christ is King!")
		fmt.Println("\nWelcome to Folklore!")
		goToYourMenu.Menu(startMenuOptions)
	}
}

func Run(program string, args ...string) {
	command := exec.Command(program, args...)
	command.Stdout = os.Stdout
	command.Run()
}

func IsValidInput(input string) bool {
	if strings.Contains(input, ";") {
		return false
	}
	return true
}

func CreateAccount() {
	username, password := CreateUsernameAndPassword()
	res := SendUsernameAndPasswordToServer(username, password)
	resStruct := struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{}
	json.NewDecoder(res.Body).Decode(&resStruct)
	fmt.Println("Status:", res.Status)
	fmt.Println("Body:", resStruct)
	prompt := bufio.NewScanner(os.Stdin)
	prompt.Scan()
}

func CreateUsernameAndPassword() (string, string) {
	prompt := bufio.NewScanner(os.Stdin)
	var username string
	var password string
	for {
		Run("clear")
		fmt.Println("Create an Account")
		fmt.Print("\nEnter a username > ")
		prompt.Scan()
		username = prompt.Text()
		if !IsValidInput(username) {
			fmt.Println("';' character not allowed...")
			prompt.Scan()
			continue
		}
		break
	}
	fmt.Print("Enter a password > ")
	prompt.Scan()
	password = prompt.Text()
	return username, password
}

func SendUsernameAndPasswordToServer(username, password string) *http.Response {
	payloadStruct := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}
	payload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Println("couldn't marshal json:", err)
		return nil
	}
	requestURL := fmt.Sprintf("%v/%v/users", hostUrl, hostVersion)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("couldn't create http request:", err)
		return nil
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("couldn't execute http request:", err)
		return nil
	}
	return res
}
