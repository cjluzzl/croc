package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// runClient spawns threads for parallel uplink/downlink via TCP
func runClient(connectionType string, codePhrase string) {
	var wg sync.WaitGroup
	wg.Add(numberConnections)
	for id := 0; id < numberConnections; id++ {
		go func(id int) {
			defer wg.Done()
			port := strconv.Itoa(27001 + id)
			connection, err := net.Dial("tcp", "localhost:"+port)
			if err != nil {
				panic(err)
			}
			defer connection.Close()

			message := receiveMessage(connection)
			fmt.Println(message)
			sendMessage(connectionType+"."+codePhrase, connection)
			if connectionType == "s" {
				message = receiveMessage(connection)
				fmt.Println(message)
				// Send file name
				sendMessage("filename", connection)
				// Send file size
				time.Sleep(3 * time.Second)
				sendMessage("filesize", connection)
				// TODO: Write data from file

				// TODO: Release from connection pool
				// POST /release
			} else {
				fileName := receiveMessage(connection)
				fileSize := receiveMessage(connection)
				fmt.Println(fileName, fileSize)
				// TODO: Pull data and write to file

				// TODO: Release from connection pool
				// POST /release
			}

		}(id)
	}
	wg.Wait()
}
