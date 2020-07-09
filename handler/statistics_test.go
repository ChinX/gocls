package handler

import "testing"

func TestStatistics(t *testing.T) {
	a, b := splitSubject("net, syscall: ")
	t.Log(a, b)
	t.Log(splitSubject("[dev.link] cmd/link/internal/loader: support cloneToExternal for aux syms"))
}
