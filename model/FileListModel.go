package model

import "io/fs"

type FilePolicy string

const (
	OverWrite FilePolicy = "overwrite"
	Ignore    FilePolicy = "ignore"
)

type FileListModel struct {
	FileId     int64       `gorm:"column:file_id;primary_key;autoIncrement:true;comment:file_id"`
	Site       string      `gorm:"column:site;comment:site"`
	FileDir    string      `gorm:"column:file_dir;comment:file_dir"`
	FileName   string      `gorm:"column:file_name;comment:file_name"`
	Hash       string      `gorm:"column:hash;comment:hash"`
	Perm       fs.FileMode `gorm:"column:perm;comment:perm"`
	Policy     FilePolicy  `gorm:"column:policy;comment:policy"`
	CreateTime int64       `gorm:"column:create_time;comment:create_time"`
	UpdateTime int64       `gorm:"column:update_time;comment:update_time"`
}

func (FileListModel) TableName() string {
	return "file_list"
}
