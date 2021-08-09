GITHUB_USER=$(gh repo view --json 'owner' --jq '.owner.login')

if [ -z $GITHUB_TOKEN ]; then
    echo "Please set GITHUB_TOKEN environment variable"
    exit 1
fi

kubectl create secret generic gh-token -n preview-env-controller-system \
    --from-literal=username=$GITHUB_USER \
    --from-literal=password=$(echo $GITHUB_TOKEN)
