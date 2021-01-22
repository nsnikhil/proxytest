package config

import (
	"github.com/stretchr/testify/mock"
	"path"
)

type LogFileConfig interface {
	FileName() string
	FileDir() string
	FilePath() string
	FileMaxSizeInMb() int
	FileMaxBackups() int
	FileMaxAge() int
	FileWithLocalTimeStamp() bool
}

type appLogFileConfig struct {
	name               string
	dir                string
	maxSizeInMb        int
	maxBackups         int
	maxAge             int
	withLocalTimeStamp bool
}

func (lfc appLogFileConfig) FileName() string {
	return lfc.name
}

func (lfc appLogFileConfig) FileDir() string {
	return lfc.dir
}

func (lfc appLogFileConfig) FilePath() string {
	return path.Join(lfc.dir, lfc.name)
}

func (lfc appLogFileConfig) FileMaxSizeInMb() int {
	return lfc.maxSizeInMb
}

func (lfc appLogFileConfig) FileMaxBackups() int {
	return lfc.maxBackups
}

func (lfc appLogFileConfig) FileMaxAge() int {
	return lfc.maxAge
}

func (lfc appLogFileConfig) FileWithLocalTimeStamp() bool {
	return lfc.withLocalTimeStamp
}

func newLogFileConfig() LogFileConfig {
	return appLogFileConfig{
		name:               getString("LOG_FILE_NAME"),
		dir:                getString("LOG_FILE_DIR"),
		maxSizeInMb:        getInt("LOG_FILE_MAX_SIZE_IN_MB"),
		maxBackups:         getInt("LOG_FILE_MAX_BACKUPS"),
		maxAge:             getInt("LOG_FILE_MAX_AGE"),
		withLocalTimeStamp: getBool("LOG_FILE_WITH_LOCAL_TIME_STAMP"),
	}
}

type MockLogFileConfig struct {
	mock.Mock
}

func (mock *MockLogFileConfig) FileName() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockLogFileConfig) FileDir() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockLogFileConfig) FilePath() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockLogFileConfig) FileMaxSizeInMb() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *MockLogFileConfig) FileMaxBackups() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *MockLogFileConfig) FileMaxAge() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *MockLogFileConfig) FileWithLocalTimeStamp() bool {
	args := mock.Called()
	return args.Bool(0)
}
