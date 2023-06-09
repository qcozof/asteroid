package service

import (
	"fmt"
	"github.com/qcozof/asteroid/global"
	"github.com/qcozof/asteroid/model"
	"github.com/qcozof/asteroid/utils"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var mpSiteDir = map[string][]model.FileListModel{}
var mpSiteFileList = map[string][]model.FileListModel{}

func getFileDirListBySite(siteDirName string) (fileDirList []model.FileListModel, err error) {
	if val, exists := mpSiteDir[siteDirName]; exists && len(val) > 0 {
		return val, nil
	}

	if err = global.GormDB.Select("file_dir").Where("site=?", siteDirName).Group("file_dir").Find(&fileDirList).Error; err != nil {
		return fileDirList, err
	}

	mpSiteDir[siteDirName] = fileDirList
	return fileDirList, err
}

func getFileListBySite(fileDir string) (fileList []model.FileListModel, err error) {
	if val, exists := mpSiteFileList[fileDir]; exists && len(val) > 0 {
		return val, nil
	}

	if err := global.GormDB.Select("file_dir,file_name,hash,perm,policy").Where("file_dir=?", fileDir).Find(&fileList).Error; err != nil {
		return fileList, err
	}

	mpSiteFileList[fileDir] = fileList
	return fileList, err
}

func MonitorService(siteModel model.SiteModel, siteDir, siteDirName string) error {

	/*	var fileDirList []model.FileListModel
		if err := global.GormDB.Select("file_dir").Where("site=?", siteDirName).Group("file_dir").Find(&fileDirList).Error; err != nil {
			return err
		}*/
	fileDirList, err := getFileDirListBySite(siteDirName)
	if err != nil {
		return err
	}

	repositoryDir := global.RepositoryDir

	for _, m := range fileDirList {
		/*		var fileList []model.FileListModel
				if err := global.GormDB.Select("file_dir,file_name,hash,perm,policy").Where("file_dir=?", m.FileDir).Find(&fileList).Error; err != nil {
					return err
				}*/
		fileList, err := getFileListBySite(m.FileDir)
		if err != nil {
			return err
		}

		//1).table files
		for _, f := range fileList {
			file := f.FileDir + "/" + f.FileName
			hs, err := utils.GetFileHash(file)
			if err != nil {
				global.NoticeToChan(fmt.Sprintf("utils.GetFileHash %s err:", file), err)
			}

			relDir := strings.ReplaceAll(f.FileDir, siteDir, "")

			if f.Hash != hs {
				global.InfoHighlightToChan("hash NOT equal:", file)

				switch f.Policy {
				case model.OverWrite:
					originalFile := strings.ReplaceAll(file, siteDir, repositoryDir+siteDirName) //fmt.Sprintf("%s/%s", repositoryDir, file)
					backupFile(file, fmt.Sprintf("%s/%s/%s/%s-%s", global.BackupDir, siteDirName, relDir, f.FileName, time.Now().Format("20060102150405")))
					replaceFile(originalFile, file, f.Perm)

				default:

				}
			} else {
				fmt.Println("hash equal:", file)
				time.Sleep(time.Nanosecond)
			}

			//2).not in table,remove
			sliFiles, err := utils.ListFolderFiles(f.FileDir)
			if err != nil {
				global.ErrorToChan("utils.ListFolderFiles err:", err)
			}

			for _, tmpFile := range sliFiles {
				if !mustInclude(tmpFile, siteModel.IncludeExt) {
					continue
				}

				if fileInTable(tmpFile, fileList) {
					continue
				}

				isolateFile(f.FileDir+"/"+tmpFile, global.IsolationDir+siteDirName+"/"+relDir+"/"+tmpFile)
			}
		}
	}

	return nil
}

func fileInTable(file string, sliGroupFilesTab []model.FileListModel) bool {
	for _, item := range sliGroupFilesTab {
		if item.FileName == file {
			return true
		}
	}
	return false
}

func replaceFile(originalFile, file string, originalPerm fs.FileMode) {
	if err := os.Chmod(file, 0644); err != nil {
		global.ErrorToChan("replaceFile.os.Chmod err:", err)
	}

	err := utils.CopyFile(originalFile, file)
	if err != nil {
		global.ErrorToChan("replaceFile.utils.CopyFile err:", err)
	} else {
		global.NoticeToChan(fmt.Sprintf("file restored:%s %s=>%s", filepath.Base(file), originalFile, file))
		if err = utils.SetPerm(file, originalPerm); err != nil {
			global.ErrorToChan("replace utils.SetPerm err:", err)
		}
	}
}

func backupFile(file, destFile string) {
	err := utils.CopyFile(file, destFile)
	if err != nil {
		global.ErrorToChan("backupFile.utils.CopyFile err:", err)
	} else {
		global.InfoToChan(fmt.Sprintf("backup file:%s => %s", file, destFile))
	}
}

func isolateFile(file, destFile string) {
	err := utils.CopyFile(file, destFile)
	if err != nil {
		global.ErrorToChan("utils.CopyFile err:", err)
	} else {
		global.NoticeToChan(fmt.Sprintf("isolate file: %s=>%s", file, destFile))
	}

	if err = utils.RemoveAll(file); err != nil {
		global.ErrorToChan("isolateFile utils.RemoveAll:", err)
	}
}
