package ipmsg

import "testing"

func TestMode(t *testing.T) {
	entry := BR_ENTRY
	if BR_ENTRY != entry.Mode() {
		t.Errorf("entry.Mode() contains not only mode part")
	}
	entry.SetOpt(SECRET)
	//fmt.Printf("entry        %#08x\n", entry)
	//fmt.Printf("entry.Mode() %#08x\n", entry.Mode())
	if BR_ENTRY != entry.Mode() {
		t.Errorf("entry.Mode() contains not only mode part")
	}
}

func TestModeName(t *testing.T) {
	entry := BR_ENTRY
	name := entry.ModeName()
	if "BR_ENTRY" != name {
		t.Errorf("entry.ModeName() is not 'BR_ENTRY'")
	}
}

func TestGet(t *testing.T) {
	entry := BR_ENTRY
	if entry.Get(BR_EXIT) {
		t.Errorf("BR_ENTRY contains BR_EXIT")
	}
	if !entry.Get(BR_ENTRY) {
		t.Errorf("BR_ENTRY does not contains BR_ENTRY")
	}
}

func TestSetOpt(t *testing.T) {
	entry := BR_ENTRY
	if entry.Get(SECRET) {
		t.Errorf("BR_ENTRY does contains SECRET before set it up")
	}
	//fmt.Printf("entry  %#08x\n", entry)
	//fmt.Printf("SECRET %#08x\n", SECRET)
	entry.SetOpt(SECRET)
	//fmt.Printf("entry  %#08x\n", entry)
	// output as hex digit like 'entry  0x00000001'
	//                          'SECRET 0x00000200'
	//                          'entry  0x00000201'
	if !entry.Get(BR_ENTRY) {
		t.Errorf("BR_ENTRY does not contains BR_ENTRY")
	}
	if !entry.Get(SECRET) {
		t.Errorf("BR_ENTRY does not contains SECRET after set it up")
	}
}
