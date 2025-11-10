package cmd

import (
	"fmt"
	"gli/internal/vault"
	"log"
	"os"
	"strconv"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"golang.design/x/clipboard"
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
	fmt.Println("Password copied to clipboard")
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
