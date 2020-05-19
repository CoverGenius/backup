package helpers

import (
	"github.com/CoverGenius/backup/base"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	Log       *logrus.Logger
	Notifiers []base.Notifier
	Config    *base.Config
)

func LogError(err error) {
	if err != nil {
		Config.Notifier.Status.Status = StringP("WARN")
		Log.Error(err)
	}
}

func LogErrorExit(err error) {
	if err != nil {
		Log.Error(err)
		RemoveDir(Config.TmpDir)
		Config.Notifier.Status.Status = StringP("FAILURE")
		Config.Notifier.Status.EndTime = TimeP(time.Now())
		Config.Notifier.Status.CalculateDuration()
		for _, n := range Notifiers {
			n.Notify(Config)
			n.Post(Config)
		}
		Log.Fatal(err)
	}
}
