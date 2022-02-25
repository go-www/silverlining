package main

import (
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/go-www/h1"
	"github.com/go-www/silverlining"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lemon-mint/envaddr"
)

func main() {
	ln, err := net.Listen("tcp", envaddr.Get(":8080"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on http://%s", ln.Addr())

	defer ln.Close()

	srv := silverlining.Server{
		//ReadTimeout: time.Minute,
	}

	data := []byte("Hello, World!")
	healthz := []byte("OK")
	jsonData := map[string]string{"hello": "world"}

	srv.Handler = func(r *silverlining.Context) {
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
		case "/httpbin/get":
			if r.Method() != h1.MethodGET {
				r.WriteFullBodyString(http.StatusMethodNotAllowed, "Method not allowed")
				return
			}

			qps := r.QueryParams()
			hs := r.RequestHeaders().List()

			type HttpRequest struct {
				Args    map[string]string `json:"args"`
				Headers map[string]string `json:"headers"`
			}

			reqData := HttpRequest{
				Args:    make(map[string]string),
				Headers: make(map[string]string),
			}

			for _, h := range hs {
				reqData.Headers[string(h.Name)] = string(h.RawValue)
			}

			for _, qp := range qps {
				v, err := url.QueryUnescape(string(qp.RawValue))
				if err != nil {
					continue
				}
				reqData.Args[string(qp.Key)] = v
			}

			r.WriteJSONIndent(200, reqData, "", "  ")
		default:
			r.WriteFullBody(404, nil)
		}
	}

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
