package dto

import "encoding/json"

type SavePdfRequest struct {
	UserId                 int64
	CartId                 int64
	PublicationId          int64
	Logo                   json.RawMessage
	ExecutorParameters     json.RawMessage
	PresentationParameters json.RawMessage
	StyleTemplate          json.RawMessage
	Count                  int
	isSave                 bool
}

type LogoText struct {
	Value   string `json:"value"`
	Font    string `json:"name_font"`
	Bold    bool   `json:"bold"`
	Kursive bool   `json:"kursive"`
	Under   bool   `json:"under"`
}

type Logo struct {
	Square    string   `json:"logo_square"`
	Rectangle string   `json:"logo_rectangle"`
	LogoText  LogoText `json:"logo_text"`
}

type ExecutorParam struct {
	ShowLogo     bool   `json:"show_logo"`
	ShowName     string `json:"show_name"`
	ShowContacts string `json:"show_contacts"`
}

type ExecutorParameters struct {
	First ExecutorParam `json:"first"`
	All   ExecutorParam `json:"all"`
	Last  ExecutorParam `json:"last"`
}

type PresentationParameters struct {
	List     bool `json:"list"`
	OneByOne bool `json:"one_by_one"`
	Sum      bool `json:"sum"`
	Price    bool `json:"price"`
}

type StyleTemplate struct {
	TemplateID string `json:"id_template"`
	Color      string `json:"color,omitempty"` // если есть
}

type SaveRequest struct {
	UserId                 int64                  `json:"id_user"`
	CartId                 int64                  `json:"id_cart"`
	PublicationId          int64                  `json:"id_publication"`
	Logo                   Logo                   `json:"logo"`
	ExecutorParameters     ExecutorParameters     `json:"executor_parameters"`
	PresentationParameters PresentationParameters `json:"presentation_parameters"`
	StyleTemplate          StyleTemplate          `json:"style_template"`
	Count                  int                    `json:"count"`
}
