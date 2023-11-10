"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[2590],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>g});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var u=r.createContext({}),p=function(e){var t=r.useContext(u),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},d=function(e){var t=p(e.components);return r.createElement(u.Provider,{value:t},e.children)},s="mdxType",m={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},c=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,u=e.parentName,d=o(e,["components","mdxType","originalType","parentName"]),s=p(n),c=a,g=s["".concat(u,".").concat(c)]||s[c]||m[c]||i;return n?r.createElement(g,l(l({ref:t},d),{},{components:n})):r.createElement(g,l({ref:t},d))}));function g(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,l=new Array(i);l[0]=c;var o={};for(var u in t)hasOwnProperty.call(t,u)&&(o[u]=t[u]);o.originalType=e,o[s]="string"==typeof e?e:a,l[1]=o;for(var p=2;p<i;p++)l[p]=n[p];return r.createElement.apply(null,l)}return r.createElement.apply(null,n)}c.displayName="MDXCreateElement"},7161:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>u,contentTitle:()=>l,default:()=>m,frontMatter:()=>i,metadata:()=>o,toc:()=>p});var r=n(7462),a=(n(7294),n(3905));const i={description:"Customization",slug:"/customization",sidebar_position:2},l="Customization",o={unversionedId:"guide/customization",id:"guide/customization",title:"Customization",description:"Customization",source:"@site/docs/guide/customization.md",sourceDirName:"guide",slug:"/customization",permalink:"/docs/customization",draft:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/guide/customization.md",tags:[],version:"current",sidebarPosition:2,frontMatter:{description:"Customization",slug:"/customization",sidebar_position:2},sidebar:"docsSidebar",previous:{title:"Quick Start",permalink:"/docs/quick-start"},next:{title:"Binding",permalink:"/docs/binding"}},u={},p=[{value:"Debug",id:"debug",level:2},{value:"Logging",id:"logging",level:2},{value:"Log Header",id:"log-header",level:3},{value:"Available Tags",id:"available-tags",level:4},{value:"Log Output",id:"log-output",level:3},{value:"Log Level",id:"log-level",level:3},{value:"Custom Logger",id:"custom-logger",level:3},{value:"Startup Banner",id:"startup-banner",level:2},{value:"Custom Listener",id:"custom-listener",level:2},{value:"Disable HTTP/2",id:"disable-http2",level:2},{value:"Read Timeout",id:"read-timeout",level:2},{value:"Write Timeout",id:"write-timeout",level:2},{value:"Validator",id:"validator",level:2},{value:"Custom Binder",id:"custom-binder",level:2},{value:"Custom JSON Serializer",id:"custom-json-serializer",level:2},{value:"Renderer",id:"renderer",level:2},{value:"HTTP Error Handler",id:"http-error-handler",level:2}],d={toc:p},s="wrapper";function m(e){let{components:t,...n}=e;return(0,a.kt)(s,(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"customization"},"Customization"),(0,a.kt)("h2",{id:"debug"},"Debug"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Debug")," can be used to enable / disable debug mode. Debug mode sets the log level\nto ",(0,a.kt)("inlineCode",{parentName:"p"},"DEBUG"),"."),(0,a.kt)("h2",{id:"logging"},"Logging"),(0,a.kt)("p",null,"The default format for logging is JSON, which can be changed by modifying the header."),(0,a.kt)("h3",{id:"log-header"},"Log Header"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Logger.SetHeader(string)")," can be used to set the header for\nthe logger. Default value:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-js"},'{"time":"${time_rfc3339_nano}","level":"${level}","prefix":"${prefix}","file":"${short_file}","line":"${line}"}\n')),(0,a.kt)("p",null,(0,a.kt)("em",{parentName:"p"},"Example")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},'import "github.com/labstack/gommon/log"\n\n/* ... */\n\nif l, ok := e.Logger.(*log.Logger); ok {\n  l.SetHeader("${time_rfc3339} ${level}")\n}\n')),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-sh"},"2018-05-08T20:30:06-07:00 INFO info\n")),(0,a.kt)("h4",{id:"available-tags"},"Available Tags"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"time_rfc3339")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"time_rfc3339_nano")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"level")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"prefix")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"long_file")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"short_file")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"line"))),(0,a.kt)("h3",{id:"log-output"},"Log Output"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Logger.SetOutput(io.Writer)")," can be used to set the output destination for\nthe logger. Default value is ",(0,a.kt)("inlineCode",{parentName:"p"},"os.Stdout")),(0,a.kt)("p",null,"To completely disable logs use ",(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Logger.SetOutput(io.Discard)")," or ",(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Logger.SetLevel(log.OFF)")),(0,a.kt)("h3",{id:"log-level"},"Log Level"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Logger.SetLevel(log.Lvl)")," can be used to set the log level for the logger.\nDefault value is ",(0,a.kt)("inlineCode",{parentName:"p"},"ERROR"),". Possible values:"),(0,a.kt)("ul",null,(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"DEBUG")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"INFO")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"WARN")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"ERROR")),(0,a.kt)("li",{parentName:"ul"},(0,a.kt)("inlineCode",{parentName:"li"},"OFF"))),(0,a.kt)("h3",{id:"custom-logger"},"Custom Logger"),(0,a.kt)("p",null,"Logging is implemented using ",(0,a.kt)("inlineCode",{parentName:"p"},"echo.Logger")," interface which allows you to register\na custom logger using ",(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Logger"),"."),(0,a.kt)("h2",{id:"startup-banner"},"Startup Banner"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#HideBanner")," can be used to hide the startup banner."),(0,a.kt)("h2",{id:"custom-listener"},"Custom Listener"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#*Listener")," can be used to run a custom listener."),(0,a.kt)("p",null,(0,a.kt)("em",{parentName:"p"},"Example")),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-go"},'l, err := net.Listen("tcp", ":1323")\nif err != nil {\n  e.Logger.Fatal(err)\n}\ne.Listener = l\ne.Logger.Fatal(e.Start(""))\n')),(0,a.kt)("h2",{id:"disable-http2"},"Disable HTTP/2"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#DisableHTTP2")," can be used disable HTTP/2 protocol."),(0,a.kt)("h2",{id:"read-timeout"},"Read Timeout"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#*Server#ReadTimeout")," can be used to set the maximum duration before timing out read\nof the request."),(0,a.kt)("h2",{id:"write-timeout"},"Write Timeout"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#*Server#WriteTimeout")," can be used to set the maximum duration before timing out write\nof the response."),(0,a.kt)("h2",{id:"validator"},"Validator"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Validator")," can be used to register a validator for performing data validation\non request payload."),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"/docs/request#validate-data"},"Learn more")),(0,a.kt)("h2",{id:"custom-binder"},"Custom Binder"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Binder")," can be used to register a custom binder for binding request payload."),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"/docs/request/#custom-binder"},"Learn more")),(0,a.kt)("h2",{id:"custom-json-serializer"},"Custom JSON Serializer"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#JSONSerializer")," can be used to register a custom JSON serializer."),(0,a.kt)("p",null,"Have a look at ",(0,a.kt)("inlineCode",{parentName:"p"},"DefaultJSONSerializer")," on ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/labstack/echo/blob/master/json.go"},"json.go"),"."),(0,a.kt)("h2",{id:"renderer"},"Renderer"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#Renderer")," can be used to register a renderer for template rendering."),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"/docs/templates"},"Learn more")),(0,a.kt)("h2",{id:"http-error-handler"},"HTTP Error Handler"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"Echo#HTTPErrorHandler")," can be used to register a custom http error handler."),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"/docs/error-handling"},"Learn more")))}m.isMDXComponent=!0}}]);