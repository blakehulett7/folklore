package main

import (
	"bufio"
	"fmt"
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
	Run("clear")
	fmt.Println("Create an Account")
	prompt := bufio.NewScanner(os.Stdin)
	fmt.Print("\nEnter a username > ")
	prompt.Scan()
	username := prompt.Text()
	if !IsValidInput(username) {
		fmt.Println("';' character not allowed...")
		prompt.Scan()
		return
	}
	fmt.Println(username)
	fmt.Print("Enter a password > ")
	prompt.Scan()
	password := prompt.Text()
	fmt.Println(password)
	prompt.Scan()
}
