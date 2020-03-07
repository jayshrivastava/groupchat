package authentication

type Authenticator interface {
	GenerateToken() string
	AuthenticateToken(candidateToken string, candidateUsername string) bool
	IsTokenValid(candidateToken string) bool
	DeleteToken(token string, username string) error
}
