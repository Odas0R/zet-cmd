package fs

import "os"

func EditorHasSession() bool {
	socket := os.Getenv("NVIM_SOCKET")
	return Exec("test -S "+socket) == nil
}
