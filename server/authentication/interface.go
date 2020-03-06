package authentication

type Authenticator interface {
	GenerateToken() string
	Authenticate(candidateToken string, candidateUsername string) bool
	IsTokenValid(candidateToken string) bool
	DeleteToken(token string, username string) error
}
