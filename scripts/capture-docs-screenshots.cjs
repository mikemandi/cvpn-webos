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
const language=process.argv.includes('--lang=ru')?'ru':'en';
const isRU=language==='ru';
function screenshotName(name){return isRU?name.replace(/\.png$/,'-ru.png'):name;}
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
 // Use the app's OWN i18n.js, not a separate screenshot translation table.
 await page.evaluate(()=>{
  if(!window.CVPNLocale)throw new Error('Real app localization did not load');
  window.CVPNLocale.set(document.documentElement.dataset.docsLanguage||'en');
  window.CVPNLocale.refresh();
  // Closed language dropdown: use English for its inaccessible off-screen option.
  if(document.documentElement.dataset.docsLanguage==='en')document.querySelector('#settings-language option[value=ru]').textContent='Russian';
  document.getElementById('settings-panel-address').textContent='—';
 });
}
async function preparePhone(page,origin,tab){
 await page.route('**/app.js*',r=>r.fulfill({status:200,contentType:'application/javascript',body:''}));
 await page.addInitScript(lang=>{localStorage.setItem('cvpn_webos_language',lang);document.addEventListener('DOMContentLoaded',()=>{document.documentElement.dataset.docsLanguage=lang;});},language);
 await page.goto(origin+'/',{waitUntil:'networkidle'});
 await page.evaluate(({name,lang})=>{
  const $=id=>document.getElementById(id);
  $('login').hidden=true;$('dashboard').hidden=false;
  const ru=lang==='ru';
  $('panel-health').textContent=ru?'Локальная панель':'Local panel';
  $('connection-state').textContent=ru?'VPN отключён':'VPN disconnected';
  $('connection-detail').textContent=ru?'Выберите сервер для подключения':'Choose a server to connect';
  $('profiles').innerHTML='<div class="empty-state"><strong>'+(ru?'Добавьте свою первую конфигурацию':'Add your first configuration')+'</strong><p>'+(ru?'Конфигураций пока нет. Добавьте свою ссылку либо начните с CVPN.':'No configurations yet. Add your own link or get started with CVPN.')+'</p></div>';
  $('add-toggle').disabled=false;$('subs-add').disabled=false;
  $('home-empty-promo').hidden=true;$('subs-promo').hidden=true;
  for(const key of ['home','subscriptions','settings']){
   $('tab-'+key).hidden=(name!==key);
   $('nav-'+key).classList.toggle('active',name===key);
   if(name===key)$('nav-'+key).setAttribute('aria-current','page');
   else $('nav-'+key).removeAttribute('aria-current');
  }
  $('subscription-list').innerHTML='<div class="empty-state"><strong>'+(ru?'Конфигураций пока нет':'No configurations yet')+'</strong><p>'+(ru?'Добавьте свою ссылку подписки HTTPS, VLESS или JSON.':'Import your own HTTPS subscription, VLESS link or JSON configuration.')+'</p></div>';
  window.scrollTo(0,0);
 },{name:tab,lang:language});
 await translate(page);
}
function assertLanguage(visible,screen){
 if(isRU){assert(/[А-Яа-яЁё]/.test(visible),'No Russian text visible on '+screen);return;}
 assert(!/[А-Яа-яЁё]/.test(visible),'Non-English visible text on '+screen+': '+(visible.match(/.{0,25}[А-Яа-яЁё].{0,25}/)||[])[0]);
}
async function phone(browser,origin,tab,file,modal){
 const ctx=await browser.newContext({viewport:{width:430,height:932},deviceScaleFactor:2,isMobile:true,hasTouch:true,locale:'en-US',colorScheme:'dark',reducedMotion:'reduce'});
 const page=await ctx.newPage();
 await preparePhone(page,origin,tab);
 if(modal){
  await page.evaluate(()=>{document.getElementById('modal-overlay').hidden=false;document.getElementById('import-view').hidden=false;document.getElementById('import-content').value='';});
  await translate(page);
 }
 assertLanguage(await page.locator('body').innerText(),file);
 await page.screenshot({path:path.join(outDir,screenshotName(file)),animations:'disabled'});
 await ctx.close();
}
async function tv(browser,origin,mode,file){
 const ctx=await browser.newContext({viewport:{width:1280,height:720},deviceScaleFactor:1,locale:'en-US',reducedMotion:'reduce'});
 const page=await ctx.newPage();
 // Mock only the Homebrew callback boundary so genuine TV setup UI stays in its checking state.
 await page.addInitScript(()=>{window.PalmServiceBridge=function(){this.call=function(){};};});
 await page.route('http://127.0.0.1:8765/**',r=>r.fulfill({status:200,contentType:'text/javascript',body:''}));
 await page.goto(origin+'/tv-launcher',{waitUntil:'domcontentloaded'});
 await page.evaluate(({mode,lang})=>{
  document.body.style.setProperty('background','linear-gradient(115deg,#1b2941,#0a1527 72%)','important');
  window.cvpnLanguageInfo({localePreference:lang,resolvedTVLocale:lang});
  if(mode==='setup'){
   // Let the real launcher set every translated title, step and footer.
   document.dispatchEvent(new CustomEvent('webOSLaunch',{detail:{mode:'normal'}}));
  } else {
   document.dispatchEvent(new CustomEvent('webOSLaunch',{detail:{mode:'qr'}}));
   // Example address reserved by RFC 5737, non-functional sample code.
   window.cvpnQRInfo({ok:true,code:'000000',address:'http://192.0.2.10:8765',
     pairURL:'http://192.0.2.10:8765/#pair=000000',
     localePreference:lang,resolvedTVLocale:lang,
     expires:new Date(Date.now()+180000).toISOString()});
  }
 },{mode,lang:language});
 assertLanguage(await page.locator('body').innerText(),file);
 await page.screenshot({path:path.join(outDir,screenshotName(file)),animations:'disabled'});
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
  for(const file of fs.readdirSync(outDir))if(file.endsWith('.png')&&file.endsWith(isRU?'-ru.png':'.png'))console.log(file,fs.statSync(path.join(outDir,file)).size,'bytes');
 } finally {await browser.close();server.close();}
})().catch(e=>{console.error(e);server.close();process.exitCode=1;});
