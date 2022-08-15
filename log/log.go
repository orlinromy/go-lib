package log

import (
	"time"
	"errors"
	"fmt"
	"runtime/debug"
	"encoding/json"
)

type Log struct {
	config	string
	json	bool
}

type Line struct {
	Ts	time.Time	`json:"ts"`
	Scope	string		`json:"scope"`
	Msg	string		`json:"msg"`
	Stack	string		`json:"stack,omitempty"`
}

func New(logtype string) (Log, error) {
	var l Log
	var e error
	if !isValid(logtype) {
		e = errors.New("Invalid log type")
		return l, e
	}
	l.config = logtype
	l.json = true
	return l, e
}

func (l Log) Out(scope string, msg string) {
	if l.config != "empty" && l.config != "erroronly" {
		logPrint(scope, msg, l.json)
	}
}

func (l Log) Debug(scope string, msg string) {
	if l.config != "empty" {
		logPrint(scope, msg, l.json)
	}
}

func (l Log) Error(scope string, err error) {
	if err != nil && l.config != "empty" {
		ts := time.Now()
		if !l.json {
			fmt.Println(ts.Format(time.RFC3339), scope, err)
			fmt.Println(string(debug.Stack()));
		} else {
			jline, e := log2Json(ts, scope, err.Error(), string(debug.Stack()))
			if e != nil {
				fmt.Println(ts.Format(time.RFC3339), scope, err.Error())
			} else {
				fmt.Println(jline)
			}
		}
	}
}

func (l *Log) JsonDisable() {
	l.json = false
}

func (l *Log) JsonEnable() {
	l.json = true
}

func logPrint(scope string, msg string, j bool) {
	ts := time.Now()
	if !j {
		fmt.Println(ts.Format(time.RFC3339), scope, msg)
	} else {
		jline, e := log2Json(ts, scope, msg, "")
		if e != nil {
			fmt.Println(ts.Format(time.RFC3339), scope, msg)
		} else {
			fmt.Println(jline)
		}
	}
}

func isValid(logtype string) bool {
	switch logtype {
	case "", "standard", "empty", "erroronly":
		return true
	}
	return false
}

func log2Json(ts time.Time, scope string, msg string, stack string) (string, error) {
	line := &Line{
		Ts: ts,
		Scope: scope,
		Msg: msg,
	}
	if (stack != "") {
		line.Stack = stack
	}
	j, e := json.Marshal(line)
	if e != nil {
		return "", e
	}
	return string(j), nil
}
