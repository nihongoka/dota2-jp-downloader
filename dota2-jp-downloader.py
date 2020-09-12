#!/usr/bin/python
# -*- coding: utf-8 -*-

import os.path
import vpk
import urllib.request
import io
import zipfile
import json
import vdf
import shutil

url = 'https://nihongoka.games/download/dota2/?format=zip'

default = (lambda data: {"lang": {"Language": "japanese", "Tokens": data}})
def simple(key):
    return lambda data: {key: data}

files = {
    'abilities_japanese.json': (default, 'abilities_japanese.txt'),
    'broadcastfacts_japanese.json': (default, 'broadcastfacts_japanese.txt'),
    'chat_japanese.json': (default, 'chat_japanese.txt'),
    'dota_japanese.json': (default, 'dota_japanese.txt'),
    'gameui_japanese.json': (default, 'dota_japanese.txt'),
    'hero_chat_wheel_japanese.json': (simple('hero_chat_wheel'), 'hero_chat_wheel_japanese.txt'),
    'hero_lore_japanese.json': (default, 'hero_lore_japanese.txt'),
    'leagues_japanese.json': (simple('leagues'), 'leagues_japanese.txt'),
    'richpresence_japanese.json': (default, 'richpresence_japanese.txt'),
}

os.makedirs('game/dota_japanese/', exist_ok=True)

os.chdir('game/dota_japanese/')

shutil.rmtree('pak01', ignore_errors=True)

os.makedirs('pak01/resource/localization/patchnotes', exist_ok=True)

print('ダウンロード中…')
req = urllib.request.Request(url)
with urllib.request.urlopen(req) as res:
    print('ダウンロード完了')
    with zipfile.ZipFile(io.BytesIO(res.read())) as zip:
        print('展開完了')
        for file in zip.filelist:
            basename = os.path.basename(file.filename)
            if not basename in files:
                continue
            fi = files[basename]
            data = json.load(zip.open(file.filename, 'r'))
            fp = 'pak01/resource/localization/' + fi[1]
            with open(fp, 'w', encoding='utf-8') as wf:
                wf.write(vdf.dumps(fi[0](data), pretty=True))
print('書き出し中…')
vpk.new('./pak01').save('pak01_dir.vpk')
print('完了!')