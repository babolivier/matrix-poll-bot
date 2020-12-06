package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/matrix-org/gomatrix"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// emojiRegexpStr is the regexp for detecting emojis (as a string). Copied from https://github.com/urakozz/go-emoji
const emojiRegexpStr = "[\\x{2712}\\x{2714}\\x{2716}\\x{271d}\\x{2721}\\x{2728}\\x{2733}\\x{2734}\\x{2744}\\x{2747}\\x{274c}\\x{274e}\\x{2753}-\\x{2755}\\x{2757}\\x{2763}\\x{2764}\\x{2795}-\\x{2797}\\x{27a1}\\x{27b0}\\x{27bf}\\x{2934}\\x{2935}\\x{2b05}-\\x{2b07}\\x{2b1b}\\x{2b1c}\\x{2b50}\\x{2b55}\\x{3030}\\x{303d}\\x{1f004}\\x{1f0cf}\\x{1f170}\\x{1f171}\\x{1f17e}\\x{1f17f}\\x{1f18e}\\x{1f191}-\\x{1f19a}\\x{1f201}\\x{1f202}\\x{1f21a}\\x{1f22f}\\x{1f232}-\\x{1f23a}\\x{1f250}\\x{1f251}\\x{1f300}-\\x{1f321}\\x{1f324}-\\x{1f393}\\x{1f396}\\x{1f397}\\x{1f399}-\\x{1f39b}\\x{1f39e}-\\x{1f3f0}\\x{1f3f3}-\\x{1f3f5}\\x{1f3f7}-\\x{1f4fd}\\x{1f4ff}-\\x{1f53d}\\x{1f549}-\\x{1f54e}\\x{1f550}-\\x{1f567}\\x{1f56f}\\x{1f570}\\x{1f573}-\\x{1f579}\\x{1f587}\\x{1f58a}-\\x{1f58d}\\x{1f590}\\x{1f595}\\x{1f596}\\x{1f5a5}\\x{1f5a8}\\x{1f5b1}\\x{1f5b2}\\x{1f5bc}\\x{1f5c2}-\\x{1f5c4}\\x{1f5d1}-\\x{1f5d3}\\x{1f5dc}-\\x{1f5de}\\x{1f5e1}\\x{1f5e3}\\x{1f5ef}\\x{1f5f3}\\x{1f5fa}-\\x{1f64f}\\x{1f680}-\\x{1f6c5}\\x{1f6cb}-\\x{1f6d0}\\x{1f6e0}-\\x{1f6e5}\\x{1f6e9}\\x{1f6eb}\\x{1f6ec}\\x{1f6f0}\\x{1f6f3}\\x{1f910}-\\x{1f918}\\x{1f980}-\\x{1f984}\\x{1f9c0}\\x{3297}\\x{3299}\\x{a9}\\x{ae}\\x{203c}\\x{2049}\\x{2122}\\x{2139}\\x{2194}-\\x{2199}\\x{21a9}\\x{21aa}\\x{231a}\\x{231b}\\x{2328}\\x{2388}\\x{23cf}\\x{23e9}-\\x{23f3}\\x{23f8}-\\x{23fa}\\x{24c2}\\x{25aa}\\x{25ab}\\x{25b6}\\x{25c0}\\x{25fb}-\\x{25fe}\\x{2600}-\\x{2604}\\x{260e}\\x{2611}\\x{2614}\\x{2615}\\x{2618}\\x{261d}\\x{2620}\\x{2622}\\x{2623}\\x{2626}\\x{262a}\\x{262e}\\x{262f}\\x{2638}-\\x{263a}\\x{2648}-\\x{2653}\\x{2660}\\x{2663}\\x{2665}\\x{2666}\\x{2668}\\x{267b}\\x{267f}\\x{2692}-\\x{2694}\\x{2696}\\x{2697}\\x{2699}\\x{269b}\\x{269c}\\x{26a0}\\x{26a1}\\x{26aa}\\x{26ab}\\x{26b0}\\x{26b1}\\x{26bd}\\x{26be}\\x{26c4}\\x{26c5}\\x{26c8}\\x{26ce}\\x{26cf}\\x{26d1}\\x{26d3}\\x{26d4}\\x{26e9}\\x{26ea}\\x{26f0}-\\x{26f5}\\x{26f7}-\\x{26fa}\\x{26fd}\\x{2702}\\x{2705}\\x{2708}-\\x{270d}\\x{270f}]|\\x{23}\\x{20e3}|\\x{2a}\\x{20e3}|\\x{30}\\x{20e3}|\\x{31}\\x{20e3}|\\x{32}\\x{20e3}|\\x{33}\\x{20e3}|\\x{34}\\x{20e3}|\\x{35}\\x{20e3}|\\x{36}\\x{20e3}|\\x{37}\\x{20e3}|\\x{38}\\x{20e3}|\\x{39}\\x{20e3}|\\x{1f1e6}[\\x{1f1e8}-\\x{1f1ec}\\x{1f1ee}\\x{1f1f1}\\x{1f1f2}\\x{1f1f4}\\x{1f1f6}-\\x{1f1fa}\\x{1f1fc}\\x{1f1fd}\\x{1f1ff}]|\\x{1f1e7}[\\x{1f1e6}\\x{1f1e7}\\x{1f1e9}-\\x{1f1ef}\\x{1f1f1}-\\x{1f1f4}\\x{1f1f6}-\\x{1f1f9}\\x{1f1fb}\\x{1f1fc}\\x{1f1fe}\\x{1f1ff}]|\\x{1f1e8}[\\x{1f1e6}\\x{1f1e8}\\x{1f1e9}\\x{1f1eb}-\\x{1f1ee}\\x{1f1f0}-\\x{1f1f5}\\x{1f1f7}\\x{1f1fa}-\\x{1f1ff}]|\\x{1f1e9}[\\x{1f1ea}\\x{1f1ec}\\x{1f1ef}\\x{1f1f0}\\x{1f1f2}\\x{1f1f4}\\x{1f1ff}]|\\x{1f1ea}[\\x{1f1e6}\\x{1f1e8}\\x{1f1ea}\\x{1f1ec}\\x{1f1ed}\\x{1f1f7}-\\x{1f1fa}]|\\x{1f1eb}[\\x{1f1ee}-\\x{1f1f0}\\x{1f1f2}\\x{1f1f4}\\x{1f1f7}]|\\x{1f1ec}[\\x{1f1e6}\\x{1f1e7}\\x{1f1e9}-\\x{1f1ee}\\x{1f1f1}-\\x{1f1f3}\\x{1f1f5}-\\x{1f1fa}\\x{1f1fc}\\x{1f1fe}]|\\x{1f1ed}[\\x{1f1f0}\\x{1f1f2}\\x{1f1f3}\\x{1f1f7}\\x{1f1f9}\\x{1f1fa}]|\\x{1f1ee}[\\x{1f1e8}-\\x{1f1ea}\\x{1f1f1}-\\x{1f1f4}\\x{1f1f6}-\\x{1f1f9}]|\\x{1f1ef}[\\x{1f1ea}\\x{1f1f2}\\x{1f1f4}\\x{1f1f5}]|\\x{1f1f0}[\\x{1f1ea}\\x{1f1ec}-\\x{1f1ee}\\x{1f1f2}\\x{1f1f3}\\x{1f1f5}\\x{1f1f7}\\x{1f1fc}\\x{1f1fe}\\x{1f1ff}]|\\x{1f1f1}[\\x{1f1e6}-\\x{1f1e8}\\x{1f1ee}\\x{1f1f0}\\x{1f1f7}-\\x{1f1fb}\\x{1f1fe}]|\\x{1f1f2}[\\x{1f1e6}\\x{1f1e8}-\\x{1f1ed}\\x{1f1f0}-\\x{1f1ff}]|\\x{1f1f3}[\\x{1f1e6}\\x{1f1e8}\\x{1f1ea}-\\x{1f1ec}\\x{1f1ee}\\x{1f1f1}\\x{1f1f4}\\x{1f1f5}\\x{1f1f7}\\x{1f1fa}\\x{1f1ff}]|\\x{1f1f4}\\x{1f1f2}|\\x{1f1f5}[\\x{1f1e6}\\x{1f1ea}-\\x{1f1ed}\\x{1f1f0}-\\x{1f1f3}\\x{1f1f7}-\\x{1f1f9}\\x{1f1fc}\\x{1f1fe}]|\\x{1f1f6}\\x{1f1e6}|\\x{1f1f7}[\\x{1f1ea}\\x{1f1f4}\\x{1f1f8}\\x{1f1fa}\\x{1f1fc}]|\\x{1f1f8}[\\x{1f1e6}-\\x{1f1ea}\\x{1f1ec}-\\x{1f1f4}\\x{1f1f7}-\\x{1f1f9}\\x{1f1fb}\\x{1f1fd}-\\x{1f1ff}]|\\x{1f1f9}[\\x{1f1e6}\\x{1f1e8}\\x{1f1e9}\\x{1f1eb}-\\x{1f1ed}\\x{1f1ef}-\\x{1f1f4}\\x{1f1f7}\\x{1f1f9}\\x{1f1fb}\\x{1f1fc}\\x{1f1ff}]|\\x{1f1fa}[\\x{1f1e6}\\x{1f1ec}\\x{1f1f2}\\x{1f1f8}\\x{1f1fe}\\x{1f1ff}]|\\x{1f1fb}[\\x{1f1e6}\\x{1f1e8}\\x{1f1ea}\\x{1f1ec}\\x{1f1ee}\\x{1f1f3}\\x{1f1fa}]|\\x{1f1fc}[\\x{1f1eb}\\x{1f1f8}]|\\x{1f1fd}\\x{1f1f0}|\\x{1f1fe}[\\x{1f1ea}\\x{1f1f9}]|\\x{1f1ff}[\\x{1f1e6}\\x{1f1f2}\\x{1f1fc}]"

