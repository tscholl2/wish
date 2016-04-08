CLIENT_AUTHOR = Date.now() +""+ Math.random()

var snapShots = [""];
var dmp = new diff_match_patch();
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
   document.dispatchEvent(new CustomEvent("newPatch",{detail:{
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
  console.log("new snapshot", event);
  newText = event.detail.text;
  if (snapShots.push(newText) > 10) {
    snapShots = snapShots.slice(-10);
  }
  document.getElementById("input").value = newText;
});

document.addEventListener("newPatch", function(event){
  console.log("new patch", event);
  patches = event.detail.patches;
  author = event.detail.a;
  var newText = dmp.patch_apply(patches, snapShots[snapShots.length-1])[0];
  if (snapShots.push(newText) > 10) {
    snapShots = snapShots.slice(-10);
  }
  if (author != CLIENT_AUTHOR) {
    document.getElementById("input").value = newText;
  }
});

var timeoutHandle;
window.onload = function() {
  document.getElementById("input").onkeydown = function(event) {
    window.clearTimeout(timeoutHandle);
    timeoutHandle = window.setTimeout(function(){
      if (snapShots[snapShots.length-1] == document.getElementById("input").value) {
        return
      }
      var patches = dmp.patch_make(
        snapShots[snapShots.length-1],
        document.getElementById("input").value
      );
      console.log("sending patches: ", patches);
      conn.send(JSON.stringify({
        t: "p",
        d: new Date().toJSON(),
        p: {
          a: CLIENT_AUTHOR,
          p: patches,
        },
      }));
    }, 500);
  };
};
