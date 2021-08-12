package logging

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

// writer is an io.Writer that splits its input stream into lines and writes
// those lines to an underlying logger.
type writer struct {
	// callback is the logging callback.
	callback func(string)
	// buffer is any incomplete line fragment left over from a previous write.
	buffer []byte
}

// trimCarriageReturn trims any single trailing carriage return from the end of
// a byte slice.
func trimCarriageReturn(buffer []byte) []byte {
	if len(buffer) > 0 && buffer[len(buffer)-1] == '\r' {
		return buffer[:len(buffer)-1]
	}
	return buffer
}

// Write implements io.Writer.Write.
func (w *writer) Write(buffer []byte) (int, error) {
	// Append the data to our internal buffer.
	w.buffer = append(w.buffer, buffer...)

	// Process all lines in the buffer, tracking the number of bytes that we
	// process.
	var processed int
	remaining := w.buffer
	for {
		// Find the index of the next newline character.
		index := bytes.IndexByte(remaining, '\n')
		if index == -1 {
			break
		}

		// Process the line.
		w.callback(string(trimCarriageReturn(remaining[:index])))

		// Update the number of bytes that we've processed.
		processed += index + 1

		// Update the remaining slice.
		remaining = remaining[index+1:]
	}

	// If we managed to process bytes, then truncate our internal buffer.
	if processed > 0 {
		// Compute the number of leftover bytes.
		leftover := len(w.buffer) - processed

		// If there are leftover bytes, then shift them to the front of the
		// buffer.
		if leftover > 0 {
			copy(w.buffer[:leftover], w.buffer[processed:])
		}

		// Truncate the buffer.
		w.buffer = w.buffer[:leftover]
	}

	// Done.
	return len(buffer), nil
}

// Logger is the main logger type. It has the novel property that it still
// functions if nil, but it doesn't log anything. It is designed to use the
// standard logger provided by the log package, so it respects any flags set for
// that logger. It is safe for concurrent usage.
type Logger struct {
	// prefix is any prefix specified for the logger.
	prefix string
}

// RootLogger is the root logger from which all other loggers derive.
var RootLogger = &Logger{}

// Sublogger creates a new sublogger with the specified name.
func (l *Logger) Sublogger(name string) *Logger {
	// If the logger is nil, then the sublogger will be as well.
	if l == nil {
		return nil
	}

	// Compute the new prefix.
	prefix := name
	if l.prefix != "" {
		prefix = l.prefix + "." + name
	}

	// Create the new logger.
	return &Logger{
		prefix: prefix,
	}
}

// output is the shared internal logging method.
func (l *Logger) output(level, line string) {
	// Compute the formatted line.
	if l.prefix != "" {
		line = fmt.Sprintf("[%s|%s] %s", l.prefix, level, line)
	} else {
		line = fmt.Sprintf("[%s] %s", level, line)
	}

	// Log.
	log.Output(4, line)
}

// println provides logging with formatting semantics equivalent to fmt.Println.
func (l *Logger) println(level Level, v ...interface{}) {
	if l != nil && currentLevel >= level {
		l.output(level.String(), fmt.Sprintln(v...))
	}
}

// printf provides logging with formatting semantics equivalent to fmt.Printf.
func (l *Logger) printf(level Level, format string, v ...interface{}) {
	if l != nil && currentLevel >= level {
		l.output(level.String(), fmt.Sprintf(format, v...))
	}
}

// Error logs errors with formatting semantics equivalent to fmt.Println.
func (l *Logger) Error(v ...interface{}) {
	l.println(LevelError, v...)
}

// Errorf logs errors with formatting semantics equivalent to fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.printf(LevelError, format, v...)
}

// Warning logs warnings with formatting semantics equivalent to fmt.Println.
func (l *Logger) Warning(v ...interface{}) {
	l.println(LevelWarning, v...)
}

// Warningf logs warnings with formatting semantics equivalent to fmt.Printf.
func (l *Logger) Warningf(format string, v ...interface{}) {
	l.printf(LevelWarning, format, v...)
}

// Info logs information with formatting semantics equivalent to fmt.Println.
func (l *Logger) Info(v ...interface{}) {
	l.println(LevelInfo, v...)
}

// Infof logs information with formatting semantics equivalent to fmt.Printf.
func (l *Logger) Infof(format string, v ...interface{}) {
	l.printf(LevelInfo, format, v...)
}

// Debug logs debug information with formatting semantics equivalent to
// fmt.Println.
func (l *Logger) Debug(v ...interface{}) {
	l.println(LevelDebug, v...)
}

// Debugf logs debug information with formatting semantics equivalent to
// fmt.Printf.
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.printf(LevelDebug, format, v...)
}

// Trace logs tracing information with formatting semantics equivalent to
// fmt.Println.
func (l *Logger) Trace(v ...interface{}) {
	l.println(LevelTrace, v...)
}

// Tracef logs tracing information with formatting semantics equivalent to
// fmt.Printf.
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.printf(LevelTrace, format, v...)
}

// Writer returns an io.Writer that logs output lines using the specified level.
func (l *Logger) Writer(level Level) io.Writer {
	// If the logger is nil or the current logging level is set lower than the
	// requested level, then we can just discard input since it won't be logged
	// anyway. This saves us the overhead of scanning lines.
	if l == nil || currentLevel < level {
		return ioutil.Discard
	}

	// Create the writer.
	return &writer{
		callback: func(s string) {
			l.println(level, s)
		},
	}
}
