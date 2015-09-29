package ipmsg

import "fmt"

type EvFunc func(cd *ClientData, ipmsg *IPMSG) error
type EventHandler struct {
	String   string
	Handlers map[Command]EvFunc
	Debug    bool
}

func NewEventHandler() *EventHandler {
	ev := &EventHandler{
		Handlers: make(map[Command]EvFunc),
	}
	return ev
}

func (ev *EventHandler) Regist(cmd Command, evfunc EvFunc) {
	handlers := ev.Handlers
	handlers[cmd] = evfunc
}

func (ev *EventHandler) Run(cd *ClientData, ipmsg *IPMSG) error {
	if ev.Debug {
		ev.RunDebug(cd)
	}
	cmd := cd.Command.Mode()
	evfunc := ev.Handlers[cmd]
	if evfunc == nil {
		// just do nothing when handler is undefined
		return nil
	} else {
		return (evfunc(cd, ipmsg))
	}
}

func (ev *EventHandler) RunDebug(cd *ClientData) {
	cmdstr := cd.Command.Mode().String()
	fmt.Println("EventHandler.RunDebug cmdstr=", cmdstr)
	fmt.Println("EventHandler.RunDebug key=", cd.Key())
}

//func (ev EventHandler) Run(cd *ClientData, ipmsg *IPMSG) error {
//	cmdstr := cd.Command.String()
//	v := reflect.ValueOf(&ev)
//	method := v.MethodByName(cmdstr)
//	if !method.IsValid() {
//		err := fmt.Errorf("method for Command(%v) not defined", cmdstr)
//		return err
//	}
//	in := []reflect.Value{reflect.ValueOf(cd), reflect.ValueOf(ipmsg)}
//	err := method.Call(in)[0].Interface()
//	// XXX only works if you sure about the return value is always type(error)
//	if err == nil {
//		return nil
//	}
//	return err.(error)
//	//reflect.ValueOf(&ev).MethodByName(cmdstr).Call(in)
//}
//
//func (ev *EventHandler) BR_ENTRY(cd *ClientData, ipmsg *IPMSG) error {
//	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), ANSENTRY)
//	return nil
//}
