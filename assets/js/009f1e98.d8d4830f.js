"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[3436],{3905:function(e,n,t){t.d(n,{Zo:function(){return c},kt:function(){return m}});var a=t(7294);function r(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function i(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){r(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function s(e,n){if(null==e)return{};var t,a,r=function(e,n){if(null==e)return{};var t,a,r={},o=Object.keys(e);for(a=0;a<o.length;a++)t=o[a],n.indexOf(t)>=0||(r[t]=e[t]);return r}(e,n);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)t=o[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(r[t]=e[t])}return r}var l=a.createContext({}),p=function(e){var n=a.useContext(l),t=n;return e&&(t="function"==typeof e?e(n):i(i({},n),e)),t},c=function(e){var n=p(e.components);return a.createElement(l.Provider,{value:n},e.children)},u={inlineCode:"code",wrapper:function(e){var n=e.children;return a.createElement(a.Fragment,{},n)}},d=a.forwardRef((function(e,n){var t=e.components,r=e.mdxType,o=e.originalType,l=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),d=p(t),m=r,f=d["".concat(l,".").concat(m)]||d[m]||u[m]||o;return t?a.createElement(f,i(i({ref:n},c),{},{components:t})):a.createElement(f,i({ref:n},c))}));function m(e,n){var t=arguments,r=n&&n.mdxType;if("string"==typeof e||r){var o=t.length,i=new Array(o);i[0]=d;var s={};for(var l in n)hasOwnProperty.call(n,l)&&(s[l]=n[l]);s.originalType=e,s.mdxType="string"==typeof e?e:r,i[1]=s;for(var p=2;p<o;p++)i[p]=t[p];return a.createElement.apply(null,i)}return a.createElement.apply(null,t)}d.displayName="MDXCreateElement"},2783:function(e,n,t){t.r(n),t.d(n,{frontMatter:function(){return s},contentTitle:function(){return l},metadata:function(){return p},toc:function(){return c},default:function(){return d}});var a=t(7462),r=t(3366),o=(t(7294),t(3905)),i=["components"],s={id:"configuration",title:"Configurations"},l=void 0,p={unversionedId:"getting-started/configuration",id:"getting-started/configuration",isDocsHomePage:!1,title:"Configurations",description:"Optimus can be configured with .optimus.yaml file. An example of such is:",source:"@site/docs/getting-started/configuration.md",sourceDirName:"getting-started",slug:"/getting-started/configuration",permalink:"/optimus/docs/getting-started/configuration",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/getting-started/configuration.md",tags:[],version:"current",lastUpdatedBy:"Dery Rahman Ahaddienata",lastUpdatedAt:1647944766,formattedLastUpdatedAt:"3/22/2022",frontMatter:{id:"configuration",title:"Configurations"},sidebar:"docsSidebar",previous:{title:"Installation",permalink:"/optimus/docs/getting-started/installation"},next:{title:"Using Optimus to create a Job",permalink:"/optimus/docs/guides/create-job"}},c=[],u={toc:c};function d(e){var n=e.components,t=(0,r.Z)(e,i);return(0,o.kt)("wrapper",(0,a.Z)({},u,t,{components:n,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Optimus can be configured with ",(0,o.kt)("inlineCode",{parentName:"p"},".optimus.yaml")," file. An example of such is:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\n\n# used to connect optimus service\nhost: localhost:9100 \n\n# project specification\nproject:\n  \n  # name of the Optimus project\n  name: sample_project\n  \n  # project level variables usable in specifications\n  config: {}\n\n# namespace specification of the jobs and resources\nnamespace:\n  \n  # namespace name \n  name: sample_namespace\n  \n  jobs:\n    # folder where job specifications are stored\n    path: \"job\"\n    \n  datastore:\n    # optimus is capable of supporting multiple datastores\n    type: bigquery\n    # path where resource spec for BQ are stored\n    path: \"bq\"\n    # backup configurations of a datastore\n    backup:\n      # backup result age until expired - default '720h'\n      ttl: 168h\n      # where backup result should be located - default 'optimus_backup'\n      dataset: archive\n      # backup result prefix table name - default 'backup'\n      prefix: archive\n    \n  # namespace level variables usable in specifications\n  config: {}\n\n# for configuring optimus service locally\nserve:\n  \n  # port to listen on\n  port: 9100\n  \n  # host to listen on\n  host: localhost\n  \n  # this gets injected in compiled dags to reach back out to optimus service\n  # when they run\n  ingress_host: optimus.example.io:80\n  \n  # 32 char hash used for encrypting secrets\n  app_key: Yjo4a0jn1NvYdq79SADC/KaVv9Wu0Ffc\n  \n  # database configurations\n  db:\n    # database connection string\n    dsn: postgres://user:password@localhost:5432/database?sslmode=disable\n    \n    max_idle_connection: 5\n    max_open_connection: 10\n\n# logging configuration\nlog:\n  # debug, info, warning, error, fatal - default 'info'\n  level: debug  \n\n")),(0,o.kt)("p",null,"This configuration file should not be checked in version control. All the configs can also be passed as environment\nvariables using ",(0,o.kt)("inlineCode",{parentName:"p"},"OPTIMUS_<CONFIGNAME>")," convention, for example to set client host ",(0,o.kt)("inlineCode",{parentName:"p"},"OPTIMUS_HOST=localhost:9100")," to set\ndatabase connection string ",(0,o.kt)("inlineCode",{parentName:"p"},"OPTIMUS_SERVE_DB_DSN=postgres://dbconenctionurl"),"."),(0,o.kt)("p",null,"Assuming the following configuration layout:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"host: localhost:9100\nserve:\n  port: 9100\n  app_key: randomhash\n")),(0,o.kt)("p",null,"Key ",(0,o.kt)("inlineCode",{parentName:"p"},"host")," can be set as an environment variable by upper-casing its path, using ",(0,o.kt)("inlineCode",{parentName:"p"},"_")," as the\npath delimiter and prefixing with ",(0,o.kt)("inlineCode",{parentName:"p"},"OPTIMUS_"),":"),(0,o.kt)("p",null,(0,o.kt)("inlineCode",{parentName:"p"},"serve.port")," -> ",(0,o.kt)("inlineCode",{parentName:"p"},"OPTIMUS_SERVE_PORT=9100"),"\nor:\n",(0,o.kt)("inlineCode",{parentName:"p"},"host")," -> ",(0,o.kt)("inlineCode",{parentName:"p"},"OPTIMUS_HOST=localhost:9100")),(0,o.kt)("p",null,"Environment variables always override values from the configuration file. Here are some more examples:"),(0,o.kt)("table",null,(0,o.kt)("thead",{parentName:"table"},(0,o.kt)("tr",{parentName:"thead"},(0,o.kt)("th",{parentName:"tr",align:null},"Configuration key"),(0,o.kt)("th",{parentName:"tr",align:null},"Environment variable"))),(0,o.kt)("tbody",{parentName:"table"},(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"host"),(0,o.kt)("td",{parentName:"tr",align:null},"OPTIMUS_HOST")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},"serve.app_key"),(0,o.kt)("td",{parentName:"tr",align:null},"OPTIMUS_SERVE_APP_KEY")))),(0,o.kt)("p",null,"App key is used to encrypt credentials and can be randomly generated using"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-shell"},"head -c 50 /dev/random | base64\n")),(0,o.kt)("p",null,"Just take the first 32 characters of the string."),(0,o.kt)("p",null,"Configuration file can be stored in following locations:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-shell"},"./\n<exec>/\n~/.optimus/\n")))}d.isMDXComponent=!0}}]);