# iam_check
Checks GCP Projects for permissions directly assigned to user identities

## Available CLI Arguments

* -projects
  * Comma separated list of projects you want to check for user assigned permissions
  * eg. -projects project_1,project_2
* -file
  * full name of a file in the same folder as iam_check that contains a list of projects to scan, each one on a new line
  * eg. -file project_list.csv
* -output
  * Name of file to output scan results to. Automatically adds .csv to name
  * defaults to iam.csv
  * eg. -output results
* -limit
  * Stops scan after the specified number of projects. Useful if testing functinality and excluding a list of projects
  * if excluded, all specified project, or all accessible projects, are scanned
  * eg. -limit 1

## How to Use

If both -projects and -file are excluded, will run scan against all projects you have access to. Handles lack of permissions to Asset Inventory gracefully and notes in output

iam_check.exe -projects project_1,project_2

go run . -projects project_1,project_2


