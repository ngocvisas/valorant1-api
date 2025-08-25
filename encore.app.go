// encore.app.go
//
// Valorant1 Agent & Loadout Explorer — Encore Backend
// Public endpoints cho agents/weapons + authenticated CRUD cho loadouts.
// Ready cho local dev (encore run) và Encore Cloud (encore app deploy).

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"encore.dev/beta/auth"
)

// Agent model
type Agent struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Description string   `json:"description"`
	Abilities   []string `json:"abilities"`
	ImageURL    string   `json:"imageUrl"`
}

// Weapon model
type Weapon struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"` // Primary, Sidearm
	Cost     int    `json:"cost"`
	Damage   int    `json:"damage"`
	Accuracy int    `json:"accuracy"`
	ImageURL string `json:"imageUrl"`
}

// Loadout model
type Loadout struct {
	ID      int       `json:"id"`
	UserID  string    `json:"userId"`
	Agent   string    `json:"agent"`
	Primary string    `json:"primary"`
	Sidearm string    `json:"sidearm"`
	Created time.Time `json:"created"`
}

// --- Static seed data ---
var agents = []Agent{
	{ID: "jett", Name: "Jett", Role: "Duelist", Description: "Jett's agile and evasive fighting style lets her take risks no one else can.", Abilities: []string{"Updraft", "Tailwind", "Cloudburst", "Blade Storm"}, ImageURL: "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt5ebf40a2dfaffb4e/5f21297f5f0cb0629a5bfcb9/V_AGENTS_587x900_Jett.png"},
	{ID: "sova", Name: "Sova", Role: "Initiator", Description: "Sova tracks, finds, and eliminates enemies with ruthless efficiency.", Abilities: []string{"Owl Drone", "Shock Bolt", "Recon Bolt", "Hunter's Fury"}, ImageURL: "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt181ad63adc9976a4/5f2129b2e0999b628bc8eb4e/V_AGENTS_587x900_Sova.png"},
	{ID: "sage", Name: "Sage", Role: "Sentinel", Description: "Sage creates safety for herself and her team wherever she goes.", Abilities: []string{"Barrier Orb", "Slow Orb", "Healing Orb", "Resurrection"}, ImageURL: "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt2a1c7b18aa5b1a6b/5f21297f078a8b626859f4a8/V_AGENTS_587x900_Sage.png"},
	{ID: "omen", Name: "Omen", Role: "Controller", Description: "Omen hunts in the shadows. He renders enemies blind, teleports across the field.", Abilities: []string{"Shrouded Step", "Paranoia", "Dark Cover", "From the Shadows"}, ImageURL: "https://images.contentstack.io/v3/assets/bltb6530b271fddd0b1/blt94dd043bce7fc9f2/5f21297f2ef66062fb6aa96c/V_AGENTS_587x900_Omen.png"},
}

var weapons = []Weapon{
	{ID: "classic", Name: "Classic", Type: "Sidearm", Cost: 0, Damage: 78, Accuracy: 85},
	{ID: "sheriff", Name: "Sheriff", Type: "Sidearm", Cost: 800, Damage: 159, Accuracy: 79},
	{ID: "spectre", Name: "Spectre", Type: "Primary", Cost: 1600, Damage: 78, Accuracy: 74},
	{ID: "vandal", Name: "Vandal", Type: "Primary", Cost: 2900, Damage: 160, Accuracy: 73},
	{ID: "phantom", Name: "Phantom", Type: "Primary", Cost: 2900, Damage: 156, Accuracy: 79},
	{ID: "operator", Name: "Operator", Type: "Primary", Cost: 4700, Damage: 255, Accuracy: 76},
}

// --- APIs ---

//encore:api public method=GET path=/agents
func GetAgents(ctx context.Context, params *GetAgentsParams) (*GetAgentsResponse, error) {
	var filtered []Agent
	for _, a := range agents {
		if params.Role != "" && a.Role != params.Role {
			continue
		}
		if params.Search != "" {
			s := strings.ToLower(params.Search)
			if !strings.Contains(strings.ToLower(a.Name), s) && !strings.Contains(strings.ToLower(a.Description), s) {
				continue
			}
		}
		filtered = append(filtered, a)
	}
	return &GetAgentsResponse{Agents: filtered, Total: len(filtered)}, nil
}

type GetAgentsParams struct {
	Role   string `query:"role"`
	Search string `query:"search"`
}
type GetAgentsResponse struct {
	Agents []Agent `json:"agents"`
	Total  int     `json:"total"`
}

//encore:api public method=GET path=/weapons
func GetWeapons(ctx context.Context, params *GetWeaponsParams) (*GetWeaponsResponse, error) {
	var out []Weapon
	for _, w := range weapons {
		if params.Type != "" && w.Type != params.Type {
			continue
		}
		if params.MaxCost > 0 && w.Cost > params.MaxCost {
			continue
		}
		if params.Search != "" {
			s := strings.ToLower(params.Search)
			if !strings.Contains(strings.ToLower(w.Name), s) {
				continue
			}
		}
		out = append(out, w)
	}
	return &GetWeaponsResponse{Weapons: out, Total: len(out)}, nil
}

type GetWeaponsParams struct {
	Type    string `query:"type"`
	MaxCost int    `query:"maxCost"`
	Search  string `query:"search"`
}
type GetWeaponsResponse struct {
	Weapons []Weapon `json:"weapons"`
	Total   int      `json:"total"`
}

//encore:api auth method=POST path=/loadouts
func CreateLoadout(ctx context.Context, req *CreateLoadoutRequest) (*CreateLoadoutResponse, error) {
	userID, _ := auth.UserID()
	id, err := createLoadout(ctx, userID, req.Agent, req.Primary, req.Sidearm)
	if err != nil {
		return nil, fmt.Errorf("failed to create loadout: %w", err)
	}
	return &CreateLoadoutResponse{ID: id, Message: "Loadout saved successfully"}, nil
}

type CreateLoadoutRequest struct {
	Agent   string `json:"agent"`
	Primary string `json:"primary"`
	Sidearm string `json:"sidearm"`
}
type CreateLoadoutResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

//encore:api auth method=GET path=/loadouts
func GetUserLoadouts(ctx context.Context) (*GetLoadoutsResponse, error) {
	userID, _ := auth.UserID()
	list, err := listLoadouts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loadouts: %w", err)
	}
	return &GetLoadoutsResponse{Loadouts: list, Total: len(list)}, nil
}
type GetLoadoutsResponse struct {
	Loadouts []Loadout `json:"loadouts"`
	Total    int       `json:"total"`
}

//encore:api public method=GET path=/health
func HealthCheck(ctx context.Context) (*HealthResponse, error) {
	return &HealthResponse{Status: "healthy", Message: "Valorant1 API is running", Timestamp: time.Now(), Version: "1.0.0"}, nil
}
type HealthResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

//encore:api public method=GET path=/stats
func GetStats(ctx context.Context) (*StatsResponse, error) {
	total, err := countLoadouts(ctx)
	if err != nil {
		total = 0
	}
	return &StatsResponse{TotalAgents: len(agents), TotalWeapons: len(weapons), TotalLoadouts: total, PopularAgent: "Jett"}, nil
}
type StatsResponse struct {
	TotalAgents   int    `json:"totalAgents"`
	TotalWeapons  int    `json:"totalWeapons"`
	TotalLoadouts int    `json:"totalLoadouts"`
	PopularAgent  string `json:"popularAgent"`
}

