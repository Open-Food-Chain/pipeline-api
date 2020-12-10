package pipeline

import (
	"fmt"
	"github.com/jmoiron/jsonq"
	"github.com/unchain/pipeline/pkg/actions/fileparser_action"
	"github.com/unchain/pipeline/pkg/actions/http_action"
	"github.com/unchain/pipeline/pkg/actions/templater_action"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchain/pipeline/pkg/triggers/api_trigger"
	"github.com/unchainio/pkg/errors"
)

/*
	The start function is blocking.
*/
func (p *Pipeline) Start() error {
	// Initialize trigger
	trigger := &api_trigger.Trigger{}
	err := trigger.Init(p.log, []byte(p.cfg.Trigger.Config))
	if err != nil {
		return errors.Wrap(err, "could not init trigger")
	}
	p.log.Debugf("Initialized pipeline trigger")

	p.start(trigger)

	return nil
}

func (p *Pipeline) start(trigger domain.Trigger) {
	// start infinite loop to process messages
	for {
		select {
		case <-p.stopChannel:
			return
		default:
			p.handleNextMessage(trigger)
		}
	}
}

func (p *Pipeline) handleNextMessage(trigger domain.Trigger) {
	tag, message, err := trigger.NextMessage()
	if err != nil {
		p.handleError(trigger, tag, err)
		return
	}
	p.log.Debugf("Next message with tag %v\n", tag)
	body, exists := message["body"].(string)
	if !exists {
		p.log.Errorf("could not get body")
	}
	// parse JSON
	fileparserOutput, err := fileparser_action.Invoke(p.log, map[string]interface{}{
		fileparser_action.FileType:  p.cfg.Actions.FileparserAction.Filetype,
		fileparser_action.Header:    p.cfg.Actions.FileparserAction.Header,
		fileparser_action.Delimiter: p.cfg.Actions.FileparserAction.Delimiter,
		fileparser_action.File:      []byte(body),
	})
	if err != nil {
		p.handleError(trigger, tag, err)
		return
	}
	// cast fileparser output to array of map[string]interface{}
	records, ok := fileparserOutput["messages"].([]map[string]interface{})
	if !ok {
		p.handleError(trigger, tag, errors.Errorf("could not cast fileparser messages output to array"))
		return
	}
	err = p.handleRecords(records)
	if err != nil {
		p.handleError(trigger, tag, err)
		return
	}

	// call respond to finish processing
	err = trigger.Respond(tag, nil, err)
	if err != nil {
		p.handleError(trigger, tag, err)
		return
	}
}

// loop for record handling
func (p *Pipeline) handleRecords(records []map[string]interface{}) error {
	for index, record := range records {
		// data transformation
		inputVariables := GetInputVariables(jsonq.NewQuery(record), p.cfg.Actions.TemplaterAction.Variables)
		templaterOutput, err := templater_action.Invoke(p.log, map[string]interface{}{
			templater_action.InputTemplate:  p.cfg.Actions.TemplaterAction.Template,
			templater_action.InputVariables: inputVariables,
		})
		if err != nil {
			return errors.Wrapf(err, "could not transform data for record with index %v", index)
		}

		// call import-api
		httpOutput, err := http_action.Invoke(p.log, map[string]interface{}{
			http_action.RequestBody: []byte(fmt.Sprintf("%s", templaterOutput[templater_action.TemplateResult])),
			http_action.Url:         p.cfg.Actions.HttpAction.Url,
			http_action.Method:      p.cfg.Actions.HttpAction.Method,
		})
		if err != nil {
			return errors.Wrapf(err, "could not call import-api for record with ID %v \n HTTP response: %v", index, httpOutput)
		}
	}
	return nil
}
