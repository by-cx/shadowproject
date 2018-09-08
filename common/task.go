package common

type Task struct {
	Driver ContainerDriver // Container driver for managing Containers

	UUID       string   // Identification of the task
	Containers []string // Container's docker IDs
	Domains    []string // Domain list on which this tasks listens
	Image      string   // Docker image
	Command    []string // Command to run
}

// Adds new container for this task. Returns container ID and error.
func (t *Task) AddContainer() (string, error) {
	containerId, err := t.Driver.Start(t.UUID, t.Image, t.Command)
	if err == nil {
		t.Containers = append(t.Containers, containerId)
	}
	return containerId, err
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
