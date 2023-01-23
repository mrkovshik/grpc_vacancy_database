package serveDB

import (
	proto "github.com/mrkovshik/grpc_vacancy_database/grpc/proto"
"context"
)
type GRPCServer struct {
}

func (s *GRPCServer) Read (ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	return &proto.ReadResponse{
		ReadResult: ReadFunction(req.GetQuery()),
		}, nil
}

func (s *GRPCServer) Insert (ctx context.Context, req *proto.InsertRequest) (*proto.InsertResponse, error) {
	return &proto.InsertResponse{
InsertResult: InsertFunction(req.GetNewVac()),
		}, nil
}
