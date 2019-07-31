package pkger

import (
	"io"
	"strings"
)

func createFile(i *index, p string, body ...string) (*File, error) {
	pt, err := i.Parse(p)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		body = append(body, radio)
	}

	f, err := i.Create(pt)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, strings.NewReader(strings.Join(body, "\n\n")))
	if err != nil {
		return nil, err
	}

	if err := f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

const radio = `I was tuning in the shine on the late night dial
Doing anything my radio advised
With every one of those late night stations
Playing songs bringing tears to my eyes
I was seriously thinking about hiding the receiver
When the switch broke 'cause it's old
They're saying things that I can hardly believe
They really think we're getting out of control
Radio is a sound salvation
Radio is cleaning up the nation
They say you better listen to the voice of reason
But they don't give you any choice 'cause they think that it's treason
So you had better do as you are told
You better listen to the radio
I wanna bite the hand that feeds me
I wanna bite that hand so badly
I want to make them wish they'd never seen me
Some of my friends sit around every evening
And they worry about the times ahead
But everybody else is overwhelmed by indifference
And the promise of an early bed
You either shut up or get cut up; they don't wanna hear about it
It's only inches on the reel-to-reel
And the radio is in the hands of such a lot of fools
Tryin' to anesthetize the way that you feel
Radio is a sound salvation
Radio is cleaning up the nation
They say you better listen to the voice of reason
But they don't give you any choice 'cause they think that it's treason
So you had better do as you are told
You better listen to the radio
Wonderful radio
Marvelous radio
Wonderful radio
Radio, radio`
