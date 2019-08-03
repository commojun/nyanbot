package sheet

import (
	"fmt"
	"log"
	"net/http"

	"github.com/commojun/nyanbot/app/constant"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type Sheet struct {
	Client *http.Client
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

	return &Sheet{
		Client: cfg.Client(oauth2.NoContext),
	}, nil
}

func (sheet *Sheet) Load() error {
	spreadsheetId := constant.AlarmSheetID

	sheetService, err := sheets.New(sheet.Client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	_, err = sheetService.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Fatalf("Unable to get Spreadsheets. %v", err)
	}

	fmt.Printf("success!\n")

	return err
}
