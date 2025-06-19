package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func generateHeaders(allhrs []string) map[string]string {
	headersMap := make(map[string]string)

	for i := range allhrs {
		if allhrs[i] == "" {
			continue
		}
		hrs := strings.Split(allhrs[i], ": ")
		headersMap[hrs[0]] = hrs[1]
	}
	return headersMap
}

func handleConnection(conn net.Conn, directory string) {

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading the request: ", err.Error())
	}
	req := string(buf[:n])

	readReq := strings.Split(req, "\r\n")
	method := strings.Split(readReq[0], " ")[0]
	path := strings.Split(readReq[0], " ")[1]

	allHeaders := readReq[1 : len(readReq)-1]

	headersMap := generateHeaders(allHeaders)

	if (method == "GET") && (path == "/") {
		indexOk := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\n")
		conn.Write([]byte(indexOk))
	} else if (method == "GET") && (strings.HasPrefix(path, "/echo/")) {
		echo := path[6:]
		echoOk := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echo), echo)
		conn.Write([]byte(echoOk))
	} else if (method == "POST") && (strings.HasPrefix(path, "/files/")) {
		file := path[len("/files/"):]
		fullPath := filepath.Join(directory, file)
		//splitN splits the headers and the body; ignores CRLF within the body.
		//parts[0] contains the headers; parts[1] contains the content
		parts := strings.SplitN(req, "\r\n\r\n", 2)
		expectedLength, err := strconv.Atoi(headersMap["Content-Length"])
		if err != nil {
			fmt.Println("malformed content length")
		}
		// finish reading the rest of the body
		if len(parts[1]) < expectedLength {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("unable to read the rest of the body")
			}
			parts[1] += string(buf[:n])
		}
		createdFile, err := os.Create(fullPath)
		if err != nil {
			fmt.Println("unable to create the file")
		}
		createdFile.Write([]byte(parts[1]))
		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
	} else if (method == "GET") && (strings.HasPrefix(path, "/user-agent")) {
		userAgent := headersMap["User-Agent"]
		agentTypeOK := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		conn.Write([]byte(agentTypeOK))
	} else if (method == "GET") && (strings.HasPrefix(path, "/files/")) {
		file := path[len("/files/"):]
		fullPath := filepath.Join(directory, file)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			notFound := fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
			conn.Write([]byte(notFound))
		} else {
			fileFound := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(content), content)
			conn.Write([]byte(fileFound))
		}
		fileMissing := fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write([]byte(fileMissing))

	} else {
		notFound := fmt.Sprintf("HTTP/1.1 404 Not Found\r\n\r\n")
		conn.Write([]byte(notFound))
	}

	defer conn.Close()
}

func main() {
	directory := flag.String("directory", "", "Directory to serve files from")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn, *directory)
	}
}
