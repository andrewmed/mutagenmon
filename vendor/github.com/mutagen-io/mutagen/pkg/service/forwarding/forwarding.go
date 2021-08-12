package forwarding

import (
	"github.com/pkg/errors"

	"github.com/mutagen-io/mutagen/pkg/selection"
	"github.com/mutagen-io/mutagen/pkg/url"
)

// ensureValid verifies that a CreationSpecification is valid.
func (s *CreationSpecification) ensureValid() error {
	// A nil creation specification is not valid.
	if s == nil {
		return errors.New("nil creation specification")
	}

	// Verify that the source URL is valid and is a forwarding URL.
	if err := s.Source.EnsureValid(); err != nil {
		return errors.Wrap(err, "invalid source URL")
	} else if s.Source.Kind != url.Kind_Forwarding {
		return errors.New("source URL is not a forwarding URL")
	}

	// Verify that the destination URL is valid and is a forwarding URL.
	if err := s.Destination.EnsureValid(); err != nil {
		return errors.Wrap(err, "invalid destination URL")
	} else if s.Destination.Kind != url.Kind_Forwarding {
		return errors.New("destination URL is not a forwarding URL")
	}

	// Verify that the configuration is valid.
	if err := s.Configuration.EnsureValid(false); err != nil {
		return errors.Wrap(err, "invalid session configuration")
	}

	// Verify that the source-specific configuration is valid.
	if err := s.ConfigurationSource.EnsureValid(true); err != nil {
		return errors.Wrap(err, "invalid source-specific configuration")
	}

	// Verify that the destination-specific configuration is valid.
	if err := s.ConfigurationDestination.EnsureValid(true); err != nil {
		return errors.Wrap(err, "invalid destination-specific configuration")
	}

	// Verify that the name is valid.
	if err := selection.EnsureNameValid(s.Name); err != nil {
		return errors.Wrap(err, "invalid name")
	}

	// Verify that labels are valid.
	for k, v := range s.Labels {
		if err := selection.EnsureLabelKeyValid(k); err != nil {
			return errors.Wrap(err, "invalid label key")
		} else if err = selection.EnsureLabelValueValid(v); err != nil {
			return errors.Wrap(err, "invalid label value")
		}
	}

	// There's no need to validate the Paused field - either value is valid.

	// Success.
	return nil
}

// ensureValid verifies that a CreateRequest is valid.
func (r *CreateRequest) ensureValid(first bool) error {
	// A nil create request is not valid.
	if r == nil {
		return errors.New("nil create request")
	}

	// Handle validation based on whether or not this is the first request in
	// the stream.
	if first {
		// Verify that the creation specification is valid.
		if err := r.Specification.ensureValid(); err != nil {
			return err
		}

		// Verify that the response field is empty.
		if r.Response != "" {
			return errors.New("non-empty prompt response")
		}
	} else {
		// Verify that the creation specification is nil.
		if r.Specification != nil {
			return errors.New("creation specification present")
		}

		// We can't really validate the response field, and an empty value may
		// be appropriate. It's up to the process performing the prompting to
		// decide.
	}

	// Success.
	return nil
}

// EnsureValid verifies that a CreateResponse is valid.
func (r *CreateResponse) EnsureValid() error {
	// A nil create response is not valid.
	if r == nil {
		return errors.New("nil create response")
	}

	// Count the number of fields that are set.
	var fieldsSet uint
	if r.Session != "" {
		fieldsSet++
	}
	if r.Message != "" {
		fieldsSet++
	}
	if r.Prompt != "" {
		fieldsSet++
	}

	// Enforce that exactly one field is set.
	if fieldsSet != 1 {
		return errors.New("incorrect number of fields set")
	}

	// Success.
	return nil
}

// ensureValid verifies that a ListRequest is valid.
func (r *ListRequest) ensureValid() error {
	// A nil list request is not valid.
	if r == nil {
		return errors.New("nil list request")
	}

	// Validate the session specification.
	if err := r.Selection.EnsureValid(); err != nil {
		return errors.Wrap(err, "invalid session specification")
	}

	// There's no need to validate the state index - any value is valid.

	// Success.
	return nil
}

