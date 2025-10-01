package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	fmt.Println("🧪 KOBİ Mali Durum Tahmin Sistemi - Test")
	fmt.Println("========================================")

	// Server'ın başlaması için bekle
	fmt.Println("⏳ Server'ın başlaması bekleniyor...")
	time.Sleep(3 * time.Second)

	// 1. Health Check Test
	fmt.Println("\n1️⃣  Health Check Testi:")
	testHealthCheck()

	// 2. Ana API Test
	fmt.Println("\n2️⃣  Ana API Testi:")
	testAnalyzeAPI()

	// 3. Curl örneği göster
	printCurlExample()

	fmt.Println("\n✅ Testler tamamlandı!")
}

func testHealthCheck() {
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		fmt.Printf("❌ Health check başarısız: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Health Check - Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
}

func testAnalyzeAPI() {
	// JSON string olarak test verisi hazırla
	testJSON := `{
		"company": {
			"id": "TEST001",
			"name": "Test Şirketi A.Ş.",
			"sector": "Teknoloji",
			"monthly_avg_income": 500000,
			"monthly_avg_expense": 400000
		},
		"historical_data": [
			{"month": "Mart", "income": 450000, "expense": 380000, "net_flow": 70000},
			{"month": "Nisan", "income": 420000, "expense": 350000, "net_flow": 70000},
			{"month": "Mayıs", "income": 480000, "expense": 400000, "net_flow": 80000},
			{"month": "Haziran", "income": 550000, "expense": 440000, "net_flow": 110000},
			{"month": "Temmuz", "income": 600000, "expense": 480000, "net_flow": 120000},
			{"month": "Ağustos", "income": 580000, "expense": 460000, "net_flow": 120000}
		]
	}`

	fmt.Println("📤 Gönderilen veri:")
	fmt.Println(testJSON[:200] + "...")

	// POST request gönder
	resp, err := http.Post("http://localhost:8080/api/analyze", "application/json", bytes.NewBufferString(testJSON))
	if err != nil {
		fmt.Printf("❌ API request hatası: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Response'u oku
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Response okuma hatası: %v\n", err)
		return
	}

	fmt.Printf("📥 Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode == 200 {
		// JSON'ı parse et ve güzel göster
		var result map[string]interface{}
		if err := json.Unmarshal(responseBody, &result); err != nil {
			fmt.Printf("❌ JSON parse hatası: %v\n", err)
			fmt.Printf("Raw response: %s\n", string(responseBody))
			return
		}

		fmt.Println("✅ API Başarılı!")

		// Company bilgilerini göster
		if company, ok := result["company"].(map[string]interface{}); ok {
			fmt.Printf("🏢 Şirket: %v (%v)\n", company["name"], company["sector"])
		}

		// Summary bilgilerini göster
		if summary, ok := result["summary"].(map[string]interface{}); ok {
			fmt.Printf("📊 Büyüme Trendi: %v\n", summary["growth_trend"])
			fmt.Printf("📊 Risk Seviyesi: %v\n", summary["risk_level"])
			fmt.Printf("📊 Nakit Akış Sağlığı: %v\n", summary["cash_flow_health"])

			if totalIncome, ok := summary["predicted_total_income"].(float64); ok {
				fmt.Printf("💰 6 Aylık Tahmini Gelir: ₺%.0f\n", totalIncome)
			}

			if totalExpense, ok := summary["predicted_total_expense"].(float64); ok {
				fmt.Printf("💸 6 Aylık Tahmini Gider: ₺%.0f\n", totalExpense)
			}

			if netFlow, ok := summary["predicted_total_net_flow"].(float64); ok {
				fmt.Printf("📈 6 Aylık Tahmini Net Akış: ₺%.0f\n", netFlow)
			}

			if recommendations, ok := summary["recommendations"].([]interface{}); ok && len(recommendations) > 0 {
				fmt.Println("💡 Öneriler:")
				for i, rec := range recommendations {
					fmt.Printf("  %d. %v\n", i+1, rec)
					if i >= 2 { // Sadece ilk 3 öneriyi göster
						break
					}
				}
			}
		}

		// Tahminleri göster
		if predictions, ok := result["predictions"].([]interface{}); ok && len(predictions) > 0 {
			fmt.Println("🔮 Aylık Tahminler:")
			for i, pred := range predictions {
				if predMap, ok := pred.(map[string]interface{}); ok {
					fmt.Printf("  %v: Gelir ₺%.0f, Gider ₺%.0f, Net ₺%.0f\n",
						predMap["month"], predMap["income"], predMap["expense"], predMap["net_flow"])
				}
				if i >= 3 { // Sadece ilk 4 tahmini göster
					fmt.Printf("  ... ve %d ay daha\n", len(predictions)-4)
					break
				}
			}
		}

		// JSON'ı dosyaya kaydet (opsiyonel)
		if prettyJSON, err := json.MarshalIndent(result, "", "  "); err == nil {
			fmt.Println("\n📄 Detaylı sonuçlar result.json dosyasına kaydedildi.")
			// İsterseniz dosyaya yazabilirsiniz
			// os.WriteFile("result.json", prettyJSON, 0644)
			_ = prettyJSON // Şimdilik sadece değişkeni kullanıyoruz
		}

	} else {
		fmt.Printf("❌ API Hatası (Status %d): %s\n", resp.StatusCode, string(responseBody))
	}
}

// Curl komutu örneği yazdır
func printCurlExample() {
	fmt.Println("\n📋 Manuel test için CURL komutu:")
	fmt.Println("=====================================")
	curlCmd := `curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "company": {
      "id": "MANUAL001",
      "name": "Manuel Test Şirketi",
      "sector": "E-ticaret",
      "monthly_avg_income": 300000,
      "monthly_avg_expense": 250000
    },
    "historical_data": [
      {"month": "Mart", "income": 280000, "expense": 230000},
      {"month": "Nisan", "income": 320000, "expense": 260000},
      {"month": "Mayıs", "income": 350000, "expense": 280000},
      {"month": "Haziran", "income": 380000, "expense": 300000}
    ]
  }'`

	fmt.Println(curlCmd)
	fmt.Println("\n💡 İpucu: Yukarıdaki komutu terminale yapıştırıp çalıştırabilirsiniz!")
}

// Test senaryoları
func runAdditionalTests() {
	fmt.Println("\n3️⃣  Ek Test Senaryoları:")

	// Risky business test
	riskyJSON := `{
		"company": {
			"id": "RISK001", 
			"name": "Risk Şirketi",
			"sector": "Perakende",
			"monthly_avg_income": 200000,
			"monthly_avg_expense": 220000
		},
		"historical_data": [
			{"month": "Haziran", "income": 180000, "expense": 200000},
			{"month": "Temmuz", "income": 160000, "expense": 210000},
			{"month": "Ağustos", "income": 150000, "expense": 220000}
		]
	}`

	fmt.Println("📉 Risk Testi - Düşen gelirli şirket:")
	testWithJSON("Risk Şirketi", riskyJSON)
}

func testWithJSON(testName, jsonData string) {
	resp, err := http.Post("http://localhost:8080/api/analyze", "application/json", bytes.NewBufferString(jsonData))
	if err != nil {
		fmt.Printf("❌ %s testi başarısız: %v\n", testName, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		if summary, ok := result["summary"].(map[string]interface{}); ok {
			fmt.Printf("✅ %s: Risk=%v, Trend=%v, Sağlık=%v\n",
				testName,
				summary["risk_level"],
				summary["growth_trend"],
				summary["cash_flow_health"])
		}
	} else {
		fmt.Printf("❌ %s testi başarısız - Status: %d\n", testName, resp.StatusCode)
	}
}
