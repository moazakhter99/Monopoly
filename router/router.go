package router

type Router interface {
	HandleFunc(path string, f func(Request, chan []byte)) *Route
	router(req Request)
	Run()
	inMap(action string) (ok bool)
	readResponse() (msg []byte)
}

type Route struct {
	wrChan     chan []byte
	handlerMap map[string]func(Request, chan []byte)
}

func (r Route) Run() {
	select {}

}

func (r Route) HandleFunc(path string, f func(Request, chan []byte)) *Route {
	r.handlerMap[path] = f
	return &r
}

func (r Route) router(req Request) {
	f := r.handlerMap[req.Action]
	f(req, r.wrChan)
}

func (r Route) inMap(action string) (ok bool) {
	_, ok = r.handlerMap[action]
	return
}

func (r Route) readResponse() (msg []byte) {
	return <- r.wrChan
}

func NewRouter(bufferSize int) Route {
	return Route{
		wrChan:     make(chan []byte, bufferSize),
		handlerMap: make(map[string]func(Request, chan []byte)),
	}
}