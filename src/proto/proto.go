package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"server/src/common"
)

// Encode кодирует сообщение
func Encode(pack common.Package) ([]byte, error) {
	// Считываем длину сообщения и преобразуем его в тип int32 (занимает 4 байта)
	var message = pack.Marshal()
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// пишем заголовок сообщения
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// записываем объект сообщения
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode декодирует сообщение
func Decode(reader *bufio.Reader) (common.Package, error) {
	// читаем длину сообщения
	lengthByte, _ := reader.Peek(4) // Считываем первые 4 байта данных
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return common.Package{}, err
	}

	// Buffered возвращает количество байтов, которые можно прочитать в буфере.
	if int32(reader.Buffered()) < length+4 {
		return common.Package{}, err
	}

	// читаем реальные данные сообщения
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return common.Package{}, err
	}
	return common.UnMarshal(pack[4:]), nil
}
