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
