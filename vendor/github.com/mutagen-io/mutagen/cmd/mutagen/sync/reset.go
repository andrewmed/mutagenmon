package sync

import (
	"context"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/mutagen-io/mutagen/cmd"
	"github.com/mutagen-io/mutagen/cmd/mutagen/daemon"
	"github.com/mutagen-io/mutagen/pkg/grpcutil"
	"github.com/mutagen-io/mutagen/pkg/prompting"
	"github.com/mutagen-io/mutagen/pkg/selection"
	synchronizationsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
)

// ResetWithLabelSelector is an orchestration convenience method invokes the
// reset command using the specified label selector.
func ResetWithLabelSelector(labelSelector string) error {
	resetConfiguration.labelSelector = labelSelector
	return resetMain(nil, nil)
}

func resetMain(command *cobra.Command, arguments []string) error {
	// Create session selection specification.
	selection := &selection.Selection{
		All:            resetConfiguration.all,
		Specifications: arguments,
		LabelSelector:  resetConfiguration.labelSelector,
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

	// Invoke the session reset method. The stream will close when the
	// associated context is cancelled.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := sessionService.Reset(ctx)
	if err != nil {
		return errors.Wrap(grpcutil.PeelAwayRPCErrorLayer(err), "unable to invoke reset")
	}

	// Send the initial request.
	request := &synchronizationsvc.ResetRequest{
		Selection: selection,
	}
	if err := stream.Send(request); err != nil {
		return errors.Wrap(grpcutil.PeelAwayRPCErrorLayer(err), "unable to send reset request")
	}

	// Create a status line printer.
	statusLinePrinter := &cmd.StatusLinePrinter{}

	// Receive and process responses until we're done.
	for {
		if response, err := stream.Recv(); err != nil {
			statusLinePrinter.BreakIfNonEmpty()
			return errors.Wrap(grpcutil.PeelAwayRPCErrorLayer(err), "reset failed")
		} else if err = response.EnsureValid(); err != nil {
			statusLinePrinter.BreakIfNonEmpty()
			return errors.Wrap(err, "invalid reset response received")
		} else if response.Message == "" && response.Prompt == "" {
			statusLinePrinter.Clear()
			return nil
		} else if response.Message != "" {
			statusLinePrinter.Print(response.Message)
			if err := stream.Send(&synchronizationsvc.ResetRequest{}); err != nil {
				statusLinePrinter.BreakIfNonEmpty()
				return errors.Wrap(grpcutil.PeelAwayRPCErrorLayer(err), "unable to send message response")
			}
		} else if response.Prompt != "" {
			statusLinePrinter.BreakIfNonEmpty()
			if response, err := prompting.PromptCommandLine(response.Prompt); err != nil {
				return errors.Wrap(err, "unable to perform prompting")
			} else if err = stream.Send(&synchronizationsvc.ResetRequest{Response: response}); err != nil {
				return errors.Wrap(grpcutil.PeelAwayRPCErrorLayer(err), "unable to send prompt response")
			}
		}
	}
}

var resetCommand = &cobra.Command{
	Use:          "reset [<session>...]",
	Short:        "Reset synchronization session history",
	RunE:         resetMain,
	SilenceUsage: true,
}

var resetConfiguration struct {
	// help indicates whether or not to show help information and exit.
	help bool
	// all indicates whether or not all sessions should be reset.
	all bool
	// labelSelector encodes a label selector to be used in identifying which
	// sessions should be paused.
	labelSelector string
}

func init() {
	// Grab a handle for the command line flags.
	flags := resetCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&resetConfiguration.help, "help", "h", false, "Show help information")

	// Wire up reset flags.
	flags.BoolVarP(&resetConfiguration.all, "all", "a", false, "Reset all sessions")
	flags.StringVar(&resetConfiguration.labelSelector, "label-selector", "", "Reset sessions matching the specified label selector")
}
