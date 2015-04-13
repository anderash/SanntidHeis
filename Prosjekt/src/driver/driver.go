package driver

/*
#cgo LDFLAGS: -lcomedi -lm
#cgo CFLAGS: -std=c99
#include "io.h"
*/
import "C"

import (
	"fmt"
	"encoding/json"
	"time"
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
	DOOR_LAMP = 2
	*/

	BUTTON_TYPE int
	/*
	BUTTON_CALL_UP = 0
    BUTTON_CALL_DOWN = 1
    BUTTON_COMMAND = 2
    NO_BUTTON = -1
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

// Denne bør helst endres til const (hvis mulig)
var button_matrix = [N_FLOORS][3]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}

var buttonlight_matrix = [N_FLOORS][3]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var button_status = [N_FLOORS][3]int{
		{0,0,0},
		{0,0,0},
		{0,0,0}}

var floor_status = [N_FLOORS]int{0,0,0,0}


func InitDriver(c_buttonEvents chan []byte, c_floorEvents chan int, c_outputs chan []byte) {
	Io_init()

	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		if i != 0{
			Set_button_lamp(BUTTON_CALL_DOWN, i, 0)
		}
		if i != (N_FLOORS - 1) {
			Set_button_lamp(BUTTON_CALL_UP, i, 0)
		}
		Set_button_lamp(BUTTON_COMMAND, i, 0)
	}

	// Zero door open lamp
	Set_door_open_lamp(0)

	// Make sure motor is dead
	Set_motor_direction(0)


	go Check_floor(c_floorEvents)
	go Check_buttons(c_buttonEvents)
	go Send_output(c_outputs)

	fmt.Printf("Initiated!\n")
}

// Funker
func get_floor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) == 1 {
		return 0
	}
	if Io_read_bit(SENSOR_FLOOR2) == 1 {
		return 1
	}
	if Io_read_bit(SENSOR_FLOOR3) == 1 {
		return 2
	}
	if Io_read_bit(SENSOR_FLOOR4) == 1 {
		return 3
	}
	return -1
}

// Denne vil ikke få med seg knappetrykk dersom noen holder en knapp inne i en lavere etg. Må fikses!
func Get_button_signal() (int, int) {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < 3; j++ {
			if (button_matrix[i][j] != -1) && (Io_read_bit(button_matrix[i][j]) == 1){
				if button_status[i][j] == 0 {
					button_status[i][j] = 1
					return i, j
					// i tilsvarer etage (0 = 1. etg, 1 = 2. etg. osv)
					// j tilsvarer type knapp. (0 = opp-knapp, 1 = ned-knapp, 2 = knapp inne i heis)
				}
			}else{
				button_status[i][j] = 0
			}
		}
	}
	return -1, -1
	//Dette returneres hvis ingen knapp detektert

}

// Funker
func Set_button_lamp(button int, floor int, value int) {
	if value == 1 {
		Io_set_bit(buttonlight_matrix[floor][button])
	} else {
		Io_clear_bit(buttonlight_matrix[floor][button])
	}
}

// Funker
func Set_door_open_lamp(value int) {
	if value == 1 {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

// Funker
func Set_floor_indicator(floor int) {
	// Passer her på at ett lys alltid er tent
	if floor&0x02 == 0x02 {
		Io_set_bit(LIGHT_FLOOR_IND1)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor&0x01 == 0x01 {
		Io_set_bit(LIGHT_FLOOR_IND2)
	} else {
		Io_clear_bit(LIGHT_FLOOR_IND2)
	}
	
}

// Funker
func Set_motor_direction(direction int) {
	if direction == 0 {
		Io_write_analog(MOTOR, 0)
	} else if direction > 0 {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
	} else if direction < 0 {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, 2800)
	}

}



// Funker.
func Check_buttons(c_input chan []byte) {

	for {
		if floor, button_type := Get_button_signal(); floor != -1 {
			input := Input{BUTTON, button_type, floor}
			encoded_input, err := json.Marshal(input)
			if err != nil{
				fmt.Println("error: ", err)
			}
			c_input <- encoded_input
		}

		time.Sleep(10*time.Millisecond)
	}
}


func Check_floor(c_io_floor chan int) {
	floor := get_floor_signal()
	prevFloor := -1

	c_io_floor <- floor

	for {
		floor = get_floor_signal()
		if floor != prevFloor  &&  floor != -1 {
			c_io_floor <- floor
		}
		prevFloor = floor


		time.Sleep(10*time.Millisecond)
	}
}


//Funker, men med variabel reaksjonstid
func Send_output(c_output chan []byte) {
	var decoded_output Output
	for{
		select{
		case output := <- c_output:
			err3 := json.Unmarshal(output, &decoded_output)
			if err3 != nil{
				fmt.Println("error: ", err3)
			}

			if decoded_output.OUTPUT_TYPE == LIGHT_OUTPUT {
				if decoded_output.BUTTON_TYPE == NOT_A_BUTTON {
					if decoded_output.LIGHT_TYPE == FLOOR_INDICATOR{
						Set_floor_indicator(decoded_output.FLOOR)
					}else if decoded_output.LIGHT_TYPE == DOOR_LAMP{
						Set_door_open_lamp(decoded_output.VALUE)
					}
					
				} else{
					Set_button_lamp(decoded_output.BUTTON_TYPE, decoded_output.FLOOR, decoded_output.VALUE)
				}
			}


			if decoded_output.OUTPUT_TYPE == MOTOR_OUTPUT {
				Set_motor_direction(decoded_output.OUTPUT_DIRECTION)
			}
		}
	}

}

