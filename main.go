package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/blakehulett7/goToYourMenu"
	"github.com/joho/godotenv"
)

var startMenuOptions = []goToYourMenu.MenuOption{
	{
		Name:    "Login",
		Command: Login,
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
const maxPermissions = 0777

var errBadToken = errors.New("bad token")

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

func main() {
	godotenv.Load()
	token := os.Getenv("JWT")
	var user User
	var err error
	if token != "" {
		user, err = GetUser()
	}
	fmt.Println(user, err)
	bufio.NewScanner(os.Stdin).Scan()
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

func CreateAccount() {
	username, password := CreateUsernameAndPassword()
	res := SendUsernameAndPasswordToServer(username, password, "users")
	resStruct := struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{}
	json.NewDecoder(res.Body).Decode(&resStruct)
	fmt.Println("Status:", res.Status)
	credentials := fmt.Sprintf("JWT=%v\nREFRESH_TOKEN=%v", resStruct.Token, resStruct.RefreshToken)
	os.WriteFile(".env", []byte(credentials), maxPermissions)
	prompt := bufio.NewScanner(os.Stdin)
	prompt.Scan()
}

func IsValidInput(input string) bool {
	if strings.Contains(input, ";") {
		return false
	}
	return true
}

func UsernameIsUnique(username string) bool {
	// The thought just occurred to me that I could also just use the UNIQUE constraint in the Sqlite DB instead...
	url := fmt.Sprintf("%v/%v/users/%v", hostUrl, hostVersion, username)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error checking if username is unique:", err)
		return false
	}
	if res.StatusCode != 200 {
		return false
	}
	return true
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
		if !UsernameIsUnique(username) {
			fmt.Println("username not available...")
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

func SendUsernameAndPasswordToServer(username, password, endpoint string) *http.Response {
	payloadStruct := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}
	payload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Println("couldn't marshal json:", err)
		return nil
	}
	requestURL := fmt.Sprintf("%v/%v/%v", hostUrl, hostVersion, endpoint)
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

func Login() {
	prompt := bufio.NewScanner(os.Stdin)
	username, password := GetUsernameAndPassword()
	if username == "" {
		return
	}
	res := SendUsernameAndPasswordToServer(username, password, "login")
	if res.StatusCode == 401 {
		fmt.Println("Incorrect password...")
		prompt.Scan()
		return
	}
	credentials := struct {
		JWT          string `json:"jwt"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err := json.NewDecoder(res.Body).Decode(&credentials)
	if err != nil {
		fmt.Println("error can't decode response:", err)
		prompt.Scan()
		return
	}
	fmt.Println("Login successful!")
	envString := fmt.Sprintf("JWT=%v\nREFRESH_TOKEN=%v", credentials.JWT, credentials.RefreshToken)
	os.WriteFile(".env", []byte(envString), maxPermissions)
	prompt.Scan()
}

func GetUsernameAndPassword() (string, string) {
	prompt := bufio.NewScanner(os.Stdin)
	var username string
	var password string
	for {
		Run("clear")
		fmt.Println("Login")
		fmt.Print("\nUsername > ")
		prompt.Scan()
		username = prompt.Text()
		if !IsValidInput(username) {
			fmt.Println("';' character not allowed")
			prompt.Scan()
			continue
		}
		if UsernameIsUnique(username) {
			fmt.Println("username not found... try 'Create Account'")
			prompt.Scan()
			return "", ""
		}
		fmt.Print("Password > ")
		prompt.Scan()
		password = prompt.Text()
		break
	}
	return username, password
}

func GetUser() (User, error) {
	token := os.Getenv("JWT")
	url := fmt.Sprintf("%v/%v/users", hostUrl, hostVersion)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("error creating request:", err)
	}
	req.Header.Add("Authorization", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error executing request:", err)
	}
	if res.StatusCode == 401 {
		return User{}, errBadToken
	}
	user := User{}
	err = json.NewDecoder(res.Body).Decode(&user)
	return user, nil
}