// config is the configuration structure.
type config struct {
	Matrix struct {
		AccessToken string `yaml:"access_token"`
		UserID      string `yaml:"user_id"`
		HSURL       string `yaml:"hs_url"`
		SkipFilter  bool   `yaml:"skip_filter"`
	}
}

// reaction is the content of m.reaction events.
type reaction struct {
	RelatesTo struct {
		RelType string `json:"rel_type"`
		Key     string `json:"key"`
		EventID string `json:"event_id"`
	} `json:"m.relates_to"`
}

// handler defines the handlers to call when processing incoming Matrix events, along with some util functions.
type handler struct {
	Client               *gomatrix.Client
	EmojiRegexp          *regexp.Regexp
	StartingSpacesRegexp *regexp.Regexp
}

func newHandler(homeserverURL, userID, accessToken string) (h *handler, err error) {
	h = new(handler)

	h.Client, err = gomatrix.NewClient(homeserverURL, userID, accessToken)
	if err != nil {
		return
	}

	h.EmojiRegexp, err = regexp.Compile(emojiRegexpStr)
	if err != nil {
		return
	}

	h.StartingSpacesRegexp, err = regexp.Compile("^ +")
	if err != nil {
		return
	}

	return
}

// setupFilter configures the /sync filter on the Matrix homeserver.
// Returns the filter ID.
func (h *handler) setupFilter() string {
	filter := gomatrix.Filter{
		Room: gomatrix.RoomFilter{
			Timeline: gomatrix.FilterPart{
				Types: []string{
					"m.room.message",
					"m.room.member",
				},
			},
		},
		EventFields: []string{
			"type",               // Needed to register the handlers.
			"event_id",           // Needed for logging.
			"room_id",            // Needed to manage the rooms we're in.
			"state_key",          // Needed for the syncer to manage the room's state.
			"sender",             // Needed to mention the poll's author.
			"content.body",       // Needed to process messages.
			"content.membership", // Needed to process invites.
		},
	}

	filterJSON, err := json.Marshal(filter)
	if err != nil {
		panic(err)
	}

	resp, err := h.Client.CreateFilter(filterJSON)
	if err != nil {
		panic(err)
	}

	return resp.FilterID
}

