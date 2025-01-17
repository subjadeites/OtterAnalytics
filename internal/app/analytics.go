package app

import (
	"OtterAnalytics/handlers"
	"OtterAnalytics/models"
	"OtterAnalytics/pkg/errors"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
)

type CollectJSON struct {
	ClientID       string         `json:"client_id"`
	UserID         string         `json:"user_id"`
	UserProperties UserProperties `json:"user_properties"`
	Events         []Events       `json:"events"`
}

type UserProperties struct {
	HomeWorld            PropertyValueString     `json:"HomeWorld"`
	CheatBannedHashValid PropertyValueBool       `json:"Cheat_Banned_Hash_Valid"`
	Client               PropertyValueString     `json:"Client"`
	OS                   PropertyValueString     `json:"os"`
	DalamudVersion       PropertyValueString     `json:"dalamud_version"`
	IsTesting            PropertyValueBool       `json:"is_testing"`
	PluginCount          PropertyValueInt        `json:"plugin_count"`
	MachineID            PropertyValueString     `json:"machine_id"`
	Plugin3rdList        PropertyValueListString `json:"plugin_3rd_list"`
}

type PropertyValueInt struct {
	Value int `json:"value"`
}
type PropertyValueString struct {
	Value string `json:"value"`
}
type PropertyValueBool struct {
	Value bool `json:"value"`
}
type PropertyValueListString struct {
	Value []string `json:"value"`
}

type Events struct {
	Name   string      `json:"name"`
	Params EventParams `json:"params"`
}
type EventParams struct {
	ServerId           string `json:"server_id"`
	EngagementTimeMsec string `json:"engagement_time_msec"`
	SessionId          string `json:"session_id"`
}

func handlerCollectPost(h *handlers.Handler, r *http.Request, conn net.Conn) {
	var readData handlers.RequestDataReader = &handlers.ReadPostData{}
	body, contentType, err := readData.ReadData(r)
	if err != nil {
		handlers.WriteMethodNotAllowedResponse(conn)
		return
	}

	if contentType != "json" {
		handlers.WriteServerErrorResponse(conn)
		return
	}

	var collectJSON CollectJSON
	if err := json.Unmarshal(body, &collectJSON); err != nil {
		handlers.WriteNotFoundResponse(conn)
		return
	}

	plugin3rdListJSON, err := json.Marshal(collectJSON.UserProperties.Plugin3rdList.Value)
	if err != nil {
		log.Printf("Error marshalling plugin3rdList: %v", err)
		handlers.WriteServerErrorResponse(conn)
		return
	}

	homeWorld, err := strconv.Atoi(collectJSON.UserProperties.HomeWorld.Value)
	errors.Normal(err, "Error converting HomeWorld to int: %v")

	machineID := collectJSON.UserProperties.MachineID.Value
	if machineID == "" {
		machineID = "None"
	}

	machineIDPlugin := models.MachineIDPlugin{
		MachineID:      machineID,
		Plugin3rdList:  string(plugin3rdListJSON),
		Plugin3rdCount: collectJSON.UserProperties.PluginCount.Value,
	}

	event := models.Event{
		ClientID:             collectJSON.ClientID,
		EventType:            collectJSON.Events[0].Name,
		UserID:               collectJSON.UserID,
		HomeWorld:            homeWorld,
		CheatBannedHashValid: collectJSON.UserProperties.CheatBannedHashValid.Value,
		Client:               collectJSON.UserProperties.Client.Value,
		OS:                   collectJSON.UserProperties.OS.Value,
		DalamudVersion:       collectJSON.UserProperties.DalamudVersion.Value,
		IsTesting:            collectJSON.UserProperties.IsTesting.Value,
		Plugin3rdCount:       collectJSON.UserProperties.PluginCount.Value,
		MachineID:            collectJSON.UserProperties.MachineID.Value,
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
			if r != nil {
				panic(r)
			}
		} else {
			tx.Commit()
		}
	}()

	if err := tx.Save(&machineIDPlugin).Error; err != nil {
		log.Printf("Error inserting machineIDPlugin: %v", err)
		return
	}

	if err := tx.Create(&event).Error; err != nil {
		log.Printf("Error inserting event: %v", err)
		return
	}

	if err := handlers.WriteResponse(conn, handlers.NewHeader(), "success"); err != nil {
		log.Println("Error writing response:", err)
	}

}

func RegisterAnalyticsRoutes() {
	handlers.Routes[http.MethodPost]["/collect"] = handlerCollectPost
}
