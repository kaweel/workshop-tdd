package main

import (
	"fmt"
	"sync"
)

func main() {
	Dummy(123, "This is a dummy")

	repo := NewFakeUserRepo()
	repo.Save(1, "Alice")
	println("Fake User:", repo.Get(1))

	fmt.Printf("Stub : %v\n", GetUserStub(1))
	fmt.Printf("Stub : %v\n", GetUserStub(2))

	spy := &EmailSpy{}
	spy.SendEmail("test@example.com")
	spy.SendEmail("admin@example.com")
	fmt.Printf("Spy : %v\n", spy.Calls)
}

func Dummy(id int, _ string) {
	fmt.Println("Dummy ID:", id)
}

type FakeUserRepo struct {
	data map[int]string
	mu   sync.Mutex
}

func NewFakeUserRepo() *FakeUserRepo {
	return &FakeUserRepo{data: make(map[int]string)}
}

func (f *FakeUserRepo) Save(id int, name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.data[id] = name
}

func (f *FakeUserRepo) Get(id int) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.data[id]
}

func GetUserStub(id int) string {
	if id == 1 {
		return "Alice"
	}
	return "Unknown"
}

type EmailSpy struct {
	Calls []string
}

func (s *EmailSpy) SendEmail(email string) {
	s.Calls = append(s.Calls, email)
}
