var CLIENT_AUTHOR = Date.now() +""+ Math.random();
var SERVER_COPY = "";


var conn = new WebSocket(
  "ws://" + window.location.hostname + ":" + window.location.port + "/ws/"
);
conn.onopen = function(e) {
  console.log("socket opened");
}
conn.onclose = function(e) {
  console.log("socket closed");
}
conn.onmessage = function(e) {
  var msg = JSON.parse(e.data);
  console.log("new msg", msg);
  switch (msg.t) {
    case "p":
      document.dispatchEvent(new CustomEvent("newPatches",{detail:{
        patches: msg.p.p,
        author: msg.p.a,
      }}));
      break;
    case "s":
      document.dispatchEvent(new CustomEvent("newSnapshot",{detail:{
        text: msg.p.t,
        }}));
      break;
    default:
      console.log("unknown message type", msg);
  }
};

document.addEventListener("newSnapshot", function(event){
  document.getElementById("input").value = event.detail.text;;
});

document.addEventListener("newPatches", function(event){
  var patches = event.detail.patches;
  var input = document.getElementById("input");
  var cursorStart = input.selectionStart;
  var cursorEnd = input.selectionEnd;
  for (var i = 0; i < patches.length; i++) {
    var p = patches[i];
    SERVER_COPY = applyPatch(SERVER_COPY,p);
    document.getElementById("input").value = applyPatch(document.getElementById("input").value,p);
    if (Math.max(p["1"],p["2"]) < Math.min(cursorStart,cursorEnd)) {
      input.selectionStart = cursorStart - Math.abs(p["2"] - p["1"]) + p["n"].length;
      input.selectionEnd = cursorEnd -  Math.abs(p["2"] - p["1"]) + p["n"].length;
    } else if (Math.min(p["1"],p["2"]) > Math.max(cursorStart,cursorEnd)) {
      input.selectionStart = cursorStart;
      input.selectionEnd = cursorEnd;
    } else {
      input.selectionStart = Math.min(p["1"],p["2"]);
      input.selectionEnd = input.selectionStart;
    }
  }
});

var patches = [];
sendPatches = function() {
  if (patches.length === 0) {
    return
  }
  conn.send(JSON.stringify({
    t: "p",
    d: new Date().toJSON(),
    p: {a: CLIENT_AUTHOR, p: patches},
  }));
  patches = [];
};

applyPatch = function(s,p) {
  return s.substr(0,p["1"]) + p["n"] + s.substr(p["2"]);
}

window.onload = function() {
  var input = document.getElementById("input");
  var timeoutHandle;
  input.onkeydown = function(event) {
    if (event.keyCode === 8) { // Backspace
      if (input.selectionStart === input.selectionEnd)
        patches.push({"1":input.selectionStart-1,"2":input.selectionEnd,"s":""});
      else
        patches.push({"1":input.selectionStart,"2":input.selectionEnd,"s":""});
      event.preventDefault();
    }
    if (event.keyCode === 46) // Delete
      event.preventDefault();
    window.clearTimeout(timeoutHandle);
    timeoutHandle = window.setTimeout(sendPatches,500);
  }
  input.onkeypress = function(event) {
    var char = String.fromCharCode(event.which || event.keyCode || event.charCode);
    if (char === 0) { return }
    if (patches.length > 0
        && patches[patches.length-1]["1"] === Math.min(input.selectionStart,input.selectionEnd)
        && patches[patches.length-1]["2"] === Math.max(input.selectionStart,input.selectionEnd))
        patches[patches.length-1]["n"] +=char;
    else
      patches.push({
        "1":Math.min(input.selectionStart,input.selectionEnd),
        "2":Math.max(input.selectionStart,input.selectionEnd),
        "s":char});
    event.preventDefault();
    window.clearTimeout(timeoutHandle);
    timeoutHandle = window.setTimeout(sendPatches,500);
  };
};
