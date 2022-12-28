package common

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"server/src/controller"
	"strings"
)

const (
	pAuth             int32 = 201
	pAuthSuccessful         = 202
	pAuthFailed             = 203
	pSignIn                 = 204
	pSignInSuccessful       = 205
	pSignInFail             = 206
)

type AuthServer struct {
	Addr               string
	Clients            map[int]net.Conn
	CurrentID          int
	DataBaseController *controller.DataBaseController
}

func (s *AuthServer) Init(addr string) {
	s.Addr = addr
	s.Clients = make(map[int]net.Conn)
	s.DataBaseController = controller.DataBaseConnect(os.Getenv("DATABASE_STRING"))
}

func (s *AuthServer) ListenAndServe() error {
	Listener, err := net.Listen("tcp", s.Addr+":30001")
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Println("Auth server started ", s.Addr)

	for {
		conn, err := Listener.Accept()
		if err != nil {
			log.Println("accept failed, err:", err)
			continue
		}
		s.Clients[s.CurrentID] = conn
		log.Println("New Connection: " + conn.LocalAddr().String())
		go s.ConnectionHandler(s.CurrentID)
		s.CurrentID++
	}
}

func (s AuthServer) ConnectionHandler(id int) error {

	defer s.CloseConnection(id)

	//welcomeData = Encode(Package{UUIDPackage, "", s.Clients[id].Uuid.String()})

	for {
		reader := bufio.NewReader(s.Clients[id])
		pack, err := Decode(reader)
		if err == io.EOF {
			log.Println("Error:", err)
			return err
		}
		if err != nil {
			log.Println("decode msg failed, err:", err)
			return err
		}

		err = s.handleReceivedAuthPacket(pack, id)
		if err != nil {
			log.Println("Handler Error:", err)
			return err
		}
	}
}

func (s AuthServer) handleReceivedAuthPacket(pack Package, id int) error {
	switch pack.Code {
	case pAuth:
		email := strings.Split(pack.Data, ":")[0]
		password := strings.Split(pack.Data, ":")[1]
		token := s.DataBaseController.CheckUser(email, password)
		if token != "" {
			accessPackage := Package{pAuthSuccessful, token, ""}
			s.SingleCast(accessPackage, id)
		} else {
			failedPackage := Package{pAuthFailed, "", ""}
			s.SingleCast(failedPackage, id)
			s.Clients[id].Close()
		}
		break
	case pSignIn:
		email := strings.Split(pack.Data, ":")[0]
		password := strings.Split(pack.Data, ":")[1]
		if s.DataBaseController.CreateUser(email, password) {
			successfulPackage := Package{pSignInSuccessful, "", ""}
			s.SingleCast(successfulPackage, id)
		} else {
			failedPackage := Package{pSignInFail, "", ""}
			s.SingleCast(failedPackage, id)
		}
		break
	}
	return nil
}

func (s AuthServer) CloseConnection(id int) {
	err := s.Clients[id].Close()
	log.Println("Close connection id: ", id)
	if err != nil {
		log.Println("Error:", err)
	}
	delete(s.Clients, id)
}

func (s AuthServer) SingleCast(pack Package, id int) error {
	data, err := Encode(pack)
	if err != nil {
		log.Println("Encode Error:", err)
		return err
	}
	_, err = s.Clients[id].Write(data)
	if err != nil {
		log.Println("Write Error:", err.Error())
		return err
	}
	return nil
}
