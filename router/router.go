package router

type Router interface {
	HandleFunc(path string, f func([]byte, chan []byte)) *Route
	router(action string, msg []byte)
	Run()
	inMap(action string) (ok bool)
	readResponse() (msg []byte)
}

type Route struct {
	wrChan     chan []byte
	handlerMap map[string]func([]byte, chan []byte)
}

func (r Route) Run() {
	select {}

}

func (r Route) HandleFunc(path string, f func([]byte, chan []byte)) *Route {
	r.handlerMap[path] = f
	return &r
}

func (r Route) router(action string, msg []byte) {
	f := r.handlerMap[action]
	f(msg, r.wrChan)
}

func (r Route) inMap(action string) (ok bool) {
	_, ok = r.handlerMap[action]
	return
}

func (r Route) readResponse() (msg []byte) {
	return <- r.wrChan
}

func NewRouter() Route {
	return Route{
		wrChan:     make(chan []byte, 3),
		handlerMap: make(map[string]func([]byte, chan []byte)),
	}
}