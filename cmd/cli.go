package cmd

import (
	"fmt"
	"gli/internal/vault"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"golang.org/x/term"
)

const (
	Up = iota
	Down
	Add
	Edit
	Remove
	View
	None
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Bold  = "\033[1m"
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
		log.Fatal(fullVault)
	}

	printMainMenu(fullVault, index)
	keyCode := keyboardActionDetection()
	switch keyCode {
	case Up:
		index++
	case Down:
		index--
	case Add:
		ClearScreen()
		AddMenu(password)
	}

	ClearScreen()
	MainMenu(password, index)
}

func keyboardActionDetection() int {
	keyCode := None

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Up {
			keyCode = Up
		} else if key.Code == keys.Down {
			keyCode = Down
		} else if key.String() == "a" {
			keyCode = Add
		} else if key.Code == keys.CtrlC {
			os.Exit(0)
		} else {
			return true, nil
		}

		return false, nil
	})

	return keyCode
}

func printMainMenu(vault []vault.VaultItem, count int) {
	for index, vault := range vault {
		if index == count {
			fmt.Print(Blue + ">> " + Reset)
		}
		fmt.Println(strconv.Itoa(index) + ". " + vault.Name)
	}

	if len(vault) == 0 {
		fmt.Println("No items in vault")
	}

	fmt.Println(generateOptionsMenu([]string{"^Up", "vDown", "Add", "Edit", "Remove", "View"}))
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
		vault, err := vault.ReadVault(string(password))
		if err != nil {
			log.Fatal(vault)
		}
		log.Fatal(err)
	}

	fullVault = append(fullVault, vaultItem)
	vault.WriteVault(fullVault, string(password))
}

func generateOptionsMenu(options []string) string {
	menu := ""

	for index, option := range options {
		menu += Blue + string(option[0]) + Reset + option[1:]

		if index != len(options)-1 {
			menu += " || "
		}
	}

	return menu
}
