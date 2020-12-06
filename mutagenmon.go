package mutagenmon

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/mutagen-io/mutagen/cmd/mutagen/daemon"
	"github.com/mutagen-io/mutagen/cmd/mutagen/sync"
	daemon2 "github.com/mutagen-io/mutagen/pkg/daemon"
	"github.com/mutagen-io/mutagen/pkg/project"
	"github.com/mutagen-io/mutagen/pkg/selection"
	serviceSync "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization/core"
	"google.golang.org/grpc"
)

const IntervalSec = 2

var fatal = map[synchronization.Status]struct{}{
	synchronization.Status_HaltedOnRootEmptied:    {},
	synchronization.Status_HaltedOnRootDeletion:   {},
	synchronization.Status_HaltedOnRootTypeChange: {},
}

var disconnected = map[synchronization.Status]struct{}{
	synchronization.Status_Disconnected:           {},
	synchronization.Status_ConnectingBeta:         {},
}

var watching = map[synchronization.Status]struct{}{
	synchronization.Status_WaitingForRescan: {},
	synchronization.Status_Watching:         {},
}

var syncing = map[synchronization.Status]struct{}{
	synchronization.Status_ConnectingAlpha: {},
	synchronization.Status_Scanning:        {},
	synchronization.Status_Reconciling:     {},
	synchronization.Status_StagingAlpha:    {},
	synchronization.Status_StagingBeta:     {},
	synchronization.Status_Transitioning:   {},
	synchronization.Status_Saving:          {},
}

type Peer struct {
	menu     *systray.MenuItem
	state    *synchronization.State
	callback chan struct{}
}

type MutagenMon struct {
	peers     map[string]Peer
	callbacks map[string]chan struct{}
	daemon    *grpc.ClientConn
	interval  time.Duration
	bad       int
	conflict  int
	total     int
}

func is(state *synchronization.State, scope map[synchronization.Status]struct{}) bool {
	if state == nil {
		return false
	}
	var ok bool
	_, ok = scope[state.Status]
	return ok
}

func New() (*MutagenMon, error) {
	lock, err := daemon2.AcquireLock()
	if err == nil {
		// should not be here if daemon is running
		err2 := lock.Release()
		if err2 != nil {
			panic(err2)
		}
		return nil, fmt.Errorf("no daemon is running")
	}
	connection, err := daemon.CreateClientConnection(true, false)
	if err != nil {
		return nil, fmt.Errorf("connect to mutagen daemon: %v", err)
	}
	mutagenMon := MutagenMon{
		peers:    map[string]Peer{},
		daemon:   connection,
		interval: IntervalSec * time.Second,
	}
	return &mutagenMon, nil
}

func (self *MutagenMon) SessionStates(ctx context.Context) (map[string]*synchronization.State, error) {
	synchronizationService := serviceSync.NewSynchronizationClient(self.daemon)
	request := &serviceSync.ListRequest{
		Selection: &selection.Selection{All: true},
	}
	response, err := synchronizationService.List(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("get list of mutagen sessions: %v", err)
	}
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}
	states := map[string]*synchronization.State{}
	for _, state := range response.SessionStates {
		if state == nil {
			continue
		}
		states[state.Session.Identifier] = state
	}
	return states, nil
}

func (self *MutagenMon) Scheduler() {
	ctx := context.Background()
	ticker := time.NewTicker(self.interval)
	for range ticker.C {
		states, err := self.SessionStates(ctx)
		if err != nil {
			log.Printf("[WARN] get states: %s", err)
		}
		err = self.CheckStates(ctx, states)
		if err != nil {
			log.Printf("[WARN] check states: %s", err)
		}
	}
}

func hasConflicts(state *synchronization.State) bool {
	if state == nil{
		return false
	}
	return len(state.Conflicts) > 0
}

func Icon(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return defaultIcon
	}
	return b
}

func (self *MutagenMon) UpdateMenuItem(item *systray.MenuItem, state *synchronization.State) {
	if state == nil || item == nil {
		return
	}
	log.Printf("[DEBUG] update menu item")

	var msg string

	if hasConflicts(state) && state.Session.Configuration.SynchronizationMode == core.SynchronizationMode_SynchronizationModeOneWaySafe {
		msg = "Conflicts: (click to overwrite)\n"
		msg += "⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺⎺\n"
		var n int
		for _, conflict := range state.GetConflicts() {
			if n >= 20 {
				msg += fmt.Sprintf("... and %d more\n", len(state.Conflicts)-n)
				break
			}
			n++
			if conflict == nil {
				continue
			}
			for _, change := range conflict.BetaChanges {
				if change == nil {
					continue
				}
				path := change.Path
				if len(path) > 40 {
					path = path[:20] + " ... " + path[len(path)-15:]
				}
				msg += fmt.Sprintf("%s\n", path)
			}
		}
	} else {
		msg = fmt.Sprintf("%s\n", state.Status.String())
		msg += fmt.Sprintf("Click to force flush\n")
	}
	item.SetTooltip(msg)
	log.Printf("[INFO] state changed to: %s", msg)

	if is(state, disconnected) {
		item.SetIcon(Icon("disconnected.png"))
	} else if is(state, fatal) {
		item.SetIcon(Icon("fatal.png"))
	} else if is(state, syncing) {
		SetIfNoConflict(state, item, "syncing.png")
	} else if is(state, watching) {
		SetIfNoConflict(state, item, "ok.png")
	} else {
		SetIfNoConflict(state, item, "unknown.png")
	}
}

