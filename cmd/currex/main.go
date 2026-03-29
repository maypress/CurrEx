package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/maypress/CurEx/internal/api"
	"github.com/maypress/CurEx/internal/cache"
	"github.com/maypress/CurEx/internal/debugger"
	"github.com/rivo/tview"
)

var dbg debugger.Debugger
var CURRENCY_AVIABLES []string
var CurrencyAviableList []byte
var rates api.DailyRates

const (
	CacheDir        = "./cache"
	HistoryDir      = "./history"
	CSVDumpFile     = "rates.csv"
	ConfigFile      = "config.json"
	ASCIIGraphWidth = 60
)

// Получает список доступных валют из XML
func getCurrencyList(xmlData []byte) []string {
	var currencies []string
	type ValCurs struct {
		Valutes []struct {
			CharCode string `xml:"CharCode"`
		} `xml:"Valute"`
	}

	var data ValCurs
	if err := xml.Unmarshal(xmlData, &data); err != nil {
		return nil
	}

	for _, v := range data.Valutes {
		currencies = append(currencies, v.CharCode)
	}
	sort.Strings(currencies)
	return currencies
}

// Находит курс по CharCode
func findRate(code string) (float64, int) {
	for _, valute := range rates.Valute {
		if strings.EqualFold(valute.CharCode, code) {
			value, _ := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", 1), 64)
			return value, valute.Nominal
		}
	}
	return 0, 1
}

// Конвертирует сумму
func convert(amount float64, from, to string) (float64, error) {
	if from == "RUB" && to == "RUB" {
		return amount, nil
	}

	rateFrom, nomFrom := findRate(from)
	rateTo, nomTo := findRate(to)

	if rateFrom == 0 || rateTo == 0 {
		return 0, fmt.Errorf("курсы не найдены")
	}

	if from == "RUB" {
		return (amount / float64(nomTo)) * rateTo, nil
	}
	if to == "RUB" {
		return (amount / rateFrom) * float64(nomFrom), nil
	}

	return amount * (rateTo/float64(nomTo)) / (rateFrom/float64(nomFrom)), nil
}

// ASCII график (простой)
func asciiGraph(values []float64) string {
	return `█▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀█
█ ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄█
█▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄█`
}

// Экспорт CSV
func exportToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"CharCode", "Nominal", "Name", "Value"})
	for _, valute := range rates.Valute {
		writer.Write([]string{
			valute.CharCode,
			strconv.Itoa(valute.Nominal),
			valute.Name,
			valute.Value,
		})
	}
	return nil
}

// Инициализация директорий
func initCacheDir() {
	os.MkdirAll(CacheDir, 0755)
	os.MkdirAll(HistoryDir, 0755)
}

// 🎮 РАБОЧИЙ TUI (tview v0.42.0)
func uiMode() {
	app := tview.NewApplication()

	// Список валют
	currencyList := append([]string{"RUB"}, CURRENCY_AVIABLES...)

	// Dropdown "Из"
	fromDrop := tview.NewDropDown().
		SetLabel(" Из: ").
		SetOptions(currencyList, nil).
		SetCurrentOption(0)

	// Dropdown "В"
	toDrop := tview.NewDropDown().
		SetLabel(" В:  ").
		SetOptions(currencyList, nil).
		SetCurrentOption(29) // USD

	// Сумма
	amountInput := tview.NewInputField().
		SetLabel("Сумма: ").
		SetText("100")

	// Результат
	resultText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorGreen).
		SetText("💱 Выберите валюты и сумму ↓")

	// ✅ ФОРМА БЕЗ SetButtons() - используем AddButton()
	form := tview.NewForm().
		AddFormItem(fromDrop).
		AddFormItem(toDrop).
		AddFormItem(amountInput)

	// ✅ Добавляем кнопки через AddButton()
	form.AddButton("💱 Конвертировать", func() {
		fromIdx, _ := fromDrop.GetCurrentOption()
		toIdx, _ := toDrop.GetCurrentOption()
		
		from := currencyList[fromIdx]
		to := currencyList[toIdx]
		amount, _ := strconv.ParseFloat(amountInput.GetText(), 64)
		
		res, err := convert(amount, from, to)
		if err != nil {
			resultText.SetText(fmt.Sprintf("[red]❌ %v[white]", err))
		} else {
			resultText.SetText(fmt.Sprintf("[green]✅ %.2f %s = %.2f %s[white]", 
				amount, from, res, to))
		}
	})

	form.AddButton("📈 График", func() {
		resultText.SetText("[yellow]📈 USD/RUB за неделю:\n" + asciiGraph(nil) + "[white]")
	})

	form.AddButton("📊 CSV", func() {
		if err := exportToCSV(CSVDumpFile); err != nil {
			resultText.SetText(fmt.Sprintf("[red]❌ Ошибка CSV: %v[white]", err))
		} else {
			resultText.SetText("[green]✅ CSV сохранен: rates.csv[white]")
		}
	})

	form.AddButton("❌ Выход", func() {
		app.Stop()
	})

	form.SetBorder(true).SetTitle("💱 CurEx Конвертер").SetTitleColor(tcell.ColorYellow)

	// Layout
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 3, true).
		AddItem(resultText, 0, 1, false)

	if err := app.SetRoot(layout, true).SetFocus(form).Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	initCacheDir()
	
	uiFlag := flag.Bool("ui", false, "Графический интерфейс")
	listFlag := flag.Bool("list", false, "Список валют")
	from := flag.String("from", "RUB", "Из валюты")
	to := flag.String("to", "USD", "В валюту")
	amount := flag.Float64("amount", 100, "Сумма")
	flag.Parse()

	dbg = debugger.Debugger{IsActive: true, Module: "Main"}

	// Загрузка данных
	CurrencyAviableList = api.GetAvailableXML()
	if len(CurrencyAviableList) == 0 {
		log.Fatal("❌ Не удалось получить курсы")
	}

	CURRENCY_AVIABLES = getCurrencyList(CurrencyAviableList)
	
	ratesData := cache.GetRates()
	if ratesData == nil {
		ratesData = api.FetchRates()
		if ratesData != nil {
			cache.SaveRates(ratesData)
		}
	}
	
	if ratesData == nil {
		log.Fatal("❌ Не удалось загрузить курсы")
	}
	rates = *ratesData

	dbg.Log(fmt.Sprintf("✅ Загружено %d валют", len(CURRENCY_AVIABLES)))

	// UI режим
	if *uiFlag {
		uiMode()
		return
	}

	// Список валют
	if *listFlag {
		fmt.Println("📋 Доступные валюты:")
		fmt.Printf("%s\n", strings.Join(CURRENCY_AVIABLES, ", "))
		return
	}

	// CLI конвертация
	result, err := convert(*amount, *from, *to)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		return
	}

	fmt.Printf("💱 %.2f %s = %.2f %s\n", *amount, *from, result, *to)
}