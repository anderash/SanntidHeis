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
)

type Input struct{
	INPUT_TYPE int
	/*
	button = 0
	floor = 1
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
	LIGHT = 0
	MOTOR = 1 
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

var status_matrix = [N_FLOORS][3]


func Initiate(c_input chan []byte, c_output chan []byte) {
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


	go Check_input(c_input)
	go Send_output(c_output)

	fmt.Printf("Initiated!\n")
}

// Funker
func Get_floor_signal() int {
	if Io_read_bit(SENSOR_FLOOR1) == 1 {
		return 0
	} else if Io_read_bit(SENSOR_FLOOR2) == 1 {
		return 1
	} else if Io_read_bit(SENSOR_FLOOR3) == 1 {
		return 2
	} else if Io_read_bit(SENSOR_FLOOR4) == 1{
		return 3
	} else {
		return -1
	}
}

// Denne vil ikke få med seg knappetrykk dersom noen holder en knapp inne i en lavere etg. Må fikses!
func Get_button_signal() (int, int) {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < 3; j++ {
			if (button_matrix[i][j] != -1) && (Io_read_bit(button_matrix[i][j]) == 1){
				return i, j
				// i tilsvarer etage (0 = 1. etg, 1 = 2. etg. osv)
				// j tilsvarer type knapp. (0 = opp-knapp, 1 = ned-knapp, 2 = knapp inne i heis)
			}else{
				continue
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
func Check_input(c_input chan []byte) {
	var button_status [N_FLOORS][3]int := {
		{0,0,0},
		{0,0,0},
		{0,0,0}}

	for {
		if floor, button_type := Get_button_signal(); floor != -1 && {
			input := Input{BUTTON, button_type, floor}
			encoded_input, err := json.Marshal(input)
			if err != nil{
				fmt.Println("error: ", err)
			}
			c_input <- encoded_input
		}
		if floor := Get_floor_signal(); floor != -1 {
			input := Input{FLOOR_SENSOR, NOT_A_BUTTON, floor}
			encoded_input, err2 := json.Marshal(input)
			if err2 != nil{
				fmt.Println("error: ", err2)
			}
			c_input <- encoded_input
		} 
	}
}



// IKKE KOMPLETT
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
					Set_floor_indicator(decoded_output.FLOOR)
				} else{
					Set_button_lamp(decoded_output.BUTTON_TYPE, decoded_output.FLOOR, decoded_output.VALUE)
				}
			}

		}
	}

}

