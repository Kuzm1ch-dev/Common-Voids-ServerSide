package main

import (
	"bufio"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"os"
	"sync"
)

const (
	CONN_TYPE = "tcp"
	CONN_HOST = "localhost"
	CONN_PORT = "25565"
)

var connMap = &sync.Map{}

type usersList struct {
	uid  string
	conn net.Conn
}

func main() {
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		logger.Error("Error listening: ", zap.Error(err))
		os.Exit(0)
	}
	defer l.Close()

	logger.Info("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error("Error accepting ", zap.Error(err))
			continue
		}
		logger.Info("New connection: " + conn.LocalAddr().String())

		id := uuid.New().String()
		connMap.Store(id, conn)

		go handleRequest(id, conn, connMap, logger)
	}
}

func handleRequest(id string, conn net.Conn, connMap *sync.Map, logger *zap.Logger) {
	defer func() {
		conn.Close()
		connMap.Delete(id)
	}()

	for {
		userInput, err := bufio.NewReader(conn).ReadString('\n')
		logger.Info("MSG:" + userInput)
		if err != nil {
			logger.Error("Error reading from client", zap.Error(err))
			return
		}

		connMap.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if _, err := conn.Write([]byte(userInput)); err != nil {
					logger.Error("Error on writing to connection", zap.Error(err))
				}
			}

			return true
		})
	}
}
