package driver 

/*
#cgo LDFLAGS: -lcomedi -lm
#cgo CFLAGS: -std=c99
#include "io.h"
*/
import "C"

import(
	"fmt"
)

// Denne bør helst endres til const!!
var button_matrix = [4][3]int{
				[3]int{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    			[3]int{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
   				[3]int{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    			[3]int{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}

func Initiate(){
	// Init hardware
	C.io_init()

	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < 3; j++{
			Set_button_lamp(button_matrix[i][j], 0)
		}
	}

	fmt.Printf("Initiated!\n")
}

func Get_floor_signal() int {
	    if (C.io_read_bit(C.int(SENSOR_FLOOR1)))
        return 0;
    else if (C.io_read_bit(C.int(SENSOR_FLOOR2)))
        return 1;
    else if (C.io_read_bit(C.int(SENSOR_FLOOR3)))
        return 2;
    else if (C.io_read_bit(C.int(SENSOR_FLOOR4)))
        return 3;
    else
        return -1;
}

func Get_button_signal() {
	
}



func Set_button_lamp(button int, value int){
	if button == -1{
		return
	}
	if value == 1 {
		C.io_set_bit(C.int(button))
	}else {
		C.io_clear_bit(C.int(button))
	}
}

func Set_door_open_lamp(value int){
	if value == 1 {
		C.io_set_bit(LIGHT_DOOR_OPEN)
	}else {
		C.io_clear_bit(LIGHT_DOOR_OPEN)
	}
}



// func main(){
// 	initiate()
// 	set_door_open_lamp(1)

// }