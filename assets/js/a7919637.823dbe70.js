"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[2375],{3073:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>a,contentTitle:()=>c,default:()=>h,frontMatter:()=>s,metadata:()=>i,toc:()=>d});var o=t(5893),r=t(1151);const s={description:"Context in Echo",slug:"/context",sidebar_position:4},c="Context",i={id:"guide/context",title:"Context",description:"Context in Echo",source:"@site/docs/guide/context.md",sourceDirName:"guide",slug:"/context",permalink:"/docs/context",draft:!1,unlisted:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/guide/context.md",tags:[],version:"current",sidebarPosition:4,frontMatter:{description:"Context in Echo",slug:"/context",sidebar_position:4},sidebar:"docsSidebar",previous:{title:"Binding",permalink:"/docs/binding"},next:{title:"Cookies",permalink:"/docs/cookies"}},a={},d=[{value:"Extending",id:"extending",level:2},{value:"Concurrency",id:"concurrency",level:2},{value:"Solution",id:"solution",level:3}];function l(e){const n={a:"a",admonition:"admonition",code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",...(0,r.a)(),...e.components};return(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)(n.header,{children:(0,o.jsx)(n.h1,{id:"context",children:"Context"})}),"\n",(0,o.jsxs)(n.p,{children:[(0,o.jsx)(n.code,{children:"echo.Context"})," represents the context of the current HTTP request. It holds request and\nresponse reference, path, path parameters, data, registered handler and APIs to read\nrequest and write response. As Context is an interface, it is easy to extend it with\ncustom APIs."]}),"\n",(0,o.jsx)(n.h2,{id:"extending",children:"Extending"}),"\n",(0,o.jsx)(n.p,{children:(0,o.jsx)(n.strong,{children:"Define a custom context"})}),"\n",(0,o.jsx)(n.pre,{children:(0,o.jsx)(n.code,{className:"language-go",children:'type CustomContext struct {\n\techo.Context\n}\n\nfunc (c *CustomContext) Foo() {\n\tprintln("foo")\n}\n\nfunc (c *CustomContext) Bar() {\n\tprintln("bar")\n}\n'})}),"\n",(0,o.jsx)(n.p,{children:(0,o.jsx)(n.strong,{children:"Create a middleware to extend default context"})}),"\n",(0,o.jsx)(n.pre,{children:(0,o.jsx)(n.code,{className:"language-go",children:"e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {\n\treturn func(c echo.Context) error {\n\t\tcc := &CustomContext{c}\n\t\treturn next(cc)\n\t}\n})\n"})}),"\n",(0,o.jsx)(n.admonition,{type:"caution",children:(0,o.jsx)(n.p,{children:"This middleware should be registered before any other middleware."})}),"\n",(0,o.jsx)(n.admonition,{type:"caution",children:(0,o.jsx)(n.p,{children:"Custom context cannot be defined in a middleware before the router ran (Pre)"})}),"\n",(0,o.jsx)(n.p,{children:(0,o.jsx)(n.strong,{children:"Use in handler"})}),"\n",(0,o.jsx)(n.pre,{children:(0,o.jsx)(n.code,{className:"language-go",children:'e.GET("/", func(c echo.Context) error {\n\tcc := c.(*CustomContext)\n\tcc.Foo()\n\tcc.Bar()\n\treturn cc.String(200, "OK")\n})\n'})}),"\n",(0,o.jsx)(n.h2,{id:"concurrency",children:"Concurrency"}),"\n",(0,o.jsxs)(n.admonition,{type:"caution",children:[(0,o.jsxs)(n.p,{children:[(0,o.jsx)(n.code,{children:"Context"})," must not be accessed out of the goroutine handling the request. There are two reasons:"]}),(0,o.jsxs)(n.ol,{children:["\n",(0,o.jsxs)(n.li,{children:[(0,o.jsx)(n.code,{children:"Context"})," has functions that are dangerous to execute from multiple goroutines. Therefore, only one goroutine should access it."]}),"\n",(0,o.jsxs)(n.li,{children:["Echo uses a pool to create ",(0,o.jsx)(n.code,{children:"Context"}),"'s. When the request handling finishes, Echo returns the ",(0,o.jsx)(n.code,{children:"Context"})," to the pool."]}),"\n"]}),(0,o.jsxs)(n.p,{children:["See issue ",(0,o.jsx)(n.a,{href:"https://github.com/labstack/echo/issues/1908",children:"1908"}),' for a "cautionary tale" caused by this reason. Concurrency is complicated. Beware of this pitfall when working with goroutines.']})]}),"\n",(0,o.jsx)(n.h3,{id:"solution",children:"Solution"}),"\n",(0,o.jsx)(n.p,{children:"Use a channel"}),"\n",(0,o.jsx)(n.pre,{children:(0,o.jsx)(n.code,{className:"language-go",children:'func(c echo.Context) error {\n\tca := make(chan string, 1) // To prevent this channel from blocking, size is set to 1.\n\tr := c.Request()\n\tmethod := r.Method\n\n\tgo func() {\n\t\t// This function must not touch the Context.\n\n\t\tfmt.Printf("Method: %s\\n", method)\n\n\t\t// Do some long running operations...\n\n\t\tca <- "Hey!"\n\t}()\n\n\tselect {\n\tcase result := <-ca:\n\t\treturn c.String(http.StatusOK, "Result: "+result)\n\tcase <-c.Request().Context().Done(): // Check context.\n\t\t// If it reaches here, this means that context was canceled (a timeout was reached, etc.).\n\t\treturn nil\n\t}\n}\n'})})]})}function h(e={}){const{wrapper:n}={...(0,r.a)(),...e.components};return n?(0,o.jsx)(n,{...e,children:(0,o.jsx)(l,{...e})}):l(e)}},1151:(e,n,t)=>{t.d(n,{Z:()=>i,a:()=>c});var o=t(7294);const r={},s=o.createContext(r);function c(e){const n=o.useContext(s);return o.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function i(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:c(e.components),o.createElement(s.Provider,{value:n},e.children)}}}]);