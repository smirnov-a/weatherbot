package message

import (
	"bytes"
	"context"
	"github.com/chromedp/chromedp"
	"html/template"
	"os"
	"path/filepath"
	"time"
	"weatherbot/i18n"
	"weatherbot/internal/weather"
)

func GenerateWeatherHtm(data *weather.WeatherData, templatePath string) (string, error) {
	t, err := template.New("base").Funcs(template.FuncMap{
		"T": func(messageID string) string {
			return i18n.Translate(messageID)
		},
		"greaterThan": func(a, b float64) bool {
			return a > b
		},
	}).ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	templateName := filepath.Base(templatePath)
	var tplBuffer bytes.Buffer
	if err := t.ExecuteTemplate(&tplBuffer, templateName, data); err != nil {
		return "", err
	}

	return tplBuffer.String(), nil
}

func RenderHTMLToImage(htmlContent string, outputPath string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	tmpFile, err := os.CreateTemp("", "*.html")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	_, err = tmpFile.Write([]byte(htmlContent))
	if err != nil {
		return err
	}

	tmpFilePath := tmpFile.Name()

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("file://" + tmpFilePath),
		chromedp.Sleep(3 * time.Second), // Wait for the page to render
		chromedp.FullScreenshot(&buf, 90),
	}); err != nil {
		return err
	}

	// Save the screenshot to the specified output path
	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return err
	}

	return nil
}
