package disc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// usuage

/*
func cliCommand(ctx cli.Context) error {
	// some code
}

func main() {
	app := cli.NewApp()
	log.Fatal()
}

use the normal usage of the cli package, just switch the import
so I need to modify cli.App
	// use the discord type as an api to use the discord server

// changes to cli.Context
	- cli.Context should wrap around the discordgo.Session

// change the cli.NewApp function to make a discord cli app

1) server boots up and connects to discord
2) wrap the needed discord api stuff into the cli.Context
3) pass args to the app.RunContext(ctx, args)
4) run like normal!



*/

// Server listens for discord messages and forwards incoming commands via the
// cli.App.RunContext()
type Server struct {
	disc *discordgo.Session
	ctx  context.Context
	name string
	Msgs chan Convo
}

// New inits a connection to the discord server with provided creds
func New(name string) (*Server, error) {
	crds, err := creds()
	if err != nil {
		return nil, errors.Wrap(err, "failure to read creds during server init")
	}
	dg, err := discordgo.New("Bot " + crds.Token)
	if err != nil {
		fmt.Println("failure to create Discord session,", err)
		return nil, err
	}
	out := Server{
		disc: dg,
	}
	return &out, nil
}

// Boot is the cli command that is added (or at least should be added)to all
// apps in order to switch the discord server on to begin forwarding parsed messages to app.RunContext
func (s *Server) Boot(ctx *cli.Context) error {
	// listen for ctrl + c
	mngr := NewManager(contetxt.Background)
	s.ctx = mngr.Ctx
	go mngr.Listen()
	s.disc.AddHandler(s.mainHandler)
	// Listen for discord Messages
	go s.Listen()
	mngr.WG.Wait()
	return nil
}

func (s *Server) mainHandler(ss *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore everything that doesn't include the command
	if !strings.Contains(m.Content, fmt.Sprintf("!%s", s.name)) {
		return
	}

}

//////////////////////////////////////////////////
//  	Conversations
///////////////////////////////////////////////
/* Conversions allow commands to wait for user input

msg in -> msg
write a response back
out <- msg
and wait for another response
 in -> msg

either that or make a specific pipe reader and writer?

*/

// Convo allows commands to wait for user input
type Convo struct {
	id    string
	User  string
	Write chan string
	Read  chan string
}

func (c *Convo) Write(msg string) {

}

// type Phrase struct {
// 	ChanID string

// }

//////////////////////////////////////////////////
//  fetching credentials
///////////////////////////////////////////////
type cred struct {
	Permissions int    `json:"BOT_PERMISSIONS"`
	Token       string `json:"TOKEN"`
	ClientID    string `json:"CLIENT_ID"`
}

// TODO: change this
func creds() (*cred, error) {
	var out cred
	// Ask to unlock credentials

	jsonFile, err := ioutil.ReadFile("/home/evan/.creds/discord.json")
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
func (s *Server) Listen(ctx context.Context) {
	s.disc.Open()
	<-ctx.Done()
	s.disc.Close()
}
