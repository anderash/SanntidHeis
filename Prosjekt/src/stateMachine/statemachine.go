package stateMachine

import (
	"fmt"
	"time"
	"encoding/json"
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
init:
	for {
		select {
		case floorInput := <-c_io_floor:
			if floorInput == 0 {
				state = "idle"
				break init
			}
		case <-time.After(100 * time.Millisecond):
			if !run {
				goDown := Output{1, -1, -1, -1, -1, -1}
				encoded_output, err := json.Marshal(goDown)
				if err != nil {
					fmt.Println("init JSON error: ", err)
				}
				c_SM_output <- encoded_output
				run = true
			}
		}
	}

	go stateMachine(c_queMan_destination, c_io_floor)
}

func stateMachine(c_queMan_destination chan int, c_io_floor chan int) {

	for {
		select {
		case dest := <-c_queMan_destination:
			destination = dest
			dest_pos := destination*2 - 2

			switch {
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

			case state == "move":

			case state == "at_floor":
			}

		case input := <-c_io_floor:
			fmt.Println(input)
			switch {
			case state == "idle":

			case state == "move":

			case state == "at_floor":
			}
		}
	}
}
