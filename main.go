package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

// FinancialData represents monthly financial data
type FinancialData struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	NetFlow float64 `json:"net_flow"`
}

// CompanyProfile represents the company's basic info
type CompanyProfile struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Sector            string  `json:"sector"`
	MonthlyAvgIncome  float64 `json:"monthly_avg_income"`
	MonthlyAvgExpense float64 `json:"monthly_avg_expense"`
}

// FinancialAnalysis represents the complete financial analysis
type FinancialAnalysis struct {
	Company        CompanyProfile  `json:"company"`
	HistoricalData []FinancialData `json:"historical_data"`
	Predictions    []FinancialData `json:"predictions"`
	Summary        AnalysisSummary `json:"summary"`
	CreatedAt      time.Time       `json:"created_at"`
}

// AnalysisSummary provides key insights
type AnalysisSummary struct {
	TotalHistoricalIncome  float64  `json:"total_historical_income"`
	TotalHistoricalExpense float64  `json:"total_historical_expense"`
	TotalHistoricalNetFlow float64  `json:"total_historical_net_flow"`
	PredictedTotalIncome   float64  `json:"predicted_total_income"`
	PredictedTotalExpense  float64  `json:"predicted_total_expense"`
	PredictedTotalNetFlow  float64  `json:"predicted_total_net_flow"`
	GrowthTrend            string   `json:"growth_trend"`
	RiskLevel              string   `json:"risk_level"`
	CashFlowHealth         string   `json:"cash_flow_health"`
	Recommendations        []string `json:"recommendations"`
}

// AnalysisRequest represents the input data structure
type AnalysisRequest struct {
	Company        CompanyProfile  `json:"company"`
	HistoricalData []FinancialData `json:"historical_data"`
}

// FinancialAnalyzer handles the prediction logic
type FinancialAnalyzer struct {
	// In a real application, this could connect to a database
}

// PredictNext6Months generates predictions based on historical data
func (fa *FinancialAnalyzer) PredictNext6Months(historical []FinancialData) []FinancialData {
	predictions := make([]FinancialData, 6)

	if len(historical) == 0 {
		return predictions
	}

	// Calculate trends and seasonal patterns
	incomeGrowthRate := fa.calculateGrowthRate(historical, "income")
	expenseGrowthRate := fa.calculateGrowthRate(historical, "expense")

	// Get the last known values as baseline
	lastData := historical[len(historical)-1]
	baseIncome := lastData.Income
	baseExpense := lastData.Expense

	// Add seasonal adjustment
	seasonalFactors := fa.getSeasonalFactors(historical)

	for i := 0; i < 6; i++ {
		monthIndex := (len(historical) + i) % 12
		seasonalFactor := seasonalFactors[monthIndex]

		// Apply growth rate and seasonal adjustment
		predictedIncome := baseIncome * math.Pow(1+incomeGrowthRate, float64(i+1)) * seasonalFactor
		predictedExpense := baseExpense * math.Pow(1+expenseGrowthRate, float64(i+1))

		// Add some volatility (random factor between 0.9-1.1)
		volatilityFactor := 0.95 + (float64(i%3) * 0.05) // Simplified volatility
		predictedIncome *= volatilityFactor
		predictedExpense *= (2.0 - volatilityFactor) // Inverse for expenses

		predictions[i] = FinancialData{
			Month:   fa.getMonthName(time.Now().AddDate(0, i+1, 0)),
			Income:  math.Round(predictedIncome*100) / 100,
			Expense: math.Round(predictedExpense*100) / 100,
			NetFlow: math.Round((predictedIncome-predictedExpense)*100) / 100,
		}
	}

	return predictions
}

// calculateGrowthRate calculates monthly growth rate
func (fa *FinancialAnalyzer) calculateGrowthRate(data []FinancialData, field string) float64 {
	if len(data) < 2 {
		return 0.02 // Default 2% growth
	}

	var totalGrowth float64
	var validPeriods int

	for i := 1; i < len(data); i++ {
		var current, previous float64

		if field == "income" {
			current = data[i].Income
			previous = data[i-1].Income
		} else {
			current = data[i].Expense
			previous = data[i-1].Expense
		}

		if previous > 0 {
			growth := (current - previous) / previous
			totalGrowth += growth
			validPeriods++
		}
	}

	if validPeriods == 0 {
		return 0.02
	}

	avgGrowthRate := totalGrowth / float64(validPeriods)

	// Cap growth rate between -20% and +30% monthly
	if avgGrowthRate > 0.30 {
		avgGrowthRate = 0.30
	} else if avgGrowthRate < -0.20 {
		avgGrowthRate = -0.20
	}

	return avgGrowthRate
}

