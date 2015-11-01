package gochatwork

import (
	"fmt"
	"net/url"
)

// Rooms return rooms response by []Room
func (c *Client) Rooms() ([]Room, error) {
	var rooms []Room

	b, err := c.RoomsRaw()
	err = setSturctFromJSON(b, &rooms, err)
	return rooms, err
}

// RoomsRaw return rooms response by []byte
func (c *Client) RoomsRaw() ([]byte, error) {
	return c.connection.Get("rooms", url.Values{}, c.config)
}

// Room return rooms/room_id response by Room
func (c *Client) Room(roomID int64) (Room, error) {
	var room Room

	b, err := c.RoomRaw(roomID)
	err = setSturctFromJSON(b, &room, err)
	return room, err
}

// RoomRaw return rooms/room_id response by []byte
func (c *Client) RoomRaw(roomID int64) ([]byte, error) {
	return c.connection.Get(fmt.Sprintf("rooms/%d", roomID), url.Values{}, c.config)
}

// PutRooms return PUT rooms/room_id response by int64
func (c *Client) PutRooms(roomID int64, description string, iconPreset string, name string) (int64, error) {
	var responseJSON = struct {
		RoomID int64 `json:"room_id"`
	}{}

	b, err := c.PutRoomsRaw(roomID, description, iconPreset, name)
	err = setSturctFromJSON(b, &responseJSON, err)
	return responseJSON.RoomID, err
}

// PutRoomsRaw return PUT rooms/room_id response by []byte
func (c *Client) PutRoomsRaw(roomID int64, description string, iconPreset string, name string) ([]byte, error) {
	params := url.Values{}
	if description != "" {
		params.Add("description", description)
	}

	if iconPreset != "" {
		params.Add("icon_preset", iconPreset)
	}

	if name != "" {
		params.Add("name", name)
	}

	return c.connection.Put(fmt.Sprintf("rooms/%d", roomID), params, c.config)
}

// DeleteRooms send DELETE rooms/room_id response, this api don't return response
func (c *Client) DeleteRooms(roomID int64, actionType string) error {
	params := url.Values{}
	if actionType != "" {
		params.Add("action_type", actionType)
	}

	_, err := c.connection.Delete(fmt.Sprintf("rooms/%d", roomID), params, c.config)
	return err
}
