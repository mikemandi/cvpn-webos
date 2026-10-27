/* Real inline TV Launcher JS behavior, without external apps or packages.
   Exercise actual translation, Homebrew timeouts, resume/poll, QR state. */
'use strict';
const fs=require('fs'),path=require('path'),vm=require('vm'),assert=require('assert');
const html=fs.readFileSync(path.join(__dirname,'../launcher/index.html'),'utf8');
const source=html.split('<script>')[1].split('</script>')[0];

function fixture(responses,language){
 const elements={},calls=[],bridges=[],timers=[],handlers={},scripts=[];
 const steps=Array.from({length:3},()=>({
   className:'setup-step',
   mark:{textContent:''},name:{textContent:''},
   querySelector(selector){
     if(selector==='.mark')return this.mark;
     if(selector==='.step-name')return this.name;
     throw new Error('unexpected selector '+selector);
   }
 }));
 function el(id) {
   if(!elements[id]){
     elements[id]={
       hidden:true,textContent:'',className:'',firstChild:{nodeValue:''},
       attributes:{},setAttribute(k,v){this.attributes[k]=v;},
       focus(){this.focused=true},
       getContext(){return{fillStyle:'',fillRect(){}}}
     };
   }
   return elements[id];
 }
 const document={
   documentElement:{lang:'ru'},getElementById:el,
   querySelectorAll(){return steps;},
   addEventListener(name,fn){handlers[name]=fn;},
   createElement(){return{src:'',onerror:null};},
   head:{appendChild(node){scripts.push(node)}},
 };
 let nextId=1,closed=false;
 function timeout(fn,ms){const item={id:nextId++,fn,ms,cancelled:false};timers.push(item);return item.id;}
 function clear(id){timers.forEach(t=>{if(t.id===id)t.cancelled=true});}
 function fire(ms) {
   const t=timers.find(t=>!t.cancelled&&t.ms===ms);
   assert(t,'No active '+ms+'ms timer');
   t.cancelled=true;t.fn();
 }
 class PalmServiceBridge {
   constructor(){bridges.push(this)}
   call(url,params){
     assert.strictEqual(url,'luna://org.webosbrew.hbchannel.service/exec');
     const command=JSON.parse(params).command;
     assert(/payload\/install\.sh (status|install)$/.test(command));
     calls.push(command);
     const value=responses.shift();
     if(value!==null&&value!==undefined)this.onservicecallback(JSON.stringify(value));
   }
 }
 const window={
   PalmServiceBridge,
   close(){closed=true;},
   webOSSystem:{launchParams(){return '{}';}},
   CVPNQRCode:{Level:{M:1},QRCode:function(){
     this.addData=function(){};this.make=function(){};
     this.getModuleCount=function(){return 21};
     this.isDark=function(){return false};
   }}
 };
 vm.runInNewContext(source,{document,window,Date,JSON,Math,
   navigator:{language:language||'ru-RU'},
   Image:function(){this.src='';},
   setTimeout:timeout,clearTimeout:clear,
   setInterval(){return 99},clearInterval(){}},{timeout:1000});
 return {
   el,steps,calls,bridges,scripts,fire,
   get closed(){return closed},
   launch(mode){handlers.webOSLaunch({detail:{mode}})},
   backendLanguage(pref,resolved){window.cvpnLanguageInfo({localePreference:pref,resolvedTVLocale:resolved})},
   sendQR(data){window.cvpnQRInfo(data)},
   pageLanguage(){return document.documentElement.lang},
 };
}
const yes={returnValue:true,stdoutString:'BUNDLE_READY=YES\n'};
const no={returnValue:true,stdoutString:'BUNDLE_READY=NO\n'};
function run(title,fn){try{fn();process.stdout.write('ok '+title+'\n');}
 catch(e){process.stderr.write('FAIL '+title+'\n'+e.stack+'\n');process.exitCode=1;}}

