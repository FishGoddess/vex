// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/13 00:09:25

package vex

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	ProtocolVersion = uint8(1)

	okMark = uint8(0)
	errorMark = uint8(1)
)

var (
	ProtocolVersionMismatchErr = errors.New("the version of protocol used in client and server doesn't match")
)

type request struct {
	version uint8
	command string
	args    [][]byte
}

func readRequest(conn net.Conn) (*request, error) {

	reader := bufio.NewReader(conn)

	version, err := readVersion(reader)
	if err != nil {
		return nil, err
	}

	command, err := readCommand(reader)
	if err != nil {
		return nil, err
	}

	args, err := readArgs(reader)
	if err != nil {
		return nil, err
	}

	return &request{
		version: version,
		command: command,
		args:    args,
	}, nil
}

func readVersion(reader *bufio.Reader) (uint8, error) {
	version, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// Check version
	if version != ProtocolVersion {
		return 0, ProtocolVersionMismatchErr
	}
	return version, nil
}

func readCommand(reader *bufio.Reader) (string, error) {

	// Read the length of command
	buffer := make([]byte, 2)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return "", err
	}

	commandLength := binary.BigEndian.Uint16(buffer)

	// Read the command by length
	buffer = make([]byte, commandLength)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func readArgs(reader *bufio.Reader) ([][]byte, error) {

	// Read the length of args
	buffer := make([]byte, 4)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}

	// Read all args by length
	argsLength := binary.BigEndian.Uint32(buffer)
	args := make([][]byte, argsLength)
	for i := uint32(0); i < argsLength; i++ {

		_, err = io.ReadFull(reader, buffer)
		if err != nil {
			return nil, err
		}

		// Add arg to args
		argLength := binary.BigEndian.Uint32(buffer)
		arg := make([]byte, argLength)
		_, err = io.ReadFull(reader, arg)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}

	return args, nil
}

func writeRequest(conn net.Conn, req *request) error {

	command := []byte(req.command)
	commandLength := make([]byte, 2)

	// Notice: len returns an int value which may overflow in uint16
	// However, I think this will not happen?
	binary.BigEndian.PutUint16(commandLength, uint16(len(command)))

	// Write version, command
	buffer := make([]byte, 0, len(command) + len(commandLength) + 1)
	buffer = append(buffer, req.version)
	buffer = append(buffer, commandLength...)
	buffer = append(buffer, command...)
	_, err := conn.Write(buffer)
	if err != nil {
		return err
	}

	// Notice: len returns an int value which may overflow in uint32
	buffer = make([]byte, 4)
	binary.BigEndian.PutUint32(buffer, uint32(len(req.args)))
	_, err = conn.Write(buffer)
	if err != nil {
		return err
	}

	for _, arg := range req.args {
		buffer = make([]byte, 4, len(arg) + 4)
		binary.BigEndian.PutUint32(buffer, uint32(len(arg)))
		buffer = append(buffer, arg...)
		_, err = conn.Write(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func readResponse(conn net.Conn) ([]byte, error) {

	reader := bufio.NewReader(conn)

	// For checking version
	_, err := readVersion(reader)
	if err != nil {
		return nil, err
	}

	mark, err := readMark(reader)
	if err != nil {
		return nil, err
	}

	content, err := readContent(reader)
	if err != nil {
		return nil, err
	}

	// Return an error if this response is an error
	if mark == errorMark {
		return content, errors.New(string(content))
	}
	return content, nil
}

func readMark(reader *bufio.Reader) (uint8, error) {
	mark, err := reader.ReadByte()
	if err != nil {
		return errorMark, err
	}
	return mark, nil
}

func readContent(reader *bufio.Reader) ([]byte, error) {

	// Read the length of content
	buffer := make([]byte, 8)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}

	contentLength := binary.BigEndian.Uint64(buffer)

	// Read the content by length
	buffer = make([]byte, contentLength)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func writeResponse(conn net.Conn, mark uint8, content []byte) error {

	contentLength := make([]byte, 8)
	binary.BigEndian.PutUint64(contentLength, uint64(len(content)))

	_, err := conn.Write([]byte{
		ProtocolVersion, mark,
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(contentLength)
	if err != nil {
		return err
	}

	_, err = conn.Write(content)
	if err != nil {
		return err
	}
	return nil
}
