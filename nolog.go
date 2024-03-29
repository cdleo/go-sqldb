package sqldb

import (
	"io"
	"time"

	"github.com/cdleo/go-commons/logger"
)

type nolog struct {
}

func NewNoLogLogger() logger.Logger {
	return &nolog{}
}

func (nl *nolog) SetLogLevel(level string) error {
	return nil
}

func (nl *nolog) SetOutput(w io.Writer) {
}

func (nl *nolog) SetTimestampFunc(f func() time.Time) {
}

func (nl *nolog) WithRefID(refID string) logger.Logger {
	return nl
}

func (nl *nolog) Show(msg string) {
}
func (nl *nolog) Showf(msg string, v ...interface{}) {
}

func (nl *nolog) Fatal(err error, msg string) {
}
func (nl *nolog) Fatalf(err error, msg string, v ...interface{}) {
}

func (nl *nolog) Error(err error, msg string) {
}
func (nl *nolog) Errorf(err error, msg string, v ...interface{}) {
}

func (nl *nolog) Warn(msg string) {
}
func (nl *nolog) Warnf(msg string, v ...interface{}) {
}

func (nl *nolog) Info(msg string) {
}
func (nl *nolog) Infof(msg string, v ...interface{}) {
}

func (nl *nolog) Bus(msg string) {
}
func (nl *nolog) Busf(msg string, v ...interface{}) {
}

func (nl *nolog) Msg(msg string) {
}
func (nl *nolog) Msgf(msg string, v ...interface{}) {
}

func (nl *nolog) Dbg(msg string) {
}
func (nl *nolog) Dbgf(msg string, v ...interface{}) {
}

func (nl *nolog) Qry(msg string) {
}
func (nl *nolog) Qryf(msg string, v ...interface{}) {
}

func (nl *nolog) Trace(msg string) {
}
func (nl *nolog) Tracef(msg string, v ...interface{}) {
}
