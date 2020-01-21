## API responses: List releases for a repository

> This returns a list of releases, which does not include regular Git tags that have not been associated with a release. To get a list of Git tags, use the Repository Tags API.

https://developer.github.com/v3/repos/releases/#list-releases-for-a-repository


## Storing a new sample file

There is a quick bash script that uses curl to retrieve sample data from the unauthenticated GitHub API. Simply `./sample.sh <owner> <repo>`.
