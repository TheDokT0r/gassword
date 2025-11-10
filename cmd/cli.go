package cmd

import (
	"fmt"
	"gli/internal/vault"
	"log"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/term"
)

func CreatePassword() string {
	fmt.Println("It seems like you still don't have a master password. Please insert one now")

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	if len(password) == 0 {
		ClearScreen()
		fmt.Println("Cannot use empty password. Please try again")
		CreatePassword()
	}

	return string(password)
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func Login() []byte {
	fmt.Print("Welcome to Gassword. Please enter your password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}

	_, err = vault.ReadVault(string(password))

	if err != nil {
		ClearScreen()
		fmt.Println("Incorrect password. Please try again")
		password = Login()
	}

	return password
}

func MainMenu(password string, index int) {
	fullVault, err := vault.ReadVault(password)
	if err != nil {
		log.Fatal(err)
	}

	printMainMenu(fullVault, index)
	keyCode := keyboardActionDetection()
	switch keyCode {
	case Up:
		index--
	case Down:
		index++
	case Add:
		ClearScreen()
		AddMenu(password)
	case Remove:
		vault.RemoveItemFromVault(password, index)
	case Copy:
		copyPasswordToClipboard(fullVault[index])
	case Edit:
		ClearScreen()
		EditMenu(fullVault, index, password)
	}

	if index < 0 {
		index = len(fullVault) - 1
	} else if index >= len(fullVault) {
		index = 0
	}

	ClearScreen()
	MainMenu(password, index)
}

func AddMenu(masterPass string) {
	var name string
	var email string

	fmt.Print("Username: ")
	fmt.Scan(&name)

	fmt.Print("Email: ")
	fmt.Scan(&email)

	fmt.Print("Password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {

		log.Fatal(err)
	}

	vaultItem := vault.NewVaultItem(name, email, string(password))

	fullVault, err := vault.ReadVault(masterPass)
	if err != nil {
		log.Fatal(err)
	}

	fullVault = append(fullVault, vaultItem)
	vault.WriteVault(fullVault, string(password))
}

func EditMenu(fullVault []vault.VaultItem, index int, password string) {
	vaultItem := fullVault[index]

	fmt.Print("Username: ")
	fmt.Scan(&vaultItem.Name)
	fmt.Print("Email: ")
	fmt.Scan(&vaultItem.Email)

	fullVault[index] = vaultItem
	vault.WriteVault(fullVault, password)
}
