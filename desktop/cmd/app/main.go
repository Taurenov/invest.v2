// Минимальное десктоп-окно (Fyne).
// В проде: Tauri + React, данные с Go API (см. desktop/ui — макет веб-слоя).
package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type txResponse struct {
	Data []struct {
		Kind        string  `json:"kind"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		Description string  `json:"description"`
	} `json:"data"`
}

type row struct {
	Kind, Description, Currency string
	Amount                       float64
}

var defaultRows = []row{
	{Kind: "income", Description: "Зарплата", Amount: 85000, Currency: "₽"},
	{Kind: "expense", Description: "Продукты", Amount: 3200, Currency: "₽"},
	{Kind: "expense", Description: "Транспорт", Amount: 890, Currency: "₽"},
}

func main() {
	rows := loadTransactions()

	a := app.NewWithID("com.finhelper.desktop")
	a.Settings().SetTheme(&finTheme{})
	w := a.NewWindow("Fin Helper")
	w.Resize(fyne.NewSize(960, 640))
	w.SetMaster()

	balance := widget.NewLabel("124 500 ₽")
	balance.TextStyle = fyne.TextStyle{Bold: true}
	balance.Alignment = fyne.TextAlignCenter

	balanceCard := container.NewVBox(
		widget.NewLabelWithStyle("Баланс", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		balance,
		widget.NewLabel("Текущий баланс"),
	)

	list := widget.NewList(
		func() int { return len(rows) },
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), widget.NewLabel(""))
		},
		func(i int, o fyne.CanvasObject) {
			item := rows[i]
			box := o.(*fyne.Container)
			left := box.Objects[0].(*widget.Label)
			right := box.Objects[1].(*widget.Label)
			sign := "+"
			if item.Kind == "expense" {
				sign = "−"
			}
			left.SetText(item.Description)
			right.SetText(fmt.Sprintf("%s%s %s", sign, formatMoney(item.Amount), item.Currency))
		},
	)

	header := widget.NewLabelWithStyle("Последние транзакции", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	w.SetContent(container.NewPadded(container.NewBorder(
		container.NewPadded(balanceCard),
		nil, nil, nil,
		container.NewBorder(header, nil, nil, nil, list),
	)))
	w.ShowAndRun()
}

func loadTransactions() []row {
	apiURL := os.Getenv("FIN_API_URL")
	token := os.Getenv("API_TOKEN")
	if apiURL == "" || token == "" {
		return defaultRows
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return defaultRows
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return defaultRows
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var parsed txResponse
	if json.Unmarshal(body, &parsed) != nil || len(parsed.Data) == 0 {
		return defaultRows
	}

	out := make([]row, 0, len(parsed.Data))
	for _, t := range parsed.Data {
		cur := t.Currency
		if cur == "RUB" {
			cur = "₽"
		}
		out = append(out, row{Kind: t.Kind, Description: t.Description, Amount: t.Amount, Currency: cur})
	}
	return out
}

func formatMoney(v float64) string {
	return fmt.Sprintf("%.0f", v)
}

type finTheme struct{}

func (t *finTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x0f, G: 0x14, B: 0x19, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x3b, G: 0x82, B: 0xf6, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xe8, G: 0xed, B: 0xf4, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x1a, G: 0x23, B: 0x32, A: 0xff}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *finTheme) Font(style fyne.TextStyle) fyne.Resource  { return theme.DefaultTheme().Font(style) }
func (t *finTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(name) }
func (t *finTheme) Size(name fyne.ThemeSizeName) float32       { return theme.DefaultTheme().Size(name) }
