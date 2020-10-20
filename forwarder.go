package main

import (
	"bytes"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/labstack/echo"
)

var (
	client = &http.Client{Timeout: time.Second * 5}
)

func forward(c echo.Context, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			c.Logger().Error("%s", string(debug.Stack()))
		}
	}()
	r, err := client.Post(otel, echo.MIMEApplicationJSON, bytes.NewReader(data))
	if err != nil {
		c.Logger().Error("err forwarding", err)
		return
	}
	r.Body.Close()
	c.Logger().Info("forward status = ", r.StatusCode, ", data size = ", len(data))
}
