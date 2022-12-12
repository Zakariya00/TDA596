package main

import "time"

/* Maintenance Functions, Runs in the background */

var stabilizationDelay time.Duration
var predeccesorCheckDelay time.Duration
var fixFingersDelay time.Duration

func (chord *ChordNode) backGroundProcesses() {
	go chord.backGroundStabilize()
	go chord.backGroundFix()
	go chord.backGroundCheck()
}

func (chord *ChordNode) backGroundStabilize() {
	for {
		chord.stabilize()
		time.Sleep(stabilizationDelay * time.Millisecond)
	}
}

func (chord *ChordNode) backGroundCheck() {

	for {
		chord.check_predecessor()
		time.Sleep(predeccesorCheckDelay * time.Millisecond)
	}
}

func (chord *ChordNode) backGroundFix() {
	for {
		chord.fix_fingers()
		time.Sleep(fixFingersDelay * time.Millisecond)
	}
}
