
## 作成

```
kubectl apply -f all.yaml
```

## 削除 

```
kubectl delete -f all.yaml
```

## podの一覧

```
kubectl get po -l app=nginx -o wide
NAME                                READY   STATUS    RESTARTS   AGE     IP           NODE                 NOMINATED NODE   READINESS GATES
nginx-deployment-574b87c764-fvqlg   1/1     Running   0          7m17s   10.244.0.9   kind-control-plane   <none>           <none>
nginx-deployment-574b87c764-jsz22   1/1     Running   0          7m17s   10.244.0.7   kind-control-plane   <none>           <none>
nginx-deployment-574b87c764-zt97b   1/1     Running   0          7m17s   10.244.0.8   kind-control-plane   <none>           <none>
```

### ホストからcurlでpodのIPにリクエストを飛ばしてみる

```bash
curl http://10.244.0.9:80 
# もちろん通らない
```

## port-forward

ホストの8080にpodの80をport forward

```
kubectl port-forward nginx-deployment-574b87c764-fvqlg 8080:80
```
