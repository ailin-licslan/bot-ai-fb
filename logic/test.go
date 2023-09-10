package logic

import (
	"fmt"
	"net/http"
)

// Test 测试请求
func Test(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, err := fmt.Fprint(w, "HELLO WORLD!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
