package stateMachine

import (
	"encoding/json"
	"fmt"
	"time"
)

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

type ElevState struct {
	POSITION int

	DIRECTION int
	/*
		opp = 1, ned = -1, stillestående = 0
	*/
	DESTINATION int
	/*
		1. etg = 0
		2. etg = 1
		3. etg = 2
		4. etg = 3
	*/
	MOVING bool
}

var elevatorState ElevState

var state string

func InitStatemachine(c_queMan_destination chan int, c_io_floor chan int, c_stMachine_output chan []byte, c_stMachine_state chan []byte) {
	// run := false
	goDown := Output{1, -1, -1, -1, -1, -1}
	stopMotor := Output{1, -1, -1, -1, -1, 0}
	openDoor := Output{0, 2, -1, -1, 1, -1}

	elevatorState.POSITION = <-c_io_floor
	if elevatorState.POSITION != 0 {
		sendOutput(goDown, c_stMachine_output)
	}
	elevatorState.DIRECTION = -1
	elevatorState.DESTINATION = 0

	for {
		elevatorState.POSITION = <-c_io_floor
		fmt.Println("FLOOR SENSOR SIGNAL:", elevatorState.POSITION)
		if elevatorState.POSITION == 0 {
			break
		}
	}
	sendOutput(stopMotor, c_stMachine_output)
	sendOutput(openDoor, c_stMachine_output)

	sendOutput(Output{0, 1, -1, elevatorState.POSITION, 1, -1}, c_stMachine_output) // Floor indicator lamp

	elevatorState.DIRECTION = 0
	elevatorState.MOVING = false
	sendState(elevatorState, c_stMachine_state)
	state = "at_floor"

	fmt.Printf("Statemachine operational\n")
	go statemachine(c_queMan_destination, c_io_floor, c_stMachine_output, c_stMachine_state)
}

func statemachine(c_queMan_destination chan int, c_io_floor chan int, c_stMachine_output chan []byte, c_stMachine_state chan []byte) {

	goUp := Output{1, -1, -1, -1, -1, 1}
	goDown := Output{1, -1, -1, -1, -1, -1}
	stopMotor := Output{1, -1, -1, -1, -1, 0}

	openDoor := Output{0, 2, -1, -1, 1, -1}
	closeDoor := Output{0, 2, -1, -1, 0, -1}

	doorTimer := time.NewTimer(3 * time.Second)

	for {
		select {
		case elevatorState.DESTINATION = <-c_queMan_destination:
			switch state {

			case "at_floor":
				<-doorTimer.C
				sendOutput(closeDoor, c_stMachine_output)
				fallthrough

			case "idle":
				if elevatorState.DESTINATION > elevatorState.POSITION {
					elevatorState.DIRECTION = 1
					elevatorState.MOVING = true
					state = "move"
					sendOutput(goUp, c_stMachine_output)
					sendState(elevatorState, c_stMachine_state)

				} else if elevatorState.DESTINATION < elevatorState.POSITION {
					elevatorState.DIRECTION = -1
					elevatorState.MOVING = true
					state = "move"
					sendOutput(goDown, c_stMachine_output)
					sendState(elevatorState, c_stMachine_state)
				} else {
					elevatorState.DIRECTION = 0
					elevatorState.MOVING = false
					state = "at_floor"
					sendOutput(openDoor, c_stMachine_output)
					sendOutput(stopMotor, c_stMachine_output)
					sendState(elevatorState, c_stMachine_state)
					doorTimer.Reset(3 * time.Second)
				}
			}

		case elevatorState.POSITION = <-c_io_floor:
			sendOutput(Output{0, 1, -1, elevatorState.POSITION, 1, -1}, c_stMachine_output) // Lighting floor lamp

			switch state {
			case "move":
				if elevatorState.POSITION == elevatorState.DESTINATION {
					sendOutput(stopMotor, c_stMachine_output)
					fmt.Printf("SM: Arrived at floor %d \n", elevatorState.POSITION)
					sendOutput(openDoor, c_stMachine_output)
					doorTimer.Reset(3 * time.Second)
					state = "at_floor"
					elevatorState.MOVING = false

				} else {
					state = "move"
					elevatorState.MOVING = true
				}

			case "at_floor":
				state = "idle"
				elevatorState.MOVING = false

			}
			sendState(elevatorState, c_stMachine_state)

		case <-doorTimer.C:
			switch state {
			case "at_floor":
				sendOutput(closeDoor, c_stMachine_output)
				state = "idle"
				elevatorState.DIRECTION = 0
				elevatorState.MOVING = false
				sendState(elevatorState, c_stMachine_state)

			}
		}
	}
}

func sendOutput(output Output, c_stMachine_output chan []byte) {
	encoded_output, err := json.Marshal(output)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	c_stMachine_output <- encoded_output
}

func sendState(elevatorState ElevState, channel chan []byte) {
	encoded_output, err := json.Marshal(elevatorState)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output
}
