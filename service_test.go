package sproto

import (
	"bytes"
	"testing"
)

type Test int

var inst = Test(5)

var barCalled bool

func (t *Test) Foobar(req *FoobarRequest, resp *FoobarResponse) {
	if t != &inst {
		panic(t)
	}
	resp.What = req.What
}

func (t *Test) Foo(resp *FooResponse) {
	if t != &inst {
		panic(t)
	}
	resp.Ok = Bool(true)
}

func (t *Test) Bar() {
	if t != &inst {
		panic(t)
	}
	barCalled = true
}

func TestFoobarService(t *testing.T) {
	name := "test.foobar"
	input := "hello"
	rw := bytes.NewBuffer(nil)

	// client
	client, _ := NewService(rw, protocols, HEAD_UINT16)
	req := FoobarRequest{
		What: &input,
	}
	call, err := client.Go(name, &req, nil)
	if err != nil {
		t.Fatalf("client call failed:%s", err)
	}

	// server
	server, _ := NewService(rw, protocols, HEAD_UINT16)
	if err := server.Register(&inst); err != nil {
		t.Fatalf("register service failed:%s", err)
	}
	if err := server.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}

	//
	if err := client.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}
	<-call.Done
	resp := call.Resp.(*FoobarResponse)
	if resp.What == nil || *resp.What != input {
		t.Fatalf("unexpected response:%v", resp.What)
	}
}

func TestFooService(t *testing.T) {
	name := "test.foo"
	rw := bytes.NewBuffer(nil)

	// client
	client, _ := NewService(rw, protocols, HEAD_UINT16)
	call, err := client.Go(name, nil, nil)
	if err != nil {
		t.Fatalf("client call failed:%s", err)
	}

	// server
	server, _ := NewService(rw, protocols, HEAD_UINT16)
	if err := server.Register(&inst); err != nil {
		t.Fatalf("register service failed:%s", err)
	}
	if err := server.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}

	//
	if err := client.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}
	<-call.Done
	resp := call.Resp.(*FooResponse)
	if resp.Ok == nil || !*resp.Ok {
		t.Fatalf("unexpected response:%v", resp.Ok)
	}
}

func TestBarService(t *testing.T) {
	name := "test.bar"
	rw := bytes.NewBuffer(nil)

	// client
	client, _ := NewService(rw, protocols, HEAD_UINT16)
	err := client.Invoke(name, nil)
	if err != nil {
		t.Fatalf("client call failed:%s", err)
	}

	// server
	barCalled = false
	server, _ := NewService(rw, protocols, HEAD_UINT16)
	if err := server.Register(&inst); err != nil {
		t.Fatalf("register service failed:%s", err)
	}
	if err := server.DispatchOnce(); err != nil {
		t.Fatalf("dispatch service failed:%s", err)
	}

	//
	if !barCalled {
		t.Fatal("unexpected dispatch")
	}
}
