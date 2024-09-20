package main

import (
	"os"
	"os/exec"

	"github.com/blakehulett7/goToYourMenu"
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

var reviewLanguageOptions = []goToYourMenu.MenuOption{
	{
		Name:    "Listen to some Folklore",
		Command: func() {},
	},
	{
		Name:    "Review top 100 words",
		Command: func() {},
	},
	{
		Name:    "Go Back",
		Command: func() {},
	},
}

func Logout() {
	exec.Command("rm", ".env").Run()
}
