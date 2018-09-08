package common

type Task struct {
	Driver ContainerDriver // Container driver for managing Containers

	UUID       string   // Identification of the task
	Containers []string // Container's docker IDs
}

func (t *Task) AddContainer() error {
	containerId, err := t.Driver.Start(t.UUID)
	if err == nil {
		t.Containers = append(t.Containers, containerId)
	}
	return err
}

func (t *Task) DestroyAll() error {
	var remainingIds []string
	var lastErr error

	for _, containerId := range t.Containers {
		err := t.Driver.Kill(containerId)
		if err != nil {
			lastErr = err
			remainingIds = append(remainingIds, containerId)
		}
	}

	t.Containers = remainingIds

	return lastErr
}
