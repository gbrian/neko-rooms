package room

import (
	"fmt"
	"strconv"

	"github.com/m1k1o/neko-rooms/internal/types"
)

type RoomLabels struct {
	Name      string
	URL       string
	Mux       bool
	Epr       EprPorts
	NekoImage string

	BrowserPolicy *BrowserPolicyLabels

	ConstrainsLabels *ConstrainsLabels
}

type BrowserPolicyLabels struct {
	Type types.BrowserPolicyType
	Path string
}

type ConstrainsLabels struct {
	TimeToLiveConstrain types.TimeToLiveConstrain
	EmptyRoomConstrain types.EmptyRoomConstrain
}

func (manager *RoomManagerCtx) extractLabels(labels map[string]string) (*RoomLabels, error) {
	name, ok := labels["m1k1o.neko_rooms.name"]
	if !ok {
		return nil, fmt.Errorf("damaged container labels: name not found")
	}

	url, ok := labels["m1k1o.neko_rooms.url"]
	if !ok {
		// TODO: It should be always available.
		url = manager.config.GetRoomUrl(name)
		//return nil, fmt.Errorf("damaged container labels: url not found")
	}

	nekoImage, ok := labels["m1k1o.neko_rooms.neko_image"]
	if !ok {
		return nil, fmt.Errorf("damaged container labels: neko_image not found")
	}

	var mux bool
	var epr EprPorts

	muxStr, ok := labels["m1k1o.neko_rooms.mux"]
	if ok {
		muxPort, err := strconv.ParseUint(muxStr, 10, 16)
		if err != nil {
			return nil, err
		}

		mux = true
		epr = EprPorts{
			Min: uint16(muxPort),
			Max: uint16(muxPort),
		}
	} else {
		eprMinStr, ok := labels["m1k1o.neko_rooms.epr.min"]
		if !ok {
			return nil, fmt.Errorf("damaged container labels: epr.min not found")
		}

		eprMin, err := strconv.ParseUint(eprMinStr, 10, 16)
		if err != nil {
			return nil, err
		}

		eprMaxStr, ok := labels["m1k1o.neko_rooms.epr.max"]
		if !ok {
			return nil, fmt.Errorf("damaged container labels: epr.max not found")
		}

		eprMax, err := strconv.ParseUint(eprMaxStr, 10, 16)
		if err != nil {
			return nil, err
		}

		mux = false
		epr = EprPorts{
			Min: uint16(eprMin),
			Max: uint16(eprMax),
		}
	}

	var browserPolicy *BrowserPolicyLabels
	if val, ok := labels["m1k1o.neko_rooms.browser_policy"]; ok && val == "true" {
		policyType, ok := labels["m1k1o.neko_rooms.browser_policy.type"]
		if !ok {
			return nil, fmt.Errorf("damaged container labels: browser_policy.type not found")
		}

		policyPath, ok := labels["m1k1o.neko_rooms.browser_policy.path"]
		if !ok {
			return nil, fmt.Errorf("damaged container labels: browser_policy.path not found")
		}

		browserPolicy = &BrowserPolicyLabels{
			Type: types.BrowserPolicyType(policyType),
			Path: policyPath,
		}
	}

	// var constrainsLabels *ConstrainsLabels
	// constains_path := "m1k1o.neko_rooms.constrains"
	// if val, ok := labels[constains_path]; ok && val == "true" {
	// 	if time_to_live_time, ok := strconv.ParseUint(labels[constains_path + "time_to_live.time"], 10, 16); ok {
	// 		if time_to_live_action, ok := labels[constains_path + "time_to_live.action"]; !ok {
	// 			time_to_live_action = "stop"
	// 		}
	// 	}

	// 	if empty_room_threshold, ok := strconv.ParseUint(labels[constains_path + "empty_room.threshold"], 10, 16); ok {
	// 		if empty_room_action, ok := labels[constains_path + "empty_room.action"]; !ok {
	// 			empty_room_action = "stop"
	// 		}
	// 	}

	// 	constrainsLabels = &ConstrainsLabels{
	// 		TimeToLiveConstrain: &types.TimeToLiveConstrain{
	// 			TimeToLive: time_to_live_time,
	// 			Action: time_to_live_action,
	// 		},
	// 		EmptyRoomConstrain: &types.EmptyRoomConstrain{
	// 			ThresholTime: empty_room_threshold,
	// 			Action: empty_room_action,
	// 		},
	// 	}
	// }

	return &RoomLabels{
		Name:      name,
		URL:       url,
		NekoImage: nekoImage,
		Mux:       mux,
		Epr:       epr,

		BrowserPolicy: browserPolicy,
		// ConstrainsLabelss: ConstrainsLabelss,
	}, nil
}

func (manager *RoomManagerCtx) serializeLabels(labels RoomLabels) map[string]string {
	labelsMap := map[string]string{
		"m1k1o.neko_rooms.name":       labels.Name,
		"m1k1o.neko_rooms.url":        manager.config.GetRoomUrl(labels.Name),
		"m1k1o.neko_rooms.instance":   manager.config.InstanceName,
		"m1k1o.neko_rooms.neko_image": labels.NekoImage,
	}

	if labels.Mux && labels.Epr.Min == labels.Epr.Max {
		labelsMap["m1k1o.neko_rooms.mux"] = fmt.Sprintf("%d", labels.Epr.Min)
	} else {
		labelsMap["m1k1o.neko_rooms.epr.min"] = fmt.Sprintf("%d", labels.Epr.Min)
		labelsMap["m1k1o.neko_rooms.epr.max"] = fmt.Sprintf("%d", labels.Epr.Max)
	}

	if labels.BrowserPolicy != nil {
		labelsMap["m1k1o.neko_rooms.browser_policy"] = "true"
		labelsMap["m1k1o.neko_rooms.browser_policy.type"] = string(labels.BrowserPolicy.Type)
		labelsMap["m1k1o.neko_rooms.browser_policy.path"] = labels.BrowserPolicy.Path
	}

	// if labels.ConstrainsLabels != nil {
	// 	constains_path := "m1k1o.neko_rooms.constrains"
	// 	labelsMap[constains_path] = "true"

	// 	if labels.ConstrainsLabelss.TimeToLiveConstrain.TimeToLive != nil {
	// 		labelsMap[constains_path + "time_to_live.time"] = string(labels.ConstrainsLabelss.TimeToLiveConstrain.TimeToLive)
	// 		labelsMap[constains_path + "time_to_live.action"] = string(labels.ConstrainsLabelss.TimeToLiveConstrain.Action)
	// 	}
	// 	if labels.ConstrainsLabelss.EmptyConstrainsLabels.EmptyConstrainsLabels != nil {
	// 		labelsMap[constains_path + "empty_room.threshold"] = string(labels.ConstrainsLabelss.EmptyConstrainsLabels.Threshold)
	// 		labelsMap[constains_path + "empty_room.action"] = string(labels.ConstrainsLabelss.EmptyConstrainsLabels.Action)
	// 	}
	// }	

	return labelsMap
}
