package main

import "time"

/* Maintenance Functions, Runs in the background */

var stabilizationDelay time.Duration
var predeccesorCheckDelay time.Duration
var fixFingersDelay time.Duration
var backupTimeDelay time.Duration

// backGroundProcesses spawn a new thread for methods to run in the background
func (chord *ChordNode) backGroundProcesses() {
	go chord.backGroundStabilize()
	go chord.backGroundFix()
	go chord.backGroundCheck()
	go chord.backupHandler()
}

// backGroundStabilize runs stabilize in loop til program termination. With set delay
func (chord *ChordNode) backGroundStabilize() {
	for {
		chord.stabilize()
		time.Sleep(stabilizationDelay * time.Millisecond)
	}
}

// backGroundCheck runs check_predecessor in loop til program termination. With set delay
func (chord *ChordNode) backGroundCheck() {

	for {
		chord.check_predecessor()
		time.Sleep(predeccesorCheckDelay * time.Millisecond)
	}
}

// backGroundFix runs fix_fingers in loop til program termination. With set delay
func (chord *ChordNode) backGroundFix() {
	for {
		chord.fix_fingers()
		time.Sleep(fixFingersDelay * time.Millisecond)
	}
}

// backupHandler runs backupFiles every set minutes
func (chord *ChordNode) backupHandler() {
	for {
		chord.backupFiles()
		chord.movBackups()
		time.Sleep(backupTimeDelay * time.Minute)
	}
}