// ensureValid verifies that a ListResponse is valid.
func (r *ListResponse) EnsureValid() error {
	// A nil list response is not valid.
	if r == nil {
		return errors.New("nil list response")
	}

	// Ensure that all states are valid.
	for _, s := range r.SessionStates {
		if err := s.EnsureValid(); err != nil {
			return errors.Wrap(err, "invalid session state")
		}
	}

	// Success.
	return nil
}

// ensureValid verifies that a PauseRequest is valid.
func (r *PauseRequest) ensureValid(first bool) error {
	// A nil pause request is not valid.
	if r == nil {
		return errors.New("nil pause request")
	}

	// Handle validation based on whether or not this is the first request in
	// the stream.
	if first {
		// Validate the session selection specification.
		if err := r.Selection.EnsureValid(); err != nil {
			return errors.Wrap(err, "invalid session selection specification")
		}
	} else {
		// Ensure that no session selection specification is present when
		// acknowledging messages.
		if r.Selection != nil {
			return errors.New("non-nil session selection specification on message acknowledgement")
		}
	}

	// Success.
	return nil
}

// EnsureValid verifies that a PauseResponse is valid.
func (r *PauseResponse) EnsureValid() error {
	// A nil pause response is not valid.
	if r == nil {
		return errors.New("nil pause response")
	}

	// We can't really verify the message field. Even an empty value may be
	// valid.

	// Success.
	return nil
}

// ensureValid verifies that a ResumeRequest is valid.
func (r *ResumeRequest) ensureValid(first bool) error {
	// A nil resume request is not valid.
	if r == nil {
		return errors.New("nil resume request")
	}

	// Handle validation based on whether or not this is the first request in
	// the stream.
	if first {
		// Validate the session selection specification.
		if err := r.Selection.EnsureValid(); err != nil {
			return errors.Wrap(err, "invalid session selection specification")
		}

		// Verify that the response field is empty.
		if r.Response != "" {
			return errors.New("non-empty prompt response")
		}
	} else {
		// Ensure that no session selection specification is present when
		// acknowledging messages.
		if r.Selection != nil {
			return errors.New("non-nil session selection specification on message acknowledgement")
		}

		// We can't really validate the response field, and an empty value may
		// be appropriate. It's up to the process performing the prompting to
		// decide.
	}

	// Success.
	return nil
}

// EnsureValid verifies that a ResumeResponse is valid.
func (r *ResumeResponse) EnsureValid() error {
	// A nil resume response is not valid.
	if r == nil {
		return errors.New("nil resume response")
	}

	// Count the number of fields that are set.
	var fieldsSet uint
	if r.Message != "" {
		fieldsSet++
	}
	if r.Prompt != "" {
		fieldsSet++
	}

	// Enforce that at most a single field is set. Unlike CreateResponse, we
	// allow neither to be set, which indicates completion. In CreateResponse,
	// this completion is indicated by the session identifier being set.
	if fieldsSet > 1 {
		return errors.New("multiple fields set")
	}

	// Success.
	return nil
}

// ensureValid verifies that a TerminateRequest is valid.
func (r *TerminateRequest) ensureValid(first bool) error {
	// A nil terminate request is not valid.
	if r == nil {
		return errors.New("nil terminate request")
	}

	// Handle validation based on whether or not this is the first request in
	// the stream.
	if first {
		// Validate the session selection specification.
		if err := r.Selection.EnsureValid(); err != nil {
			return errors.Wrap(err, "invalid session selection specification")
		}
	} else {
		// Ensure that no session selection specification is present when
		// acknowledging messages.
		if r.Selection != nil {
			return errors.New("non-nil session selection specification on message acknowledgement")
		}

		// We can't really validate the response field, and an empty value may
		// be appropriate, especially if this is just a message acknowledgement.
	}

	// Success.
	return nil
}

// EnsureValid verifies that a TerminateResponse is valid.
func (r *TerminateResponse) EnsureValid() error {
	// A nil terminate response is not valid.
	if r == nil {
		return errors.New("nil terminate response")
	}

	// We can't really verify the message field. Even an empty value may be
	// valid.

	// Success.
	return nil
}
