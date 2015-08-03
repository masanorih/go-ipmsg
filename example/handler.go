package main

import (
	"strconv"

	goipmsg "github.com/masanorih/go-ipmsg"
)

var Users = make(map[string]*goipmsg.ClientData)
var Messages = []*goipmsg.ClientData{}

func RECEIVE_BR_ENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	Users[cd.Key()] = cd
	ipmsg.SendMSG(cd.Addr, ipmsg.Myinfo(), goipmsg.ANSENTRY)
	return nil
}

func RECEIVE_ANSENTRY(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	Users[cd.Key()] = cd
	return nil
}

func RECEIVE_SENDMSG(cd *goipmsg.ClientData, ipmsg *goipmsg.IPMSG) error {
	sl := append(Messages, cd)
	Messages = sl

	cmd := cd.Command
	if cmd.Get(goipmsg.SENDCHECK) {
		num := cd.PacketNum
		err := ipmsg.SendMSG(cd.Addr, strconv.Itoa(num), goipmsg.RECVMSG)
		return err
	}
	return nil
}
