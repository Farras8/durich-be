package services

import (
	"context"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
)

type TraceabilityService interface {
	TraceLot(ctx context.Context, lotID string) (*response.TraceLotResponse, error)
	TraceFruit(ctx context.Context, fruitID string) (*response.TraceFruitResponse, error)
	TraceShipment(ctx context.Context, shipmentID string) (*response.TraceShipmentResponse, error)
}

type traceabilityService struct {
	repo repository.TraceabilityRepository
}

func NewTraceabilityService(repo repository.TraceabilityRepository) TraceabilityService {
	return &traceabilityService{repo: repo}
}

func (s *traceabilityService) TraceLot(ctx context.Context, lotID string) (*response.TraceLotResponse, error) {
	return s.repo.TraceLot(ctx, lotID)
}

func (s *traceabilityService) TraceFruit(ctx context.Context, fruitID string) (*response.TraceFruitResponse, error) {
	return s.repo.TraceFruit(ctx, fruitID)
}

func (s *traceabilityService) TraceShipment(ctx context.Context, shipmentID string) (*response.TraceShipmentResponse, error) {
	return s.repo.TraceShipment(ctx, shipmentID)
}
