package daemon

import (
	"fmt"
	beatrix "github.com/eXtern-OS/Beatrix"
	"github.com/eXtern-OS/core9-common/db"
	"github.com/eXtern-OS/core9-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
)

/*
This package implements eXternOS App Store daemon, that will continuously update database with the newest applications
*/

type Daemon struct {
	SnapLocked    bool
	FlatpakLocked bool
}

var D Daemon

func (d *Daemon) Run() {
	log.Println("[DAEMON] Starting")
	d.SnapLocked = true
	log.Println("[DAEMON] Fetching snaps")
	snaps, err := FetchSnap()
	if err != nil {
		go beatrix.SendError(fmt.Sprintf("Error occured while fetching snaps: %e", err), "AppStore4.Daemon.Run")
	}

	if err = db.DefaultClient.DeleteMany(bson.M{}, "AppStore", "Snaps"); err != nil {
		go beatrix.SendError(fmt.Sprintf("Error occured while deleting snaps %e", err), "AppStore4.Daemon.Run")
	}

	if err = db.DefaultClient.InsertMany(utils.ArrayToInterface(snaps), "AppStore", "Snaps"); err != nil {
		go beatrix.SendError(fmt.Sprintf("Error occured while inserting snaps %e", err), "AppStore4.Daemon.Run")
	}
	d.SnapLocked = false
	d.FlatpakLocked = true
	log.Println("[DAEMON] Fetching flatpaks")
	flatpaks, err := FetchFlatpak()
	if err != nil {
		go beatrix.SendError(fmt.Sprintf("Error occured while fetching flatpaks: %e", err), "AppStore4.Daemon.Run")
	}
	if err = db.DefaultClient.DeleteMany(bson.M{}, "AppStore", "Flatpaks"); err != nil {
		go beatrix.SendError(fmt.Sprintf("Error occured while deleting flatpaks %e", err), "AppStore4.Daemon.Run")
	}

	if err = db.DefaultClient.InsertMany(utils.ArrayToInterface(flatpaks), "AppStore", "Flatpaks"); err != nil {
		go beatrix.SendError(fmt.Sprintf("Error occured while inserting flatpaks %e", err), "AppStore4.Daemon.Run")
	}
	d.FlatpakLocked = false
	log.Println("[DAEMON] Finished")

}

func (d *Daemon) CheckRunning() {
	for d.SnapLocked || d.FlatpakLocked {

	}
}

// Exit function exists to prevent disrupting ongoing processes
func (d *Daemon) Exit() {
	for d.SnapLocked {

	}
	d.SnapLocked = true
	for d.FlatpakLocked {

	}
	d.FlatpakLocked = true
}

func StartDaemon() {
	D.Run()
	for utils.SleepHours(1) {
		D.Run()
	}
}
