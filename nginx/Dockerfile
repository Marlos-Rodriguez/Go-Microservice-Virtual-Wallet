FROM nginx:latest

RUN apt-get update
RUN apt-get install sudo
RUN sudo apt-get install nano

COPY ./nginx.conf /etc/nginx/sites-available/nginx.conf
COPY ./nginx.conf /etc/nginx/conf.d/nginx.conf