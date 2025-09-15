package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func WriteCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	//Add BOM for Excel Compatibility
	file.WriteString("\xEF\xBB\xBF")
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return nil
	}
	for page, pageData := range pages {
		csv_data := make([]string, 0)
		csv_data = append(csv_data, page)
		csv_data = append(csv_data, pageData.H1)
		csv_data = append(csv_data, pageData.FirstParagraph)
		csv_data = append(csv_data, strings.Join(pageData.OutgoingLinks, ";"))
		csv_data = append(csv_data, strings.Join(pageData.ImageURLs, ";"))
		err = writer.Write(csv_data)
		if err != nil {
			return err
		}
	}
	return nil
}
