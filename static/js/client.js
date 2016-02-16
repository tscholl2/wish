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
   switch (msg.t) {
      case "p":
         document.dispatchEvent(new Event("newPatch",{detail:{
            patches: msg.p.p,
            author: msg.p.a,
         }}));
         break;
      case "t":
         document.dispatchEvent(new Event("newSnapshot",{detail:{
            text: msg.p.p.t,
         }}));
         break;
      default:
         console.log("unknown message type", msg);
   }
};

document.addEventListener("newSnapshot", function(event){
   newText = event.details.text;
   if (snapShots.push(newText) > 10) {
      snapShots = snapShots.slice(-10);
   }
   document.getElementById("input").value = newText;
});

document.addEventListener("newPatch", function(event){
   patches = event.details.patches;
   author = event.details.a;
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
         conn.send(JSON.stringify({
            t: "p",
            p: {
               a: CLIENT_AUTHOR,
               d: new Date().toJSON(),
               p: patches,
            },
         }));
      }, 500);
   };
};
