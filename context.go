// Copyright 2020 Ye Zi Jie.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2020/10/12 23:39:30

package vex

import "net"

type Context struct {
	conn    net.Conn
	req *request
}

func newContext(conn net.Conn, req *request) *Context {
	return &Context{
		conn: conn,
		req: req,
	}
}

func (c *Context) Write(p []byte) (n int, err error) {
	err = writeResponse(c.conn, okMark, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *Context) WriteError(msg string) error {
	err := writeResponse(c.conn, errorMark, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) Command() string {
	return c.req.command
}

func (c *Context) Args() [][]byte {
	return c.req.args
}

func (c *Context) Arg(i int) []byte {
	return c.req.args[i]
}

func (c *Context) StringArg(i int) string {
	return string(c.req.args[i])
}

func (c *Context) Version() uint8 {
	return c.req.version
}
