package server

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/James-Ren/Go-001/tree/main/Week04/api/article/v1"
	"github.com/James-Ren/Go-001/tree/main/Week04/internal/service"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	service *service.Service
}

func NewServer(s *service.Service) *Server {
	return &Server{service: s}
}

func (srv *Server) Run() error {
	lis, err := net.Listen("tcp", viper.GetString("grpc.port"))
	if err != nil {
		return err
	}
	g, ctx := errgroup.WithContext(context.Background())
	s := grpc.NewServer()
	g.Go(func() error {
		go func() {
			<-ctx.Done()
			s.GracefulStop()
			log.Printf("Shutdown Server")
		}()
		pb.RegisterArticleServer(s, srv.service)
		return s.Serve(lis)
	})
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			select {
			case <-ctx.Done():
				return nil
			case s := <-c:
				log.Printf("get a signal %s", s.String())
				switch s {
				case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
					return errors.New("Close by signal " + s.String())
				case syscall.SIGHUP:
				default:
					return errors.New("Undefined signal")
				}
			}
		}
	})
	return g.Wait()
}
