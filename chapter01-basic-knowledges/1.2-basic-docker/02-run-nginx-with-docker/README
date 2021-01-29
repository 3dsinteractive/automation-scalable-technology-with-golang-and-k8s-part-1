Run Nginx with Docker

0. Look at slides 1.6 Basic Docker (Why we need Docker?)

1. Open vscode at chapter01-basic-knowledges/1.2-basic-docker/02-run-nginx-with-docker

2. Run command
$ docker-compose up -d
Creating network "1202-run-nginx-with-docker_default" with the default driver
Pulling nginx (3dsinteractive/nginx:1.12)...
1.12: Pulling from 3dsinteractive/nginx
21a7e111de8c: Pull complete
65e50e28a4a4: Pull complete
44e1ce638145: Pull complete
1f0a6ff9558f: Pull complete
0f754eafdcc4: Pull complete
fc4e1aa60163: Pull complete
c1c4fbdf6fa4: Pull complete
44bac6404904: Pull complete
5ddecd5c929f: Pull complete
71bba2db9089: Pull complete
74a4030dfede: Pull complete
3fe62741569a: Pull complete
Digest: sha256:03200b49c087d4338f6860bad8cd2e7d02371bd5d2d2e04f3c8901c47f0ff35f
Status: Downloaded newer image for 3dsinteractive/nginx:1.12
Creating nginx ... done

3. Run command
$ docker ps
CONTAINER ID   IMAGE                       COMMAND                  CREATED          STATUS          PORTS                                         NAMES
ebece58ff823   3dsinteractive/nginx:1.12   "/app-entrypoint.sh …"   26 seconds ago   Up 25 seconds   0.0.0.0:8080->8080/tcp, 0.0.0.0:8443->8443/tcp   nginx

4. Run command to request http from nginx
$ curl -X GET "http://localhost:8080"
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>

5. Run command to stop nginx
$ docker-compose down
Stopping nginx ... done
Removing nginx ... done
Removing network 1202-run-nginx-with-docker_default

6. Run command to check nginx is stop and remove
$ docker ps
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
