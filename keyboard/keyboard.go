package keyboard

import (
	"bufio"
	"os"
	"syscall"
	"unsafe"
)

func KhbitUnix() (byte, bool, error) {
	fd := int(os.Stdin.Fd())
	var termios syscall.Termios

	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&termios)),
		0, 0, 0,
	)

	if errno != 0 {
		khbitFallback()
	}

	oldTermios := termios

	termios.Lflag &^= syscall.ICANON | syscall.ECHO
	termios.Cc[syscall.VMIN] = 0
	termios.Cc[syscall.VTIME] = 0

	_, _, errno = syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&termios)),
		0, 0, 0,
	)

	if errno != 0 {
		khbitFallback()
	}

	defer func () {
		syscall.Syscall6(
			syscall.SYS_IOCTL,
			uintptr(fd),
			uintptr(syscall.TCSETS),
			uintptr(unsafe.Pointer(&oldTermios)),
			0, 0, 0,
		)
	}()

	var buf [1]byte
	n, err := syscall.Read(fd, buf[:])

	if err != nil || n == 0{
		return 0, false , nil
	}
	return buf[0], true, nil
}

func khbitFallback() bool { 
	reader := bufio.NewReader(os.Stdin)

	_, err := reader.Peek(1)
	return err == nil
}