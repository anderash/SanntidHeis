// Main for testing driver
package main

import(
	"fmt"
	"./driver"
	"encoding/json"
	// "time"
)


type Input struct{
	INPUT_TYPE int
	/*
	BUTTON = 0
	FLOOR_SENSOR = 1
	*/
	BUTTON_TYPE int
	/*
	BUTTON_CALL_UP = 0
    BUTTON_CALL_DOWN = 1
    BUTTON_COMMAND = 2
    NO_BUTTON = -1
	*/
	FLOOR int

}

type Output struct{
	OUTPUT_TYPE int
	/*
	LIGHT_OUTPUT = 0
	MOTOR_OUTPUT = 1 
	*/

	LIGHT_TYPE int
	/*
	BUTTON_LAMP = 0
	FLOOR_INDICATOR = 1
	*/

	BUTTON_TYPE int
	/*
	BUTTON_CALL_UP = 0
    BUTTON_CALL_DOWN = 1
    BUTTON_COMMAND = 2
    NOT_A_BUTTON = -1
	*/

	FLOOR int

	VALUE int
	/*
	on = 1
	off = 0
	*/

	OUTPUT_DIRECTION int
	/*
	UP = 1
	STOP = 0
	DOWN = -1
	*/

}



func main() {
	c_input := make(chan []byte)
	c_output := make(chan []byte)
	var decoded_input Input

	driver.Initiate(c_input, c_output)

	for {
		select{
		case byte_input := <-c_input:
			err := json.Unmarshal(byte_input, &decoded_input)
			if err != nil{
				fmt.Println("error: ", err)
			}
			fmt.Println("Input:", decoded_input.INPUT_TYPE, "Button type:", decoded_input.BUTTON_TYPE, "Floor:", decoded_input.FLOOR)

			if decoded_input.BUTTON_TYPE == driver.BUTTON{
					output := Output{driver.LIGHT_OUTPUT, driver.BUTTON_LAMP, decoded_input.BUTTON_TYPE, decoded_input.FLOOR, 1, 0}
					encoded_output, err2 := json.Marshal(output)
					if err2 != nil{
						fmt.Println("error: ", err2)
					}
					c_output <- encoded_output
			}
		}
	}

// 	output := Output{driver.LIGHT_OUTPUT, driver.FLOOR_INDICATOR, driver.NOT_A_BUTTON, 1, 1, 0}
// 	encoded_output, err2 := json.Marshal(output)
// 	if err2 != nil{
// 		fmt.Println("error: ", err2)
// 	}
// 	c_output <- encoded_output
}