1. Install docker from https://docs.docker.com/get-docker/

2. Double click Docker icon in Application directory, wait until docker is running in menu bar

3. After install open terminal and run command
$ docker version
Client: Docker Engine - Community
 Version:           20.10.2
 API version:       1.41
 Go version:        go1.13.15
 Git commit:        2291f61
 Built:             Mon Dec 28 16:12:42 2020
 OS/Arch:           darwin/amd64
 Context:           default
 Experimental:      true

Server: Docker Engine - Community
 Engine:
  Version:          20.10.2
  API version:      1.41 (minimum version 1.12)
  Go version:       go1.13.15
  Git commit:       8891c58
  Built:            Mon Dec 28 16:15:28 2020
  OS/Arch:          linux/amd64
  Experimental:     false
 containerd:
  Version:          1.4.3
  GitCommit:        269548fa27e0089a8b8278fc4fc781d7f65a939b
 runc:
  Version:          1.0.0-rc92
  GitCommit:        ff819c7e9184c13b7c2607fe6c30ae19403a7aff
 docker-init:
  Version:          0.19.0
  GitCommit:        de40ad0

4. Open Docker Preference > Kubernetes > Check box at Enable Kubernetes, 
wait until Kubernetes is running (If Kubernetes starting forever, try reset factory default)

5. After Kubernetes is running, Run command
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"16", GitVersion:"v1.16.0", GitCommit:"2bd9643cee5b3b3a5ecbd3af49d09018f0773c77", GitTreeState:"clean", BuildDate:"2019-09-18T14:36:53Z", GoVersion:"go1.12.9", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.3", GitCommit:"1e11e4a2108024935ecfcb2912226cedeafd99df", GitTreeState:"clean", BuildDate:"2020-10-14T12:41:49Z", GoVersion:"go1.15.2", Compiler:"gc", Platform:"linux/amd64"}

6. Run command
$ kubectl get node
NAME             STATUS   ROLES    AGE     VERSION
docker-desktop   Ready    master   3m44s   v1.19.3

7. Install nginx-ingress follow the step from this URL
https://kubernetes.github.io/ingress-nginx/deploy/#docker-for-mac
TLDR; Run this command
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.43.0/deploy/static/provider/cloud/deploy.yaml

8. Run command
$ kubectl get ns
You must see ingress-nginx in the list of namespace

9. Run command
$ kubectl get po -n ingress-nginx
NAME                                        READY   STATUS      RESTARTS   AGE
ingress-nginx-admission-create-mlwmr        0/1     Completed   0          5m19s
ingress-nginx-admission-patch-fpk4v         0/1     Completed   0          5m19s
ingress-nginx-controller-56c75d774d-xwjhq   1/1     Running     0          5m19s

10. Run command
$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.109.191.187   localhost     80:31170/TCP,443:31054/TCP   5m54s
ingress-nginx-controller-admission   ClusterIP      10.106.142.243   <none>        443/TCP                      5m54s