package rpcdef

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Uuid []byte

type User struct {
	Uuid       Uuid
	Login      string
	Registered time.Time
}

type Users struct {
	sync.Mutex
	u []User
	m []*sync.Mutex
}

/*using a slice with mutexes instead of database/sql or whatever
to avoid messing with heavy DBs in this simple example

server itself is not affected*/

func (us *Users) Add(arg *string, reply *User) error { //register a new user
	reply.Uuid = make([]byte, 16, 16)
	if _, err := rand.Read(reply.Uuid); err != nil {
		return errors.New("Unable to generate UUID")
	}
	reply.Login = *arg
	reply.Registered = time.Now()
	us.Lock() //locking the entire slice for thread safety
	us.m = append(us.m, new(sync.Mutex))
	us.u = append(us.u, *reply)
	us.Unlock()
	return nil
}

func (us *Users) Get(arg *Uuid, reply *User) error { //find user by uuid
	var i int
	for i, _ = range us.u { //will only need to lock one user later
		if bytes.Compare(us.u[i].Uuid, *arg) == 0 { //no need to lock here, UUID is not modified
			break
		} else {
			if i == len(us.u)-1 {
				return errors.New("No user with that UUID")
			}
		}
	}
	us.m[i].Lock() //getting the user safely
	u := us.u[i]
	us.m[i].Unlock()
	reply.Uuid = make([]byte, 16, 16)
	copy(reply.Uuid, u.Uuid)
	reply.Login = u.Login
	reply.Registered = u.Registered

	return nil
}

func (us *Users) Change(arg *User, reply *User) error { //change login
	var i int
	for i, _ = range us.u { //will only need to lock one user later
		if bytes.Compare(us.u[i].Uuid, arg.Uuid) == 0 { //no need to lock here, UUID is not modified
			break
		} else {
			if i == len(us.u)-1 {
				return errors.New("No user with that UUID")
			}
		}
	}
	us.m[i].Lock() //accessing the user safely
	us.u[i].Login = arg.Login

	*reply = us.u[i]
	us.m[i].Unlock()
	return nil
}

func (us *Users) Init(filename string) error { //load users from file
	f := make([]byte, 0)
	f, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		us.u = make([]User, 0)
		us.m = make([]*sync.Mutex, 0)
		return nil
	}
	if err != nil {
		return err
	}
	us.u = make([]User, 0)
	json.Unmarshal(f, &us.u)
	us.m = make([]*sync.Mutex, len(us.u))
	return nil
}

func (us *Users) Finalize(filename string) error { //save users to file
	f, err := json.Marshal(us.u)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, f, 0666)
	return err
}
