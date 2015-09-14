package proto

import (
	"errors"
	"fmt"
)

type ProtoMsger interface {
	Read(*Packet) error
	WriteMsgID(*Packet)
	Write(*Packet) error
}

func ParseProto(bin []byte) (msgID uint16, msg interface{}, err error) {
	p := NewReader(bin)
	msgID, err = p.readUint16()
	if err != nil {
		return
	}
	mid := msgID / 100
	switch mid {
	case 100:
		switch msgID {
		case SC_ACCOUNT_KICK:
			v := &Sc_account_kick{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case CS_ACCOUNT_HEART:
			v := &Cs_account_heart{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_HEART:
			v := &Sc_account_heart{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case CS_ACCOUNT_LOGOUT:
			v := &Cs_account_logout{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case CS_ACCOUNT_LOGIN:
			v := &Cs_account_login{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_LOGIN:
			v := &Sc_account_login{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case CS_ACCOUNT_CREATE:
			v := &Cs_account_create{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_CREATE:
			v := &Sc_account_create{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		default:
			err = errors.New(fmt.Sprintf("error: invalid msgID %d", msgID))
		}
	case 10:
		switch msgID {
		case CS_TEST:
			v := &Cs_test{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case P_TEST:
			v := &P_test{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		default:
			err = errors.New(fmt.Sprintf("error: invalid msgID %d", msgID))
		}
	default:
		err = errors.New(fmt.Sprintf("error: invalid msgID %d", msgID))
	}
	return
}
