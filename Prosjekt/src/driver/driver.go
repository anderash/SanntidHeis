package driver

import "C"

import (
	"encoding/json"
	"fmt"
	"time"
	"os"
)

type Input struct {
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

type Output struct {
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
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0}}

var floor_status = [N_FLOORS]int{0, 0, 0, 0}

func InitDriver(c_io_button chan []byte, c_io_floor chan int, c_stateMach_output chan []byte, c_queMan_output chan []byte) {
	Io_init()

	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		if i != 0 {
			setButtonLamp(BUTTON_CALL_DOWN, i, 0)
		}
		if i != (N_FLOORS - 1) {
			setButtonLamp(BUTTON_CALL_UP, i, 0)
		}
		setButtonLamp(BUTTON_COMMAND, i, 0)
	}

	// Zero door open lamp
	setDoorOpenLamp(0)

	// Make sure motor is dead
	setMotorDirection(0)

	go checkFloor(c_io_floor)
	go checkButtons(c_io_button)
	go Send_output(c_stateMach_output, c_queMan_output)

	fmt.Printf("Driver initiated!\n")
}

func getFloorSignal() int {
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
func getButtonSignal() (int, int) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for button := 0; button < 3; button++ {
			if (button_matrix[floor][button] != -1) && (Io_read_bit(button_matrix[floor][button]) == 1) {
				if button_status[floor][button] == 0 {
					button_status[floor][button] = 1
					return floor, button
					// 
				}
			} else {
				button_status[floor][button] = 0
			}
		}
	}
	return -1, -1
	// Return value if no button aktivatet
}


func setButtonLamp(button int, floor int, value int) {
	if value == 1 {
		Io_set_bit(buttonlight_matrix[floor][button])
	} else {
		Io_clear_bit(buttonlight_matrix[floor][button])
	}
}

func setDoorOpenLamp(value int) {
	if value == 1 {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func setFloorIndicator(floor int) {
	// Making sure allways one floor light is lit
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

func setMotorDirection(direction int) {
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

func checkButtons(c_input chan []byte) {

	for {
		if floor, button_type := getButtonSignal(); floor != -1 {
			input := Input{BUTTON, button_type, floor}
			encoded_input, err := json.Marshal(input)
			if err != nil {
				fmt.Println("error: ", err)
			}
			c_input <- encoded_input
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func checkFloor(c_io_floor chan int) {
	floor := getFloorSignal()
	prevFloor := -1

	c_io_floor <- floor

	for {
		floor = getFloorSignal()
		if floor != prevFloor && floor != -1 {
			c_io_floor <- floor
		}
		prevFloor = floor

		time.Sleep(10 * time.Millisecond)
	}
}

func Send_output(c_stateMach_output chan []byte, c_queMan_output chan []byte) {
	var output Output
	for {
		select {
		case enc_output := <-c_stateMach_output:
			err3 := json.Unmarshal(enc_output, &output)
			if err3 != nil {
				fmt.Println("Driver_o JSON error: ", err3)
			}

			if output.OUTPUT_TYPE == LIGHT_OUTPUT {
				if output.LIGHT_TYPE == FLOOR_INDICATOR {
					setFloorIndicator(output.FLOOR)
				} else if output.LIGHT_TYPE == DOOR_LAMP {
					if getFloorSignal() == -1 {
						fmt.Printf("Invalid position. Cannot open door. Call maintenance")
						os.Exit(2)
					}
					setDoorOpenLamp(output.VALUE)
				}

			}

			if output.OUTPUT_TYPE == MOTOR_OUTPUT {
				setMotorDirection(output.OUTPUT_DIRECTION)
			}

		case enc_light_output := <-c_queMan_output:
			err3 := json.Unmarshal(enc_light_output, &output)
			if err3 != nil {
				fmt.Println("Driver_l JSON error: ", err3)
			}
			setButtonLamp(output.BUTTON_TYPE, output.FLOOR, output.VALUE)
				
		}
	}

}
