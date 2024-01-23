"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[7246],{3740:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>l,contentTitle:()=>d,default:()=>h,frontMatter:()=>t,metadata:()=>a,toc:()=>c});var i=r(5893),s=r(1151);const t={description:"Binding request data",slug:"/binding",sidebar_position:3},d="Binding",a={id:"guide/binding",title:"Binding",description:"Binding request data",source:"@site/docs/guide/binding.md",sourceDirName:"guide",slug:"/binding",permalink:"/docs/binding",draft:!1,unlisted:!1,editUrl:"https://github.com/labstack/echox/blob/master/website/docs/guide/binding.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{description:"Binding request data",slug:"/binding",sidebar_position:3},sidebar:"docsSidebar",previous:{title:"Customization",permalink:"/docs/customization"},next:{title:"Context",permalink:"/docs/context"}},l={},c=[{value:"Struct Tag Binding",id:"struct-tag-binding",level:2},{value:"Data Sources",id:"data-sources",level:3},{value:"Data Types",id:"data-types",level:3},{value:"Multiple Sources",id:"multiple-sources",level:3},{value:"Direct Source",id:"direct-source",level:3},{value:"Security",id:"security",level:3},{value:"Example",id:"example",level:3},{value:"JSON Data",id:"json-data",level:4},{value:"Form Data",id:"form-data",level:4},{value:"Query Parameters",id:"query-parameters",level:4},{value:"Fluent Binding",id:"fluent-binding",level:2},{value:"Error Handling",id:"error-handling",level:3},{value:"Example",id:"example-1",level:3},{value:"Supported Data Types",id:"supported-data-types",level:3},{value:"Custom Binding",id:"custom-binding",level:2}];function o(e){const n={a:"a",code:"code",em:"em",h1:"h1",h2:"h2",h3:"h3",h4:"h4",li:"li",ol:"ol",p:"p",pre:"pre",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,s.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"binding",children:"Binding"}),"\n",(0,i.jsxs)(n.p,{children:["Parsing request data is a crucial part of a web application.  In Echo this is done with a process called ",(0,i.jsx)(n.em,{children:"binding"}),". This is done with information passed by the client in the following parts of an HTTP request:"]}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"URL Path parameter"}),"\n",(0,i.jsx)(n.li,{children:"URL Query parameter"}),"\n",(0,i.jsx)(n.li,{children:"Header"}),"\n",(0,i.jsx)(n.li,{children:"Request body"}),"\n"]}),"\n",(0,i.jsx)(n.p,{children:"Echo provides different ways to perform binding, each described in the sections below."}),"\n",(0,i.jsx)(n.h2,{id:"struct-tag-binding",children:"Struct Tag Binding"}),"\n",(0,i.jsxs)(n.p,{children:["With struct binding you define a Go struct with tags specifying the data source and corresponding key. In your request handler you simply call ",(0,i.jsx)(n.code,{children:"Context#Bind(i interface{})"})," with a pointer to your struct. The tags tell the binder everything it needs to know to load data from the request."]}),"\n",(0,i.jsxs)(n.p,{children:["In this example a struct type ",(0,i.jsx)(n.code,{children:"User"})," tells the binder to bind the query string parameter ",(0,i.jsx)(n.code,{children:"id"})," to its string field ",(0,i.jsx)(n.code,{children:"ID"}),":"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n  ID string `query:"id"`\n}\n\n// in the handler for /users?id=<userID>\nvar user User\nerr := c.Bind(&user); if err != nil {\n    return c.String(http.StatusBadRequest, "bad request")\n}\n'})}),"\n",(0,i.jsx)(n.h3,{id:"data-sources",children:"Data Sources"}),"\n",(0,i.jsx)(n.p,{children:"Echo supports the following tags specifying data sources:"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"query"})," - query parameter"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"param"})," - path parameter (also called route)"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"header"})," - header parameter"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"json"})," - request body. Uses builtin Go ",(0,i.jsx)(n.a,{href:"https://golang.org/pkg/encoding/json/",children:"json"})," package for unmarshalling."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"xml"})," - request body. Uses builtin Go ",(0,i.jsx)(n.a,{href:"https://golang.org/pkg/encoding/xml/",children:"xml"})," package for unmarshalling."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"form"})," - form data. Values are taken from query and request body. Uses Go standard library form parsing."]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"data-types",children:"Data Types"}),"\n",(0,i.jsxs)(n.p,{children:["When decoding the request body, the following data types are supported as specified by the ",(0,i.jsx)(n.code,{children:"Content-Type"})," header:"]}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:(0,i.jsx)(n.code,{children:"application/json"})}),"\n",(0,i.jsx)(n.li,{children:(0,i.jsx)(n.code,{children:"application/xml"})}),"\n",(0,i.jsx)(n.li,{children:(0,i.jsx)(n.code,{children:"application/x-www-form-urlencoded"})}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["When binding path parameter, query parameter, header, or form data, tags must be explicitly set on each struct field. However, JSON and XML binding is done on the struct field name if the tag is omitted. This is according to the behaviour of ",(0,i.jsx)(n.a,{href:"https://pkg.go.dev/encoding/json#Unmarshal",children:"Go's json package"}),"."]}),"\n",(0,i.jsxs)(n.p,{children:["For form data, Echo uses Go standard library form parsing. This parses form data from both the request URL and body if content type is not ",(0,i.jsx)(n.code,{children:"MIMEMultipartForm"}),". See documentation for ",(0,i.jsx)(n.a,{href:"https://golang.org/pkg/net/http/#Request.ParseForm",children:"non-MIMEMultipartForm"}),"and ",(0,i.jsx)(n.a,{href:"https://golang.org/pkg/net/http/#Request.ParseMultipartForm",children:"MIMEMultipartForm"})]}),"\n",(0,i.jsx)(n.h3,{id:"multiple-sources",children:"Multiple Sources"}),"\n",(0,i.jsx)(n.p,{children:"It is possible to specify multiple sources on the same field. In this case request data is bound in this order:"}),"\n",(0,i.jsxs)(n.ol,{children:["\n",(0,i.jsx)(n.li,{children:"Path parameters"}),"\n",(0,i.jsx)(n.li,{children:"Query parameters (only for GET/DELETE methods)"}),"\n",(0,i.jsx)(n.li,{children:"Request body"}),"\n"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n  ID string `param:"id" query:"id" form:"id" json:"id" xml:"id"`\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["Note that binding at each stage will overwrite data bound in a previous stage. This means if your JSON request contains the query param ",(0,i.jsx)(n.code,{children:"name=query"})," and body ",(0,i.jsx)(n.code,{children:'{"name": "body"}'})," then the result will be ",(0,i.jsx)(n.code,{children:'User{Name: "body"}'}),"."]}),"\n",(0,i.jsx)(n.h3,{id:"direct-source",children:"Direct Source"}),"\n",(0,i.jsx)(n.p,{children:"It is also possible to bind data directly from a specific source:"}),"\n",(0,i.jsx)(n.p,{children:"Request body:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"err := (&DefaultBinder{}).BindBody(c, &payload)\n"})}),"\n",(0,i.jsx)(n.p,{children:"Query parameters:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"err := (&DefaultBinder{}).BindQueryParams(c, &payload)\n"})}),"\n",(0,i.jsx)(n.p,{children:"Path parameters:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"err := (&DefaultBinder{}).BindPathParams(c, &payload)\n"})}),"\n",(0,i.jsx)(n.p,{children:"Header parameters:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"err := (&DefaultBinder{}).BindHeaders(c, &payload)\n"})}),"\n",(0,i.jsxs)(n.p,{children:["Note that headers is not one of the included sources with ",(0,i.jsx)(n.code,{children:"Context#Bind"}),". The only way to bind header data is by calling ",(0,i.jsx)(n.code,{children:"BindHeaders"})," directly."]}),"\n",(0,i.jsx)(n.h3,{id:"security",children:"Security"}),"\n",(0,i.jsx)(n.p,{children:"To keep your application secure, avoid passing bound structs directly to other methods if these structs contain fields that should not be bindable. It is advisable to have a separate struct for binding and map it explicitly to your business struct."}),"\n",(0,i.jsxs)(n.p,{children:["Consider what will happen if your bound struct has an Exported field ",(0,i.jsx)(n.code,{children:"IsAdmin bool"})," and the request body contains ",(0,i.jsx)(n.code,{children:'{IsAdmin: true, Name: "hacker"}'}),"."]}),"\n",(0,i.jsx)(n.h3,{id:"example",children:"Example"}),"\n",(0,i.jsxs)(n.p,{children:["In this example we define a ",(0,i.jsx)(n.code,{children:"User"})," struct type with field tags to bind from ",(0,i.jsx)(n.code,{children:"json"}),", ",(0,i.jsx)(n.code,{children:"form"}),", or ",(0,i.jsx)(n.code,{children:"query"})," request data:"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'type User struct {\n  Name  string `json:"name" form:"name" query:"name"`\n  Email string `json:"email" form:"email" query:"email"`\n}\n\ntype UserDTO struct {\n  Name    string\n  Email   string\n  IsAdmin bool\n}\n'})}),"\n",(0,i.jsxs)(n.p,{children:["And a handler at the POST ",(0,i.jsx)(n.code,{children:"/users"})," route binds request data to the struct:"]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'e.POST("/users", func(c echo.Context) (err error) {\n  u := new(User)\n  if err := c.Bind(u); err != nil {\n    return c.String(http.StatusBadRequest, "bad request")\n  }\n\n  // Load into separate struct for security\n  user := UserDTO{\n    Name: u.Name,\n    Email: u.Email,\n    IsAdmin: false // avoids exposing field that should not be bound\n  }\n\n  executeSomeBusinessLogic(user)\n  \n  return c.JSON(http.StatusOK, u)\n})\n'})}),"\n",(0,i.jsx)(n.h4,{id:"json-data",children:"JSON Data"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-sh",children:'curl -X POST http://localhost:1323/users \\\n  -H \'Content-Type: application/json\' \\\n  -d \'{"name":"Joe","email":"joe@labstack"}\'\n'})}),"\n",(0,i.jsx)(n.h4,{id:"form-data",children:"Form Data"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-sh",children:"curl -X POST http://localhost:1323/users \\\n  -d 'name=Joe' \\\n  -d 'email=joe@labstack.com'\n"})}),"\n",(0,i.jsx)(n.h4,{id:"query-parameters",children:"Query Parameters"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-sh",children:"curl -X GET 'http://localhost:1323/users?name=Joe&email=joe@labstack.com'\n"})}),"\n",(0,i.jsx)(n.h2,{id:"fluent-binding",children:"Fluent Binding"}),"\n",(0,i.jsxs)(n.p,{children:["Echo provides an interface to bind explicit data types from a specified source. It uses method chaining, also known as a ",(0,i.jsx)(n.a,{href:"https://en.wikipedia.org/wiki/Fluent_interface",children:"Fluent Interface"}),"."]}),"\n",(0,i.jsx)(n.p,{children:"The following methods provide a handful of methods for binding to Go data type. These binders offer a fluent syntax and can be chained to configure & execute binding, and handle errors."}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"echo.QueryParamsBinder(c)"})," - binds query parameters (source URL)"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"echo.PathParamsBinder(c)"})," - binds path parameters (source URL)"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"echo.FormFieldBinder(c)"})," - binds form fields (source URL + body). See also ",(0,i.jsx)(n.a,{href:"https://golang.org/pkg/net/http/#Request.ParseForm",children:"Request.ParseForm"}),"."]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"error-handling",children:"Error Handling"}),"\n",(0,i.jsxs)(n.p,{children:["A binder is usually completed by calling ",(0,i.jsx)(n.code,{children:"BindError()"})," or ",(0,i.jsx)(n.code,{children:"BindErrors()"}),". If any errors have occurred, ",(0,i.jsx)(n.code,{children:"BindError()"})," returns the first error encountered, while",(0,i.jsx)(n.code,{children:"BindErrors()"})," returns all bind errors. Any errors stored in the binder are also reset."]}),"\n",(0,i.jsxs)(n.p,{children:["With ",(0,i.jsx)(n.code,{children:"FailFast(true)"})," the binder can be configured to stop binding on the first error, or with ",(0,i.jsx)(n.code,{children:"FailFast(false)"})," execute the entire binder call chain. Fail fast is enabled by default and should be disabled when using ",(0,i.jsx)(n.code,{children:"BindErrors()"}),"."]}),"\n",(0,i.jsx)(n.h3,{id:"example-1",children:"Example"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'// url =  "/api/search?active=true&id=1&id=2&id=3&length=25"\nvar opts struct {\n  IDs []int64\n  Active bool\n}\nlength := int64(50) // default length is 50\n\n// creates query params binder that stops binding at first error\nerr := echo.QueryParamsBinder(c).\n  Int64("length", &length).\n  Int64s("ids", &opts.IDs).\n  Bool("active", &opts.Active).\n  BindError() // returns first binding error\n'})}),"\n",(0,i.jsx)(n.h3,{id:"supported-data-types",children:"Supported Data Types"}),"\n",(0,i.jsxs)(n.table,{children:[(0,i.jsx)(n.thead,{children:(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.th,{children:"Data Type"}),(0,i.jsx)(n.th,{children:"Notes"})]})}),(0,i.jsxs)(n.tbody,{children:[(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"bool"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"float32"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"float64"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"int"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"int8"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"int16"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"int32"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"int64"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"uint"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"uint8/byte"})}),(0,i.jsxs)(n.td,{children:["Does not support ",(0,i.jsx)(n.code,{children:"bytes()"}),". Use ",(0,i.jsx)(n.code,{children:"BindUnmarshaler"}),"/",(0,i.jsx)(n.code,{children:"CustomFunc"})," to convert value from base64 etc to ",(0,i.jsx)(n.code,{children:"[]byte{}"}),"."]})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"uint16"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"uint32"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"uint64"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"string"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"time"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"duration"})}),(0,i.jsx)(n.td,{})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"BindUnmarshaler()"})}),(0,i.jsx)(n.td,{children:"binds to a type implementing BindUnmarshaler interface"})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"TextUnmarshaler()"})}),(0,i.jsx)(n.td,{children:"binds to a type implementing encoding.TextUnmarshaler interface"})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"JsonUnmarshaler()"})}),(0,i.jsx)(n.td,{children:"binds to a type implementing json.Unmarshaler interface"})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"UnixTime()"})}),(0,i.jsxs)(n.td,{children:["converts Unix time (integer) to ",(0,i.jsx)(n.code,{children:"time.Time"})]})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"UnixTimeMilli()"})}),(0,i.jsxs)(n.td,{children:["converts Unix time with millisecond precision (integer) to ",(0,i.jsx)(n.code,{children:"time.Time"})]})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"UnixTimeNano()"})}),(0,i.jsxs)(n.td,{children:["converts Unix time with nanosecond precision (integer) to ",(0,i.jsx)(n.code,{children:"time.Time"})]})]}),(0,i.jsxs)(n.tr,{children:[(0,i.jsx)(n.td,{children:(0,i.jsx)(n.code,{children:"CustomFunc()"})}),(0,i.jsx)(n.td,{children:"callback function for your custom conversion logic"})]})]})]}),"\n",(0,i.jsx)(n.p,{children:"Each supported type has the following methods:"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:'<Type>("param", &destination)'})," - if parameter value exists then binds it to given destination of that type i.e ",(0,i.jsx)(n.code,{children:"Int64(...)"}),"."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:'Must<Type>("param", &destination)'})," - parameter value is required to exist, binds it to given destination of that type i.e ",(0,i.jsx)(n.code,{children:"MustInt64(...)"}),"."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:'<Type>s("param", &destination)'})," - (for slices) if parameter values exists then binds it to given destination of that type i.e ",(0,i.jsx)(n.code,{children:"Int64s(...)"}),"."]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:'Must<Type>s("param", &destination)'})," - (for slices) parameter value is required to exist, binds it to given destination of that type i.e ",(0,i.jsx)(n.code,{children:"MustInt64s(...)"}),"."]}),"\n"]}),"\n",(0,i.jsxs)(n.p,{children:["For certain slice types ",(0,i.jsx)(n.code,{children:'BindWithDelimiter("param", &dest, ",")'})," supports splitting parameter values before type conversion is done. For example binding an integer slice from the URL ",(0,i.jsx)(n.code,{children:"/api/search?id=1,2,3&id=1"})," will result in ",(0,i.jsx)(n.code,{children:"[]int64{1,2,3,1}"}),"."]}),"\n",(0,i.jsx)(n.h2,{id:"custom-binding",children:"Custom Binding"}),"\n",(0,i.jsxs)(n.p,{children:["A custom binder can be registered using ",(0,i.jsx)(n.code,{children:"Echo#Binder"}),"."]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:"type CustomBinder struct {}\n\nfunc (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {\n  // You may use default binder\n  db := new(echo.DefaultBinder)\n  if err := db.Bind(i, c); err != echo.ErrUnsupportedMediaType {\n    return\n  }\n\n  // Define your custom implementation here\n  return\n}\n"})})]})}function h(e={}){const{wrapper:n}={...(0,s.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(o,{...e})}):o(e)}},1151:(e,n,r)=>{r.d(n,{Z:()=>a,a:()=>d});var i=r(7294);const s={},t=i.createContext(s);function d(e){const n=i.useContext(t);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:d(e.components),i.createElement(t.Provider,{value:n},e.children)}}}]);