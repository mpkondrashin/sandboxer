package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "examen/pkg/grpc/proto"
	"examen/pkg/state"
	"examen/pkg/task"
)

func Submit(path string) error {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewExamenSvcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Submit(ctx, &pb.SubmitRequest{Name: path})
	if err != nil {
		return err
	}
	if r.Message != "Ok" {
		return errors.New(r.Message)
	}
	return nil
}

type StatusCallbackFunc func(*task.Task)

func Status(From int32, Count int32, callback StatusCallbackFunc) error {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewExamenSvcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Status(ctx, &pb.StatusRequest{
		From:  From,
		Count: Count,
	})
	if err != nil {
		return err
	}
	for _, t := range r.Tasks {
		task := task.NewTask(t.Path)
		task.SetState(state.State(t.Status))
		callback(task)
	}
	return nil
}
