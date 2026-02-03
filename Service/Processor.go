package service


type Processor interface {
	Validate(data []byte) (req any, err error)
	ProcessMsg(req any) (resp []byte, err error)
}