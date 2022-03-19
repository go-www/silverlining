package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-www/h1"
	"github.com/go-www/silverlining"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lemon-mint/envaddr"
)

//var json = jsoniter.ConfigFastest

var (
	prefork = flag.Bool("prefork", false, "use prefork")
)

func main() {
	flag.Parse()

	data := []byte("Hello, World!")
	healthz := []byte("OK")
	jsonData := map[string]string{"hello": "world"}

	Handler := func(r *silverlining.Context) {
		switch string(r.Path()) {
		case "/":
			r.ResponseHeaders().Set("Content-Type", "text/plain")
			r.WriteFullBody(200, data)
		case "/json":
			r.WriteJSON(200, jsonData)
		case "/healthz":
			r.ResponseHeaders().Set("Content-Type", "text/plain")
			r.WriteFullBody(200, healthz)
		case "/redirect":
			r.Redirect(http.StatusSeeOther, "/")
		case "/bind_query":
			type User struct {
				Name string `query:"name"`
				Age  uint8  `query:"age"`
			}
			u := &User{}
			if err := r.BindQuery(u); err != nil {
				r.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
				return
			}

			r.WriteJSONIndent(200, u, "", "  ")
		case "/ws":
			conn, err := r.UpgradeWebSocket(ws.OpBinary)
			if err != nil {
				//println(err.Error())
				r.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
				return
			}

			go func() {
				defer conn.Close()
				for {
					msg, op, err := wsutil.ReadClientData(conn)
					if err != nil {
						return
					}

					if err := wsutil.WriteServerMessage(conn, op, msg); err != nil {
						return
					}
				}
			}()
		case "/httpbin":
			// inspired by https://httpbin.org/

			Origin, ok := r.RequestHeaders().Get("Origin")
			if !ok {
				Origin = "*"
			}
			r.ResponseHeaders().Set("Access-Control-Allow-Origin", Origin)
			r.ResponseHeaders().Set("Access-Control-Allow-Credentials", "true")
			r.ResponseHeaders().Set("Vary", "Origin")

			// Handle CORS preflight request
			if r.Method() == h1.MethodOPTIONS {
				RequestMethod, ok := r.RequestHeaders().Get("Access-Control-Request-Method")
				if ok {
					r.ResponseHeaders().Set("Access-Control-Allow-Methods", RequestMethod)
				}

				RequestHeaders, ok := r.RequestHeaders().Get("Access-Control-Request-Headers")
				if ok {
					r.ResponseHeaders().Set("Access-Control-Allow-Headers", RequestHeaders)
				}

				r.ResponseHeaders().Set("Access-Control-Max-Age", "86400")

				r.WriteFullBody(http.StatusNoContent, nil)
				return
			}

			qps := r.QueryParams()
			hs := r.RequestHeaders().List()

			type HttpRequest struct {
				Method string `json:"method"`

				Args    map[string]string      `json:"args"`
				Data    string                 `json:"data"`
				JSON    map[string]interface{} `json:"json"`
				Form    map[string]string      `json:"form"`
				Headers map[string]string      `json:"headers"`
			}

			reqData := HttpRequest{
				Method:  r.Method().String(),
				Args:    make(map[string]string),
				Headers: make(map[string]string),
			}

			for _, h := range hs {
				reqData.Headers[string(h.Name)] = string(h.RawValue)
			}

			for _, qp := range qps {
				reqData.Args[string(qp.Key)] = string(qp.Value)
			}

			body, err := r.FastBodyUnsafe(1024 * 1024 * 5) // 5MB
			if err != nil {
				r.WriteJSONIndent(500, map[string]string{"error": err.Error()}, "", "  ")
				return
			}
			reqData.Data = string(body)
			json.Unmarshal(body, &reqData.JSON)

			qfurl := h1.ParseRawQuery(body, nil)
			reqData.Form = make(map[string]string)
			for _, qp := range qfurl {
				reqData.Form[string(qp.Key)] = string(qp.Value)
			}

			r.WriteJSONIndent(200, reqData, "", "  ")
		case "/chunked":
			w := r.ChunkedBodyWriter()
			defer w.Close()
			w.WriteString("Hello, World!")
		case "/sse":
			w := r.ServerSentEventWriter()
			var i int
			err := w.WriteHeartbeat()
			if err != nil {
				return
			}
			for {
				i++
				err := w.Send(strconv.Itoa(i), "time", time.Now().Format(time.RFC3339))
				if err != nil {
					break
				}
				time.Sleep(time.Second)
			}
		case "/kill":
			r.KillConn()
		default:
			r.WriteFullBody(404, nil)
		}
	}

	var err error
	if *prefork {
		var id int
		id, err = silverlining.PreforkChildID()
		if err != nil {
			log.Fatalln(err)
		}

		if id == 0 {
			log.Println("Starting prefork leader process")
		} else {
			log.Printf("Starting prefork replica process %d", id)
		}
		err = silverlining.ListenAndServePrefork(envaddr.Get(":8080"), Handler)
	} else {
		err = silverlining.ListenAndServe(envaddr.Get(":8080"), Handler)
	}
	if err != nil {
		log.Fatalln(err)
	}
}
