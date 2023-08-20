package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var errInvalidFormat = errors.New("invalid format")
var attendeeRole = "1136878110646743170"

// https://yourbasic.org/golang/time-change-convert-location-timezone/
// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func main() {
	//discord, err := discordgo.New("Bot <TOKEN>")
	discord, err := discordgo.New("Bot MTE0MjUxMTI1NDgwODgyNTk3Nw.Gy4fnF.xWD-S-4IwFUSgFRwwgkYWXzaqYVtnVbOEyz9B8")

	if err != nil {
		log.Fatal(err)
	}

	// Add event handler
	discord.AddHandler(newMessage)

	// Open session
	discord.Open()
	defer discord.Close()

	// Run until code is terminated
	fmt.Println("Bot running...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore bot messaage
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Respond to messages
	switch {
	case strings.HasPrefix(message.Content, "$dumphim"):
		// send "\:dumphim:" as yourself to get the emoji id needed for the bot to send the emoji
		discord.ChannelMessageSend(message.ChannelID, "<:dumphim:1136426428565569567>")
	case strings.HasPrefix(message.Content, "$talkshit"):
		discord.ChannelMessageSend(message.ChannelID, "post fit")
	case strings.HasPrefix(message.Content, "$skillissue"):
		discord.ChannelMessageSend(message.ChannelID, "skill issue")
	case strings.HasPrefix(message.Content, "$dogepoint"):
		discord.ChannelMessageSend(message.ChannelID, "<:dogekek:1135750882584182925> 👉")
	case message.Content == "$horse":
		discord.ChannelMessageSend(message.ChannelID, "🐴")
	case strings.HasPrefix(message.Content, "$horses"):
		discord.ChannelMessageSend(message.ChannelID, "🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴")
	case strings.HasPrefix(message.Content, "$color"):
		createColorRole(message, discord)
	case strings.HasPrefix(message.Content, "$tacobell"):
		discord.ChannelMessageSend(message.ChannelID, "https://i.imgur.com/TkZbs3J.png")
	case strings.HasPrefix(message.Content, "$wednesday"):
		// NYC supremacy
		t, _ := TimeIn(time.Now(), "America/New_York")
		weekday := t.Weekday()
		if int(weekday) == 3 {
			discord.ChannelMessageSend(message.ChannelID, "https://giphy.com/gifs/filmeditor-mean-girls-movie-3otPozZKy1ALqGLoVG")
		} else {
			discord.ChannelMessageSend(message.ChannelID, "it's not Wednesday numbnuts it is "+weekday.String())
		}

	case strings.HasPrefix(message.Content, "$rsvp"):
		discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, attendeeRole)
		discord.ChannelMessageSend(message.ChannelID, "Thanks, "+message.Author.Username+" for RSVPing to WWD24! Hope to see you there!")
	}
}

// https://github.com/bwmarrin/discordgo/wiki/FAQ
func sendImage(discord *discordgo.Session, message *discordgo.MessageCreate, imageFilename string) {
	//const attachment = new &discordgo.MessageAttachment('');
	imageFilename = "images/" + imageFilename
	f, err := os.Open(imageFilename)
	if err != nil {
		return
	}
	defer f.Close()

	ms := &discordgo.MessageSend{
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   imageFilename,
				Reader: f,
			},
		},
	}
	discord.ChannelMessageSendComplex(message.ChannelID, ms)
}

func createColorRole(message *discordgo.MessageCreate, discord *discordgo.Session) {
	c := message.Content[7:14]
	fmt.Println(c)
	_, e := ParseHexColorFast(c)

	if e != nil { //if the color is a bad hex code throw this error
		discord.ChannelMessageSend(message.ChannelID, "Sorry, I couldn't find that color. Please try again, using a hex code (starts with #)")
		fmt.Println("fail 1")
		return
	}
	//the role params struct needs a decimal INT for some reason, so we need to convert it here.
	cInt, convErr := strconv.ParseInt(c[1:], 16, 64)
	cIntPoi := int(cInt)
	if convErr != nil { //if the int parser fails for some godforsaken reason
		discord.ChannelMessageSend(message.ChannelID, "Sorry, I had issues creating that color. Please try again or contact @synanasthesia")
		fmt.Println("fail 2")
		return
	}
	newRole := discordgo.RoleParams{Name: c, Color: &cIntPoi} //this creates the role parameters - currently, the name is set to just the color string and color is obvious.
	role, er := discord.GuildRoleCreate(message.GuildID, &newRole)
	if er != nil {
		discord.ChannelMessageSend(message.ChannelID, "Sorry, I couldn't create a role with that color. Please try again or contact @synanasthesia")
		fmt.Println("fail 3")
		return
	}
	//somewhere here we either need to remove old roles or reorder the roles.
	discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, role.ID)
	discord.ChannelMessageSend(message.ChannelID, "Done! How's that?")
}

// shout out to https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func ParseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}

//token
