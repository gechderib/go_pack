package main

import (
	"fmt"
	"math"
	"strings"
)

func sayGreetings(name string) {
	fmt.Println("hello,", name)
}

func cycleName(n []string, f func(string)) {
	for _, value := range n {
		f(value)
	}
}

func circleArea(r float64) float64 {
	return math.Pi * r * r
}

func getInitials(n string) (string, string) {
	s := strings.ToUpper(n)
	names := strings.Split(s, " ")

	var initials []string

	for _, v := range names {
		initials = append(initials, v[:1])
	}

	if len(initials) > 1 {
		return initials[0], initials[1]
	}

	return initials[0], "_"

}

var points = []float32{89, 98.87, 76, 90.8, 87}

func sayHello() {
	fmt.Println(score)
	fmt.Println("hello mario ")
}

// func main() {

// x := 0

// for x < 5 {
// 	fmt.Println("value of x is: ", x)
// 	x++
// }

// for i := 0; i < 5; i++ {
// 	fmt.Println("the value is: ", i)
// }

// names := []string{"mario", "luigi", "jack"}
// cycleName(names, sayGreetings)
// for _, value := range names {
// 	fmt.Println(value)
// }

// a1 := circleArea(10.5)
// a2 := circleArea(10)

// fmt.Printf("%.2f, %.2f ", a1, a2)

// fn1, sn1 := getInitials("tifa lockhart")

// fmt.Println(fn1, sn1)

// }
