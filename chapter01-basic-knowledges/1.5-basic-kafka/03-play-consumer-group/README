Play Consumer Group
===

Kafkacat URL
https://github.com/edenhill/kafkacat

1. Start all k8s objects
$ kubectl apply -f .

2. Check if kafka and zookeeper deployment is Running
$ kubectl get po -n basic-kafka
NAME                    READY   STATUS    RESTARTS   AGE
client-util             1/1     Running   0          5m28s
kfk1-86886b6b84-75p7p   1/1     Running   2          5m29s
kfk2-5b69dfcdb4-wk9lt   1/1     Running   2          5m29s
kfk3-6d4c8874c6-47dxp   1/1     Running   2          5m29s
zk1-76cc547698-8865g    1/1     Running   0          5m29s
zk2-7bb59d6788-5dfxx    1/1     Running   0          5m29s
zk3-566db54d6b-xqjzh    1/1     Running   0          5m29s

3. Exec into client-util
$ kubectl exec -it client-util -n basic-kafka -- bash
root@client-util:/#

4. Start Producer and send some messages
$ kafkacat -P -b "kfk1,kfk2,kfk3" -t "mytopic1"

$ message1
$ message2
$ message3

5. Open another tab in terminal, then exec into client-util shell
$ kubectl exec -it client-util -n basic-kafka -- bash
root@client-util:/#

6. Start Consumer 1 in new terminal
$ kafkacat -C -b "kfk1,kfk2,kfk3" -G mygroup mytopic1

7. Open another tab in terminal, then exec into client-util shell
$ kubectl exec -it client-util -n basic-kafka -- bash
root@client-util:/#

8. Start Consumer 2 in new terminal
$ kafkacat -C -b "kfk1,kfk2,kfk3" -G mygroup mytopic1

9. Send some more messages in Producer
$ message4
$ message5
$ message6
$ message7
$ message8

10. Check each Consumers will receive new messages

12. Exit from Consumer using (Ctrl + c)

13. When Consumer exit, see the rebalance process happen

14. Exit from all consumer using (Ctrl + c)

15. Exit from all client-util
$ exit

16. Cleanup workshop
$ kubectl delete ns basic-kafka



