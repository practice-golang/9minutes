import{s as _l,f as We,n as he,r as _t,e as dl,o as hl,a as pl}from"../chunks/scheduler.b95eede2.js";import{S as ml,i as bl,g as s,s as B,m as fe,h as f,y as L,c as w,j as N,f as u,n as _e,k as A,A as G,a as m,B as pe,z as h,C as q,o as Qe,D as se,e as oe,E as Ke,F as vl,G as gl}from"../chunks/index.ac44e828.js";import{e as M,i as ft}from"../chunks/table.24a228dd.js";import{p as Cl}from"../chunks/stores.20ae8124.js";const qt=10;async function kl(){let l=[];const t=await fetch("/api/admin/board-grades",{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});if(t.ok){let e=Object.entries(await t.json()).sort((n,a)=>n[1].rank-a[1].rank);for(let n of e)l.push(n[1])}return l}async function Tl(l,t,e){let n={},a=`/api/admin/board?page=${l}&list-count=${t}`;e!=""&&(a+=`&search=${e}`);const o=await fetch(a,{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});return o.ok&&(n=await o.json()),n["board-list"]==null&&(n["board-list"]=[]),n}const yl=async({url:l,fetch:t})=>{const e=Number(l.searchParams.get("list-count"))||qt,n=Number(l.searchParams.get("page"))||1,a=l.searchParams.get("search")||"";return{"default-count":qt,columns:[{"display-name":"Index","column-code":"idx","column-name":"IDX"},{"display-name":"Name","column-code":"board-name","column-name":"BOARD_NAME"},{"display-name":"Code","column-code":"board-code","column-name":"BOARD_CODE"},{"display-name":"Type","column-code":"board-type","column-name":"BOARD_TYPE"},{"display-name":"Board table","column-code":"board-table","column-name":"BOARD_TABLE"},{"display-name":"Comment table","column-code":"comment-table","column-name":"COMMENT_TABLE"},{"display-name":"Grant read","column-code":"grant-read","column-name":"GRANT_READ"},{"display-name":"Grant write","column-code":"grant-write","column-name":"GRANT_WRITE"},{"display-name":"Grant comment","column-code":"grant-comment","column-name":"GRANT_COMMENT"},{"display-name":"Grant upload","column-code":"grant-upload","column-name":"GRANT_UPLOAD"}],grades:kl(),"boardlist-data":Tl(n,e,a)}},Wl=Object.freeze(Object.defineProperty({__proto__:null,load:yl},Symbol.toStringTag,{value:"Module"}));function Mt(l,t,e){const n=l.slice();return n[49]=t[e],n}function Ht(l,t,e){const n=l.slice();return n[49]=t[e],n}function Jt(l,t,e){const n=l.slice();return n[54]=t[e],n[56]=e,n}function Xt(l,t,e){const n=l.slice();return n[57]=t[e],n}function Vt(l,t,e){const n=l.slice();return n[60]=t[e][0],n[61]=t[e][1],n}function Yt(l,t,e){const n=l.slice();return n[57]=t[e],n[58]=t,n[59]=e,n}function Ft(l,t,e){const n=l.slice();return n[60]=t[e][0],n[61]=t[e][1],n}function Kt(l,t,e){const n=l.slice();return n[57]=t[e],n[68]=t,n[69]=e,n}function Wt(l,t,e){const n=l.slice();return n[60]=t[e][0],n[61]=t[e][1],n}function Qt(l,t,e){const n=l.slice();return n[57]=t[e],n}function El(l,t,e){const n=l.slice();return n[74]=t[e],n}function Nl(l){let t,e;return{c(){t=s("option"),e=fe(l[74]),this.h()},l(n){t=f(n,"OPTION",{});var a=N(t);e=_e(a,l[74]),a.forEach(u),this.h()},h(){t.__value=l[74],G(t,t.__value)},m(n,a){m(n,t,a),h(t,e)},p:he,d(n){n&&u(t)}}}function Zt(l){let t,e=l[57]["display-name"]+"",n;return{c(){t=s("th"),n=fe(e)},l(a){t=f(a,"TH",{});var o=N(t);n=_e(o,e),o.forEach(u)},m(a,o){m(a,t,o),h(t,n)},p:he,d(a){a&&u(t)}}}function xt(l){let t,e,n,a,o,r,c,p,T="Cancel",d,C,I="Save",D,O,j=M(l[13]),b=[];for(let P=0;P<j.length;P+=1)b[P]=el(Kt(l,j,P));return{c(){t=s("tr"),e=s("td"),n=B(),a=s("td"),o=B();for(let P=0;P<b.length;P+=1)b[P].c();r=B(),c=s("td"),p=s("button"),p.textContent=T,d=B(),C=s("button"),C.textContent=I,this.h()},l(P){t=f(P,"TR",{id:!0});var R=N(t);e=f(R,"TD",{}),N(e).forEach(u),n=w(R),a=f(R,"TD",{}),N(a).forEach(u),o=w(R);for(let g=0;g<b.length;g+=1)b[g].l(R);r=w(R),c=f(R,"TD",{});var y=N(c);p=f(y,"BUTTON",{type:!0,"data-svelte-h":!0}),L(p)!=="svelte-1gpezhr"&&(p.textContent=T),d=w(y),C=f(y,"BUTTON",{type:!0,"data-svelte-h":!0}),L(C)!=="svelte-1wprywb"&&(C.textContent=I),y.forEach(u),R.forEach(u),this.h()},h(){A(p,"type","button"),A(C,"type","button"),A(t,"id","add-board")},m(P,R){m(P,t,R),h(t,e),h(t,n),h(t,a),h(t,o);for(let y=0;y<b.length;y+=1)b[y]&&b[y].m(t,null);h(t,r),h(t,c),h(c,p),h(c,d),h(c,C),D||(O=[q(p,"click",l[20]),q(C,"click",l[21])],D=!0)},p(P,R){if(R[0]&57345){j=M(P[13]);let y;for(y=0;y<j.length;y+=1){const g=Kt(P,j,y);b[y]?b[y].p(g,R):(b[y]=el(g),b[y].c(),b[y].m(t,r))}for(;y<b.length;y+=1)b[y].d(1);b.length=j.length}},d(P){P&&u(t),se(b,P),D=!1,_t(O)}}}function Ol(l){let t,e,n,a;function o(){l[36].call(e,l[57])}return{c(){t=s("td"),e=s("input"),this.h()},l(r){t=f(r,"TD",{});var c=N(t);e=f(c,"INPUT",{type:!0,placeholder:!0}),c.forEach(u),this.h()},h(){A(e,"type","text"),A(e,"placeholder",l[57]["display-name"])},m(r,c){m(r,t,c),h(t,e),G(e,l[0][l[57]["column-code"]]),n||(a=q(e,"input",o),n=!0)},p(r,c){l=r,c[0]&24577&&e.value!==l[0][l[57]["column-code"]]&&G(e,l[0][l[57]["column-code"]])},d(r){r&&u(t),n=!1,a()}}}function Pl(l){let t;return{c(){t=s("td")},l(e){t=f(e,"TD",{}),N(t).forEach(u)},m(e,n){m(e,t,n)},p:he,d(e){e&&u(t)}}}function Al(l){let t,e,n,a,o=M(Object.entries(l[14])),r=[];for(let p=0;p<o.length;p+=1)r[p]=$t(Wt(l,o,p));function c(){l[35].call(e,l[57])}return{c(){t=s("td"),e=s("select");for(let p=0;p<r.length;p+=1)r[p].c();this.h()},l(p){t=f(p,"TD",{});var T=N(t);e=f(T,"SELECT",{});var d=N(e);for(let C=0;C<r.length;C+=1)r[C].l(d);d.forEach(u),T.forEach(u),this.h()},h(){l[0][l[57]["column-code"]]===void 0&&We(c)},m(p,T){m(p,t,T),h(t,e);for(let d=0;d<r.length;d+=1)r[d]&&r[d].m(e,null);pe(e,l[0][l[57]["column-code"]],!0),n||(a=q(e,"change",c),n=!0)},p(p,T){if(l=p,T[0]&16384){o=M(Object.entries(l[14]));let d;for(d=0;d<o.length;d+=1){const C=Wt(l,o,d);r[d]?r[d].p(C,T):(r[d]=$t(C),r[d].c(),r[d].m(e,null))}for(;d<r.length;d+=1)r[d].d(1);r.length=o.length}T[0]&24577&&pe(e,l[0][l[57]["column-code"]])},d(p){p&&u(t),se(r,p),n=!1,a()}}}function Dl(l){let t,e,n,a;function o(){l[34].call(e,l[57])}return{c(){t=s("td"),e=s("input"),this.h()},l(r){t=f(r,"TD",{});var c=N(t);e=f(c,"INPUT",{}),c.forEach(u),this.h()},h(){e.disabled=!0},m(r,c){m(r,t,c),h(t,e),G(e,l[0][l[57]["column-code"]]),n||(a=q(e,"input",o),n=!0)},p(r,c){l=r,c[0]&24577&&e.value!==l[0][l[57]["column-code"]]&&G(e,l[0][l[57]["column-code"]])},d(r){r&&u(t),n=!1,a()}}}function Bl(l){let t,e,n,a="Board",o,r="Gallery",c,p;function T(){l[33].call(e,l[57])}return{c(){t=s("td"),e=s("select"),n=s("option"),n.textContent=a,o=s("option"),o.textContent=r,this.h()},l(d){t=f(d,"TD",{});var C=N(t);e=f(C,"SELECT",{});var I=N(e);n=f(I,"OPTION",{"data-svelte-h":!0}),L(n)!=="svelte-ueeluy"&&(n.textContent=a),o=f(I,"OPTION",{"data-svelte-h":!0}),L(o)!=="svelte-79feqe"&&(o.textContent=r),I.forEach(u),C.forEach(u),this.h()},h(){n.__value="board",G(n,n.__value),o.__value="gallery",G(o,o.__value),l[0][l[57]["column-code"]]===void 0&&We(T)},m(d,C){m(d,t,C),h(t,e),h(e,n),h(e,o),pe(e,l[0][l[57]["column-code"]],!0),c||(p=q(e,"change",T),c=!0)},p(d,C){l=d,C[0]&24577&&pe(e,l[0][l[57]["column-code"]])},d(d){d&&u(t),c=!1,p()}}}function wl(l){let t="",e;return{c(){e=fe(t)},l(n){e=_e(n,t)},m(n,a){m(n,e,a)},p:he,d(n){n&&u(e)}}}function $t(l){let t,e=l[61].name+"",n;return{c(){t=s("option"),n=fe(e),this.h()},l(a){t=f(a,"OPTION",{});var o=N(t);n=_e(o,e),o.forEach(u),this.h()},h(){t.__value=l[61].code,G(t,t.__value)},m(a,o){m(a,t,o),h(t,n)},p:he,d(a){a&&u(t)}}}function el(l){let t;function e(o,r){return o[57]["column-code"]=="idx"?wl:o[57]["column-code"]=="board-type"?Bl:o[57]["column-code"]=="board-table"||o[57]["column-code"]=="comment-table"?Dl:o[15].includes(o[57]["column-code"])?Al:o[57]["column-code"]=="regdate"?Pl:Ol}let a=e(l)(l);return{c(){a.c(),t=oe()},l(o){a.l(o),t=oe()},m(o,r){a.m(o,r),m(o,t,r)},p(o,r){a.p(o,r)},d(o){o&&u(t),a.d(o)}}}function Sl(l){let t,e,n,a=!1,o,r,c,p,T="View",d,C,I="Edit",D,O,j="Delete",b,P,R,y=M(l[13]),g=[];for(let S=0;S<y.length;S+=1)g[S]=nl(Xt(l,y,S));function E(){return l[43](l[56])}function we(){return l[44](l[56])}function ge(){return l[45](l[56])}return b=gl(l[42][0]),{c(){t=s("td"),e=s("input"),o=B();for(let S=0;S<g.length;S+=1)g[S].c();r=B(),c=s("td"),p=s("button"),p.textContent=T,d=B(),C=s("button"),C.textContent=I,D=B(),O=s("button"),O.textContent=j,this.h()},l(S){t=f(S,"TD",{});var X=N(t);e=f(X,"INPUT",{type:!0}),X.forEach(u),o=w(S);for(let V=0;V<g.length;V+=1)g[V].l(S);r=w(S),c=f(S,"TD",{});var U=N(c);p=f(U,"BUTTON",{type:!0,"data-svelte-h":!0}),L(p)!=="svelte-4gqy8p"&&(p.textContent=T),d=w(U),C=f(U,"BUTTON",{type:!0,"data-svelte-h":!0}),L(C)!=="svelte-1t2ew4d"&&(C.textContent=I),D=w(U),O=f(U,"BUTTON",{type:!0,"data-svelte-h":!0}),L(O)!=="svelte-kppxo3"&&(O.textContent=j),U.forEach(u),this.h()},h(){A(e,"type","checkbox"),e.__value=n=l[54].idx,G(e,e.__value),A(p,"type","button"),A(C,"type","button"),A(O,"type","button"),b.p(e)},m(S,X){m(S,t,X),h(t,e),e.checked=~(l[7]||[]).indexOf(e.__value),m(S,o,X);for(let U=0;U<g.length;U+=1)g[U]&&g[U].m(S,X);m(S,r,X),m(S,c,X),h(c,p),h(c,d),h(c,C),h(c,D),h(c,O),P||(R=[q(e,"change",l[41]),q(p,"click",E),q(C,"click",we),q(O,"click",ge)],P=!0)},p(S,X){if(l=S,X[0]&1024&&n!==(n=l[54].idx)&&(e.__value=n,G(e,e.__value),a=!0),(a||X[0]&1152)&&(e.checked=~(l[7]||[]).indexOf(e.__value)),X[0]&58368){y=M(l[13]);let U;for(U=0;U<y.length;U+=1){const V=Xt(l,y,U);g[U]?g[U].p(V,X):(g[U]=nl(V),g[U].c(),g[U].m(r.parentNode,r))}for(;U<g.length;U+=1)g[U].d(1);g.length=y.length}},d(S){S&&(u(t),u(o),u(r),u(c)),se(g,S),b.r(),P=!1,_t(R)}}}function Il(l){let t,e,n,a,o,r="Cancel",c,p,T="Save",d,C,I=M(l[13]),D=[];for(let O=0;O<I.length;O+=1)D[O]=al(Yt(l,I,O));return{c(){t=s("td"),e=B();for(let O=0;O<D.length;O+=1)D[O].c();n=B(),a=s("td"),o=s("button"),o.textContent=r,c=B(),p=s("button"),p.textContent=T,this.h()},l(O){t=f(O,"TD",{}),N(t).forEach(u),e=w(O);for(let b=0;b<D.length;b+=1)D[b].l(O);n=w(O),a=f(O,"TD",{});var j=N(a);o=f(j,"BUTTON",{type:!0,"data-svelte-h":!0}),L(o)!=="svelte-ufbtnv"&&(o.textContent=r),c=w(j),p=f(j,"BUTTON",{type:!0,"data-svelte-h":!0}),L(p)!=="svelte-x4yzzr"&&(p.textContent=T),j.forEach(u),this.h()},h(){A(o,"type","button"),A(p,"type","button")},m(O,j){m(O,t,j),m(O,e,j);for(let b=0;b<D.length;b+=1)D[b]&&D[b].m(O,j);m(O,n,j),m(O,a,j),h(a,o),h(a,c),h(a,p),d||(C=[q(o,"click",l[24]),q(p,"click",l[25])],d=!0)},p(O,j){if(j[0]&57346){I=M(O[13]);let b;for(b=0;b<I.length;b+=1){const P=Yt(O,I,b);D[b]?D[b].p(P,j):(D[b]=al(P),D[b].c(),D[b].m(n.parentNode,n))}for(;b<D.length;b+=1)D[b].d(1);D.length=I.length}},d(O){O&&(u(t),u(e),u(n),u(a)),se(D,O),d=!1,_t(C)}}}function jl(l){let t=l[54][l[57]["column-code"]]+"",e;return{c(){e=fe(t)},l(n){e=_e(n,t)},m(n,a){m(n,e,a)},p(n,a){a[0]&1024&&t!==(t=n[54][n[57]["column-code"]]+"")&&Qe(e,t)},d(n){n&&u(e)}}}function Ll(l){let t,e=M(Object.entries(l[14])),n=[];for(let a=0;a<e.length;a+=1)n[a]=ll(Vt(l,e,a));return{c(){for(let a=0;a<n.length;a+=1)n[a].c();t=oe()},l(a){for(let o=0;o<n.length;o+=1)n[o].l(a);t=oe()},m(a,o){for(let r=0;r<n.length;r+=1)n[r]&&n[r].m(a,o);m(a,t,o)},p(a,o){if(o[0]&25600){e=M(Object.entries(a[14]));let r;for(r=0;r<e.length;r+=1){const c=Vt(a,e,r);n[r]?n[r].p(c,o):(n[r]=ll(c),n[r].c(),n[r].m(t.parentNode,t))}for(;r<n.length;r+=1)n[r].d(1);n.length=e.length}},d(a){a&&u(t),se(n,a)}}}function tl(l){let t=l[61].name+"",e;return{c(){e=fe(t)},l(n){e=_e(n,t)},m(n,a){m(n,e,a)},p:he,d(n){n&&u(e)}}}function ll(l){let t,e=l[61].code==l[54][l[57]["column-code"]]&&tl(l);return{c(){e&&e.c(),t=oe()},l(n){e&&e.l(n),t=oe()},m(n,a){e&&e.m(n,a),m(n,t,a)},p(n,a){n[61].code==n[54][n[57]["column-code"]]?e?e.p(n,a):(e=tl(n),e.c(),e.m(t.parentNode,t)):e&&(e.d(1),e=null)},d(n){n&&u(t),e&&e.d(n)}}}function nl(l){let t;function e(o,r){return o[15].includes(o[57]["column-code"])?Ll:jl}let a=e(l)(l);return{c(){t=s("td"),a.c()},l(o){t=f(o,"TD",{});var r=N(t);a.l(r),r.forEach(u)},m(o,r){m(o,t,r),a.m(t,null)},p(o,r){a.p(o,r)},d(o){o&&u(t),a.d()}}}function Ul(l){let t,e,n,a;function o(){l[40].call(e,l[57])}return{c(){t=s("td"),e=s("input"),this.h()},l(r){t=f(r,"TD",{});var c=N(t);e=f(c,"INPUT",{type:!0,placeholder:!0}),c.forEach(u),this.h()},h(){A(e,"type","text"),A(e,"placeholder",l[57]["display-name"])},m(r,c){m(r,t,c),h(t,e),G(e,l[1][l[57]["column-code"]]),n||(a=q(e,"input",o),n=!0)},p(r,c){l=r,c[0]&24578&&e.value!==l[1][l[57]["column-code"]]&&G(e,l[1][l[57]["column-code"]])},d(r){r&&u(t),n=!1,a()}}}function Gl(l){let t;return{c(){t=s("td")},l(e){t=f(e,"TD",{}),N(t).forEach(u)},m(e,n){m(e,t,n)},p:he,d(e){e&&u(t)}}}function Rl(l){let t,e,n,a,o=M(Object.entries(l[14])),r=[];for(let p=0;p<o.length;p+=1)r[p]=ol(Ft(l,o,p));function c(){l[39].call(e,l[57])}return{c(){t=s("td"),e=s("select");for(let p=0;p<r.length;p+=1)r[p].c();this.h()},l(p){t=f(p,"TD",{});var T=N(t);e=f(T,"SELECT",{});var d=N(e);for(let C=0;C<r.length;C+=1)r[C].l(d);d.forEach(u),T.forEach(u),this.h()},h(){l[1][l[57]["column-code"]]===void 0&&We(c)},m(p,T){m(p,t,T),h(t,e);for(let d=0;d<r.length;d+=1)r[d]&&r[d].m(e,null);pe(e,l[1][l[57]["column-code"]],!0),n||(a=q(e,"change",c),n=!0)},p(p,T){if(l=p,T[0]&16384){o=M(Object.entries(l[14]));let d;for(d=0;d<o.length;d+=1){const C=Ft(l,o,d);r[d]?r[d].p(C,T):(r[d]=ol(C),r[d].c(),r[d].m(e,null))}for(;d<r.length;d+=1)r[d].d(1);r.length=o.length}T[0]&24578&&pe(e,l[1][l[57]["column-code"]])},d(p){p&&u(t),se(r,p),n=!1,a()}}}function zl(l){let t,e,n,a;function o(){l[38].call(e,l[57])}return{c(){t=s("td"),e=s("input"),this.h()},l(r){t=f(r,"TD",{});var c=N(t);e=f(c,"INPUT",{}),c.forEach(u),this.h()},h(){e.disabled=!0},m(r,c){m(r,t,c),h(t,e),G(e,l[1][l[57]["column-code"]]),n||(a=q(e,"input",o),n=!0)},p(r,c){l=r,c[0]&24578&&e.value!==l[1][l[57]["column-code"]]&&G(e,l[1][l[57]["column-code"]])},d(r){r&&u(t),n=!1,a()}}}function ql(l){let t,e,n,a="Board",o,r="Gallery",c,p;function T(){l[37].call(e,l[57])}return{c(){t=s("td"),e=s("select"),n=s("option"),n.textContent=a,o=s("option"),o.textContent=r,this.h()},l(d){t=f(d,"TD",{});var C=N(t);e=f(C,"SELECT",{});var I=N(e);n=f(I,"OPTION",{"data-svelte-h":!0}),L(n)!=="svelte-ueeluy"&&(n.textContent=a),o=f(I,"OPTION",{"data-svelte-h":!0}),L(o)!=="svelte-79feqe"&&(o.textContent=r),I.forEach(u),C.forEach(u),this.h()},h(){n.__value="board",G(n,n.__value),o.__value="gallery",G(o,o.__value),l[1][l[57]["column-code"]]===void 0&&We(T)},m(d,C){m(d,t,C),h(t,e),h(e,n),h(e,o),pe(e,l[1][l[57]["column-code"]],!0),c||(p=q(e,"change",T),c=!0)},p(d,C){l=d,C[0]&24578&&pe(e,l[1][l[57]["column-code"]])},d(d){d&&u(t),c=!1,p()}}}function Ml(l){let t,e=l[1].idx+"",n;return{c(){t=s("td"),n=fe(e)},l(a){t=f(a,"TD",{});var o=N(t);n=_e(o,e),o.forEach(u)},m(a,o){m(a,t,o),h(t,n)},p(a,o){o[0]&2&&e!==(e=a[1].idx+"")&&Qe(n,e)},d(a){a&&u(t)}}}function ol(l){let t,e=l[61].name+"",n;return{c(){t=s("option"),n=fe(e),this.h()},l(a){t=f(a,"OPTION",{});var o=N(t);n=_e(o,e),o.forEach(u),this.h()},h(){t.__value=l[61].code,G(t,t.__value)},m(a,o){m(a,t,o),h(t,n)},p:he,d(a){a&&u(t)}}}function al(l){let t;function e(o,r){return o[57]["column-code"]=="idx"?Ml:o[57]["column-code"]=="board-type"?ql:o[57]["column-code"]=="board-table"||o[57]["column-code"]=="comment-table"?zl:o[15].includes(o[57]["column-code"])?Rl:o[57]["column-code"]=="regdate"?Gl:Ul}let a=e(l)(l);return{c(){a.c(),t=oe()},l(o){a.l(o),t=oe()},m(o,r){a.m(o,r),m(o,t,r)},p(o,r){a.p(o,r)},d(o){o&&u(t),a.d(o)}}}function rl(l){let t,e;function n(r,c){return r[8]==r[56]?Il:Sl}let a=n(l),o=a(l);return{c(){t=s("tr"),o.c(),e=B()},l(r){t=f(r,"TR",{});var c=N(t);o.l(c),e=w(c),c.forEach(u)},m(r,c){m(r,t,c),o.m(t,null),h(t,e)},p(r,c){a===(a=n(r))&&o?o.p(r,c):(o.d(1),o=a(r),o&&(o.c(),o.m(t,e)))},d(r){r&&u(t),o.d()}}}function il(l){let t,e=l[49]+"",n,a;return{c(){t=s("a"),n=fe(e),this.h()},l(o){t=f(o,"A",{href:!0});var r=N(t);n=_e(r,e),r.forEach(u),this.h()},h(){A(t,"href",a=`?page=${l[49]}&list-count=${l[4]}`+(l[5]!=""?`&search=${l[5]}`:""))},m(o,r){m(o,t,r),h(t,n)},p(o,r){r[0]&4&&e!==(e=o[49]+"")&&Qe(n,e),r[0]&52&&a!==(a=`?page=${o[49]}&list-count=${o[4]}`+(o[5]!=""?`&search=${o[5]}`:""))&&A(t,"href",a)},d(o){o&&u(t)}}}function ul(l){let t,e=l[49]>=1&&il(l);return{c(){e&&e.c(),t=oe()},l(n){e&&e.l(n),t=oe()},m(n,a){e&&e.m(n,a),m(n,t,a)},p(n,a){n[49]>=1?e?e.p(n,a):(e=il(n),e.c(),e.m(t.parentNode,t)):e&&(e.d(1),e=null)},d(n){n&&u(t),e&&e.d(n)}}}function cl(l){let t,e=l[49]+"",n,a;return{c(){t=s("a"),n=fe(e),this.h()},l(o){t=f(o,"A",{href:!0});var r=N(t);n=_e(r,e),r.forEach(u),this.h()},h(){A(t,"href",a=`?page=${l[49]}&list-count=${l[4]}`+(l[5]!=""?`&search=${l[5]}`:""))},m(o,r){m(o,t,r),h(t,n)},p(o,r){r[0]&4&&e!==(e=o[49]+"")&&Qe(n,e),r[0]&52&&a!==(a=`?page=${o[49]}&list-count=${o[4]}`+(o[5]!=""?`&search=${o[5]}`:""))&&A(t,"href",a)},d(o){o&&u(t)}}}function sl(l){let t,e=l[49]<=l[3]&&cl(l);return{c(){e&&e.c(),t=oe()},l(n){e&&e.l(n),t=oe()},m(n,a){e&&e.m(n,a),m(n,t,a)},p(n,a){n[49]<=n[3]?e?e.p(n,a):(e=cl(n),e.c(),e.m(t.parentNode,t)):e&&(e.d(1),e=null)},d(n){n&&u(t),e&&e.d(n)}}}function Hl(l){let t,e="Boards",n,a,o="Add board",r,c,p="|",T,d,C="Delete selected boards",I,D,O="|",j,b,P="Search:",R,y,g,E,we="Search",ge,S,X="|",U,V,Ze="List:",Ee,Q,Ne,Z,Ce,te,ke,de,Ve,Se,Te,xe="Control",Ye,Ie,me,je,x,ae,$e="Admin",re,et="Manager",ie,tt="Regular user",v,$="Pending user",H,ue="Banned user",ce,yt="Guest",lt,ye,be,Et="Board",ve,Nt="Gallery",nt,z,Oe,Le,Ot="«",ot,dt,Pe,Ue,Pt="<",at,ht,Ge,At="..",pt,rt,Fe,it,mt,ut,Re,Dt="..",bt,Ae,ze,Bt=">",ct,vt,De,qe,wt="»",st,gt,St,fl=M([5,10,20,30,50,80]),Me=[];for(let i=0;i<6;i+=1)Me[i]=Nl(El(l,fl,i));let He=M(l[13]),F=[];for(let i=0;i<He.length;i+=1)F[i]=Zt(Qt(l,He,i));let W=l[9]&&xt(l),Je=M(l[10]),K=[];for(let i=0;i<Je.length;i+=1)K[i]=rl(Jt(l,Je,i));let Ct=M([l[2]-2,l[2]-1]),le=[];for(let i=0;i<2;i+=1)le[i]=ul(Ht(l,Ct,i));let kt=M([l[2]+1,l[2]+2]),ne=[];for(let i=0;i<2;i+=1)ne[i]=sl(Mt(l,kt,i));return{c(){t=s("h1"),t.textContent=e,n=B(),a=s("button"),a.textContent=o,r=B(),c=s("span"),c.textContent=p,T=B(),d=s("button"),d.textContent=C,I=B(),D=s("span"),D.textContent=O,j=B(),b=s("label"),b.textContent=P,R=B(),y=s("input"),g=B(),E=s("button"),E.textContent=we,ge=B(),S=s("span"),S.textContent=X,U=B(),V=s("label"),V.textContent=Ze,Ee=B(),Q=s("select");for(let i=0;i<6;i+=1)Me[i].c();Ne=B(),Z=s("table"),Ce=s("thead"),te=s("tr"),ke=s("th"),de=s("input"),Ve=B();for(let i=0;i<F.length;i+=1)F[i].c();Se=B(),Te=s("th"),Te.textContent=xe,Ye=B(),W&&W.c(),Ie=B(),me=s("tbody");for(let i=0;i<K.length;i+=1)K[i].c();je=B(),x=s("datalist"),ae=s("option"),ae.textContent=$e,re=s("option"),re.textContent=et,ie=s("option"),ie.textContent=tt,v=s("option"),v.textContent=$,H=s("option"),H.textContent=ue,ce=s("option"),ce.textContent=yt,lt=B(),ye=s("datalist"),be=s("option"),be.textContent=Et,ve=s("option"),ve.textContent=Nt,nt=B(),z=s("div"),Oe=s("a"),Le=s("span"),Le.textContent=Ot,dt=B(),Pe=s("a"),Ue=s("span"),Ue.textContent=Pt,ht=B(),Ge=s("span"),Ge.textContent=At,pt=B();for(let i=0;i<2;i+=1)le[i].c();rt=B(),Fe=s("b"),it=fe(l[2]),mt=B();for(let i=0;i<2;i+=1)ne[i].c();ut=B(),Re=s("span"),Re.textContent=Dt,bt=B(),Ae=s("a"),ze=s("span"),ze.textContent=Bt,vt=B(),De=s("a"),qe=s("span"),qe.textContent=wt,this.h()},l(i){t=f(i,"H1",{"data-svelte-h":!0}),L(t)!=="svelte-1um8syh"&&(t.textContent=e),n=w(i),a=f(i,"BUTTON",{type:!0,"data-svelte-h":!0}),L(a)!=="svelte-4lgeqm"&&(a.textContent=o),r=w(i),c=f(i,"SPAN",{"data-svelte-h":!0}),L(c)!=="svelte-1e2i4m"&&(c.textContent=p),T=w(i),d=f(i,"BUTTON",{type:!0,"data-svelte-h":!0}),L(d)!=="svelte-23zzud"&&(d.textContent=C),I=w(i),D=f(i,"SPAN",{"data-svelte-h":!0}),L(D)!=="svelte-1e2i4m"&&(D.textContent=O),j=w(i),b=f(i,"LABEL",{for:!0,"data-svelte-h":!0}),L(b)!=="svelte-1tn9osg"&&(b.textContent=P),R=w(i),y=f(i,"INPUT",{type:!0,id:!0,placeholder:!0}),g=w(i),E=f(i,"BUTTON",{type:!0,"data-svelte-h":!0}),L(E)!=="svelte-k4n4il"&&(E.textContent=we),ge=w(i),S=f(i,"SPAN",{"data-svelte-h":!0}),L(S)!=="svelte-1e2i4m"&&(S.textContent=X),U=w(i),V=f(i,"LABEL",{for:!0,"data-svelte-h":!0}),L(V)!=="svelte-17tcrp1"&&(V.textContent=Ze),Ee=w(i),Q=f(i,"SELECT",{id:!0});var k=N(Q);for(let J=0;J<6;J+=1)Me[J].l(k);k.forEach(u),Ne=w(i),Z=f(i,"TABLE",{id:!0});var _=N(Z);Ce=f(_,"THEAD",{});var ee=N(Ce);te=f(ee,"TR",{});var Xe=N(te);ke=f(Xe,"TH",{});var It=N(ke);de=f(It,"INPUT",{type:!0}),It.forEach(u),Ve=w(Xe);for(let J=0;J<F.length;J+=1)F[J].l(Xe);Se=w(Xe),Te=f(Xe,"TH",{"data-svelte-h":!0}),L(Te)!=="svelte-10rrrw5"&&(Te.textContent=xe),Xe.forEach(u),ee.forEach(u),Ye=w(_),W&&W.l(_),Ie=w(_),me=f(_,"TBODY",{id:!0});var jt=N(me);for(let J=0;J<K.length;J+=1)K[J].l(jt);jt.forEach(u),_.forEach(u),je=w(i),x=f(i,"DATALIST",{id:!0});var Be=N(x);ae=f(Be,"OPTION",{"data-svelte-h":!0}),L(ae)!=="svelte-17op260"&&(ae.textContent=$e),re=f(Be,"OPTION",{"data-svelte-h":!0}),L(re)!=="svelte-1ob7gsw"&&(re.textContent=et),ie=f(Be,"OPTION",{"data-svelte-h":!0}),L(ie)!=="svelte-oxxqbx"&&(ie.textContent=tt),v=f(Be,"OPTION",{"data-svelte-h":!0}),L(v)!=="svelte-163nbr3"&&(v.textContent=$),H=f(Be,"OPTION",{"data-svelte-h":!0}),L(H)!=="svelte-p0330j"&&(H.textContent=ue),ce=f(Be,"OPTION",{"data-svelte-h":!0}),L(ce)!=="svelte-1j47qc6"&&(ce.textContent=yt),Be.forEach(u),lt=w(i),ye=f(i,"DATALIST",{id:!0});var Tt=N(ye);be=f(Tt,"OPTION",{"data-svelte-h":!0}),L(be)!=="svelte-ueeluy"&&(be.textContent=Et),ve=f(Tt,"OPTION",{"data-svelte-h":!0}),L(ve)!=="svelte-79feqe"&&(ve.textContent=Nt),Tt.forEach(u),nt=w(i),z=f(i,"DIV",{id:!0});var Y=N(z);Oe=f(Y,"A",{href:!0});var Lt=N(Oe);Le=f(Lt,"SPAN",{"data-svelte-h":!0}),L(Le)!=="svelte-1z054it"&&(Le.textContent=Ot),Lt.forEach(u),dt=w(Y),Pe=f(Y,"A",{href:!0});var Ut=N(Pe);Ue=f(Ut,"SPAN",{"data-svelte-h":!0}),L(Ue)!=="svelte-1kd6by1"&&(Ue.textContent=Pt),Ut.forEach(u),ht=w(Y),Ge=f(Y,"SPAN",{"data-svelte-h":!0}),L(Ge)!=="svelte-1v1zlza"&&(Ge.textContent=At),pt=w(Y);for(let J=0;J<2;J+=1)le[J].l(Y);rt=w(Y),Fe=f(Y,"B",{});var Gt=N(Fe);it=_e(Gt,l[2]),Gt.forEach(u),mt=w(Y);for(let J=0;J<2;J+=1)ne[J].l(Y);ut=w(Y),Re=f(Y,"SPAN",{"data-svelte-h":!0}),L(Re)!=="svelte-1v1zlza"&&(Re.textContent=Dt),bt=w(Y),Ae=f(Y,"A",{href:!0});var Rt=N(Ae);ze=f(Rt,"SPAN",{"data-svelte-h":!0}),L(ze)!=="svelte-x0xyl0"&&(ze.textContent=Bt),Rt.forEach(u),vt=w(Y),De=f(Y,"A",{href:!0});var zt=N(De);qe=f(zt,"SPAN",{"data-svelte-h":!0}),L(qe)!=="svelte-131q397"&&(qe.textContent=wt),zt.forEach(u),Y.forEach(u),this.h()},h(){A(a,"type","button"),A(d,"type","button"),A(b,"for","search"),A(y,"type","text"),A(y,"id","search"),A(y,"placeholder","Search for..."),A(E,"type","button"),A(V,"for","set-list-count"),A(Q,"id","set-list-count"),l[4]===void 0&&We(()=>l[31].call(Q)),A(de,"type","checkbox"),A(me,"id","boards-list-body"),A(Z,"id","boards-list-container"),ae.__value="admin",G(ae,ae.__value),re.__value="manager",G(re,re.__value),ie.__value="user_active",G(ie,ie.__value),v.__value="user_hold",G(v,v.__value),H.__value="user_banned",G(H,H.__value),ce.__value="guest",G(ce,ce.__value),A(x,"id","grant-list"),be.__value="board",G(be,be.__value),ve.__value="gallery",G(ve,ve.__value),A(ye,"id","board-types"),A(Oe,"href",ot=`?page=1&list-count=${l[4]}`+(l[5]!=""?`&search=${l[5]}`:"")),A(Pe,"href",at=`?page=${l[12]}&list-count=${l[4]}`+(l[5]!=""?`&search=${l[5]}`:"")),A(Ae,"href",ct=`?page=${l[11]}&list-count=${l[4]}`+(l[5]!=""?`&search=${l[5]}`:"")),A(De,"href",st=`?page=${l[3]}&list-count=${l[4]}`+(l[5]!=""?`&search=${l[5]}`:"")),A(z,"id","page-container")},m(i,k){m(i,t,k),m(i,n,k),m(i,a,k),m(i,r,k),m(i,c,k),m(i,T,k),m(i,d,k),m(i,I,k),m(i,D,k),m(i,j,k),m(i,b,k),m(i,R,k),m(i,y,k),G(y,l[5]),m(i,g,k),m(i,E,k),m(i,ge,k),m(i,S,k),m(i,U,k),m(i,V,k),m(i,Ee,k),m(i,Q,k);for(let _=0;_<6;_+=1)Me[_]&&Me[_].m(Q,null);pe(Q,l[4],!0),m(i,Ne,k),m(i,Z,k),h(Z,Ce),h(Ce,te),h(te,ke),h(ke,de),de.checked=l[6],h(te,Ve);for(let _=0;_<F.length;_+=1)F[_]&&F[_].m(te,null);h(te,Se),h(te,Te),h(Z,Ye),W&&W.m(Z,null),h(Z,Ie),h(Z,me);for(let _=0;_<K.length;_+=1)K[_]&&K[_].m(me,null);m(i,je,k),m(i,x,k),h(x,ae),h(x,re),h(x,ie),h(x,v),h(x,H),h(x,ce),m(i,lt,k),m(i,ye,k),h(ye,be),h(ye,ve),m(i,nt,k),m(i,z,k),h(z,Oe),h(Oe,Le),h(z,dt),h(z,Pe),h(Pe,Ue),h(z,ht),h(z,Ge),h(z,pt);for(let _=0;_<2;_+=1)le[_]&&le[_].m(z,null);h(z,rt),h(z,Fe),h(Fe,it),h(z,mt);for(let _=0;_<2;_+=1)ne[_]&&ne[_].m(z,null);h(z,ut),h(z,Re),h(z,bt),h(z,Ae),h(Ae,ze),h(z,vt),h(z,De),h(De,qe),gt||(St=[q(a,"click",l[29]),q(d,"click",l[27]),q(y,"input",l[30]),q(y,"keypress",l[18]),q(E,"click",l[17]),q(Q,"change",l[31]),q(Q,"change",l[19]),q(de,"change",l[32]),q(de,"change",l[16])],gt=!0)},p(i,k){if(k[0]&32&&y.value!==i[5]&&G(y,i[5]),k[0]&16&&pe(Q,i[4]),k[0]&64&&(de.checked=i[6]),k[0]&8192){He=M(i[13]);let _;for(_=0;_<He.length;_+=1){const ee=Qt(i,He,_);F[_]?F[_].p(ee,k):(F[_]=Zt(ee),F[_].c(),F[_].m(te,Se))}for(;_<F.length;_+=1)F[_].d(1);F.length=He.length}if(i[9]?W?W.p(i,k):(W=xt(i),W.c(),W.m(Z,Ie)):W&&(W.d(1),W=null),k[0]&130082178){Je=M(i[10]);let _;for(_=0;_<Je.length;_+=1){const ee=Jt(i,Je,_);K[_]?K[_].p(ee,k):(K[_]=rl(ee),K[_].c(),K[_].m(me,null))}for(;_<K.length;_+=1)K[_].d(1);K.length=Je.length}if(k[0]&48&&ot!==(ot=`?page=1&list-count=${i[4]}`+(i[5]!=""?`&search=${i[5]}`:""))&&A(Oe,"href",ot),k[0]&4144&&at!==(at=`?page=${i[12]}&list-count=${i[4]}`+(i[5]!=""?`&search=${i[5]}`:""))&&A(Pe,"href",at),k[0]&52){Ct=M([i[2]-2,i[2]-1]);let _;for(_=0;_<2;_+=1){const ee=Ht(i,Ct,_);le[_]?le[_].p(ee,k):(le[_]=ul(ee),le[_].c(),le[_].m(z,rt))}for(;_<2;_+=1)le[_].d(1)}if(k[0]&4&&Qe(it,i[2]),k[0]&60){kt=M([i[2]+1,i[2]+2]);let _;for(_=0;_<2;_+=1){const ee=Mt(i,kt,_);ne[_]?ne[_].p(ee,k):(ne[_]=sl(ee),ne[_].c(),ne[_].m(z,ut))}for(;_<2;_+=1)ne[_].d(1)}k[0]&2096&&ct!==(ct=`?page=${i[11]}&list-count=${i[4]}`+(i[5]!=""?`&search=${i[5]}`:""))&&A(Ae,"href",ct),k[0]&56&&st!==(st=`?page=${i[3]}&list-count=${i[4]}`+(i[5]!=""?`&search=${i[5]}`:""))&&A(De,"href",st)},i:he,o:he,d(i){i&&(u(t),u(n),u(a),u(r),u(c),u(T),u(d),u(I),u(D),u(j),u(b),u(R),u(y),u(g),u(E),u(ge),u(S),u(U),u(V),u(Ee),u(Q),u(Ne),u(Z),u(je),u(x),u(lt),u(ye),u(nt),u(z)),se(Me,i),se(F,i),W&&W.d(),se(K,i),se(le,i),se(ne,i),gt=!1,_t(St)}}}async function Jl(l,t,e){let n={},a=`/api/admin/board?page=${l}&list-count=${t}`;e!=""&&(a+=`&search=${e}`);const o=await fetch(a,{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});return o.ok&&(n=await o.json()),n["board-list"]==null&&(n["board-list"]=[]),n}function Xl(l,t,e){let n,a,o,r,c,p;dl(l,Cl,v=>e(47,p=v));let{data:T}=t;const d=T["default-count"],C=T.columns,I=T.grades;let D=Number(p.url.searchParams.get("list-count"))||d,O=-1,j=p.url.searchParams.get("search")||"",b=!1,P=[],R=-1,y=!1,g={},E={};const we=["grant-read","grant-write","grant-comment","grant-upload"];function ge(){b?e(7,P=[]):e(7,P=n.map(v=>v.idx))}async function S(){const v=await Jl(1,D,j);e(28,T["boardlist-data"]=v,T)}function X(v){v.key==="Enter"&&S()}function U(){location.href=`?list-count=${D}`+(j!=""?`&search=${j}`:"")}function V(){e(0,g={}),e(9,y=!1)}async function Ze(){const $=await fetch("/api/admin/board",{method:"POST",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(g)});$.ok||alert(await $.text()),V(),ft()}function Ee(v){window.open("/board/list?board_code="+n[v]["board-code"],"_blank")}function Q(v){e(8,R=v),e(1,E={});for(const $ in n[v])e(1,E[$]=n[v][$],E);console.log(E)}function Ne(){e(1,E={}),e(8,R=-1)}async function Z(){for(const H in E)typeof E[H]==null||E[H]==null||typeof E[H]!="string"&&e(1,E[H]=E[H].toString(),E);const $=await fetch("/api/admin/board",{method:"PUT",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([E])});$.ok||alert(await $.text()),Ne(),ft()}async function Ce(v){const $=typeof n[v].idx=="number"?n[v].idx.toString():n[v].idx,ue=await fetch("/api/admin/board",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([{idx:$}])});ue.ok||alert(await ue.text()),e(7,P=[]),ft()}async function te(){if(P.length==0){alert("Selected nothing");return}const v=[];for(let ue=0;ue<P.length;ue++){const ce=typeof P[ue]=="number"?P[ue].toString():P[ue];v.push({idx:ce})}const H=await fetch("/api/admin/board",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(v)});H.ok||alert(await H.text()),e(7,P=[]),ft()}hl(()=>{}),pl(()=>{O!=a&&(e(7,P=[]),O=a),e(6,b=P.length==n.length)});const ke=[[]],de=()=>{e(0,g={"board-type":"board"}),e(9,y=!0)};function Ve(){j=this.value,e(5,j)}function Se(){D=Ke(this),e(4,D)}function Te(){b=this.checked,e(6,b)}function xe(v){g[v["column-code"]]=Ke(this),e(0,g),e(1,E),e(14,I)}function Ye(v){g[v["column-code"]]=this.value,e(0,g),e(1,E),e(14,I)}function Ie(v){g[v["column-code"]]=Ke(this),e(0,g),e(1,E),e(14,I)}function me(v){g[v["column-code"]]=this.value,e(0,g),e(1,E),e(14,I)}function je(v){E[v["column-code"]]=Ke(this),e(1,E),e(0,g),e(14,I)}function x(v){E[v["column-code"]]=this.value,e(1,E),e(0,g),e(14,I)}function ae(v){E[v["column-code"]]=Ke(this),e(1,E),e(0,g),e(14,I)}function $e(v){E[v["column-code"]]=this.value,e(1,E),e(0,g),e(14,I)}function re(){P=vl(ke[0],this.__value,this.checked),e(7,P)}const et=v=>{Ee(v)},ie=v=>{Q(v)},tt=v=>{Ce(v)};return l.$$set=v=>{"data"in v&&e(28,T=v.data)},l.$$.update=()=>{l.$$.dirty[0]&268435456&&e(10,n=T["boardlist-data"]["board-list"]),l.$$.dirty[0]&268435456&&e(2,a=T["boardlist-data"]["current-page"]),l.$$.dirty[0]&268435456&&e(3,o=T["boardlist-data"]["total-page"]),l.$$.dirty[0]&4&&e(12,r=a-5>1?a-5:1),l.$$.dirty[0]&12&&e(11,c=a+5<o?a+5:o),l.$$.dirty[0]&3&&(g["board-code"]!=null&&g["board-code"].length>0&&(e(0,g["board-table"]="board_"+g["board-code"].toLowerCase(),g),e(0,g["comment-table"]="comment_"+g["board-code"].toLowerCase(),g)),E["board-code"]!=null&&E["board-code"].length>0&&(e(1,E["board-table"]="board_"+E["board-code"].toLowerCase(),E),e(1,E["comment-table"]="comment_"+E["board-code"].toLowerCase(),E)))},[g,E,a,o,D,j,b,P,R,y,n,c,r,C,I,we,ge,S,X,U,V,Ze,Ee,Q,Ne,Z,Ce,te,T,de,Ve,Se,Te,xe,Ye,Ie,me,je,x,ae,$e,re,ke,et,ie,tt]}class Ql extends ml{constructor(t){super(),bl(this,t,Xl,Hl,_l,{data:28},null,[-1,-1,-1])}}export{Ql as component,Wl as universal};
