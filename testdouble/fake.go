package main

// type FakeUserRepo struct {
// 	data map[int]string
// 	mu   sync.Mutex
// }

// func NewFakeUserRepo() *FakeUserRepo {
// 	return &FakeUserRepo{data: make(map[int]string)}
// }

// func (f *FakeUserRepo) Save(id int, name string) {
// 	f.mu.Lock()
// 	defer f.mu.Unlock()
// 	f.data[id] = name
// }

// func (f *FakeUserRepo) Get(id int) string {
// 	f.mu.Lock()
// 	defer f.mu.Unlock()
// 	return f.data[id]
// }
