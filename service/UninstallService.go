package service

import (
	"github.com/qcozof/asteroid/global"
	"github.com/qcozof/asteroid/model"
	"github.com/qcozof/asteroid/utils"
)

func Uninstall(siteDirName string) error {
	var err error

	if err = global.GormDB.Where(" site=?", siteDirName).Delete(&model.FileListModel{}).Error; err != nil {
		return err
	}

	if err = global.GormDB.Exec("VACUUM").Error; err != nil {
		return err
	}

	if err = utils.RemoveAll(global.RepositoryDir + siteDirName); err != nil {
		return err
	}

	if err = utils.RemoveAll(global.BackupDir + siteDirName); err != nil {
		return err
	}

	return utils.RemoveAll(global.IsolationDir + siteDirName)
}

func Reset(initSql string) error {
	var err error

	if err = global.GormDB.Exec(initSql).Error; err != nil {
		return err
	}

	if err = utils.RemoveAll(global.RepositoryDir); err != nil {
		return err
	}

	if err = utils.RemoveAll(global.BackupDir); err != nil {
		return err
	}

	if err = global.LogFile.Close(); err == nil {
		if err = utils.RemoveAll(global.LogDir); err != nil {
			return err
		}
	}

	return utils.RemoveAll(global.IsolationDir)
}
