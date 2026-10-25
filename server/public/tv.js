(function(){'use strict';
var expiry=document.getElementById('expiry');
var seconds=180;
var reloading=false;
function renew(){if(!reloading){reloading=true;window.location.reload();}}
function tick(){seconds--;if(seconds<=0){renew();return;}expiry.textContent='Код действителен ещё '+seconds+' сек.';}
function check(){fetch('/api/pairing/tv-pending',{cache:'no-store'}).then(function(r){if(!r.ok)throw Error('status');return r.json();}).then(function(data){if(typeof data.remaining==='number'){seconds=Math.min(seconds,data.remaining);if(data.remaining<=0)renew();}}).catch(function(){});}
setInterval(tick,1000);
setInterval(check,2500);
})();
