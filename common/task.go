package common

type Task struct {
	Driver ContainerDriver // Container driver for managing containers

	UUID       string   // Identification of the task
	containers []string // Container's docker IDs
}

func (t *Task) AddContainer() error {
	containerId, err := t.Driver.Start(t.UUID)
	if err == nil {
		t.containers = append(t.containers, containerId)
	}
	return err
}

func (t *Task) DestroyAll() error {
	var remainingIds []string
	var lastErr error

	for _, containerId := range t.containers {
		err := t.Driver.Kill(containerId)
		if err != nil {
			lastErr = err
			remainingIds = append(remainingIds, containerId)
		}
	}

	t.containers = remainingIds

	return lastErr
}
