"use strict";(self.webpackChunkoptimus=self.webpackChunkoptimus||[]).push([[5254],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return f}});var r=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,i=function(e,t){if(null==e)return{};var n,r,i={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var s=r.createContext({}),p=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},u=function(e){var t=p(e.components);return r.createElement(s.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},c=r.forwardRef((function(e,t){var n=e.components,i=e.mdxType,o=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),c=p(n),f=i,m=c["".concat(s,".").concat(f)]||c[f]||d[f]||o;return n?r.createElement(m,a(a({ref:t},u),{},{components:n})):r.createElement(m,a({ref:t},u))}));function f(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=n.length,a=new Array(o);a[0]=c;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:i,a[1]=l;for(var p=2;p<o;p++)a[p]=n[p];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}c.displayName="MDXCreateElement"},5500:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return l},contentTitle:function(){return s},metadata:function(){return p},toc:function(){return u},default:function(){return c}});var r=n(7462),i=n(3366),o=(n(7294),n(3905)),a=["components"],l={id:"profiling-auditing",title:"Profiling and Auditing Bigquery"},s="Profiling and Auditing BigQuery",p={unversionedId:"guides/profiling-auditing",id:"guides/profiling-auditing",isDocsHomePage:!1,title:"Profiling and Auditing Bigquery",description:"To enable Profiler and Auditor (Predator), answer the related questions in Job specification.",source:"@site/docs/guides/predator.md",sourceDirName:"guides",slug:"/guides/profiling-auditing",permalink:"/optimus/docs/guides/profiling-auditing",editUrl:"https://github.com/odpf/optimus/edit/master/docs/docs/guides/predator.md",tags:[],version:"current",lastUpdatedBy:"Anwar Hidayat",lastUpdatedAt:1647946380,formattedLastUpdatedAt:"3/22/2022",frontMatter:{id:"profiling-auditing",title:"Profiling and Auditing Bigquery"}},u=[],d={toc:u};function c(e){var t=e.components,n=(0,i.Z)(e,a);return(0,o.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"profiling-and-auditing-bigquery"},"Profiling and Auditing BigQuery"),(0,o.kt)("p",null,"To enable Profiler and Auditor (Predator), answer the related questions in Job specification.\nNote: this is not available for public use at the moment"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"\u2022 $ optimus job addhook\n? What is the job name? test_job\n? Who is the Owner of this job? de@go-jek.com\n? Specify the start date (YYYY-MM-DD) 2021-01-01\n? Specify the interval (in crontab notation) @daily\n? Enable profile for the destination table? true\n? Enable audit for the destination table? true\n? Filter expression for profiling? (empty for always do full scan profiling) \nevent_timestamp >= '{{.DSTART}}' AND event_timestamp < '{{.DEND}}'\n? Specify the profile/audit result grouping field (empty to not group the result) __PARTITION__\n? Choose the profiling mode complete\n")),(0,o.kt)("p",null,"Configs:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Filter expression"),": Expression is used as a where clause to restrict the number of rows to only profile the ones\nthat needed to be profiled.\nExpression can be templated with: DSTART and DEND. These will be replaced with the window for which the current\ntransformation is getting executed. EXECUTION_TIME will be replaced with job execution time that is being\nused by the transformation task. ",(0,o.kt)("inlineCode",{parentName:"li"},"__PARTITION__")," represents the partitioning field of the table and the type of\npartition. If it is a daily partition using field ",(0,o.kt)("inlineCode",{parentName:"li"},"event_timestamp"),", then the macros is equal to date\n",(0,o.kt)("inlineCode",{parentName:"li"},"event_timestamp"),"."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Group"),": Represent the column on which the records will be grouped for profiling. Can be ",(0,o.kt)("inlineCode",{parentName:"li"},"__PARTITION__")," or any other\nfield in the target table."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"Mode"),": Mode represents the profiling strategy used with the above configurations, it doesn\u2019t affect the profile\nresults. ",(0,o.kt)("inlineCode",{parentName:"li"},"complete")," means all the records in a given group are considered for profiling. \u2018incremental\u2019 only the newly added records for the given group are considered for profiling. This input is needed when DataQuality results are shown in UI.")),(0,o.kt)("p",null,"Here is a sample DAG specification that has Predator enabled."),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"version: 1\nname: test_job\nowner: de@go-jek.com\nschedule:\n  start_date: \"2021-02-26\"\n  interval: 0 2 * * *\nbehavior:\n  depends_on_past: false\n  catch_up: true\ntask:\n  name: bq2bq\n  config:\n    DATASET: playground\n    LOAD_METHOD: REPLACE\n    PROJECT: gcp-project\n    SQL_TYPE: STANDARD\n    TABLE: hello_test_table\n  window:\n    size: 24h\n    offset: \"0\"\n    truncate_to: d\ndependencies: []\nhooks:\n  - name: predator\n    config:\n      AUDIT_TIME: '{{.EXECUTION_TIME}}'\n      BQ_DATASET: '{{.TASK__DATASET}}'\n      BQ_PROJECT: '{{.TASK__PROJECT}}'\n      BQ_TABLE: '{{.TASK__TABLE}}'\n      FILTER: 'event_timestamp >= \"{{.DSTART}}\" AND event_timestamp < \"{{.DEND}}\"'\n      GROUP: __PARTITION__\n      MODE: complete\n      PREDATOR_URL: '{{.GLOBAL__PREDATOR_HOST}}'\n      SUB_COMMAND: profile_audit\n")),(0,o.kt)("p",null,"After the Job is created, create a Data Quality Spec of the particular table and\nplace it in the Optimus jobs repository, inside the Predator directory.\nDetail of quality spec creation is available in Predator documentation."))}c.isMDXComponent=!0}}]);