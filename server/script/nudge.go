package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/helpers"
	"github.com/elegios/topdown/server/types"
)

const (
	nudgeHEALTH = iota
	nudgeMAX_HEALTH
	nudgeVIEW_DIST
	nudgeACTIONS
	nudgeRECOVERY_SPEED
	nudgeRECOVERY_MAX
)

func (v *vm) nudge(L *lua.State) int {
	var origin *types.Character
	if L.GetTop() == 4 {
		origin = L.ToGoStruct(1).(*types.Character)
	}
	target := L.ToGoStruct(2).(*types.Character)
	nudgeType := L.CheckString(3)
	amount := float32(L.CheckNumber(4))

	switch nudgeType {
	case "health":
		prev := target.Health
		target.Health = helpers.Min(target.Health+int(amount), target.MaxHealth)
		amount = float32(target.Health - prev)

	case "maxhealth":
		target.MaxHealth += int(amount)
		target.Health = helpers.Min(target.Health, target.MaxHealth)

	case "view distance":
		target.ViewDist += int(amount)

	case "actions":
		target.Actions += amount

	case "recovery speed":
		target.RecoverySpeed += amount

	case "recovery max":
		target.RecoveryMax += int(amount)
	}
	content := types.NudgeUpdate{
		Nudge:  nudgeType,
		Amount: amount,
		Target: target.Id,
	}
	if origin != nil {
		content.Character = origin.Id
	}
	p := target.Pos
	v.world.Updates = append(v.world.Updates, types.Update{
		Pos:     &p,
		Content: content,
	})

	return 0
}
