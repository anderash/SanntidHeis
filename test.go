package main

import "fmt"
type Input struct{
	SENSOR_TYPE int
	BUTTON_TYPE int
	FLOOR int
}

func main() {
	test_struct := Input{1,2,3} 
	fmt.Println(test_struct.SENSOR_TYPE)
}