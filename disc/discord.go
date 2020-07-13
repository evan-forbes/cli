package disc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Server listens for discord messages and forwards incoming commands via the
// cli.App.RunContext()
type Server struct {
	disc    *discordgo.Session
	ctx     context.Context
	name    string
	waiting map[string]*Slug
	config  Config
	Sink    chan Slug
}

// New inits a connection to the discord server with provided creds
func New(ctx context.Context, name string, path string) (*Server, error) {
	crds, err := configs(path)
	if err != nil {
		return nil, errors.Wrap(err, "failure to read config during server init")
	}
	dg, err := discordgo.New("Bot " + crds.Token)
	if err != nil {
		fmt.Println("failure to create Discord session,", err)
		return nil, err
	}
	out := Server{
		disc: dg,
		ctx:  ctx,
		Sink: make(chan Slug, 10),
		name: name,
	}
	out.disc.AddHandler(out.mainHandler)
	return &out, nil
}

// mainHandler filters messages from discord passing all qualifying messages as
// command line input into an instance of the app
func (s *Server) mainHandler(ss *discordgo.Session, m *discordgo.MessageCreate) {
	// discord app name
	name := fmt.Sprintf("!%s", s.name)

	// ignore self posts
	if m.Author.ID == ss.State.User.ID {
		return
	}

	// check is the server is waiting on a response from the user
	waitingSlug, has := s.waiting[fmt.Sprintf("%s%s", m.ChannelID, m.Author.Username)]
	if has {
		// forward response to the slug
		waitingSlug.response <- m.Content
		return
	}

	// ignore everything else that doesn't include the command
	if !strings.Contains(m.Content, name) {
		return
	}

	// do a quick parse for args
	index := strings.Index(m.Content, name)
	args := strings.Split(m.Content[index+len(name):], " ")

	// create a new slug
	slug := s.NewSlug(m, args)

	// pass the args to app.RunContext
	s.Sink <- *slug
}

func parseArgs(input string) []string {
	return strings.Split(input)
}

//////////////////////////////////////////////////
//  	Conversations
///////////////////////////////////////////////

// Slug contains data pertaining to a conversation with a user and fullfills the
// io.Reader and io.Writer interfaces
type Slug struct {
	context.Context
	ChanID   string
	User     string
	Args     []string
	srv      *Server
	response chan string
}

// NewSlug issues a *Slug
func (s *Server) NewSlug(m *discordgo.MessageCreate, args []string) *Slug {
	id := fmt.Sprintf("%s%s", m.Author.Username, m.ChannelID)
	return &Slug{Context: s.ctx, ChanID: id, User: m.Author.Username, srv: s, Args: args}
}

// ID combines the user's name and channel id to identify a conversion
func (s *Slug) ID() string {
	return fmt.Sprintf("%s%s", s.ChanID, s.User)
}

// Write fullfills the io.Writer interface, writing the provided []byte to the
// discord channel of origin
func (s *Slug) Write(in []byte) (int, error) {
	_, err := s.srv.disc.ChannelMessageSend(s.ChanID, string(in))
	return len(in), err
}

// Read fullfills the io.Reader interface, which asks the discord channel of
// origin for input
func (s *Slug) Read(out []byte) (int, error) {
	// notify the server
	s.srv.waiting[s.ID()] = s
	select {
	case resp := <-s.response:
		out = []byte(resp)
		return len(out), nil
	case <-time.After(time.Minute):
		s.Write([]byte("no input detected, aborting"))
		return 0, errors.New("user did not respond within 1 minute, aborting")
	}
}

// func (s *Slug) WritePNG() {}

//////////////////////////////////////////////////
//  fetching credentials
///////////////////////////////////////////////
type config struct {
	Permissions int    `json:"BOT_PERMISSIONS"`
	Token       string `json:"TOKEN"`
	ClientID    string `json:"CLIENT_ID"`
}

func configs(path string) (*config, error) {
	var out config
	// Ask to unlock credentials

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonFile, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Listen opens websocket streaming from discord
func (s *Server) Listen(mngr *Manager) {
	defer mngr.WG.Done()
	s.disc.Open()
	<-mngr.Ctx.Done()
	close(s.Sink)
	s.disc.Close()
}
