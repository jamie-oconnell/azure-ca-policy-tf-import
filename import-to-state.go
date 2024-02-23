package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func import_policy_to_tfstate(policy models.ConditionalAccessPolicy) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		log.Fatalf("error getting Terraform executable path: %s", err)
	}
	workingDir := "generated"
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	resource_name := *policy.GetDisplayName()
	resource_name = strings.ToLower(strings.ReplaceAll(resource_name, " ", "_"))
	err = tf.Import(context.Background(), fmt.Sprintf("azuread_conditional_access_policy.%s", resource_name), *policy.GetId())
	if err != nil {
		log.Printf("error running Import: %s", err)
	}

}
