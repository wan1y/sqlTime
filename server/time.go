package server

import (
	"fmt"
	_ "gorm.io/driver/mysql"
	"io"
	"math"
	"os"
	"time"
)

type sqlTimes struct {
	firstTime  float64
	secondTime float64
	thirdTime  float64
}

func startCompareTime(s *dsnAndName, sqls string) {

	var sqlTimes map[string]sqlTimes
	for i := range s.servers {
		if i == 0 {
			sqlTimes = getStandardTime(s.getServers(i), sqls)
		} else {
			getTime(s.getServers(i), sqlTimes, sqls)
		}
	}
}

func getStandardTime(servers *servers, sqls string) map[string]sqlTimes{
	//db, err := sql.Open("mysql", servers.dsn)

	var allTime string
	var errorTime string
	sqlAndTime := map[string]sqlTimes{}
	times := sqlTimes{}
	for j := 0; j < 3; j++ {
		t1 := time.Now()
		if err := servers.dsnDB.Exec(sqls); err != nil {
			fmt.Println(err)
		}
		elapsed := time.Since(t1)
		if j == 0 {
			times.firstTime = elapsed.Seconds()
		} else if j == 1 {
			times.secondTime = elapsed.Seconds()
		} else if j == 2 {
			times.thirdTime = elapsed.Seconds()
		}
	}

	sqlAndTime[sqls] = times


	/**
	Time comparison standard
	*/
	for k, v := range sqlAndTime{
		allTime = "sql: " + k + "\n" + "first: " + fmt.Sprintf("%.3f", v.firstTime) +
			"s    second: " + fmt.Sprintf("%.3f", v.secondTime) + "s    third: " +
			fmt.Sprintf("%.3f", v.thirdTime)+ "s\n\n"

		if math.Abs(v.firstTime - v.secondTime) > 3 || math.Abs(v.firstTime - v.thirdTime) > 3 || math.Abs(v.secondTime - v.thirdTime) > 3 {
			errorTime = "sql: " + k + "\n" + "first: " + fmt.Sprintf("%.3f", v.firstTime) +
				"s    second: " + fmt.Sprintf("%.3f", v.secondTime) + "s    third: " +
				fmt.Sprintf("%.3f", v.thirdTime)+ "s\n\n"
		}
	}
	var err error
	var f *os.File
	if checkFileIsExist("StandardTime.txt") {
		f, err = os.OpenFile("StandardTime.txt", os.O_APPEND|os.O_WRONLY, 0666)
		fmt.Println("StandardTime.txt exists")
	} else {
		f, err = os.Create("StandardTime.txt")
		fmt.Println("StandardTime.txt don't exists")
	}
	defer f.Close()
	 _, err = io.WriteString(f, allTime)
	 if err != nil {
	 	panic(err)
	 }
	if checkFileIsExist(servers.dsnFileName) {
		f, err = os.OpenFile(servers.dsnFileName, os.O_APPEND|os.O_WRONLY, 0666)
		fmt.Println(servers.dsnFileName + " exists")
	} else {
		f, err = os.Create(servers.dsnFileName)
		fmt.Println(servers.dsnFileName + " don't exists")
	}
	defer f.Close()
	_, err = io.WriteString(f, errorTime);
	if err != nil {
		panic(err)
	}
	return sqlAndTime
}


func getTime(servers *servers, standardTime map[string]sqlTimes, sqls string){

	var errorTime string
	sqlAndTime := map[string]sqlTimes{}
	times := sqlTimes{}
	for j := 0; j < 3; j++ {
		t1 := time.Now()
		if err := servers.dsnDB.Exec(sqls); err != nil {
			fmt.Println(err)
		}
		elapsed := time.Since(t1)
		if j == 0 {
			times.firstTime = elapsed.Seconds()
		} else if j == 1 {
			times.secondTime = elapsed.Seconds()
		} else if j == 2 {
			times.thirdTime = elapsed.Seconds()
		}
	}
	sqlAndTime[sqls] = times


	/**
	Time comparison
	*/
	for k, v := range sqlAndTime{
		if math.Abs(v.firstTime - v.secondTime) > 3 || math.Abs(v.firstTime - v.thirdTime) > 3 || math.Abs(v.secondTime - v.thirdTime) > 3 ||
			math.Abs(v.firstTime - standardTime[k].firstTime) > 3 || math.Abs(v.secondTime - standardTime[k].secondTime) > 3 ||
			math.Abs(v.thirdTime - standardTime[k].thirdTime) > 3{
			errorTime = "sql: " + k + "\n" + "firstExp: " + fmt.Sprintf("%.3f", standardTime[k].firstTime) +
				"s    secondExp: " + fmt.Sprintf("%.3f", standardTime[k].secondTime) + "s    thirdExp: " +
				fmt.Sprintf("%.3f", standardTime[k].thirdTime)+ "s\n" + "firstGot: " + fmt.Sprintf("%.3f", v.firstTime) +
				"s    secondGot: " + fmt.Sprintf("%.3f", v.secondTime) + "s    thirdGot: " +
				fmt.Sprintf("%.3f", v.thirdTime)+ "s\n\n"
		}
	}

	var err error
	var f *os.File
	if checkFileIsExist(servers.dsnFileName) {
		f, err = os.OpenFile(servers.dsnFileName, os.O_APPEND|os.O_WRONLY, 0666)
		fmt.Println(servers.dsnFileName + " exists")
	} else {
		f, err = os.Create(servers.dsnFileName)
		fmt.Println(servers.dsnFileName + " don't exists")
	}
	defer f.Close()
	_, err = io.WriteString(f, errorTime);
	if err != nil {
		panic(err)
	}
}


func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}