package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Report represents a grade report.
type Report struct {
	User    string
	Grade   float64
	Comment string
	Files   []string
}

// Assignment represents an assignment and all reports for it.
type Assignment struct {
	Name    string
	Reports []Report
}

// FilterUser filters out the reports that are not for the
// given user, and returns the results as a new Assignment.
func (a *Assignment) FilterUser(user string) *Assignment {
	aa := &Assignment{Name: a.Name}
	for _, r := range a.Reports {
		if r.User != user {
			continue
		}
		rr := Report{
			User:    r.User,
			Grade:   r.Grade,
			Comment: r.Comment,
			Files:   append([]string(nil), r.Files...),
		}
		aa.Reports = append(aa.Reports, rr)
	}
	return aa
}

// DB represents the entire database.
type DB struct {
	Asgn []Assignment
}

// FilterDB filters out the assignments that do not include
// a report for the given user, and for those that do, filters
// out all of the reports that are not for the given user;
// the result is returned as a new DB.
func (d *DB) FilterUser(user string) *DB {
	dd := new(DB)
	for _, a := range d.Asgn {
		aa := a.FilterUser(user)
		if len(aa.Reports) > 0 {
			dd.Asgn = append(dd.Asgn, *aa)
		}
	}
	return dd
}

// ReadDB reads the database from the named file.
func ReadDB(path string) (*DB, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	var db DB
	err = d.Decode(&db)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

// WriteDB writes the database back to the named
// file, overwriting the old database if it existed.
func WriteDB(db *DB, path string) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	e := json.NewEncoder(f)
	return e.Encode(*db)
}

// DestroyDB reinitializes the database to empty.
func DestroyDB(path string) error {
	os.Remove(path)
	os.Create(path)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	// print a blank JSON object so the json parser doesn't get sad
	// the next time we try to parse this
	fmt.Fprintln(f, "{}")
	return nil
}
