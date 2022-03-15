"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[2082],{3905:function(e,n,t){t.d(n,{Zo:function(){return c},kt:function(){return m}});var r=t(7294);function a(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function i(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){a(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function s(e,n){if(null==e)return{};var t,r,a=function(e,n){if(null==e)return{};var t,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)t=o[r],n.indexOf(t)>=0||(a[t]=e[t]);return a}(e,n);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)t=o[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(a[t]=e[t])}return a}var l=r.createContext({}),p=function(e){var n=r.useContext(l),t=n;return e&&(t="function"==typeof e?e(n):i(i({},n),e)),t},c=function(e){var n=p(e.components);return r.createElement(l.Provider,{value:n},e.children)},u={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},d=r.forwardRef((function(e,n){var t=e.components,a=e.mdxType,o=e.originalType,l=e.parentName,c=s(e,["components","mdxType","originalType","parentName"]),d=p(t),m=a,f=d["".concat(l,".").concat(m)]||d[m]||u[m]||o;return t?r.createElement(f,i(i({ref:n},c),{},{components:t})):r.createElement(f,i({ref:n},c))}));function m(e,n){var t=arguments,a=n&&n.mdxType;if("string"==typeof e||a){var o=t.length,i=new Array(o);i[0]=d;var s={};for(var l in n)hasOwnProperty.call(n,l)&&(s[l]=n[l]);s.originalType=e,s.mdxType="string"==typeof e?e:a,i[1]=s;for(var p=2;p<o;p++)i[p]=t[p];return r.createElement.apply(null,i)}return r.createElement.apply(null,t)}d.displayName="MDXCreateElement"},9807:function(e,n,t){t.r(n),t.d(n,{frontMatter:function(){return s},contentTitle:function(){return l},metadata:function(){return p},toc:function(){return c},default:function(){return d}});var r=t(7462),a=t(3366),o=(t(7294),t(3905)),i=["components"],s={id:"organising-specifications",title:"Organising specifications"},l=void 0,p={unversionedId:"guides/organising-specifications",id:"guides/organising-specifications",isDocsHomePage:!1,title:"Organising specifications",description:"Optimus supports two ways to deploy specifications",source:"@site/docs/guides/organising-specifcations.md",sourceDirName:"guides",slug:"/guides/organising-specifications",permalink:"/optimus/docs/guides/organising-specifications",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/guides/organising-specifcations.md",tags:[],version:"current",lastUpdatedBy:"Sandeep Bhardwaj",lastUpdatedAt:1647329409,formattedLastUpdatedAt:"3/15/2022",frontMatter:{id:"organising-specifications",title:"Organising specifications"},sidebar:"docsSidebar",previous:{title:"Create bigquery external table",permalink:"/optimus/docs/guides/create-bigquery-external-table"},next:{title:"Starting Optimus Server",permalink:"/optimus/docs/guides/optimus-serve"}},c=[],u={toc:c};function d(e){var n=e.components,t=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},u,t,{components:n,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Optimus supports two ways to deploy specifications"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"REST/GRPC"),(0,o.kt)("li",{parentName:"ul"},"Optimus CLI deploy command")),(0,o.kt)("p",null,"When using Optimus cli to deploy, either manually or from a CI pipeline, it is advised to use version control system like git. Here is a simple directory structure that can be used as a template for jobs and datastore resources."),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre"},".\n\u251c\u2500\u2500 .optimus.yaml\n\u251c\u2500\u2500 README.md\n\u251c\u2500\u2500 datastore\n\u2502\xa0\xa0 \u251c\u2500\u2500 bigquery\n\u2502\xa0\xa0 \u2502\xa0\xa0 \u251c\u2500\u2500 project1\n\u2502\xa0\xa0 \u2502\xa0\xa0 \u2502\xa0\xa0 \u251c\u2500\u2500 dataset1\n\u2502\xa0\xa0 \u2502\xa0\xa0 \u2502\xa0\xa0 \u2502\xa0\xa0 \u251c\u2500\u2500 table1\n\u2502\xa0\xa0 \u2502\xa0\xa0 \u2502\xa0\xa0 \u2502\xa0\xa0 \u2514\u2500\u2500 table2\n\u2502\xa0\xa0 \u2502\xa0\xa0 \u2502\xa0\xa0 \u2514\u2500\u2500 this.yaml\n\u2502\xa0\xa0 \u2502\xa0\xa0 \u2514\u2500\u2500 project2\n\u2502\xa0\xa0 \u2502\xa0\xa0     \u2514\u2500\u2500 dataset1\n\u2502\xa0\xa0 \u2502\xa0\xa0         \u2514\u2500\u2500 table1\n\u2502\xa0\xa0 \u2514\u2500\u2500 postgres\n\u2502\xa0\xa0     \u2514\u2500\u2500 table1\n\u2514\u2500\u2500 jobs\n    \u251c\u2500\u2500 project1\n    \u2502\xa0\xa0 \u251c\u2500\u2500 job1\n    \u2502\xa0\xa0 \u251c\u2500\u2500 job2\n    \u2502\xa0\xa0 \u2514\u2500\u2500 this.yaml\n    \u251c\u2500\u2500 project2\n    \u2502\xa0\xa0 \u251c\u2500\u2500 job1\n    \u2502\xa0\xa0 \u2514\u2500\u2500 job2\n    \u2514\u2500\u2500 this.yaml\n")),(0,o.kt)("p",null,"A sample ",(0,o.kt)("inlineCode",{parentName:"p"},".optimus.yaml")," would look like"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\nhost: localhost:9100\njob:\n  path: jobs\ndatastore:\n- type: bigquery\n  path: datastore/bigquery\n- type: postgres\n  path: datastore/postgres\nconfig:\n  global:\n    environment: integration\n    storage_path: gs://example-bucket/test    \n  local: {}\n")),(0,o.kt)("p",null,"You might have also noticed there are ",(0,o.kt)("inlineCode",{parentName:"p"},"this.yaml")," files being used in some directories. This file is used to share a single set of configuration across multiple sub directories. For example if I create a file at ",(0,o.kt)("inlineCode",{parentName:"p"},"/jobs/project1/this.yaml"),", then all sub directories inside ",(0,o.kt)("inlineCode",{parentName:"p"},"/jobs/project1")," will inherit this config as defaults. If same config is specified in sub directory, then sub directory will override the parent defaults. For example a ",(0,o.kt)("inlineCode",{parentName:"p"},"this.yaml")," in ",(0,o.kt)("inlineCode",{parentName:"p"},"/jobs/project/")),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\nschedule:\n  interval: @daily\nbehavior:\n  depends_on_past: false\n  catch_up: true\n  retry:\n    count: 1\n    delay: 5s\nlabels:\n  owner: overlords\n  transform: sql\n")),(0,o.kt)("p",null,"and a ",(0,o.kt)("inlineCode",{parentName:"p"},"job.yaml")," in ",(0,o.kt)("inlineCode",{parentName:"p"},"/jobs/project/job1/")),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'name: sample_replace\nowner: optimus@example.io\nschedule:\n  start_date: "2020-09-25"\n  interval: 0 10 * * *\nbehavior:\n  depends_on_past: true\ntask:\n  name: bq2bq\n  config:\n    project: project_name\n    dataset: project_dataset\n    table: sample_replace\n    load_method: REPLACE\n  window:\n    size: 48h\n    offset: 24h\nlabels:\n  process: bq\n')),(0,o.kt)("p",null,"will result in final computed ",(0,o.kt)("inlineCode",{parentName:"p"},"job.yaml")," during deployment as"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'version: 1\nname: sample_replace\nowner: optimus@example.io\nschedule:\n  start_date: "2020-10-06"\n  interval: 0 10 * * *\nbehavior:\n  depends_on_past: true\n  catch_up: true\n  retry:\n    count: 1\n    delay: 5s\ntask:\n  name: bq2bq\n  config:\n    project: project_name\n    dataset: project_dataset\n    table: sample_replace\n    load_method: REPLACE\n  window:\n    size: 48h\n    offset: 24h\nlabels:\n  process: bq\n  owner: overlords\n  transform: sql\n')))}d.isMDXComponent=!0}}]);