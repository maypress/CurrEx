// internal/api/api.go
package api

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

const API_URL = "https://www.cbr-xml-daily.ru/daily_utf8.xml"

type Valute struct {
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type DailyRates struct {
	Date   time.Time `xml:"Date"`
	Valute []Valute  `xml:"Valute"`
}

// Получает XML курсы (с User-Agent)
func GetAvailableXML() []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", API_URL, nil)
	req.Header.Set("User-Agent", "CurEx/1.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

// Парсит и возвращает структуру курсов
func FetchRates() *DailyRates {
	data := GetAvailableXML()
	if data == nil {
		return nil
	}

	var rates DailyRates
	if err := xml.Unmarshal(data, &rates); err != nil {
		return nil
	}

	// Парсим дату
	rates.Date, _ = time.Parse("02/01/2006", rates.Date.Format("02/01/2006"))
	return &rates
}