package action

type Handler interface {
	Action(Action)
}

type HandlerFunc func(Action)

func (h HandlerFunc) Action(action Action) {
	h(action)
}
