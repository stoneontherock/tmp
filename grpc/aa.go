package grpc

import (
	"aa/config"
	"aa/grpc/pb"
	"aa/httpapi"
	"aa/panicerr"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

type Server config.GRPCConf

func (s *Server) Authenticate(ctx context.Context, in *pb.AuthenticateIn) (*pb.AuthenticateOut, error) {
	claims, err := httpapi.ParseToken(in.Token)
	if err != nil {
		return &pb.AuthenticateOut{}, errors.New("parseToken:" + err.Error())
	}

	var dom string
	if in.Domain == "" {
		dom = claims.Domain
	} else {
		dom = in.Domain
	}

	if !httpapi.Enforcer.Enforce(claims.Username, dom+"@"+in.Resource, in.Action) {
		logrus.Errorf("user=%s dom@rsc=%s action=%s 没有权限", claims.Username, dom+"@"+in.Resource, in.Action)
		return &pb.AuthenticateOut{}, fmt.Errorf("无权限(u=%s)", claims.Username)
	}

	return &pb.AuthenticateOut{}, nil
}

func (s *Server) Authorize(ctx context.Context, in *pb.AuthorizeIn) (*pb.AuthorizeOut, error) {
	token, err := httpapi.CheckUser(in.UserName, in.Password)
	if err != nil {
		return &pb.AuthorizeOut{}, errors.New("CheckUser:" + err.Error())
	}

	return &pb.AuthorizeOut{Token: token}, nil
}

func Serve() {
	lis, err := net.Listen("tcp", config.C.GRPC.Addr)
	panicerr.PE(err, "创建grpc监听socket失败")

	cred, err := credentials.NewServerTLSFromFile(config.C.GRPC.Certificate, config.C.GRPC.Key)
	panicerr.PE(err, "读取grpc的TLS配置失败")

	instance := grpc.NewServer(grpc.Creds(cred))
	pb.RegisterAAServer(instance, (*Server)(&config.C.GRPC))
	err = instance.Serve(lis)
	panicerr.PE(err, "grpc服务端实例启动服务异常")
}
