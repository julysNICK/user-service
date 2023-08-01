package token

type Maker interface {
	CreateToken(username string) (string, *Payload, error)

	VerifyToken(token string) (*Payload, error)
}
