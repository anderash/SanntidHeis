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
	c_io_floor := make(chan int)
	var decoded_input Input

	driver.InitDriver(c_input, c_output, c_io_floor)

	for {
		select{
		case byte_input := <-c_input:
			err := json.Unmarshal(byte_input, &decoded_input)
			if err != nil{
				fmt.Println("error: ", err)
			}
			fmt.Println("Input:", decoded_input.INPUT_TYPE, "Button type:", decoded_input.BUTTON_TYPE, "Floor:", decoded_input.FLOOR)

			if decoded_input.INPUT_TYPE == driver.BUTTON{
					output := Output{driver.LIGHT_OUTPUT, driver.BUTTON_LAMP, decoded_input.BUTTON_TYPE, decoded_input.FLOOR, 1, 0}
					encoded_output, err2 := json.Marshal(output)
					if err2 != nil{
						fmt.Println("error: ", err2)
					}
					c_output <- encoded_output
					if decoded_input.BUTTON_TYPE == driver.BUTTON_COMMAND && decoded_input.FLOOR == 2{
						output := Output{driver.MOTOR_OUTPUT, 0, 0, 0, 0, 1}
						encoded_output, err3 := json.Marshal(output)
						if err3 != nil{
							fmt.Println("error: ", err3)
						}
						c_output <- encoded_output
					}
			}
			if decoded_input.INPUT_TYPE == driver.FLOOR_SENSOR && decoded_input.FLOOR == 2{
				output := Output{driver.MOTOR_OUTPUT, 0, 0, 0, 0, 0}
				encoded_output, err4 := json.Marshal(output)
				if err4 != nil{
					fmt.Println("error: ", err4)
				}
				c_output <- encoded_output				
			}
		case floor := <- c_io_floor:
			fmt.Println(floor)
			
		}
	}

	output := Output{driver.MOTOR_OUTPUT, 0, driver.NOT_A_BUTTON, 0, 0, 1}
	encoded_output, err2 := json.Marshal(output)
	if err2 != nil{
		fmt.Println("error: ", err2)
	}
	c_output <- encoded_output
}