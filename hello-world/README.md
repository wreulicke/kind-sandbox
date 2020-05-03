# hello world

## kindのインストール

```bash
curl -LO https://github.com/kubernetes-sigs/kind/releases/download/v0.7.0/kind-darwin-amd64
chmod +x kind-darwin-amd64
mv kind-darwin-amd64 /usr/local/bin/
```

## podの作成

```
kubectl apply -f pod.yaml
```

## podの削除 

```
kubectl delete -f pod.yaml
```

## port-forward

ホストの8080にpodの80をport forward

```
kubectl port-forward pod/sample-pod 8080:80
```

## exec 

```
kubectl exec -it sample-pod bash
```