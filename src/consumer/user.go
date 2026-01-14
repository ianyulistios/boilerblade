package consumer

import (
	"boilerblade/config/amqp"
	"boilerblade/constants"
	"boilerblade/helper"
	"boilerblade/src/dto"
	"boilerblade/src/usecase"
	"encoding/json"
	"fmt"

	amqplib "github.com/streadway/amqp"
)

// UserConsumer handles AMQP messages for user operations
type UserConsumer struct {
	subConnection amqp.IAMQPChannel
	userUsecase   usecase.UserUsecase
}

// NewUserConsumer creates a new user consumer instance
// It sets up the channel and declares the exchange
func NewUserConsumer(amqpConn amqp.IAMQPConnection, userUsecase usecase.UserUsecase) (*UserConsumer, error) {
	// Get AMQP channel
	channel, err := amqpConn.Channel()
	if err != nil {
		helper.LogError("Failed to get AMQP channel for UserConsumer", err, "", map[string]interface{}{
			"source": "NewUserConsumer",
		})
		return nil, err
	}

	// Declare exchange for user events
	exchangeName := constants.UserExchangeName
	if err := channel.DeclareExchange(exchangeName, "direct"); err != nil {
		helper.LogError("Failed to declare exchange in UserConsumer", err, exchangeName, map[string]interface{}{
			"source": "NewUserConsumer",
		})
		channel.Close()
		return nil, err
	}

	helper.LogInfo("UserConsumer initialized with channel and exchange", map[string]interface{}{
		"source":        "NewUserConsumer",
		"exchange":      exchangeName,
		"exchange_type": "direct",
	})

	return &UserConsumer{
		subConnection: channel,
		userUsecase:   userUsecase,
	}, nil
}

// UserCreateMessage represents the message payload for creating a user
type UserCreateMessage struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Close closes the AMQP channel
func (c *UserConsumer) Close() error {
	if c.subConnection != nil {
		return c.subConnection.Close()
	}
	return nil
}

// ProcessMessage processes messages (placeholder)
func (c *UserConsumer) ProcessMessage() {
	go c.ProcessUserCreated()
	go c.ProcessUserUpdated()
}

