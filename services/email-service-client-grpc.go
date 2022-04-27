package services

import (
	"context"
	"log"
	"time"

	"github.com/oceano-dev/microservices-go-common/proto"

	"github.com/go-playground/validator/v10"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type EmailServiceClientGrpc struct {
}

type passwordCode struct {
	Email string
	Code  string
}

func (s *EmailServiceClientGrpc) SendPasswordCode(client proto.EmailServiceClient, email string, code string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, span := trace.NewSpan(ctx, "emailServiceGrpc.SendPasswordCodeReq")
	defer span.End()

	req := &proto.PasswordCodeReq{
		Email: email,
		Code:  code,
	}

	validator := validator.New()
	if err := validator.StructCtx(ctx, req); err != nil {
		trace.AddSpanError(span, err)
		log.Printf("emailServiceGrpc.SendPasswordCodeReq: %v", err)
		return err
	}

	_, err := client.SendPasswordCode(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailServiceClientGrpc) SendSupportMessage(client proto.EmailServiceClient, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, span := trace.NewSpan(ctx, "emailServiceGrpc.SendSupportMessageReq")
	defer span.End()

	req := &proto.SupportMessageReq{
		Message: message,
	}

	validator := validator.New()
	if err := validator.StructCtx(ctx, req); err != nil {
		trace.AddSpanError(span, err)
		log.Printf("emailServiceGrpc.SendSupportMessageReq: %v", err)
		return err
	}

	_, err := client.SendSupportMessage(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
