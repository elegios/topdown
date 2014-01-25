function network(url, handlers, onopen, onerror) {
  var sock = new WebSocket(url)

  sock.onopen = onopen
  sock.onerror = onerror

  sock.onmessage = function(mess) {
    var obj = JSON.parse(mess.data)
    for (var key in handlers) {
      if (!handlers.hasOwnProperty(key))
        continue

      if (!obj[key])
        continue

      if (handlers[key](obj, send))
        return
    }
  }

  function send(obj) {
    sock.send(JSON.stringify(obj))
  }

  return send
}