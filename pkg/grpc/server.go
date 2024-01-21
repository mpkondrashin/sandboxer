package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "examen/pkg/grpc/proto"
	"examen/pkg/task"
)

const GRPCPort = 10034

type SubmitFunc func(name string) error
type StatusFunc func(from, count int32) ([]*task.Task, error)

type server struct {
	submit SubmitFunc
	status StatusFunc
	pb.UnimplementedExamenSvcServer
}

var ErrNotImpelmentd = errors.New("not implemented")

func (s *server) Status(ctx context.Context, in *pb.StatusRequest) (*pb.StatusReply, error) {
	if s.status == nil {
		return nil, ErrNotImpelmentd
	}
	tasks, err := s.status(in.GetFrom(), in.GetCount())
	if err != nil {
		return nil, err
	}
	reply := &pb.StatusReply{}
	for _, t := range tasks {
		pbTask := &pb.Task{
			Path:   t.Path,
			Status: int32(t.State),
		}
		reply.Tasks = append(reply.Tasks, pbTask)
	}
	return reply, nil
}

func (s *server) Submit(ctx context.Context, request *pb.SubmitRequest) (*pb.SubmitReply, error) {
	if s.submit == nil {
		return nil, ErrNotImpelmentd
	}
	err := s.submit(request.GetName())
	if err != nil {
		return &pb.SubmitReply{Message: err.Error()}, nil
	}
	return &pb.SubmitReply{Message: "Ok"}, nil
}

func RunServer(submit SubmitFunc, status StatusFunc) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", GRPCPort))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterExamenSvcServer(s, &server{
		submit: submit,
		status: status,
	})
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}
