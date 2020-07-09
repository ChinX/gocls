package main

import (
	"fmt"
	"time"

	"github.com/chinx/gocls/api/google/source"
	"github.com/chinx/gocls/api/google/source/change"
	"github.com/chinx/gocls/handler"
	"github.com/tealeg/xlsx"
)

func main() {
	source.SetProxy("socks5://127.0.0.1:1080")
	xf := xlsx.NewFile()
	ch := make(chan error)
	go openCLs(ch, xf)
	go mergedCLs(ch, xf)
	go abandonedCLs(ch, xf)

	for num := 3; num > 0; num-- {
		if err := <-ch; err != nil {
			fmt.Println(err)
		}
	}

	path := fmt.Sprintf("/root/chan/go/src/github.com/chinx/gocls/go_cls_%s.xlsx",
		time.Now().Format("2006-01-02_15:04:05"))

	if err := xf.Save(path); err != nil {
		fmt.Println(err)
		return
	}
}

func openCLs(c chan error, xf *xlsx.File) {
	query := change.NewQuery()
	query.Project("go")
	query.Open()
	c <- statisticsCLs(query, xf)
}

func mergedCLs(c chan error, xf *xlsx.File) {
	query := change.NewQuery()
	query.Project("go")
	query.Merged()
	c <- statisticsCLs(query, xf)
}

func abandonedCLs(c chan error, xf *xlsx.File) {
	query := change.NewQuery()
	query.Project("go")
	query.Abandoned()
	c <- statisticsCLs(query, xf)
}

func statisticsCLs(query *change.Query, xf *xlsx.File) error {
	for {
		changeList := change.List{}
		err := source.Do(query, &changeList)
		if err != nil {
			fmt.Println(err)
			break
		}

		err = handler.Statistics(changeList, xf)
		if err != nil {
			return err
		}
		query.NextPage()
	}
	return nil
}