// ProcessUserCreated processes user creation messages from AMQP
func (c *UserConsumer) ProcessUserCreated() {
	var messages <-chan amqplib.Delivery

	// Declare queue
	que, err := c.subConnection.NewQueue(
		constants.UserExchangeName,
		constants.UserCreatedQueueName,
		constants.QueueType,
		constants.UserCreatedRouteKey,
		constants.UserCreatedQueueInterval,
	)
	if err != nil {
		helper.LogError("(UserCreated) Error when declare queue", err, constants.UserCreatedQueueName, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Read messages
	if messages, err = c.subConnection.ReadMessage(que); err != nil {
		fmt.Println(err)
		return
	}

	helper.LogInfo("Started consuming user.created messages", map[string]interface{}{
		"source": "UserConsumer.ProcessUserCreated",
		"queue":  constants.UserCreatedQueueName,
	})

	// Process messages
	for msg := range messages {
		if err := c.handleUserCreatedMessage(msg); err != nil {
			helper.LogError("Failed to process user.created message", err, "", map[string]interface{}{
				"source":     "UserConsumer.ProcessUserCreated",
				"message_id": msg.MessageId,
			})
			// Nack message to retry
			msg.Nack(false, true)
			continue
		}
		// Ack message on success
		msg.Ack(false)
	}
}

// ProcessUserUpdated processes user update messages from AMQP
func (c *UserConsumer) ProcessUserUpdated() {
	var messages <-chan amqplib.Delivery

	// Declare queue
	que, err := c.subConnection.NewQueue(
		constants.UserExchangeName,
		constants.UserUpdatedQueueName,
		constants.QueueType,
		constants.UserUpdatedRouteKey,
		constants.UserUpdatedQueueInterval,
	)
	if err != nil {
		helper.LogError("(UserUpdated) Error when declare queue", err, constants.UserUpdatedQueueName, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Read messages
	if messages, err = c.subConnection.ReadMessage(que); err != nil {
		fmt.Println(err)
		return
	}

	helper.LogInfo("Started consuming user.updated messages", map[string]interface{}{
		"source": "UserConsumer.ProcessUserUpdated",
		"queue":  constants.UserUpdatedQueueName,
	})

	// Process messages
	for msg := range messages {
		if err := c.handleUserUpdatedMessage(msg); err != nil {
			helper.LogError("Failed to process user.updated message", err, "", map[string]interface{}{
				"source":     "UserConsumer.ProcessUserUpdated",
				"message_id": msg.MessageId,
			})
			// Nack message to retry
			msg.Nack(false, true)
			continue
		}
		// Ack message on success
		msg.Ack(false)
	}
}

// handleUserCreatedMessage processes a single user creation message
func (c *UserConsumer) handleUserCreatedMessage(msg amqplib.Delivery) error {
	helper.LogInfo("Processing user creation message", map[string]interface{}{
		"source":      "UserConsumer.handleUserCreatedMessage",
		"message_id":  msg.MessageId,
		"routing_key": msg.RoutingKey,
	})

	// Parse message body
	var userMsg UserCreateMessage
	if err := json.Unmarshal(msg.Body, &userMsg); err != nil {
		helper.LogError("Failed to unmarshal user message", err, "", map[string]interface{}{
			"source":     "UserConsumer.handleUserCreatedMessage",
			"message_id": msg.MessageId,
			"body":       string(msg.Body),
		})
		return err
	}

	// Validate required fields
	if userMsg.Name == "" || userMsg.Email == "" {
		helper.LogError("Invalid user message: missing required fields", nil, "", map[string]interface{}{
			"source":     "UserConsumer.handleUserCreatedMessage",
			"message_id": msg.MessageId,
			"user_msg":   userMsg,
		})
		return nil // Return nil to ack the message even if invalid
	}

	// Create user using usecase
	createReq := &dto.CreateUserRequest{
		Name:     userMsg.Name,
		Email:    userMsg.Email,
		Password: userMsg.Password,
	}

	userResponse, err := c.userUsecase.CreateUser(createReq)
	if err != nil {
		// Check if error is due to email already exists
		if err.Error() == "email already exists" {
			helper.LogInfo("User already exists, skipping", map[string]interface{}{
				"source":     "UserConsumer.handleUserCreatedMessage",
				"message_id": msg.MessageId,
				"email":      userMsg.Email,
			})
			return nil // User already exists, ack the message
		}

		helper.LogError("Failed to create user from message", err, "", map[string]interface{}{
			"source":     "UserConsumer.handleUserCreatedMessage",
			"message_id": msg.MessageId,
			"email":      userMsg.Email,
		})
		return err // Return error to nack and retry
	}

	helper.LogInfo("User created successfully from message", map[string]interface{}{
		"source":     "UserConsumer.handleUserCreatedMessage",
		"message_id": msg.MessageId,
		"user_id":    userResponse.ID,
		"email":      userResponse.Email,
	})

	return nil // Success, message will be acked
}

// handleUserUpdatedMessage processes a single user update message
func (c *UserConsumer) handleUserUpdatedMessage(msg amqplib.Delivery) error {
	helper.LogInfo("Processing user update message", map[string]interface{}{
		"source":     "UserConsumer.handleUserUpdatedMessage",
		"message_id": msg.MessageId,
	})

	// Parse message body
	var userMsg struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(msg.Body, &userMsg); err != nil {
		helper.LogError("Failed to unmarshal user update message", err, "", map[string]interface{}{
			"source":     "UserConsumer.handleUserUpdatedMessage",
			"message_id": msg.MessageId,
		})
		return err
	}

	// Validate user ID
	if userMsg.ID == 0 {
		helper.LogError("Invalid user ID in update message", nil, "", map[string]interface{}{
			"source":     "UserConsumer.handleUserUpdatedMessage",
			"message_id": msg.MessageId,
		})
		return nil // Return nil to ack the message even if invalid
	}

	// Update user using usecase
	updateReq := &dto.UpdateUserRequest{
		Name:     userMsg.Name,
		Email:    userMsg.Email,
		Password: userMsg.Password,
	}

	userResponse, err := c.userUsecase.UpdateUser(userMsg.ID, updateReq)
	if err != nil {
		// Check if error is due to user not found
		if err.Error() == "user not found" {
			helper.LogInfo("User not found for update, skipping", map[string]interface{}{
				"source":     "UserConsumer.handleUserUpdatedMessage",
				"message_id": msg.MessageId,
				"user_id":    userMsg.ID,
			})
			return nil // User not found, ack message
		}

		helper.LogError("Failed to update user from message", err, "", map[string]interface{}{
			"source":     "UserConsumer.handleUserUpdatedMessage",
			"message_id": msg.MessageId,
			"user_id":    userMsg.ID,
		})
		return err
	}

	helper.LogInfo("User updated successfully from message", map[string]interface{}{
		"source":     "UserConsumer.handleUserUpdatedMessage",
		"message_id": msg.MessageId,
		"user_id":    userResponse.ID,
	})

	return nil
}
