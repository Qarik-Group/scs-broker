package logger

import "code.cloudfoundry.org/lager"

type null struct{}

func (logger *null) Info(msg string, data ...lager.Data) {}

func (logger *null) Error(msg string, err error, data ...lager.Data) {}

func (logger *null) Fatal(msg string, err error, data ...lager.Data) {
	// the conventional expectation across logging suites is that a fatal
	// log does a panic. since control flow may rely on this, we're doing
	// it here.
	panic("(null logger) " + msg + " - " + err.Error())
}
