package web

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	handl          *GameHandler
	authHandler    *AuthHandler
	userHandler    *UserHandler
	authMiddleware Authenticator
	port           string
	server         *http.Server
}

func NewServer(handl *GameHandler, authHandler *AuthHandler, userHandler *UserHandler, authMiddleware Authenticator, port string) *Server {
	return &Server{
		handl:          handl,
		authHandler:    authHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
		port:           port,
		server: &http.Server{
			Addr: ":" + port,
		},
	}
}

func (s *Server) DefiningMethods() {
	http.HandleFunc("POST /sign-up", s.authHandler.SignUp)
	http.HandleFunc("POST /sign-in", s.authHandler.SignIn)

	http.Handle("POST /game", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.CreateGame)))
	http.Handle("GET /game/list", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.ListAvailableGames)))
	http.Handle("GET /game/{id}", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.GetGame)))
	http.Handle("POST /game/{id}/move", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.MakeMove)))
	http.Handle("POST /game/{id}/join", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.JoinGame)))
	http.Handle("GET /users/{id}", s.authMiddleware.Middleware(http.HandlerFunc(s.userHandler.GetUser)))
	http.Handle("GET /user/current", s.authMiddleware.Middleware(http.HandlerFunc(s.userHandler.GetCurrentUser)))
	http.Handle("GET /games/history", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.GetCompletedGames)))
	http.Handle("GET /leaderboard", s.authMiddleware.Middleware(http.HandlerFunc(s.handl.GetTopPlayers)))
	http.HandleFunc("POST /refresh-access", s.authHandler.RefreshAccessToken)
	http.HandleFunc("POST /refresh-refresh", s.authHandler.RefreshRefreshToken)

}

func (s *Server) Start() error {
	s.DefiningMethods()
	log.Printf("Запуск %s", s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	s.server.Shutdown(ctx)
}
