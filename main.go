package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func SayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!!!")
	// io.WriteString(w, "Hello World!!!\n")
}

func StartHttpServer(hs *http.Server) error {
	http.HandleFunc("/", SayHello)
	fmt.Println("服务启动")
	if err := hs.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	c, cancel := context.WithCancel(ctx)
	group, errCtx := errgroup.WithContext(c)
	hsr := &http.Server{Addr: ":8080"}

	group.Go(func() error {
		return StartHttpServer(hsr)
	})

	group.Go(func() error {
		<-errCtx.Done()
		fmt.Println("服务关闭")
		return hsr.Shutdown(errCtx)
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			case <-ch:
				cancel()
			}
		}
	})

	if err := group.Wait(); err != nil {
		fmt.Println("group wait error:", err)
	}
	fmt.Println("所有线程结束")
}
