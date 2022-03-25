"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[8932],{3905:function(e,t,n){n.d(t,{Zo:function(){return c},kt:function(){return m}});var r=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function s(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function a(e,t){if(null==e)return{};var n,r,i=function(e,t){if(null==e)return{};var n,r,i={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var u=r.createContext({}),p=function(e){var t=r.useContext(u),n=t;return e&&(n="function"==typeof e?e(t):s(s({},t),e)),n},c=function(e){var t=p(e.components);return r.createElement(u.Provider,{value:t},e.children)},l={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,i=e.mdxType,o=e.originalType,u=e.parentName,c=a(e,["components","mdxType","originalType","parentName"]),d=p(n),m=i,f=d["".concat(u,".").concat(m)]||d[m]||l[m]||o;return n?r.createElement(f,s(s({ref:t},c),{},{components:n})):r.createElement(f,s({ref:t},c))}));function m(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=n.length,s=new Array(o);s[0]=d;var a={};for(var u in t)hasOwnProperty.call(t,u)&&(a[u]=t[u]);a.originalType=e,a.mdxType="string"==typeof e?e:i,s[1]=a;for(var p=2;p<o;p++)s[p]=n[p];return r.createElement.apply(null,s)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},4244:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return a},contentTitle:function(){return u},metadata:function(){return p},toc:function(){return c},default:function(){return d}});var r=n(7462),i=n(3366),o=(n(7294),n(3905)),s=["components"],a={id:"optimus-serve",title:"Starting Optimus Server"},u=void 0,p={unversionedId:"guides/optimus-serve",id:"guides/optimus-serve",isDocsHomePage:!1,title:"Starting Optimus Server",description:"Once the optimus binary is installed, it can be started in serve mode using",source:"@site/docs/guides/optimus-serve.md",sourceDirName:"guides",slug:"/guides/optimus-serve",permalink:"/optimus/docs/guides/optimus-serve",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/guides/optimus-serve.md",tags:[],version:"current",lastUpdatedBy:"Siddhanta Rath",lastUpdatedAt:1648199863,formattedLastUpdatedAt:"3/25/2022",frontMatter:{id:"optimus-serve",title:"Starting Optimus Server"},sidebar:"docsSidebar",previous:{title:"Organising specifications",permalink:"/optimus/docs/guides/organising-specifications"},next:{title:"Bigquery to bigquery transformation task",permalink:"/optimus/docs/guides/task-bq2bq"}},c=[],l={toc:c};function d(e){var t=e.components,n=(0,i.Z)(e,s);return(0,o.kt)("wrapper",(0,r.Z)({},l,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Once the optimus binary is installed, it can be started in serve mode using"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-shell"},"optimus serve\n")),(0,o.kt)("p",null,"It needs few ",(0,o.kt)("a",{parentName:"p",href:"/optimus/docs/getting-started/configuration"},"configurations")," as prerequisites, create a ",(0,o.kt)("inlineCode",{parentName:"p"},".optimus.yaml")," file with"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\nhost: localhost:9100\nserve:\n  port: 9100\n  host: localhost\n  ingress_host: optimus.example.io:80\n  app_key: 32charrandomhash32charrandomhash\n  db:\n    dsn: postgres://user:password@localhost:5432/optimus?sslmode=disable\n")),(0,o.kt)("p",null,"You will need to change ",(0,o.kt)("inlineCode",{parentName:"p"},"dsn")," and ",(0,o.kt)("inlineCode",{parentName:"p"},"app_key")," according to your installation."),(0,o.kt)("p",null,"Once the server is up and running, before it is ready to deploy ",(0,o.kt)("inlineCode",{parentName:"p"},"jobs")," we need to"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Register an optimus project"),(0,o.kt)("li",{parentName:"ul"},"Register a namespace under project"),(0,o.kt)("li",{parentName:"ul"},"Register required secrets under project")),(0,o.kt)("p",null,"This needs to be done in order using REST/GRPC endpoints provided by the server."))}d.isMDXComponent=!0}}]);