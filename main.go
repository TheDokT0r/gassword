package main

import (
	"gli/cmd"
	"gli/internal/vault"
)

func main() {
	if !vault.VaultExists() {
		password := cmd.CreatePassword()
		vault.CreateVault(password)
	}

	cmd.ClearScreen()
	password := cmd.Login()

	cmd.ClearScreen()
	cmd.MainMenu(string(password), 0)
}
