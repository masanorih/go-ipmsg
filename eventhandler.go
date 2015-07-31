package ipmsg

import "fmt"

type EventHandler struct {
	String   string
	Handlers map[Command]func(cd *ClientData, ipmsg *IPMSG)
}

func NewEventHandler() *EventHandler {
	ev := &EventHandler{
		Handlers: make(map[Command]func(cd *ClientData, ipmsg *IPMSG)),
	}
	return ev
}

func (ev *EventHandler) Regist(cmd Command, handler func(cd *ClientData, ipmsg *IPMSG)) {
	//cmdstr := cmd.String()
	handlers := ev.Handlers
	handlers[cmd] = handler
}

func (ev *EventHandler) Run(cd *ClientData, ipmsg *IPMSG) error {
	cmd := cd.Command
	handler := ev.Handlers[cmd]
	if handler == nil {
		err := fmt.Errorf("func for Command(%v) not defined", cmd.String())
		return err
	} else {
		handler(cd, ipmsg)
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
