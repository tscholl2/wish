console.log("hi");

window.onload = function() {
   //conn = new WebSocket("ws://" window.location.hostname + ":" + window.location.port + "/ws");
   conn = new WebSocket("ws://localhost:8080/ws/");
   conn.onmessage = function(e) {
      document.getElementById("output").innerHTML += "<p>"+e.data+"</p>";
   };
   document.getElementById("input").onkeydown = function(event) {
      if (event.keyCode == 13)
         conn.send(document.getElementById("input").value);
   };
};
