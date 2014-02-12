function renderer(model) {
  var FOG_ALPHA = 0.5
  var SCALE = 1
  var TILEX = 8; var TILEY = 8

  var canvas = document.querySelector(".view")
  var context = canvas.getContext("2d")
  var tileset = document.querySelector(".tileset")

  canvas.width = document.body.clientWidth/SCALE
  canvas.height = document.body.clientHeight/SCALE

  var data = {}
  data.center = {x: 0, y: 0}
  data.mapname = ""

  // semi global helpers, set at the start of drawmap
  var xoff; var yoff
  var map

  var sets = [blue, cave, grass, water, bridge]

  function drawmap() {
    context.clearRect(0, 0, canvas.width, canvas.height)
    data.center.x = model.currChar().x
    data.center.y = model.currChar().y
    data.mapname = model.currChar().mapname
    xoff = Math.round(data.center.x*TILEX - canvas.width/2)
    yoff = Math.round(data.center.y*TILEY - canvas.height/2)

    if (!data.mapname)
      return
    map = model.maps[data.mapname]

    drawTiles()

    drawProps()

    drawItems()

    drawCharacters()

    drawFog()
  }

  function collides(x, y) {
    return map[y] && map[y][x] && map[y][x].collide
  }

  function drawTiles() {
    var maxY = Math.ceil((yoff + canvas.height)/8)
    var maxX = Math.ceil((xoff + canvas.width)/8)
    for (var j = Math.floor(yoff/8); j < maxY; j++) {
      if (!map[j])
        continue

      for (var i = Math.floor(xoff/8); i < maxX; i++) {
        if (!map[j][i])
          continue

        drawTile(i, j, map[j][i].setid, map[j][i].collide, map[j][i].variation)
      }
    }
  }

  function drawCharacters() {
    model.characters.forEach(function(character) {
      if (character.mapname !== data.mapname)
        return

      var tile = {x: character.variation % 5, y: 8}
      drawFromTileset(tile, character)
    })
  }

  function drawProps() {
    if (!model.props[data.mapname])
      return

    model.props[data.mapname].forEach(function(prop) {
      var tile = (function() {
        var v = prop.variation % 25
        if (v < 15)
          return {x: v % 5, y: Math.floor(v/5) + 6}
        return {x: (v-15) % 6, y: Math.floor((v-15)/6) + 13}
      }())
      drawFromTileset(tile, prop)
    })
  }

  function drawItems() {
    if (!model.items[data.mapname])
      return

    model.items[data.mapname].forEach(function(itempos) {
      if (!model.blueprints[itempos.id]) {
        drawFromTileset({x: 5, y: 6}, itempos)
        return
      }

      switch (model.blueprints[itempos.id].type) {
      case "potion":
        drawFromTileset({x: 0, y: 6}, itempos)
        return

      default:
        drawFromTileset({x: 5, y:6}, itempos)
        return
      }
    })
  }

  function drawFog() {
    context.save()
    context.beginPath()
    context.rect(0, 0, canvas.width, canvas.height)
    model.visible.forEach(function(coord) {
      context.rect(coord.x*TILEX - xoff + TILEX, coord.y*TILEY - yoff, -TILEX, TILEY)
    })
    context.clip()
    context.globalAlpha = FOG_ALPHA
    context.fillStyle = "black"
    context.fillRect(0, 0, canvas.width, canvas.height)
    context.restore()
  }

  function drawTile(x, y, setid, wall, variation) {
    var tile = sets[setid](x, y, wall, variation)
    drawFromTileset(tile, {x: x, y: y})
  }

  function blue(x, y, wall, variation) {
    if (!wall)
      return variation ? {x: 5, y: 3} : {x: 0, y: 0}

    if (collides(x, y+1))
      return {x: 0, y: 3}

    return {x: 1+(variation%4), y: 3}
  }

  function cave(x, y, wall, variation) {
    if (!wall)
      return {x: variation ? 1 : 2, y: 4}

    return {x: collides(x, y+1) ? 0 : 1, y: 1}
  }

  function grass(x, y, wall, variation) {
    if (wall)
      return {x: 4, y: 2}

    variation %= 6
    if (variation == 4)
      return {x: 5, y: 2}
    if (variation == 5)
      return {x: 0, y: 4}

    return {x: variation, y: 2}
  }

  function water(x, y, wall, variation) {
    return [
      {x: 3, y: 5},
      {x: 0, y: 5},
      {x: 1, y: 5}
    ][variation % 3]
  }

  function bridge(x, y, wall, variation) {
    if (variation)
      return {x: 3, y: 1}
    return {x: 2, y: 1}
  }

  function drawFromTileset(tile, pos) {
    context.drawImage(tileset,
      tile.x*TILEX, tile.y*TILEY,
      TILEX, TILEY,
      pos.x*TILEX - xoff, pos.y*TILEY - yoff,
      TILEX, TILEY)
  }

  setInterval(drawmap, 100)

  return data
}