package webserver

import (
    "github.com/rs/cors"
	"github.com/gorilla/mux"
	// "github.com/gorilla/con text"   // don' use !!!!
	builtin_context "context"
	"fmt"
	"github.com/mebusy/goweb/tools"
	"log"
	"net/http"
	_ "net/http/pprof" // pprof server
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

/*
func clearStaleRequest() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("clearStaleRequest goroutine paniced:", r)
			clearStaleRequest()
		}
	}()

	for {
		time.Sleep(3 * time.Minute)
		context.Purge(45)
	}
}
//*/

var gitcommit string

func infoHandler(w http.ResponseWriter, r *http.Request) {
	// time.Sleep( 10 * time.Second ) // test gracefully shutdown
	fmt.Fprintf(w, "git commit: %s\n", gitcommit)
	fmt.Fprintf(w, "host: %s\n", tools.GetIP())

	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("Error Getting Rlimit ", err)
	}
	fmt.Fprintf(w, "rlimt cur:%d , max: %d \n", rLimit.Cur, rLimit.Max)
}

func StartServer(r *mux.Router, listenPort, verbose int, _GitCommit string) {

	gitcommit = _GitCommit
	tools.MaxOpenFiles()

	if runtime.GOOS == "darwin" || verbose == 1 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	r.HandleFunc("/info", infoHandler)

	// profile server
	go func() {
		port := listenPort + 1
		log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

    corsobj := cors.New(cors.Options{
        AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "auth-ts", "auth-sig"} , 
    })
    _ = corsobj


	// main server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", listenPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		/*
			Handler:      context.ClearHandler(r), // Pass our instance of gorilla/mux in.
			/*/
		// Handler: cors.Default().Handler(r),
		Handler: corsobj.Handler( r ), 
		//*/
	}

	go func() {
		log.Println("listening on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Println(err) // do NOT fatal if you want gracefully shutdown
			} else {
				log.Fatal(err)
			}
		}
	}()

	// clean gorilla/con text leap
	// go clearStaleRequest()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := builtin_context.WithTimeout(builtin_context.Background(), 15*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
