"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[465],{4293:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>r,default:()=>u,frontMatter:()=>s,metadata:()=>o,toc:()=>d});var i=t(5893),a=t(1151);const s={description:"Basic auth middleware"},r="Basic Auth",o={id:"middleware/basic-auth",title:"Basic Auth",description:"Basic auth middleware",source:"@site/docs/middleware/basic-auth.md",sourceDirName:"middleware",slug:"/middleware/basic-auth",permalink:"/docs/middleware/basic-auth",draft:!1,unlisted:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/middleware/basic-auth.md",tags:[],version:"current",frontMatter:{description:"Basic auth middleware"},sidebar:"docsSidebar",previous:{title:"Middleware",permalink:"/docs/category/middleware"},next:{title:"Body Dump",permalink:"/docs/middleware/body-dump"}},c={},d=[{value:"Usage",id:"usage",level:2},{value:"Custom Configuration",id:"custom-configuration",level:2},{value:"Usage",id:"usage-1",level:3},{value:"Configuration",id:"configuration",level:2},{value:"Default Configuration",id:"default-configuration",level:3}];function l(e){const n={code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,a.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"basic-auth",children:"Basic Auth"}),"\n",(0,i.jsx)(n.p,{children:"Basic auth middleware provides an HTTP basic authentication."}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"For valid credentials it calls the next handler."}),"\n",(0,i.jsx)(n.li,{children:'For missing or invalid credentials, it sends "401 - Unauthorized" response.'}),"\n"]}),"\n",(0,i.jsx)(n.h2,{id:"usage",children:"Usage"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {\n\t// Be careful to use constant time comparison to prevent timing attacks\n\tif subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&\n\t\tsubtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {\n\t\treturn true, nil\n\t}\n\treturn false, nil\n}))\n'})}),"\n",(0,i.jsx)(n.h2,{id:"custom-configuration",children:"Custom Configuration"}),"\n",(0,i.jsx)(n.h3,{id:"usage-1",children:"Usage"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{}))\n"})}),"\n",(0,i.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'BasicAuthConfig struct {\n  // Skipper defines a function to skip middleware.\n  Skipper Skipper\n\n  // Validator is a function to validate BasicAuth credentials.\n  // Required.\n  Validator BasicAuthValidator\n\n  // Realm is a string to define realm attribute of BasicAuth.\n  // Default value "Restricted".\n  Realm string\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"default-configuration",children:"Default Configuration"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"DefaultBasicAuthConfig = BasicAuthConfig{\n\tSkipper: DefaultSkipper,\n}\n"})})]})}function u(e={}){const{wrapper:n}={...(0,a.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(l,{...e})}):l(e)}},1151:(e,n,t)=>{t.d(n,{Z:()=>o,a:()=>r});var i=t(7294);const a={},s=i.createContext(a);function r(e){const n=i.useContext(s);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:r(e.components),i.createElement(s.Provider,{value:n},e.children)}}}]);