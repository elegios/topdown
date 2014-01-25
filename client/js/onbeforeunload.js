function onbeforeunload() {
  var args = arguments
  window.addEventListener("beforeunload", function(e) {
    Array.forEach(args, function(arg) {
      if (typeof arg.onbeforeunload === "function")
        arg.onbeforeunload(e)
    })
  })
  console.log("onbeforeunload has been set")
}