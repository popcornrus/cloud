package users

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type UserGrpcInterface interface {
	Authorize() (*AuthorizeUserResponse, error)
}

type UserGrpc struct {
	UserGrpcInterface
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
