package utils

import (
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ReadFile 读取文件
func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	chunks := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}
	return string(chunks)
}

// Watcher 监控文件夹
func Watcher(path string, callback func(fileName string)) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
		return watcher
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		// 捕获异常
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				log.Printf("Panic: %v\n", err)
			}
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				// done <- true
				if !ok {
					log.Println("not ok event:", event)
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					time.Sleep(3 * time.Second)
					callback(event.Name)
				}
			case err, ok := <-watcher.Errors:
				log.Println("watch error:", err)
				if !ok {
					return
				}
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return watcher
}
