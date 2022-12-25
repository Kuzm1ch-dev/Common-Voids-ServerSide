package common

import (
	"bufio"
	"io"
	"log"
	"net"
	"server/src/controller"
	"strings"
)

const (
	pAuth             int32 = 201
	pSuccessfulAuth         = 202
	pFailedAuth             = 203
	pSignIn                 = 204
	pSignInFail             = 205
	pSignInSuccessful       = 206
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
	s.DataBaseController = controller.DataBaseConnect("user=postgres password=root dbname=accounts sslmode=disable")
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
			log.Fatal("accept failed, err:", err)
			continue
		}
		s.Clients[s.CurrentID] = conn
		log.Println("New Connection: " + conn.LocalAddr().String())
		go s.ConnectionHandler(s.CurrentID)
		s.CurrentID++
	}
}

func (s AuthServer) ConnectionHandler(id int) error {

	defer s.closeConnection(id)

	//welcomeData = Encode(Package{UUIDPackage, "", s.Clients[id].Uuid.String()})

	for {
		reader := bufio.NewReader(s.Clients[id])
		pack, err := Decode(reader)
		if err == io.EOF {
			log.Fatal("Error:", err)
			return err
		}
		if err != nil {
			log.Fatal("decode msg failed, err:", err)
			return err
		}

		err = s.handleReceivedAuthPacket(pack)
		if err != nil {
			log.Fatal("Error:", err)
			return err
		}
	}
}

func (s AuthServer) handleReceivedAuthPacket(pack Package) error {
	switch pack.Code {
	case pAuth:
		email := strings.Split(pack.Data, ":")[0]
		password := strings.Split(pack.Data, ":")[1]
		s.DataBaseController.CheckUser(email, password)
	case pSignIn:
		email := strings.Split(pack.Data, ":")[0]
		password := strings.Split(pack.Data, ":")[1]
		s.DataBaseController.CreateUser(email, password)
	}
	return nil
}

func (s AuthServer) closeConnection(id int) {
	err := s.Clients[id].Close()
	if err != nil {
		log.Fatal("Error:", err)
	}
	delete(s.Clients, id)
}
