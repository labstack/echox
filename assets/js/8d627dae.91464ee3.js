"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[9960],{731:(e,t,o)=>{o.r(t),o.d(t,{assets:()=>a,contentTitle:()=>n,default:()=>h,frontMatter:()=>i,metadata:()=>s,toc:()=>d});const s=JSON.parse('{"id":"cookbook/auto-tls","title":"Auto TLS","description":"Automatic TLS certificates from Let\'s Encrypt recipe","source":"@site/docs/cookbook/auto-tls.md","sourceDirName":"cookbook","slug":"/cookbook/auto-tls","permalink":"/docs/cookbook/auto-tls","draft":false,"unlisted":false,"editUrl":"https://github.com/labstack/echox/blob/master/website/docs/cookbook/auto-tls.md","tags":[],"version":"current","frontMatter":{"description":"Automatic TLS certificates from Let\'s Encrypt recipe"},"sidebar":"docsSidebar","previous":{"title":"Cookbook","permalink":"/docs/category/cookbook"},"next":{"title":"CORS","permalink":"/docs/cookbook/cors"}}');var r=o(4848),c=o(8453);const i={description:"Automatic TLS certificates from Let's Encrypt recipe"},n="Auto TLS",a={},d=[{value:"Server",id:"server",level:2}];function l(e){const t={a:"a",admonition:"admonition",code:"code",h1:"h1",h2:"h2",header:"header",li:"li",p:"p",pre:"pre",ul:"ul",...(0,c.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(t.header,{children:(0,r.jsx)(t.h1,{id:"auto-tls",children:"Auto TLS"})}),"\n",(0,r.jsxs)(t.p,{children:["This recipe demonstrates how to obtain TLS certificates for a domain automatically from\nLet's Encrypt. ",(0,r.jsx)(t.code,{children:"Echo#StartAutoTLS"})," accepts an address which should listen on port ",(0,r.jsx)(t.code,{children:"443"}),"."]}),"\n",(0,r.jsxs)(t.p,{children:["Browse to ",(0,r.jsx)(t.code,{children:"https://<DOMAIN>"}),". If everything goes fine, you should see a welcome\nmessage with TLS enabled on the website."]}),"\n",(0,r.jsx)(t.admonition,{type:"tip",children:(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsx)(t.li,{children:"For added security you should specify host policy in auto TLS manager"}),"\n",(0,r.jsxs)(t.li,{children:["Cache certificates to avoid issues with rate limits (",(0,r.jsx)(t.a,{href:"https://letsencrypt.org/docs/rate-limits",children:"https://letsencrypt.org/docs/rate-limits"}),")"]}),"\n",(0,r.jsxs)(t.li,{children:["To redirect HTTP traffic to HTTPS, you can use ",(0,r.jsx)(t.a,{href:"../middleware/redirect#https-redirect",children:"redirect middleware"})]}),"\n"]})}),"\n",(0,r.jsx)(t.h2,{id:"server",children:"Server"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-go",metastring:"reference",children:"https://github.com/labstack/echox/blob/master/cookbook/auto-tls/server.go\n"})})]})}function h(e={}){const{wrapper:t}={...(0,c.R)(),...e.components};return t?(0,r.jsx)(t,{...e,children:(0,r.jsx)(l,{...e})}):l(e)}},8453:(e,t,o)=>{o.d(t,{R:()=>i,x:()=>n});var s=o(6540);const r={},c=s.createContext(r);function i(e){const t=s.useContext(c);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function n(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:i(e.components),s.createElement(c.Provider,{value:t},e.children)}}}]);