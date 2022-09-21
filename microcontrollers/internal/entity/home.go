package entity

import "errors"

type Home struct {
	ID           string `json:"id,omitempty"`
	ClientId     string `json:"client_id"`
	Temperature  string `json:"temperature"`
	IsGateOpened bool   `json:"is_gate_opened"`
	IsRobbery    bool   `json:"is_robbery"`
	IsLedTurned  bool   `json:"is_led_turned"`
	SecureMode   bool   `json:"secure_mode"`
}

type CreateHomeInput struct {
	HomeId   *string `json:"home_id"`
	ClientId *string `json:"client_id"`
}

type UpdateHomeInput struct {
	Temperature *string `json:"temperature"`
	IsRobbery   *bool   `json:"is_robbery"`
}

type UpdateHomeCommandInput struct {
	NewClientId *string `json:"new_client_id""`
	SecureMode  *bool   `json:"secure_mode"`
	OpenGate    *bool   `json:"open_gate"`
	LedTurn     *bool   `json:"turn_led"`
}

func (i UpdateHomeInput) Validate() error {
	if i.Temperature == nil && i.IsRobbery == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

func (i UpdateHomeCommandInput) Validate() error {
	if i.SecureMode == nil && i.NewClientId == nil && i.OpenGate == nil && i.LedTurn == nil {
		return errors.New("update security structure has no values")
	}

	return nil
}

func (i CreateHomeInput) Validate() error {
	if i.ClientId == nil || i.HomeId == nil {
		return errors.New("no client id")
	}

	return nil
}
