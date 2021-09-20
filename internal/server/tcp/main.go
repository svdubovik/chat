package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/rs/zerolog"
	"svdubovik.com/chat/internal/config"
	"svdubovik.com/chat/internal/models"
)

type Client struct {
	conn       net.Conn
	chOut      chan string
	chIn       chan string
	authorized bool
	user       models.User
}

type Server struct {
	cfg     *config.Config
	logger  *zerolog.Logger
	clients []*Client
}

func NewServer(cfg *config.Config, logger *zerolog.Logger) *Server {
	return &Server{
		cfg:     cfg,
		logger:  logger,
		clients: make([]*Client, 0),
	}
}

func (s *Server) Run() (string, error) {
	ln, err := net.Listen("tcp", s.cfg.BindAddress)
	if err != nil {
		return fmt.Sprintf("Cannot start on %s", s.cfg.BindAddress), err
	}
	defer ln.Close()

	s.logger.Info().Msgf("TCP Server started on %s", s.cfg.BindAddress)
	defer s.logger.Info().Msg("Server stopped")

	for {
		s.logger.Debug().Msg("Starting accept conections")
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Error().Err(err).Msg("Cannot accept new conection")
			continue
		}

		client := &Client{
			conn:  conn,
			chOut: make(chan string),
			chIn:  make(chan string),
		}
		s.clients = append(s.clients, client)

		go s.handleConnection(client)
	}
}

func (s *Server) handleConnection(client *Client) {
	s.logger.Debug().Msg(fmt.Sprintf("Connection is accepted; RemoteAddr: %s", client.conn.RemoteAddr().String()))
	defer func() {
		client.conn.Close()
		s.logger.Debug().Msg("Connection is closed")
	}()

	s.login(client)

	if client.authorized {
		fmt.Fprintf(client.conn, "Hello, %s\n", client.user.Login)
	} else {
		fmt.Fprintf(client.conn, "You are not authorized")
		return
	}

	go func() {
		for {
			fmt.Fprintf(client.conn, "# ")
			str, err := bufio.NewReader(client.conn).ReadString('\n')
			if err != nil {
				s.logger.Error().Err(err).Msg("Cannot read network message")
				close(client.chOut)
				break
			}

			client.chOut <- str
		}
	}()

	go func() {
		for msg := range client.chIn {
			fmt.Fprint(client.conn, msg)
		}
	}()

	for msg := range client.chOut {
		for _, c := range s.clients {
			if c.conn.RemoteAddr().String() != client.conn.RemoteAddr().String() {
				c.chIn <- fmt.Sprintf("%s > %s", client.user.Login, msg)
			}
		}
	}
}

func (s *Server) login(client *Client) {
	s.logger.Debug().Msg("Try to login")

	fmt.Fprintf(client.conn, "Type your username: \n")
	str, err := bufio.NewReader(client.conn).ReadString('\n')
	if err != nil {
		s.logger.Error().Err(err).Msg("Cannot read network message")
		return
	}

	// Authentification
	client.user.Login = strings.TrimSpace(str)

	// Authorization
	client.authorized = true

	s.logger.Debug().Msgf("%s is authorized", client.user.Login)
}
