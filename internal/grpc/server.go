package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

	commonGRPCServer "github.com/crypto-bundle/bc-wallet-common-lib-grpc/pkg/server"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	logger            *zap.Logger
	grpcServer        *grpc.Server
	grpcServerOptions []grpc.ServerOption
	handlers          pbApi.HdWalletApiServer
	configSvc         configService

	sockFilePath string
	listener     net.Listener
}

func (s *Server) Init(_ context.Context) error {
	options := commonGRPCServer.DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(commonGRPCServer.DefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(commonGRPCServer.DefaultServerMaxSendMessageSize),
	}
	options = append(options, msgSizeOptions...)
	options = append(options, grpc.StatsHandler(otelgrpc.NewServerHandler()))

	s.grpcServerOptions = options

	return nil
}

func (s *Server) shutdown() error {
	s.logger.Info("start close instances")

	s.grpcServer.GracefulStop()

	err := os.Remove(s.sockFilePath)
	if err != nil {
		return err
	}

	s.logger.Info("grpc server shutdown completed")

	return nil
}

func (s *Server) ListenAndServe(ctx context.Context) (err error) {
	tf, err := os.CreateTemp(s.configSvc.GetConnectionPath(), s.configSvc.GetUnixFileNameTemplate())
	if err != nil {
		return err
	}

	path := tf.Name()

	// Close the file and remove it because it has to not exist for
	// the domain socket.
	err = tf.Close()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	listenConn, err := net.Listen("unix", path)
	if err != nil {
		s.logger.Error("unable to listen", zap.Error(err),
			zap.String("path", s.configSvc.GetConnectionPath()))

		return err
	}

	s.sockFilePath = path
	s.listener = listenConn

	s.grpcServer = grpc.NewServer(s.grpcServerOptions...)
	if (s.configSvc.IsDev() || s.configSvc.IsLocal()) && s.configSvc.IsDebug() {
		reflection.Register(s.grpcServer)
	}

	go s.serve(ctx)

	return nil
}

func (s *Server) serve(ctx context.Context) {
	newCtx, causeFunc := context.WithCancelCause(ctx)
	pbApi.RegisterHdWalletApiServer(s.grpcServer, s.handlers)

	s.logger.Info("grpc serve success")
	go func() {
		err := s.grpcServer.Serve(s.listener)
		if err != nil {
			s.logger.Error("unable to start serving", zap.Error(err),
				zap.String("path", s.configSvc.GetConnectionPath()))

			causeFunc(err)
		}
	}()

	<-newCtx.Done()
	intErr := newCtx.Err()
	if !errors.Is(intErr, context.Canceled) {
		s.logger.Error("ctx cause errors", zap.Error(intErr))
	}

	err := s.shutdown()
	if err != nil {
		s.logger.Error("unable to graceful shutdown", zap.Error(err))
	}

	return
}

// nolint:revive // fixme
func NewServer(ctx context.Context,
	loggerSrv *zap.Logger,
	cfg configService,
	handlers pbApi.HdWalletApiServer,
) (*Server, error) {
	l := loggerSrv.Named("grpc.server")

	srv := &Server{
		logger:    l,
		configSvc: cfg,
		handlers:  handlers,
	}

	return srv, nil
}
