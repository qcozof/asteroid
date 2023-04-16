package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func Mkdir(path string) error {
	return os.MkdirAll(path, 0777)
}

func RemoveAll(path string) error {
	err := os.RemoveAll(path)
	return err
}

func ExistsDir(dir string) bool {
	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}

	if !fi.IsDir() {
		return false
	}

	return true
}

func ExistsFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func GetLastDirName(src string, isDir bool) (string, error) {
	// 获取源文件夹的信息
	fi, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	if isDir && !fi.IsDir() {
		return "", errors.New(src + " is not a dir")
	}

	return fi.Name(), nil
}

func ListAllFiles(dir string) ([]string, error) {
	// 定义一个切片来存储文件列表
	var files []string

	// Walk函数会遍历指定目录下的所有文件，并调用匿名函数
	// 该匿名函数将接收到的文件名添加到上面定义的files列表中
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func CopyFile(src, dst string) error {
	// 检查源文件是否存在
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 检查目标文件夹是否存在，如果不存在则创建
	dstDir := filepath.Dir(dst)
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return err
		}
	}

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 拷贝文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func GetFileHash(path string) (string, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	// 确保函数退出时关闭文件
	defer file.Close()

	// 创建一个SHA256 hasher对象
	hasher := sha256.New()

	// 将文件数据写入哈希器
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	// 计算哈希值并转为16进制字符串
	hashBytes := hasher.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)

	return hashStr, nil
}

func GetPerm(name string) (fs.FileMode, error) {
	fileInfo, err := os.Stat(name)
	return fileInfo.Mode().Perm(), err
}

func SetPerm(name string, perm fs.FileMode) error {
	return os.Chmod(name, perm)
}

func ListFolderFiles(dir string) ([]string, error) {
	var sli []string
	// 获取源文件夹的信息
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	// 获取源文件夹中的所有文件和文件夹
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			sli = append(sli, entry.Name())
		}
	}

	return sli, nil
}
