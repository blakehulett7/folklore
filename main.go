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
	"reflect"
	"slices"
	"strings"

	"github.com/blakehulett7/goToYourMenu"
	"github.com/joho/godotenv"
)

const hostUrl = "http://localhost:8080"
const hostVersion = "v1"
const maxPermissions = 0777

var languages = []string{"Italian", "Spanish", "French"}
var errBadToken = errors.New("bad token")

type User struct {
	Id              string   `json:"id"`
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	RefreshToken    string   `json:"refresh_token"`
	ListeningStreak string   `json:"listening_streak"`
	Languages       []string `json:"languages"`
}

func main() {
	for {
		godotenv.Load()
		token := os.Getenv("JWT")
		var user User
		var err error
		if token != "" {
			user, err = GetUser(token)
		}
		if err != nil {
			fmt.Println("couldn't get user, error:", err)
		}
		if !reflect.DeepEqual(user, User{}) {
			LaunchDashboard(user)
		}
		for {
			Run("clear")
			fmt.Println("Christ is King!")
			fmt.Println("\nWelcome to Folklore!")
			command := goToYourMenu.Menu(startMenuOptions)
			if command != "Login" {
				continue
			}
			if err != nil {
				fmt.Println("couldn't get user, error:", err)
			}
			break
		}
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
	os.Setenv("JWT", credentials.JWT)
	os.Setenv("REFRESH_TOKEN", credentials.RefreshToken)
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

func GetUser(token string) (User, error) {
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
		//user the refresh token and try once more
		return User{}, errBadToken
	}
	user := User{}
	err = json.NewDecoder(res.Body).Decode(&user)
	return user, nil
}

func LaunchDashboard(user User) {
	for {
		Run("clear")
		fmt.Println("Christ is King!")
		fmt.Println("\nWelcome to Folklore,", user.Username)
		fmt.Println("Listening Streak:", user.ListeningStreak)
		fmt.Println(user.Languages)
		fmt.Println("\nWhat would you like to do?")
		var dashboardOptions = []goToYourMenu.MenuOption{
			{
				Name:    "Review Your Languages",
				Command: func() { fmt.Println("Not implemented") },
			},
			{
				Name:    "Add New Language",
				Command: user.AddLanguage,
			},
			{
				Name:    "Remove a Language",
				Command: func() { fmt.Println("Not implemented") },
			},
			{
				Name:    "Logout",
				Command: Logout,
			},
			{
				Name:    "Exit",
				Command: func() { os.Exit(0) },
			},
		}
		command := goToYourMenu.Menu(dashboardOptions)
		if command == "Logout" {
			break
		}
		bufio.NewScanner(os.Stdin).Scan()
	}
}

func (user *User) AddLanguage() {
	Run("clear")
	fmt.Println("Christ is King!")
	fmt.Println("\nWelcome to Folklore,", user.Username)
	fmt.Println("Listening Streak:", user.ListeningStreak)
	fmt.Println(user.Languages)
	fmt.Println("\nChoose a language to add")
	options := []goToYourMenu.MenuOption{}
	for _, language := range languages {
		if !slices.Contains(user.Languages, language) {
			options = append(options, goToYourMenu.MenuOption{
				Name:    language,
				Command: func() {},
			})
		}
	}
	options = append(options, goToYourMenu.MenuOption{Name: "Go Back", Command: func() {}})
	languageToAdd := goToYourMenu.Menu(options)
	if languageToAdd == "Go Back" {
		return
	}
	_, err := SendLanguageRequest(languageToAdd)
	if err != nil {
		fmt.Println("Couldn't add language, error:", err)
		return
	}
	user.Languages = append(user.Languages, languageToAdd)
}

func SendLanguageRequest(languagetoAdd string) (User, error) {
	token := os.Getenv("JWT")
	payloadStruct := struct {
		Name string `json:"name"`
	}{languagetoAdd}
	payload, err := json.Marshal(payloadStruct)
	if err != nil {
		fmt.Println("couldn't marshal json:", err)
		return User{}, err
	}
	requestURL := fmt.Sprintf("%v/%v/users_languages", hostUrl, hostVersion)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", token)
	if err != nil {
		fmt.Println("couldn't create http request:", err)
		return User{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("couldn't execute http request:", err)
		return User{}, err
	}
	if res.StatusCode == 401 {
		fmt.Println("Error authenticating")
		return User{}, errBadToken
	}
	user := User{}
	json.NewDecoder(res.Body).Decode(&user)
	return user, nil
}
