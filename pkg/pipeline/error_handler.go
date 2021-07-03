package pipeline

import (
	"fmt"
	"github.com/unchain/pipeline/pkg/actions/smtp_action"
	"github.com/unchain/pipeline/pkg/domain"
)

func (p *Pipeline) handleError(trigger domain.Trigger, tag string, err error) {
	
	err = trigger.Respond(tag, map[string]interface{}{
		"error": err,
	}, err)
	if err != nil {
		p.log.Errorf("Could not handle error, msg: %v", err)
	}
	// send alert email
	_, err = smtp_action.Invoke(p.log, map[string]interface{}{
		"username":   p.cfg.Actions.SmtpAction.Username,
		"password":   p.cfg.Actions.SmtpAction.Password,
		"hostname":   p.cfg.Actions.SmtpAction.Hostname,
		"port":       p.cfg.Actions.SmtpAction.Port,
		"from":       p.cfg.Actions.SmtpAction.From,
		"recipients": p.cfg.Actions.SmtpAction.Recipients,
		"message":    []byte(fmt.Sprintf(`
Organization: %s;
Error message: %v
`, p.cfg.Organization, err)),
	})
	if err != nil {
		p.log.Errorf("Could not handle error, msg: %v", err)
	}

	
	
}
