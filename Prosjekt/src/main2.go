// Main for testing driver
package main

import(
	"fmt"
	"./driver"
	//"time"
)

func main() {
	fmt.Printf("Starting driver\n")
	driver.Initiate()
	//driver.Set_button_lamp(driver.LIGHT_DOWN2,1)
	//time.Sleep(1000 * time.Millisecond)
	//driver.Set_door_open_lamp(1)
}