run('already prepared: one quick status; no reinstall, auto close',()=>{
 const f=fixture([yes]);f.fire(250);
 assert.strictEqual(f.calls.length,1);
 assert(f.el('setup-title').textContent.includes('готов'));
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(650);assert.strictEqual(f.el('setup-screen').hidden,true);
});
run('clean install: only status + install and auto ready',()=>{
 const f=fixture([no,yes]);f.fire(250);
 assert.strictEqual(f.calls.length,2);
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(650);assert.strictEqual(f.el('setup-screen').hidden,true);
});
run('English auto: non-Russian TV translates every stage and QR',()=>{
 const f=fixture([no,yes],'en-US');f.fire(250);
 assert.strictEqual(f.pageLanguage(),'en');
 assert.strictEqual(f.el('setup-title').textContent,'CVPN is ready');
 assert.strictEqual(f.steps[0].name.textContent,'Checking device');
 f.launch('qr');
 assert.strictEqual(f.el('qr-instruction').textContent,'Point your phone camera at the QR code');
 assert.strictEqual(f.el('close-qr').textContent,'Close');
 f.sendQR({ok:true,code:'123456',address:'http://192.168.1.2:8765',
   pairURL:'http://192.168.1.2:8765/#pair=123456',
   expires:new Date(Date.now()+94000).toISOString(),localePreference:'auto',resolvedTVLocale:'en'});
 assert(f.el('qr-hint').textContent.includes('Code expires in:'));
 assert(f.el('qr-hint').textContent.includes('Home network only'));
});
run('Manual language preference takes priority over TV language',()=>{
 const f=fixture([no,yes],'ru-RU');f.fire(250);
 f.backendLanguage('en','ru');
 assert.strictEqual(f.pageLanguage(),'en');
 assert.strictEqual(f.el('setup-title').textContent,'CVPN is ready');
 f.launch('qr');
 f.sendQR({ok:false,localePreference:'ru',resolvedTVLocale:'en'});
 assert.strictEqual(f.pageLanguage(),'ru');
 assert(f.el('qr-hint').textContent.includes('Код истёк'));
});
run('English denied Homebrew: bounded, error with useful controls',()=>{
 const f=fixture([{returnValue:false,error:'Permission denied'}],'fr-FR');
 f.fire(250);
 assert(f.el('setup-status').textContent.includes('Homebrew'));
 assert.strictEqual(f.el('setup-retry').hidden,false);
 assert.strictEqual(f.el('setup-close').hidden,false);
 assert.strictEqual(f.el('setup-retry').textContent,'Retry');
});
run('Slow Homebrew status is retried read-only without confusing user',()=>{
 const f=fixture([null,yes]);f.fire(250);f.fire(15000);
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(1200);
 assert.strictEqual(f.calls.length,2);
 assert(f.calls.every(c=>/ status$/.test(c)));
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(650);assert.strictEqual(f.el('setup-screen').hidden,true);
});
run('Persistent silent status eventually shows a real error, no install',()=>{
 const f=fixture([null,null,null]);f.fire(250);
 f.fire(15000);f.fire(1200);f.fire(15000);f.fire(1200);f.fire(15000);
 assert.strictEqual(f.calls.length,3);
 assert.strictEqual(f.el('setup-retry').hidden,false);
 assert(f.calls.every(c=>/ status$/.test(c)));
});
run('Install callback timeout: background status polling, no Continue button',()=>{
 const f=fixture([no,null,yes]);f.fire(250);
 assert.strictEqual(f.calls.length,2);
 f.fire(90000);
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(2400);
 assert.strictEqual(f.calls.length,3);
 assert(/ status$/.test(f.calls[2]));
 assert.strictEqual(f.el('setup-title').textContent,'CVPN готов к работе');
 f.fire(650);assert.strictEqual(f.el('setup-screen').hidden,true);
});
run('In-progress install from previous launch: status-only polling',()=>{
 const f=fixture([{returnValue:true,stdoutString:'BUNDLE_READY=NO\nINSTALL_IN_PROGRESS=YES\n'},yes]);
 f.fire(250);assert.strictEqual(f.calls.length,1);
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(2400);assert.strictEqual(f.calls.length,2);
 assert(f.calls.every(c=>/ status$/.test(c)));
});
run('Stalled installer: bounded poll, then Retry on actual interruption',()=>{
 const responses=[no,null,no,no,no];
 const f=fixture(responses);f.fire(250);f.fire(90000);
 assert.strictEqual(f.el('setup-retry').hidden,true);
 f.fire(2400);f.fire(2400);f.fire(2400);
 assert.strictEqual(f.el('setup-retry').hidden,false);
 assert.strictEqual(f.calls.filter(c=>/ install$/.test(c)).length,1);
});
run('VPN active on update: refuse replacement, no second install',()=>{
 const f=fixture([{returnValue:true,stdoutString:'BUNDLE_READY=NO\nUPGRADE_BLOCKED=VPN_RUNNING\n'}]);
 f.fire(250);
 assert.strictEqual(f.calls.length,1);
 assert(f.el('setup-status').textContent.includes('отключите VPN'));
});
run('Malformed root service reply: no unsafe installation',()=>{
 const f=fixture([{returnValue:true}]);f.fire(250);
 assert.strictEqual(f.calls.length,1);
 assert.strictEqual(f.el('setup-retry').hidden,false);
});
run('Opening QR while installing discards stale installer callbacks',()=>{
 const f=fixture([null]);f.fire(250);f.launch('qr');
 f.bridges[0].onservicecallback(JSON.stringify(yes));
 assert.strictEqual(f.el('qr-screen').hidden,false);
 assert.strictEqual(f.el('setup-screen').hidden,true);
});
if(process.exitCode)process.exit(process.exitCode);
