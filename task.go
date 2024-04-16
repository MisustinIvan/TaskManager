package main

import (
	"github.com/charmbracelet/bubbles/list"
)

type Status int

const (
	Todo Status = iota
	Doing
	Done
)

func NextStatus(s Status) Status {
	if s == Done {
		return Todo
	} else {
		return s + 1
	}
}

func PrevStatus(s Status) Status {
	if s == Todo {
		return Done
	} else {
		return s - 1
	}
}

type ExportableModel struct {
	Lists []ExportableTaskList `json:"lists"`
}

func ExportModel(m Model) ExportableModel {
	lists := make([]ExportableTaskList, len(m.lists))

	for i, l := range m.lists {
		lists[i] = ExportTaskList(l)
	}

	return ExportableModel{
		Lists: lists,
	}
}

// ExportableTaskList is a struct that is used to export a TaskList
type ExportableTaskList struct {
	Title string           `json:"title"`
	Items []ExportableTask `json:"items"`
}

func ExportTaskList(l list.Model) ExportableTaskList {
	items := make([]Task, len(l.Items()))
	for i, item := range l.Items() {
		items[i] = item.(Task)
	}

	el := ExportableTaskList{
		Title: l.Title,
	}

	for _, item := range items {
		el.Items = append(el.Items, ExportTask(item))
	}

	return el
}

// ExportableTask is a struct that is used to export a Task
type ExportableTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func ExportTask(t Task) ExportableTask {
	return ExportableTask{
		Title:       t.title,
		Description: t.description,
		Status:      t.status,
	}
}

// Task implements the list.Item interface
// by implementing the FilterValue method
type Task struct {
	title       string
	description string
	status      Status
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
