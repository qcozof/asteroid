package global

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/qcozof/asteroid/model"
	"github.com/qcozof/asteroid/utils"

	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var (
	GormDB *gorm.DB
	Config *model.ConfigModel
)

var (
	AsteroidDir string
	MiscDir     string

	RepositoryDir string
	IsolationDir  string
	BackupDir     string

	ConfigFile string
)

var (
	BroadcastInfoList          = make(chan string)
	BroadcastNoticeList        = make(chan string)
	BroadcastErrorList         = make(chan string)
	BroadcastHighlightInfoList = make(chan string)
)

const LogDir = "logs/"
const dbName = "asteroid.db"

var mutex sync.Mutex

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
	ConfigFile = MiscDir + "config.yaml"
	if !utils.ExistsFile(ConfigFile) {
		return errors.New(ConfigFile + " not exists.")
	}

	bt, err := os.ReadFile(ConfigFile)
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

func InitLog(runOnce bool) error {
	var err error
	logDir := AsteroidDir + LogDir

	if runOnce {
		err = os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var logFile os.File
	for {
		if !runOnce {
			now := time.Now()
			tomorrow0Ux := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(time.Duration(tomorrow0Ux.Sub(now).Nanoseconds()) * time.Nanosecond) //10000s
			logFile.Close()
		}

		logFileFmt := fmt.Sprintf("%s%s.log", logDir, time.Now().Format("2006-01-02"))
		logFile, err := openFile(logFileFmt)
		if err != nil {
			fmt.Println("open log file failed, err:", err)
			return err
		}

		//defer logFile.Close()

		//both write to console and file
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		//only write to file
		//log.SetOutput(logFile)

		log.SetPrefix(" [asteroid] ")
		log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

		if runOnce {
			break
		}
	}

	return err
}

func openFile(logFileFmt string) (*os.File, error) {
	mutex.Lock()
	defer mutex.Unlock()
	return os.OpenFile(logFileFmt, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
}

func ErrorToChan(msg string, err error) {
	BroadcastErrorList <- msg + err.Error()
}

func InfoToChan(obj ...interface{}) {
	BroadcastInfoList <- fmt.Sprintf("%+v", obj)
}

func InfoHighlightToChan(obj ...interface{}) {
	BroadcastHighlightInfoList <- fmt.Sprintf("%+v", obj)
}

func NoticeToChan(obj ...interface{}) {
	BroadcastNoticeList <- fmt.Sprintf("%+v", obj)
}
