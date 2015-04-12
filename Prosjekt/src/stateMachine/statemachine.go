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

var position int

/*
   Etg.			Pos. nr.
    1 ............ 0
  	  ............ 1
	2 ............ 2
	  ............ 3
	3 ............ 4
	  ............ 5
	4 ............ 6
*/

var direction int

/*
	opp = 1, ned = -1, stillestående = 0
*/

var destination int

/*
	1. etg = 0
	2. etg = 1
	3. etg = 2
	4. etg = 3
*/

var state string

// Må på en eller annen måte sørge for at heisen går ned til 1. etg ved oppstart
func InitStateMachine(c_queMan_destination chan int, c_io_floor chan int, c_SM_output chan []byte) {

	run := false
	goDown := Output{1, -1, -1, -1, -1, -1}
	stopMotor := Output{1, -1, -1, -1, -1, 0}
init:
	for {
		select {
		case floorInput := <-c_io_floor:
			if floorInput == 0 {
				state = "idle"
				fmt.Printf("Arrived at floor 0, stopping motor")
				encoded_output, err := json.Marshal(stopMotor)
				if err != nil {
					fmt.Println("init JSON error: ", err)
				}
				c_SM_output <- encoded_output
				break init
			}
		case <-time.After(100 * time.Millisecond):
			if !run {
				fmt.Printf("Starting elevator")

				encoded_output, err := json.Marshal(goDown)
				if err != nil {
					fmt.Println("init JSON error: ", err)
				}
				c_SM_output <- encoded_output
				run = true
			}
		}
	}

	go stateMachine(c_queMan_destination, c_io_floor, c_SM_output)
}

func stateMachine(c_queMan_destination chan int, c_io_floor chan int, c_SM_output chan []byte) {

//	goUp := Output{1, -1, -1, -1, -1, 1}
//	goDown := Output{1, -1, -1, -1, -1, -1}
	stopMotor := Output{1, -1, -1, -1, -1, 0}

	openDoor := Output{0, 2, -1, -1, 1, -1}
	closeDoor := Output{0, 2, -1, -1, 0, -1}

	doorTimer := time.NewTimer(3 * time.Second)

	for {
		select {
		case dest := <-c_queMan_destination:
			destination = dest
			dest_pos := destination*2 - 2

			switch {

			case state == "move":

			case state == "at_floor":
				<-doorTimer.C
				encoded_output, err := json.Marshal(closeDoor)
				if err != nil {
					fmt.Println("SM JSON error: ", err)
				}
				c_SM_output <- encoded_output
				state = "idle"
				fallthrough

			case state == "idle":
				if dest_pos > position {
					direction = 1
					state = "move"
				} else if dest_pos < position {
					direction = -1
					state = "move"
				} else {
					direction = 0
					state = "at_floor"
				}

			}

		case floorInput := <-c_io_floor:
			fmt.Println(floorInput)
			switch {
			case state == "idle":

			case state == "move":
				if floorInput == destination {
					encoded_output, err := json.Marshal(stopMotor)
					if err != nil {
						fmt.Println("SM JSON error: ", err)
					}
					c_SM_output <- encoded_output

					fmt.Printf("Arrived at floor %d", floorInput)
					encoded_output, err = json.Marshal(openDoor)
					if err != nil {
						fmt.Println("SM JSON error: ", err)
					}
					c_SM_output <- encoded_output
					doorTimer.Reset(3 * time.Second)

					state = "at_floor"
				}

			case state == "at_floor":

			}

		case <-doorTimer.C:
			switch state{
			case "at_floor":
				encoded_output, err := json.Marshal(closeDoor)
				if err != nil {
					fmt.Println("SM JSON error: ", err)
				}
				c_SM_output <- encoded_output
			}

		}
	}
}
