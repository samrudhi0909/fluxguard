package main

import (
	"context"
	"log"
	"time"

	pb "fluxguard/api" // Import generated proto code
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. Connect to FluxGuard Server
	// We use "WithTransportCredentials(insecure)" because we are on localhost (no SSL)
	// OLD: conn, err := grpc.NewClient("localhost:50051", ...)
	// NEW:
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewRateLimiterClient(conn)

	// 2. Define the Test Case
	// "User Sam" has a bucket size of 5 tokens. Refill rate is 1 token/sec.
	request := &pb.Request{
		UserId:     "user:sam",
		Capacity:   5,
		RefillRate: 1,
	}

	// 3. Simulate an "Attack" (10 requests instantly)
	log.Println("--- STARTING BURST TRAFFIC ---")
	for i := 1; i <= 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		
		res, err := client.ShouldAllow(ctx, request)
		if err != nil {
			log.Fatalf("RPC Error: %v", err)
		}
		
		// Print the result
		if res.Allowed {
			log.Printf("Request %d: ✅ ALLOWED", i)
		} else {
			log.Printf("Request %d: 🛑 BLOCKED (Rate Limit Exceeded)", i)
		}
		
		cancel()
	}
	log.Println("--- END OF BURST ---")
}
