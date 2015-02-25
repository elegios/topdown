-- Constants
QUESTNAME = "The Opening of the Temple"
BRIDGE_PARTIAL = "bridge"
BRIDGE_POS = {
  { map = "othertest", x = 42, y = 42 },
  { map = "othertest", x = 42, y = 42 },
  { map = "othertest", x = 42, y = 42 },
}
TURRET_AI = "turret_ai"
TURRETS = {
  { pos = { map = "othertest", x = 42, y = 42 }, key = "turret1id" },
  { pos = { map = "othertest", x = 42, y = 42 }, key = "turret2id" },
  { pos = { map = "othertest", x = 42, y = 42 }, key = "turret3id" },
}
LOOT = {
  { pos = { map = "othertest", x = 42, y = 42 }, blueprint = "something" },
  { pos = { map = "othertest", x = 42, y = 42 }, blueprint = "something" },
  { pos = { map = "othertest", x = 42, y = 42 }, blueprint = "something" },
}

-- Variables
turretsAlive = 0

function main(firstRun)
  auto_end(false)
	announce("The Opening of the Temple", "The temple has been opened.")
  if firstRun then
    setup()
  else
    recover()
  end
end

function setup()
  -- setup turrets
  turretsAlive = table.getn(TURRETS) -- TODO: table
  for i,info in ipairs(TURRETS) do
    local turret = create_character(info.pos) -- TODO: what should be specified in this function?
    -- ai selection persists through recovery, but not its actual state, it
    -- restarts when the world is reloaded, thus the ai itself should be
    -- capable of returning to the same state easily
    set_ai(turret, TURRET_AI) -- TODO: supply
    on_death(turret, turretDeath(i)) -- TODO: supply
    store_value(info.key, turret.Id)
  end

  -- create loot
  for _,info in ipairs(LOOT) do
    item(info.blueprint, info.pos)
  end

  -- TODO: how to create a patrol with a list of points? Regardless, it should
  -- probably be ai internal, or at most a "fire story event" behavior (which
  -- would then require setting up in recover)

  -- TODO: create boss like creatures in inner chamber
end

function recover()
  for i,info in ipairs(TURRETS) do
    local turretId = retrieve_value(info.key)
    if turretId then
      turretsAlive = turretsAlive + 1
      on_death(turretId, turretDeath(i)) -- TODO: can take a character or an id
    end
  end
end

function turretDeath(index)
  return function(turret)
    apply_partial(BRIDGE_PARTIAL, BRIDGE_POS[index])
    -- Storing nil is the same as deleting the value
    store_value(turret.Id, nil)
    announce(QUESTNAME, "A bridge lowers over the water by the entrance.")
    checkEnd()
  end
end

function checkEnd()
  if turretsAlive > 0 then
    return
  end

  announce(QUESTNAME, "The way to the inner chamber is now open.")

  end_story()
end
