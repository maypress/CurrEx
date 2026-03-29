# CurEx 💱 Конвертер валют ЦБ РФ

[![Go](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**Быстрый CLI + TUI конвертер валют по официальным курсам Центрального банка РФ.**

Получает данные с [ЦБ РФ](https://www.cbr.ru/scripts/XML_daily.asp) через зеркало `cbr-xml-daily.ru`. 50+ валют, кэширование, CSV-экспорт.

## ✨ Быстрый старт

```bash
# Клонировать и запустить
git clone https://github.com/maypress/CurEx.git
cd CurEx
go mod tidy
```

## 🚀 Использование

### ✅ **CLI с флагами (основной синтаксис)**
```bash
# Конвертация USD → EUR
PS C:\CurrEx> go run cmd/currex/main.go --amount 100 --from USD --to EUR
💱 100.00 USD = 92.50 EUR

# Другие примеры
go run cmd/currex/main.go --amount 50 --from RUB --to USD
go run cmd/currex/main.go --amount 1000 --from EUR --to RUB
```

### 📱 **Графический интерфейс**
```bash
go run cmd/currex/main.go --ui
```

### 📋 **Полезные команды**
```bash
go run cmd/currex/main.go --list      # Список всех валют
go run cmd/currex/main.go --csv       # Экспорт курсов в CSV
```

## 🎮 **Управление в TUI (--ui)**

```
🔽/🔼 Стрелки     - листать валюты в dropdown
⏎ Enter         - подтвердить выбор валюты
Tab             - переключение полей
←→ Стрелки      - выбор кнопок
```

```
┌─ 💱 CurEx Конвертер ──────────────┐
│ Из: [USD ▼]                       │
│ В:  [RUB ▼]                       │
│ Сумма:                       │ [myfin](https://myfin.by/currency/cb-rf)
│ [💱 Конверт] [📈 График] [📊 CSV]  │
└───────────────────────────────────┘
✅ 100.00 USD = 9200.50 RUB
```

## 🛠 **Структура проекта**

```
CurEx/
├── cmd/currex/
│   └── main.go           # Точка входа CLI + TUI
├── internal/
│   ├── api/     api.go   # HTTP ЦБ РФ + XML парсинг
│   ├── cache/   cache.go # Кэш курсов (1 час)
│   └── debugger debugger.go # Debug логи
├── cache/        (auto)  # Кэш-файлы
├── go.mod                # Зависимости
└── README.md
```

## 📋 **Полные примеры команд**

| Команда | Результат |
|---------|-----------|
| `--amount 100 --from USD --to RUB` | `100 USD = 9200 RUB` |
| `--amount 50 --from EUR --to USD` | `50 EUR = 54.20 USD` |
| `--ui` | Графический интерфейс |
| `--list` | AUD, AZN, EUR, GBP, USD... |
| `--csv` | `rates.csv` с курсами |

## 🔧 **Установка зависимостей**

```bash
go mod tidy
```

**go.mod:**
```go
module github.com/maypress/CurEx
go 1.22

require (
    github.com/gdamore/tcell/v2 v2.6.0
    github.com/rivo/tview v0.42.0
)
```

## 🌐 **API источник**

**XML ЦБ РФ:** `https://www.cbr-xml-daily.ru/daily_utf8.xml` [web:6]

```xml
<Valute>
  <CharCode>USD</CharCode>
  <Nominal>1</Nominal>
  <Value>92,50</Value>  <!-- за 1 USD -->
</Valute>
```

## 📁 **Автогенерируемые файлы**

| Файл | Описание |
|------|----------|
| `cache/rates.cache` | Кэш курсов (обновление каждые 60 мин) |
| `rates.csv` | Экспорт всех 50+ курсов |

## 🐛 **ЧаВо**

| ❌ Проблема | ✅ Решение |
|-------------|------------|
| `go: command-line arguments` | `go run cmd/currex/main.go --amount 100 --from USD --to RUB` |
| `TUI глючит` | `export TERM=xterm-256color` |
| `Курсы = 0` | `rm -rf cache/; go run ...` |
| `403 Forbidden` | Уже использовано зеркало CBR |

## 🤝 **Контрибьютинг**

```bash
# 1. Форк + клон
git clone YOUR-FORK
cd CurEx

# 2. Фича-бранч
git checkout -b feat/historical-dates

# 3. Commit + PR
git commit -m "feat: add --date YYYY-MM-DD support"
git push origin feat/historical-dates
```

---
