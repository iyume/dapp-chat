package server

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/utils"
)

const maxBodySize = 1 * 1024 * 1024

type HTTPConfig struct {
	// Token is http authorization token. If this field is empty string, then generate
	// 16-byte random hex string
	Token string

	// Must be hostname:port
	Address string

	// TODO: add reverse http event poster
}

var DefaultHTTPConfig = HTTPConfig{
	Token:   "",
	Address: "127.0.0.1:8080",
}

type httpServer struct {
	HTTPConfig
	backend *api.Backend
	caller  Caller
}

type Getter interface {
	Get(string) any

	// GetString returns empty string if not string type
	GetString(string) string

	Has(string) bool

	Require(...string) string
}

type paramsGetter struct {
	json  map[string]any // could be lazy?
	query url.Values
}

func (p paramsGetter) Has(key string) bool {
	if _, ok := p.json[key]; ok || p.query.Has(key) {
		return true
	}
	return false
}

func (p paramsGetter) Get(key string) any {
	val, ok := p.json[key]
	if ok {
		return val
	}
	val = p.query.Get(key)
	if val != "" {
		return val
	}
	return nil
}

func (p paramsGetter) GetString(key string) string {
	val := p.Get(key)
	strval, _ := val.(string)
	return strval
}

// Require checks multiple parameters and returns the missing key
func (p paramsGetter) Require(params ...string) string {
	for _, key := range params {
		if !p.Has(key) {
			return key
		}
	}
	return ""
}

// main loop of http server
func StartHTTPServer(backend *api.Backend, config HTTPConfig) (*http.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, nil, err
	}
	var handler http.Handler = httpServer{
		HTTPConfig: config,
		backend:    backend,
		caller:     NewCaller(backend),
	}
	token := config.Token
	if token == "" {
		token = hex.EncodeToString(utils.RandomBytes(16))
	}
	handler = NewHTTPStack(handler, token)
	srv := &http.Server{
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Println("Started http server at", "http://"+listener.Addr().String(),
		"with token", token)
	go srv.Serve(listener)
	return srv, listener, nil
}

func (srv httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("handle http request", *r)
	getter := &paramsGetter{}
	switch r.Method {
	case http.MethodPost:
		contentType := r.Header.Get("Content-Type")
		// may charset but ignored
		if !strings.Contains(contentType, "application/json") {
			http.Error(w, "only support application/json", http.StatusBadRequest)
			return
		}
		jsonBytes, err := readBody(r)
		if err != nil {
			http.Error(w, "too large body", http.StatusBadRequest)
			return
		}
		if !bytes.HasPrefix(jsonBytes, []byte("{")) || !json.Valid(jsonBytes) {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(jsonBytes, &getter.json); err != nil {
			http.Error(w, "parameter invalid", http.StatusBadRequest)
			return
		}
		fallthrough // Support post URL query
	case http.MethodGet:
		getter.query = r.URL.Query()
	default:
		log.Printf("已拒绝 %s %s，请求方法错误\n", r.Method, r.RemoteAddr)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// call action
	action := strings.TrimPrefix(r.URL.Path, "/")
	log.Println("接收 API 调用:", action, *getter)
	resp := srv.caller.Call(action, getter)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("写入响应数据出错:", err)
		log.Println(resp)
	}
}

// read POST body with limited length, returns error if exceeded
func readBody(r *http.Request) ([]byte, error) {
	bytes, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		return nil, err
	}
	if len(bytes) >= maxBodySize {
		return nil, errors.New("exceeded max body size")
	}
	r.Body.Close()
	return bytes, err
}
