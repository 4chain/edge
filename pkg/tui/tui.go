package tui

import (
	"context"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/echogy-io/echogy/pkg/logger"
	"github.com/echogy-io/echogy/pkg/stat"
	"github.com/gliderlabs/ssh"
	"io"
	"strings"
)

type HttpReversProxyPty struct {
	cancel context.CancelFunc
	*tea.Program
}

func (t *HttpReversProxyPty) Update() {
	t.Program.Send(tea.ShowCursor())
}

func (t *HttpReversProxyPty) Start() error {
	_, err := t.Run()
	t.cancel()
	if err != nil {
		logger.Error("run pty", err, map[string]interface{}{
			"module": "tui",
		})
		return err
	}
	return nil
}

const (
	maxWidth  = 860
	maxHeight = 480
)

// setupTerminalEnv configures the terminal environment variables for proper color support
func setupTerminalEnv(term string, environ []string) []string {
	environ = append(environ, fmt.Sprintf("TERM=%s", term))

	// Add color support environment variables
	environ = append(environ, "CLICOLOR=1")
	environ = append(environ, "CLICOLOR_FORCE=1")

	// Set color term based on terminal type
	termType := strings.ToLower(term)
	if strings.Contains(termType, "256color") {
		environ = append(environ, "COLORTERM=truecolor")
	} else if strings.Contains(termType, "color") || strings.Contains(termType, "xterm") {
		environ = append(environ, "COLORTERM=color")
	}

	return environ
}

func setupProgram(ctx context.Context, rw io.ReadWriter, term string, env []string, model tea.Model) *tea.Program {
	// Setup terminal environment
	environ := setupTerminalEnv(term, env)

	p := tea.NewProgram(
		model,
		tea.WithEnvironment(environ),
		tea.WithAltScreen(),
		tea.WithOutput(rw),
		tea.WithInput(rw),
		tea.WithContext(ctx),
	)
	return p
}

// NewHttpReverseProxyPty creates a new terminal UI instance
func NewHttpReverseProxyPty(sess ssh.Session, addr string) (*HttpReversProxyPty, error) {
	pty, windowCh, hasPty := sess.Pty()
	if !hasPty {
		return nil, errors.New("no pty")
	}

	ctx := sess.Context()

	// Setup terminal environment
	stdCtx, cancelFunc := context.WithCancel(ctx)

	queue := stat.GetQueue(ctx)

	s := stat.GetStat(ctx)

	m := newDashboard(queue, s, addr, pty.Window.Width, pty.Window.Height)

	program := setupProgram(ctx, sess, pty.Term, sess.Environ(), m)

	// Initialize dashboard
	// Start window size monitoring
	go func() {
		for {
			select {
			case <-stdCtx.Done():
				return
			case newSize := <-windowCh:
				if newSize.Height == 0 || newSize.Width == 0 {
					continue
				}
				program.Send(tea.WindowSizeMsg{
					Width:  min(newSize.Width, maxWidth),
					Height: min(newSize.Height, maxHeight),
				})
			}
		}
	}()

	return &HttpReversProxyPty{
		Program: program,
		cancel:  cancelFunc,
	}, nil
}

type RegisterPty struct {
	*tea.Program
	cancel context.CancelFunc
}

func (r *RegisterPty) Start() error {
	_, err := r.Program.Run()
	r.cancel()
	return err
}

func NewRegisterPty(sess ssh.Session, doReg doReg) (*RegisterPty, error) {
	pty, windowCh, hasPty := sess.Pty()
	if !hasPty {
		return nil, errors.New("no pty")
	}

	// Setup terminal environment
	ctx := sess.Context()

	stdCtx, cancelFunc := context.WithCancel(ctx)

	m := newRegister(doReg)

	program := setupProgram(ctx, sess, pty.Term, sess.Environ(), m)

	// Start window size monitoring
	go func() {
		for {
			select {
			case <-stdCtx.Done():
				program.Quit()
				return
			case newSize := <-windowCh:
				if newSize.Height == 0 || newSize.Width == 0 {
					continue
				}
				program.Send(tea.WindowSizeMsg{
					Width:  min(newSize.Width, maxWidth),
					Height: min(newSize.Height, maxHeight),
				})
			}
		}
	}()

	return &RegisterPty{
		Program: program,
		cancel:  cancelFunc}, nil
}
