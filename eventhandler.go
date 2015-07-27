package ipmsg

import (
	"fmt"
	"reflect"
)

type EventHandler struct {
	String string
}

func (ev *EventHandler) Debug() {
}

func (ev EventHandler) Run(cd *ClientData, ipmsg *IPMSG) error {
	cmdstr := cd.Command.String()
	v := reflect.ValueOf(&ev)
	method := v.MethodByName(cmdstr)
	if !method.IsValid() {
		err := fmt.Errorf("method for cmdstr(%v) not defined", cmdstr)
		return err
	}
	//pp.Println("method=", method)
	in := []reflect.Value{reflect.ValueOf(cd), reflect.ValueOf(ipmsg)}
	err := method.Call(in)[0].Interface()
	// XXX only works if you are sure about the return value is always type(error)
	if err == nil {
		return nil
	}
	return err.(error)
	//reflect.ValueOf(&ev).MethodByName(cmdstr).Call(in)
}

func (ev *EventHandler) BR_ENTRY(cd *ClientData, ipmsg *IPMSG) error {
	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), ANSENTRY)
	//pp.Println("BR_ENTRY in=", cd)
	return nil
}

//func (ev *EventHandler) BR_EXIT(cd *ClientData) error {
//	err := errors.New("always error")
//	return err
//}
