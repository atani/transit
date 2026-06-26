package main

import (
	"testing"

	"github.com/atani/transit/internal/transit"
)

func TestPlanIsWalkOnly(t *testing.T) {
	walkOnly := transit.PlanResponse{Journeys: []transit.Journey{
		{Legs: []transit.Leg{{Kind: "walk"}}},
	}}
	if !planIsWalkOnly(walkOnly) {
		t.Error("a walk-only itinerary should be detected as walk-only")
	}

	withTransit := transit.PlanResponse{Journeys: []transit.Journey{
		{Legs: []transit.Leg{{Kind: "walk"}, {Kind: "rail"}}},
	}}
	if planIsWalkOnly(withTransit) {
		t.Error("an itinerary with a rail leg is not walk-only")
	}

	noJourneys := transit.PlanResponse{}
	if !planIsWalkOnly(noJourneys) {
		t.Error("no journeys means no transit route was found")
	}
}
