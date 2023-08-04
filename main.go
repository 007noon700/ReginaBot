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

	"github.com/bwmarrin/discordgo"
)

var errInvalidFormat = errors.New("invalid format")
var attendeeRole = "1136878110646743170"

func main() {
	discord, err := discordgo.New("Bot MTEzNjgxMDAyMDc4NTM2MDk5OQ.GFsmm5.FmiDtK3qAmmVrvyvdyI-L4hYRvuzrFdkZ9hiNY")
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
		discord.ChannelMessageSend(message.ChannelID, ":dumphim:")
	case strings.HasPrefix(message.Content, "$talkshit"):
		discord.ChannelMessageSend(message.ChannelID, "post fit")
	case strings.HasPrefix(message.Content, "$skillissue"):
		discord.ChannelMessageSend(message.ChannelID, "skill issue")
	case strings.HasPrefix(message.Content, "$color"):
		createColorRole(message, discord)
	case strings.HasPrefix(message.Content, "$rsvp"):
		discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, attendeeRole)
		discord.ChannelMessageSend(message.ChannelID, "Thanks, "+message.Author.Username+" for RSVPing to WWD24! Hope to see you there!")
	}
}

func createColorRole(message *discordgo.MessageCreate, discord *discordgo.Session) {
	c := message.Content[7:14]
	fmt.Println(c)
	_, e := ParseHexColorFast(c)

	if e != nil {
		discord.ChannelMessageSend(message.ChannelID, "Sorry, I couldn't find that color. Please try again, using a hex code (starts with #)")
		fmt.Println("fail 1")
		return
	}
	cInt, convErr := strconv.ParseInt(c[1:], 16, 64)
	cIntPoi := int(cInt)
	if convErr != nil {
		discord.ChannelMessageSend(message.ChannelID, "Sorry, I had issues creating that color. Please try again or contact @synanasthesia")
		fmt.Println("fail 2")
		return
	}
	newRole := discordgo.RoleParams{Name: c, Color: &cIntPoi}
	role, er := discord.GuildRoleCreate(message.GuildID, &newRole)
	if er != nil {
		discord.ChannelMessageSend(message.ChannelID, "Sorry, I couldn't create a role with that color. Please try again or contact @synanasthesia")
		fmt.Println("fail 3")
		return
	}
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
