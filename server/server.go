package server

import (
	"bufio"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type dsnAndName struct {
	threads  int
	logLevel logger.LogLevel
	servers []*servers
	sqlFile string
}

type servers struct{
	dsn  string
	dsnFileName  string
}
func New(threads int, logLevel string, sqlFile string) *dsnAndName {
	rand.Seed(time.Now().UnixNano())

	level := logger.Info
	switch logLevel {
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		level = logger.Error
	case "silent":
		level = logger.Silent
	default:
		panic("unknown log level: " + logLevel)
	}

	logger.Default = logger.New(log.New(os.Stdout, "\n", log.LstdFlags), logger.Config{
		SlowThreshold: 100 * time.Millisecond,
		LogLevel:      level,
		Colorful:      true,
	})

	s := &dsnAndName{
		threads:  threads,
		logLevel: level,
		servers:   nil,
		sqlFile:  sqlFile,
	}
	return s
}

func (s *dsnAndName) SetDsnAndFileNames(dsns []string, fileNames []string) {
	for i := range dsns {
		var dsnAndName servers

		dsnAndName.dsn = dsns[i]
		dsnAndName.dsnFileName = fileNames[i]
		s.servers = append(s.servers, &dsnAndName)
	}
}


func (s *dsnAndName) getSqlFile() string {
	return s.sqlFile
}

func (s *dsnAndName) getServers(num int) *servers {
	return s.servers[num]
}

func (s *dsnAndName) CompareTime() {

	os.Remove("StandardTime.txt")
	for i := range s.servers {
		os.Remove(s.servers[i].dsnFileName)
	}

	file ,err := os.Open(s.sqlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	/*
		ScanLines (默认)
		ScanWords
		ScanRunes (遍历UTF-8字符非常有用)
		ScanBytes
	*/
	var sqls [] string
	for scanner.Scan() {
		sqls = append(sqls, scanner.Text())
	}

	var wg sync.WaitGroup
	ch := make(chan string, s.threads+1)
	for i := 0; i < s.threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case sql, ok := <-ch:
					if !ok {
						return
					}
					startCompareTime(s, sql)
				}
			}
		}()
	}
	for i := 0; i < len(sqls); i++ {
		ch <- sqls[i]
	}
	close(ch)
	wg.Wait()
}