package synchronization

import (
	"github.com/pkg/errors"
)

// createStreamPrompter implements Prompter on top of a
// Synchronization_CreateServer stream.
type createStreamPrompter struct {
	// stream is the underlying Synchronization_CreateServer stream.
	stream Synchronization_CreateServer
}

// sendReceive performs a send/receive cycle by sending a CreateResponse and
// receiving a CreateRequest.
func (p *createStreamPrompter) sendReceive(request *CreateResponse) (*CreateRequest, error) {
	// Send the request.
	if err := p.stream.Send(request); err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	// Receive the response.
	if response, err := p.stream.Recv(); err != nil {
		return nil, errors.Wrap(err, "unable to receive response")
	} else if err = response.ensureValid(false); err != nil {
		return nil, errors.Wrap(err, "invalid response received")
	} else {
		return response, nil
	}
}

// Message implements the Message method of Prompter.
func (p *createStreamPrompter) Message(message string) error {
	_, err := p.sendReceive(&CreateResponse{Message: message})
	return err
}

// Prompt implements the Prompt method of Prompter.
func (p *createStreamPrompter) Prompt(prompt string) (string, error) {
	if response, err := p.sendReceive(&CreateResponse{Prompt: prompt}); err != nil {
		return "", err
	} else {
		return response.Response, nil
	}
}

// flushStreamPrompter implements Prompter on top of a
// Synchronization_FlushServer stream.
type flushStreamPrompter struct {
	// stream is the underlying Synchronization_FlushServer stream.
	stream Synchronization_FlushServer
}

// sendReceive performs a send/receive cycle by sending a FlushResponse and
// receiving a FlushRequest.
func (p *flushStreamPrompter) sendReceive(request *FlushResponse) (*FlushRequest, error) {
	// Send the request.
	if err := p.stream.Send(request); err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	// Receive the response.
	if response, err := p.stream.Recv(); err != nil {
		return nil, errors.Wrap(err, "unable to receive response")
	} else if err = response.ensureValid(false); err != nil {
		return nil, errors.Wrap(err, "invalid response received")
	} else {
		return response, nil
	}
}

// Message implements the Message method of Prompter.
func (p *flushStreamPrompter) Message(message string) error {
	_, err := p.sendReceive(&FlushResponse{Message: message})
	return err
}

// Prompt implements the Prompt method of Prompter.
func (p *flushStreamPrompter) Prompt(_ string) (string, error) {
	return "", errors.New("prompting not supported on flush message streams")
}

// pauseStreamPrompter implements Prompter on top of a
// Synchronization_PauseServer stream.
type pauseStreamPrompter struct {
	// stream is the underlying Synchronization_PauseServer stream.
	stream Synchronization_PauseServer
}

// sendReceive performs a send/receive cycle by sending a PauseResponse and
// receiving a PauseRequest.
func (p *pauseStreamPrompter) sendReceive(request *PauseResponse) (*PauseRequest, error) {
	// Send the request.
	if err := p.stream.Send(request); err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	// Receive the response.
	if response, err := p.stream.Recv(); err != nil {
		return nil, errors.Wrap(err, "unable to receive response")
	} else if err = response.ensureValid(false); err != nil {
		return nil, errors.Wrap(err, "invalid response received")
	} else {
		return response, nil
	}
}

// Message implements the Message method of Prompter.
func (p *pauseStreamPrompter) Message(message string) error {
	_, err := p.sendReceive(&PauseResponse{Message: message})
	return err
}

// Prompt implements the Prompt method of Prompter.
func (p *pauseStreamPrompter) Prompt(_ string) (string, error) {
	return "", errors.New("prompting not supported on pause message streams")
}

// resetStreamPrompter implements Prompter on top of a
// Synchronization_ResetServer stream.
type resetStreamPrompter struct {
	// stream is the underlying Synchronization_ResetServer stream.
	stream Synchronization_ResetServer
}

// sendReceive performs a send/receive cycle by sending a ResetResponse and
// receiving a ResetRequest.
func (p *resetStreamPrompter) sendReceive(request *ResetResponse) (*ResetRequest, error) {
	// Send the request.
	if err := p.stream.Send(request); err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	// Receive the response.
	if response, err := p.stream.Recv(); err != nil {
		return nil, errors.Wrap(err, "unable to receive response")
	} else if err = response.ensureValid(false); err != nil {
		return nil, errors.Wrap(err, "invalid response received")
	} else {
		return response, nil
	}
}

// Message implements the Message method of Prompter.
func (p *resetStreamPrompter) Message(message string) error {
	_, err := p.sendReceive(&ResetResponse{Message: message})
	return err
}

// Prompt implements the Prompt method of Prompter.
func (p *resetStreamPrompter) Prompt(prompt string) (string, error) {
	if response, err := p.sendReceive(&ResetResponse{Prompt: prompt}); err != nil {
		return "", err
	} else {
		return response.Response, nil
	}
}

// resumeStreamPrompter implements Prompter on top of a
// Synchronization_ResumeServer stream.
type resumeStreamPrompter struct {
	// stream is the underlying Synchronization_ResumeServer stream.
	stream Synchronization_ResumeServer
}

// sendReceive performs a send/receive cycle by sending a ResumeResponse and
// receiving a ResumeRequest.
func (p *resumeStreamPrompter) sendReceive(request *ResumeResponse) (*ResumeRequest, error) {
	// Send the request.
	if err := p.stream.Send(request); err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	// Receive the response.
	if response, err := p.stream.Recv(); err != nil {
		return nil, errors.Wrap(err, "unable to receive response")
	} else if err = response.ensureValid(false); err != nil {
		return nil, errors.Wrap(err, "invalid response received")
	} else {
		return response, nil
	}
}

// Message implements the Message method of Prompter.
func (p *resumeStreamPrompter) Message(message string) error {
	_, err := p.sendReceive(&ResumeResponse{Message: message})
	return err
}

// Prompt implements the Prompt method of Prompter.
func (p *resumeStreamPrompter) Prompt(prompt string) (string, error) {
	if response, err := p.sendReceive(&ResumeResponse{Prompt: prompt}); err != nil {
		return "", err
	} else {
		return response.Response, nil
	}
}

// terminateStreamPrompter implements Prompter on top of a
// Synchronization_TerminateServer stream.
type terminateStreamPrompter struct {
	// stream is the underlying Synchronization_TerminateServer stream.
	stream Synchronization_TerminateServer
}

// sendReceive performs a send/receive cycle by sending a TerminateResponse and
// receiving a TerminateRequest.
func (p *terminateStreamPrompter) sendReceive(request *TerminateResponse) (*TerminateRequest, error) {
	// Send the request.
	if err := p.stream.Send(request); err != nil {
		return nil, errors.Wrap(err, "unable to send request")
	}

	// Receive the response.
	if response, err := p.stream.Recv(); err != nil {
		return nil, errors.Wrap(err, "unable to receive response")
	} else if err = response.ensureValid(false); err != nil {
		return nil, errors.Wrap(err, "invalid response received")
	} else {
		return response, nil
	}
}

// Message implements the Message method of Prompter.
func (p *terminateStreamPrompter) Message(message string) error {
	_, err := p.sendReceive(&TerminateResponse{Message: message})
	return err
}

// Prompt implements the Prompt method of Prompter.
func (p *terminateStreamPrompter) Prompt(_ string) (string, error) {
	return "", errors.New("prompting not supported on terminate message streams")
}
