package ipmsg

import "fmt"

type Message int

const (
	// COMMAND
	NOOPERATION     Message = 0x00000000
	BR_ENTRY        Message = 0x00000001
	BR_EXIT         Message = 0x00000002
	ANSENTRY        Message = 0x00000003
	BR_ABSENCE      Message = 0x00000004
	BR_ISGETLIST    Message = 0x00000010
	OKGETLIST       Message = 0x00000011
	GETLIST         Message = 0x00000012
	ANSLIST         Message = 0x00000013
	BR_ISGETLIST2   Message = 0x00000018
	SENDMSG         Message = 0x00000020
	RECVMSG         Message = 0x00000021
	READMSG         Message = 0x00000030
	DELMSG          Message = 0x00000031
	ANSREADMSG      Message = 0x00000032
	GETINFO         Message = 0x00000040
	SENDINFO        Message = 0x00000041
	GETABSENCEINFO  Message = 0x00000050
	SENDABSENCEINFO Message = 0x00000051
	GETFILEDAT      Message = 0x00000060
	RELEASEFIL      Message = 0x00000061
	GETDIRFILE      Message = 0x00000062
	GETPUBKEY       Message = 0x00000072
	ANSPUBKEY       Message = 0x00000073
	// MODE
	MODE Message = 0x000000ff
	// OPTION
	//ABSENCE    Message = 0x00000100
	//SERVER     Message = 0x00000200
	//DIALUP     Message = 0x00010000
	SENDCHECK  Message = 0x00000100
	SECRET     Message = 0x00000200
	BROADCAST  Message = 0x00000400
	MULTICAST  Message = 0x00000800
	NOPOPUP    Message = 0x00001000
	AUTORET    Message = 0x00002000
	RETRY      Message = 0x00004000
	PASSWORD   Message = 0x00008000
	NOLOG      Message = 0x00020000
	NEWMUTI    Message = 0x00040000
	NOADDLIST  Message = 0x00080000
	READCHECK  Message = 0x00100000
	FILEATTACH Message = 0x00200000
	ENCRYPT    Message = 0x00400000
)

func (msg Message) Mode() Message {
	if 0 == msg {
		return 0
	}
	return msg & MODE
}

func (msg Message) ModeName() string {
	mode := msg.Mode()
	str := fmt.Sprint(mode)
	return str
}

func (msg Message) Get(cmd Message) bool {
	return msg&cmd != 0
}

func (msg *Message) SetOpt(cmd Message) {
	*msg |= cmd
}
