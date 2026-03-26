package main

import (
	"errors"
	"fmt"
)

type Address struct {
	City  string
	State string
}

type BaseModel struct {
	Id int `json:"id"`
}
type User struct {
	BaseModel
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
}

var score = 988

func updateUser(user *User) {
	user.Name = "updatedname"
	user.Id = 99
}

func (user User) Greet() string {
	return "hello " + user.Name
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero is not allowed")
	}
	return a / b, nil
}

func main() {

	// address := Address{
	// 	City:  "Gondar",
	// 	State: "Amhara",
	// }

	// user := User{
	// 	// Id:90,
	// 	BaseModel: BaseModel{Id: 90},
	// 	Name:      "Getachew",
	// 	Email:     "gechderib@gmail.com",
	// 	Address:   address,
	// }

	// ans := user.Greet()
	// fmt.Println(ans)
	// updateUser(&user)
	// fmt.Println(user)
	// fmt.Println(user.Id)

	// test()
	// area()

	// circle := Circle{Radius: 5}
	// rectangle := Rectangle{Width: 10, Height: 5}
	// triangle := Triangle{Base: 10, Height: 5}

	// println("Circle Area:", calculateArea(&circle))
	// println("Rectangle Area:", calculateArea(&rectangle))
	// println("Triangle Area:", calculateArea(&triangle))

	// println("Circle Perimeter:", CalculatePerimeter(&circle))
	// println("Rectangle Perimeter:", CalculatePerimeter(&rectangle))
	// println("Triangle Perimeter:", CalculatePerimeter(&triangle))

	resutt, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", resutt)
	}
}
