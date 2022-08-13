package ServerEndpoint

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/kxg3030/shermie-nat-traversal/Concrate"
	"net"
)

type ConnPeer struct {
	conn         net.Conn
	readerWriter *bufio.ReadWriter
	uniqueId     int
}

func NewConnPeer(conn net.Conn) *ConnPeer {
	return &ConnPeer{
		conn: conn,
		readerWriter: bufio.NewReadWriter(
			bufio.NewReader(conn),
			bufio.NewWriter(conn),
		),
	}
}

// action(1)-uniqueId(4)-to(4)-bodyLen(4)-body
func (i *ConnPeer) ReadMessage() (*Concrate.Message, error) {
	message := Concrate.NewMessage()
	// 读取命令
	action := make([]byte, 1)
	_, err := i.readerWriter.Read(action)
	if err != nil {
		return nil, err
	}
	message.Action = action[0]
	uniqueIdByte := make([]byte, 4)
	_, err = i.readerWriter.Read(uniqueIdByte)
	// 读取唯一id
	message.UniqueId = int32(binary.LittleEndian.Uint32(uniqueIdByte))
	// 读取发送给谁
	toByte := make([]byte, 4)
	_, err = i.readerWriter.Read(toByte)
	message.To = int32(binary.LittleEndian.Uint32(toByte))
	// 读取消息长度
	headerLenByte := make([]byte, 4)
	_, err = i.readerWriter.Read(headerLenByte)
	message.HeaderLen = binary.LittleEndian.Uint32(headerLenByte)
	// 读取消息体
	message.Data = make([]byte, message.HeaderLen)
	_, err = i.readerWriter.Read(message.Data)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (i *ConnPeer) WriteMessage(data []byte, uniqueId, to int) error {
	buffer := &bytes.Buffer{}
	_ = binary.Write(buffer, binary.LittleEndian, uint8(Concrate.ActionData))
	_ = binary.Write(buffer, binary.LittleEndian, int32(uniqueId))
	_ = binary.Write(buffer, binary.LittleEndian, int32(to))
	_ = binary.Write(buffer, binary.LittleEndian, int32(len(data)))
	_ = binary.Write(buffer, binary.LittleEndian, data)
	_, err := i.conn.Write(buffer.Bytes())
	return err
}

func (i *ConnPeer) SetUniqueId(id int) {
	i.uniqueId = id
}

func (i *ConnPeer) GetTag() int {
	return i.uniqueId
}

func (i *ConnPeer) Close() {
	_ = i.conn.Close()
}
