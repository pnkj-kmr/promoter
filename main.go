package main

import (
	"flag"
	"log"
	"net"

	// "net/http"
	// _ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/pnkj-kmr/promoter/handler"
	"github.com/pnkj-kmr/promoter/internal"
	"github.com/pnkj-kmr/promoter/internal/settings"
	"github.com/pnkj-kmr/promoter/medium/pb"
	"google.golang.org/grpc"
)

func main() {
	f := flag.String("c", "app", "configuration file")
	debug := flag.Bool("debug", false, "Application Debug Mode")
	flag.Parse()
	settings.Init(*f, *debug)

	promote := handler.New()
	grpcServer := grpc.NewServer()
	pb.RegisterPromoteServer(grpcServer, promote)
	// // TO DEBUG THE gRPC SERVICE with help to
	// // EVANS Client --- https://github.com/ktr0731/evans
	// reflection.Register(grpcServer)

	address := settings.C.Bind
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start the server agent ", err)
	}

	// Graceful shutdwon of server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("listen: %s\n", err)
			done <- syscall.SIGTERM
		}
	}()

	internal.InitiateLeader()
	internal.InitiateBroker()
	log.Printf("server started [%s]... debug[%v]", address, *debug)

	// // test func to enable the pprofiging
	// go func() {
	// 	if *f == "app" {
	// 		log.Printf("start http server at http://localhost:8000/ port")
	// 		http.ListenAndServe("localhost:8000", nil)
	// 	} else {
	// 		log.Printf("start http server at http://localhost:8001/ port")
	// 		http.ListenAndServe("localhost:8001", nil)
	// 	}
	// }()

	<-done
	grpcServer.GracefulStop()
	log.Print("server exited properly")

}
