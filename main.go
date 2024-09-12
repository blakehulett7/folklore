package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

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

func CreateAccount() {
	Run("clear")
	fmt.Println("Create an Account")
	prompt := bufio.NewScanner(os.Stdin)
	fmt.Print("\nEnter a username > ")
	prompt.Scan()
	input := prompt.Text()
	fmt.Println(input)
	prompt.Scan()
}

func Run(program string, args ...string) {
	command := exec.Command(program, args...)
	command.Stdout = os.Stdout
	command.Run()
}
