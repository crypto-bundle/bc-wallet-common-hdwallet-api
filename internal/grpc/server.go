/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package grpc

import (
	"context"
	"errors"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	pbApi "github.com/crypto-bundle/bc-wallet-common-hdwallet-controller/pkg/grpc/hdwallet"

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
	//options := commonGRPCServer.DefaultServeOptions()
	msgSizeOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(1024 * 1024 * 3),
		grpc.MaxSendMsgSize(1024 * 1024 * 3),
	}
	//options = append(options, msgSizeOptions...)
	//options = append(options, grpc.StatsHandler(otelgrpc.NewServerHandler()))

	s.grpcServerOptions = msgSizeOptions

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

	resolved, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return err
	}

	listenConn, err := net.ListenUnix("unix", resolved)
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
