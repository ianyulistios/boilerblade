package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

// ConsumerGen generates a general-purpose RabbitMQ consumer from name and title only.
// No entity/fields or usecase/DTO coupling.
type ConsumerGen struct {
	Name        string // user-provided name (e.g. OrderEvents, payment)
	Title       string // human-readable title (e.g. "Order Events"); used in comments/logs
	Identifier  string // snake_case for files and const values (e.g. order_events)
	StructName  string // PascalCase for type names (e.g. OrderEvents) -> XConsumer
	ConstPrefix string // PascalCase for constant names (e.g. OrderEvents)
}

// NewConsumerGen creates a consumer generator from the name and optional title.
// Name is normalized to identifier (snake_case) and struct/const prefix (PascalCase).
// Title defaults to StructName if empty.
func NewConsumerGen(name, title string) *ConsumerGen {
	name = strings.TrimSpace(name)
	title = strings.TrimSpace(title)
	identifier := nameToIdentifier(name)
	structName := nameToStructName(identifier)
	if title == "" {
		title = structName
	}
	return &ConsumerGen{
		Name:        name,
		Title:       title,
		Identifier:  identifier,
		StructName:  structName,
		ConstPrefix: structName,
	}
}

// nameToIdentifier converts "OrderEvents" or "order_events" to "order_events".
func nameToIdentifier(s string) string {
	if s == "" {
		return "consumer"
	}
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' {
			b.WriteByte('_')
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	// Collapse multiple underscores
	result := b.String()
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	result = strings.Trim(result, "_")
	if result == "" {
		return "consumer"
	}
	return result
}

// nameToStructName converts "order_events" to "OrderEvents".
func nameToStructName(identifier string) string {
	parts := strings.Split(identifier, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}
	return strings.Join(parts, "")
}

// GenerateAMQPConstants writes constants/amqp_<identifier>.go.
func (c *ConsumerGen) GenerateAMQPConstants() error {
	tmpl := `package constants

// {{.ConstPrefix}} Exchange and Queues (generated — {{.Title}})
const (
	{{.ConstPrefix}}ExchangeName           = "{{.ExchangeName}}"
	{{.ConstPrefix}}CreatedQueueName       = "{{.CreatedQueueName}}"
	{{.ConstPrefix}}UpdatedQueueName       = "{{.UpdatedQueueName}}"
	{{.ConstPrefix}}CreatedRouteKey        = "{{.CreatedRouteKey}}"
	{{.ConstPrefix}}UpdatedRouteKey        = "{{.UpdatedRouteKey}}"
	{{.ConstPrefix}}CreatedQueueInterval   = {{.QueueInterval}}
	{{.ConstPrefix}}UpdatedQueueInterval   = {{.QueueInterval}}
)
`
	data := map[string]interface{}{
		"ConstPrefix":        c.ConstPrefix,
		"Title":              c.Title,
		"ExchangeName":       c.Identifier + "_events",
		"CreatedQueueName":   c.Identifier + "_created_queue",
		"UpdatedQueueName":   c.Identifier + "_updated_queue",
		"CreatedRouteKey":    c.Identifier + ".created",
		"UpdatedRouteKey":    c.Identifier + ".updated",
		"QueueInterval":      3000,
	}
	return c.writeFile("constants/amqp_"+c.Identifier+".go", tmpl, data)
}

