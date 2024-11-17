package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authGRPC "sso/internal/grpc/auth"

	//"time"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
	//storagePath string,
	//tokenTTL time.Duration,
) *App {
	gRPCSever := grpc.NewServer()
	authGRPC.Register(gRPCSever)

	return &App{
		log:        log,
		gRPCServer: gRPCSever,
		port:       port,
	}
}

// MustRun runs the gRPC server and panics if an error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs the gRPC server.
func (a *App) Run() error {
	const operation = "grpcapp.Run"

	log := a.log.With(
		slog.String("operation", operation),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s, %w", operation, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s, %w", operation, err)
	}

	return nil
}

// Stop stops the gRPC server
func (a *App) Stop() {
	const operation = "grpcapp.Stop"

	a.log.With(
		slog.String("operation", operation),
	).Info("stopping grpc server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
