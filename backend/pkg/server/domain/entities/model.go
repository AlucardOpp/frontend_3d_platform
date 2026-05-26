package entities

import (
	"html"
	"strings"
	"time"

	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
)

type Model struct {
	ID                    uint64    `json:"id"`
	UserID                uint64    `json:"user_id"`
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
	Keywords              string    `json:"keywords"`
	FullDescriptionFileID *uint64   `json:"full_description_file_id,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ModelRequest struct {
	Title                 string   `json:"title"`
	Description           string   `json:"description"`
	Keywords              string   `json:"keywords"`
	FilesId               []uint64 `json:"files_id"`
	FullDescriptionFileId *uint64  `json:"full_description_file_id"`
}

type ModelData struct {
	Model           Model       `json:"model"`
	User            PublicUser  `json:"author"`
	Files           SortedFiles `json:"files"`
	FullDescription *File       `json:"full_description,omitempty"`
}

func (modelReq *ModelRequest) NewModel() *Model {
	return &Model{
		Title:       modelReq.Title,
		Description: modelReq.Description,
		Keywords:    modelReq.Keywords,
	}
}

func (m *Model) BeforeUpdate() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.Keywords = html.EscapeString(strings.TrimSpace(m.Keywords))
	m.UpdatedAt = time.Now()
}

func (m *Model) Prepare() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.Keywords = html.EscapeString(strings.TrimSpace(m.Keywords))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

func (model *ModelRequest) Validate() string {
	if model.Title == "" {
		return constants.ModelTitleCantBeEmpty
	}
	return ""
}

func (model *Model) Validate() string {
	if model.Title == "" {
		return constants.ModelTitleCantBeEmpty
	}
	if model.UserID == 0 {
		return constants.UserIDInvalid
	}
	return ""
}
