/* Lightweight behavioral smoke tests of actual inline TV Launcher JavaScript.
   No npm packages, network, LG device, or privileged operations required. */
'use strict';
const fs=require('fs');
const path=require('path');
const vm=require('vm');
const assert=require('assert');
const html=fs.readFileSync(path.join(__dirname,'../launcher/index.html'),'utf8');
const source=html.split('<script>')[1].split('</script>')[0];

function fixture(responses) {
  const els={};
  const calls=[], bridges=[], timers=[], handlers={};
  const steps=Array.from({length:3},()=>({
    className:'setup-step',
    mark:{textContent:''},
    querySelector(x){ assert.strictEqual(x,'.mark'); return this.mark; }
  }));
  function element(id) {
    if(!els[id]){
      els[id]={hidden:true,textContent:'',className:'',
        setAttribute(k,v){ this[k]=v; },focus(){this.focused=true;},
        getContext(){return {fillStyle:'',fillRect(){}}}
      };
    }
    return els[id];
  }
  const document={
    getElementById:element,querySelectorAll(){return steps;},
    addEventListener(event,cb){handlers[event]=cb;},createElement(){return {src:'',onerror:null};},
    head:{appendChild(){}},
  };
  let closed=false, nextTimer=1;
  function setTimeoutFake(fn,ms){const id=nextTimer++;timers.push({id,fn,ms,cancelled:false});return id;}
  function clearTimeoutFake(id){timers.forEach(t=>{if(t.id===id)t.cancelled=true});}
  function fire(ms){
    const t=timers.find(t=>!t.cancelled&&t.ms===ms);
    assert(t,'Missing active '+ms+'ms timer');
    t.cancelled=true;t.fn();
  }
  class PalmServiceBridge {
    constructor(){bridges.push(this)}
    call(url,json){
      assert.strictEqual(url,'luna://org.webosbrew.hbchannel.service/exec');
      const data=JSON.parse(json);
      assert(/payload\/install\.sh (status|install)$/.test(data.command));
      calls.push(data.command);
      const value=responses.shift();
      if(value!==undefined&&value!==null) this.onservicecallback(JSON.stringify(value));
    }
  }
  const window={
    PalmServiceBridge,close(){closed=true},
    webOSSystem:{launchParams(){return '{}'}},
  };
  const context={
    document,window,Date,JSON,Math,Image:function(){this.src='';},
    setTimeout:setTimeoutFake,clearTimeout:clearTimeoutFake,
    clearInterval(){},setInterval(){return 40;},
  };
  vm.runInNewContext(source,context,{timeout:1000});
  return {element,steps,calls,bridges,fire,launch(mode){handlers.webOSLaunch({detail:{mode:mode}});},get closed(){return closed}};
}
const yes={returnValue:true,stdoutString:'BUNDLE_READY=YES\n'};
const no={returnValue:true,stdoutString:'BUNDLE_READY=NO\n'};

function run(title,fn){
  try{fn();process.stdout.write('ok: '+title+'\n');}
  catch(e){process.stderr.write('FAIL: '+title+'\n'+e.stack+'\n');process.exitCode=1;}
}
run('installed: check status only, display completed preparation',()=>{
  const f=fixture([yes]);f.fire(250);assert.strictEqual(f.calls.length,1);
  assert(f.element('setup-title').textContent.includes('готов'));
  f.fire(700);assert.strictEqual(f.element('setup-screen').hidden,true);
});
run('first launch: one status and one install then ready',()=>{
  const f=fixture([no,yes]);f.fire(250);
  assert.strictEqual(f.calls.length,2);
  assert(f.element('setup-title').textContent.includes('готов'));
  f.fire(700);assert.strictEqual(f.element('setup-screen').hidden,true);
});
run('denied root-service: visible error and retry controls',()=>{
  const f=fixture([{returnValue:false,error:'Denied method call exec'}]);f.fire(250);
  assert(f.element('setup-status').textContent.includes('Homebrew'));
  assert.strictEqual(f.element('setup-retry').hidden,false);
  assert.strictEqual(f.element('setup-close').hidden,false);
});
run('silent status: finite 20-second timeout',()=>{
  const f=fixture([null]);f.fire(250);f.fire(20000);
  assert(f.element('setup-status').textContent.includes('не ответил'));
  assert.strictEqual(f.element('setup-retry').hidden,false);
});
run('silent install: retry performs status only, never a second install',()=>{
  const f=fixture([no,null,{returnValue:true,stdoutString:'BUNDLE_READY=NO\nINSTALL_IN_PROGRESS=YES\n'}]);
  f.fire(250);assert.strictEqual(f.calls.length,2);
  f.fire(90000);assert.strictEqual(f.element('setup-retry').hidden,false);
  f.element('setup-retry').onclick();
  assert.strictEqual(f.calls.length,3);
  assert(/ status$/.test(f.calls[2]));
  assert(f.element('setup-status').textContent.includes('уже выполняется'));
});
run('upgrade while VPN runs: never call installer',()=>{
  const f=fixture([{returnValue:true,stdoutString:'BUNDLE_READY=NO\nUPGRADE_BLOCKED=VPN_RUNNING\n'}]);
  f.fire(250);assert.strictEqual(f.calls.length,1);
  assert(f.element('setup-status').textContent.includes('отключите VPN'));
});
run('malformed service reply: fail without starting installer',()=>{
  const f=fixture([{returnValue:true}]);
  f.fire(250);
  assert.strictEqual(f.calls.length,1);
  assert.strictEqual(f.element('setup-retry').hidden,false);
  assert.strictEqual(f.element('setup-screen').hidden,false);
});
run('QR relaunch invalidates obsolete status callback',()=>{
  const f=fixture([null]);f.fire(250);assert.strictEqual(f.bridges.length,1);
  f.launch('qr');
  assert.strictEqual(f.element('qr-screen').hidden,false);
  f.bridges[0].onservicecallback(JSON.stringify(yes));
  assert.strictEqual(f.element('qr-screen').hidden,false);
  assert.strictEqual(f.element('setup-screen').hidden,true);
});
if(process.exitCode) process.exit(process.exitCode);
