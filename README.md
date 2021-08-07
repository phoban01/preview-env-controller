# GitOps Preview Environments

⚠️ Proof of Concept - not production ready

## Flux based GitOps PreviewEnvironments

This project is a proof-of-concept set of controller that implement Preview Environments
in a GitOps compatible manner. The basic principle is to define a set of manifests
which will be used as a template. Then the controller will watch a specified repository for
new branches matching a specified rule.

If a branch matches the rule, a `PreviewEnvironment` will be created using the specified template
with the branch name and commit SHA substituted by the controller.

Rules could be extended to match Pull Request status, comments, labels or tags etc...

Under the hood this makes use of Flux and the GitOps Toolkit components.

### Caveats
Please do not use this in production. This is very much a proof of concept and the
controllers are very primitive without many safeguards and no tests.

### Instructions

1. Ensure you have a local `kind` cluster running (https://kind.sigs.k8s.io/docs/user/quick-start/)

2. Ensure you have the flux components installed (https://fluxcd.io/docs/installation/)

3. Create a secret with your GitHub credentials `$ ./scripts/create_gh_secret.sh`

4. Install the CRDs `$ make install`

5. Deploy the controller `$ make deploy`

6. Create a `PreviewEnvironmentManager` `$ kubectl apply -f examples/preview-env-manager.yaml`

7. Verify the environment has been created `$ kubectl get previewenvironments -A`

8. Verify the deployment `$ kubectl -n default get deploy server -oyaml`

