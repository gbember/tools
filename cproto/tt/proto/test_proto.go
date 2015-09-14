package proto

const (
	CS_TEST = uint16(1001)
	P_TEST  = uint16(1002)
)

//消息注释
type Cs_test struct {
	//字段注释1
	ID       int32
	Name     string
	BList    []bool
	AddrList []string
	Tt       P_test
	Tts      []P_test
}

func (r *Cs_test) Read(p *Packet) error {
	value0, err := p.readInt32()
	if err != nil {
		return err
	}
	r.ID = value0
	value1, err := p.readString()
	if err != nil {
		return err
	}
	r.Name = value1
	value2, err := p.readArrayBool()
	if err != nil {
		return err
	}
	r.BList = value2
	value3, err := p.readArrayString()
	if err != nil {
		return err
	}
	r.AddrList = value3
	r.Tt = P_test{}
	err = r.Tt.Read(p)
	if err != nil {
		return err
	}
	len5, err := p.readUint16()
	if err != nil {
		return err
	}
	r.Tts = make([]P_test, 0, int(len5))
	for i := uint16(0); i < len5; i++ {
		v := P_test{}
		err = v.Read(p)
		if err != nil {
			return err
		}
		r.Tts = append(r.Tts, v)
	}
	return nil
}
func (r *Cs_test) WriteMsgID(p *Packet) {
	p.writeUint16(CS_TEST)
}
func (r *Cs_test) Write(p *Packet) error {
	p.writeInt32(r.ID)
	p.writeString(r.Name)
	p.writeArrayBool(r.BList)
	p.writeArrayString(r.AddrList)
	err := r.Tt.Write(p)
	if err != nil {
		return err
	}
	p.writeUint16(uint16(len(r.Tts)))
	for i := 0; i < len(r.Tts); i++ {
		err := r.Tts[i].Write(p)
		if err != nil {
			return err
		}
	}
	return nil
}

type P_test struct {
	Size uint16
	A    int32
	B    int32
	C    int32
	D    int32
	E    int32
	F    int32
	I    int32
	J    int32
}

func (r *P_test) Read(p *Packet) error {
	value0, err := p.readUint16()
	if err != nil {
		return err
	}
	r.Size = value0
	value1, err := p.readInt32()
	if err != nil {
		return err
	}
	r.A = value1
	value2, err := p.readInt32()
	if err != nil {
		return err
	}
	r.B = value2
	value3, err := p.readInt32()
	if err != nil {
		return err
	}
	r.C = value3
	value4, err := p.readInt32()
	if err != nil {
		return err
	}
	r.D = value4
	value5, err := p.readInt32()
	if err != nil {
		return err
	}
	r.E = value5
	value6, err := p.readInt32()
	if err != nil {
		return err
	}
	r.F = value6
	value7, err := p.readInt32()
	if err != nil {
		return err
	}
	r.I = value7
	value8, err := p.readInt32()
	if err != nil {
		return err
	}
	r.J = value8
	return nil
}
func (r *P_test) WriteMsgID(p *Packet) {
	p.writeUint16(P_TEST)
}
func (r *P_test) Write(p *Packet) error {
	p.writeUint16(r.Size)
	p.writeInt32(r.A)
	p.writeInt32(r.B)
	p.writeInt32(r.C)
	p.writeInt32(r.D)
	p.writeInt32(r.E)
	p.writeInt32(r.F)
	p.writeInt32(r.I)
	p.writeInt32(r.J)
	return nil
}
