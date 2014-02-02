function mapmodel() {
  var COLLIDE = 1 << 8
  var SETID = (1 << 4) - 1
  var VARIATION = ((1 << 4) - 1) << 4

  var data = {maps: {}, visible: [], blueprints: {}, selectedId: ""}

  if (localStorage.maps) {
    data.maps = JSON.parse(localStorage.maps)
    data.blueprints = JSON.parse(localStorage.blueprints)
  }

  data.onbeforeunload = function(e) {
    localStorage.maps = JSON.stringify(data.maps)
    localStorage.blueprints = JSON.stringify(data.blueprints)
  }

  function setpos(mapname, pos) {
    if (!data.maps[mapname])
      data.maps[mapname] = []

    if (!data.maps[mapname][pos.y])
      data.maps[mapname][pos.y] = []

    data.maps[mapname][pos.y][pos.x] = {
      collide: (pos.data & COLLIDE) != 0,
      setid: pos.data & SETID,
      variation: (pos.data & VARIATION) >> 4,
    }
    data.visible.push({x: pos.x, y: pos.y})
  }

  function ensure(id, send) {
    if (data.blueprints[id] || !id)
      return

    send({
      command: "itemrequest",
      id: id,
    })
  }

  data.tickUpdate = function(update, send) {
    clearVisible()
    for (var mapname in update.maps) {
      if (!update.maps.hasOwnProperty(mapname))
        return

      update.maps[mapname].forEach(function(pos) {
        setpos(mapname, pos)
      })
    }

    data.controllable = update.controllable
    data.characters = update.characters
    data.props = update.props
    data.items = update.items
    data.updates = update.updates

    data.characters.forEach(function(character) {
      ensure(character.weapon, send)
      ensure(character.armor, send)
      if (!character.inventory)
        return

      character.inventory.forEach(function(item) {
        ensure(item.id, send)
      })
    })
    data.items.forEach(function(item) {
      ensure(item.id, send)
    })
    if (!data.currChar())
      data.selectedId = data.controllable[0]
  }

  data.itemUpdate = function(update, net) {
    data.blueprints[update.id] = update
  }

  function clearVisible() {
    data.visible = []
  }

  data.currChar = function() {
    if (!data.characters || !data.selectedId)
      return null
    for (var i = 0; i < data.characters.length; i++)
      if (data.characters[i].id === data.selectedId)
        return data.characters[i]
    console.log(data.characters)
    console.log(data.controllable)
    console.log(data.selectedId)
  }

  return data
}