// GenerateConsumer writes src/consumer/<identifier>.go (generic handler, no usecase/DTO).
func (c *ConsumerGen) GenerateConsumer() error {
	tmpl := `package consumer

import (
	"boilerblade/config/amqp"
	"boilerblade/constants"
	"boilerblade/helper"
	"encoding/json"
	"fmt"

	amqplib "github.com/streadway/amqp"
)

// {{.StructName}}Consumer handles AMQP messages for {{.Title}}.
// Add your business logic in handleCreatedMessage and handleUpdatedMessage.
type {{.StructName}}Consumer struct {
	subConnection amqp.IAMQPChannel
}

// New{{.StructName}}Consumer creates a new {{.Title}} consumer instance.
func New{{.StructName}}Consumer(amqpConn amqp.IAMQPConnection) (*{{.StructName}}Consumer, error) {
	channel, err := amqpConn.Channel()
	if err != nil {
		helper.LogError("Failed to get AMQP channel for {{.StructName}}Consumer", err, "", map[string]interface{}{
			"source": "New{{.StructName}}Consumer",
		})
		return nil, err
	}

	exchangeName := constants.{{.ConstPrefix}}ExchangeName
	if err := channel.DeclareExchange(exchangeName, "direct"); err != nil {
		helper.LogError("Failed to declare exchange in {{.StructName}}Consumer", err, exchangeName, map[string]interface{}{
			"source": "New{{.StructName}}Consumer",
		})
		channel.Close()
		return nil, err
	}

	helper.LogInfo("{{.StructName}}Consumer initialized", map[string]interface{}{
		"source":        "New{{.StructName}}Consumer",
		"exchange":      exchangeName,
		"exchange_type": "direct",
	})

	return &{{.StructName}}Consumer{
		subConnection: channel,
	}, nil
}

// Close closes the AMQP channel.
func (c *{{.StructName}}Consumer) Close() error {
	if c.subConnection != nil {
		return c.subConnection.Close()
	}
	return nil
}

// ProcessMessage starts all {{.Title}} consumers.
func (c *{{.StructName}}Consumer) ProcessMessage() {
	go c.ProcessCreated()
	go c.ProcessUpdated()
}

// ProcessCreated consumes {{.Identifier}}.created messages.
func (c *{{.StructName}}Consumer) ProcessCreated() {
	var messages <-chan amqplib.Delivery

	que, err := c.subConnection.NewQueue(
		constants.{{.ConstPrefix}}ExchangeName,
		constants.{{.ConstPrefix}}CreatedQueueName,
		constants.QueueType,
		constants.{{.ConstPrefix}}CreatedRouteKey,
		constants.{{.ConstPrefix}}CreatedQueueInterval,
	)
	if err != nil {
		helper.LogError("({{.StructName}}Created) Error when declare queue", err, constants.{{.ConstPrefix}}CreatedQueueName, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if messages, err = c.subConnection.ReadMessage(que); err != nil {
		fmt.Println(err)
		return
	}

	helper.LogInfo("Started consuming {{.Identifier}}.created messages", map[string]interface{}{
		"source": "{{.StructName}}Consumer.ProcessCreated",
		"queue":  constants.{{.ConstPrefix}}CreatedQueueName,
	})

	for msg := range messages {
		if err := c.handleCreatedMessage(msg); err != nil {
			helper.LogError("Failed to process {{.Identifier}}.created message", err, "", map[string]interface{}{
				"source":     "{{.StructName}}Consumer.ProcessCreated",
				"message_id": msg.MessageId,
			})
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
	}
}

// ProcessUpdated consumes {{.Identifier}}.updated messages.
func (c *{{.StructName}}Consumer) ProcessUpdated() {
	var messages <-chan amqplib.Delivery

	que, err := c.subConnection.NewQueue(
		constants.{{.ConstPrefix}}ExchangeName,
		constants.{{.ConstPrefix}}UpdatedQueueName,
		constants.QueueType,
		constants.{{.ConstPrefix}}UpdatedRouteKey,
		constants.{{.ConstPrefix}}UpdatedQueueInterval,
	)
	if err != nil {
		helper.LogError("({{.StructName}}Updated) Error when declare queue", err, constants.{{.ConstPrefix}}UpdatedQueueName, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if messages, err = c.subConnection.ReadMessage(que); err != nil {
		fmt.Println(err)
		return
	}

	helper.LogInfo("Started consuming {{.Identifier}}.updated messages", map[string]interface{}{
		"source": "{{.StructName}}Consumer.ProcessUpdated",
		"queue":  constants.{{.ConstPrefix}}UpdatedQueueName,
	})

	for msg := range messages {
		if err := c.handleUpdatedMessage(msg); err != nil {
			helper.LogError("Failed to process {{.Identifier}}.updated message", err, "", map[string]interface{}{
				"source":     "{{.StructName}}Consumer.ProcessUpdated",
				"message_id": msg.MessageId,
			})
			msg.Nack(false, true)
			continue
		}
		msg.Ack(false)
	}
}

// handleCreatedMessage processes a single .created message. Add your logic here.
func (c *{{.StructName}}Consumer) handleCreatedMessage(msg amqplib.Delivery) error {
	helper.LogInfo("Processing {{.Title}} created message", map[string]interface{}{
		"source":       "{{.StructName}}Consumer.handleCreatedMessage",
		"message_id":   msg.MessageId,
		"routing_key":  msg.RoutingKey,
	})

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		helper.LogError("Failed to unmarshal message", err, "", map[string]interface{}{
			"source":     "{{.StructName}}Consumer.handleCreatedMessage",
			"message_id": msg.MessageId,
			"body":       string(msg.Body),
		})
		return err
	}

	// TODO: add your business logic (e.g. call usecase, persist, notify)
	_ = payload

	return nil
}

// handleUpdatedMessage processes a single .updated message. Add your logic here.
func (c *{{.StructName}}Consumer) handleUpdatedMessage(msg amqplib.Delivery) error {
	helper.LogInfo("Processing {{.Title}} updated message", map[string]interface{}{
		"source":     "{{.StructName}}Consumer.handleUpdatedMessage",
		"message_id": msg.MessageId,
	})

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		helper.LogError("Failed to unmarshal message", err, "", map[string]interface{}{
			"source":     "{{.StructName}}Consumer.handleUpdatedMessage",
			"message_id": msg.MessageId,
		})
		return err
	}

	// TODO: add your business logic
	_ = payload

	return nil
}
`
	data := map[string]interface{}{
		"StructName":  c.StructName,
		"ConstPrefix": c.ConstPrefix,
		"Title":       c.Title,
		"Identifier": c.Identifier,
	}
	return c.writeFile("src/consumer/"+c.Identifier+".go", tmpl, data)
}

func (c *ConsumerGen) writeFile(filePath string, tmplStr string, data interface{}) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: %s", filePath)
	}
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
