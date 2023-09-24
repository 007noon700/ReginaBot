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
	discord, err := discordgo.New("Bot <TOKEN>")

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

func helpMessage() (s string) {
	var message = `
$color <hex code>: I'll set your color to <hex code> (example: #FFFFFF).
$help: Make me repeat this message for some ungodly reason
$rsvp: RSVP to WWD24!

$date: Get today's date
$dogepoint: <:dogekek:1135750882584182925> 👉
$dumphim: <:dumphim:1136426428565569567>
$horse: 🐴
$horses: Several horses
$tacobell: I'll express my feelings about Taco Bell
$talkshit: I'll say "post fit"
$skillissue: I'll say "skill issue"
$uck: Sometimes people get up to this I guess
$wednesday: If it's Wednesday I'll let you know
$white: <:white:1136047355997720747>
`
	return message
}

func wednesdayMessage() (s string) {
	// NYC supremacy
	t, _ := TimeIn(time.Now(), "America/New_York")
	weekday := t.Weekday()
	if int(weekday) == 3 {
		return "https://giphy.com/gifs/filmeditor-mean-girls-movie-3otPozZKy1ALqGLoVG"
	} else {
		return "It's not Wednesday numbnuts it is " + weekday.String()
	}
}

func dateMessage() (s string) {
	t, _ := TimeIn(time.Now(), "America/New_York")
	dateString := t.Format("January 2")
	if dateString != "October 3" {
		return dateString
	}
	// It's October 3rd.
	return "https://tenor.com/view/crush-diary-october3rd-mean-girls-lindsay-lohan-gif-9906172"
}

func checkForStrings(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// OK we don't want to spam so let's kind of sort these by funniest highest to lowest.
	// and send only one per message (even if has two.)
	if strings.Contains(message.Content, "fetch") {
		// Stop trying to make fetch happen
		discord.ChannelMessageSend(message.ChannelID,
			"https://tenor.com/view/fetch-mean-girls-gif-19691105")
	} else if strings.Contains(message.Content, "how much") || strings.Contains(message.Content, "how many") {
		// The limit does not exist
		discord.ChannelMessageSend(message.ChannelID,
			"https://tenor.com/view/mean-girls-karen-gif-9300840")
	} else if strings.Contains(message.Content, "=") || strings.Contains(message.Content, "+") {
		// You can't join mathletes. It's social suicide.
		discord.ChannelMessageSend(message.ChannelID,
			"https://y.yarn.co/ea1dd776-80ed-43fb-a53d-9e77520bf781_text.gif")
	}
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	fmt.Println(message.Content)

	// Ignore bot messaage
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Don't bother with all the below switch logic if there isn't a $ in the front.
	// Conversely, don't check for strings if there's maybe a command.
	if message.Content[0:1] != "$" {
		checkForStrings(discord, message)
		return
	}

	tokens := strings.Fields(message.Content)

	// Respond to messages
	switch tokens[0] {
	// Somewhat more involved commands
	case "$color":
		createColorRole(message, discord, tokens)
	case "$help":
		discord.ChannelMessageSend(message.ChannelID, helpMessage())
	case "$rsvp":
		discord.GuildMemberRoleAdd(message.GuildID, message.Author.ID, attendeeRole)
		discord.ChannelMessageSend(message.ChannelID,
			"Thanks, "+message.Author.Username+" for RSVPing to WWD24! Hope to see you there!")

	// Memes
	case "$date":
		discord.ChannelMessageSend(message.ChannelID, dateMessage())
	case "$dogepoint":
		discord.ChannelMessageSend(message.ChannelID, "<:dogekek:1135750882584182925> 👉")
	case "$dumphim":
		discord.ChannelMessageSend(message.ChannelID, "<:dumphim:1136426428565569567>")
	case "$horse":
		discord.ChannelMessageSend(message.ChannelID, "🐴")
	case "$horses":
		discord.ChannelMessageSend(message.ChannelID,
			"🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴🐴")
	case "$tacobell":
		discord.ChannelMessageSend(message.ChannelID, "https://i.imgur.com/TkZbs3J.png")
	case "$talkshit":
		discord.ChannelMessageSend(message.ChannelID, "post fit")
	case "$skillissue":
		discord.ChannelMessageSend(message.ChannelID, "skill issue")
	case "$uck":
		discord.ChannelMessageSend(message.ChannelID,
			"https://cdn.discordapp.com/attachments/1136094668136915066/1143333955555303568/image0.gif")
	case "$wednesday":
		discord.ChannelMessageSend(message.ChannelID, wednesdayMessage())
	case "$waluigi":
		waluigi(message, discord, tokens)
	case "$white":
		discord.ChannelMessageSend(message.ChannelID, "<:white:1136047355997720747>")
	default:
		// There wasn't a command so let's just check for strings.
		checkForStrings(discord, message)
		return
	}
}

// so far all this does is print the profile picture of either the message author
// (if no one mentioned) or the user mentioned. I can't figure out how to invert stuff.
func waluigi(message *discordgo.MessageCreate, discord *discordgo.Session, tokens []string) {
	var url = ""
	if len(tokens) == 1 {
		url = message.Author.AvatarURL("")
	} else {
		// convert e.g. <@425544164365565962> to 425544164365565962
		userId := tokens[1][2 : len(tokens[1])-1]
		mentionedUser, err := discord.User(userId)
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID,
				"Literally who are you even talking about?")
			return
		}
		url = mentionedUser.AvatarURL("")
	}
	discord.ChannelMessageSend(message.ChannelID, "Ok, just imagine this is in inverted colors.")
	discord.ChannelMessageSend(message.ChannelID, url)
}

func createColorRole(message *discordgo.MessageCreate, discord *discordgo.Session, tokens []string) {
	if len(tokens) == 1 {
		discord.ChannelMessageSend(message.ChannelID,
			`Are you fucking with me? You need to type "$color <hex code>" (e.g. "$color #FFFFFF")`)
		return
	} else if !strings.Contains(tokens[1], "#") {
		discord.ChannelMessageSend(message.ChannelID,
			`Are you fucking with me? You need to type "$color <hex code>" (e.g. "$color #FFFFFF"). YOU NEED THE # IN THERE TOO.`)
		return
	}

	c := tokens[1]
	fmt.Println(c)
	_, e := ParseHexColorFast(c)

	if e != nil { //if the color is a bad hex code throw this error
		discord.ChannelMessageSend(message.ChannelID,
			`I couldn't find that color. Try again using a hex code (e.g. "$color #FFFFFF")`)
		fmt.Println("fail 1")
		return
	}
	//the role params struct needs a decimal INT for some reason, so we need to convert it here.
	cInt, convErr := strconv.ParseInt(c[1:], 16, 64)
	cIntPoi := int(cInt)
	if convErr != nil { //if the int parser fails for some godforsaken reason
		discord.ChannelMessageSend(message.ChannelID,
			"I had issues creating that color. Is it even real?? Try again or contact @synanasthesia")
		fmt.Println("fail 2")
		return
	}
	newRole := discordgo.RoleParams{Name: c, Color: &cIntPoi} //this creates the role parameters - currently, the name is set to just the color string and color is obvious.
	role, er := discord.GuildRoleCreate(message.GuildID, &newRole)
	if er != nil {
		discord.ChannelMessageSend(message.ChannelID,
			"Sorry, I couldn't create a role with that color. Please try again or contact @synanasthesia")
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
