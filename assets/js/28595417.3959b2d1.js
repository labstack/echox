"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[7404],{3905:(e,t,n)=>{n.d(t,{Zo:()=>u,kt:()=>f});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var c=r.createContext({}),l=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},u=function(e){var t=l(e.components);return r.createElement(c.Provider,{value:t},e.children)},d="mdxType",p={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,c=e.parentName,u=s(e,["components","mdxType","originalType","parentName"]),d=l(n),m=a,f=d["".concat(c,".").concat(m)]||d[m]||p[m]||i;return n?r.createElement(f,o(o({ref:t},u),{},{components:n})):r.createElement(f,o({ref:t},u))}));function f(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,o=new Array(i);o[0]=m;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s[d]="string"==typeof e?e:a,o[1]=s;for(var l=2;l<i;l++)o[l]=n[l];return r.createElement.apply(null,o)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},9612:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>p,frontMatter:()=>i,metadata:()=>s,toc:()=>l});var r=n(7462),a=(n(7294),n(3905));const i={description:"Casbin auth middleware"},o="Casbin Auth",s={unversionedId:"middleware/casbin-auth",id:"middleware/casbin-auth",title:"Casbin Auth",description:"Casbin auth middleware",source:"@site/docs/middleware/casbin-auth.md",sourceDirName:"middleware",slug:"/middleware/casbin-auth",permalink:"/docs/middleware/casbin-auth",draft:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/middleware/casbin-auth.md",tags:[],version:"current",frontMatter:{description:"Casbin auth middleware"},sidebar:"docsSidebar",previous:{title:"Body Limit",permalink:"/docs/middleware/body-limit"},next:{title:"CORS",permalink:"/docs/middleware/cors"}},c={},l=[{value:"Dependencies",id:"dependencies",level:2},{value:"Usage",id:"usage",level:2},{value:"Custom Configuration",id:"custom-configuration",level:2},{value:"Usage",id:"usage-1",level:3},{value:"Configuration",id:"configuration",level:2},{value:"Default Configuration",id:"default-configuration",level:3}],u={toc:l},d="wrapper";function p(e){let{components:t,...n}=e;return(0,a.kt)(d,(0,r.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"casbin-auth"},"Casbin Auth"),(0,a.kt)("admonition",{type:"note"},(0,a.kt)("p",{parentName:"admonition"},"Echo community contribution")),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"https://github.com/casbin/casbin"},"Casbin")," is a powerful and efficient open-source access control library for Go. It provides support for enforcing authorization based on various models. So far, the access control models supported by Casbin are:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},"ACL (Access Control List)"),(0,a.kt)("li",{parentName:"ul"},"ACL with superuser"),(0,a.kt)("li",{parentName:"ul"},"ACL without users: especially useful for systems that don't have authentication or user log-ins."),(0,a.kt)("li",{parentName:"ul"},"ACL without resources: some scenarios may target for a type of resources instead of an individual resource by using permissions like write-article, read-log. It doesn't control the access to a specific article or log."),(0,a.kt)("li",{parentName:"ul"},"RBAC (Role-Based Access Control)"),(0,a.kt)("li",{parentName:"ul"},"RBAC with resource roles: both users and resources can have roles (or groups) at the same time."),(0,a.kt)("li",{parentName:"ul"},"RBAC with domains/tenants: users can have different role sets for different domains/tenants."),(0,a.kt)("li",{parentName:"ul"},"ABAC (Attribute-Based Access Control)"),(0,a.kt)("li",{parentName:"ul"},"RESTful"),(0,a.kt)("li",{parentName:"ul"},"Deny-override: both allow and deny authorizations are supported, deny overrides the allow.")),(0,a.kt)("admonition",{type:"info"},(0,a.kt)("p",{parentName:"admonition"},"Currently, only HTTP basic authentication is supported.")),(0,a.kt)("h2",{id:"dependencies"},"Dependencies"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},'import (\n  "github.com/casbin/casbin"\n  casbin_mw "github.com/labstack/echo-contrib/casbin"\n)\n')),(0,a.kt)("h2",{id:"usage"},"Usage"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},'e := echo.New()\nenforcer, err := casbin.NewEnforcer("casbin_auth_model.conf", "casbin_auth_policy.csv")\ne.Use(casbin_mw.Middleware(enforcer))\n')),(0,a.kt)("p",null,"For syntax, see: ",(0,a.kt)("a",{parentName:"p",href:"https://casbin.org/docs/syntax-for-models"},"Syntax for Models"),"."),(0,a.kt)("h2",{id:"custom-configuration"},"Custom Configuration"),(0,a.kt)("h3",{id:"usage-1"},"Usage"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},'e := echo.New()\nce := casbin.NewEnforcer("casbin_auth_model.conf", "")\nce.AddRoleForUser("alice", "admin")\nce.AddPolicy(...)\ne.Use(casbin_mw.MiddlewareWithConfig(casbin_mw.Config{\n  Enforcer: ce,\n}))\n')),(0,a.kt)("h2",{id:"configuration"},"Configuration"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},"// Config defines the config for CasbinAuth middleware.\nConfig struct {\n  // Skipper defines a function to skip middleware.\n  Skipper middleware.Skipper\n\n  // Enforcer CasbinAuth main rule.\n  // Required.\n  Enforcer *casbin.Enforcer\n}\n")),(0,a.kt)("h3",{id:"default-configuration"},"Default Configuration"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},"// DefaultConfig is the default CasbinAuth middleware config.\nDefaultConfig = Config{\n  Skipper: middleware.DefaultSkipper,\n}\n")))}p.isMDXComponent=!0}}]);