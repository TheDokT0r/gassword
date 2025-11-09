package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"gli/internal/encryption"
	"log"
	"os"
	"path"
	"sort"
)

func getVaultLocation() string {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	return path.Join(homeDir, ".gassword", ".vault")
}

func VaultExists() bool {
	vaultPath := getVaultLocation()
	_, err := os.Stat(vaultPath)
	return !errors.Is(err, os.ErrNotExist)
}

func CreateVault(password string) {
	vaultPath := getVaultLocation()

	err := os.Mkdir(path.Dir(vaultPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(vaultPath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var emptyVault [0]VaultItem
	encryptedVault, err := encryption.Encrypt([]byte(password), vaultToBytes(emptyVault[:]))

	if err != nil {
		log.Fatal(err)
	}

	f.Write(encryptedVault)

	f.Close()
}

func vaultToBytes(vault []VaultItem) []byte {
	b, err := json.Marshal(vault)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func bytesToVault(vaultBytes []byte) []VaultItem {
	var vault []VaultItem
	json.Unmarshal(vaultBytes, &vault)

	return vault
}

func ReadVault(password string) ([]VaultItem, error) {
	data, err := os.ReadFile(getVaultLocation())

	if err != nil {
		log.Fatal(err)
	}

	decryptData, err := encryption.Decrypt([]byte(password), data)

	if err != nil {
		return nil, err
	}

	vault := bytesToVault(decryptData)
	return vault, nil
}

func WriteVault(vault []VaultItem, password string) {
	// vault = sortVault(vault)

	vaultPath := getVaultLocation()

	b := vaultToBytes(vault)
	encryptedData, err := encryption.Encrypt([]byte(password), b)

	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
	}

	err = os.WriteFile(vaultPath, encryptedData, 0644)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
	}
}

func RemoveItemFromVault(password string, index int) {
	vaultItems, err := ReadVault(password)
	if err != nil {
		log.Fatal(err)
	}

	vaultCopy := vaultItems[:index]
	vaultCopy = append(vaultCopy, vaultItems[index+1:]...)

	WriteVault(vaultCopy, password)
}

func sortVault(vaultItems []VaultItem) []VaultItem {
	sort.Slice(vaultItems, func(i, j int) bool {
		return vaultItems[i].Name < vaultItems[j].Name
	})

	return vaultItems
}
