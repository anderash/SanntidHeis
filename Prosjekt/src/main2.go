// Main for testing driver
package main

import(
	"fmt"
	"./driver"
	// "time"
)

func main() {
	fmt.Printf("Starting driver\n")
	driver.Initiate()
	driver.Set_motor_direction(0)

	// for {
	// 	fmt.Println(driver.Get_button_signal())
	// }
	driver.Set_button_lamp(driver.BUTTON_CALL_DOWN, 1 ,1)
	// driver.Set_door_open_lamp(1)
	// time.Sleep(1000 * time.Millisecond)
	// driver.Set_door_open_lamp(0)}

}