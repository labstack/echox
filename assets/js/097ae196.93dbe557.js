"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[3189],{9639:(e,n,i)=>{i.r(n),i.d(n,{assets:()=>c,contentTitle:()=>t,default:()=>p,frontMatter:()=>r,metadata:()=>d,toc:()=>a});var s=i(5893),o=i(1151);const r={description:"Decompress middleware"},t="Decompress",d={id:"middleware/decompress",title:"Decompress",description:"Decompress middleware",source:"@site/docs/middleware/decompress.md",sourceDirName:"middleware",slug:"/middleware/decompress",permalink:"/docs/middleware/decompress",draft:!1,unlisted:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/middleware/decompress.md",tags:[],version:"current",frontMatter:{description:"Decompress middleware"},sidebar:"docsSidebar",previous:{title:"CSRF",permalink:"/docs/middleware/csrf"},next:{title:"Gzip",permalink:"/docs/middleware/gzip"}},c={},a=[{value:"Usage",id:"usage",level:2},{value:"Custom Configuration",id:"custom-configuration",level:2},{value:"Usage",id:"usage-1",level:3},{value:"Configuration",id:"configuration",level:2},{value:"Default Configuration",id:"default-configuration",level:3}];function l(e){const n={admonition:"admonition",code:"code",h1:"h1",h2:"h2",h3:"h3",p:"p",pre:"pre",...(0,o.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h1,{id:"decompress",children:"Decompress"}),"\n",(0,s.jsx)(n.p,{children:"Decompress middleware decompresses HTTP request if Content-Encoding header is set to gzip."}),"\n",(0,s.jsx)(n.admonition,{type:"note",children:(0,s.jsx)(n.p,{children:"The body will be decompressed in memory and consume it for the lifetime of the request (and garbage collection)."})}),"\n",(0,s.jsx)(n.h2,{id:"usage",children:"Usage"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"e.Use(middleware.Decompress())\n"})}),"\n",(0,s.jsx)(n.h2,{id:"custom-configuration",children:"Custom Configuration"}),"\n",(0,s.jsx)(n.h3,{id:"usage-1",children:"Usage"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"e := echo.New()\ne.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{\n  Skipper: Skipper\n}))\n"})}),"\n",(0,s.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"DecompressConfig struct {\n  // Skipper defines a function to skip middleware.\n  Skipper Skipper\n}\n"})}),"\n",(0,s.jsx)(n.h3,{id:"default-configuration",children:"Default Configuration"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-go",children:"DefaultDecompressConfig = DecompressConfig{\n  Skipper: DefaultSkipper,\n}\n"})})]})}function p(e={}){const{wrapper:n}={...(0,o.a)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(l,{...e})}):l(e)}},1151:(e,n,i)=>{i.d(n,{Z:()=>d,a:()=>t});var s=i(7294);const o={},r=s.createContext(o);function t(e){const n=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function d(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:t(e.components),s.createElement(r.Provider,{value:n},e.children)}}}]);