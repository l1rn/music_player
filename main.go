package main

import (
	"fmt"
	context "music-player/player"
)

func main() {
	kl := context.NewKeyboardLifecycle();
	context.Start(kl)
	defer context.Stop(kl)

	keyEvents := kl.GetKeyEvents()

	for {
		select {
		case key, ok := <-keyEvents:
			if !ok{
				fmt.Println("Key events channel closed")
				return
			}

			fmt.Printf("Key pressed: %c (ASCII: %d, Hex: 0x%02x)\n", key, key, key)
			
			if key == 'q' || key == 'Q' {
				fmt.Println("Quit key pressed, exiting...")
				return
			}

			if key == 'k' || key == 'K' {
				context.PlayerInit()
				fmt.Println("Enjoy the sound!")
			}
			
			if key == 's' || key == 'S' {
				context.Pause()
			}

			if key == 'w' || key == 'W' {
				context.Unpause()
			}

		case <-kl.Done():
			fmt.Println("Keyboard lifecycle ended")
			return
			
		}	
	}
}
