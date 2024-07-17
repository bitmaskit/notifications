package channel

type Channel string

func (c Channel) String() string {
	return string(c)
}

const (
	SMS   Channel = "sms"
	Email Channel = "email"
	Slack Channel = "slack"
)
