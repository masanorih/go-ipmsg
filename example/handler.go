package main

import (
	goipmsg "github.com/masanorih/go-ipmsg"
)

var Userlist = make(map[string]*goipmsg.ClientData)

func RECEIVE_BR_ENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	Userlist[cd.Key()] = cd
	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), goipmsg.ANSENTRY)
	return nil
}

func RECEIVE_ANSENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	Userlist[cd.Key()] = cd
	return nil
}
