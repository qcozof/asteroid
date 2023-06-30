package service

import (
	"fmt"
	"github.com/qcozof/asteroid/global"
	"github.com/qcozof/asteroid/model"
	"github.com/qcozof/asteroid/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func InitService(siteModel model.SiteModel, siteDirName string) error {
	//I.store original files to repository
	//remove repo siteDirName data

	global.InfoToChan(siteDirName, " cleaning dir and data ...")
	if err := Uninstall(siteDirName); err != nil {
		return err
	}
	global.InfoToChan(siteDirName, " cleaned")

	//make repo && isolation dir
	if err := utils.Mkdir(global.RepositoryDir + siteDirName); err != nil {
		return err
	}

	if err := utils.Mkdir(global.IsolationDir + siteDirName); err != nil {
		return err
	}

	if err := utils.Mkdir(global.BackupDir + siteDirName); err != nil {
		return err
	}

	global.InfoToChan(siteDirName, " indexing files...")
	//add siteDirName data to repo

	beginTime := time.Now()
	global.InfoToChan(siteDirName, " it may take some minutes ...")
	if err := copyToRepository(siteModel, siteDirName, siteModel.SiteDir, global.RepositoryDir+siteDirName); err != nil {
		return err
	}
	minutes := time.Now().Sub(beginTime).Minutes()

	global.InfoHighlightToChan(siteDirName, " indexed.", "Spent", minutes, "minutes.")
	return nil
}

func copyToRepository(siteModel model.SiteModel, siteDirName, src, dest string) error {
	// 获取源文件夹的信息
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	/*	// 创建目标文件夹
		if err := os.MkdirAll(dest, fi.Mode()); err != nil {
			return err
		}*/
	// 获取源文件夹中的所有文件和文件夹
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {

			//排除的目录
			if mustExcludeDir(srcPath, siteModel.ExcludeDir) {
				continue
			}

			// 递归拷贝文件夹
			if err := copyToRepository(siteModel, siteDirName, srcPath, destPath); err != nil {
				return err
			}
		} else {

			//包含的后缀
			if !mustInclude(srcPath, siteModel.IncludeExt) {
				continue
			}

			// 创建目标文件夹(放這邊，避免創建空文件夾)
			if err := os.MkdirAll(dest, fi.Mode()); err != nil {
				return err
			}

			// 拷贝文件
			if err := utils.CopyFile(srcPath, destPath); err != nil {
				return err
			}
			fmt.Println(srcPath, " -> ", destPath)
			addToDb(srcPath, siteDirName)
		}
	}

	return nil
}

func addToDb(file, siteName string) {
	hash, err := utils.GetFileHash(file)
	if err != nil {
		global.ErrorToChan("GetFileHash:", err)
	}

	perm, err := utils.GetPerm(file)
	if err != nil {
		global.ErrorToChan("GetPerm:", err)
	}

	fileName, err := utils.GetLastDirName(file, false)
	if err != nil {
		global.ErrorToChan(fmt.Sprintf("utils.GetLastDirName %s err:", file), err)
	}

	fileList := model.FileListModel{
		Site:       siteName,
		FileDir:    filepath.Dir(file),
		FileName:   fileName,
		Hash:       hash,
		Perm:       perm,
		Policy:     "overwrite",
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	err = global.GormDB.Create(&fileList).Error
	if err != nil {
		global.ErrorToChan("global.GormDB.Create err:", err)
	}
}

func ListFolderFiles(dir string) []string {
	var sli []string
	// 获取源文件夹的信息
	_, err := os.Stat(dir)
	if err != nil {
		global.ErrorToChan("ListFolderFiles.os.Stat", err)
	}

	// 获取源文件夹中的所有文件和文件夹
	entries, err := os.ReadDir(dir)
	if err != nil {
		global.ErrorToChan("ListFolderFiles.os.ReadDir", err)
	}

	for _, entry := range entries {

		if !entry.IsDir() {

			//file := filepath.Join(dir, entry.Name())
			//sli = append(sli, file)

			sli = append(sli, entry.Name())
		}
	}

	return sli
}

func mustInclude(file string, sliIncludeExt []string) bool {
	if len(sliIncludeExt) == 0 {
		return true
	}

	file = strings.TrimSpace(file)
	for _, ext := range sliIncludeExt {
		if strings.HasSuffix(file, strings.TrimSpace(ext)) {
			return true
		}
	}
	return false
}

func mustExcludeDir(dir string, sliIncludeDir []string) bool {
	if len(sliIncludeDir) == 0 {
		return false
	}

	dir = strings.TrimSpace(dir)
	for _, exDir := range sliIncludeDir {
		if dir == strings.TrimSpace(exDir) {
			return true
		}
	}
	return false
}
