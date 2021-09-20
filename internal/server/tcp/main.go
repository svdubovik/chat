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
	conn net.Conn
	user models.User
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
	defer s.logger.Info().Msg("TCP Server stopped")

	for {
		s.logger.Debug().Msg("Accept conections...")
		c, err := ln.Accept()
		if err != nil {
			s.logger.Error().Err(err).Msg("Cannot accept new conection")
			continue
		}

		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(nc net.Conn) {
	remoteAddr := nc.RemoteAddr().String()
	s.logger.Info().Str("RemoteAddr", remoteAddr).Msgf("Connection is accepted from: %s", remoteAddr)
	defer func() {
		nc.Close()
		s.logger.Info().Str("RemoteAddr", remoteAddr).Msgf("Connection from %s is closed", remoteAddr)
	}()

	client, err := s.login(nc)
	if err != nil {
		fmt.Fprintf(nc, "You are not authorized")
		s.logger.Error().Err(err)
		return
	}

	s.clients = append(s.clients, client)
	s.fanMsg(fmt.Sprintf("*** %s join to chat ***", client.user.Login))
	s.directMsg(client, fmt.Sprintf("> Hello, %s\n", client.user.Login))

	for {
		str, err := bufio.NewReader(client.conn).ReadString('\n')
		if err != nil {
			s.logger.Error().Err(err)
			break
		}
		s.fanMsg(fmt.Sprintf("%s> %s", client.user.Login, str))
	}

	s.leave(client)
}

func (s *Server) directMsg(client *Client, msg string) {
	fmt.Fprintf(client.conn, "%s\n", msg)
}

func (s *Server) fanMsg(msg string) {
	for _, client := range s.clients {
		s.directMsg(client, msg)
	}
}

func (s *Server) leave(client *Client) {
	s.fanMsg(fmt.Sprintf("*** %s leave the chat ***", client.user.Login))
}

func (s *Server) login(nc net.Conn) (*Client, error) {
	s.logger.Debug().Msg("Try to login")

	fmt.Fprintf(nc, "> Type your username: \n")
	str, err := bufio.NewReader(nc).ReadString('\n')
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn: nc,
	}
	client.user.Login = strings.TrimSpace(str)

	s.logger.Debug().Msgf("%s is authorized", client.user.Login)
	return client, nil
}
