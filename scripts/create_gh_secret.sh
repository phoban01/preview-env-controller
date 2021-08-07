GITHUB_USER=phoban01

kubectl create secret generic gh-token -n preview-env-controller-system \
    --from-literal=username=$GITHUB_USER \
    --from-literal=password=$(echo $GITHUB_TOKEN)
