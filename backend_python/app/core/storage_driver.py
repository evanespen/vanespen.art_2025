from os import makedirs, path

from PIL import Image

from app.core import config


for d in [
        config.PICTURES_DIR,
        config.THUMBS_DIR,
        config.HALFRES_DIR,
        config.TEMP_DIR
]:
    if not path.isdir(d):
        makedirs(d)


def make_thumbs(picture):
    file_path = path.join(config.PICTURES_DIR, picture.path)

    thumb_path = path.join(config.THUMBS_DIR, picture.path)
    if not path.isfile(thumb_path):
        im = Image.open(file_path)
        im.thumbnail(config.THUMB_SIZE)
        im.save(thumb_path)

    halfres_path = path.join(config.HALFRES_DIR, picture.path)
    if not path.isfile(halfres_path):
        im = Image.open(file_path)
        im.thumbnail(config.HALFRES_SIZE)
        im.save(halfres_path)
