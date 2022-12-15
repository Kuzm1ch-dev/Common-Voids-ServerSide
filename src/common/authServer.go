package common

import (
	"bufio"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"server/src/game"
)

type AuthServer struct {
	Addr      string
	Clients   map[int]Client
	CurrentID int
}

func (s *AuthServer) Init(addr string) {
	s.Addr = addr
	s.Clients = make(map[int]Client)
}

func (s *AuthServer) ListenAndServe() error {
	Listener, err := net.Listen("tcp", s.Addr+":30001")
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Println("Server started %s", s.Addr)

	for {
		conn, err := Listener.Accept()
		if err != nil {
			log.Fatal("accept failed, err:", err)
			continue
		}
		s.Clients[s.CurrentID] = Client{uuid.New(), conn, game.Player{}}
		log.Println("New Connection: " + conn.LocalAddr().String())
		go s.ConnectionHandler(s.CurrentID)
		s.CurrentID++
	}
}

func (s AuthServer) ConnectionHandler(id int) error {

	defer s.closeConnection(id)

	//Отправляем UUID игроку
	welcomeData, err := Encode(Package{UUIDPackage, "", s.Clients[id].Uuid.String()})
	log.Println("welcome: ", string(welcomeData))
	log.Println("Игроку присвоен uuid: " + s.Clients[id].Uuid.String())

	if err != nil {
		log.Fatal("Error:", err)
		return err
	}

	_, err = s.Clients[id].Conn.Write(welcomeData)
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}

	for {
		reader := bufio.NewReader(s.Clients[id].Conn)
		pack, err := Decode(reader)
		if err == io.EOF {
			log.Fatal("Error:", err)
			return err
		}
		if err != nil {
			log.Fatal("decode msg failed, err:", err)
			return err
		}

		err = s.handleReceivedPacket(pack, id)
		if err != nil {
			log.Fatal("Error:", err)
			return err
		}
	}
}

func (s AuthServer) handleReceivedPacket(pack Package, id int) error {

	return nil
}

func (s AuthServer) closeConnection(id int) {
	err := s.Clients[id].Conn.Close()
	if err != nil {
		log.Fatal("Error:", err)
	}
	delete(s.Clients, id)
}
