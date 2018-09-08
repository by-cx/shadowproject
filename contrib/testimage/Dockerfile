FROM ubuntu:bionic

# Image access points
VOLUME /srv
EXPOSE 8000

# App dependencies
RUN apt-get update && apt-get upgrade -y

# Clean things we don't need
RUN apt-get clean

# Image related stuff
RUN useradd -d /srv -U app

# User everything should run under
USER app
