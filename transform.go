package uniottrans

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type TraceNative interface {
	Foreach(handler func(native SpanNative))
}

type SpanNative interface {
	GetTraceID() int64
	GetParentID() int64
	GetSpanID() int64
	GetService() string
	GetOperation() string
	GetMeta() map[string]string
	GetMetrics() map[string]*Numeric
	GetSpanStatus() SpanStatus
	GetStartTime() int64
	GetEndTime() int64
}

type Decoder func(body io.Reader) SpanNative

type SpanTransform interface {
	RegisterDecoder(contentType string, decoder Decoder)
	BuildSpanFrom(native SpanNative) Span
}

type SpanTransformServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (svr *SpanTransformServerConfig) newSpanTransformServer() *SpanTransformServer {

}

type SpanTransformServer struct {
	mux *http.ServeMux
}

func (sts *SpanTransformServer) BuildSpanFrom(native SpanNative) Span {
	svr := http.NewServeMux()
	http.ListenAndServe("", svr)
}

const (
	config_path = "OTTF_CONFIG_PATH"
	address     = "OTTF_ADDRESS"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	v, ok := os.LookupEnv(config_path)
	if !ok {
		v = "./config.json"
	}
	if _, err := os.Stat(v); os.IsNotExist(err) {
		log.Fatalln("config file not exists")
	}

	config := &SpanTransformServerConfig{}
	if bts, err := os.ReadFile(v); err != nil {
		log.Fatalln(err.Error())
	} else {
		if err = json.Unmarshal(bts, config); err != nil {
			log.Println("load config file failed with error: %s", err.Error())
		}
	}
}
