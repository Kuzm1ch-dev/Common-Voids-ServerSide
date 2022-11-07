package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode кодирует сообщение
func Encode(message string) ([]byte, error) {
	// Считываем длину сообщения и преобразуем его в тип int32 (занимает 4 байта)
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// пишем заголовок сообщения
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// записываем объект сообщения
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode декодирует сообщение
func Decode(reader *bufio.Reader) (string, error) {
	// читаем длину сообщения
	lengthByte, _ := reader.Peek(4) // Считываем первые 4 байта данных
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// Buffered возвращает количество байтов, которые можно прочитать в буфере.
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// читаем реальные данные сообщения
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}
