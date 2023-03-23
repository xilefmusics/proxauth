# proxauth

## Rules

|Param|Default|Explanation|
|-|-|-|
|fromScheme|`http`|The scheme which is used for redirects and logging. It is not responisble for matching the rule.|
|fromHost|`*`|Selects the rule based on the host. Use `*` to match all hosts.|
|fromPath|`/`|The subpath which is used to match the rule.|
|toScheme|`http`|The url scheme which is used to forward the request.|
|toHost|`localhost`|The host the request is forwarded to.|
|toPort||The port the request is forwarded to.|
|toPath|`/`|The path the request is forwarded to.|
|loginPath|`/login`|The path which is used for login. (GET & POST is blocked)|
|loginPath|`/logout`|The path which points to the logout page. (GET is blocked)|
|allowedUsers|`[]`|List of the users that are allowed to access the endpoints. Use `[]` if authentication is disabled and all users can access it.|
|redirectToLogin|`false`|If this is set to true the unauthorized respone is catched and a redirection to the loginpage is made instead.|
|backgroundColor|`#000000`|Color scheme used in html pages surfed by proxauth itself.|
|textColor|`#002200`|Color scheme used in html pages surfed by proxauth itself.|
|primaryColor|`#002200`|Color scheme used in html pages surfed by proxauth itself.|
|title|`Proxauth`|Title of the html pages surfed by proxauth itself.|