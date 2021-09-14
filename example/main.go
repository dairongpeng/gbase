package main

import "fmt"

func main() {
	t := GetTicket("你好", "世界")
	fmt.Println(t.A)
	fmt.Println(t.B)
}

type Ticket struct {
	A string
	B string
}

func GetTicket(hello ,world string) *Ticket {
	t := Ticket{
		A: hello,
		B: world,
	}
	return &t
}