package http_srv

import (
	"context"
	"fmt"
	"net/http"

	"time"

	client "github.com/IsaqueB/ps-klever/pkg/client"
	"github.com/IsaqueB/ps-klever/proto"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HttpSrv interface {
	GetServer() *http.Server
}

type httpSvr struct {
	svr    *http.Server
	client client.Client
}

func NewHTTPServer() (HttpSrv, error) {
	grpc_client, err := client.NewGrpcClient()
	if err != nil {
		return nil, err
	}
	server := httpSvr{
		svr: &http.Server{
			Addr:        ":3000",
			ReadTimeout: 5 * time.Second,
		},
		client: grpc_client,
	}
	server.setRoutes()
	return &server, nil
}

func (s *httpSvr) GetServer() *http.Server {
	return s.svr
}

func (s *httpSvr) setServer(http_svr *http.Server) {
	s.svr = http_svr
}

func (s *httpSvr) setRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/insert", func(writer http.ResponseWriter, req *http.Request) {
		id, err := s.client.Insert(context.Background(), &proto.VoteStruct{
			Video:  primitive.NewObjectID().Hex(),
			User:   primitive.NewObjectID().Hex(),
			Upvote: false,
		})
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("%v", err)))
			return
		}
		writer.Write([]byte(id))
	})
	s.svr.Handler = r
}

// func insertRoute(writer http.ResponseWriter, req *http.Request) {
// 	writer.Write([]byte(insert()))
// }
