package srv

import "go.uber.org/fx"

// "log"
// "net"
// "go.uber.org/fx"
// "google.golang.org/grpc/reflection"
// "google.golang.org/grpc"

type GrpcSrv struct {
	// server *grpc.Server
}

func NewGrpcSrv(lc fx.Lifecycle) *GrpcSrv {
	// //举例：
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()

	// lc.Append(fx.Hook{
	// 	OnStart: func(ctx context.Context) error {
	// 		pb.RegisterGreeterServer(s, &server{})
	// 		reflection.Register(s)
	// 		if err := s.Serve(lis); err != nil {
	// 			log.Fatalf("failed to serve: %v", err)
	// 			return err
	// 		}
	// 		return nil
	// 	},
	// 	OnStop: func(ctx context.Context) error {
	// 		return s.Shutdown(ctx)
	// 	},
	// })

	// return GrpcSrv{server: s}
	return nil
}
