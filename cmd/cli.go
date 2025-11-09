package cmd

import (
	"fmt"
	"gli/internal/vault"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"golang.design/x/clipboard"
	"golang.org/x/term"
)

const (
	Up = iota
	Down
	Add
	Edit
	Remove
	View
	Copy
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

func keyboardActionDetection() int {
	keyCode := None

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Up {
			keyCode = Up
		} else if key.Code == keys.Down {
			keyCode = Down
		} else if strings.ToLower(key.String()) == "a" {
			keyCode = Add
		} else if strings.ToLower(key.String()) == "r" {
			keyCode = Remove
		} else if strings.ToLower(key.String()) == "c" {
			keyCode = Copy
		} else if strings.ToLower(key.String()) == "e" {
			keyCode = Edit
		} else if key.Code == keys.CtrlC {
			os.Exit(0)
		} else {
			return false, nil
		}

		return true, nil
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

	fmt.Println(generateOptionsMenu([]string{"^Up", "vDown", "Add", "Edit", "Remove", "View", "Copy Password"}))
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

func copyPasswordToClipboard(item vault.VaultItem) {
	err := clipboard.Init()
	if err != nil {
		log.Fatal(err)
	}

	clipboard.Write(clipboard.FmtText, []byte(item.Password))
}
