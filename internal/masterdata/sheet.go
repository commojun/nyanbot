package masterdata

import (
	"context"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

type sheetConfig struct {
	Email        string
	PrivateKey   string
	PrivateKeyID string
	TokenURL     string
}

type Sheet struct {
	Service *sheets.Service
}

func newSheet(ctx context.Context, cfg sheetConfig) (*Sheet, error) {
	jwtCfg := &jwt.Config{
		Email:        cfg.Email,
		PrivateKey:   []byte(cfg.PrivateKey),
		PrivateKeyID: cfg.PrivateKeyID,
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets"},
		TokenURL:     cfg.TokenURL,
	}
	if jwtCfg.TokenURL == "" {
		jwtCfg.TokenURL = google.JWTTokenURL
	}

	sheetService, err := sheets.NewService(ctx, option.WithHTTPClient(jwtCfg.Client(ctx)))
	if err != nil {
		return &Sheet{}, err
	}

	return &Sheet{
		Service: sheetService,
	}, nil
}

func (sheet *Sheet) Get(ctx context.Context, spreadsheetID string, sheetName string) ([][]any, error) {
	res, err := sheet.Service.Spreadsheets.Values.Get(spreadsheetID, sheetName).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return res.Values, nil
}
