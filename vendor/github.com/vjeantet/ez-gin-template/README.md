# Easy Gin Template

When i start using Gin i struggling with template system. I find out many people have the same problem.

I have found some code in Github ([multitemplate.go](https://github.com/gin-gonic/contrib/tree/master/renders/multitemplate), [gin_html_render.go](https://gist.github.com/madhums/4340cbeb36871e227905)) which help but not everything i need is supported like Template helpers.

**Check it out** my package in official gin contribute repository: [gin-gonic/contrib](https://github.com/gin-gonic/contrib)

### Feature
- Simple rendering syntax for the template

```go
  // suppose "app/views/articles/list.html" is your file to be rendered
  c.HTML(http.StatusOK, "articles/list", "")
```

- Configure layout file
- Configure template file extension
- Configure templates directory
- Feels friendlier for people coming from communities like rails, express or django.
- **Template Helpers** ([gin_html_render.go](https://gist.github.com/madhums/4340cbeb36871e227905) is not support yet)

### How to use

Suppose your structure is
```go
|-- app/views/
    |-- layouts/
        |--- base.html
    |-- blogs/
        |--- index.html          
        |--- show.html

See in "example" folder
```

##### 1. Download package to your workspace
```go
go get https://github.com/michelloworld/ez-gin-template
```

##### 2. Import package to your application (*Import with alias)
```go
import eztemplate "github.com/michelloworld/ez-gin-template"
```

##### 3. Enjoy
```go
  r := gin.Default()

  render := eztemplate.New()

  // render.TemplatesDir = "app/views/" // default

  // render.Layout = "layouts/base"     // default

  // render.Ext = ".html"               // default

  // render.Debug = false               // default

  // render.TemplateFuncMap = template.FuncMap{}

  r.HTMLRender = render.Init()
  r.Run(":9000")
```

### Note

I hope this package will resolve your problem about template system.
and give you an idea about how to use **template helpers** in [Gin framework](https://github.com/gin-gonic/gin)

Thanks, [multitemplate.go](https://github.com/gin-gonic/contrib/tree/master/renders/multitemplate), [gin_html_render.go](https://gist.github.com/madhums/4340cbeb36871e227905) for the idea
