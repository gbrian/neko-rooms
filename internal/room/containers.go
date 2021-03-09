package room

import (
	"context"
	"fmt"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"

	"m1k1o/neko_rooms/internal/types"
)

func (manager *RoomManagerCtx) containerToEntry(container dockerTypes.Container) (*types.RoomEntry, error) {
	roomName, ok := container.Labels["m1k1o.neko_rooms.name"]
	if !ok {
		return nil, fmt.Errorf("Damaged container labels: name not found.")
	}

	URL, ok := container.Labels["m1k1o.neko_rooms.url"]
	if !ok {
		return nil, fmt.Errorf("Damaged container labels: url not found.")
	}

	epr, err := manager.getEprFromLabels(container.Labels)
	if err != nil {
		return nil, err
	}

	return &types.RoomEntry{
		ID:             container.ID,
		URL:            URL,
		Name:           roomName,
		MaxConnections: epr.Max - epr.Min + 1,
		Image:          container.Image,
		Running:        container.State == "running",
		Status:         container.Status,
		Created:        time.Unix(container.Created, 0),
	}, nil
}

func (manager *RoomManagerCtx) listContainers() ([]dockerTypes.Container, error) {
	args := filters.NewArgs(
		filters.Arg("label", "m1k1o.neko_rooms.instance"),
	)

	containers, err := manager.client.ContainerList(context.Background(), dockerTypes.ContainerListOptions{
		All:     true,
		Filters: args,
	})

	if err != nil {
		return nil, err
	}

	result := []dockerTypes.Container{}
	for _, container := range containers {
		val, ok := container.Labels["m1k1o.neko_rooms.instance"]
		if !ok || val != manager.config.InstanceName {
			continue
		}

		result = append(result, container)
	}

	return result, nil
}

func (manager *RoomManagerCtx) containerInfo(id string) (*dockerTypes.Container, error) {
	args := filters.NewArgs(
		filters.Arg("id", id),
		filters.Arg("label", "m1k1o.neko_rooms.instance"),
	)

	containers, err := manager.client.ContainerList(context.Background(), dockerTypes.ContainerListOptions{
		All:     true,
		Filters: args,
	})

	if err != nil {
		return nil, err
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("Container not found.")
	}

	container := containers[0]

	val, ok := container.Labels["m1k1o.neko_rooms.instance"]
	if !ok || val != manager.config.InstanceName {
		return nil, fmt.Errorf("This container does not belong to neko_rooms.")
	}

	return &container, nil
}

func (manager *RoomManagerCtx) inspectContainer(id string) (*dockerTypes.ContainerJSON, error) {
	container, err := manager.client.ContainerInspect(context.Background(), id)
	if err != nil {
		return nil, err
	}

	val, ok := container.Config.Labels["m1k1o.neko_rooms.instance"]
	if !ok || val != manager.config.InstanceName {
		return nil, fmt.Errorf("This container does not belong to neko_rooms.")
	}

	return &container, nil
}
