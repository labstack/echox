"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[3740],{7605:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>c,contentTitle:()=>r,default:()=>h,frontMatter:()=>i,metadata:()=>a,toc:()=>l});var t=s(5893),o=s(1151);const i={description:"Session middleware"},r="Session",a={id:"middleware/session",title:"Session",description:"Session middleware",source:"@site/docs/middleware/session.md",sourceDirName:"middleware",slug:"/middleware/session",permalink:"/docs/middleware/session",draft:!1,unlisted:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/middleware/session.md",tags:[],version:"current",frontMatter:{description:"Session middleware"},sidebar:"docsSidebar",previous:{title:"Secure",permalink:"/docs/middleware/secure"},next:{title:"Static",permalink:"/docs/middleware/static"}},c={},l=[{value:"Dependencies",id:"dependencies",level:2},{value:"Usage",id:"usage",level:2},{value:"Example usage",id:"example-usage",level:3},{value:"Custom Configuration",id:"custom-configuration",level:2},{value:"Usage",id:"usage-1",level:3},{value:"Configuration",id:"configuration",level:2},{value:"Default Configuration",id:"default-configuration",level:3}];function d(e){const n={a:"a",admonition:"admonition",code:"code",h1:"h1",h2:"h2",h3:"h3",p:"p",pre:"pre",...(0,o.a)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(n.h1,{id:"session",children:"Session"}),"\n",(0,t.jsxs)(n.p,{children:["Session middleware facilitates HTTP session management backed by ",(0,t.jsx)(n.a,{href:"https://github.com/gorilla/sessions",children:"gorilla sessions"}),". The default implementation provides cookie and\nfilesystem based session store; however, you can take advantage of ",(0,t.jsx)(n.a,{href:"https://github.com/gorilla/sessions#store-implementations",children:"community maintained\nimplementation"})," for various backends."]}),"\n",(0,t.jsx)(n.admonition,{type:"note",children:(0,t.jsx)(n.p,{children:"Echo community contribution"})}),"\n",(0,t.jsx)(n.h2,{id:"dependencies",children:"Dependencies"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-go",children:'import (\n  "github.com/gorilla/sessions"\n  "github.com/labstack/echo-contrib/session"\n)\n'})}),"\n",(0,t.jsx)(n.h2,{id:"usage",children:"Usage"}),"\n",(0,t.jsxs)(n.p,{children:["This example exposes two endpoints: ",(0,t.jsx)(n.code,{children:"/create-session"})," creates new session and ",(0,t.jsx)(n.code,{children:"/read-session"})," read value from\nsession if request contains session id."]}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-go",children:'import (\n    "errors"\n    "fmt"\n    "github.com/gorilla/sessions"\n    "github.com/labstack/echo-contrib/session"\n    "github.com/labstack/echo/v4"\n    "log"\n    "net/http"\n)\n\nfunc main() {\n\te := echo.New()\n\te.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))\n\n\te.GET("/create-session", func(c echo.Context) error {\n\t\tsess, err := session.Get("session", c)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tsess.Options = &sessions.Options{\n\t\t\tPath:     "/",\n\t\t\tMaxAge:   86400 * 7,\n\t\t\tHttpOnly: true,\n\t\t}\n\t\tsess.Values["foo"] = "bar"\n\t\tif err := sess.Save(c.Request(), c.Response()); err != nil {\n\t\t\treturn err\n\t\t}\n\t\treturn c.NoContent(http.StatusOK)\n\t})\n\n\te.GET("/read-session", func(c echo.Context) error {\n\t\tsess, err := session.Get("session", c)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\treturn c.String(http.StatusOK, fmt.Sprintf("foo=%v\\n", sess.Values["foo"]))\n\t})\n\n\tif err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {\n\t\tlog.Fatal(err)\n\t}\n}\n'})}),"\n",(0,t.jsx)(n.h3,{id:"example-usage",children:"Example usage"}),"\n",(0,t.jsxs)(n.p,{children:["Requesting ",(0,t.jsx)(n.code,{children:"/read-session"})," without providing session it will output nil as ",(0,t.jsx)(n.code,{children:"foo"})," value"]}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-bash",children:"x@x:~/$ curl -v http://localhost:8080/read-session\n* processing: http://localhost:8080/read-session\n*   Trying [::1]:8080...\n* Connected to localhost (::1) port 8080\n> GET /read-session HTTP/1.1\n> Host: localhost:8080\n> User-Agent: curl/8.2.1\n> Accept: */*\n> \n< HTTP/1.1 200 OK\n< Content-Type: text/plain; charset=UTF-8\n< Date: Thu, 25 Apr 2024 09:15:14 GMT\n< Content-Length: 10\n< \nfoo=<nil>\n"})}),"\n",(0,t.jsxs)(n.p,{children:["Requesting ",(0,t.jsx)(n.code,{children:"/create-session"})," creates new session"]}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-bash",children:'x@x:~/$ curl -v -c cookies.txt http://localhost:8080/create-session\n* processing: http://localhost:8080/create-session\n*   Trying [::1]:8080...\n* Connected to localhost (::1) port 8080\n> GET /create-session HTTP/1.1\n> Host: localhost:8080\n> User-Agent: curl/8.2.1\n> Accept: */*\n> \n< HTTP/1.1 200 OK\n* Added cookie session="MTcxNDAzNjYyMHxEWDhFQVFMX2dBQUJFQUVRQUFBZ180QUFBUVp6ZEhKcGJtY01CUUFEWm05dkJuTjBjbWx1Wnd3RkFBTmlZWEk9fHJQxR5fJDUEV-6iHSWuyVzjYX2f9F5tVaMGV6pjIE1Y" for domain localhost, path /, expire 1714641420\n< Set-Cookie: session=MTcxNDAzNjYyMHxEWDhFQVFMX2dBQUJFQUVRQUFBZ180QUFBUVp6ZEhKcGJtY01CUUFEWm05dkJuTjBjbWx1Wnd3RkFBTmlZWEk9fHJQxR5fJDUEV-6iHSWuyVzjYX2f9F5tVaMGV6pjIE1Y; Path=/; Expires=Thu, 02 May 2024 09:17:00 GMT; Max-Age=604800; HttpOnly\n< Date: Thu, 25 Apr 2024 09:17:00 GMT\n< Content-Length: 0\n<\n* Connection #0 to host localhost left intact\n'})}),"\n",(0,t.jsxs)(n.p,{children:["Using session cookie from previous response and requesting ",(0,t.jsx)(n.code,{children:"/read-session"})," will output ",(0,t.jsx)(n.code,{children:"foo"})," value from session."]}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-bash",children:"x@x:~/$ curl -v -b cookies.txt http://localhost:8080/read-session\n* processing: http://localhost:8080/read-session\n*   Trying [::1]:8080...\n* Connected to localhost (::1) port 8080\n> GET /read-session HTTP/1.1\n> Host: localhost:8080\n> User-Agent: curl/8.2.1\n> Accept: */*\n> Cookie: session=MTcxNDAzNjYyMHxEWDhFQVFMX2dBQUJFQUVRQUFBZ180QUFBUVp6ZEhKcGJtY01CUUFEWm05dkJuTjBjbWx1Wnd3RkFBTmlZWEk9fHJQxR5fJDUEV-6iHSWuyVzjYX2f9F5tVaMGV6pjIE1Y\n> \n< HTTP/1.1 200 OK\n< Content-Type: text/plain; charset=UTF-8\n< Date: Thu, 25 Apr 2024 09:18:56 GMT\n< Content-Length: 8\n< \nfoo=bar\n* Connection #0 to host localhost left intact\n"})}),"\n",(0,t.jsx)(n.h2,{id:"custom-configuration",children:"Custom Configuration"}),"\n",(0,t.jsx)(n.h3,{id:"usage-1",children:"Usage"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-go",children:"e := echo.New()\ne.Use(session.MiddlewareWithConfig(session.Config{}))\n"})}),"\n",(0,t.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-go",children:"Config struct {\n  // Skipper defines a function to skip middleware.\n  Skipper middleware.Skipper\n\n  // Session store.\n  // Required.\n  Store sessions.Store\n}\n"})}),"\n",(0,t.jsx)(n.h3,{id:"default-configuration",children:"Default Configuration"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-go",children:"DefaultConfig = Config{\n  Skipper: DefaultSkipper,\n}\n"})})]})}function h(e={}){const{wrapper:n}={...(0,o.a)(),...e.components};return n?(0,t.jsx)(n,{...e,children:(0,t.jsx)(d,{...e})}):d(e)}},1151:(e,n,s)=>{s.d(n,{Z:()=>a,a:()=>r});var t=s(7294);const o={},i=t.createContext(o);function r(e){const n=t.useContext(i);return t.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:r(e.components),t.createElement(i.Provider,{value:n},e.children)}}}]);