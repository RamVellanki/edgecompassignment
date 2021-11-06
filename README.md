This is a simple golang application that pushes metrics to Prometheus which are availbale in Grafana.

#### Steps to run in microk8s
1. Install `microk8s`
2. Enable `registry` and `dashboard` plugins  in `microk8s`
3. Run `microk8s inspect` - as you have enabled registry it will ask you to update config, complete the same
4. Create dockerimage using `docker build -f Dockerfile -t localhost:32000/gometrics:v1` 
5. Run `docker push localhost:32000/gometrics:v1`
6. Now run `kubectl apply -f kompose/golang-deployment.yaml`
7. Now run `kubectl apply -f kompose/golang-service.yaml`
8. Install `helm`
9. Run `helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`
10. Run `helm repo update`
11. Run `helm install prometheus prometheus-community/kube-prometheus-stack`
12. Now you have the application, grafana (via dashboard plugin) and prometheus enabled
13. Run `microk8s dashboard-proxy` to see the kubernetes dashboard
14. You should be able to access the services via their IP that you see in Kuberenetes dashboard

 

#### Steps to run in minikube
1. Install `minikube`
2. Run `eval $(minikube docker-env)
3. Create dockerimage using `docker build -f Dockerfile -t gometrics:v1` 
4. Now run `kubectl apply -f kompose/golang-deployment.yaml`
5. Now run `kubectl apply -f kompose/golang-service.yaml`
6. Install `helm`
7. Run `helm repo add prometheus-community https://prometheus-community.github.io/helm-charts`
8. Run `helm repo update`
9. Run `helm install prometheus prometheus-community/kube-prometheus-stack`
10. Now you have the application, grafana (via dashboard plugin) and prometheus enabled
11. Run `minikube dashboard &` to view the Kubernetes dashboard
12. Run `minikube service --url gometrics`, this provides an URL where the gometrics application is running. Use `/metrics` path to view the captured metrics
13. Similarly you can run `grafana` and `prometheus` using the above command