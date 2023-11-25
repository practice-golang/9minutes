import{s as xt,n as pe,r as nt,e as $t,o as el,a as tl,f as ot}from"../chunks/scheduler.b95eede2.js";import{S as ll,i as nl,g as d,s as A,m as de,h as _,y as w,c as B,j as k,f as c,n as _e,k as O,A as R,a as E,z as u,B as H,o as qe,C as se,e as me,D as tt,E as ol,F as be,G as al}from"../chunks/index.e4e63e43.js";import{e as X,i as lt}from"../chunks/table.a906765d.js";import{p as rl}from"../chunks/stores.9b28410a.js";const il=async({url:t,fetch:l})=>{const e=Number(t.searchParams.get("list-count"))||20,n=Number(t.searchParams.get("page"))||1,a=t.searchParams.get("search")||"",o=[{"display-name":"Index","column-code":"idx","column-name":"IDX"},{"display-name":"Name","column-code":"board-name","column-name":"BOARD_NAME"},{"display-name":"Code","column-code":"board-code","column-name":"BOARD_CODE"},{"display-name":"Type","column-code":"board-type","column-name":"BOARD_TYPE"},{"display-name":"Board table","column-code":"board-table","column-name":"BOARD_TABLE"},{"display-name":"Comment table","column-code":"comment-table","column-name":"COMMENT_TABLE"},{"display-name":"Grant read","column-code":"grant-read","column-name":"GRANT_READ"},{"display-name":"Grant write","column-code":"grant-write","column-name":"GRANT_WRITE"},{"display-name":"Grant comment","column-code":"grant-comment","column-name":"GRANT_COMMENT"},{"display-name":"Grant upload","column-code":"grant-upload","column-name":"GRANT_UPLOAD"}];async function r(){let s=[];const C=await l("/api/admin/user-grades-for-grant",{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});if(C.ok){let h=Object.entries(await C.json()).sort((m,j)=>m[1].point-j[1].point);for(let m of h)s.push(m[1])}return s}async function f(s,C,h){let m={},j=`/api/admin/board?page=${s}&list-count=${C}`;h!=""&&(j+=`&search=${h}`);const D=await l(j,{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});return D.ok&&(m=await D.json()),m["board-list"]==null&&(m["board-list"]=[]),m}return{columns:o,grades:r(),"boardlist-data":f(n,e,a)}},Dl=Object.freeze(Object.defineProperty({__proto__:null,load:il},Symbol.toStringTag,{value:"Module"}));function St(t,l,e){const n=t.slice();return n[43]=l[e],n}function wt(t,l,e){const n=t.slice();return n[43]=l[e],n}function It(t,l,e){const n=t.slice();return n[48]=l[e],n[50]=e,n}function Ut(t,l,e){const n=t.slice();return n[51]=l[e],n}function jt(t,l,e){const n=t.slice();return n[51]=l[e],n[52]=l,n[53]=e,n}function Lt(t,l,e){const n=t.slice();return n[54]=l[e][0],n[55]=l[e][1],n}function Gt(t,l,e){const n=t.slice();return n[51]=l[e],n[60]=l,n[61]=e,n}function Rt(t,l,e){const n=t.slice();return n[54]=l[e][0],n[55]=l[e][1],n}function zt(t,l,e){const n=t.slice();return n[51]=l[e],n}function Mt(t){let l,e=t[51]["display-name"]+"",n;return{c(){l=d("th"),n=de(e)},l(a){l=_(a,"TH",{});var o=k(l);n=_e(o,e),o.forEach(c)},m(a,o){E(a,l,o),u(l,n)},p:pe,d(a){a&&c(l)}}}function qt(t){let l,e,n,a,o,r,f,s,C="Cancel",h,m,j="Save",D,T,I=X(t[11]),v=[];for(let P=0;P<I.length;P+=1)v[P]=Jt(Gt(t,I,P));return{c(){l=d("tr"),e=d("td"),n=A(),a=d("td"),o=A();for(let P=0;P<v.length;P+=1)v[P].c();r=A(),f=d("td"),s=d("button"),s.textContent=C,h=A(),m=d("button"),m.textContent=j,this.h()},l(P){l=_(P,"TR",{id:!0});var z=k(l);e=_(z,"TD",{}),k(e).forEach(c),n=B(z),a=_(z,"TD",{}),k(a).forEach(c),o=B(z);for(let b=0;b<v.length;b+=1)v[b].l(z);r=B(z),f=_(z,"TD",{});var g=k(f);s=_(g,"BUTTON",{type:!0,"data-svelte-h":!0}),w(s)!=="svelte-1gpezhr"&&(s.textContent=C),h=B(g),m=_(g,"BUTTON",{type:!0,"data-svelte-h":!0}),w(m)!=="svelte-1wprywb"&&(m.textContent=j),g.forEach(c),z.forEach(c),this.h()},h(){O(s,"type","button"),O(m,"type","button"),O(l,"id","add-board")},m(P,z){E(P,l,z),u(l,e),u(l,n),u(l,a),u(l,o);for(let g=0;g<v.length;g+=1)v[g]&&v[g].m(l,null);u(l,r),u(l,f),u(f,s),u(f,h),u(f,m),D||(T=[H(s,"click",t[17]),H(m,"click",t[18])],D=!0)},p(P,z){if(z[0]&38913){I=X(P[11]);let g;for(g=0;g<I.length;g+=1){const b=Gt(P,I,g);v[g]?v[g].p(b,z):(v[g]=Jt(b),v[g].c(),v[g].m(l,r))}for(;g<v.length;g+=1)v[g].d(1);v.length=I.length}},d(P){P&&c(l),se(v,P),D=!1,nt(T)}}}function ul(t){let l,e,n,a;function o(){t[31].call(e,t[51])}return{c(){l=d("td"),e=d("input"),this.h()},l(r){l=_(r,"TD",{});var f=k(l);e=_(f,"INPUT",{type:!0,placeholder:!0}),f.forEach(c),this.h()},h(){O(e,"type","text"),O(e,"placeholder",t[51]["display-name"])},m(r,f){E(r,l,f),u(l,e),R(e,t[0][t[51]["column-code"]]),n||(a=H(e,"input",o),n=!0)},p(r,f){t=r,f[0]&6145&&e.value!==t[0][t[51]["column-code"]]&&R(e,t[0][t[51]["column-code"]])},d(r){r&&c(l),n=!1,a()}}}function cl(t){let l;return{c(){l=d("td")},l(e){l=_(e,"TD",{}),k(l).forEach(c)},m(e,n){E(e,l,n)},p:pe,d(e){e&&c(l)}}}function sl(t){let l,e,n,a,o=X(Object.entries(t[12])),r=[];for(let s=0;s<o.length;s+=1)r[s]=Ht(Rt(t,o,s));function f(){t[30].call(e,t[51])}return{c(){l=d("td"),e=d("select");for(let s=0;s<r.length;s+=1)r[s].c();this.h()},l(s){l=_(s,"TD",{});var C=k(l);e=_(C,"SELECT",{});var h=k(e);for(let m=0;m<r.length;m+=1)r[m].l(h);h.forEach(c),C.forEach(c),this.h()},h(){t[0][t[51]["column-code"]]===void 0&&ot(f)},m(s,C){E(s,l,C),u(l,e);for(let h=0;h<r.length;h+=1)r[h]&&r[h].m(e,null);be(e,t[0][t[51]["column-code"]],!0),n||(a=H(e,"change",f),n=!0)},p(s,C){if(t=s,C[0]&4096){o=X(Object.entries(t[12]));let h;for(h=0;h<o.length;h+=1){const m=Rt(t,o,h);r[h]?r[h].p(m,C):(r[h]=Ht(m),r[h].c(),r[h].m(e,null))}for(;h<r.length;h+=1)r[h].d(1);r.length=o.length}C[0]&6145&&be(e,t[0][t[51]["column-code"]])},d(s){s&&c(l),se(r,s),n=!1,a()}}}function dl(t){let l,e,n,a;function o(){t[29].call(e,t[51])}return{c(){l=d("td"),e=d("input"),this.h()},l(r){l=_(r,"TD",{});var f=k(l);e=_(f,"INPUT",{}),f.forEach(c),this.h()},h(){e.disabled=!0},m(r,f){E(r,l,f),u(l,e),R(e,t[0][t[51]["column-code"]]),n||(a=H(e,"input",o),n=!0)},p(r,f){t=r,f[0]&6145&&e.value!==t[0][t[51]["column-code"]]&&R(e,t[0][t[51]["column-code"]])},d(r){r&&c(l),n=!1,a()}}}function _l(t){let l,e,n,a="Board",o,r="Gallery",f,s;function C(){t[28].call(e,t[51])}return{c(){l=d("td"),e=d("select"),n=d("option"),n.textContent=a,o=d("option"),o.textContent=r,this.h()},l(h){l=_(h,"TD",{});var m=k(l);e=_(m,"SELECT",{});var j=k(e);n=_(j,"OPTION",{"data-svelte-h":!0}),w(n)!=="svelte-ueeluy"&&(n.textContent=a),o=_(j,"OPTION",{"data-svelte-h":!0}),w(o)!=="svelte-79feqe"&&(o.textContent=r),j.forEach(c),m.forEach(c),this.h()},h(){n.__value="board",R(n,n.__value),o.__value="gallery",R(o,o.__value),t[0][t[51]["column-code"]]===void 0&&ot(C)},m(h,m){E(h,l,m),u(l,e),u(e,n),u(e,o),be(e,t[0][t[51]["column-code"]],!0),f||(s=H(e,"change",C),f=!0)},p(h,m){t=h,m[0]&6145&&be(e,t[0][t[51]["column-code"]])},d(h){h&&c(l),f=!1,s()}}}function fl(t){let l="",e;return{c(){e=de(l)},l(n){e=_e(n,l)},m(n,a){E(n,e,a)},p:pe,d(n){n&&c(e)}}}function Ht(t){let l,e=t[55].name+"",n;return{c(){l=d("option"),n=de(e),this.h()},l(a){l=_(a,"OPTION",{});var o=k(l);n=_e(o,e),o.forEach(c),this.h()},h(){l.__value=t[55].code,R(l,l.__value)},m(a,o){E(a,l,o),u(l,n)},p:pe,d(a){a&&c(l)}}}function Jt(t){let l;function e(o,r){return o[51]["column-code"]=="idx"?fl:o[51]["column-code"]=="board-type"?_l:o[51]["column-code"]=="board-table"||o[51]["column-code"]=="comment-table"?dl:o[15].includes(o[51]["column-code"])?sl:o[51]["column-code"]=="regdate"?cl:ul}let a=e(t)(t);return{c(){a.c(),l=me()},l(o){a.l(o),l=me()},m(o,r){a.m(o,r),E(o,l,r)},p(o,r){a.p(o,r)},d(o){o&&c(l),a.d(o)}}}function hl(t){let l,e,n,a=!1,o,r,f,s,C="View",h,m,j="Edit",D,T,I="Delete",v,P,z,g=X(t[11]),b=[];for(let U=0;U<g.length;U+=1)b[U]=Vt(Ut(t,g,U));function ve(){return t[38](t[50])}function x(){return t[39](t[50])}function ge(){return t[40](t[50])}return v=al(t[37][0]),{c(){l=d("td"),e=d("input"),o=A();for(let U=0;U<b.length;U+=1)b[U].c();r=A(),f=d("td"),s=d("button"),s.textContent=C,h=A(),m=d("button"),m.textContent=j,D=A(),T=d("button"),T.textContent=I,this.h()},l(U){l=_(U,"TD",{});var L=k(l);e=_(L,"INPUT",{type:!0}),L.forEach(c),o=B(U);for(let q=0;q<b.length;q+=1)b[q].l(U);r=B(U),f=_(U,"TD",{});var S=k(f);s=_(S,"BUTTON",{type:!0,"data-svelte-h":!0}),w(s)!=="svelte-4gqy8p"&&(s.textContent=C),h=B(S),m=_(S,"BUTTON",{type:!0,"data-svelte-h":!0}),w(m)!=="svelte-1t2ew4d"&&(m.textContent=j),D=B(S),T=_(S,"BUTTON",{type:!0,"data-svelte-h":!0}),w(T)!=="svelte-kppxo3"&&(T.textContent=I),S.forEach(c),this.h()},h(){O(e,"type","checkbox"),e.__value=n=t[48].idx,R(e,e.__value),O(s,"type","button"),O(m,"type","button"),O(T,"type","button"),v.p(e)},m(U,L){E(U,l,L),u(l,e),e.checked=~(t[5]||[]).indexOf(e.__value),E(U,o,L);for(let S=0;S<b.length;S+=1)b[S]&&b[S].m(U,L);E(U,r,L),E(U,f,L),u(f,s),u(f,h),u(f,m),u(f,D),u(f,T),P||(z=[H(e,"change",t[36]),H(s,"click",ve),H(m,"click",x),H(T,"click",ge)],P=!0)},p(U,L){if(t=U,L[0]&256&&n!==(n=t[48].idx)&&(e.__value=n,R(e,e.__value),a=!0),(a||L[0]&288)&&(e.checked=~(t[5]||[]).indexOf(e.__value)),L[0]&2304){g=X(t[11]);let S;for(S=0;S<g.length;S+=1){const q=Ut(t,g,S);b[S]?b[S].p(q,L):(b[S]=Vt(q),b[S].c(),b[S].m(r.parentNode,r))}for(;S<b.length;S+=1)b[S].d(1);b.length=g.length}},d(U){U&&(c(l),c(o),c(r),c(f)),se(b,U),v.r(),P=!1,nt(z)}}}function pl(t){let l,e,n,a,o,r="Cancel",f,s,C="Save",h,m,j=X(t[11]),D=[];for(let T=0;T<j.length;T+=1)D[T]=Yt(jt(t,j,T));return{c(){l=d("td"),e=A();for(let T=0;T<D.length;T+=1)D[T].c();n=A(),a=d("td"),o=d("button"),o.textContent=r,f=A(),s=d("button"),s.textContent=C,this.h()},l(T){l=_(T,"TD",{}),k(l).forEach(c),e=B(T);for(let v=0;v<D.length;v+=1)D[v].l(T);n=B(T),a=_(T,"TD",{});var I=k(a);o=_(I,"BUTTON",{type:!0,"data-svelte-h":!0}),w(o)!=="svelte-ufbtnv"&&(o.textContent=r),f=B(I),s=_(I,"BUTTON",{type:!0,"data-svelte-h":!0}),w(s)!=="svelte-x4yzzr"&&(s.textContent=C),I.forEach(c),this.h()},h(){O(o,"type","button"),O(s,"type","button")},m(T,I){E(T,l,I),E(T,e,I);for(let v=0;v<D.length;v+=1)D[v]&&D[v].m(T,I);E(T,n,I),E(T,a,I),u(a,o),u(a,f),u(a,s),h||(m=[H(o,"click",t[21]),H(s,"click",t[22])],h=!0)},p(T,I){if(I[0]&38914){j=X(T[11]);let v;for(v=0;v<j.length;v+=1){const P=jt(T,j,v);D[v]?D[v].p(P,I):(D[v]=Yt(P),D[v].c(),D[v].m(n.parentNode,n))}for(;v<D.length;v+=1)D[v].d(1);D.length=j.length}},d(T){T&&(c(l),c(e),c(n),c(a)),se(D,T),h=!1,nt(m)}}}function Vt(t){let l,e=t[48][t[51]["column-code"]]+"",n;return{c(){l=d("td"),n=de(e)},l(a){l=_(a,"TD",{});var o=k(l);n=_e(o,e),o.forEach(c)},m(a,o){E(a,l,o),u(l,n)},p(a,o){o[0]&256&&e!==(e=a[48][a[51]["column-code"]]+"")&&qe(n,e)},d(a){a&&c(l)}}}function ml(t){let l,e,n,a;function o(){t[35].call(e,t[51])}return{c(){l=d("td"),e=d("input"),this.h()},l(r){l=_(r,"TD",{});var f=k(l);e=_(f,"INPUT",{type:!0,placeholder:!0}),f.forEach(c),this.h()},h(){O(e,"type","text"),O(e,"placeholder",t[51]["display-name"])},m(r,f){E(r,l,f),u(l,e),R(e,t[1][t[51]["column-code"]]),n||(a=H(e,"input",o),n=!0)},p(r,f){t=r,f[0]&6146&&e.value!==t[1][t[51]["column-code"]]&&R(e,t[1][t[51]["column-code"]])},d(r){r&&c(l),n=!1,a()}}}function bl(t){let l;return{c(){l=d("td")},l(e){l=_(e,"TD",{}),k(l).forEach(c)},m(e,n){E(e,l,n)},p:pe,d(e){e&&c(l)}}}function vl(t){let l,e,n,a,o=X(Object.entries(t[12])),r=[];for(let s=0;s<o.length;s+=1)r[s]=Xt(Lt(t,o,s));function f(){t[34].call(e,t[51])}return{c(){l=d("td"),e=d("select");for(let s=0;s<r.length;s+=1)r[s].c();this.h()},l(s){l=_(s,"TD",{});var C=k(l);e=_(C,"SELECT",{});var h=k(e);for(let m=0;m<r.length;m+=1)r[m].l(h);h.forEach(c),C.forEach(c),this.h()},h(){t[1][t[51]["column-code"]]===void 0&&ot(f)},m(s,C){E(s,l,C),u(l,e);for(let h=0;h<r.length;h+=1)r[h]&&r[h].m(e,null);be(e,t[1][t[51]["column-code"]],!0),n||(a=H(e,"change",f),n=!0)},p(s,C){if(t=s,C[0]&4096){o=X(Object.entries(t[12]));let h;for(h=0;h<o.length;h+=1){const m=Lt(t,o,h);r[h]?r[h].p(m,C):(r[h]=Xt(m),r[h].c(),r[h].m(e,null))}for(;h<r.length;h+=1)r[h].d(1);r.length=o.length}C[0]&6146&&be(e,t[1][t[51]["column-code"]])},d(s){s&&c(l),se(r,s),n=!1,a()}}}function gl(t){let l,e,n,a;function o(){t[33].call(e,t[51])}return{c(){l=d("td"),e=d("input"),this.h()},l(r){l=_(r,"TD",{});var f=k(l);e=_(f,"INPUT",{}),f.forEach(c),this.h()},h(){e.disabled=!0},m(r,f){E(r,l,f),u(l,e),R(e,t[1][t[51]["column-code"]]),n||(a=H(e,"input",o),n=!0)},p(r,f){t=r,f[0]&6146&&e.value!==t[1][t[51]["column-code"]]&&R(e,t[1][t[51]["column-code"]])},d(r){r&&c(l),n=!1,a()}}}function Cl(t){let l,e,n,a="Board",o,r="Gallery",f,s;function C(){t[32].call(e,t[51])}return{c(){l=d("td"),e=d("select"),n=d("option"),n.textContent=a,o=d("option"),o.textContent=r,this.h()},l(h){l=_(h,"TD",{});var m=k(l);e=_(m,"SELECT",{});var j=k(e);n=_(j,"OPTION",{"data-svelte-h":!0}),w(n)!=="svelte-ueeluy"&&(n.textContent=a),o=_(j,"OPTION",{"data-svelte-h":!0}),w(o)!=="svelte-79feqe"&&(o.textContent=r),j.forEach(c),m.forEach(c),this.h()},h(){n.__value="board",R(n,n.__value),o.__value="gallery",R(o,o.__value),t[1][t[51]["column-code"]]===void 0&&ot(C)},m(h,m){E(h,l,m),u(l,e),u(e,n),u(e,o),be(e,t[1][t[51]["column-code"]],!0),f||(s=H(e,"change",C),f=!0)},p(h,m){t=h,m[0]&6146&&be(e,t[1][t[51]["column-code"]])},d(h){h&&c(l),f=!1,s()}}}function Tl(t){let l,e=t[1].idx+"",n;return{c(){l=d("td"),n=de(e)},l(a){l=_(a,"TD",{});var o=k(l);n=_e(o,e),o.forEach(c)},m(a,o){E(a,l,o),u(l,n)},p(a,o){o[0]&2&&e!==(e=a[1].idx+"")&&qe(n,e)},d(a){a&&c(l)}}}function Xt(t){let l,e=t[55].name+"",n;return{c(){l=d("option"),n=de(e),this.h()},l(a){l=_(a,"OPTION",{});var o=k(l);n=_e(o,e),o.forEach(c),this.h()},h(){l.__value=t[55].point,R(l,l.__value)},m(a,o){E(a,l,o),u(l,n)},p:pe,d(a){a&&c(l)}}}function Yt(t){let l;function e(o,r){return o[51]["column-code"]=="idx"?Tl:o[51]["column-code"]=="board-type"?Cl:o[51]["column-code"]=="board-table"||o[51]["column-code"]=="comment-table"?gl:o[15].includes(o[51]["column-code"])?vl:o[51]["column-code"]=="regdate"?bl:ml}let a=e(t)(t);return{c(){a.c(),l=me()},l(o){a.l(o),l=me()},m(o,r){a.m(o,r),E(o,l,r)},p(o,r){a.p(o,r)},d(o){o&&c(l),a.d(o)}}}function Ft(t){let l,e;function n(r,f){return r[6]==r[50]?pl:hl}let a=n(t),o=a(t);return{c(){l=d("tr"),o.c(),e=A()},l(r){l=_(r,"TR",{});var f=k(l);o.l(f),e=B(f),f.forEach(c)},m(r,f){E(r,l,f),o.m(l,null),u(l,e)},p(r,f){a===(a=n(r))&&o?o.p(r,f):(o.d(1),o=a(r),o&&(o.c(),o.m(l,e)))},d(r){r&&c(l),o.d()}}}function Kt(t){let l,e=t[43]+"",n,a;return{c(){l=d("a"),n=de(e),this.h()},l(o){l=_(o,"A",{href:!0});var r=k(l);n=_e(r,e),r.forEach(c),this.h()},h(){O(l,"href",a=`?page=${t[43]}&list-count=${t[13]}`+(t[14]!=""?`&search=${t[14]}`:""))},m(o,r){E(o,l,r),u(l,n)},p(o,r){r[0]&4&&e!==(e=o[43]+"")&&qe(n,e),r[0]&4&&a!==(a=`?page=${o[43]}&list-count=${o[13]}`+(o[14]!=""?`&search=${o[14]}`:""))&&O(l,"href",a)},d(o){o&&c(l)}}}function Wt(t){let l,e=t[43]>=1&&Kt(t);return{c(){e&&e.c(),l=me()},l(n){e&&e.l(n),l=me()},m(n,a){e&&e.m(n,a),E(n,l,a)},p(n,a){n[43]>=1?e?e.p(n,a):(e=Kt(n),e.c(),e.m(l.parentNode,l)):e&&(e.d(1),e=null)},d(n){n&&c(l),e&&e.d(n)}}}function Qt(t){let l,e=t[43]+"",n,a;return{c(){l=d("a"),n=de(e),this.h()},l(o){l=_(o,"A",{href:!0});var r=k(l);n=_e(r,e),r.forEach(c),this.h()},h(){O(l,"href",a=`?page=${t[43]}&list-count=${t[13]}`+(t[14]!=""?`&search=${t[14]}`:""))},m(o,r){E(o,l,r),u(l,n)},p(o,r){r[0]&4&&e!==(e=o[43]+"")&&qe(n,e),r[0]&4&&a!==(a=`?page=${o[43]}&list-count=${o[13]}`+(o[14]!=""?`&search=${o[14]}`:""))&&O(l,"href",a)},d(o){o&&c(l)}}}function Zt(t){let l,e=t[43]<=t[3]&&Qt(t);return{c(){e&&e.c(),l=me()},l(n){e&&e.l(n),l=me()},m(n,a){e&&e.m(n,a),E(n,l,a)},p(n,a){n[43]<=n[3]?e?e.p(n,a):(e=Qt(n),e.c(),e.m(l.parentNode,l)):e&&(e.d(1),e=null)},d(n){n&&c(l),e&&e.d(n)}}}function yl(t){let l,e="Boards",n,a,o,r="Add board",f,s,C="|",h,m,j="Delete selected boards",D,T,I="|",v,P,z="Search:",g,b,ve,x,ge="Search",U,L,S,q,Ce,le,Re,Te,fe,He="Control",ze,Oe,ie,Pe,Q,ne,Je="Admin",oe,Ve="Manager",ae,Xe="Regular user",re,y="Pending user",M,$="Banned user",J,Ye="Guest",Fe,he,ue,pt="Board",ce,mt="Gallery",Ke,G,De,Ae,bt="«",at,ye,Be,vt="<",We,rt,Se,gt="..",it,Qe,Me,Ze,ut,xe,we,Ct="..",ct,ke,Ie,Tt=">",$e,st,Ee,Ue,yt="»",et,dt,kt,je=X(t[11]),Y=[];for(let i=0;i<je.length;i+=1)Y[i]=Mt(zt(t,je,i));let K=t[7]&&qt(t),Le=X(t[8]),F=[];for(let i=0;i<Le.length;i+=1)F[i]=Ft(It(t,Le,i));let _t=X([t[2]-2,t[2]-1]),ee=[];for(let i=0;i<2;i+=1)ee[i]=Wt(wt(t,_t,i));let ft=X([t[2]+1,t[2]+2]),te=[];for(let i=0;i<2;i+=1)te[i]=Zt(St(t,ft,i));return{c(){l=d("h1"),l.textContent=e,n=A(),a=d("div"),o=d("button"),o.textContent=r,f=A(),s=d("span"),s.textContent=C,h=A(),m=d("button"),m.textContent=j,D=A(),T=d("span"),T.textContent=I,v=A(),P=d("label"),P.textContent=z,g=A(),b=d("input"),ve=A(),x=d("button"),x.textContent=ge,U=A(),L=d("table"),S=d("thead"),q=d("tr"),Ce=d("th"),le=d("input"),Re=A();for(let i=0;i<Y.length;i+=1)Y[i].c();Te=A(),fe=d("th"),fe.textContent=He,ze=A(),K&&K.c(),Oe=A(),ie=d("tbody");for(let i=0;i<F.length;i+=1)F[i].c();Pe=A(),Q=d("datalist"),ne=d("option"),ne.textContent=Je,oe=d("option"),oe.textContent=Ve,ae=d("option"),ae.textContent=Xe,re=d("option"),re.textContent=y,M=d("option"),M.textContent=$,J=d("option"),J.textContent=Ye,Fe=A(),he=d("datalist"),ue=d("option"),ue.textContent=pt,ce=d("option"),ce.textContent=mt,Ke=A(),G=d("div"),De=d("a"),Ae=d("span"),Ae.textContent=bt,at=A(),ye=d("a"),Be=d("span"),Be.textContent=vt,rt=A(),Se=d("span"),Se.textContent=gt,it=A();for(let i=0;i<2;i+=1)ee[i].c();Qe=A(),Me=d("b"),Ze=de(t[2]),ut=A();for(let i=0;i<2;i+=1)te[i].c();xe=A(),we=d("span"),we.textContent=Ct,ct=A(),ke=d("a"),Ie=d("span"),Ie.textContent=Tt,st=A(),Ee=d("a"),Ue=d("span"),Ue.textContent=yt,this.h()},l(i){l=_(i,"H1",{"data-svelte-h":!0}),w(l)!=="svelte-1um8syh"&&(l.textContent=e),n=B(i),a=_(i,"DIV",{});var N=k(a);o=_(N,"BUTTON",{type:!0,"data-svelte-h":!0}),w(o)!=="svelte-16bu93y"&&(o.textContent=r),f=B(N),s=_(N,"SPAN",{"data-svelte-h":!0}),w(s)!=="svelte-1e2i4m"&&(s.textContent=C),h=B(N),m=_(N,"BUTTON",{type:!0,"data-svelte-h":!0}),w(m)!=="svelte-1ookiet"&&(m.textContent=j),D=B(N),T=_(N,"SPAN",{"data-svelte-h":!0}),w(T)!=="svelte-1e2i4m"&&(T.textContent=I),v=B(N),P=_(N,"LABEL",{for:!0,"data-svelte-h":!0}),w(P)!=="svelte-1tn9osg"&&(P.textContent=z),g=B(N),b=_(N,"INPUT",{type:!0,id:!0,onkeyup:!0,placeholder:!0}),ve=B(N),x=_(N,"BUTTON",{type:!0,onclick:!0,"data-svelte-h":!0}),w(x)!=="svelte-155jhv4"&&(x.textContent=ge),N.forEach(c),U=B(i),L=_(i,"TABLE",{id:!0});var p=k(L);S=_(p,"THEAD",{});var Z=k(S);q=_(Z,"TR",{});var Ge=k(q);Ce=_(Ge,"TH",{});var Et=k(Ce);le=_(Et,"INPUT",{type:!0}),Et.forEach(c),Re=B(Ge);for(let W=0;W<Y.length;W+=1)Y[W].l(Ge);Te=B(Ge),fe=_(Ge,"TH",{"data-svelte-h":!0}),w(fe)!=="svelte-10rrrw5"&&(fe.textContent=He),Ge.forEach(c),Z.forEach(c),ze=B(p),K&&K.l(p),Oe=B(p),ie=_(p,"TBODY",{id:!0});var Nt=k(ie);for(let W=0;W<F.length;W+=1)F[W].l(Nt);Nt.forEach(c),p.forEach(c),Pe=B(i),Q=_(i,"DATALIST",{id:!0});var Ne=k(Q);ne=_(Ne,"OPTION",{"data-svelte-h":!0}),w(ne)!=="svelte-17op260"&&(ne.textContent=Je),oe=_(Ne,"OPTION",{"data-svelte-h":!0}),w(oe)!=="svelte-1ob7gsw"&&(oe.textContent=Ve),ae=_(Ne,"OPTION",{"data-svelte-h":!0}),w(ae)!=="svelte-ielivd"&&(ae.textContent=Xe),re=_(Ne,"OPTION",{"data-svelte-h":!0}),w(re)!=="svelte-311eip"&&(re.textContent=y),M=_(Ne,"OPTION",{"data-svelte-h":!0}),w(M)!=="svelte-1ot027p"&&(M.textContent=$),J=_(Ne,"OPTION",{"data-svelte-h":!0}),w(J)!=="svelte-1j47qc6"&&(J.textContent=Ye),Ne.forEach(c),Fe=B(i),he=_(i,"DATALIST",{id:!0});var ht=k(he);ue=_(ht,"OPTION",{"data-svelte-h":!0}),w(ue)!=="svelte-ueeluy"&&(ue.textContent=pt),ce=_(ht,"OPTION",{"data-svelte-h":!0}),w(ce)!=="svelte-79feqe"&&(ce.textContent=mt),ht.forEach(c),Ke=B(i),G=_(i,"DIV",{id:!0});var V=k(G);De=_(V,"A",{href:!0});var Ot=k(De);Ae=_(Ot,"SPAN",{"data-svelte-h":!0}),w(Ae)!=="svelte-1z054it"&&(Ae.textContent=bt),Ot.forEach(c),at=B(V),ye=_(V,"A",{href:!0});var Pt=k(ye);Be=_(Pt,"SPAN",{"data-svelte-h":!0}),w(Be)!=="svelte-1kd6by1"&&(Be.textContent=vt),Pt.forEach(c),rt=B(V),Se=_(V,"SPAN",{"data-svelte-h":!0}),w(Se)!=="svelte-1v1zlza"&&(Se.textContent=gt),it=B(V);for(let W=0;W<2;W+=1)ee[W].l(V);Qe=B(V),Me=_(V,"B",{});var Dt=k(Me);Ze=_e(Dt,t[2]),Dt.forEach(c),ut=B(V);for(let W=0;W<2;W+=1)te[W].l(V);xe=B(V),we=_(V,"SPAN",{"data-svelte-h":!0}),w(we)!=="svelte-1v1zlza"&&(we.textContent=Ct),ct=B(V),ke=_(V,"A",{href:!0});var At=k(ke);Ie=_(At,"SPAN",{"data-svelte-h":!0}),w(Ie)!=="svelte-x0xyl0"&&(Ie.textContent=Tt),At.forEach(c),st=B(V),Ee=_(V,"A",{href:!0});var Bt=k(Ee);Ue=_(Bt,"SPAN",{"data-svelte-h":!0}),w(Ue)!=="svelte-131q397"&&(Ue.textContent=yt),Bt.forEach(c),V.forEach(c),this.h()},h(){O(o,"type","button"),O(m,"type","button"),O(P,"for","search"),O(b,"type","text"),O(b,"id","search"),O(b,"onkeyup","pressEnter()"),O(b,"placeholder","Search for..."),O(x,"type","button"),O(x,"onclick","search()"),O(le,"type","checkbox"),O(ie,"id","boards-list-body"),O(L,"id","boards-list-container"),ne.__value="admin",R(ne,ne.__value),oe.__value="manager",R(oe,oe.__value),ae.__value="regular_user",R(ae,ae.__value),re.__value="pending_user",R(re,re.__value),M.__value="banned_user",R(M,M.__value),J.__value="guest",R(J,J.__value),O(Q,"id","grant-list"),ue.__value="board",R(ue,ue.__value),ce.__value="gallery",R(ce,ce.__value),O(he,"id","board-types"),O(De,"href",`?page=1&list-count=${t[13]}`+(t[14]!=""?`&search=${t[14]}`:"")),O(ye,"href",We=`?page=${t[10]}&list-count=${t[13]}`+(t[14]!=""?`&search=${t[14]}`:"")),O(ke,"href",$e=`?page=${t[9]}&list-count=${t[13]}`+(t[14]!=""?`&search=${t[14]}`:"")),O(Ee,"href",et=`?page=${t[3]}&list-count=${t[13]}`+(t[14]!=""?`&search=${t[14]}`:"")),O(G,"id","page-container")},m(i,N){E(i,l,N),E(i,n,N),E(i,a,N),u(a,o),u(a,f),u(a,s),u(a,h),u(a,m),u(a,D),u(a,T),u(a,v),u(a,P),u(a,g),u(a,b),u(a,ve),u(a,x),E(i,U,N),E(i,L,N),u(L,S),u(S,q),u(q,Ce),u(Ce,le),le.checked=t[4],u(q,Re);for(let p=0;p<Y.length;p+=1)Y[p]&&Y[p].m(q,null);u(q,Te),u(q,fe),u(L,ze),K&&K.m(L,null),u(L,Oe),u(L,ie);for(let p=0;p<F.length;p+=1)F[p]&&F[p].m(ie,null);E(i,Pe,N),E(i,Q,N),u(Q,ne),u(Q,oe),u(Q,ae),u(Q,re),u(Q,M),u(Q,J),E(i,Fe,N),E(i,he,N),u(he,ue),u(he,ce),E(i,Ke,N),E(i,G,N),u(G,De),u(De,Ae),u(G,at),u(G,ye),u(ye,Be),u(G,rt),u(G,Se),u(G,it);for(let p=0;p<2;p+=1)ee[p]&&ee[p].m(G,null);u(G,Qe),u(G,Me),u(Me,Ze),u(G,ut);for(let p=0;p<2;p+=1)te[p]&&te[p].m(G,null);u(G,xe),u(G,we),u(G,ct),u(G,ke),u(ke,Ie),u(G,st),u(G,Ee),u(Ee,Ue),dt||(kt=[H(o,"click",t[26]),H(m,"click",t[24]),H(le,"change",t[27]),H(le,"change",t[16])],dt=!0)},p(i,N){if(N[0]&16&&(le.checked=i[4]),N[0]&2048){je=X(i[11]);let p;for(p=0;p<je.length;p+=1){const Z=zt(i,je,p);Y[p]?Y[p].p(Z,N):(Y[p]=Mt(Z),Y[p].c(),Y[p].m(q,Te))}for(;p<Y.length;p+=1)Y[p].d(1);Y.length=je.length}if(i[7]?K?K.p(i,N):(K=qt(i),K.c(),K.m(L,Oe)):K&&(K.d(1),K=null),N[0]&16292194){Le=X(i[8]);let p;for(p=0;p<Le.length;p+=1){const Z=It(i,Le,p);F[p]?F[p].p(Z,N):(F[p]=Ft(Z),F[p].c(),F[p].m(ie,null))}for(;p<F.length;p+=1)F[p].d(1);F.length=Le.length}if(N[0]&1024&&We!==(We=`?page=${i[10]}&list-count=${i[13]}`+(i[14]!=""?`&search=${i[14]}`:""))&&O(ye,"href",We),N[0]&24580){_t=X([i[2]-2,i[2]-1]);let p;for(p=0;p<2;p+=1){const Z=wt(i,_t,p);ee[p]?ee[p].p(Z,N):(ee[p]=Wt(Z),ee[p].c(),ee[p].m(G,Qe))}for(;p<2;p+=1)ee[p].d(1)}if(N[0]&4&&qe(Ze,i[2]),N[0]&24588){ft=X([i[2]+1,i[2]+2]);let p;for(p=0;p<2;p+=1){const Z=St(i,ft,p);te[p]?te[p].p(Z,N):(te[p]=Zt(Z),te[p].c(),te[p].m(G,xe))}for(;p<2;p+=1)te[p].d(1)}N[0]&512&&$e!==($e=`?page=${i[9]}&list-count=${i[13]}`+(i[14]!=""?`&search=${i[14]}`:""))&&O(ke,"href",$e),N[0]&8&&et!==(et=`?page=${i[3]}&list-count=${i[13]}`+(i[14]!=""?`&search=${i[14]}`:""))&&O(Ee,"href",et)},i:pe,o:pe,d(i){i&&(c(l),c(n),c(a),c(U),c(L),c(Pe),c(Q),c(Fe),c(he),c(Ke),c(G)),se(Y,i),K&&K.d(),se(F,i),se(ee,i),se(te,i),dt=!1,nt(kt)}}}function kl(t,l,e){let n,a,o,r,f,s;$t(t,rl,y=>e(42,s=y));let{data:C}=l;const h=C.columns,m=C.grades;let j=Number(s.url.searchParams.get("list-count"))||10,D=-1,T=s.url.searchParams.get("search")||"",I=!1,v=[],P=-1,z=!1,g={},b={};const ve=["grant-read","grant-write","grant-comment","grant-upload"];function x(){I?e(5,v=[]):e(5,v=n.map(y=>y.idx))}function ge(){e(0,g={}),e(7,z=!1)}async function U(){const M=await fetch("/api/admin/board",{method:"POST",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(g)});M.ok||alert(await M.text()),ge(),lt()}function L(y){window.open("/board/list?board_code="+n[y]["board-code"],"_blank")}function S(y){e(6,P=y),e(1,b={});for(const M in n[y])e(1,b[M]=n[y][M],b)}function q(){e(1,b={}),e(6,P=-1)}async function Ce(){for(const $ in b)typeof b[$]==null||b[$]==null||typeof b[$]!="string"&&e(1,b[$]=b[$].toString(),b);const M=await fetch("/api/admin/board",{method:"PUT",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([b])});M.ok||alert(await M.text()),q(),lt()}async function le(y){const M=typeof n[y].idx=="number"?n[y].idx.toString():n[y].idx,J=await fetch("/api/admin/board",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([{idx:M}])});J.ok||alert(await J.text()),e(5,v=[]),lt()}async function Re(){if(v.length==0){alert("Selected nothing");return}const y=[];for(let J=0;J<v.length;J++){const Ye=typeof v[J]=="number"?v[J].toString():v[J];y.push({idx:Ye})}const $=await fetch("/api/admin/board",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(y)});$.ok||alert(await $.text()),e(5,v=[]),lt()}el(()=>{}),tl(()=>{D!=a&&(e(5,v=[]),D=a),e(4,I=v.length==n.length)});const Te=[[]],fe=()=>{e(0,g={"board-type":"board"}),e(7,z=!0)};function He(){I=this.checked,e(4,I)}function ze(y){g[y["column-code"]]=tt(this),e(0,g),e(1,b),e(12,m)}function Oe(y){g[y["column-code"]]=this.value,e(0,g),e(1,b),e(12,m)}function ie(y){g[y["column-code"]]=tt(this),e(0,g),e(1,b),e(12,m)}function Pe(y){g[y["column-code"]]=this.value,e(0,g),e(1,b),e(12,m)}function Q(y){b[y["column-code"]]=tt(this),e(1,b),e(0,g),e(12,m)}function ne(y){b[y["column-code"]]=this.value,e(1,b),e(0,g),e(12,m)}function Je(y){b[y["column-code"]]=tt(this),e(1,b),e(0,g),e(12,m)}function oe(y){b[y["column-code"]]=this.value,e(1,b),e(0,g),e(12,m)}function Ve(){v=ol(Te[0],this.__value,this.checked),e(5,v)}const ae=y=>{L(y)},Xe=y=>{S(y)},re=y=>{le(y)};return t.$$set=y=>{"data"in y&&e(25,C=y.data)},t.$$.update=()=>{t.$$.dirty[0]&33554432&&e(8,n=C["boardlist-data"]["board-list"]),t.$$.dirty[0]&33554432&&e(2,a=C["boardlist-data"]["current-page"]),t.$$.dirty[0]&33554432&&e(3,o=C["boardlist-data"]["total-page"]),t.$$.dirty[0]&4&&e(10,r=a-5>1?a-5:1),t.$$.dirty[0]&12&&e(9,f=a+5<o?a+5:o),t.$$.dirty[0]&3&&(g["board-code"]!=null&&g["board-code"].length>0&&(e(0,g["board-table"]="board_"+g["board-code"].toLowerCase(),g),e(0,g["comment-table"]="comment_"+g["board-code"].toLowerCase(),g)),b["board-code"]!=null&&b["board-code"].length>0&&(e(1,b["board-table"]="board_"+b["board-code"].toLowerCase(),b),e(1,b["comment-table"]="comment_"+b["board-code"].toLowerCase(),b)))},[g,b,a,o,I,v,P,z,n,f,r,h,m,j,T,ve,x,ge,U,L,S,q,Ce,le,Re,C,fe,He,ze,Oe,ie,Pe,Q,ne,Je,oe,Ve,Te,ae,Xe,re]}class Al extends ll{constructor(l){super(),nl(this,l,kl,yl,xt,{data:25},null,[-1,-1,-1])}}export{Al as component,Dl as universal};
