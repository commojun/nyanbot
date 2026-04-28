package masterdata

import (
	"testing"

	"github.com/commojun/nyanbot/internal/masterdata/key_value"
	"github.com/commojun/nyanbot/internal/masterdata/table"
)

func TestMasterData(t *testing.T) {
	global = nil

	t.Run("Before initialization returns empty", func(t *testing.T) {
		if len(GetTables().Alarms) != 0 {
			t.Errorf("GetTables().Alarms should be empty before initialization")
		}
		if len(GetTables().Anniversaries) != 0 {
			t.Errorf("GetTables().Anniversaries should be empty before initialization")
		}
		if _, err := GetKeyVals().RoomID("test"); err == nil {
			t.Errorf("RoomID() should return error before initialization")
		}
		if _, err := GetKeyVals().Nickname("test"); err == nil {
			t.Errorf("Nickname() should return error before initialization")
		}
	})

	SetTestData(&MasterData{
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
			Nicknames: map[string]string{
				"user123": "Test User",
			},
		},
	})

	t.Run("After initialization returns data", func(t *testing.T) {
		alarms := GetTables().Alarms
		if len(alarms) != 1 || alarms[0].ID != "1" {
			t.Errorf("unexpected alarms: %v", alarms)
		}

		anniversaries := GetTables().Anniversaries
		if len(anniversaries) != 1 || anniversaries[0].ID != "1" {
			t.Errorf("unexpected anniversaries: %v", anniversaries)
		}

		roomID, err := GetKeyVals().RoomID("test_room")
		if err != nil || roomID != "room123" {
			t.Errorf("unexpected room: id=%s, err=%v", roomID, err)
		}
		_, err = GetKeyVals().RoomID("not_exist")
		if err == nil {
			t.Errorf("non-existent room key should return error")
		}

		nickname, err := GetKeyVals().Nickname("user123")
		if err != nil || nickname != "Test User" {
			t.Errorf("unexpected nickname: %s, err=%v", nickname, err)
		}
		_, err = GetKeyVals().Nickname("not_exist")
		if err == nil {
			t.Errorf("non-existent user should return error")
		}
	})
}
