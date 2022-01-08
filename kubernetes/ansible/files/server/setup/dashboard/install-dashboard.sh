#!/bin/bash
LOG_FILE=install.log
function clean_dashboard() {
  sudo k3s kubectl delete ns kubernetes-dashboard  >> "$LOG_FILE"
  sudo k3s kubectl delete clusterrolebinding kubernetes-dashboard  >> "$LOG_FILE"
  sudo k3s kubectl delete clusterrole kubernetes-dashboard  >> "$LOG_FILE"
  sudo kubectl delete clusterrolebinding admin-user >> "$LOG_FILE"
}

[[ -f "$LOG_FILE" ]] && rm "$LOG_FILE"
[[ $(sudo kubectl get namespace | grep kubernetes-dashboard) ]]; clean_dashboard   >> "$LOG_FILE"

GITHUB_URL=https://github.com/kubernetes/dashboard/releases
VERSION_KUBE_DASHBOARD=$(curl -w '%{url_effective}' -I -L -s -S ${GITHUB_URL}/latest -o /dev/null | sed -e 's|.*/||')
sudo k3s kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/${VERSION_KUBE_DASHBOARD}/aio/deploy/recommended.yaml   >> "$LOG_FILE"
sudo k3s kubectl create -f admin-user.yml -f admin-user-role.yml   >> "$LOG_FILE"

PROXY=/usr/local/bin/start-proxy
cat <<EOF > "$PROXY"
  echo 'URL to server is http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/'
  sudo k3s kubectl -n kubernetes-dashboard describe secret admin-user-token | grep '^token'
  sudo kubectl proxy
EOF
sudo chmod 755 "$PROXY"

