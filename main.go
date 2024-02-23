package main

import "log"

func main() {
	var policies, err = getExistingPolicies()
	if err != nil {
		log.Fatalf("error getting existing policies: %v", err)
	}
	for _, value := range policies {
		create_azurecapolicy(value)
	}
	for _, value := range policies {
		import_policy_to_tfstate(value)
	}

}
