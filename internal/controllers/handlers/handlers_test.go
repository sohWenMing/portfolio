package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

var server = &http.Server{}
var pingChan = make(chan bool, 1)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	r := chi.NewRouter()
	r.Get("/testping", TestPing)
	r.Get("/test_usercreate", TestUserCreateHandler)
	server.Addr = ":8000"
	server.Handler = r

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("server error: ", err)
		}
	}()
	go testPing(pingChan)

	isServerStarted := <-pingChan
	if !isServerStarted {
		fmt.Println("server never started!")
		os.Exit(1)
	}

	exitCode := m.Run()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("err occured during server shutdown %v", err)
	}
	cancel()
	os.Exit(exitCode)
}

func TestPingGetRequest(t *testing.T) {
	statusCode, err := ping()
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	if statusCode != 200 {
		t.Errorf("got %d\nwant %d", 200, statusCode)
	}
}

func TestCreateUserRequest(t *testing.T) {
	req := httptest.NewRequest("POST", "/TestUserCreateHandler", nil)
	req.PostForm = url.Values{
		"email":    {"wenming.soh@gmail.com"},
		"password": {"Holoq123holoq123"},
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	TestUserCreateHandler(rec, req)
	res := rec.Result()
	if res.StatusCode != 200 {
		t.Errorf("got %d\nwant %d", 200, res.StatusCode)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	expected := "username: wenming.soh@gmail.com password: Holoq123holoq123"
	got := string(bodyBytes)
	if got != expected {
		t.Errorf("got %s\nwant %s", got, expected)
		return
	}
}

func TestInMemoryPingHandler(t *testing.T) {
	testReq := httptest.NewRequest("GET", "/testping", nil)
	rec := httptest.NewRecorder()
	TestPing(rec, testReq)
	code := rec.Result().StatusCode
	if code != 200 {
		t.Errorf("got %d\nwant%d", code, 200)
		return
	}
	bodyBytes, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	gotString := string(bodyBytes)
	wantString := "ping successful"
	if gotString != wantString {
		t.Errorf("got %s\nwant %s", gotString, wantString)
	}
}

func testPing(pingChan chan<- bool) {
	for i := 1; i < 4; i++ {
		fmt.Println("Server Ping Attempt: ", i)
		statusCode, err := ping()
		if err != nil {
			fmt.Println("error occured: ", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if statusCode == 200 {
			pingChan <- true
			return
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	pingChan <- false
}

func ping() (int, error) {
	testReq, err := http.NewRequest("GET", "http://localhost:8000/testping", nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(testReq)
	if err != nil {
		return 0, err
	}
	statusCode := res.StatusCode
	return statusCode, nil
}
