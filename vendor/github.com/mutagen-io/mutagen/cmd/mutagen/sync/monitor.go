package sync

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/fatih/color"

	"github.com/mutagen-io/mutagen/cmd"
	"github.com/mutagen-io/mutagen/cmd/mutagen/daemon"
	"github.com/mutagen-io/mutagen/pkg/grpcutil"
	selectionpkg "github.com/mutagen-io/mutagen/pkg/selection"
	synchronizationsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
	"github.com/mutagen-io/mutagen/pkg/synchronization"
)

func computeMonitorStatusLine(state *synchronization.State) string {
	// Build the status line.
	status := "Status: "
	if state.Session.Paused {
		status += color.YellowString("[Paused]")
	} else {
		// Add a conflict flag if there are conflicts.
		if len(state.Conflicts) > 0 {
			status += color.RedString("[Conflicts] ")
		}

		// Add a problems flag if there are problems.
		if len(state.AlphaProblems) > 0 || len(state.BetaProblems) > 0 {
			status += color.RedString("[Problems] ")
		}

		// Add an error flag if there is one present.
		if state.LastError != "" {
			status += color.RedString("[Errored] ")
		}

		// Add the status.
		status += state.Status.Description()

		// If we're staging and have sane statistics, add them.
		if (state.Status == synchronization.Status_StagingAlpha ||
			state.Status == synchronization.Status_StagingBeta) &&
			state.StagingStatus != nil {
			status += fmt.Sprintf(
				": %.0f%% (%d/%d)",
				100.0*float32(state.StagingStatus.Received)/float32(state.StagingStatus.Total),
				state.StagingStatus.Received,
				state.StagingStatus.Total,
			)
		}
	}

	// Done.
	return status
}

func monitorMain(command *cobra.Command, arguments []string) error {
	// Create a session selection specification that will select our initial
	// batch of sessions. From this batch, we'll determine which session to
	// monitor based on creation date. In any case, we only allow one
	// specification to be provided in order to enforce the notion that this is
	// a single-session command.
	if len(arguments) > 1 {
		return errors.New("multiple session specifications not allowed")
	}
	selection := &selectionpkg.Selection{
		All:            len(arguments) == 0 && monitorConfiguration.labelSelector == "",
		Specifications: arguments,
		LabelSelector:  monitorConfiguration.labelSelector,
	}
	if err := selection.EnsureValid(); err != nil {
		return errors.Wrap(err, "invalid session selection specification")
	}

	// Connect to the daemon and defer closure of the connection.
	daemonConnection, err := daemon.CreateClientConnection(true, true)
	if err != nil {
		return errors.Wrap(err, "unable to connect to daemon")
	}
	defer daemonConnection.Close()

	// Create a session service client.
	sessionService := synchronizationsvc.NewSynchronizationClient(daemonConnection)

	// Create a status line printer and defer a break.
	statusLinePrinter := &cmd.StatusLinePrinter{}
	defer statusLinePrinter.BreakIfNonEmpty()

	// Loop and print monitoring information indefinitely.
	var identifier string
	var previousStateIndex uint64
	sessionInformationPrinted := false
	for {
		// Create the list request. If there's no session specified, then we
		// need to grab all sessions and identify the most recently created one
		// for future queries.
		request := &synchronizationsvc.ListRequest{
			Selection:          selection,
			PreviousStateIndex: previousStateIndex,
		}

		// Invoke list.
		response, err := sessionService.List(context.Background(), request)
		if err != nil {
			return errors.Wrap(grpcutil.PeelAwayRPCErrorLayer(err), "list failed")
		} else if err = response.EnsureValid(); err != nil {
			return errors.Wrap(err, "invalid list response received")
		}

		// Validate the response and extract the relevant session state. If we
		// haven't already selected our target monitoring session, then we
		// choose the last session in the batch (which will be the one with the
		// most recent creation date).
		var state *synchronization.State
		previousStateIndex = response.StateIndex
		if identifier == "" {
			if len(response.SessionStates) == 0 {
				err = errors.New("no matching sessions exist")
			} else {
				state = response.SessionStates[len(response.SessionStates)-1]
				identifier = state.Session.Identifier
				selection = &selectionpkg.Selection{
					Specifications: []string{identifier},
				}
			}
		} else if len(response.SessionStates) != 1 {
			err = errors.New("invalid list response")
		} else {
			state = response.SessionStates[0]
		}
		if err != nil {
			return err
		}

		// Print session information the first time through the loop.
		if !sessionInformationPrinted {
			// Print session information.
			printSession(state, monitorConfiguration.long)

			// Print endpoint URLs, but only if not in long mode (where they're
			// already printed in the session metadata).
			if !monitorConfiguration.long {
				fmt.Println("Alpha:", state.Session.Alpha.Format("\n\t"))
				fmt.Println("Beta:", state.Session.Beta.Format("\n\t"))
			}

			// Mark session information as printed.
			sessionInformationPrinted = true
		}

		// Compute the status line.
		statusLine := computeMonitorStatusLine(state)

		// Print the status line.
		statusLinePrinter.Print(statusLine)
	}
}

var monitorCommand = &cobra.Command{
	Use:          "monitor [<session>]",
	Short:        "Show a dynamic status display for a single session",
	RunE:         monitorMain,
	SilenceUsage: true,
}

var monitorConfiguration struct {
	// help indicates whether or not to show help information and exit.
	help bool
	// long indicates whether or not to use long-format monitoring.
	long bool
	// labelSelector encodes a label selector to be used in identifying which
	// sessions should be paused.
	labelSelector string
}

func init() {
	// Grab a handle for the command line flags.
	flags := monitorCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&monitorConfiguration.help, "help", "h", false, "Show help information")

	// Wire up monitor flags.
	flags.BoolVarP(&monitorConfiguration.long, "long", "l", false, "Show detailed session information")
	flags.StringVar(&monitorConfiguration.labelSelector, "label-selector", "", "Monitor the most recently created session matching the specified label selector")
}
