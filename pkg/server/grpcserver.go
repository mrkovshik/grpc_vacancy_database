package serveDB

import (
	proto "github.com/mrkovshik/grpc_vacancy_database/grpc/proto"
"context"
)
type GRPCServer struct {
}

func (s *GRPCServer) Read (ctx context.Context, req *proto.ReadRequest) (*proto.ReadResponse, error) {
	return &proto.ReadResponse{
		ReadResult: ReadFunction(req.GetReadQuery()),
		}, nil
}

func (s *GRPCServer) Insert (ctx context.Context, req *proto.InsertRequest) (*proto.InsertResponse, error) {
	return &proto.InsertResponse{
InsertResult: InsertFunction(req.GetNewVac()),
		}, nil
}

func (s *GRPCServer) Delete (ctx context.Context, req *proto.DeleteRequest) (*proto.DeleteResponse, error) {
	return &proto.DeleteResponse{
		DeleteResult: DeleteFunction(req.GetDeleteTarget()),
		}, nil
}