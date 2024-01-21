package users

import (
	"cloud/internal/root/model"
	"cloud/internal/root/service/user_service"
	"context"
	"github.com/patrickmn/go-cache"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strings"
)

type (
	UserGrpcServer struct {
		user user_service.UserService
		log  *slog.Logger
		ch   *cache.Cache
	}
)

func (u *UserGrpcServer) Get(ctx context.Context, request *GetUserRequest) (*GetUserResponse, error) {
	var userResponse = &GetUserResponse{}

	user, err := u.user.FindUserByUUID(ctx, request.Uuid)

	if err != nil {
		return nil, err
	}

	userResponse.Id = user.ID
	userResponse.Uuid = user.UUID
	userResponse.Email = user.Email
	userResponse.Name = user.Username
	userResponse.CreatedAt = user.CreatedAt.String()
	userResponse.UpdatedAt = user.UpdatedAt.String()

	return userResponse, nil
}

func (u *UserGrpcServer) Authorize(ctx context.Context, request *AuthorizeUserRequest) (*AuthorizeUserResponse, error) {
	token := strings.Replace(request.Token, "Bearer ", "", 1)

	if v, ok := u.ch.Get(token); ok {
		u.log.Info("user was found in cache", slog.Any("user", v))
	}

	user := _getUserFromCacheByToken(u.ch, token)
	if user == nil {
		u.log.Warn("user wasn't found")

		return nil, nil
	}

	return &AuthorizeUserResponse{
		Id:        user.ID,
		Uuid:      user.UUID,
		Email:     user.Email,
		Name:      user.Username,
		UpdatedAt: user.UpdatedAt.String(),
	}, nil
}

func (u *UserGrpcServer) mustEmbedUnimplementedUserServer() {
	panic("implement me")
}

func NewUserGRPC(
	u user_service.UserService,
	log *slog.Logger,
	ch *cache.Cache,
) *UserGrpcServer {
	return &UserGrpcServer{
		user: u,
		log:  log,
		ch:   ch,
	}
}

func _getUserFromCacheByToken(ch *cache.Cache, token string) *model.User {
	if v, ok := ch.Get(token); ok {
		if user, ok := v.(model.User); ok {
			return &user
		}

		slog.Error("failed to cast user from cache", slog.Any("user", v))
	}

	return nil
}

func RunUserGRPCServer(
	lc fx.Lifecycle,
	log *slog.Logger,
	ch *cache.Cache,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				lis, err := net.Listen("tcp", ":50051")
				if err != nil {
					panic(err)
				}

				grpcServer := grpc.NewServer()
				RegisterUserServer(grpcServer, &UserGrpcServer{
					log: log,
					ch:  ch,
				})

				if err := grpcServer.Serve(lis); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Error("Server stopped")
			return nil
		},
	})
}