// getSeasonalFactors returns seasonal adjustment factors
func (fa *FinancialAnalyzer) getSeasonalFactors(data []FinancialData) []float64 {
	// Default seasonal factors (can be calculated from historical data)
	factors := []float64{
		1.0, 0.95, 1.05, 1.1, 1.15, 1.2, // Jan-Jun
		1.25, 1.2, 1.1, 1.05, 1.0, 1.3, // Jul-Dec (Dec higher for year-end)
	}

	// In a more sophisticated version, calculate actual seasonal patterns
	if len(data) >= 12 {
		// Calculate seasonal patterns from historical data
		monthlyAvgs := make([]float64, 12)
		monthlyCounts := make([]int, 12)

		for _, d := range data {
			month := fa.getMonthIndex(d.Month)
			if month >= 0 && month < 12 {
				monthlyAvgs[month] += d.Income
				monthlyCounts[month]++
			}
		}

		totalAvg := 0.0
		validMonths := 0

		for i := 0; i < 12; i++ {
			if monthlyCounts[i] > 0 {
				monthlyAvgs[i] /= float64(monthlyCounts[i])
				totalAvg += monthlyAvgs[i]
				validMonths++
			}
		}

		if validMonths > 0 {
			totalAvg /= float64(validMonths)

			for i := 0; i < 12; i++ {
				if monthlyCounts[i] > 0 {
					factors[i] = monthlyAvgs[i] / totalAvg
				}
			}
		}
	}

	return factors
}

// getMonthIndex returns month index (0-11)
func (fa *FinancialAnalyzer) getMonthIndex(monthName string) int {
	months := []string{
		"Ocak", "Şubat", "Mart", "Nisan", "Mayıs", "Haziran",
		"Temmuz", "Ağustos", "Eylül", "Ekim", "Kasım", "Aralık",
	}

	for i, month := range months {
		if month == monthName {
			return i
		}
	}
	return -1
}

// getMonthName returns Turkish month name
func (fa *FinancialAnalyzer) getMonthName(t time.Time) string {
	months := []string{
		"Ocak", "Şubat", "Mart", "Nisan", "Mayıs", "Haziran",
		"Temmuz", "Ağustos", "Eylül", "Ekim", "Kasım", "Aralık",
	}
	return months[t.Month()-1]
}

// GenerateAnalysis creates a complete financial analysis
func (fa *FinancialAnalyzer) GenerateAnalysis(req AnalysisRequest) *FinancialAnalysis {
	predictions := fa.PredictNext6Months(req.HistoricalData)
	summary := fa.generateSummary(req.HistoricalData, predictions)

	return &FinancialAnalysis{
		Company:        req.Company,
		HistoricalData: req.HistoricalData,
		Predictions:    predictions,
		Summary:        summary,
		CreatedAt:      time.Now(),
	}
}

// generateSummary creates analysis summary
func (fa *FinancialAnalyzer) generateSummary(historical, predicted []FinancialData) AnalysisSummary {
	var histIncome, histExpense, histNetFlow float64
	var predIncome, predExpense, predNetFlow float64

	// Calculate totals
	for _, h := range historical {
		histIncome += h.Income
		histExpense += h.Expense
		histNetFlow += h.NetFlow
	}

	for _, p := range predicted {
		predIncome += p.Income
		predExpense += p.Expense
		predNetFlow += p.NetFlow
	}

	// Determine trends and health
	growthTrend := "Stabil"
	if predIncome > histIncome*1.1 {
		growthTrend = "Yükseliş"
	} else if predIncome < histIncome*0.9 {
		growthTrend = "Düşüş"
	}

	riskLevel := "Orta"
	if predNetFlow < 0 {
		riskLevel = "Yüksek"
	} else if predNetFlow > histNetFlow*1.2 {
		riskLevel = "Düşük"
	}

	cashFlowHealth := "Normal"
	avgNetFlow := predNetFlow / 6
	if avgNetFlow < 0 {
		cashFlowHealth = "Risk"
	} else if avgNetFlow > histNetFlow/float64(len(historical))*1.5 {
		cashFlowHealth = "Güçlü"
	}

	// Generate recommendations
	recommendations := fa.generateRecommendations(growthTrend, riskLevel, cashFlowHealth, predNetFlow)

	return AnalysisSummary{
		TotalHistoricalIncome:  math.Round(histIncome*100) / 100,
		TotalHistoricalExpense: math.Round(histExpense*100) / 100,
		TotalHistoricalNetFlow: math.Round(histNetFlow*100) / 100,
		PredictedTotalIncome:   math.Round(predIncome*100) / 100,
		PredictedTotalExpense:  math.Round(predExpense*100) / 100,
		PredictedTotalNetFlow:  math.Round(predNetFlow*100) / 100,
		GrowthTrend:            growthTrend,
		RiskLevel:              riskLevel,
		CashFlowHealth:         cashFlowHealth,
		Recommendations:        recommendations,
	}
}

