package gameplay

import (
)





type Game interface {
	Validate(reqMsg []byte) (payload any, err error)
	Play(payload any) (resp any, err error)
	// StateManagement(reqType, state string) (err error)
	// Response(respMsg any) (response []byte, err error)
}


