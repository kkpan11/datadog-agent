// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

// Package helpers provides utility functions for gRPC servers.
package helpers

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	grpccontext "github.com/DataDog/datadog-agent/pkg/util/grpc/context"
)

// NewMuxedGRPCServer returns an http.Server that multiplexes connections
// between a gRPC server and an HTTP handler.
func NewMuxedGRPCServer(addr string, tlsConfig *tls.Config, grpcServer http.Handler, httpHandler http.Handler, timeout time.Duration) *http.Server {
	// our gRPC clients do not handle protocol negotiation, so we need to force
	// HTTP/2

	if timeout > 0 {
		httpHandler = TimeoutHandlerFunc(httpHandler, timeout)
	}

	var handler http.Handler
	// when HTTP/2 traffic that is not TLS being sent we need to create a new handler which
	// is able to handle the pre fix magic sent
	handler = h2c.NewHandler(handlerWithFallback(grpcServer, httpHandler), &http2.Server{})
	if tlsConfig != nil {
		tlsConfig.NextProtos = []string{"h2"}
		handler = handlerWithFallback(grpcServer, httpHandler)
	}

	return &http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: tlsConfig,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			// Store the connection in the context so requests can reference it if needed
			return context.WithValue(ctx, grpccontext.ConnContextKey, c)
		},
	}
}

// TimeoutHandlerFunc returns an HTTP handler that times out after a duration.
// This is useful for muxed gRPC servers where http.Server cannot have a
// timeout when handling streaming, long running connections.
func TimeoutHandlerFunc(httpHandler http.Handler, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deadline := time.Now().Add(timeout)

		conn := r.Context().Value(grpccontext.ConnContextKey).(net.Conn)
		_ = conn.SetWriteDeadline(deadline)

		httpHandler.ServeHTTP(w, r)
	})
}

// handlerWithFallback returns an http.Handler that delegates to grpcServer on
// incoming gRPC connections or httpServer otherwise. Copied from
// cockroachdb.
func handlerWithFallback(grpcServer http.Handler, httpServer http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks
		// https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpServer.ServeHTTP(w, r)
		}
	})
}
