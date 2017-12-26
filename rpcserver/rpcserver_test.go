package rpcserver_test

import (
	"testing"
	"rpcdef"
	"rpcserver"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"reflect"
)

func TestServer(t *testing.T){
	s := new(rpcserver.Server)
	err := s.Launch("127.0.0.1:1234", "./users.txt") //launching server
	if err != nil {
		t.Fatal("Cannot initialize server: " + err.Error())
	}
	
	client, err := net.Dial("tcp", "127.0.0.1:1234") //connecting
	if err != nil {
		t.Fatal("Cannot connect to server: " + err.Error())
	}
	
	var u1, u2, u3 rpcdef.User
	
	c := jsonrpc.NewClient(client)
	l := "user1"
	err = c.Call("Users.Add", &l, &u1) //adding one user
	if err != nil {
		t.Fatal("Cannot add user: " + err.Error())
	}
	if u1.Login != l{
		t.Errorf("Wrong reply: expected \"%s\", got \"%s\"", l, u1.Login)
	}
	
	u1.Login = "user2"
	err = c.Call("Users.Change", &u1, &u2) //changing login
	if err != nil {
		t.Error("Cannot change user: " + err.Error())
	}
	if !reflect.DeepEqual(u1, u2) {
		t.Error("Failed to modify user correctly")
	}
	
	l = "user3"
	err = c.Call("Users.Add", &l, &u2) //adding second user
	if err != nil {
		t.Fatal("Cannot add user: " + err.Error())
	}
	if u2.Login != l{
		t.Errorf("Wrong reply: expected \"%s\", got \"%s\"", l, u2.Login)
	}
	
	err = c.Call("Users.Get", &u1.Uuid, &u3) //retrieving first user
	if err != nil {
		t.Error("Cannot retrieve user: " + err.Error())
	}
	if !reflect.DeepEqual(u1, u3) {
		t.Error("Failed to retrieve user correctly")
	}
	
	err = c.Call("Users.Get", &u2.Uuid, &u3) //retrieving second user
	if err != nil {
		t.Error("Cannot retrieve user: " + err.Error())
	}
	if !reflect.DeepEqual(u2, u3) {
		t.Error("Failed to retrieve user correctly")
	}
	
	s.Shutdown("./users.txt") //shutting down server
	
	os.Remove("./users.txt") //removing temporary "users base" file
}