package tools

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func GetSheetsServiceWithServiceAccount() (*sheets.Service, error) {
	// Path to your Service Account credentials JSON file
	credentialsFile := "google_service_account.json"

	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
		return nil, err
	}
	return sheetsService, nil
}

func CreateNewSheet(spreadsheetID, sheetName string) error {
	ggSheetService, err := GetSheetsServiceWithServiceAccount()
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}
	// Create a new sheet
	requests := []*sheets.Request{
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: sheetName,
				},
			},
		},
	}

	_, err = ggSheetService.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		log.Printf("Unable to create sheet: %v", err)
		return err
	}
	return nil
}

func WriteSheetAtRange(spreadsheetID string, writeRange string, data [][]interface{}) error {
	ggSheetService, err := GetSheetsServiceWithServiceAccount()
	if err != nil {
		log.Printf("Unable to create Sheets service: %v", err)
		return err
	}
	// Write the data to the spreadsheet
	valueRange := &sheets.ValueRange{
		Values: data,
	}

	// Update the range in the sheet
	_, err = ggSheetService.Spreadsheets.Values.Update(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		log.Printf("Unable to write data to sheet: %v", err)
		return err
	}

	log.Printf("Data written successfully! %s\n", writeRange)
	return nil
}

func ResizeSheetFitData(spreadsheetID string, sheetIds ...int64) error {
	ggSheetService, err := GetSheetsServiceWithServiceAccount()
	if err != nil {
		log.Printf("Unable to create Sheets service: %v", err)
		return err
	}

	sheetReqs := make([]*sheets.Request, 0)
	for _, sheetId := range sheetIds {
		sheetReqs = append(sheetReqs, &sheets.Request{
			AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
				Dimensions: &sheets.DimensionRange{
					SheetId:    sheetId, // Use the sheet ID obtained earlier
					Dimension:  "COLUMNS",
					StartIndex: 0,  // Start column index (0-based)
					EndIndex:   20, // End column index (non-inclusive)
				},
			},
		})
	}
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: sheetReqs,
	}
	_, err = ggSheetService.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Do()
	if err != nil {
		log.Fatalf("Unable to auto-resize columns: %v", err)
	}
	fmt.Println("Columns resized successfully!")
	return nil
}
