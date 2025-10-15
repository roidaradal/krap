package authn

import (
	"time"

	"github.com/roidaradal/krap"
	"github.com/roidaradal/krap/daemon"
)

// ArchiveExpiredSessions Daemon
func Daemon_ArchiveExpiredSessions(interval int, timeScale time.Duration) {
	name := "ArchiveExpiredSessions"
	task := func() {
		rq, err := archiveExpiredSessions()
		krap.DisplayOutput(rq, err)
	}
	daemon.Run(name, task, interval, timeScale)
}

// DeleteArchivedSessions Daemon
func Daemon_DeleteArchivedSessions(marginDays uint, interval int, timeScale time.Duration) {
	name := "DeleteArchivedSessions"
	task := func() {
		rq, err := deleteArchivedSessions(marginDays)
		krap.DisplayOutput(rq, err)
	}
	daemon.Run(name, task, interval, timeScale)
}

// ExtendSessions Daemon
func Daemon_ExtendSessions(interval int, timeScale time.Duration) {
	name := "ExtendSessions"
	task := func() {
		rq, err := extendSessions()
		krap.DisplayOutput(rq, err)
	}
	daemon.Run(name, task, interval, timeScale)
}
