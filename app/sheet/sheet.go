package sheet

import (
	"fmt"
	"log"

	"github.com/commojun/nyanbot/app/constant"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type Sheet struct {
	Service *sheets.Service
}

func New() (*Sheet, error) {
	cfg := &jwt.Config{
		Email:        constant.GoogleClientEmail,
		PrivateKey:   []byte(constant.GooglePrivateKey),
		PrivateKeyID: constant.GooglePrivateKeyID,
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets"},
		TokenURL:     constant.GoogleTokenURL,
	}
	if cfg.TokenURL == "" {
		cfg.TokenURL = google.JWTTokenURL
	}

	sheetService, err := sheets.New(cfg.Client(oauth2.NoContext))
	if err != nil {
		return &Sheet{}, err
	}

	return &Sheet{
		Service: sheetService,
	}, nil
}

func (sheet *Sheet) Get(spreadsheetID string, sheetName string) (*sheets.ValueRange, error) {
	res, err := sheet.Service.Spreadsheets.Values.Get(spreadsheetID, sheetName).Do()
	if err != nil {
		return &sheets.ValueRange{}, err
	}

	return res, err
}

func (sheet *Sheet) Load() error {
	spreadsheetId := constant.SheetID

	_, err := sheet.Service.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Fatalf("Unable to get Spreadsheets. %v", err)
	}

	fmt.Printf("success!\n")

	return err
}
