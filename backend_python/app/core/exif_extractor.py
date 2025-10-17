import os
import re

import pendulum


def get_exif(file_path):
    file_path = file_path.replace('(', '\(')
    file_path = file_path.replace(')', '\)')
    RGX = {
        'cam_model': 'Exif.Image.Model (.*)',
        'timestamp': 'Exif.Photo.DateTimeOriginal (.*)',
        'exposure': 'Exif.Photo.ExposureTime (.*)',
        'aperture': 'Exif.Photo.FNumber (.*)',
        'mode': 'Exif.Photo.ExposureProgram (.*)',
        'iso': 'Exif.Photo.ISOSpeedRatings (.*)',
        'flash': 'Exif.Photo.Flash (.*)',
        'focal': 'Exif.Photo.FocalLength (.*)',
        'lens': 'Exif.NikonLd3.LensIDNumber (.*)',
        'lens2': 'Exif.Photo.LensModel (.*)',
        'focal_equiv': 'Exif.Photo.FocalLengthIn35mmFilm (.*)'
    }

    cmd = "exiv2 -q -p a file_path | sed -E $'s/[[:space:]]{2,}/\#/g' | sed -e $'s/ /_/g' | sed -E $'s/\#/ /g' | awk '{printf \"%s %s\\n\", $1, $4}'".replace(
        'file_path', file_path)
    p = os.popen(cmd)
    infos = p.read()

    data = {
        'cam_model': '',
        'timestamp': '',
        'exposure': '',
        'aperture': '',
        'mode': '',
        'iso': '',
        'flash': '',
        'focal': '',
        'lens': '',
        'lens2': '',
        'focal_equiv': ''
    }

    for k, v in RGX.items():
        for l in infos.split('\n'):
            try:
                if k == 'timestamp':
                    m = re.search(v, l)
                    datestr = re.sub(
                        r' {2,}', '', m.group(1)).replace('_', ' ')
                    date = pendulum.parse(datestr)
                    if date.year == 2018:
                        date = date.add(years=1)
                    data[k] = date
                else:
                    m = re.search(v, l)
                    data[k] = re.sub(
                        r' {2,}', '', m.group(1)).replace('_', ' ')
            except AttributeError:
                pass
            except Exception as e:
                print(e)
                pass

    if data['timestamp'] == '':
        file_date = os.stat(file_path).st_mtime
        data['timestamp'] = pendulum.from_timestamp(int(file_date))

    print(
        '!!!!!!!!!!', data['lens2'])
    if data['lens'] == '':
        if data['lens2'] == '200.0-500.0 mm f/5.6':
            data['lens'] = 'Nikon 200.0-500.0 mm f/5.6 ED VR'

    data.pop('lens2')

    return data
