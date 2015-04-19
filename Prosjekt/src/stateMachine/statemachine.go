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
}

var elevatorState ElevState

var state string


func InitStatemachine(c_queMan_destination chan int, c_io_floor chan int, c_SM_output chan []byte, c_SM_state chan []byte) {
	// run := false
	goDown := Output{1, -1, -1, -1, -1, -1}
	stopMotor := Output{1, -1, -1, -1, -1, 0}

	elevatorState.POSITION = <-c_io_floor
	if elevatorState.POSITION != 0 {
		sendOutput(goDown, c_SM_output)
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
	sendOutput(stopMotor, c_SM_output)
	// Floor indicator lamp
	sendOutput(Output{0,1,-1, elevatorState.POSITION,1,-1}, c_SM_output)

	elevatorState.DIRECTION = 0
	sendState(elevatorState, c_SM_state)
	state = "idle"

	fmt.Printf("Statemachine operational\n")
	go statemachine(c_queMan_destination, c_io_floor, c_SM_output, c_SM_state)
}

func statemachine(c_queMan_destination chan int, c_io_floor chan int, c_SM_output chan []byte, c_SM_state chan []byte) {

	goUp := Output{1, -1, -1, -1, -1, 1}
	goDown := Output{1, -1, -1, -1, -1, -1}
	stopMotor := Output{1, -1, -1, -1, -1, 0}

	openDoor := Output{0, 2, -1, -1, 1, -1}
	closeDoor := Output{0, 2, -1, -1, 0, -1}

	doorTimer := time.NewTimer(3 * time.Second)
	doorTimer.Stop()

	for {
		select {
		case elevatorState.DESTINATION = <-c_queMan_destination:
			fmt.Printf("SM: Destination \n")
			switch state {

			case "move":

			case "at_floor":		// Hvis man får en ny DEST før door timeren går ut.
				<-doorTimer.C
				sendOutput(closeDoor, c_SM_output)
				if elevatorState.DESTINATION > elevatorState.POSITION {
					elevatorState.DIRECTION = 1
					state = "move"
					sendOutput(goUp, c_SM_output)
					sendState(elevatorState, c_SM_state)
				} else if elevatorState.DESTINATION < elevatorState.POSITION {
					elevatorState.DIRECTION = -1
					state = "move"
					sendOutput(goDown, c_SM_output)
					sendState(elevatorState, c_SM_state)
				} else {
					elevatorState.DIRECTION = 0
					state = "at_floor"
					sendOutput(openDoor, c_SM_output)
					sendOutput(stopMotor, c_SM_output)
					sendState(elevatorState, c_SM_state)
					doorTimer.Reset(3 * time.Second)
				}
				fallthrough

			case "idle":
				if elevatorState.DESTINATION > elevatorState.POSITION {
					elevatorState.DIRECTION = 1
					state = "move"
					sendOutput(goUp, c_SM_output)
					sendState(elevatorState, c_SM_state)

				} else if elevatorState.DESTINATION < elevatorState.POSITION {
					elevatorState.DIRECTION = -1
					state = "move"
					sendOutput(goDown, c_SM_output)
					sendState(elevatorState, c_SM_state)
				} else {
					elevatorState.DIRECTION = 0
					state = "at_floor"
					sendOutput(openDoor, c_SM_output)
					sendOutput(stopMotor, c_SM_output)
					sendState(elevatorState, c_SM_state)
					doorTimer.Reset(3 * time.Second)
				}

			}
			fmt.Println(state)

		case elevatorState.POSITION = <-c_io_floor:
			fmt.Printf("SM: Floorinput \n")
			fmt.Println(elevatorState.POSITION)
			sendOutput(Output{0,1,-1, elevatorState.POSITION,1,-1}, c_SM_output)	// Tenner etg.-lys

			switch state {
			case "idle": //Skal ikke skje

			case "move":
				if elevatorState.POSITION == elevatorState.DESTINATION {
					sendOutput(stopMotor, c_SM_output)

					// elevatorState.DIRECTION = 0		// DETTE MÅ ENDRES. Skal bare sette DIR = 0 når man går i state idle.

					fmt.Printf("SM: Arrived at floor %d \n", elevatorState.POSITION)
					sendOutput(openDoor, c_SM_output)
					doorTimer.Reset(3 * time.Second)
					state = "at_floor"

				} else {
					state = "move"
				}

			case "at_floor":
				state = "idle" //Ikke tenkt noe mer over dette

			}
			sendState(elevatorState, c_SM_state)
			fmt.Println(state)

		case <-doorTimer.C:
			fmt.Printf("SM Doortimer\n")
			switch state {
			case "at_floor":
				sendOutput(closeDoor, c_SM_output)
				state = "idle"
				elevatorState.DIRECTION = 0
				sendState(elevatorState, c_SM_state)

			}
			fmt.Println(state)

		}
	}
}

func sendOutput(output Output, c_SM_output chan []byte) {
	encoded_output, err := json.Marshal(output)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	c_SM_output <- encoded_output
}

func sendState(elevatorState ElevState, channel chan []byte) {
	encoded_output, err := json.Marshal(elevatorState)
	if err != nil {
		fmt.Println("SM JSON error: ", err)
	}
	channel <- encoded_output
}
