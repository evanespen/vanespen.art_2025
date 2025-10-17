[![NO AI](https://raw.githubusercontent.com/nuxy/no-ai-badge/master/badge.svg)](https://github.com/nuxy/no-ai-badge)

## Project Description

vanespen.art is a web application for displaying and managing a photography portfolio.
The application allows for displaying pictures organized in albums, with separate public and admin interfaces.
The public interface displays the art portfolio, while the admin interface allows for managing pictures and albums.

This project is also my developer lab and portfolio as I remake the project as I discover new languages or technologies.

## Requirements of the application

- Serving 1000+ pictures without latency
- Handling albums and animals species groups of pictures
- Admin interface to upload the images
- Exif data extraction on the backend

## Contents of the repo

### `backend_iex`

> This implementation is going to be deployed to production soon.

The *Elixir* implementation of the backend API. It is built using the *Phoenix* framework and sqlite database.

### `backend_go`

> This version was never used, superseded by the Elixir implementation.

The *Go* implementation of the backend API. It is built using the *Gin* web framework and Apache *parquet* files for database.

### `fullstack_sveltekit`

> Implementation currently in production since 2020.

A full-stack implementation using *SvelteKit* that combines both frontend and backend functionality in a single application.
Uses *postgresql* as database.

### `backend_python`

> The first implementation of the backend, superseded by the SvelteKit implementation.

The *Python* implementation of the backend API, using *FastAPI* and *postgresql*.


## TODO

- [ ] New frontend (Phoenix LiveView or VueJS)
- [ ] Statistics page