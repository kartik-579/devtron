package cluster

import (
	"errors"
	"github.com/devtron-labs/devtron/internal/util"
	repository2 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/cluster/repository/mocks"
	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestForEphemeralContainers(t *testing.T) {

	t.Run("TestSaveEphemeralContainer_Success", func(t *testing.T) {
		repository := mocks.NewEphemeralContainersRepository(t)

		// Set up the expected repository method calls and return values
		repository.On("FindContainerByName", 1, "namespace-1", "pod-1", "container-1").Return(nil, nil)
		repository.On("StartTx").Return(&pg.Tx{}, nil)
		repository.On("RollbackTx", mock.AnythingOfType("*pg.Tx")).Return(nil) // Add this expectation
		repository.On("SaveData", mock.AnythingOfType("*pg.Tx"), mock.AnythingOfType("*repository.EphemeralContainerBean")).Return(nil)
		repository.On("SaveAction", mock.AnythingOfType("*pg.Tx"), mock.AnythingOfType("*repository.EphemeralContainerAction")).Return(nil)
		repository.On("CommitTx", mock.AnythingOfType("*pg.Tx")).Return(nil)
		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.SaveEphemeralContainer(request)

		assert.NoError(t, err)

		repository.AssertCalled(t, "FindContainerByName", 1, "namespace-1", "pod-1", "container-1")
		repository.AssertCalled(t, "StartTx")
		repository.AssertCalled(t, "RollbackTx", mock.AnythingOfType("*pg.Tx")) // Add this assertion
		repository.AssertCalled(t, "SaveData", mock.AnythingOfType("*pg.Tx"), mock.AnythingOfType("*repository.EphemeralContainerBean"))
		repository.AssertCalled(t, "SaveAction", mock.AnythingOfType("*pg.Tx"), mock.AnythingOfType("*repository.EphemeralContainerAction"))
		repository.AssertCalled(t, "CommitTx", mock.AnythingOfType("*pg.Tx"))
	})

	t.Run("TestSaveEphemeralContainer_FindContainerError", func(t *testing.T) {
		repository := mocks.NewEphemeralContainersRepository(t)

		repository.On("FindContainerByName", mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("error finding container"))

		logger, _ := util.NewSugardLogger()
		service := EphemeralContainerServiceImpl{
			repository: repository,
			logger:     logger,
		}

		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.SaveEphemeralContainer(request)

		assert.Error(t, err)

		repository.AssertCalled(t, "FindContainerByName", 1, "namespace-1", "pod-1", "container-1")
	})

	t.Run("TestAuditEphemeralContainerAction_Success", func(t *testing.T) {
		repository := mocks.NewEphemeralContainersRepository(t)
		tx := &pg.Tx{}
		// Set up the expected repository method calls and return values
		repository.On("FindContainerByName", 1, "namespace-1", "pod-1", "container-1").Return(nil, nil)
		repository.On("StartTx").Return(tx, nil)
		repository.On("RollbackTx", tx).Return(nil)
		repository.On("SaveData", tx, mock.AnythingOfType("*repository.EphemeralContainerBean")).Return(nil)
		repository.On("SaveAction", tx, mock.AnythingOfType("*repository.EphemeralContainerAction")).Return(nil)
		repository.On("CommitTx", tx).Return(nil)
		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.AuditEphemeralContainerAction(request, repository2.ActionAccessed)

		assert.NoError(t, err)
	})

	t.Run("TestAuditEphemeralContainerAction_FindContainerError", func(t *testing.T) {
		repository := mocks.NewEphemeralContainersRepository(t)

		repository.On("FindContainerByName", mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("error finding container"))

		logger, _ := util.NewSugardLogger()
		service := EphemeralContainerServiceImpl{
			repository: repository,
			logger:     logger,
		}

		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.AuditEphemeralContainerAction(request, repository2.ActionAccessed)

		assert.Error(t, err)

		repository.AssertCalled(t, "FindContainerByName", 1, "namespace-1", "pod-1", "container-1")
	})

	t.Run("TestSaveEphemeralContainer_ContainerAlreadyPresentError", func(t *testing.T) {
		repository := mocks.NewEphemeralContainersRepository(t)

		// Set up the expected repository method calls and return values
		container := &repository2.EphemeralContainerBean{
			PodName: "existing-pod",
		}
		repository.On("FindContainerByName", mock.AnythingOfType("int"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(container, nil)

		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.SaveEphemeralContainer(request)

		assert.Error(t, err)
		assert.EqualError(t, err, "container already present in the provided pod")

		repository.AssertCalled(t, "FindContainerByName", 1, "namespace-1", "pod-1", "container-1")
	})

	t.Run("TestAuditEphemeralContainerAction_CommitError", func(t *testing.T) {

		repository := mocks.NewEphemeralContainersRepository(t)
		tx := &pg.Tx{}
		// Set up the expected repository method calls and return values
		repository.On("FindContainerByName", 1, "namespace-1", "pod-1", "container-1").Return(nil, nil)
		repository.On("StartTx").Return(tx, nil)
		repository.On("RollbackTx", tx).Return(nil)
		repository.On("SaveData", tx, mock.AnythingOfType("*repository.EphemeralContainerBean")).Return(nil)
		repository.On("SaveAction", tx, mock.AnythingOfType("*repository.EphemeralContainerAction")).Return(nil)
		repository.On("CommitTx", tx).Return(errors.New("error committing transaction")) // Return an error during commit
		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.AuditEphemeralContainerAction(request, repository2.ActionAccessed)

		assert.Error(t, err)
		assert.EqualError(t, err, "error committing transaction")

		repository.AssertCalled(t, "FindContainerByName", 1, "namespace-1", "pod-1", "container-1")
		repository.AssertCalled(t, "StartTx")
		repository.AssertCalled(t, "RollbackTx", tx)
		repository.AssertCalled(t, "SaveData", tx, mock.AnythingOfType("*repository.EphemeralContainerBean"))
		repository.AssertCalled(t, "SaveAction", tx, mock.AnythingOfType("*repository.EphemeralContainerAction"))
		repository.AssertCalled(t, "CommitTx", tx)
	})

	t.Run("TestAuditEphemeralContainerAction_SaveActionError", func(t *testing.T) {

		repository := mocks.NewEphemeralContainersRepository(t)
		tx := &pg.Tx{}

		// Set up the expected repository method calls and return values
		repository.On("FindContainerByName", 1, "namespace-1", "pod-1", "container-1").Return(nil, nil)
		repository.On("StartTx").Return(tx, nil)
		repository.On("RollbackTx", tx).Return(nil)
		repository.On("SaveData", tx, mock.AnythingOfType("*repository.EphemeralContainerBean")).Return(nil)
		repository.On("SaveAction", tx, mock.AnythingOfType("*repository.EphemeralContainerAction")).Return(errors.New("failed to save action")) // Return an error when saving action

		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.AuditEphemeralContainerAction(request, repository2.ActionAccessed)

		assert.Error(t, err)

	})

	t.Run("TestAuditEphemeralContainerAction_SaveDataError", func(t *testing.T) {
		repository := mocks.NewEphemeralContainersRepository(t)
		tx := &pg.Tx{}

		// Set up the expected repository method calls and return values
		repository.On("FindContainerByName", 1, "namespace-1", "pod-1", "container-1").Return(nil, nil)
		repository.On("StartTx").Return(tx, nil)
		repository.On("RollbackTx", tx).Return(nil)
		repository.On("SaveData", tx, mock.AnythingOfType("*repository.EphemeralContainerBean")).Return(errors.New("error saving data")) // Return an error during SaveData
		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.AuditEphemeralContainerAction(request, repository2.ActionAccessed)

		assert.Error(t, err)

	})

	t.Run("TestAuditEphemeralContainerAction_CreateTransactionError", func(t *testing.T) {

		repository := mocks.NewEphemeralContainersRepository(t)
		tx := &pg.Tx{}
		// Set up the expected repository method calls and return values
		repository.On("FindContainerByName", 1, "namespace-1", "pod-1", "container-1").Return(nil, nil)
		repository.On("StartTx").Return(tx, errors.New("error creating transaction")) // Simulate error in creating transaction
		repository.On("RollbackTx", tx).Return(nil)
		logger, _ := util.NewSugardLogger()
		service := NewEphemeralContainerServiceImpl(repository, logger)

		// Create a sample EphemeralContainerRequest
		request := EphemeralContainerRequest{
			BasicData: &EphemeralContainerBasicData{
				ContainerName:       "container-1",
				TargetContainerName: "target-container-1",
				Image:               "image-1",
			},
			AdvancedData: &EphemeralContainerAdvancedData{
				Manifest: "manifest-1",
			},
			Namespace: "namespace-1",
			ClusterId: 1,
			PodName:   "pod-1",
			UserId:    123,
		}

		err := service.AuditEphemeralContainerAction(request, repository2.ActionAccessed)

		assert.Error(t, err)

	})

}
