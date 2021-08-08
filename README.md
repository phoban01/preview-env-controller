# GitOps Preview Environments

⚠️ Proof of Concept - not production ready

***Please do not use this in production. This is very much a proof of concept and the controllers are very primitive without many safeguards and no tests.***

## Flux based GitOps PreviewEnvironments

https://fluxcd.io/

This project is a proof-of-concept set of controllers that implement Preview Environments in a GitOps compatible fashion. The basic principle is to define a set of manifests which will be used as a template. Then the controller will watch a specified repository and create a new environment using the template when specified conditions are met.

With `strategy.type=Branch` if a branch matches the rule, a `PreviewEnvironment` will be created using the specified template
with the branch name and commit SHA substituted by the controller.

With `strategy.type=PullRequest` if a PullRequest is created on the watched repository then a `Preview Environment` will be created.

Under the hood this makes use of Flux and the GitOps Toolkit components.

### Limitations

- Currently only supports GitHub provider
- GitHub credentials must be specified using a secret with the keys `username` (GITHUB_USER) and `password` (GITHUB_PAT)
- Only branch name and commit sha are substituted into environment
- No Helm support yet (only supports Kustomize currently)

### Configuration

Configuration is handled via a `PreviewEnvironmentManager` resource.

#### Pull Request Environments

To configure a `PreviewEnvironmentManager` that will watch pull requests which have status `Ready` use the following resource:

```yaml
---
apiVersion: preview.gitops.phoban.io/v1alpha1
kind: PreviewEnvironmentManager
metadata:
  name: demo-branch-manager
  namespace: preview-env-controller-system
spec:
  interval: 2m
  limit: 2
  watch:
    url: https://github.com/phoban01/sample-app
    ref:
      branch: "main"
    credentialsRef:
      name: gh-token
  strategy:
    type: PullRequest
    rules:
      draft: false
  template:
    prefix: pr-env-
    interval: 1m
    createNamespace: true
    kustomizationSpec:
      path: './preview-template'
      prune: true
```

#### Branch Environments

The following resource manifest will configure a manager that will create a new `PreviewEnvironment` whenever a branch matching the regex pattern is pushed to the repo:

``` yaml
---
apiVersion: preview.gitops.phoban.io/v1alpha1
kind: PreviewEnvironmentManager
metadata:
  name: demo-branch-manager
  namespace: preview-env-controller-system
spec:
  interval: 2m
  limit: 2
  watch:
    url: https://github.com/phoban01/sample-app
    ref:
      branch: "main"
    credentialsRef:
      name: gh-token
  strategy:
    type: Branch
    rules:
      match: '^feature-.*$'
  template:
    prefix: pr-env-
    interval: 1m
    createNamespace: true
    kustomizationSpec:
      path: './preview-template'
      prune: true
```

### Demo

1. Ensure you have a local `kind` cluster running (https://kind.sigs.k8s.io/docs/user/quick-start/)

2. Ensure you have the flux components installed (https://fluxcd.io/docs/installation/)

3. Create a secret with your GitHub credentials `$ ./scripts/create_gh_secret.sh`

4. Install the CRDs `$ make install`

5. Deploy the controller `$ make deploy`

6. Create a `PreviewEnvironmentManager` `$ kubectl apply -f examples/branch-manager.yaml`

7. Verify the environment has been created `$ kubectl get previewenvironments -A`

8. Verify the deployment `$ kubectl -n default get deploy server -oyaml`

