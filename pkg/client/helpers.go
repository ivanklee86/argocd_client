package client

import (
	"context"
	"os"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/ivanklee86/octanap/pkg/testhelpers"
	"github.com/joho/godotenv"
)

func GenerateTestProjects() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	clientOptions := ArgoCDClientOptions{
		ServerAddr: "localhost:8080",
		Insecure:   true,
		AuthToken:  os.Getenv("ARGOCD_TOKEN"),
	}

	argoCDClient, _ := New(&clientOptions)
	projects := []*v1alpha1.AppProject{}
	for range 3 {
		project, _ := argoCDClient.CreateProject(context.Background(), testhelpers.RandomProjectName())
		projects = append(projects, project)
	}
	defer deleteTestProjects(argoCDClient, projects)

	// Set up Project 1 without SyncWindow
	project1 := projects[0]
	project1.Labels = make(map[string]string)
	project1.Labels["team"] = "Jets"
	project1.Labels["env"] = "prod"
	project1.Labels["department"] = "A"

	// Set up Project 2 with SyncWindow
	project2 := projects[1]
	project2.Labels = make(map[string]string)
	project2.Labels["team"] = "Giants"
	project2.Labels["env"] = "prod"
	project2.Labels["department"] = "B"
	project2.Spec.SyncWindows = append(project2.Spec.SyncWindows, &v1alpha1.SyncWindow{
		Kind:       "allow",
		Schedule:   "10 1 * * *",
		Duration:   "1h",
		Namespaces: []string{"*"},
	})

	// Set up Project 3
	project3 := projects[2]
	project3.Labels = make(map[string]string)
	project3.Labels["team"] = "Eagles"
	project3.Labels["env"] = "dev"
	project3.Labels["department"] = "C"

	for _, project := range projects {
		_, err := argoCDClient.UpdateProject(context.Background(), *project)
		if err != nil {
			panic(err)
		}
	}
}

func deleteTestProjects(argoCDClient IArgoCDClient, projects []*v1alpha1.AppProject) {
	for _, project := range projects {
		_, err := argoCDClient.DeleteProject(context.Background(), project.ObjectMeta.Name)
		if err != nil {
			panic(err)
		}
	}
}