// generateRecommendations creates actionable recommendations
func (fa *FinancialAnalyzer) generateRecommendations(growth, risk, health string, netFlow float64) []string {
	var recommendations []string

	if risk == "Yüksek" {
		recommendations = append(recommendations, "Acil nakit akış planı oluşturun")
		recommendations = append(recommendations, "Gereksiz giderleri kısmayı düşünün")
		recommendations = append(recommendations, "Alternatif finansman kaynaklarını araştırın")
	}

	if growth == "Düşüş" {
		recommendations = append(recommendations, "Yeni pazarlama stratejileri geliştirin")
		recommendations = append(recommendations, "Maliyet optimizasyonu yapın")
		recommendations = append(recommendations, "Ürün/hizmet portföyünüzü gözden geçirin")
	}

	if health == "Güçlü" {
		recommendations = append(recommendations, "Yatırım fırsatlarını değerlendirin")
		recommendations = append(recommendations, "Büyüme stratejileri planlayın")
		recommendations = append(recommendations, "Acil durum fonu oluşturun")
	}

	if netFlow > 0 {
		recommendations = append(recommendations, "Kâr paylaşım planı düşünün")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Mevcut performansınızı korumaya odaklanın")
	}

	return recommendations
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// HTTP Handlers
func (fa *FinancialAnalyzer) analyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Debug log
	fmt.Printf("Method: %s, URL: %s\n", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed. Use POST", http.StatusMethodNotAllowed)
		return
	}

	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate input
	if len(req.HistoricalData) == 0 {
		http.Error(w, "Historical data is required", http.StatusBadRequest)
		return
	}

	// Calculate net flows if not provided
	for i := range req.HistoricalData {
		req.HistoricalData[i].NetFlow = req.HistoricalData[i].Income - req.HistoricalData[i].Expense
	}

	analysis := fa.GenerateAnalysis(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(analysis); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// healthHandler provides health check endpoint
func (fa *FinancialAnalyzer) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Health check - Method: %s\n", r.Method)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"time":    time.Now().Format(time.RFC3339),
		"service": "KOBİ Financial Analysis API",
	})
}

// Simple home handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Home - Method: %s, Path: %s\n", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"service": "KOBİ Mali Durum Tahmin Sistemi",
		"version": "1.0.0",
		"endpoints": map[string]string{
			"analyze": "POST /api/analyze",
			"health":  "GET /api/health",
		},
		"status": "running",
		"time":   time.Now().Format("2006-01-02 15:04:05"),
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	analyzer := &FinancialAnalyzer{}

	// Setup routes without external router
	http.HandleFunc("/", corsMiddleware(homeHandler))
	http.HandleFunc("/api/analyze", corsMiddleware(analyzer.analyzeHandler))
	http.HandleFunc("/api/health", corsMiddleware(analyzer.healthHandler))

	fmt.Println("🚀 KOBİ Mali Durum Tahmin Sistemi başlatılıyor...")
	fmt.Println("🌐 Server: http://localhost:8080")
	fmt.Println("📊 API Endpoint: http://localhost:8080/api/analyze")
	fmt.Println("🔍 Health Check: http://localhost:8080/api/health")
	fmt.Println("📋 Home: http://localhost:8080/")
	fmt.Println("\n✅ Sistem hazır - test client'ını çalıştırabilirsiniz")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
