package controllers

import (
	"context"
	"regexp"

	v1alpha1 "github.com/phoban01/preview-env-controller/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *PreviewEnvironmentManagerReconciler) branchMatchingStrategy(
	ctx context.Context,
	obj *v1alpha1.PreviewEnvironmentManager,
	existingEnvs *v1alpha1.PreviewEnvironmentList,
	newEnvs map[string]string,
) ([]v1alpha1.PreviewEnvironment, error) {
	log := ctrl.LoggerFrom(ctx)

	branches, err := r.listBranches(ctx, obj)
	if err != nil {
		log.Info("failed to list branches", "error", err.Error())
		return nil, err
	}

	re, err := regexp.Compile(obj.Spec.Strategy.Rules.Match)
	if err != nil {
		return nil, err
	}

	var gcEnvs []v1alpha1.PreviewEnvironment

	if obj.Spec.Prune {
		for _, p := range existingEnvs.Items {
			if _, ok := branches[p.Spec.Branch]; !ok {
				log.Info("branch not found", "branch", p.Spec.Branch)
				gcEnvs = append(gcEnvs, p)
			}

			if ok := re.MatchString(p.Spec.Branch); !ok {
				log.Info("environment no longer matches rule",
					"branch", p.Spec.Branch,
					"rule", obj.Spec.Strategy.Rules.Match)
				gcEnvs = append(gcEnvs, p)
			}
		}
	}

	// iterate branches and check for matching rules
	for branch, sha := range branches {
		if ok := re.MatchString(branch); !ok {
			log.Info("branch doesn't match rule",
				"branch", branch,
				"rule", obj.Spec.Strategy.Rules.Match)
			continue
		}

		log.Info("branch matches rule",
			"branch", branch,
			"rule", obj.Spec.Strategy.Rules.Match)

		newEnvs[branch] = sha
	}

	return gcEnvs, nil
}

func (r *PreviewEnvironmentManagerReconciler) pullRequestMatchingStrategy(
	ctx context.Context,
	obj *v1alpha1.PreviewEnvironmentManager,
	existingEnvs *v1alpha1.PreviewEnvironmentList,
	newEnvs map[string]string,
) ([]v1alpha1.PreviewEnvironment, error) {
	log := ctrl.LoggerFrom(ctx)

	pullRequests, err := r.listPullRequests(ctx, obj)
	if err != nil {
		log.Info("failed to list pull requests", "error", err.Error())
		return nil, err
	}

	var gcEnvs []v1alpha1.PreviewEnvironment

	if obj.Spec.Prune {
		for _, p := range existingEnvs.Items {
			if _, ok := pullRequests[p.Spec.Branch]; !ok {
				log.Info("pull-request branch not found", "branch", p.Spec.Branch)
				gcEnvs = append(gcEnvs, p)
			}
		}
	}

	for branch, sha := range pullRequests {
		newEnvs[branch] = sha
	}

	return gcEnvs, nil
}

func (r *PreviewEnvironmentManagerReconciler) listBranches(ctx context.Context, obj *v1alpha1.PreviewEnvironmentManager) (map[string]string, error) {
	resp, _, err := r.RepoClient.Repositories.ListBranches(ctx, obj.GetOwner(), obj.GetRepo(), nil)
	if err != nil {
		return nil, err
	}

	branches := make(map[string]string, len(resp))

	for _, b := range resp {
		branches[*b.Name] = *b.Commit.SHA
	}

	return branches, nil
}

func (r *PreviewEnvironmentManagerReconciler) listPullRequests(ctx context.Context, obj *v1alpha1.PreviewEnvironmentManager) (map[string]string, error) {
	resp, _, err := r.RepoClient.PullRequests.List(ctx, obj.GetOwner(), obj.GetRepo(), nil)
	if err != nil {
		return nil, err
	}

	branches := make(map[string]string, len(resp))

	for _, b := range resp {
		if !obj.Spec.Strategy.Rules.Draft && *b.Draft || *b.State != "open" {
			continue
		}
		branches[*b.Head.Ref] = *b.Head.SHA
	}

	return branches, nil
}
