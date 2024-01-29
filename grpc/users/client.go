package users

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type UserGrpcInterface interface {
	Authorize() (*AuthorizeUserResponse, error)
	Get() (*GetUserResponse, error)
}

type UserGrpc struct {
	UserGrpcInterface
}

func Get(log *slog.Logger, id int64) (*GetUserResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("could not connect", slog.Any("err", err))
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Error("could not connect", slog.Any("err", err))
			return
		}
	}(conn)

	c := NewUserClient(conn)

	response, err := c.Get(context.Background(), &GetUserRequest{
		Id: id,
	})

	if err != nil {
		fmt.Println(err)
		log.Error("could not connect", slog.Any("err", err))
		return nil, err
	}

	return response, nil
}

func Authorize(log *slog.Logger, token string) (*AuthorizeUserResponse, error) {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("could not connect", slog.Any("err", err))
		return nil, err
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Error("could not connect", slog.Any("err", err))
			return
		}
	}(conn)

	c := NewUserClient(conn)

	response, err := c.Authorize(context.Background(), &AuthorizeUserRequest{
		Token: token,
	})

	if err != nil {
		log.Error("could not connect", slog.Any("err", err))
		return nil, err
	}

	return response, nil
}
