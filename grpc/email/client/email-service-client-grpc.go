package proto

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oceano-dev/microservices-go-common/config"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmailServiceClientGrpc struct {
	config *config.Config
}

func NewEmailServiceClientGrpc(
	config *config.Config,
) *EmailServiceClientGrpc {
	return &EmailServiceClientGrpc{
		config: config,
	}
}

var grpcClient EmailServiceClient

func (s *EmailServiceClientGrpc) SendPasswordCode(email string, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, span := trace.NewSpan(ctx, "emailServiceGrpc.SendPasswordCodeReq")
	defer span.End()

	log.Println("verifyClientGrpc start")

	s.verifyClientGrpc()

	log.Println("verifyClientGrpc end")

	fmt.Println(grpcClient)

	req := &PasswordCodeReq{
		Email: email,
		Code:  code,
	}

	validator := validator.New()
	if err := validator.StructCtx(ctx, req); err != nil {
		trace.AddSpanError(span, err)
		log.Printf("emailServiceGrpc.SendPasswordCodeReq: %v", err)
		return err
	}

	out, err := grpcClient.SendPasswordCode(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("out: %v", out)

	log.Println("email sent")

	return nil
}

func (s *EmailServiceClientGrpc) SendSupportMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, span := trace.NewSpan(ctx, "emailServiceGrpc.SendSupportMessageReq")
	defer span.End()

	s.verifyClientGrpc()

	req := &SupportMessageReq{
		Message: message,
	}

	validator := validator.New()
	if err := validator.StructCtx(ctx, req); err != nil {
		trace.AddSpanError(span, err)
		log.Printf("emailServiceGrpc.SendSupportMessageReq: %v", err)
		return err
	}

	_, err := grpcClient.SendSupportMessage(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailServiceClientGrpc) verifyClientGrpc() {
	if grpcClient == nil {
		s.createClientGrpc()
	}
}

func (s *EmailServiceClientGrpc) createClientGrpc() {
	conn, err := grpc.Dial(s.config.EmailService.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("EmailServiceClientGrpc error connection: %v", err)
	}

	grpcClient = NewEmailServiceClient(conn)
}
