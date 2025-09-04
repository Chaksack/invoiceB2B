package services

import (
	"encoding/json"
	"fmt"
	"time"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"gorm.io/gorm"
)

// ReportingService interface defines reporting operations
type ReportingService interface {
	GenerateReport(request *dtos.ReportRequest) (interface{}, error)
	GenerateLoanApplicationReport(request *dtos.ReportRequest) (*dtos.LoanApplicationReportResponse, error)
	GenerateInvoiceReport(request *dtos.ReportRequest) (*dtos.InvoiceReportResponse, error)
	GenerateUserReport(request *dtos.ReportRequest) (*dtos.UserReportResponse, error)
	GenerateActivityReport(request *dtos.ReportRequest) (*dtos.ActivityReportResponse, error)
	GenerateFinancialInstitutionReport(request *dtos.ReportRequest) (*dtos.FinancialInstitutionReportResponse, error)
	GenerateOverviewReport(request *dtos.ReportRequest) (*dtos.OverviewReportResponse, error)
	ExportReport(request *dtos.ReportRequest, data interface{}) (*dtos.ExportReportResponse, error)
}

type reportingService struct {
	db *gorm.DB
}

// NewReportingService creates a new reporting service
func NewReportingService(db *gorm.DB) ReportingService {
	return &reportingService{
		db: db,
	}
}

// GenerateReport is a dispatcher method that generates the appropriate report type
func (s *reportingService) GenerateReport(request *dtos.ReportRequest) (interface{}, error) {
	switch request.ReportType {
	case "loan_applications":
		return s.GenerateLoanApplicationReport(request)
	case "invoices":
		return s.GenerateInvoiceReport(request)
	case "users":
		return s.GenerateUserReport(request)
	case "activity":
		return s.GenerateActivityReport(request)
	case "financial_institutions":
		return s.GenerateFinancialInstitutionReport(request)
	case "overview":
		return s.GenerateOverviewReport(request)
	default:
		return nil, fmt.Errorf("unsupported report type: %s", request.ReportType)
	}
}

