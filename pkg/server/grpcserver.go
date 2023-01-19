package serveDB

import (
	proto "github.com/mrkovshik/grpc_vacancy_database/grpc/proto"
"context"
)
type GRPCServer struct {
}

func (s *GRPCServer) Get (ctx context.Context, req *proto.GetRequest) (*proto.GetResponse, error) {
	return &proto.GetResponse{
		GetResult: GetFunction(req.GetQuery()),
		}, nil
}

func (s *GRPCServer) Put (ctx context.Context, req *proto.PutRequest) (*proto.PutResponse, error) {
	return &proto.PutResponse{
PutResult: PutFunction(req.GetVacName(),req.GetVacDesc(),req.GetKeySkills(),int(req.GetSalary()),int(req.GetJobCode())),
		}, nil
}
