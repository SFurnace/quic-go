#!/bin/bash

openssl req -config ./cert.ini -new -x509 -sha256 -newkey rsa:2048 -nodes -keyout privkey.pem -days 365 -out fullchain.pem
