package auth

type User struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	SubscriptionType string `json:"subscription_type"`
	Role             string `json:"role"`
}
