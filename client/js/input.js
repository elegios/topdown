function input(send, model, render) {
  var data = {}

  function currentCharacterId() {
    return model.selectedId
  }

  function keypress(e) {
    switch (e.keyCode) {
      case 32:
        send({command: "tick"})
        return

      case 190:
      case 69:
        send({
          command: "pickup",
          character: currentCharacterId(),
        })
        return
    }

    var dir
    switch (e.keyCode) {
      case 37:
      case 65:
        dir = "left"
        break

      case 39:
      case 68:
        dir = "right"
        break

      case 38:
      case 87:
        dir = "up"
        break

      case 40:
      case 83:
        dir = "down"
        break
    }
    if (dir) {
      send({
        command: "move",
        character: currentCharacterId(),
        direction: dir,
      })
      send({command: "tick"})
      return
    }
    if (49 <= e.keyCode && e.keyCode <= 57) {
      var num = e.keyCode - 49
      if (num >= model.controllable.length)
        return
      model.selectedId = model.controllable[num]
    }
  }

  window.addEventListener("keydown", keypress, true)

  return data
}