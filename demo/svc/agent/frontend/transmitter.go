package frontend

import (
	"encoding/binary"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
	"github.com/gorilla/websocket"
	"io"
	"net"
)

type socketOpt interface {
	MaxPacketSize() int
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())
}

type directTCPTransmitter struct {
}

// 来自客户端的消息
func (directTCPTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	reader, ok := ses.Raw().(io.Reader)

	// 转换错误，或者连接已经关闭时退出
	if !ok || reader == nil {
		return nil, nil
	}

	opt := ses.Peer().(socketOpt)

	if conn, ok := ses.Raw().(net.Conn); ok {

		for {

			// 有读超时时，设置超时
			opt.ApplySocketReadTimeout(conn, func() {

				var msgID int
				var msgData []byte

				// 接收来自客户端的封包
				msgID, msgData, err = RecvLTVPacketData(reader, opt.MaxPacketSize())

				// 尝试透传到后台或者解析
				if err == nil {
					msg, err = ProcFrontendPacket(msgID, msgData, ses)
				}

			})

			// 有错退出
			if err != nil {
				break
			}

			// msg=nil时,透传了客户端的封包到后台, 不用传给下一个proc, 继续重新读取下一个包
		}

	}

	return
}

// 网关发往客户端的消息
func (directTCPTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) (err error) {

	writer, ok := ses.Raw().(io.Writer)

	// 转换错误，或者连接已经关闭时退出
	if !ok || writer == nil {
		return nil
	}

	opt := ses.Peer().(socketOpt)

	// 有写超时时，设置超时
	opt.ApplySocketWriteTimeout(writer.(net.Conn), func() {

		err = util.SendLTVPacket(writer, ses.(cellnet.ContextSet), msg)

	})

	return
}

const (
	MsgIDSize = 2 // uint16
)

type directWSMessageTransmitter struct {
}

func (directWSMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil, nil
	}

	var (
		messageType int
		raw         []byte
	)

	for {
		messageType, raw, err = conn.ReadMessage()

		if err != nil {
			break
		}

		switch messageType {
		case websocket.BinaryMessage:
			msgID := binary.LittleEndian.Uint16(raw)
			msgData := raw[MsgIDSize:]

			// 尝试透传到后台或者解析
			if err == nil {
				msg, err = ProcFrontendPacket(int(msgID), msgData, ses)
			}
		}

		if err != nil {
			break
		}

	}

	return

}

func (directWSMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {

	conn, ok := ses.Raw().(*websocket.Conn)

	// 转换错误，或者连接已经关闭时退出
	if !ok || conn == nil {
		return nil
	}

	var (
		msgData []byte
		msgID   int
	)

	switch m := msg.(type) {
	case *cellnet.RawPacket: // 发裸包
		msgData = m.MsgData
		msgID = m.MsgID
	default: // 发普通编码包
		var err error
		var meta *cellnet.MessageMeta

		// 将用户数据转换为字节数组和消息ID
		msgData, meta, err = codec.EncodeMessage(msg, nil)

		if err != nil {
			return err
		}

		msgID = meta.ID
	}

	pkt := make([]byte, MsgIDSize+len(msgData))
	binary.LittleEndian.PutUint16(pkt, uint16(msgID))
	copy(pkt[MsgIDSize:], msgData)

	conn.WriteMessage(websocket.BinaryMessage, pkt)

	return nil
}
