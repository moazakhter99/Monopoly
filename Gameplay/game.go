package gameplay

type Game interface {
	Validate(reqMsg []byte) (payload any, err error)
	Play(payload any, reqParam map[string]string) (targetMap map[string]any, err error)
	// StateManagement(reqType, state string) (err error)
	Response(targetMap map[string]any, reqParam map[string]string, readChan chan []byte) (err error)
}