func SetIfNoConflict(state *synchronization.State, item *systray.MenuItem, name string) {
	if hasConflicts(state) {
		item.SetIcon(Icon("conflict.png"))
		return
	}
	item.SetIcon(Icon(name))
}
func (self *MutagenMon) CheckStates(ctx context.Context, states map[string]*synchronization.State) error {
	var bad int
	var conflict int
	for id, current := range states {
		if is(current, disconnected) || is(current, fatal){
			bad++
		}
		if hasConflicts(current) {
			conflict++
		}
		peer, ok := self.peers[id]
		if !ok {
			item := systray.AddMenuItem(fmt.Sprintf("%s:%s", current.Session.Beta.Host, current.Session.Beta.Path), "")
			ch := make(chan struct{})
			go func(id string) {
				for range ch {
					self.Resolve(id)
				}
			}(id)
			item.ClickedCh = ch
			self.UpdateMenuItem(item, current)
			self.peers[id] = Peer{
				menu:     item,
				state:    current,
				callback: ch,
			}
			continue
		}

		if peer.state.Status != current.Status || len(peer.state.Conflicts) != len(current.Conflicts) {
			self.UpdateMenuItem(peer.menu, current)
			peer.state = current
			self.peers[id] = peer
		}
	}
	for id, peer := range self.peers {
		_, ok := states[id]
		if !ok {
			peer.menu.Hide()
			close(peer.callback)
			delete(self.peers, id)
		}
	}
	total := len(self.peers)
	if bad != self.bad || total != self.total || conflict != self.conflict {
		systray.SetTitle(fmt.Sprintf(`%d-%d-%d`, total-conflict-bad, total-bad, total))
	}
	self.bad = bad
	self.conflict = conflict
	self.total = total
	log.Printf("[DEBUG] states checked, goroutines: %d", runtime.NumGoroutine())
	return nil
}

func (self *MutagenMon) Flush(state *synchronization.State) {
	if state == nil {
		return
	}
	labelSelector := fmt.Sprintf("%s=%s", project.LabelKey, state.Session.GetIdentifier())
	err := sync.FlushWithLabelSelector(labelSelector, true)
	if err != nil {
		log.Printf("[ERROR] Flush: %s", err)
	}
	return
}

func (self *MutagenMon) Resolve(id string) {
	peer, ok := self.peers[id]
	if !ok || peer.state.Session == nil {
		return
	}
	log.Printf("[DEBUG] resolve")

	err := peer.state.Session.EnsureValid()
	if err != nil {
		log.Printf("[ERROR] Resolve %s", err)
		return
	}

	out, _ := ioutil.TempFile("", "mutagenmon-*")
	defer out.Close()
	paths := []string{}
	for _, conflict := range peer.state.GetConflicts() {
		if conflict == nil {
			continue
		}
		for _, change := range conflict.BetaChanges {
			if change == nil {
				continue
			}
			paths = append(paths, change.GetPath())
		}
	}
	if len(paths) == 0 {
		self.Flush(peer.state)
		return
	}

	for _, path := range paths {
		_, err := out.WriteString(path + "\r\n")
		if err != nil {
			log.Printf("[ERROR] writing to tmp path: %s", err)
			return
		}
	}

	cmd := []string{
		"rsync",
		"-r",
		fmt.Sprintf("--files-from=%s", out.Name()),
		peer.state.Session.Alpha.Path,
		fmt.Sprintf("%s:%s", peer.state.Session.Beta.GetHost(), peer.state.Session.Beta.GetPath()),
	}
	log.Printf("[DEBUG] rsync args: %s", strings.Join(cmd, ", "))
	c := exec.Command(cmd[0], cmd[1:]...)
	b, err := c.CombinedOutput()
	if err != nil {
		log.Printf("[ERROR] Rsync: %s", err)
	} else {
		log.Printf("[DEBUG] Rsync: %s", string(b))
	}
	self.Flush(peer.state)
}

func (self *MutagenMon) Run() {
	log.Printf("[INFO] Mutagenmon")
	ep, err := os.Executable()
	if err != nil {
		log.Fatalln("os.Executable:", err)
	}
	err = os.Chdir(filepath.Join(filepath.Dir(ep), "..", "Resources"))
	if err != nil {
		log.Fatalln("os.Chdir:", err)
	}
	systray.Run(self.Init, nil)
}

func (self *MutagenMon) Init() {
	systray.SetIcon(Icon("icon.png"))
	systray.SetTooltip("no conflict / connected / total")
	mQuit := systray.AddMenuItem("Quit", "")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	systray.AddSeparator()
	go self.Scheduler()
}
