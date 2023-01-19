package main
import(
	serveDB "github.com/mrkovshik/grpc_vacancy_database/pkg/server"
	proto "github.com/mrkovshik/grpc_vacancy_database/grpc/proto"
"google.golang.org/grpc"
"log"
	"net"
)


func main (){
 // Create new gRPC server instance
 s := grpc.NewServer()

 // Put function that return ADD
 srv := &serveDB.GRPCServer{}

 // Register gRPC server for handle
 proto.RegisterDBServerServer(s, srv)

 // Listen on port 8080
 l, err := net.Listen("tcp", ":8080")
 if err != nil {
	 log.Fatal(err)
 }

 // Start gRPC server
 if err := s.Serve(l); err != nil {
	 log.Fatal(err)
 }
}