// handleMessage handles incoming m.room.messages events.
func (h *handler) handleMessage(event *gomatrix.Event) {
	logger := logrus.WithFields(logrus.Fields{
		"event_id": event.ID,
		"room_id":  event.RoomID,
	})

	body := event.Content["body"].(string)

	// Extract the question and choices from the event's content's body.
	question, choices, errMsg := h.parseMessage(body)
	if len(errMsg) > 0 {
		// Send the error message as a notice to the room.
		if _, err := h.Client.SendNotice(event.RoomID, errMsg); err != nil {
			logger.WithField("errMsg", errMsg).Errorf("Failed to send error message: %s", err.Error())
		}
		return
	}

	if len(question) == 0 {
		// This message is not a poll.
		return
	}

	// Generate the HTML string to send as a notice.
	notice := h.generateNoticeHTML(event.Sender, question, choices)

	// Send the poll to the room.
	res, err := h.Client.SendMessageEvent(
		event.RoomID, "m.room.message", gomatrix.GetHTMLMessage("m.notice", notice),
	)
	if err != nil {
		logger.Errorf("Couldn't send poll to room: %s", err.Error())
		return
	}

	// Add reactions to the poll message.
	for k := range choices {
		content := reaction{}
		content.RelatesTo.RelType = "m.annotation"
		content.RelatesTo.Key = k
		content.RelatesTo.EventID = res.EventID

		if _, err := h.Client.SendMessageEvent(event.RoomID, "m.reaction", &content); err != nil {
			logger.WithField("poll_event_id", res.EventID).Errorf(
				"Couldn't send reaction to poll: %s", err.Error(),
			)
			return
		}
	}
}

