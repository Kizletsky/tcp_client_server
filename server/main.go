package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var storage db

func main() {
	storage = newInMemoryDb()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8888"
	}

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	log.Printf("Serving: %s", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	log.Printf("New connection: %s", c.RemoteAddr())
	defer c.Close()

	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Println(err)
			log.Println("Client disconnected")
			return
		}

		log.Printf("Received message: %s", data)
		req, err := parseRequest(data)
		res := performResponse(req, err)
		io.WriteString(c, res)
	}
}

func performResponse(req request, reqErr error) string {
	res := newResponse()
	res.setOk()

	if reqErr != nil {
		res.setErr(reqErr.Error())
	} else {
		switch req.command {
		case RequestGetCommand:
			if req.key == "*" {
				if records, err := storage.getAll(); err != nil {
					res.setErr(err.Error())
					log.Println(err)
				} else {
					res.records = records
				}
			} else {
				if records, err := storage.get(req.key); err != nil {
					res.setErr(err.Error())
					log.Println(err)
				} else {
					res.records = records
				}
			}
		case RequestPutCommand:
			model := model{
				key:       req.key,
				value:     req.value,
				timestamp: req.timestamp,
			}

			if err := storage.put(model); err != nil {
				res.setErr(err.Error())
				log.Println(err)
			}
		}
	}
	return res.build()
}

func parseRequest(data string) (req request, err error) {
	scanner := bufio.NewScanner(strings.NewReader(data))
	scanner.Split(bufio.ScanWords)
	req = request{}

	var index int
	for scanner.Scan() {
		index++

		token := scanner.Text()
		switch index {
		case 1:
			req.command = token
		case 2:
			req.key = token
		case 3:
			value, parseErr := strconv.ParseFloat(token, 64)
			if parseErr != nil {
				log.Println(parseErr)
				err = errors.New(BadValueFormatErr)
			}
			req.value = value
		case 4:
			timestamp, parseErr := strconv.ParseInt(token, 10, 64)
			if parseErr != nil {
				log.Println(parseErr)
				err = errors.New(BadTimestampErr)
			}
			req.timestamp = time.Unix(timestamp, 0)
		default:
			err = errors.New(BadRequestFormatErr)
		}

	}

	if err == nil {
		err = req.validate()
	}

	return
}
