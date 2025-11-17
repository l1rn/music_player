package context

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var (
	streamer	beep.StreamSeekCloser
	ctrl 		*beep.Ctrl
	isPlaying 	bool
)


func PlayerInit() error {
	if streamer != nil {
		streamer.Close()
	}

	if _, err := os.Stat("music-test/wrld.mp3"); os.IsNotExist(err){
		log.Fatalf("File wrld.mp3 does not exist in current directory: %s", err)
	}

	file, err := os.Open("music-test/wrld.mp3")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	newStreamer, format, err := mp3.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	if !isPlaying{
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			log.Fatal(err)
		}
	}
	streamer = newStreamer

	return play(streamer)
}

func play(streamer beep.StreamSeekCloser) error {	
	speaker.Clear()

	done := make(chan bool)
	ctrl = &beep.Ctrl{Streamer: beep.Seq(streamer, beep.Callback(func() {
		done <- true
	}))}

	speaker.Play(ctrl)
	isPlaying = true

	go func() {
		<- done
		isPlaying = false
		streamer.Close()
	}()
	return nil
}

func Pause() {
	speaker.Lock()
}

func Unpause() {
	speaker.Unlock()
}