/*
 *
 * Copyright 2015, Google Inc.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     * Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following disclaimer
 * in the documentation and/or other materials provided with the
 * distribution.
 *     * Neither the name of Google Inc. nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 * A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

// Package main implements a simple gRPC server that demonstrates how to use gRPC-Go libraries
// to perform unary, client streaming, server streaming and full duplex RPCs.
//
// It implements the route guide service whose definition can be found in proto/route_guide.proto.
package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	pb "google.golang.org/grpc/examples/stream_fail/foo"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile  = flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
)

type fooServer struct {
}

// RecordRoute records a route composited of a sequence of points.
//
// It gets a stream of points, and responds with statistics about the "trip":
// number of points,  number of known features visited, total distance traveled, and
// total time spent.
// func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
// 	var pointCount, featureCount, distance int32
// 	var lastPoint *pb.Point
// 	startTime := time.Now()
// 	for {
// 		point, err := stream.Recv()
// 		if err == io.EOF {
// 			endTime := time.Now()
// 			return stream.SendAndClose(&pb.RouteSummary{
// 				PointCount:   pointCount,
// 				FeatureCount: featureCount,
// 				Distance:     distance,
// 				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
// 			})
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		pointCount++
// 		for _, feature := range s.savedFeatures {
// 			if proto.Equal(feature.Location, point) {
// 				featureCount++
// 			}
// 		}
// 		if lastPoint != nil {
// 			distance += calcDistance(lastPoint, point)
// 		}
// 		lastPoint = point
// 	}
// }

func (s *fooServer) PointChat(stream pb.Foo_PointChatServer) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				return
			}
			if err != nil {
				return
			}
			grpclog.Println(in)
		}
	}()
	for {
		point := randomPoint(r)
		if err := stream.Send(point); err != nil {
			grpclog.Printf("%v.Send(%v) = %v", stream, point, err)
			return nil

		}
	}

}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterFooServer(grpcServer, new(fooServer))
	grpcServer.Serve(lis)
}
