package ipmsg

import "fmt"

type Command int

const (
	// COMMAND
	NOOPERATION     Command = 0x00000000
	BR_ENTRY        Command = 0x00000001
	BR_EXIT         Command = 0x00000002
	ANSENTRY        Command = 0x00000003
	BR_ABSENCE      Command = 0x00000004
	BR_ISGETLIST    Command = 0x00000010
	OKGETLIST       Command = 0x00000011
	GETLIST         Command = 0x00000012
	ANSLIST         Command = 0x00000013
	BR_ISGETLIST2   Command = 0x00000018
	SENDMSG         Command = 0x00000020
	RECVMSG         Command = 0x00000021
	READMSG         Command = 0x00000030
	DELMSG          Command = 0x00000031
	ANSREADMSG      Command = 0x00000032
	GETINFO         Command = 0x00000040
	SENDINFO        Command = 0x00000041
	GETABSENCEINFO  Command = 0x00000050
	SENDABSENCEINFO Command = 0x00000051
	GETFILEDAT      Command = 0x00000060
	RELEASEFIL      Command = 0x00000061
	GETDIRFILE      Command = 0x00000062
	GETPUBKEY       Command = 0x00000072
	ANSPUBKEY       Command = 0x00000073
	// MODE
	MODE Command = 0x000000ff
	// OPTION
	//ABSENCE    Command = 0x00000100
	//SERVER     Command = 0x00000200
	//DIALUP     Command = 0x00010000
	SENDCHECK  Command = 0x00000100
	SECRET     Command = 0x00000200
	BROADCAST  Command = 0x00000400
	MULTICAST  Command = 0x00000800
	NOPOPUP    Command = 0x00001000
	AUTORET    Command = 0x00002000
	RETRY      Command = 0x00004000
	PASSWORD   Command = 0x00008000
	NOLOG      Command = 0x00020000
	NEWMUTI    Command = 0x00040000
	NOADDLIST  Command = 0x00080000
	READCHECK  Command = 0x00100000
	FILEATTACH Command = 0x00200000
	ENCRYPT    Command = 0x00400000
)

func (msg Command) Mode() Command {
	if 0 == msg {
		return 0
	}
	return msg & MODE
}

func (msg Command) ModeName() string {
	mode := msg.Mode()
	str := fmt.Sprint(mode)
	return str
}

func (msg Command) Get(flg Command) bool {
	return msg&flg != 0
}

func (msg *Command) SetOpt(flg Command) {
	*msg |= flg
}