// GenerateLoanApplicationReport generates a comprehensive loan application report
func (s *reportingService) GenerateLoanApplicationReport(request *dtos.ReportRequest) (*dtos.LoanApplicationReportResponse, error) {
	var summary dtos.LoanApplicationSummary
	var trends []dtos.LoanApplicationTrendData
	var statusBreakdown []dtos.StatusBreakdownData
	var amountRanges []dtos.AmountRangeData
	var sourceBreakdown []dtos.SourceBreakdownData
	var topApplications []dtos.LoanApplicationSummaryItem

	// Build base query
	query := s.db.Model(&models.LoanApplication{})
	if request.StartDate != nil {
		query = query.Where("created_at >= ?", request.StartDate)
	}
	if request.EndDate != nil {
		query = query.Where("created_at <= ?", request.EndDate)
	}
	if request.Status != "" {
		query = query.Where("status = ?", request.Status)
	}
	if request.UserID != nil {
		query = query.Where("user_id = ?", *request.UserID)
	}

	// Calculate summary statistics
	var totalCount int64
	var totalRequested, totalApproved, totalDisbursed float64
	
	query.Count(&totalCount)
	
	query.Select("COALESCE(SUM(requested_amount), 0) as total_requested, COALESCE(SUM(approved_amount), 0) as total_approved, COALESCE(SUM(disbursed_amount), 0) as total_disbursed").
		Row().Scan(&totalRequested, &totalApproved, &totalDisbursed)

	// Calculate rates
	approvalRate := float64(0)
	disbursementRate := float64(0)
	if totalCount > 0 {
		var approvedCount, disbursedCount int64
		query.Where("status = ?", "approved").Count(&approvedCount)
		query.Where("status = ?", "disbursed").Count(&disbursedCount)
		approvalRate = float64(approvedCount) / float64(totalCount) * 100
		disbursementRate = float64(disbursedCount) / float64(totalCount) * 100
	}

	// Calculate average processing days
	var avgProcessingDays float64
	s.db.Model(&models.LoanApplication{}).
		Where("approved_at IS NOT NULL AND submitted_at IS NOT NULL").
		Select("AVG(EXTRACT(epoch FROM (approved_at - submitted_at))/86400)").
		Row().Scan(&avgProcessingDays)

	summary = dtos.LoanApplicationSummary{
		TotalApplications:     totalCount,
		TotalAmountRequested:  totalRequested,
		TotalAmountApproved:   totalApproved,
		TotalAmountDisbursed:  totalDisbursed,
		AverageProcessingDays: avgProcessingDays,
		ApprovalRate:          approvalRate,
		DisbursementRate:      disbursementRate,
	}

	// Generate status breakdown
	var statusCounts []struct {
		Status string
		Count  int64
		Amount float64
	}
	query.Select("status, count(*) as count, COALESCE(SUM(requested_amount), 0) as amount").
		Group("status").Find(&statusCounts)

	for _, sc := range statusCounts {
		percentage := float64(sc.Count) / float64(totalCount) * 100
		statusBreakdown = append(statusBreakdown, dtos.StatusBreakdownData{
			Status:     sc.Status,
			Count:      sc.Count,
			Amount:     sc.Amount,
			Percentage: percentage,
		})
	}

	// Generate source breakdown
	var sourceCounts []struct {
		Source string
		Count  int64
		Amount float64
	}
	query.Select("source, count(*) as count, COALESCE(SUM(requested_amount), 0) as amount").
		Group("source").Find(&sourceCounts)

	for _, sc := range sourceCounts {
		percentage := float64(sc.Count) / float64(totalCount) * 100
		sourceBreakdown = append(sourceBreakdown, dtos.SourceBreakdownData{
			Source:     sc.Source,
			Count:      sc.Count,
			Amount:     sc.Amount,
			Percentage: percentage,
		})
	}

	// Generate amount ranges
	amountRanges = []dtos.AmountRangeData{
		{Range: "0-1000", Count: 0, Percentage: 0},
		{Range: "1001-5000", Count: 0, Percentage: 0},
		{Range: "5001-10000", Count: 0, Percentage: 0},
		{Range: "10001-50000", Count: 0, Percentage: 0},
		{Range: "50000+", Count: 0, Percentage: 0},
	}

	for i := range amountRanges {
		var count int64
		switch i {
		case 0:
			query.Where("requested_amount <= 1000").Count(&count)
		case 1:
			query.Where("requested_amount > 1000 AND requested_amount <= 5000").Count(&count)
		case 2:
			query.Where("requested_amount > 5000 AND requested_amount <= 10000").Count(&count)
		case 3:
			query.Where("requested_amount > 10000 AND requested_amount <= 50000").Count(&count)
		case 4:
			query.Where("requested_amount > 50000").Count(&count)
		}
		amountRanges[i].Count = count
		if totalCount > 0 {
			amountRanges[i].Percentage = float64(count) / float64(totalCount) * 100
		}
	}

	// Get top applications
	var applications []models.LoanApplication
	query.Preload("User").Order("requested_amount DESC").Limit(10).Find(&applications)

	for _, app := range applications {
		topApplications = append(topApplications, dtos.LoanApplicationSummaryItem{
			ID:                   app.ID,
			ApplicationReference: app.ApplicationReference,
			CompanyName:         app.User.CompanyName,
			RequestedAmount:     app.RequestedAmount,
			Status:              string(app.Status),
			SubmittedAt:         *app.SubmittedAt,
		})
	}

	return &dtos.LoanApplicationReportResponse{
		Summary:         summary,
		Trends:          trends,
		StatusBreakdown: statusBreakdown,
		AmountRanges:    amountRanges,
		SourceBreakdown: sourceBreakdown,
		TopApplications: topApplications,
	}, nil
}

// GenerateInvoiceReport generates a comprehensive invoice report
func (s *reportingService) GenerateInvoiceReport(request *dtos.ReportRequest) (*dtos.InvoiceReportResponse, error) {
	// Implementation for invoice reports
	return &dtos.InvoiceReportResponse{
		Summary: dtos.InvoiceSummary{
			TotalInvoices:        0,
			TotalInvoiceAmount:   0,
			TotalProcessed:       0,
			AverageInvoiceAmount: 0,
			ProcessingRate:       0,
		},
		Trends:          []dtos.InvoiceTrendData{},
		StatusBreakdown: []dtos.StatusBreakdownData{},
		AmountRanges:    []dtos.AmountRangeData{},
		TopInvoices:     []dtos.InvoiceSummaryItem{},
	}, nil
}

// GenerateUserReport generates a comprehensive user report
func (s *reportingService) GenerateUserReport(request *dtos.ReportRequest) (*dtos.UserReportResponse, error) {
	// Implementation for user reports
	return &dtos.UserReportResponse{
		Summary: dtos.UserSummary{
			TotalUsers:          0,
			ActiveUsers:         0,
			VerifiedUsers:       0,
			KYCCompletedUsers:   0,
			AverageApplications: 0,
		},
		Trends:       []dtos.UserTrendData{},
		KYCBreakdown: []dtos.KYCBreakdownData{},
		TopUsers:     []dtos.UserSummaryItem{},
	}, nil
}

// GenerateActivityReport generates a comprehensive activity report
func (s *reportingService) GenerateActivityReport(request *dtos.ReportRequest) (*dtos.ActivityReportResponse, error) {
	// Implementation for activity reports
	return &dtos.ActivityReportResponse{
		Summary: dtos.ActivitySummary{
			TotalActivities: 0,
			UniqueUsers:     0,
			UniqueStaff:     0,
			TopAction:       "",
		},
		Trends:          []dtos.ActivityTrendData{},
		ActionBreakdown: []dtos.ActionBreakdownData{},
		TopActivities:   []dtos.ActivitySummaryItem{},
	}, nil
}

