"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[8618],{8833:(e,i,n)=>{n.r(i),n.d(i,{assets:()=>d,contentTitle:()=>s,default:()=>h,frontMatter:()=>r,metadata:()=>t,toc:()=>o});var a=n(5893),l=n(1151);const r={description:"Trailing slash middleware"},s="Trailing Slash",t={id:"middleware/trailing-slash",title:"Trailing Slash",description:"Trailing slash middleware",source:"@site/docs/middleware/trailing-slash.md",sourceDirName:"middleware",slug:"/middleware/trailing-slash",permalink:"/docs/middleware/trailing-slash",draft:!1,unlisted:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/middleware/trailing-slash.md",tags:[],version:"current",frontMatter:{description:"Trailing slash middleware"},sidebar:"docsSidebar",previous:{title:"Timeout",permalink:"/docs/middleware/timeout"},next:{title:"Cookbook",permalink:"/docs/category/cookbook"}},d={},o=[{value:"Add Trailing Slash",id:"add-trailing-slash",level:2},{value:"Usage",id:"usage",level:3},{value:"Remove Trailing Slash",id:"remove-trailing-slash",level:2},{value:"Usage",id:"usage-1",level:3},{value:"Custom Configuration",id:"custom-configuration",level:2},{value:"Usage",id:"usage-2",level:3},{value:"Configuration",id:"configuration",level:2},{value:"Default Configuration",id:"default-configuration",level:3}];function c(e){const i={code:"code",h1:"h1",h2:"h2",h3:"h3",header:"header",p:"p",pre:"pre",...(0,l.a)(),...e.components};return(0,a.jsxs)(a.Fragment,{children:[(0,a.jsx)(i.header,{children:(0,a.jsx)(i.h1,{id:"trailing-slash",children:"Trailing Slash"})}),"\n",(0,a.jsx)(i.h2,{id:"add-trailing-slash",children:"Add Trailing Slash"}),"\n",(0,a.jsx)(i.p,{children:"Add trailing slash middleware adds a trailing slash to the request URI."}),"\n",(0,a.jsx)(i.h3,{id:"usage",children:"Usage"}),"\n",(0,a.jsx)(i.pre,{children:(0,a.jsx)(i.code,{className:"language-go",children:"e := echo.New()\ne.Pre(middleware.AddTrailingSlash())\n"})}),"\n",(0,a.jsx)(i.h2,{id:"remove-trailing-slash",children:"Remove Trailing Slash"}),"\n",(0,a.jsx)(i.p,{children:"Remove trailing slash middleware removes a trailing slash from the request URI."}),"\n",(0,a.jsx)(i.h3,{id:"usage-1",children:"Usage"}),"\n",(0,a.jsx)(i.pre,{children:(0,a.jsx)(i.code,{className:"language-go",children:"e := echo.New()\ne.Pre(middleware.RemoveTrailingSlash())\n"})}),"\n",(0,a.jsx)(i.h2,{id:"custom-configuration",children:"Custom Configuration"}),"\n",(0,a.jsx)(i.h3,{id:"usage-2",children:"Usage"}),"\n",(0,a.jsx)(i.pre,{children:(0,a.jsx)(i.code,{className:"language-go",children:"e := echo.New()\ne.Use(middleware.AddTrailingSlashWithConfig(middleware.TrailingSlashConfig{\n  RedirectCode: http.StatusMovedPermanently,\n}))\n"})}),"\n",(0,a.jsxs)(i.p,{children:["Example above will add a trailing slash to the request URI and redirect with ",(0,a.jsx)(i.code,{children:"301 - StatusMovedPermanently"}),"."]}),"\n",(0,a.jsx)(i.h2,{id:"configuration",children:"Configuration"}),"\n",(0,a.jsx)(i.pre,{children:(0,a.jsx)(i.code,{className:"language-go",children:'TrailingSlashConfig struct {\n  // Skipper defines a function to skip middleware.\n  Skipper Skipper\n\n  // Status code to be used when redirecting the request.\n  // Optional, but when provided the request is redirected using this code.\n  RedirectCode int `json:"redirect_code"`\n}\n'})}),"\n",(0,a.jsx)(i.h3,{id:"default-configuration",children:"Default Configuration"}),"\n",(0,a.jsx)(i.pre,{children:(0,a.jsx)(i.code,{className:"language-go",children:"DefaultTrailingSlashConfig = TrailingSlashConfig{\n  Skipper: DefaultSkipper,\n}\n"})})]})}function h(e={}){const{wrapper:i}={...(0,l.a)(),...e.components};return i?(0,a.jsx)(i,{...e,children:(0,a.jsx)(c,{...e})}):c(e)}},1151:(e,i,n)=>{n.d(i,{Z:()=>t,a:()=>s});var a=n(7294);const l={},r=a.createContext(l);function s(e){const i=a.useContext(r);return a.useMemo((function(){return"function"==typeof e?e(i):{...i,...e}}),[i,e])}function t(e){let i;return i=e.disableParentContext?"function"==typeof e.components?e.components(l):e.components||l:s(e.components),a.createElement(r.Provider,{value:i},e.children)}}}]);