package frontend

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
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
