package cache

import (
	"testing"

	"github.com/commojun/nyanbot/masterdata/key_value"
	"github.com/commojun/nyanbot/masterdata/table"
)

func TestCache(t *testing.T) {
	// Reset global cache for testing
	globalCache = nil

	// Test before initialization
	t.Run("Getters before initialization", func(t *testing.T) {
		if len(GetAlarms()) != 0 {
			t.Errorf("GetAlarms() should return empty slice before initialization")
		}
		if len(GetAnniversaries()) != 0 {
			t.Errorf("GetAnniversaries() should return empty slice before initialization")
		}
		if _, err := GetRoomID("test"); err == nil {
			t.Errorf("GetRoomID() should return error before initialization")
		}
		if _, err := GetNickname("test"); err == nil {
			t.Errorf("GetNickname() should return error before initialization")
		}
	})

	// Setup mock data
	mockCache := &Cache{
		Tables: &table.Tables{
			Alarms: []table.Alarm{
				{ID: "1", Message: "Test Alarm"},
			},
			Anniversaries: []table.Anniversary{
				{ID: "1", Name: "Test Anniversary"},
			},
		},
		KeyVals: &key_value.KVs{
			Rooms: map[string]string{
				"test_room": "room123",
			},
			Nickname: map[string]string{
				"user123": "Test User",
			},
		},
	}

	// Set the mock cache
	SetTestCache(mockCache)

	// Test after setting mock cache
	t.Run("Getters with mock cache", func(t *testing.T) {
		// Test GetAlarms
		alarms := GetAlarms()
		if len(alarms) != 1 || alarms[0].ID != "1" {
			t.Errorf("GetAlarms() returned unexpected data: %v", alarms)
		}

		// Test GetAnniversaries
		anniversaries := GetAnniversaries()
		if len(anniversaries) != 1 || anniversaries[0].ID != "1" {
			t.Errorf("GetAnniversaries() returned unexpected data: %v", anniversaries)
		}

		// Test GetRoomID
		roomID, err := GetRoomID("test_room")
		if err != nil || roomID != "room123" {
			t.Errorf("GetRoomID() returned unexpected data: id=%s, err=%v", roomID, err)
		}
		_, err = GetRoomID("not_exist")
		if err == nil {
			t.Errorf("GetRoomID() should return error for non-existent key")
		}

		// Test GetNickname
		nickname, err := GetNickname("user123")
		if err != nil || nickname != "Test User" {
			t.Errorf("GetNickname() returned unexpected data: nickname=%s, err=%v", nickname, err)
		}
		_, err = GetNickname("not_exist")
		if err == nil {
			t.Errorf("GetNickname() should return error for non-existent user")
		}
	})
}
