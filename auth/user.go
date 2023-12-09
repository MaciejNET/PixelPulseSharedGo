package auth

type User struct {
	Id               string `json:"id"`
	Email            string `json:"email"`
	SubscriptionType string `json:"subscription_type"`
	Role             string `json:"role"`
}
