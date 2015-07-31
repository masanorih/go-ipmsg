package main

import (
	goipmsg "github.com/masanorih/go-ipmsg"
)

func RECEIVE_BR_ENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) {
	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), goipmsg.ANSENTRY)
	//return nil
}
