package util

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kutilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	e2e "k8s.io/kubernetes/test/e2e/framework"
)

// RemoveDeploymentConfigs deletes the given DeploymentConfigs in a namespace
func RemoveDeploymentConfigs(oc *CLI, dcs ...string) error {
	errs := []error{}
	for _, dc := range dcs {
		e2e.Logf("Removing deployment config %s/%s", oc.Namespace(), dc)
		if err := oc.AdminAppsClient().Apps().DeploymentConfigs(oc.Namespace()).Delete(dc, &metav1.DeleteOptions{}); err != nil {
			e2e.Logf("Error occurred removing deployment config: %v", err)
			errs = append(errs, err)
		}

		err := wait.PollImmediate(5*time.Second, 5*time.Minute, func() (bool, error) {
			pods, err := GetApplicationPods(oc, dc)
			if err != nil {
				e2e.Logf("Unable to get pods for dc/%s: %v", dc, err)
				return false, err
			}
			if len(pods.Items) != 0 {
				return false, nil
			}
			return true, nil
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return kutilerrors.NewAggregate(errs)
	}

	return nil
}
