#!/bin/bash

time=$(date "+%Y-%m-%d %H:%M:%S")
echo ''
echo 'start:${time}'
file='/data/nginx/libmqtt.json'
targetUrl='https://iiot-dne.geega.com/api/iot-service/file/static'
token='eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxNzE2MzU0MjkwNjkzNDQzNjAwIiwidXNlcl9pZCI6IjE3MTYzNTQyOTA2OTM0NDM2MDAiLCJhenAiOiJlbWJlZC1pYW0iLCJ1bmlxdWVfa2V5IjoiMTYwMjllOTctZWZhNi00YTg1LTg0NTMtYjhkM2Q4NDRhNjE1IiwiYWNjZXNzX2p0aSI6IjA2N2ZiNDYzLTQwZWYtNDVlOS04MzdmLTdmYmRkZmY3OWNhZCIsIm5hbWUiOiLnrqHnkIblkZgiLCJpc3MiOiJodHRwOi8vZ3VjMy1hcGktZG5lLmdlZWdhLmNvbS9hcGkvaWFtLzEiLCJ0eXAiOiJCZWFyZXIiLCJyZWFsbSI6IjEiLCJsb2dpbl9zb3VyY2UiOiJtb2JpbGUtcGFzc3dvcmQiLCJqdGkiOiIwNjdmYjQ2My00MGVmLTQ1ZTktODM3Zi03ZmJkZGZmNzljYWQiLCJpYXQiOjE3MTk3OTQ5MzgsImV4cCI6MTcxOTg4MTMzOH0.fv0ui2Bg6aEtJeFF_6y903N4DMnF4J6xUXi20u8meACEweHh9f-ZJGLUyXWItauMQD_662jMlkicta_NAznamg'
/data/nginx/test_app $file $targetUrl $token