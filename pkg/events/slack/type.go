// Copyright © 2025 Ava AI.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slack

type SlackAttachment struct {
	ID        int      `json:"id"`
	Color     string   `json:"color"`
	Fallback  string   `json:"fallback"`
	Title     string   `json:"title"`
	TitleLink string   `json:"title_link"`
	Text      string   `json:"text"`
	MrkdwnIn  []string `json:"mrkdwn_in"`
}

type SlackBlockElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackBlockSection struct {
	Type     string              `json:"type"`
	Elements []SlackBlockElement `json:"elements"`
}

type SlackBlock struct {
	Type     string              `json:"type,omitempty"`
	BlockID  string              `json:"block_id,omitempty"`
	Elements []SlackBlockSection `json:"elements,omitempty"`
}

type SlackEvent struct {
	User        string            `json:"user,omitempty"`
	BotID       string            `json:"bot_id,omitempty"`
	Subtype     string            `json:"subtype,omitempty"`
	Type        string            `json:"type"`
	TS          string            `json:"ts"`
	ClientMsgID string            `json:"client_msg_id,omitempty"`
	Text        string            `json:"text,omitempty"`
	Team        string            `json:"team,omitempty"`
	ThreadTS    string            `json:"thread_ts,omitempty"`
	Blocks      []SlackBlock      `json:"blocks,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
	Channel     string            `json:"channel"`
	EventTS     string            `json:"event_ts"`
	ChannelType string            `json:"channel_type"`
}

type SlackAuthorization struct {
	EnterpriseID        *string `json:"enterprise_id"`
	TeamID              string  `json:"team_id"`
	UserID              string  `json:"user_id"`
	IsBot               bool    `json:"is_bot"`
	IsEnterpriseInstall bool    `json:"is_enterprise_install"`
}

type ReceiveSlackEvent struct {
	Token               string               `json:"token"`
	TeamID              string               `json:"team_id"`
	ContextTeamID       string               `json:"context_team_id"`
	ContextEnterpriseID *string              `json:"context_enterprise_id"`
	Challenge           string               `json:"challenge"`
	APIAppID            string               `json:"api_app_id"`
	Event               SlackEvent           `json:"event"`
	Type                string               `json:"type"`
	EventID             string               `json:"event_id"`
	EventTime           int64                `json:"event_time"`
	Authorizations      []SlackAuthorization `json:"authorizations"`
	IsExtSharedChannel  bool                 `json:"is_ext_shared_channel"`
	EventContext        string               `json:"event_context"`
}
