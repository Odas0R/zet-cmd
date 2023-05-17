package fs

import "os"

func HasNvimSession() bool {
	socket := os.Getenv("NVIM_SOCKET")
	return Exec("test -S "+socket) == nil
}
