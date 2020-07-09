package handler

import (
	"fmt"
	"strings"

	"github.com/chinx/gocls/api/google/source/change"
	"github.com/tealeg/xlsx"
)

var rowGroup = []string{"Number", "Project", "Module", "Submodule", "Subject", "Owner", "Owner Email", "Submitter", "Submitter Email", "Created", "Submitted", "Change Url"}

func Statistics(list change.List, xf *xlsx.File) (err error) {
	if len(list) == 0 {
		return nil
	}
	status := strings.ToLower(list[0].Status)
	sheet, ok := xf.Sheet[status]
	if !ok {
		sheet, err = addSheet(xf, status)
		if err != nil {
			fmt.Printf("add new sheet failed: %s", err)
			return
		}
	}

	for _, v := range list {
		addRow(sheet, v)
	}
	return nil
}

func addSheet(xf *xlsx.File, name string) (sheet *xlsx.Sheet, err error) {
	sheet, err = xf.AddSheet(name)
	if err != nil {
		return
	}
	row := sheet.AddRow()
	for _, v := range rowGroup {
		row.AddCell().SetValue(v)
	}
	return
}

func addRow(sheet *xlsx.Sheet, val change.Change) {
	if strings.Index(val.Subject, "Revert \"") != -1 {
		return
	}

	module, submodule := splitSubject(val.Subject)
	if module == "" {
		module = val.Project
	}

	row := sheet.AddRow()
	for _, v := range rowGroup {
		item := row.AddCell()
		switch v {
		case "Number":
			item.SetValue(sheet.MaxRow - 1)
		case "Project":
			item.SetValue(val.Project)
		case "Module":
			item.SetValue(module)
		case "Submodule":
			item.SetValue(submodule)
		case "Subject":
			item.SetValue(val.Subject)
		case "Owner":
			item.SetValue(val.Owner.Name)
		case "Owner Email":
			item.SetValue(val.Owner.Email)
		case "Submitter":
			if val.Submitter != nil {
				item.SetValue(val.Submitter.Name)
			}
		case "Submitter Email":
			if val.Submitter != nil {
				item.SetValue(val.Submitter.Email)
			}
		case "Created":
			item.SetValue(val.Created)
		case "Submitted":
			item.SetValue(val.Submitted)
		case "Change Url":
			item.SetValue(val.ChangeURL())
		}
	}
	return
}

func splitSubject(value string) (module, submodule string) {
	last := len(value) - 1
	switch value[0] {
	case '`', '"', '\'':
		if value[last] == value[0] {
			value = value[1:last]
		} else {
			value = value[1:]
		}
	}

	index := strings.Index(value, ":")
	if index == -1 {
		return
	}

	module = value[:index]
	last = len(module) - 1
	switch module[last] {
	case '`', '"', '\'':
		value = value[:last]
	}

	if module[0] == '[' {
		i := strings.Index(module, "]")
		module = module[i+2:]
	}


	module = strings.Split(module, ",")[0]

	dividing := strings.Index(module, "/")
	if dividing >= 0 {
		submodule = module[dividing+1:]
		module = module[:dividing]
	}
	return
}
