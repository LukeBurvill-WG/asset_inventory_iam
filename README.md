# iam_check
Checks GCP Projects for permissions directly assigned to user identities

## Available CLI Arguments

* -projects
  * Comma separated list of projects you want to check for user assigned permissions
  * if excluded, will run scan against all projects you have access to. Handles lack of permissions to Asset Inventory gracefully and notes in output
  * eg. -projects project_1,project_2
* -limit
  * Stops scan after the specified number of projects. Useful if testing functinality and excluding a list of projects
  * if excluded, all specified project, or all accessible projects, are scanned
  * eg. -limit 1

## How to Use

iam_check.exe -projects project_1,project_2

go run . -projects project_1,project_2


