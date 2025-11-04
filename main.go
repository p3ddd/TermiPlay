package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"termiplay/go-backend/game"
	"termiplay/go-backend/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	btea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

// appModel manages the state machine between lobby and games
type appModel struct {
	current tea.Model
	state   string // "lobby", "minesweeper", "game2048"
}

func newAppModel() *appModel {
	return &appModel{
		current: models.NewLobbyModel(),
		state:   "lobby",
	}
}

func (m *appModel) Init() tea.Cmd {
	if m.current != nil {
		return m.current.Init()
	}
	return nil
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit at top level
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Update current model
	updatedModel, cmd := m.current.Update(msg)

	// Check if we need to transition states
	switch m.state {
	case "lobby":
		if lobbyModel, ok := updatedModel.(*models.LobbyModel); ok {
			if lobbyModel.IsGameChosen() {
				selected := lobbyModel.GetSelected()
				if selected == -1 {
					// User quit
					return m, tea.Quit
				}
				switch selected {
				case models.Minesweeper:
					// Transition to minesweeper
					m.current = models.NewMinesweeperModel(game.Easy)
					m.state = "minesweeper"
					return m, m.current.Init()
				case models.Game2048:
					// Transition to 2048
					m.current = models.NewGame2048Model()
					m.state = "game2048"
					return m, m.current.Init()
				}
			}
		}
	case "minesweeper":
		if minesweeperModel, ok := updatedModel.(*models.MinesweeperModel); ok {
			// Check if user wants to quit
			if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "q" {
				// Return to lobby
				m.current = models.NewLobbyModel()
				m.state = "lobby"
				return m, m.current.Init()
			}
			_ = minesweeperModel // avoid unused variable
		}
	case "game2048":
		if game2048Model, ok := updatedModel.(*models.Game2048Model); ok {
			// Check if user wants to quit
			if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "q" {
				// Return to lobby
				m.current = models.NewLobbyModel()
				m.state = "lobby"
				return m, m.current.Init()
			}
			_ = game2048Model // avoid unused variable
		}
	}

	m.current = updatedModel
	return m, cmd
}

func (m *appModel) View() string {
	if m.current != nil {
		return m.current.View()
	}
	return "Loading..."
}

// teaHandler is the function that returns our Bubble Tea program.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := newAppModel()
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/termiplay_ed25519"),
		wish.WithMiddleware(
			btea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)

	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
		os.Exit(1)
	}
}
