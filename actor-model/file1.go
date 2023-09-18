package main

import(
	"fmt"
	"net/rpc"
)

type Message struct{
	To, From, Body string
}

type Actor struct{
	Name string
	Messages []Message
}

func (a *Actor) SendMessage(m Message){
	a.Messages = append(a.Messages, m)
}

func (a *Actor) ProcessMessages(){
	for _,m := range a.Messages{
		fmt.Printf("Actor %s received message: %s\n", a.Name, m.Body)
	}
	a.Messages = nil
}

type ActorManager struct {
	Actors map[string]*Actor
}

func NewActorManager() *ActorManager {
	return &ActorManager{Actors: make(map[string]*Actor)}
}

func(s *ActorManager) RegisterActor(name string){
	s.Actors[name] = &Actor{
		Name: name,
	}
}

func (s *ActorManager) SendMessage(m Message){
	if actor, ok := s.Actors[m.To]; ok{
		actor.SendMessage(m)
	}else{
		client, err := rpc.Dial("tcp", m.To)
		if err != nil{
			fmt.Printf("Error connecting to remote actor %s: %s\n",m.To, err)
			return
		}
		defer client.Close()
		var reply string
		err = client.Call("RemoteActor.ReceiveMessage", m, &reply)
		if err != nil {
			fmt.Printf("Error sending message to remote actor %s: %s\n", m.To, err)
			return
		}
	}
}

type RemoteActor struct{

}

func (a *RemoteActor) ReceiveMessage(m Message, reply *string) error{
	fmt.Printf("Remote actor %s received message: %s\n", m.To, m.Body)
	*reply = "Message Recieved!"
	return nil
}

func main(){
	actorManager := NewActorManager()
	actorManager.RegisterActor("actor1")
	actorManager.RegisterActor("actor2")
	actorManager.SendMessage(Message{
		"actor1",
		"actor2",
		"Hello from actro2 to actor1",
	})

	actorManager.SendMessage(Message{
		"localhost:1234",
		"actor1",
		"Hello from actro1 to remote actor on localhost:1234",
	})

	for _,actor := range actorManager.Actors{
		actor.ProcessMessages()
	}
}

