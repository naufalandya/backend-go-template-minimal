package model

// Generated structs from C:\Andya\Go\standard_template\mock\db.txt

type User struct {
  Id int32 `json:"id" validate:"required,min=1"`
  Name string `json:"name" validate:"required,min=1,max=255"`
  Email string `json:"email" validate:"required,min=1,max=255"`
}

type File struct {
  Id string `json:"id" validate:"required,min=1,max=255"`
  Name string `json:"name" validate:"required,min=1,max=255"`
  FolderId int32 `json:"folder_id" validate:"required,min=1"`
  UploadedBy int32 `json:"uploaded_by" validate:"required,min=1"`
  Path string `json:"path" validate:"required,min=1,max=255"`
  Size int32 `json:"size" validate:"required,min=1"`
  Type string `json:"type" validate:"required,min=1,max=255"`
  Createdat string `json:"createdat" validate:"required,min=1,max=255"`
  Updatedat string `json:"updatedat" validate:"required,min=1,max=255"`
  TglBerlaku string `json:"tgl_berlaku" validate:"required,min=1,max=255"`
  TglPenetapan string `json:"tgl_penetapan" validate:"required,min=1,max=255"`
  Title string `json:"title" validate:"required,min=1,max=255"`
  Nomor string `json:"nomor" validate:"required,min=1,max=255"`
  Tahun string `json:"tahun" validate:"required,min=1,max=255"`
  FileName string `json:"file_name" validate:"required,min=1,max=255"`
  IsTrash bool `json:"is_trash" validate:"required"`
  Download string `json:"download" validate:"required,min=1,max=255"`
}

type Relation struct {
  ImpacterFileId string `json:"impacter_file_id" validate:"required,min=1,max=255"`
  ImpactedFileId string `json:"impacted_file_id" validate:"required,min=1,max=255"`
  Type string `json:"type" validate:"required,min=1,max=255"`
}