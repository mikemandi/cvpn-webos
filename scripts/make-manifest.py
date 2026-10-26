#!/usr/bin/env python3
"""Create accurate Homebrew Channel manifest after IPK is built and repo known."""
import argparse,hashlib,json,pathlib,re
p=argparse.ArgumentParser()
p.add_argument('--repository',required=True,help='GitHub owner/repository')
p.add_argument('--tag',default='v2.0.0')
p.add_argument('--ipk',default='dist/su.cvpn.webos_2.1.3_all.ipk')
p.add_argument('--out',default='dist/su.cvpn.webos.manifest.json')
a=p.parse_args()
if not re.fullmatch(r'[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+',a.repository):p.error('repository must be owner/repo')
root=pathlib.Path(__file__).resolve().parent.parent
ipk=root/a.ipk
if not ipk.is_file():p.error('IPK missing: '+str(ipk))
version=json.loads((root/'launcher/appinfo.json').read_text())['version']
base=f'https://github.com/{a.repository}/releases/download/{a.tag}'
data={'id':'su.cvpn.webos','version':version,'type':'web','title':'CVPN webOS','appDescription':'CVPN Connect — private network on rooted LG TVs',
      'iconUri':f'https://raw.githubusercontent.com/{a.repository}/main/assets/icon160.png',
      'sourceUrl':f'https://github.com/{a.repository}', 'rootRequired':True,
      'ipkUrl':f'{base}/{ipk.name}',
      'ipkHash':{'sha256':hashlib.sha256(ipk.read_bytes()).hexdigest()},'ipkSize':ipk.stat().st_size}
output=root/a.out
output.parent.mkdir(parents=True,exist_ok=True)
output.write_text(json.dumps(data,ensure_ascii=False,indent=2)+'\n')
print('MANIFEST_CREATED='+str(output))
