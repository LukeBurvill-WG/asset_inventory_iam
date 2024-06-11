package main

import (
	"context"
	"flag"
	// "fmt"
	"log"
	"strings"
	"encoding/csv"
	"os"

	asset "cloud.google.com/go/asset/apiv1"
	"cloud.google.com/go/asset/apiv1/assetpb"
	"google.golang.org/api/iterator"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
)

func main() {
	projects := flag.String("projects", "", "Comma Seperate list of Project IDs to Scope the search.")
	limit := flag.Int("limit", 0, "Limit the number of projects to query")
	flag.Parse()

	project_list := *projects
	project_limit := *limit

	all_projects := map[string]string{}

	ctx_projects := context.Background()
	c, err := resourcemanager.NewProjectsClient(ctx_projects)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	rqst := &resourcemanagerpb.SearchProjectsRequest{}
	it_projects := c.SearchProjects(ctx_projects, rqst)

	for {
		resp, err := it_projects.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		all_projects[resp.Name] = resp.ProjectId
	}

	var project_array []string
	if project_list == "" {
		for _, v := range all_projects {
			project_array = append(project_array, v)
		}
		project_list = strings.Join(project_array, ",")
	}

	ctx_iam := context.Background()
	client, err := asset.NewClient(ctx_iam)
	if err != nil {
		log.Fatalf("asset.NewClient: %v", err)
	}
	defer client.Close()

	results := []map[string]string{}
	count := 0

	for _, project := range strings.Split(project_list, ",") {

		count++ 
		if count > project_limit {
			if project_limit != 0 {
				break
			}
		}
		log.Printf("Processing %s", project)
		scope := "projects/" + project

		req := &assetpb.SearchAllIamPoliciesRequest{
			Scope: scope,
			Query: "",
		}
		it_iam := client.SearchAllIamPolicies(ctx_iam, req)
		for {
			policy, err := it_iam.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Print(err)
				results = append(results, map[string]string{"project": req.Scope, "resource": "Do Not Have Permission to Project", "role": "", "identity": ""})
				break
			} else {
				for _, r := range policy.Policy.Bindings {
					for _, m := range r.Members {
						if strings.Contains(m, "user:") {
							results = append(results, map[string]string{"project": policy.Project, "resource": policy.Resource, "role": r.Role, "identity": m})
						}
					}
				}
			}
		}
	}

	file, err := os.Create("iam.csv")
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Project ID", "Project Number", "Resource", "Role", "Identity"}); err != nil {
		log.Fatalln("Error Writing Record to File", err)
	}

	for _, iam := range results {
		if err := writer.Write([]string{all_projects[iam["project"]], strings.Trim(iam["project"],"projects/"), iam["resource"], iam["role"], iam["identity"]}); err != nil {
			log.Fatalln("Error Writing Record to File", err)
		}
	}

}
