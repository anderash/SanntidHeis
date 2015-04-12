package stateMachine

import (
	"fmt"
)


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
func InitStateMachine(c_dest_from_queue chan int, c_floor_from_io chan int){
	state = "idle"
	go stateMachine(c_dest_from_queue, c_floor_from_io)
}


func stateMachine(c_dest_from_queue chan int, c_floor_from_io chan int){

	for{
		select{
		case dest := <- c_dest_from_queue:
			destination = dest
			dest_pos = destination*2 - 2 

			switch{
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

		case input := <- c_floor_from_io:
			fmt.Println(input)
			switch{
			case state == "idle":

			case state == "move":

			case state == "at_floor":
			}
		}	
	}
}