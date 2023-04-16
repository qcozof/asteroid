package global

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/qcozof/asteroid/model"
	"github.com/qcozof/asteroid/utils"

	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var (
	GormDB  *gorm.DB
	Config  *model.ConfigModel
	LogDate time.Time
)

var (
	AsteroidDir string
	MiscDir     string

	RepositoryDir string
	IsolationDir  string
	BackupDir     string
)

var (
	BroadcastInfoList   = make(chan string, 100)
	BroadcastNoticeList = make(chan string, 100)
	BroadcastErrorList  = make(chan string, 100)
)

const LogDir = "logs/"
const dbName = "asteroid.db"

func InitDB() error {
	var err error
	dbPath := MiscDir + dbName
	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		return err
	}

	GormDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	if err = gormDBTables(GormDB); err != nil {
		return err
	}

	sqlDB, _ := GormDB.DB()
	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(1) //set to other conns will cause `database is locked (5) (SQLITE_BUSY)`

	return err
}

func gormDBTables(db *gorm.DB) error {
	err := db.AutoMigrate(
	// to do ...
	)
	return err
}

func InitProjDir() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	AsteroidDir = dir + "/"
	MiscDir = AsteroidDir + "misc/"
	return err
}

func InitConfig() error {
	/*	bt, err := os.ReadFile(MiscDir + "config.json")
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(bt, &Config); err != nil {
			log.Println("InitConfig.json.Unmarshal err:", err)
		}*/
	configFile := MiscDir + "config.yaml"
	if !utils.ExistsFile(configFile) {
		return errors.New(configFile + " not exists.")
	}

	bt, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(bt, &Config); err != nil {
		return err
	}

	IsolationDir = Config.AsteroidDataDir + "/isolation/"
	RepositoryDir = Config.AsteroidDataDir + "/repository/"
	BackupDir = Config.AsteroidDataDir + "/backup/"
	return err
}

func InitLog() error {
	logDir := AsteroidDir + LogDir
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		return err
	}

	LogDate = time.Now()
	logFileFmt := fmt.Sprintf("%s%s.log", logDir, LogDate.Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileFmt, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	//defer logFile.Close()

	//同时输出到控制台和文件中
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	//只输出到文件中
	//log.SetOutput(logFile)

	log.SetPrefix("[asteroid] ")
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	return err
}

func ErrorToChan(msg string, err error) {
	BroadcastErrorList <- msg + err.Error()
}

func InfoToChan(obj ...interface{}) {
	BroadcastInfoList <- fmt.Sprintf("%+v", obj)
}

func NoticeToChan(obj ...interface{}) {
	BroadcastNoticeList <- fmt.Sprintf("%+v", obj)
}
