GITHUB_USER=phoban01

kubectl create secret generic gh-token -n default \
    --from-literal=username=$GITHUB_USER \
    --from-literal=password=$(echo $GITHUB_TOKEN)
