console.log("hello");

CLIENT_ID = Date.now() +""+ Math.random()

window.onload = function() {
   var timeoutHandle;
   snapShots = [""];
   dmp = new diff_match_patch();
   conn = new WebSocket(
      "ws://" + window.location.hostname + ":" + window.location.port + "/ws/"
   );
   conn.onmessage = function(e) {
      var msg = JSON.parse(e.data);
      console.log("got msg", msg);
      var newText = dmp.patch_apply(msg.p, snapShots[snapShots.length-1])[0];
      if (snapShots.push(newText) > 10) {
         snapShots = snapShots.slice(-10);
      }
      if (msg.id != CLIENT_ID) {
         document.getElementById("input").value = newText;
      }
   };
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
            id: CLIENT_ID,
            t: "p",
            p: patches,
         }));
      },1000);
   };
};
