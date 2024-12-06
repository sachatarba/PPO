package service

import (
	"context"

	pb "github.com/sachatarba/course-db/internal/api/grpc"
	"github.com/sachatarba/course-db/internal/service"

	"github.com/google/uuid"
	"github.com/sachatarba/course-db/internal/entity"
)

type GRPCClientServer struct {
	service service.IClientService
	pb.UnimplementedClientServiceServer
}

func NewGRPCClientServer(service service.IClientService) *GRPCClientServer {
	return &GRPCClientServer{
		service: service,
	}
}

func (s *GRPCClientServer) ChangeClient(ctx context.Context, req *pb.Client) (*pb.Empty, error) {
	client := entity.Client{
		ID:        uuid.MustParse(req.Id),
		Fullname:  req.Fullname,
		Login:     req.Login,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Birthdate: req.Birthdate,
	}

	err := s.service.ChangeClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *GRPCClientServer) RegisterNewClient(ctx context.Context, req *pb.Client) (*pb.Empty, error) {
	client := entity.Client{
		ID:        uuid.MustParse(req.Id),
		Fullname:  req.Fullname,
		Login:     req.Login,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		Birthdate: req.Birthdate,
	}

	err := s.service.RegisterNewClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *GRPCClientServer) DeleteClient(ctx context.Context, req *pb.UUID) (*pb.Empty, error) {
	err := s.service.DeleteClient(ctx, uuid.MustParse(req.Value))
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

// func (s *GRPCClientServer) GetClientByID(ctx context.Context, req *pb.UUID) (*pb.ClientResponse, error) {
// 	client, err := s.service.GetClientByID(ctx, uuid.MustParse(req.Value))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &pb.ClientResponse{
// 		Client: &pb.Client{
// 			Id:       client.ID.String(),
// 			Fullname: client.Fullname,
// 			Login:    client.Login,
// 			Email:    client.Email,
// 			Birthdate: client.Birthdate,
// 			Phone: client.Phone,
// 			Password: client.Password,
// 		},
// 	}, nil
// }

func (s *GRPCClientServer) GetClientByLogin(ctx context.Context, req *pb.LoginRequest) (*pb.ClientResponse, error) {
	client, err := s.service.GetClientByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	return &pb.ClientResponse{
		Client: &pb.Client{
			Id:        client.ID.String(),
			Fullname:  client.Fullname,
			Login:     client.Login,
			Email:     client.Email,
			Birthdate: client.Birthdate,
			Phone:     client.Phone,
			Password:  client.Password,
		},
	}, nil
}

// func (s *GRPCClientServer) ListClients(ctx context.Context, req *pb.Empty) (*pb.ClientListResponse, error) {
// 	clients, err := s.service.ListClients(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var pbClients []*pb.Client
// 	for _, client := range clients {
// 		pbClients = append(pbClients, &pb.Client{
// 			Id:    client.ID,
// 			Name:  client.Name,
// 			Login: client.Login,
// 			Email: client.Email,
// 		})
// 	}

// 	return &pb.ClientListResponse{
// 		Clients: pbClients,
// 	}, nil
// }