// GenerateFinancialInstitutionReport generates a comprehensive financial institution report
func (s *reportingService) GenerateFinancialInstitutionReport(request *dtos.ReportRequest) (*dtos.FinancialInstitutionReportResponse, error) {
	// Implementation for financial institution reports
	return &dtos.FinancialInstitutionReportResponse{
		Summary: dtos.FinancialInstitutionReportSummary{
			TotalInstitutions:     0,
			ActiveInstitutions:    0,
			TotalProducts:         0,
			TotalApplicationsSent: 0,
			AverageInterestRate:   0,
		},
		InstitutionBreakdown: []dtos.InstitutionBreakdownData{},
		ProductBreakdown:     []dtos.ProductBreakdownData{},
		TopInstitutions:      []dtos.FinancialInstitutionSummaryItem{},
	}, nil
}

// GenerateOverviewReport generates a comprehensive overview report
func (s *reportingService) GenerateOverviewReport(request *dtos.ReportRequest) (*dtos.OverviewReportResponse, error) {
	// Get loan application summary
	loanReport, err := s.GenerateLoanApplicationReport(request)
	if err != nil {
		return nil, fmt.Errorf("failed to generate loan application summary: %w", err)
	}

	// Get recent activities
	var recentActivities []models.ActivityLog
	s.db.Preload("User").Order("timestamp DESC").Limit(10).Find(&recentActivities)

	var activityItems []dtos.ActivitySummaryItem
	for _, activity := range recentActivities {
		userEmail := ""
		if activity.User != nil {
			userEmail = activity.User.Email
		}
		activityItems = append(activityItems, dtos.ActivitySummaryItem{
			ID:        activity.ID,
			Action:    activity.Action,
			UserEmail: userEmail,
			Timestamp: activity.CreatedAt,
			Details:   activity.Details,
		})
	}

	return &dtos.OverviewReportResponse{
		LoanApplications: loanReport.Summary,
		Invoices: dtos.InvoiceSummary{
			TotalInvoices:        0,
			TotalInvoiceAmount:   0,
			TotalProcessed:       0,
			AverageInvoiceAmount: 0,
			ProcessingRate:       0,
		},
		Users: dtos.UserSummary{
			TotalUsers:          0,
			ActiveUsers:         0,
			VerifiedUsers:       0,
			KYCCompletedUsers:   0,
			AverageApplications: 0,
		},
		Activity: dtos.ActivitySummary{
			TotalActivities: 0,
			UniqueUsers:     0,
			UniqueStaff:     0,
			TopAction:       "",
		},
		FinancialInstitutions: dtos.FinancialInstitutionReportSummary{
			TotalInstitutions:     0,
			ActiveInstitutions:    0,
			TotalProducts:         0,
			TotalApplicationsSent: 0,
			AverageInterestRate:   0,
		},
		SystemHealth: dtos.SystemHealthData{
			TotalTransactions:   0,
			SuccessRate:         0,
			AverageResponseTime: 0,
			ExpiredDataCount:    0,
			PendingRefreshCount: 0,
		},
		RecentActivities: activityItems,
	}, nil
}

// ExportReport exports report data in the requested format
func (s *reportingService) ExportReport(request *dtos.ReportRequest, data interface{}) (*dtos.ExportReportResponse, error) {
	if request.ExportFormat == "" {
		request.ExportFormat = "json"
	}

	fileName := fmt.Sprintf("report_%s_%d.%s", request.ReportType, time.Now().Unix(), request.ExportFormat)
	
	switch request.ExportFormat {
	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return &dtos.ExportReportResponse{
			FileName:    fileName,
			Format:      "json",
			Size:        int64(len(jsonData)),
			GeneratedAt: time.Now(),
			DownloadURL: fmt.Sprintf("/api/v1/admin/reports/download/%s", fileName),
		}, nil
	case "csv":
		// CSV export implementation would go here
		return &dtos.ExportReportResponse{
			FileName:    fileName,
			Format:      "csv",
			Size:        0,
			GeneratedAt: time.Now(),
			DownloadURL: fmt.Sprintf("/api/v1/admin/reports/download/%s", fileName),
		}, nil
	case "pdf":
		// PDF export implementation would go here
		return &dtos.ExportReportResponse{
			FileName:    fileName,
			Format:      "pdf",
			Size:        0,
			GeneratedAt: time.Now(),
			DownloadURL: fmt.Sprintf("/api/v1/admin/reports/download/%s", fileName),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.ExportFormat)
	}
}