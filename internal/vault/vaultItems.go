package vault

import "github.com/google/uuid"

type VaultItem struct {
	uid      string
	Name     string
	Email    string
	Password string
}

func NewVaultItem(name string, email string, password string) VaultItem {
	return VaultItem{uid: uuid.NewString(), Name: name, Email: email, Password: password}
}
