package cli

import (
	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
)

func isMapSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}

func arrayContains(a []string, v string) bool {
	for _, item := range a {
		if item == v {
			return true
		}
	}
	return false
}

func filterProjects(appProjects *v1alpha1.AppProjectList, names []string, labels map[string]string, hasSyncWindow bool) []v1alpha1.AppProject {
	matchingProjects := []v1alpha1.AppProject{}

	if len(names) > 0 {
		for _, appProject := range appProjects.Items {
			if arrayContains(names, appProject.ObjectMeta.Name) {
				matchingProjects = append(matchingProjects, appProject)
			}
		}
	} else {
		for _, appProject := range appProjects.Items {
			if isMapSubset(appProject.ObjectMeta.Labels, labels) {
				// Additional filter to only selct projects with SyncWindows.
				if hasSyncWindow {
					if len(appProject.Spec.SyncWindows) > 0 {
						matchingProjects = append(matchingProjects, appProject)
					}
				} else {
					matchingProjects = append(matchingProjects, appProject)
				}
			}
		}
	}

	return matchingProjects
}
