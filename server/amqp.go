package server

import (
	"boilerblade/constants"
	"boilerblade/helper"
	"boilerblade/src/consumer"
	"boilerblade/src/repository"
	"boilerblade/src/usecase"
)

// AMQPServe initializes and serves AMQP consumers
// This method ensures AMQP connection is available before use
// If AMQP was disabled via ENABLE_AMQP=false, it will be force-enabled
func (a *App) AMQPServe() error {
	// Ensure AMQP connection is initialized (using method from config)
	if err := a.Config.EnsureAMQP(); err != nil {
		helper.LogError("Failed to ensure AMQP connection for AMQPServe", err, "", map[string]interface{}{
			"source": "AMQPServe",
		})
		return err
	}

	// AMQP connection is now available
	helper.LogInfo("AMQP serve started", map[string]interface{}{
		"source": "AMQPServe",
	})

	// Initialize consumer dependencies
	userRepo := repository.NewUserRepository(a.Config.Database)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Create user.created consumer (channel and exchange are set up in NewUserConsumer)
	userCreatedConsumer, err := consumer.NewUserConsumer(a.Config.AMQP, userUsecase)
	if err != nil {
		helper.LogError("Failed to create user.created consumer", err, "", map[string]interface{}{
			"source": "AMQPServe",
		})
		return err
	}

	// Create user.updated consumer (channel and exchange are set up in NewUserConsumer)
	userUpdatedConsumer, err := consumer.NewUserConsumer(a.Config.AMQP, userUsecase)
	if err != nil {
		helper.LogError("Failed to create user.updated consumer", err, "", map[string]interface{}{
			"source": "AMQPServe",
		})
		userCreatedConsumer.Close()
		return err
	}

	// Start consuming user.created messages
	go func() {
		defer userCreatedConsumer.Close()
		userCreatedConsumer.ProcessUserCreated()
	}()

	// Start consuming user.updated messages
	go func() {
		defer userUpdatedConsumer.Close()
		userUpdatedConsumer.ProcessUserUpdated()
	}()

	helper.LogInfo("All AMQP consumers started", map[string]interface{}{
		"source": "AMQPServe",
		"queues": []string{constants.UserCreatedQueueName, constants.UserUpdatedQueueName},
	})

	// Wait for interrupt signal to gracefully shutdown
	a.WaitForShutdown()

	return nil
}
