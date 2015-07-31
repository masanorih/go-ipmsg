package ipmsg

import "testing"

func TestEventHander(t *testing.T) {
	conf := NewIPMSGConf()
	ipmsg, err := NewIPMSG(conf)
	if err != nil {
		t.Errorf("ipmsg error is not nil '%v'", err)
	}
	defer ipmsg.Close()

	addr, err := ipmsg.UDPAddr()
	if err != nil {
		t.Errorf("failed to resolve to UDP '%v'", err)
	}

	ev := EventHandler{}
	clientdata := ipmsg.BuildData(addr, "hogehoge", BR_ENTRY)
	err = ev.Run(clientdata, ipmsg)
	if err != nil {
		t.Errorf("ev.Run(BR_ENTRY) failed with '%v'", err)
	}

	clientdata = ipmsg.BuildData(addr, "hogehoge", BR_EXIT)
	err = ev.Run(clientdata, ipmsg)
	if err == nil {
		t.Errorf("ev.Run(BR_EXIT) do not fail")
	}
}

func TestAddEventHandler(t *testing.T) {
	conf := NewIPMSGConf()
	ipmsg, err := NewIPMSG(conf)
	if err != nil {
		t.Errorf("ipmsg error is not nil '%v'", err)
	}
	defer ipmsg.Close()

	ev := EventHandler{String: "TestAddEventHandler"}
	ipmsg.AddEventHandler(ev)
	//pp.Println("ipmsg.Handlers=", ipmsg.Handlers)

	//clientdata := ipmsg.BuildData(addr, "hogehoge", BR_ENTRY)
	addr, err := ipmsg.UDPAddr()
	if err != nil {
		t.Errorf("ipmsg.UDPAddr() has err '%v'", err)
	}
	err = ipmsg.SendMSG(addr, "TestAddEventHandler", BR_ENTRY)
	if err != nil {
		t.Errorf("ipmsg.SendMSG() has err '%v'", err)
	}

	recv, err := ipmsg.RecvMSG()
	if err != nil {
		t.Errorf("ipmsg.RecvMSG() has err '%v'", err)
	}
	if recv == nil {
		t.Errorf("recv is nil")
	}
	//pp.Println("recv=", recv)
	//if recv != "TestAddEventHandler" {
	//	t.Errorf("received message is not what I sent")
	//}
}
