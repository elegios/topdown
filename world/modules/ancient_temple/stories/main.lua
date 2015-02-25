-- Various constants
QUESTNAME = "The Opening of the Temple"
LOCKED_VARIATION = 4
UNLOCKED_VARIATION = 5
WALL_VARIATION = 6
POS = {
  {{ map = "othertest", x = 95, y = 104 }, { map = "othertest", x = 96, y = 81 }},
  {{ map = "othertest", x = 95, y = 108 }, { map = "othertest", x = 96, y = 83 }},
  {{ map = "othertest", x = 104, y = 104 }, { map = "othertest", x = 103, y = 81 }},
  {{ map = "othertest", x = 104, y = 108 }, { map = "othertest", x = 103, y = 83 }},
}
WALL_POS = {
  map = "othertest",
  x = 98,
  y = 79,
}
GATES_OPENED_PARTIAL = "gates_opened"
CONTINUATION = "gates_opened"

-- Variables
data = {
  {open = false, key = "key1open"},
  {open = false, key = "key2open"},
  {open = false, key = "key3open"},
  {open = false, key = "key4open"},
}

-- Will be called with a bool as its only argument,
-- true if it is the first run, false if it is any
-- subsequent run, in which case the function should
-- perform some sort of recovery
function main(firstRun)
  -- Determines whether story will end when execution reaches end of main function
  -- Defaults to true
  auto_end(false)

  if not firstRun then
    for i,_ in ipairs(data) do
      data[i].open = retrieve_value(data[i].key) or false
    end
  end

  for i,open in ipairs(data) do
    if not open then
      -- Creates, removes or replaces a prop
      prop({
        variation = LOCKED_VARIATION,
        collide = true,
        effect = keyeffect(i),
      }, POS[i][1])
      prop({
        variation = LOCKED_VARIATION,
        collide = true,
      }, POS[i][2])
    end
  end
end

function keyeffect(index)
  return function(character)
    data[index].open = true
    store_value(data[index].key, true)
    prop({
      variation = UNLOCKED_VARIATION,
      collide = true,
    }, POS[index][1])
    prop({
      variation = UNLOCKED_VARIATION,
      collide = true,
    }, POS[index][2])
    checkEnd(character)
  end
end

function checkEnd(character)
  local count = 0
  for _,open in ipairs(data) do
    if open then
      count = count + 1
    end
  end
  if count == 4 then
    -- Creates an announcement that will be saved in the
    -- history of the world
    announce(QUESTNAME, character.Name .. " has broken the last seal.")
  elseif count == 3 then
    announce(QUESTNAME, character.Name .. " has broken the third seal.")
  elseif count == 2 then
    announce(QUESTNAME, character.Name .. " has broken the second seal.")
  elseif count == 1 then
    announce(QUESTNAME, character.Name .. " has broken the first seal.")
  end
  if count < table.getn(POS) then
    return
  end

  -- Puts the map in <module-dir>/partials/_<argument2>.png on top of
  -- the already existing map, with the upper left corner on
  -- the position specified by the first argument
  apply_partial(GATES_OPENED_PARTIAL, WALL_POS)
  announce(QUESTNAME, "The gates to the temple have been opened.")
  -- Starts a substory residing in <module-dir>/stories/<argument>.lua
  -- Can contain slashes to denote subdirs
  substory(CONTINUATION)
  -- Ends execution of the current story.
  -- No callbacks can be registered in the current story, thus
  -- all props must be removed or have no effects
  end_story()
end
