from os import getenv, listdir

import httpx
from dotenv import load_dotenv

load_dotenv('./.env')

BASE_URL = 'http://localhost:5000/api/'

token = httpx.post(BASE_URL + 'login/access-token',
                   data={'username': getenv('SUPERUSER'),
                         'password': getenv('SUPERUSER_PASSWORD')}).json()['access_token']

headers = {'Authorization': 'bearer ' + token}

# SPECIES = [{
#     "name": "chevreuil",
#     "scientific_name": "chevreuil",
#     "threat": "LC",
#     "info_page": "http://example.com/chevreuil"
# }, {
#     "name": "bernache du canada",
#     "scientific_name": "bernache du canada",
#     "threat": "LC",
#     "info_page": "http://example.com/bernache_du_canada"
# }, {
#     "info_page": "https://www.oiseaux.net/oiseaux/rougequeue.noir.html",
#     "name": "Rougequeue noir",
#     "scientific_name": "Phoenicurus ochruros - Black Redstart",
#     "threat": "LC"
# }, {
#     "name": "chevreuil2",
#     "scientific_name": "chevreuil2",
#     "threat": "LC",
#     "info_page": "http://example.com/chevreuil2"
# }]

# for species in SPECIES:
#     httpx.post(BASE_URL + 'species/', json=species, headers=headers)

# for i in range(1, 3):
#     files = {
#         'files': open(f'test_images/img{i}.jpg', 'rb')
#     }
#     httpx.post(BASE_URL + 'pictures/upload/species-id/' +
#                str(i), files=files, headers=headers)


# files = {
#     'files': open('test_images/img1.jpg', 'rb')
# }
# httpx.post(BASE_URL + 'pictures/upload/species-id/' +
#            str(1), files=files, headers=headers)


for f in listdir('test_images'):
    files = {
        'files': open(f'test_images/{f}', 'rb')
    }
    httpx.post(BASE_URL + 'pictures/upload/', files=files, headers=headers)
