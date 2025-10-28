package headerinspectionmiddleware

import (
	"fmt"
	"net/http"
	"time"
)

func InspectHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: To replace with loggin function when logging function is built
		fmt.Println("time of log: ", time.Now())
		origin := r.Header.Get("Origin")
		fmt.Println("origin: ", origin)
		//TODO: To replace with loggin function when logging function is built
		referer := r.Header.Get("Referer")
		fmt.Println("referer: ", referer)
		next.ServeHTTP(w, r)
	})
}