// parseMessage extracts the question and choices from the given string.
// Returns with the question and choices, and an error string to send to the room if the message isn't formatted
// correctly.
func (h *handler) parseMessage(body string) (question string, choices map[string]string, errMsg string) {
	r := strings.NewReader(body)

	scanner := bufio.NewScanner(r)

	choices = make(map[string]string)
	firstLine := true
	choiceErrMsg := "Each choice needs to be on a different line, and to start with an emoji."
	// Read the message line by line.
	for scanner.Scan() {
		line := scanner.Text()

		if firstLine {
			// If the line starts with "!poll", it's the question, otherwise it's not a message for us.
			if strings.HasPrefix(line, "!poll ") {
				question = strings.Replace(line, "!poll ", "", 1)
			} else {
				// Don't send an error message since the message isn't for us.
				return
			}
		} else {
			// handle empty lines and lines containing only spaces.
			// TODO: also remove tabs.
			if len(h.trimStartingSpaces(line)) == 0 {
				continue
			}

			// This is a choice in the poll.
			indexes := h.EmojiRegexp.FindStringIndex(line)
			if len(indexes) > 0 {
				// The line must start with an emoji.
				if indexes[0] != 0 {
					errMsg = choiceErrMsg
					return
				}

				key := line[:indexes[1]]
				// Allow each emoji only once.
				if _, exists := choices[key]; exists {
					errMsg = fmt.Sprintf("Two occurrences of emoji %s, please use each emoji only once.", key)
				}

				// Save the choice.
				value := line[indexes[1]:]
				choices[key] = h.trimStartingSpaces(value)
			} else {
				// Each choice line must contain an emoji.
				errMsg = choiceErrMsg
				return
			}
		}

		firstLine = false
	}

	return
}

// trimStartingSpaces removes the spaces at the beginning of a string.
func (h *handler) trimStartingSpaces(s string) string {
	r := regexp.MustCompile("^ +")
	indexes := r.FindStringIndex(s)
	if len(indexes) == 0 {
		return s
	}
	return s[indexes[1]:]
}

// generateNoticeHTML generates the HTML text to send as a notice to the room.
func (h *handler) generateNoticeHTML(userID string, question string, choices map[string]string) string {
	displayName := userID

	// Try to retrieve the sender's display name from the homeserver. If we can't, use the user ID instead.
	res, err := h.Client.GetDisplayName(userID)
	if err == nil {
		displayName = res.DisplayName
	}

	format := `
	<b>New poll from <a href="https://matrix.to/#/%s">%s</a>!</b><br><br>
	Question: <i><u>%s</u></i><br>
	Cast your vote by clicking on the reactions under this message (if available).<br><br>
	Choices:<br>
	`

	notice := fmt.Sprintf(format, userID, displayName, question)

	for key, value := range choices {
		notice += fmt.Sprintf("%s: %s<br>", key, value)
	}

	return notice
}

// handleMembership handles m.room.member events, and autojoins rooms it's invited in.
func (h *handler) handleMembership(event *gomatrix.Event) {
	var joinDelay time.Duration = 1

	if membership, ok := event.Content["membership"]; !ok || membership != "invite" {
		return
	}

	logrus.Infof("Trying to join room %s I was invited to", event.RoomID)

	time.Sleep(joinDelay * time.Second)
	_, err := h.Client.JoinRoom(event.RoomID, "", struct{}{})
	if err != nil {
		logrus.Errorf("Failed to join room %s: %s", event.RoomID, err.Error())
	}

	logrus.Infof("Successfully joined room %s", event.RoomID)
}

func main() {
	var cfg config

	configFile := flag.String("config", "config.yaml", "Path to the configuration file")

	flag.Parse()

	configBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(errors.Wrap(err, "Couldn't open the configuration file"))
	}

	if err := yaml.Unmarshal(configBytes, &cfg); err != nil {
		panic(errors.Wrap(err, "Couldn't read the configuration file"))
	}

	h, err := newHandler(cfg.Matrix.HSURL, cfg.Matrix.UserID, cfg.Matrix.AccessToken)
	if err != nil {
		panic(errors.Wrap(err, "Couldn't initialise the Matrix client"))
	}

	if !cfg.Matrix.SkipFilter {
		filterID := h.setupFilter()

		h.Client.Store.SaveFilterID(cfg.Matrix.UserID, filterID)
	}

	syncer := h.Client.Syncer.(*gomatrix.DefaultSyncer)

	syncer.OnEventType("m.room.message", h.handleMessage)
	syncer.OnEventType("m.room.member", h.handleMembership)

	logrus.Info("Syncing...")
	if err := h.Client.Sync(); err != nil {
		panic(errors.Wrap(err, "Sync returned with an error"))
	}
}
