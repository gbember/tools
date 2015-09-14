package proto

const (
	CS_ACCOUNT_LOGIN  = uint16(10001)
	SC_ACCOUNT_LOGIN  = uint16(10002)
	CS_ACCOUNT_CREATE = uint16(10003)
	SC_ACCOUNT_CREATE = uint16(10004)
	SC_ACCOUNT_KICK   = uint16(10005)
	CS_ACCOUNT_HEART  = uint16(10006)
	SC_ACCOUNT_HEART  = uint16(10007)
	CS_ACCOUNT_LOGOUT = uint16(10008)
)

//请求登录
type Cs_account_login struct {
	////平台用户ID
	UserID int32
	////unix时间戳
	UnixTime int32
	////平台用户账号
	AccountName string
	////验证密钥
	Ticket string
	////服务器ID
	ServerID int16
}

func (r *Cs_account_login) Read(p *Packet) error {
	value0, err := p.readInt32()
	if err != nil {
		return err
	}
	r.UserID = value0
	value1, err := p.readInt32()
	if err != nil {
		return err
	}
	r.UnixTime = value1
	value2, err := p.readString()
	if err != nil {
		return err
	}
	r.AccountName = value2
	value3, err := p.readString()
	if err != nil {
		return err
	}
	r.Ticket = value3
	value4, err := p.readInt16()
	if err != nil {
		return err
	}
	r.ServerID = value4
	return nil
}
func (r *Cs_account_login) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_LOGIN)
}
func (r *Cs_account_login) Write(p *Packet) error {
	p.writeInt32(r.UserID)
	p.writeInt32(r.UnixTime)
	p.writeString(r.AccountName)
	p.writeString(r.Ticket)
	p.writeInt16(r.ServerID)
	return nil
}

type Sc_account_login struct {
	////
	Result int8
	// 1 => 成功
	// 2 => key错误
	// 3 => key超时
	// 4 => 已经登录
	// 5 => 登录的服务器id不符
	// 6 => 服务器人生已满
	////true=已经创建主公，false=未创建主公
	IsCreate bool
}

func (r *Sc_account_login) Read(p *Packet) error {
	value0, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Result = value0
	value1, err := p.readBool()
	if err != nil {
		return err
	}
	r.IsCreate = value1
	return nil
}
func (r *Sc_account_login) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_LOGIN)
}
func (r *Sc_account_login) Write(p *Packet) error {
	p.writeInt8(r.Result)
	p.writeBool(r.IsCreate)
	return nil
}

//创建角色
type Cs_account_create struct {
	////名字
	RoleName string
	////模型
	ModelID int8
}

func (r *Cs_account_create) Read(p *Packet) error {
	value0, err := p.readString()
	if err != nil {
		return err
	}
	r.RoleName = value0
	value1, err := p.readInt8()
	if err != nil {
		return err
	}
	r.ModelID = value1
	return nil
}
func (r *Cs_account_create) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_CREATE)
}
func (r *Cs_account_create) Write(p *Packet) error {
	p.writeString(r.RoleName)
	p.writeInt8(r.ModelID)
	return nil
}

type Sc_account_create struct {
	////结果
	Result int8
}

func (r *Sc_account_create) Read(p *Packet) error {
	value0, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Result = value0
	return nil
}
func (r *Sc_account_create) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_CREATE)
}
func (r *Sc_account_create) Write(p *Packet) error {
	p.writeInt8(r.Result)
	return nil
}

//强制下线通知
type Sc_account_kick struct {
	////原因
	Reason int8
}

func (r *Sc_account_kick) Read(p *Packet) error {
	value0, err := p.readInt8()
	if err != nil {
		return err
	}
	r.Reason = value0
	return nil
}
func (r *Sc_account_kick) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_KICK)
}
func (r *Sc_account_kick) Write(p *Packet) error {
	p.writeInt8(r.Reason)
	return nil
}

//心跳包
type Cs_account_heart struct {
}

func (r *Cs_account_heart) Read(p *Packet) error {
	return nil
}
func (r *Cs_account_heart) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_HEART)
}
func (r *Cs_account_heart) Write(p *Packet) error {
	return nil
}

type Sc_account_heart struct {
	////当前服务器时间
	UnixTime int32
}

func (r *Sc_account_heart) Read(p *Packet) error {
	value0, err := p.readInt32()
	if err != nil {
		return err
	}
	r.UnixTime = value0
	return nil
}
func (r *Sc_account_heart) WriteMsgID(p *Packet) {
	p.writeUint16(SC_ACCOUNT_HEART)
}
func (r *Sc_account_heart) Write(p *Packet) error {
	p.writeInt32(r.UnixTime)
	return nil
}

//注销登录
type Cs_account_logout struct {
}

func (r *Cs_account_logout) Read(p *Packet) error {
	return nil
}
func (r *Cs_account_logout) WriteMsgID(p *Packet) {
	p.writeUint16(CS_ACCOUNT_LOGOUT)
}
func (r *Cs_account_logout) Write(p *Packet) error {
	return nil
}
