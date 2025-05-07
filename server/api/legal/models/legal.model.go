package models

import "time"

type RelationRequest struct {
	Title        string `json:"title" validate:"required"`
	Reason       string `json:"reason,omitempty"`
	TargetFileID string `json:"target_file_id,omitempty"`
}

type RelatedFile struct {
	Reason       string `json:"reason,omitempty"`
	TargetFileID string `json:"target_file_id,omitempty"`
}

type UploadFileRequest struct {
	FolderID      int           `json:"folderId,omitempty"`
	Title         string        `json:"title,omitempty"`
	Nomor         string        `json:"nomor,omitempty"`
	Tahun         string        `json:"tahun,omitempty"`
	TglPenetapan  *time.Time    `json:"tgl_penetapan,omitempty"`
	TglBerlaku    *time.Time    `json:"tgl_berlaku,omitempty"`
	Pemrakarsa_ID int           `json:"pemrakarsa_id,omitempty"`
	Revoke        []RelatedFile `json:"revoke,omitempty"`
	Modify        []RelatedFile `json:"modify,omitempty"`
}

type File struct {
	Id           string    `json:"id" validate:"required,min=1,max=255"`
	Name         string    `json:"name" validate:"required,min=1,max=255"`
	FolderId     int32     `json:"folder_id" validate:"required,min=1"`
	UploadedBy   int32     `json:"uploaded_by" validate:"required,min=1"`
	Path         string    `json:"path" validate:"required,min=1,max=255"`
	Size         int32     `json:"size" validate:"required,min=1"`
	Type         string    `json:"type" validate:"required,min=1,max=255"`
	Createdat    time.Time `json:"createdat" validate:"required"`
	Updatedat    time.Time `json:"updatedat" validate:"required"`
	TglBerlaku   time.Time `json:"tgl_berlaku" validate:"required"`
	TglPenetapan time.Time `json:"tgl_penetapan" validate:"required"`
	Title        string    `json:"title" validate:"required,min=1,max=255"`
	Nomor        string    `json:"nomor" validate:"required,min=1,max=255"`
	Tahun        string    `json:"tahun" validate:"required,min=1,max=255"`
	FileName     string    `json:"file_name" validate:"required,min=1,max=255"`
	IsTrash      bool      `json:"is_trash" validate:"required"`
	Download     string    `json:"download" validate:"required,min=1,max=255"`
	IndexText    string    `json:"index_text" validate:"required,min=1,max=255"`
}
