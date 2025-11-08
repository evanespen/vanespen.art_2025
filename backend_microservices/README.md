# Backend Microservices

Microservices based implementation of the backend.

- `api_service` contains the API endpoints for the backend.
- `resizer_service` contains worker that resizes the image.
- `exif_service` contains worker that extract exif metadata.
- `database_service` contains worker that manages the database.

Communication between services is done using *NATS* and *msgpack*.

Image files are stored in *Minio*.


## TODO
- [ ] exif service
- [ ] database service
- [ ] tests
