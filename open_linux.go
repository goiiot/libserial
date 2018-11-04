// +build linux

/*
 * Copyright Go-IIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libserial

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func sysReadBaudRate(fd uintptr) uint32 {
	tty := &unix.Termios{}
	if _, _, err := unix.Syscall(
		unix.SYS_IOCTL, fd, unix.TCGETS, uintptr(unsafe.Pointer(tty)),
	); err != 0 {
		return 0
	} else {
		return tty.Cflag & 0x100f
	}
}

func (s *SerialPort) sysOpen(f *os.File, timeout uint8) error {
	tty := &unix.Termios{
		Cflag:  unix.CREAD | unix.CLOCAL | uint32(s.baudRate) | uint32(s.controlOptions),
		Ispeed: uint32(s.baudRate),
		Ospeed: uint32(s.baudRate),
		Iflag:  uint32(s.inputOptions),
		// Lflag:  0,
		// Oflag:  0,
	}

	// non block for read
	tty.Cc[unix.VMIN] = 1
	tty.Cc[unix.VTIME] = timeout

	if _, _, err := unix.Syscall(
		unix.SYS_IOCTL, f.Fd(), unix.TCSETS, uintptr(unsafe.Pointer(tty)),
	); err != 0 {
		return fmt.Errorf("failed :%v", err.Error())
	}

	// set blocking
	if err := unix.SetNonblock(int(f.Fd()), false); err != nil {
		return err
	}

	return nil
}

var validBaudRates = map[int]uint32{
	0:       unix.B0, // detect baud rate automatically
	50:      unix.B50,
	75:      unix.B75,
	110:     unix.B110,
	134:     unix.B134,
	150:     unix.B150,
	200:     unix.B200,
	300:     unix.B300,
	600:     unix.B600,
	1200:    unix.B1200,
	1800:    unix.B1800,
	2400:    unix.B2400,
	4800:    unix.B4800,
	9600:    unix.B9600,
	19200:   unix.B19200,
	38400:   unix.B38400,
	57600:   unix.B57600,
	115200:  unix.B115200,
	230400:  unix.B230400,
	460800:  unix.B460800,
	500000:  unix.B500000,
	576000:  unix.B576000,
	921600:  unix.B921600,
	1000000: unix.B1000000,
	1152000: unix.B1152000,
	1500000: unix.B1500000,
	2000000: unix.B2000000,
	2500000: unix.B2500000,
	3000000: unix.B3000000,
	3500000: unix.B3500000,
	4000000: unix.B4000000,
}