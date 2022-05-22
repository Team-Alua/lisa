package reservation

type RequestType int64

const (
	Add RequestType = iota
	Remove
)

type Request struct {
	Type RequestType
	Value *Slot
}

type ResponseType int64

const (
	Info ResponseType = iota
	Error
	Warning
	Ready
	NotReady
)

type Response struct {
	Type ResponseType
	Msg string
}

type Slot struct {
	Id string
	Out chan<- Response
}

type System struct {
	slots [16]*Slot
	queue []*Slot
}

