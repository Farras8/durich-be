package services

import (
	"context"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"time"
)

type DashboardService interface {
	GetStokDashboard(ctx context.Context, dateFrom, dateTo string) (*response.DashboardStokResponse, error)
	GetSalesDashboard(ctx context.Context, dateFrom, dateTo string) (*response.DashboardSalesResponse, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetStokDashboard(ctx context.Context, dateFrom, dateTo string) (*response.DashboardStokResponse, error) {
	from, to, err := s.parseDateRange(dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	return s.repo.GetStokDashboard(ctx, from, to)
}

func (s *dashboardService) GetSalesDashboard(ctx context.Context, dateFrom, dateTo string) (*response.DashboardSalesResponse, error) {
	from, to, err := s.parseDateRange(dateFrom, dateTo)
	if err != nil {
		return nil, err
	}

	return s.repo.GetSalesDashboard(ctx, from, to)
}

func (s *dashboardService) parseDateRange(dateFrom, dateTo string) (time.Time, time.Time, error) {
	var from, to time.Time
	var err error

	if dateFrom == "" {
		from = time.Now().AddDate(0, 0, -30)
	} else {
		from, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return from, to, err
		}
	}

	if dateTo == "" {
		to = time.Now()
	} else {
		to, err = time.Parse("2006-01-02", dateTo)
		if err != nil {
			return from, to, err
		}
		to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	return from, to, nil
}
