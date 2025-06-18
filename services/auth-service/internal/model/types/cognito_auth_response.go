package types

type AuthResult struct {
	AccessToken         string            `json:"access_token,omitempty"`
	IdToken             string            `json:"id_token,omitempty"`
	RefreshToken        string            `json:"refresh_token,omitempty"`
	TokenType           string            `json:"token_type,omitempty"`
	ExpiresIn           int32             `json:"expires_in,omitempty"`
	ChallengeName       string            `json:"challenge_name,omitempty"`
	Session             string            `json:"session,omitempty"`
	ChallengeParameters map[string]string `json:"challenge_parameters,omitempty"`
}
