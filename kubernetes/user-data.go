package kubernetes

// UserData holds information about an authorized data
type UserData struct {
	Name string
	ClientID string
	ClientSecret string
	IDToken string
	RefreshToken string
	IssuerURL string
}
