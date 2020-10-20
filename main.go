package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"golang.org/x/crypto/ed25519"
)

const (
	serviceName = "serviceName"
	accessToken = "access_token"
	mainJS      = "dist/main.js"
	replaceME   = "replaceME"
)

var (
	url       string
	svcName   string
	otel      string
	port      int
	js        []byte
	emptySeed = make([]byte, 32)
	upGrader  = websocket.Upgrader{
		Subprotocols: []string{accessToken},

		// CheckOrigin returns true if the request Origin header is acceptable. If
		// CheckOrigin is nil, then a safe default is used: return false if the
		// Origin request header is present and the origin host is not equal to
		// request Host header.
		//
		// A CheckOrigin function should carefully validate the request origin to
		// prevent cross-site request forgery.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},

		EnableCompression: true,
	}
)

func init() {
	flag.StringVar(&url, "ws", "", "websocket url, ie ws://host/ws")
	flag.StringVar(&svcName, "svc", "ui-service", "UI service name")
	flag.StringVar(&otel, "otel", "otel:55671", "OpenTelemetry collector HTTP end point, ie, host:port")
	flag.IntVar(&port, "p", 8999, "listening port, default 8999")
	var err error
	js, err = ioutil.ReadFile(mainJS)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	otel = fmt.Sprintf("http://%s/v1/trace", otel)

	public, private, err := ed25519.GenerateKey(bytes.NewReader(emptySeed))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.File("/main.js", mainJS, returnMainJS(private))
	e.GET("/ws", handleWS(public))
	e.Logger.Info("forward to collector", "endpoint", otel)
	e.Logger.Fatal(e.Start(":" + strconv.FormatInt(int64(port), 10)))
}

func handleWS(key ed25519.PublicKey) echo.HandlerFunc {
	cnt := uint64(0)
	return func(c echo.Context) error {
		cid := atomic.AddUint64(&cnt, 1)
		c.Logger().Info("new incoming request ", cid)
		conn, err := upGrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer func() {
			conn.Close()
			c.Logger().Info("request disconnected ", cid)
		}()

		if !authenticate(c.Request(), key) {
			http.Error(c.Response(), http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return http.ErrServerClosed
		}

		var (
			mt  int
			msg []byte
		)
		for {
			mt, msg, err = conn.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					err = nil
				}
				break
			} else if mt == websocket.TextMessage {
				forward(c, msg)
			}
		}

		return err
	}
}

func returnMainJS(key ed25519.PrivateKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			svc, ok := c.Request().URL.Query()[serviceName]
			if !ok || len(svc) == 0 {
				svc = []string{svcName}
			}
			tjs, err := loadMainJs(svc[0], key)
			if err != nil {
				return err
			}
			c.Response().Header().Add("Cache-Control", "no-cache")
			return c.Blob(http.StatusOK, echo.MIMEApplicationJavaScript, tjs)
		}
	}
}

func loadMainJs(svc string, key ed25519.PrivateKey) ([]byte, error) {
	tmpl, err := template.New(path.Base(mainJS)).ParseFiles(mainJS)
	if err != nil {
		return nil, err
	}

	w := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(w, struct {
		ServiceName string
		WsUrl       string
		AccessToken string
	}{
		ServiceName: svc,
		WsUrl:       url,
		AccessToken: hex.EncodeToString(ed25519.Sign(key, []byte(replaceME))),
	})
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func authenticate(req *http.Request, key ed25519.PublicKey) bool {
	subs := websocket.Subprotocols(req)
	for i := 0; i < len(subs); i++ {
		if strings.Compare(subs[i], accessToken) == 0 {
			i++
			if i < len(subs) {
				if sig, err := hex.DecodeString(subs[i]); err == nil {
					return ed25519.Verify(key, []byte(replaceME), sig)
				}
			}
		}
	}

	return false
}
