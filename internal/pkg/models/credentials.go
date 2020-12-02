//go:generate easyjson credentials.go
package models

//easyjson:json
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
