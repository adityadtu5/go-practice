package main

import(
	"fmt"
	"net"
	"net/rpc"
)

type Message struct{
	To, From, Body string
}

type RemoteActor struct{

}

func (a *RemoteActor) ReceiveMessage(m Message, reply *string) error{
	fmt.Printf("Remote actor %s received message: %s\n", m.To, m.Body)
	*reply = "Message Recieved!"
	return nil
}

func main(){
	actor := new(RemoteActor)
	rpc.Register(actor)
	l,err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("Error listening: %s\n",err)
		return
	}
	defer l.Close()
	fmt.Println("Remote actor listening on localhost:1234")
	rpc.Accept(l)
}