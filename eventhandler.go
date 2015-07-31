package ipmsg

import "fmt"

type EvFunc func(cd *ClientData, ipmsg *IPMSG) error
type EventHandler struct {
	String   string
	Handlers map[Command]EvFunc
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
	cmd := cd.Command
	evfunc := ev.Handlers[cmd]
	if evfunc == nil {
		err := fmt.Errorf("func for Command(%v) not defined", cmd.String())
		return err
	} else {
		return (evfunc(cd, ipmsg))
	}
	return nil
}

func (ev *EventHandler) Debug(cd *ClientData) {
	cmdstr := cd.Command.String()
	fmt.Println("EventHandler.Debug cmdstr=", cmdstr)
	fmt.Println("EventHandler.Debug key=", cd.Key())
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
