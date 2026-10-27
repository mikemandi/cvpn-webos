#!/usr/bin/env node
// Deterministic screenshots of the real app with an empty, offline demo state.
// No user profiles, subscription links, credentials or live device IPs.
'use strict';
const http=require('node:http'),fs=require('node:fs'),path=require('node:path');
const assert=require('node:assert/strict');
const {chromium}=require('playwright');
const root=path.resolve(__dirname,'..');
const publicDir=path.join(root,'server/public');
const launcherDir=path.join(root,'launcher');
const outDir=path.join(root,'assets/screenshots');
const translations=JSON.parse(fs.readFileSync(path.join(publicDir,'translations.json'),'utf8'));
const types={'.html':'text/html; charset=utf-8','.css':'text/css; charset=utf-8','.js':'text/javascript; charset=utf-8','.png':'image/png'};
const server=http.createServer((req,res)=>{
 let u=new URL(req.url,'http://localhost').pathname;
 let f=u==='/tv-launcher'?path.join(launcherDir,'index.html'):u==='/qr-bundle.js'?path.join(launcherDir,'qr-bundle.js'):path.resolve(publicDir,u==='/'?'index.html':u.slice(1));
 if(!f.startsWith(publicDir+path.sep)&&!f.startsWith(launcherDir+path.sep)){res.writeHead(403);res.end();return;}
 if(!fs.existsSync(f)||!fs.statSync(f).isFile()){res.writeHead(404);res.end();return;}
 res.writeHead(200,{'Content-Type':types[path.extname(f)]||'application/octet-stream','Cache-Control':'no-store'});
 fs.createReadStream(f).pipe(res);
});
async function translate(page){
 await page.evaluate(dictionary=>{
  const pairs=Object.entries(dictionary).filter(([a,b])=>a&&b&&a!==b).sort((a,b)=>b[0].length-a[0].length);
  const walker=document.createTreeWalker(document.body,NodeFilter.SHOW_TEXT);
  let n;
  while((n=walker.nextNode())){
    if(!/[А-Яа-яЁё]/.test(n.nodeValue))continue;
    let s=n.nodeValue;
    for(const [a,b] of pairs)if(s.includes(a))s=s.split(a).join(b);
    n.nodeValue=s;
  }
  document.documentElement.lang='en';
  document.getElementById('settings-language').value='en';
  document.querySelector('#settings-language option[value=ru]').textContent='Russian';
  document.getElementById('settings-panel-address').textContent='—';
 },translations);
}
async function preparePhone(page,tab){
 await page.route('**/app.js*',r=>r.fulfill({status:200,contentType:'application/javascript',body:''}));
 await page.route('**/i18n.js*',r=>r.fulfill({status:200,contentType:'application/javascript',body:''}));
 await page.goto('/',{waitUntil:'networkidle'});
 await page.evaluate(name=>{
  const $=id=>document.getElementById(id);
  $('login').hidden=true;$('dashboard').hidden=false;
  $('panel-health').textContent='Local panel';
  $('connection-state').textContent='VPN disconnected';
  $('connection-detail').textContent='Choose a server to connect';
  $('profiles').innerHTML='<div class="empty-state"><strong>Add your first configuration</strong><p>No configurations yet. Add your own link or get started with CVPN.</p></div>';
  $('add-toggle').disabled=false;$('subs-add').disabled=false;
  $('home-empty-promo').hidden=true;$('subs-promo').hidden=true;
  for(const key of ['home','subscriptions','settings']){
   $('tab-'+key).hidden=(name!==key);
   $('nav-'+key).classList.toggle('active',name===key);
   if(name===key)$('nav-'+key).setAttribute('aria-current','page');
   else $('nav-'+key).removeAttribute('aria-current');
  }
  $('subscription-list').innerHTML='<div class="empty-state"><strong>No configurations yet</strong><p>Import your own HTTPS subscription, VLESS link or JSON configuration.</p></div>';
  window.scrollTo(0,0);
 },tab);
 await translate(page);
}
function assertEnglish(visible,screen){
 assert(!/[А-Яа-яЁё]/.test(visible),'Non-English visible text on '+screen+': '+(visible.match(/.{0,25}[А-Яа-яЁё].{0,25}/)||[])[0]);
}
async function phone(browser,origin,tab,file,modal){
 const ctx=await browser.newContext({viewport:{width:430,height:932},deviceScaleFactor:2,isMobile:true,hasTouch:true,locale:'en-US',colorScheme:'dark',reducedMotion:'reduce'});
 const page=await ctx.newPage();
 await preparePhone(page,tab);
 if(modal){
  await page.evaluate(()=>{document.getElementById('modal-overlay').hidden=false;document.getElementById('import-view').hidden=false;document.getElementById('import-content').value='';});
  await translate(page);
 }
 assertEnglish(await page.locator('body').innerText(),file);
 await page.screenshot({path:path.join(outDir,file),animations:'disabled'});
 await ctx.close();
}
async function tv(browser,origin,mode,file){
 const ctx=await browser.newContext({viewport:{width:1280,height:720},deviceScaleFactor:1,locale:'en-US',reducedMotion:'reduce'});
 const page=await ctx.newPage();
 await page.route('http://127.0.0.1:8765/**',r=>r.fulfill({status:200,contentType:'text/javascript',body:''}));
 await page.goto(origin+'/tv-launcher',{waitUntil:'domcontentloaded'});
 await page.evaluate(mode=>{
  document.body.style.setProperty('background','linear-gradient(115deg,#1b2941,#0a1527 72%)','important');
  window.cvpnLanguageInfo({localePreference:'en',resolvedTVLocale:'en'});
  if(mode==='setup'){
   document.getElementById('setup-screen').hidden=false;
   document.getElementById('qr-screen').hidden=true;
   document.getElementById('setup-title').textContent='Preparing CVPN';
   document.getElementById('setup-status').textContent='Checking the required components';
  } else {
   document.dispatchEvent(new CustomEvent('webOSLaunch',{detail:{mode:'qr'}}));
   // Example address reserved by RFC 5737, non-functional sample code.
   window.cvpnQRInfo({ok:true,code:'000000',address:'http://192.0.2.10:8765',
     pairURL:'http://192.0.2.10:8765/#pair=000000',
     localePreference:'en',resolvedTVLocale:'en',
     expires:new Date(Date.now()+180000).toISOString()});
  }
 },mode);
 assertEnglish(await page.locator('body').innerText(),file);
 await page.screenshot({path:path.join(outDir,file),animations:'disabled'});
 await ctx.close();
}
(async()=>{
 fs.mkdirSync(outDir,{recursive:true});
 await new Promise(resolve=>server.listen(0,'127.0.0.1',resolve));
 const origin='http://127.0.0.1:'+server.address().port;
 const browser=await chromium.launch({headless:true,args:['--no-sandbox']});
 try{
  await phone(browser,origin,'home','phone-home-empty.png',false);
  await phone(browser,origin,'subscriptions','phone-subscriptions-empty.png',false);
  await phone(browser,origin,'settings','phone-settings.png',false);
  await phone(browser,origin,'home','phone-import-empty.png',true);
  await tv(browser,origin,'setup','tv-setup.png');
  await tv(browser,origin,'qr','tv-pairing-demo.png');
  for(const file of fs.readdirSync(outDir))if(file.endsWith('.png'))console.log(file,fs.statSync(path.join(outDir,file)).size,'bytes');
 } finally {await browser.close();server.close();}
})().catch(e=>{console.error(e);server.close();process.exitCode=1;});
