package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"campus-delivery/api"
	"campus-delivery/clientApi"
)

const network = "tcp"

type GrpcServer struct {
	port       int
	controller clientApi.Controller
}

func NewServer(port int, controller clientApi.Controller) *GrpcServer {
	return &GrpcServer{
		port:       port,
		controller: controller,
	}
}

func (srv *GrpcServer) RegistrationUser(_ context.Context, request *api.RegistrationRequest) (*api.CodeResponse, error) {
	codeResp := srv.controller.Registration(*request)
	return &codeResp, nil
}

func (srv *GrpcServer) DeleteUser(_ context.Context, request *api.DeleteRequest) (*api.CodeResponse, error) {
	codeResp := srv.controller.DeleteUser(*request)
	return &codeResp, nil
}

func (srv *GrpcServer) AddCourier(_ context.Context, request *api.AddCourierRequest) (*api.NewCourier, error) {
	newCourier := srv.controller.AddUserWithOrder(*request)
	return &newCourier, nil
}

func (srv *GrpcServer) DeleteCourier(_ context.Context, request *api.DeleteRequest) (*api.CodeResponse, error) {
	codeResp := srv.controller.DeleteUserWithOrder(*request)
	return &codeResp, nil
}

func (srv *GrpcServer) ListenCourier(_ *api.UserRequest, stream api.DeliveryServer_ListenCourierServer) error {
	//go srv.controller.ListenCourier()
	//uc := <-srv.controller.(*controller.Controller).NewCourier
	//if err := stream.Send(&uc); err != nil {
	//	log.Printf("Error to send: %v", err)
	//}

	return nil
}

func (srv *GrpcServer) GetCourier(_ context.Context, request *api.UserRequest) (*api.GetCourierResponse, error) {
	couriers := srv.controller.GetCourier(*request)
	return &couriers, nil
}

func (srv *GrpcServer) CheckCourier(_ context.Context, request *api.UserRequest) (*api.CodeResponse, error) {
	codeResp := srv.controller.CheckCourier(*request)
	return &codeResp, nil
}

func (srv *GrpcServer) ChangeLocation(_ context.Context, request *api.RegistrationRequest) (*api.CodeResponse, error) {
	codeResp := srv.controller.ChangeLocation(*request)
	return &codeResp, nil
}

func (srv *GrpcServer) ChangeNotificationStatus(_ context.Context, request *api.RegistrationRequest) (*api.CodeResponse, error) {
	codeResp := srv.controller.ChangeNotificationStatus(*request)
	return &codeResp, nil
}

func (srv *GrpcServer) GetUserInfo(_ context.Context, request *api.UserRequest) (*api.RegistrationRequest, error) {
	userResp := srv.controller.GetUserInfo(*request)
	return &userResp, nil
}

func (srv *GrpcServer) GetUsers(_ context.Context, request *api.UserRequest) (*api.Users, error) {
	users := srv.controller.GetAllUser(*request)
	return &users, nil
}

func (srv *GrpcServer) GetAllNotificationUser(_ context.Context, request *api.UserRequest) (*api.Users, error) {
	users := srv.controller.GetAllNotificationUser(*request)
	return &users, nil
}

func (srv *GrpcServer) AddRating(context.Context, *api.RatingRequest) (*api.CodeResponse, error) {
	panic("implement me")
}

func (srv *GrpcServer) Run(opts ...grpc.ServerOption) error {
	// register service
	server := grpc.NewServer(opts...)
	api.RegisterDeliveryServerServer(server, srv)

	// graceful shutdown
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sch {
			// sig is a ^Conf, handle it
			log.Println("shutting down gRPC service...")
			server.GracefulStop()
		}
	}()

	address := fmt.Sprintf(":%d", srv.port)
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("failed to listener: %v", err)
	}

	// start gRPC service
	log.Printf("starting gRPC service on port %s", address)
	return server.Serve(listener)
}
