package scrape

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	Username string `json:"username"`
}

type LogOutRequest struct {
	Username string `json:"username"`
}
