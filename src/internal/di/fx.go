package di

import (
	"context"
	"os"
	"tic-tac-toe/internal/application/auth"
	"tic-tac-toe/internal/application/service"
	"tic-tac-toe/internal/datasource"
	"tic-tac-toe/internal/db"
	"tic-tac-toe/internal/domain/repository"
	"tic-tac-toe/internal/web"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var DatabaseModule = fx.Module("databse",
	fx.Provide(func(lc fx.Lifecycle) (*pgxpool.Pool, error) {
		pool, err := db.Connect(context.Background())
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				pool.Close()
				return nil
			},
		})
		if err := db.RunMigrations(context.Background(), pool); err != nil {
			pool.Close()
			return nil, err
		}
		return pool, nil
	}),
)
var RepositoryModule = fx.Module("repository", fx.Provide(fx.Annotate(
	datasource.NewGameRepoDB,
	fx.As(new(repository.GameRepository)))))

var UserRepositoryModule = fx.Module("userrepo", fx.Provide(fx.Annotate(
	datasource.NewUserRepoDB,
	fx.As(new(repository.UserRepository)))))

var ServiceModule = fx.Module("service", fx.Provide(fx.Annotate(
	service.NewGameService,
	fx.As(new(service.GameService)))))

var UserServiceModule = fx.Module("userservice",
	fx.Provide(fx.Annotate(service.NewService, fx.As(new(service.UserService)))),
)

var AuthModule = fx.Module("auth",
	fx.Provide(
		fx.Annotate(auth.NewAuthService, fx.As(new(auth.AuthService))),
	),
)

var JwtModule = fx.Module("jwt", fx.Provide(func() *auth.JwtProvider {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		accessSecret = "default_access_secret"
	}
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "default_refresh_secret"
	}
	return auth.NewJwtProvider(
		accessSecret,
		refreshSecret,
		15*time.Minute,
		7*24*time.Hour,
	)
}),
	fx.Provide(fx.Annotate(
		func(p *auth.JwtProvider) *web.JwtAuthenticator {
			return web.NewJwtAuthenticator(p)
		},
		fx.As(new(web.Authenticator)),
	)),
)

var WebModule = fx.Module("web", fx.Provide(
	web.NewGameHandler,
	web.NewAuthHandler,
	web.NewUserHandler,
	web.NewServer,
),
	fx.Supply("8080"),
	fx.Invoke(RegisterServer),
)

var App = fx.Module("tic-tac-toe",
	DatabaseModule,
	RepositoryModule,
	UserRepositoryModule,
	ServiceModule,
	UserServiceModule,
	AuthModule,
	JwtModule,
	WebModule,
)

func RegisterServer(lc fx.Lifecycle, srv *web.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Stop(ctx)
			return nil
		},
	})
}
