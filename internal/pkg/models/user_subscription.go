//go:generate easyjson user_subscription.go
package models

//easyjson:json
type UserSubscription struct {
	TargetId int `json:"targetId"`
}
