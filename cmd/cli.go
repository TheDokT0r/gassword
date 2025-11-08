package cmd

import (
	"fmt"
	"gli/internal/vault"
	"log"
	"os"
	"os/exec"
	"runtime"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/fatih/color"
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

const (
	Up = iota
	Down
	Add
	Edit
	Remove
	View
	None
)

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
			c := color.New(color.FgBlue)
			c.Print(">> ")
		}
		fmt.Println(string(index) + ". " + vault.Name)
	}

	if len(vault) == 0 {
		fmt.Println("No items in vault")
	}

	fmt.Println("^Up ||vDown || Add || Edit || Remove || View ")
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
