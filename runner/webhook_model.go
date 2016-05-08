package runner

type GogsPayload struct {
	After   string `json:"after"`
	Before  string `json:"before"`
	Commits []struct {
		Author struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"author"`
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
	} `json:"commits"`
	CompareURL string `json:"compare_url"`
	Pusher     struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"pusher"`
	Ref        string `json:"ref"`
	Repository struct {
		Description string `json:"description"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Owner       struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"owner"`
		Private  bool   `json:"private"`
		URL      string `json:"url"`
		Watchers int    `json:"watchers"`
		Website  string `json:"website"`
	} `json:"repository"`
	Secret string `json:"secret"`
	Sender struct {
		AvatarURL string `json:"avatar_url"`
		ID        int    `json:"id"`
		Login     string `json:"login"`
	} `json:"sender"`
}

/*
{
  "secret": "",
  "ref": "refs/heads/master",
  "before": "7fe13a75367efce156f6aa28c023cad075221841",
  "after": "7fe13a75367efce156f6aa28c023cad075221841",
  "compare_url": "",
  "commits": [
    {
      "id": "7fe13a75367efce156f6aa28c023cad075221841",
      "message": "push!\n",
      "url": "https://git.noapp.net/adrian/psmon/commit/7fe13a75367efce156f6aa28c023cad075221841",
      "author": {
        "name": "Adrian Huang",
        "email": "xinhuang327@gmail.com",
        "username": ""
      }
    }
  ],
  "repository": {
    "id": 23,
    "name": "psmon",
    "url": "https://git.noapp.net/adrian/psmon",
    "ssh_url": "root@git.noapp.net:adrian/psmon.git",
    "clone_url": "https://git.noapp.net/adrian/psmon.git",
    "description": "",
    "website": "",
    "watchers": 1,
    "owner": {
      "name": "adrian",
      "email": "xinhuang327@gmail.com",
      "username": "adrian"
    },
    "private": true,
    "default_branch": "master"
  },
  "pusher": {
    "name": "adrian",
    "email": "xinhuang327@gmail.com",
    "username": "adrian"
  },
  "sender": {
    "login": "adrian",
    "id": 1,
    "avatar_url": "/avatars/1"
  }
}
*/

type DronePayload struct {
	Build struct {
		Number      int    `json:"number"`
		Status      string `json:"status"`
		StartedAt   int    `json:"started_at"`
		FinishedAt  int    `json:"finished_at"`
		Event       string `json:"event"`
		Commit      string `json:"commit"`
		Branch      string `json:"branch"`
		Message     string `json:"message"`
		Author      string `json:"author"`
		AuthorEmail string `json:"author_email"`
		LinkURL     string `json:"link_url"`
	} `json:"build"`
	Repo struct {
		Owner    string `json:"owner"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		LinkURL  string `json:"link_url"`
		CloneURL string `json:"clone_url"`
	} `json:"repo"`
	System struct {
		LinkURL string `json:"link_url"`
	} `json:"system"`
}

/*

{
    "build": {
        "number": 22,
        "status": "success",
        "started_at": 1421029603,
        "finished_at": 1421029813,
        "event": "push",
        "commit": "7fd1a60",
        "branch": "master",
        "message": "Update README",
        "author": "octocat",
        "author_email": "octocat@github.com",
        "link_url": "https://github.com/octocat/Hello-World/commit/7fd1a60"
    },
    "repo": {
        "owner": "octocat",
        "name": "hello-world",
        "full_name": "octocat/hello-world",
        "link_url": "https://github.com/octocat/hello-world",
        "clone_url": "https://github.com/octocat/hello-world.git"
    },
    "system": {
        "link_url": "https://drone.mycompany.com"
    }
}

*/
