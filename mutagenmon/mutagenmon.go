package mutagenmon

import (
	"context"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/mutagen-io/mutagen/cmd/mutagen/daemon"
	daemon2 "github.com/mutagen-io/mutagen/pkg/daemon"
	"github.com/mutagen-io/mutagen/pkg/selection"
	serviceSync "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const IntervalSec = 2

type MutagenMon struct {
	menu     map[string]*systray.MenuItem
	states   map[string]*synchronization.State
	daemon   *grpc.ClientConn
	interval time.Duration
	bad int
	total int
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
		menu:     map[string]*systray.MenuItem{},
		states:   map[string]*synchronization.State{},
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
	states := map[string]*synchronization.State{}
	for _, state := range response.SessionStates {
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

func isBad(state *synchronization.State) bool {
	if state.Status != synchronization.Status_StagingBeta && state.Status != synchronization.Status_Watching && state.Status != synchronization.Status_WaitingForRescan {
		return true
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
	msg := state.Status.String()
	if len(state.Conflicts) == 1 {
		msg += fmt.Sprintf(", %d conflict", len(state.Conflicts))
	}
	if len(state.Conflicts) > 1 {
		msg += fmt.Sprintf(", %d conflicts", len(state.Conflicts))
	}
	item.SetTooltip(msg)
	log.Printf("[INFO] state changed to: %s", msg)
	if len(state.Conflicts) > 0 {
		item.SetIcon(Icon("conflict.png"))
		return
	}
	switch state.Status {
	case synchronization.Status_HaltedOnRootEmptied:
		fallthrough
	case synchronization.Status_HaltedOnRootDeletion:
		fallthrough
	case synchronization.Status_HaltedOnRootTypeChange:
		fallthrough
	case synchronization.Status_Disconnected:
		fallthrough
	case synchronization.Status_ConnectingAlpha:
		fallthrough
	case synchronization.Status_ConnectingBeta:
		item.SetIcon(Icon("disconnected.png"))
	case synchronization.Status_Watching:
		item.SetIcon(Icon("ok.png"))
	case synchronization.Status_Scanning:
		item.SetIcon(Icon("syncing.png"))
	case synchronization.Status_WaitingForRescan:
		item.SetIcon(Icon("ok.png"))
	case synchronization.Status_Reconciling:
		fallthrough
	case synchronization.Status_StagingAlpha:
		fallthrough
	case synchronization.Status_StagingBeta:
		fallthrough
	case synchronization.Status_Transitioning:
		fallthrough
	case synchronization.Status_Saving:
		item.SetIcon(Icon("syncing.png"))
	default:
		item.SetIcon(Icon("unknown.png"))
	}
}

func (self *MutagenMon) CheckStates(ctx context.Context, states map[string]*synchronization.State) error {
	var bad int
	for id, current := range states {
		old, ok := self.states[id]
		if !ok {
			item := systray.AddMenuItem(fmt.Sprintf("%s:%s", current.Session.Beta.Host, current.Session.Beta.Path), "")
			self.menu[id] = item
			self.UpdateMenuItem(self.menu[id], current)
			self.states[id] = current
			if isBad(current) {
				bad++
			}
			continue
		}
		if old.Status != current.Status || len(old.Conflicts) != len(current.Conflicts) {
			self.UpdateMenuItem(self.menu[id], current)
			self.states[id] = current
		}
		if isBad(current) {
			bad++
		}
	}
	for id, _ := range self.states {
		_, ok := states[id]
		if !ok {
			self.menu[id].Hide()
			delete(self.menu, id)
			delete(self.states, id)
		}
	}
	total := len(self.states)
	if bad != self.bad || total != self.total {
		systray.SetTitle(fmt.Sprintf("%d / %d", total-bad, total))
	}
	self.bad = bad
	self.total = total
	log.Printf("[DEBUG] states checked, goroutines: %d", runtime.NumGoroutine())
	return nil
}

func (self *MutagenMon) Run() {
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
	systray.SetTooltip("Mutagen Monitor by Andmed")
	mQuit := systray.AddMenuItem("Quit", "")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
	systray.AddSeparator()
	go self.Scheduler()
}
