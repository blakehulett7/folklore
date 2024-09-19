package main

import (
	"fmt"
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

var dashboardOptions = []goToYourMenu.MenuOption{
	{
		Name:    "Review Your Languages",
		Command: func() { fmt.Println("Not implemented") },
	},
	{
		Name:    "Add New Language",
		Command: func() { fmt.Println("Not implemented") },
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

func Logout() {
	exec.Command("rm", ".env").Run()
}
