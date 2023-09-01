# Simple Container Manager

## Overview

Simple Container Manager is a lightweight Go application designed to manage Linux containers. Users can start new containers, while administrators have the ability to list, stop, and manage container lifecycles. The project utilizes Linux namespaces for container isolation.

## Features

### User Functionality:

- Start a new shell process in a new container.
- Perform standard shell operations within the container.

### Admin Functionality:

- List all containers (running or stopped).
- Stop running containers.
- Resume stopped containers.
- Remove any container.

Admin actions get reflected in a JSON file for persistence.

## Usage

### For Users

Starting a Container
To start a new container, run:

```bash
sudo ./container start

```

You'll be placed into a new shell process running in a new, isolated container.

## For Admins

Accessing Admin Shell
To enter admin mode, run:

```bash
sudo ./container admin -p "password"

```

In admin mode, you can perform various container management tasks.

List Containers
This will display a list of all containers, regardless of their current status.

Stop Containers
You can stop any container by its ID.

Resume Containers
If a container is stopped, you can resume its operations.

Remove Containers
You can remove any container by its ID.

Data Persistence
All admin actions are recorded in a JSON file for data persistence.
