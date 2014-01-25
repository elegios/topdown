  var model = mapmodel()
  var render = renderer(model)
  var send = network(
    "ws://"+location.host+"/ws",
    {
      maps: model.tickUpdate,
      type: model.itemUpdate,
    },
    onopen,
    function() {console.log("error ocurred")}
    )
  var inp = input(send, model, render)

  function onopen() {
    console.log("connection opened")
    var name = prompt("Enter name for new character or cancel.", "")
    if (!name)
      return

    send({
      command: "create",
      name: name,
    })
  }

  onbeforeunload(model, render, send, inp)


  render.mapname = "testmap"
  render.center.x = 0
  render.center.y = 0
