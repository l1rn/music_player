package main

import (
	"fmt"
	"music-player/keyboard"
	"time"

)

func main() {
	fmt.Println("Keyboard detection test running...")
	fmt.Println("Press any key to test (will timeout after 10 seconds)")

	for {
		keyCode, detected, err := keyboard.KhbitUnix();
		if err != nil {
			fmt.Printf("err: %v", err)
		}
		if detected {
			switch keyCode{
			case 'q', 'Q': 
				fmt.Println("You pressed Q and quit!")
				return
			}
		}

		time.Sleep(100 * time.Millisecond)
		
		fmt.Printf(".")
	}
}