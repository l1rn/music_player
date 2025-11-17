package context

import (
	"context"
	"fmt"
	"music-player/keyboard"
	"time"
)

type KeyboardLifecycle struct {
	ctx       context.Context
	cancel    context.CancelFunc
	keyEvents chan rune
	isRunning bool
	stopChan  chan struct{}
}

func NewKeyboardLifecycle() *KeyboardLifecycle {
	ctx, cancel := context.WithCancel(context.Background())
	return &KeyboardLifecycle{
		ctx:       ctx,
		cancel:    cancel,
		keyEvents: make(chan rune, 100),
		isRunning: false,
		stopChan:  make(chan struct{}, 1),
	}
}

func Start(kl *KeyboardLifecycle) {
	if kl.isRunning {
		return
	}
	kl.isRunning = true

	go kl.lifecycleLoop()
}

func Stop(kl *KeyboardLifecycle) {
	if kl.isRunning {
		kl.cancel()
		kl.isRunning = false
		select {
		case kl.stopChan <- struct{}{}:
		default:
		}
		close(kl.keyEvents)
	}
}

func (kl *KeyboardLifecycle) Done() <-chan struct{} {
	return kl.ctx.Done()
}

func (kl *KeyboardLifecycle) lifecycleLoop() {
	defer func() {
		kl.isRunning = false
		close(kl.keyEvents)
	}()

	for {
		select {
		case <-kl.ctx.Done():
			return
		case <-kl.stopChan:
			return
		default:
			key, available, err := keyboard.KhbitUnix()
			if err != nil {
				continue
			}

			if available {
				select {
				case kl.keyEvents <- rune(key):
					fmt.Printf("\033[24mHello world!\033[0m \n")
				case <-kl.ctx.Done():
					return
				default:

				}
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (kl *KeyboardLifecycle) GetKeyEvents() <-chan rune {
	return kl.keyEvents
}
