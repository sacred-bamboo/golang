package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"io"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	srv := http.Server{
		Addr:    ":80",
		Handler: mux,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("listen: %s\n", err)
		}
	}()
	glog.Info("Server Started")
	<-done
	glog.Info("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatalf("Server Shutdown Failed:%+v", err)
	}
	glog.Info("Server Exited Properly")
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	glog.V(4).Info("Entering delay handler")
	delay := randInt(10,20)
	time.Sleep(time.Millisecond*time.Duration(delay))
	io.WriteString(writer, "===================Details of the http request header:============\n")
	req, err := http.NewRequest("GET", "http://service2", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
	lowerCaseHeader := make(http.Header)
	for key, value := range request.Header {
		lowerCaseHeader[strings.ToLower(key)] = value
	}
	glog.Info("headers:", lowerCaseHeader)
	req.Header = lowerCaseHeader
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		glog.Info("HTTP get failed with error: ", "error", err)
	} else {
		glog.Info("HTTP get succeeded")
	}
	if resp != nil {
		resp.Write(writer)
	}
	glog.V(4).Infof("Respond in %d ms", delay)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
