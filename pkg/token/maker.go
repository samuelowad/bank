package token

import "time"

//Maker is a interface fo managing token creation
type Maker interface {
	//	CreateToken creates a new token for the given user and duration
	CreateToken(username string, duration time.Duration) (string, error)
	//	VerifyToken verifies the token is valid
	VerifyToken(token string) (*PayLoad, error)
}
