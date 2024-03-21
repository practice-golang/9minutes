import{s as Dt,n as ut,r as it,f as Nt}from"../chunks/scheduler.b95eede2.js";import{S as Ot,i as It,g as u,s as E,h as r,y as Z,c as g,j as D,f as h,k as p,a as Q,z as l,B as X,C as Et,A as V,F as _t,m as et,n as lt,e as mt,D as Ut,H as gt,E as St,o as ot,G as wt}from"../chunks/index.e4e63e43.js";import{e as st,i as rt}from"../chunks/table.94dca04b.js";const jt=async({fetch:n})=>{let t=[];const e=await n("/api/admin/user-columns",{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});return e.ok&&(t=await e.json()),{columns:t}},qt=Object.freeze(Object.defineProperty({__proto__:null,load:jt},Symbol.toStringTag,{value:"Module"}));function vt(n,t,e){const s=n.slice();return s[30]=t[e],s[32]=e,s}function yt(n,t,e){const s=n.slice();return s[33]=t[e][0],s[34]=t[e][1],s}function bt(n){let t,e,s,i,a,_,y,f,b,v,d,I,T,U,O,L,M,w,F,j="Cancel",A,z,R="Save",J,tt,B=st(Object.entries(n[7])),S=[];for(let C=0;C<B.length;C+=1)S[C]=Tt(yt(n,B,C));return{c(){t=u("tr"),e=u("td"),s=E(),i=u("td"),a=u("input"),_=E(),y=u("td"),f=u("input"),b=E(),v=u("td"),d=u("select");for(let C=0;C<S.length;C+=1)S[C].c();I=E(),T=u("td"),U=u("input"),O=E(),L=u("td"),M=E(),w=u("td"),F=u("button"),F.textContent=j,A=E(),z=u("button"),z.textContent=R,this.h()},l(C){t=r(C,"TR",{id:!0});var c=D(t);e=r(c,"TD",{}),D(e).forEach(h),s=g(c),i=r(c,"TD",{});var o=D(i);a=r(o,"INPUT",{type:!0,name:!0,placeholder:!0}),o.forEach(h),_=g(c),y=r(c,"TD",{});var m=D(y);f=r(m,"INPUT",{type:!0,name:!0,placeholder:!0}),m.forEach(h),b=g(c),v=r(c,"TD",{});var G=D(v);d=r(G,"SELECT",{name:!0});var N=D(d);for(let P=0;P<S.length;P+=1)S[P].l(N);N.forEach(h),G.forEach(h),I=g(c),T=r(c,"TD",{});var H=D(T);U=r(H,"INPUT",{type:!0,name:!0,placeholder:!0}),H.forEach(h),O=g(c),L=r(c,"TD",{}),D(L).forEach(h),M=g(c),w=r(c,"TD",{});var q=D(w);F=r(q,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(F)!=="svelte-b19viv"&&(F.textContent=j),A=g(q),z=r(q,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(z)!=="svelte-1jozejf"&&(z.textContent=R),q.forEach(h),c.forEach(h),this.h()},h(){p(a,"type","text"),p(a,"name","display-name"),p(a,"placeholder","Display name"),p(f,"type","text"),p(f,"name","column-code"),p(f,"placeholder","Column code"),p(d,"name","column-type"),n[4]["column-type"]===void 0&&Nt(()=>n[20].call(d)),p(U,"type","text"),p(U,"name","column-name"),p(U,"placeholder","Column name"),p(F,"type","button"),p(z,"type","button"),p(t,"id","add-column")},m(C,c){Q(C,t,c),l(t,e),l(t,s),l(t,i),l(i,a),V(a,n[4]["display-name"]),l(t,_),l(t,y),l(y,f),V(f,n[4]["column-code"]),l(t,b),l(t,v),l(v,d);for(let o=0;o<S.length;o+=1)S[o]&&S[o].m(d,null);_t(d,n[4]["column-type"],!0),l(t,I),l(t,T),l(T,U),V(U,n[4]["column-name"]),l(t,O),l(t,L),l(t,M),l(t,w),l(w,F),l(w,A),l(w,z),J||(tt=[X(a,"input",n[18]),X(f,"input",n[19]),X(d,"change",n[20]),X(U,"input",n[21]),X(F,"click",n[9]),X(z,"click",n[10])],J=!0)},p(C,c){if(c[0]&144&&a.value!==C[4]["display-name"]&&V(a,C[4]["display-name"]),c[0]&144&&f.value!==C[4]["column-code"]&&V(f,C[4]["column-code"]),c[0]&128){B=st(Object.entries(C[7]));let o;for(o=0;o<B.length;o+=1){const m=yt(C,B,o);S[o]?S[o].p(m,c):(S[o]=Tt(m),S[o].c(),S[o].m(d,null))}for(;o<S.length;o+=1)S[o].d(1);S.length=B.length}c[0]&144&&_t(d,C[4]["column-type"]),c[0]&144&&U.value!==C[4]["column-name"]&&V(U,C[4]["column-name"])},d(C){C&&h(t),Et(S,C),J=!1,it(tt)}}}function Tt(n){let t,e=n[34]+"",s;return{c(){t=u("option"),s=et(e),this.h()},l(i){t=r(i,"OPTION",{});var a=D(t);s=lt(a,e),a.forEach(h),this.h()},h(){t.__value=n[33],V(t,t.__value)},m(i,a){Q(i,t,a),l(t,s)},p:ut,d(i){i&&h(t)}}}function Pt(n){let t,e,s,i=n[30]["display-name"]+"",a,_,y,f=n[30]["column-code"]+"",b,v,d,I=n[7][n[30]["column-type"]]+"",T,U,O,L=n[30]["column-name"]+"",M,w,F,j=n[30]["sort-order"]+"",A,z,R;function J(o,m){return o[32]<=6?Ht:Bt}let B=J(n)(n);function S(o,m){return o[32]<=6?Mt:xt}let c=S(n)(n);return{c(){t=u("tr"),B.c(),e=E(),s=u("td"),a=et(i),_=E(),y=u("td"),b=et(f),v=E(),d=u("td"),T=et(I),U=E(),O=u("td"),M=et(L),w=E(),F=u("td"),A=et(j),z=E(),c.c(),R=E(),this.h()},l(o){t=r(o,"TR",{});var m=D(t);B.l(m),e=g(m),s=r(m,"TD",{class:!0});var G=D(s);a=lt(G,i),G.forEach(h),_=g(m),y=r(m,"TD",{class:!0});var N=D(y);b=lt(N,f),N.forEach(h),v=g(m),d=r(m,"TD",{class:!0});var H=D(d);T=lt(H,I),H.forEach(h),U=g(m),O=r(m,"TD",{class:!0});var q=D(O);M=lt(q,L),q.forEach(h),w=g(m),F=r(m,"TD",{class:!0});var P=D(F);A=lt(P,j),P.forEach(h),z=g(m),c.l(m),R=g(m),m.forEach(h),this.h()},h(){p(s,"class","colField"),p(y,"class","colField"),p(d,"class","colField"),p(O,"class","colField"),p(F,"class","colFixedMid")},m(o,m){Q(o,t,m),B.m(t,null),l(t,e),l(t,s),l(s,a),l(t,_),l(t,y),l(y,b),l(t,v),l(t,d),l(d,T),l(t,U),l(t,O),l(O,M),l(t,w),l(t,F),l(F,A),l(t,z),c.m(t,null),l(t,R)},p(o,m){B.p(o,m),m[0]&2&&i!==(i=o[30]["display-name"]+"")&&ot(a,i),m[0]&2&&f!==(f=o[30]["column-code"]+"")&&ot(b,f),m[0]&2&&I!==(I=o[7][o[30]["column-type"]]+"")&&ot(T,I),m[0]&2&&L!==(L=o[30]["column-name"]+"")&&ot(M,L),m[0]&2&&j!==(j=o[30]["sort-order"]+"")&&ot(A,j),c.p(o,m)},d(o){o&&h(t),B.d(),c.d()}}}function Ft(n){let t,e,s,i,a,_,y,f,b,v,d=n[7][n[30]["column-type"]]+"",I,T,U,O,L,M,w,F,j,A,z="Cancel",R,J,tt="Save",B,S,C;return{c(){t=u("tr"),e=u("td"),s=E(),i=u("td"),a=u("input"),_=E(),y=u("td"),f=u("input"),b=E(),v=u("td"),I=et(d),T=E(),U=u("td"),O=u("input"),L=E(),M=u("td"),w=u("input"),F=E(),j=u("td"),A=u("button"),A.textContent=z,R=E(),J=u("button"),J.textContent=tt,B=E(),this.h()},l(c){t=r(c,"TR",{});var o=D(t);e=r(o,"TD",{class:!0}),D(e).forEach(h),s=g(o),i=r(o,"TD",{});var m=D(i);a=r(m,"INPUT",{type:!0,name:!0,placeholder:!0}),m.forEach(h),_=g(o),y=r(o,"TD",{});var G=D(y);f=r(G,"INPUT",{type:!0,name:!0,placeholder:!0}),G.forEach(h),b=g(o),v=r(o,"TD",{});var N=D(v);I=lt(N,d),N.forEach(h),T=g(o),U=r(o,"TD",{});var H=D(U);O=r(H,"INPUT",{type:!0,name:!0,placeholder:!0}),H.forEach(h),L=g(o),M=r(o,"TD",{});var q=D(M);w=r(q,"INPUT",{type:!0,name:!0}),q.forEach(h),F=g(o),j=r(o,"TD",{class:!0});var P=D(j);A=r(P,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(A)!=="svelte-15i7fgb"&&(A.textContent=z),R=g(P),J=r(P,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(J)!=="svelte-1mm80i1"&&(J.textContent=tt),P.forEach(h),B=g(o),o.forEach(h),this.h()},h(){p(e,"class","colFixedMin"),p(a,"type","text"),p(a,"name","display-name"),p(a,"placeholder","Display name"),p(f,"type","text"),p(f,"name","column-code"),p(f,"placeholder","Column code"),p(O,"type","text"),p(O,"name","column-name"),p(O,"placeholder","Column name"),p(w,"type","number"),p(w,"name","sort-order"),p(A,"type","button"),p(J,"type","button"),p(j,"class","colFixedMax")},m(c,o){Q(c,t,o),l(t,e),l(t,s),l(t,i),l(i,a),V(a,n[5]["display-name"]),l(t,_),l(t,y),l(y,f),V(f,n[5]["column-code"]),l(t,b),l(t,v),l(v,I),l(t,T),l(t,U),l(U,O),V(O,n[5]["column-name"]),l(t,L),l(t,M),l(M,w),V(w,n[5]["sort-order"]),l(t,F),l(t,j),l(j,A),l(j,R),l(j,J),l(t,B),S||(C=[X(a,"input",n[22]),X(f,"input",n[23]),X(O,"input",n[24]),X(w,"input",n[25]),X(A,"click",n[12]),X(J,"click",n[13])],S=!0)},p(c,o){o[0]&32&&a.value!==c[5]["display-name"]&&V(a,c[5]["display-name"]),o[0]&32&&f.value!==c[5]["column-code"]&&V(f,c[5]["column-code"]),o[0]&2&&d!==(d=c[7][c[30]["column-type"]]+"")&&ot(I,d),o[0]&32&&O.value!==c[5]["column-name"]&&V(O,c[5]["column-name"]),o[0]&32&&gt(w.value)!==c[5]["sort-order"]&&V(w,c[5]["sort-order"])},d(c){c&&h(t),S=!1,it(C)}}}function Bt(n){let t,e,s,i=!1,a,_,y;return a=wt(n[27][0]),{c(){t=u("td"),e=u("input"),this.h()},l(f){t=r(f,"TD",{class:!0});var b=D(t);e=r(b,"INPUT",{type:!0}),b.forEach(h),this.h()},h(){p(e,"type","checkbox"),e.__value=s=n[30].idx,V(e,e.__value),p(t,"class","colFixedMin"),a.p(e)},m(f,b){Q(f,t,b),l(t,e),e.checked=~(n[0]||[]).indexOf(e.__value),_||(y=X(e,"change",n[26]),_=!0)},p(f,b){b[0]&2&&s!==(s=f[30].idx)&&(e.__value=s,V(e,e.__value),i=!0),(i||b[0]&3)&&(e.checked=~(f[0]||[]).indexOf(e.__value))},d(f){f&&h(t),a.r(),_=!1,y()}}}function Ht(n){let t;return{c(){t=u("td"),this.h()},l(e){t=r(e,"TD",{class:!0}),D(t).forEach(h),this.h()},h(){p(t,"class","colFixedMin")},m(e,s){Q(e,t,s)},p:ut,d(e){e&&h(t)}}}function xt(n){let t,e,s="Edit",i,a,_="Delete",y,f;function b(){return n[28](n[32])}function v(){return n[29](n[32])}return{c(){t=u("td"),e=u("button"),e.textContent=s,i=E(),a=u("button"),a.textContent=_,this.h()},l(d){t=r(d,"TD",{class:!0});var I=D(t);e=r(I,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(e)!=="svelte-g4pdh8"&&(e.textContent=s),i=g(I),a=r(I,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(a)!=="svelte-a8j1fe"&&(a.textContent=_),I.forEach(h),this.h()},h(){p(e,"type","button"),p(a,"type","button"),p(t,"class","colFixedMax")},m(d,I){Q(d,t,I),l(t,e),l(t,i),l(t,a),y||(f=[X(e,"click",b),X(a,"click",v)],y=!0)},p(d,I){n=d},d(d){d&&h(t),y=!1,it(f)}}}function Mt(n){let t;return{c(){t=u("td"),this.h()},l(e){t=r(e,"TD",{class:!0}),D(t).forEach(h),this.h()},h(){p(t,"class","colFixedMax")},m(e,s){Q(e,t,s)},p:ut,d(e){e&&h(t)}}}function Ct(n){let t;function e(a,_){return a[2]==a[32]?Ft:Pt}let s=e(n),i=s(n);return{c(){i.c(),t=mt()},l(a){i.l(a),t=mt()},m(a,_){i.m(a,_),Q(a,t,_)},p(a,_){s===(s=e(a))&&i?i.p(a,_):(i.d(1),i=s(a),i&&(i.c(),i.m(t.parentNode,t)))},d(a){a&&h(t),i.d(a)}}}function At(n){let t,e="User columns",s,i,a="Add column",_,y,f="Delete selected columns",b,v,d,I,T,U,O,L,M,w="Display name",F,j,A="Column code",z,R,J="Column type",tt,B,S="Column name",C,c,o="Sort Order",m,G,N="Control",H,q,P,ct,ht,K=n[3]&&bt(n),nt=st(n[1]),Y=[];for(let k=0;k<nt.length;k+=1)Y[k]=Ct(vt(n,nt,k));return{c(){t=u("h1"),t.textContent=e,s=E(),i=u("button"),i.textContent=a,_=E(),y=u("button"),y.textContent=f,b=E(),v=u("div"),d=u("table"),I=u("thead"),T=u("tr"),U=u("td"),O=u("input"),L=E(),M=u("th"),M.textContent=w,F=E(),j=u("th"),j.textContent=A,z=E(),R=u("th"),R.textContent=J,tt=E(),B=u("th"),B.textContent=S,C=E(),c=u("th"),c.textContent=o,m=E(),G=u("th"),G.textContent=N,H=E(),K&&K.c(),q=E(),P=u("tbody");for(let k=0;k<Y.length;k+=1)Y[k].c();this.h()},l(k){t=r(k,"H1",{"data-svelte-h":!0}),Z(t)!=="svelte-1qzb11o"&&(t.textContent=e),s=g(k),i=r(k,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(i)!=="svelte-13ld9q5"&&(i.textContent=a),_=g(k),y=r(k,"BUTTON",{type:!0,"data-svelte-h":!0}),Z(y)!=="svelte-1kxuovt"&&(y.textContent=f),b=g(k),v=r(k,"DIV",{class:!0});var W=D(v);d=r(W,"TABLE",{id:!0});var x=D(d);I=r(x,"THEAD",{});var at=D(I);T=r(at,"TR",{});var $=D(T);U=r($,"TD",{});var ft=D(U);O=r(ft,"INPUT",{type:!0}),ft.forEach(h),L=g($),M=r($,"TH",{"data-svelte-h":!0}),Z(M)!=="svelte-5fbpcz"&&(M.textContent=w),F=g($),j=r($,"TH",{"data-svelte-h":!0}),Z(j)!=="svelte-5edxmv"&&(j.textContent=A),z=g($),R=r($,"TH",{"data-svelte-h":!0}),Z(R)!=="svelte-rfc7ms"&&(R.textContent=J),tt=g($),B=r($,"TH",{"data-svelte-h":!0}),Z(B)!=="svelte-m0r1hp"&&(B.textContent=S),C=g($),c=r($,"TH",{"data-svelte-h":!0}),Z(c)!=="svelte-30jp5e"&&(c.textContent=o),m=g($),G=r($,"TH",{"data-svelte-h":!0}),Z(G)!=="svelte-10rrrw5"&&(G.textContent=N),$.forEach(h),at.forEach(h),H=g(x),K&&K.l(x),q=g(x),P=r(x,"TBODY",{id:!0});var pt=D(P);for(let dt=0;dt<Y.length;dt+=1)Y[dt].l(pt);pt.forEach(h),x.forEach(h),W.forEach(h),this.h()},h(){p(i,"type","button"),p(y,"type","button"),p(O,"type","checkbox"),O.checked=n[6],p(P,"id","column-list-body"),p(d,"id","column-list-container"),p(v,"class","table-container")},m(k,W){Q(k,t,W),Q(k,s,W),Q(k,i,W),Q(k,_,W),Q(k,y,W),Q(k,b,W),Q(k,v,W),l(v,d),l(d,I),l(I,T),l(T,U),l(U,O),l(T,L),l(T,M),l(T,F),l(T,j),l(T,z),l(T,R),l(T,tt),l(T,B),l(T,C),l(T,c),l(T,m),l(T,G),l(d,H),K&&K.m(d,null),l(d,q),l(d,P);for(let x=0;x<Y.length;x+=1)Y[x]&&Y[x].m(P,null);ct||(ht=[X(i,"click",n[17]),X(y,"click",n[15]),X(O,"click",n[8])],ct=!0)},p(k,W){if(W[0]&64&&(O.checked=k[6]),k[3]?K?K.p(k,W):(K=bt(k),K.c(),K.m(d,q)):K&&(K.d(1),K=null),W[0]&30887){nt=st(k[1]);let x;for(x=0;x<nt.length;x+=1){const at=vt(k,nt,x);Y[x]?Y[x].p(at,W):(Y[x]=Ct(at),Y[x].c(),Y[x].m(P,null))}for(;x<Y.length;x+=1)Y[x].d(1);Y.length=nt.length}},i:ut,o:ut,d(k){k&&(h(t),h(s),h(i),h(_),h(y),h(b),h(v)),K&&K.d(),Et(Y,k),ct=!1,it(ht)}}}const kt=7;function Lt(n,t,e){let s,i,{data:a}=t,_=[],y=-1,f=!1,b={},v={};const d={text:"Text","number-integer":"Number Integer","number-real":"Number Real"};function I(){i?e(0,_=[]):e(0,_=s.filter(N=>N.idx>kt).map(N=>N.idx))}function T(){e(4,b={}),e(3,f=!1)}async function U(){const H=await fetch("/api/admin/user-columns",{method:"POST",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(b)});H.ok||alert(await H.text()),T(),rt()}function O(N){e(2,y=N),e(5,v={});for(const H in s[N])e(5,v[H]=s[N][H],v)}function L(){e(5,v={}),e(2,y=-1)}async function M(){const H=await fetch("/api/admin/user-columns",{method:"PUT",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([v])});H.ok||alert(await H.text()),L(),rt()}async function w(N){const H=s[N].idx,P=await fetch("/api/admin/user-columns",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([{idx:parseInt(H)}])});P.ok||alert(await P.text()),e(0,_=[]),rt()}async function F(){if(_.length==0){alert("Selected nothing");return}const N=[];for(let P=0;P<_.length;P++)N.push({idx:parseInt(_[P])});const q=await fetch("/api/admin/user-columns",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(N)});q.ok||alert(await q.text()),e(0,_=[]),rt()}const j=[[]],A=()=>{e(3,f=!0)};function z(){b["display-name"]=this.value,e(4,b),e(7,d)}function R(){b["column-code"]=this.value,e(4,b),e(7,d)}function J(){b["column-type"]=Ut(this),e(4,b),e(7,d)}function tt(){b["column-name"]=this.value,e(4,b),e(7,d)}function B(){v["display-name"]=this.value,e(5,v)}function S(){v["column-code"]=this.value,e(5,v)}function C(){v["column-name"]=this.value,e(5,v)}function c(){v["sort-order"]=gt(this.value),e(5,v)}function o(){_=St(j[0],this.__value,this.checked),e(0,_)}const m=N=>O(N),G=N=>w(N);return n.$$set=N=>{"data"in N&&e(16,a=N.data)},n.$$.update=()=>{n.$$.dirty[0]&65536&&e(1,s=a.columns),n.$$.dirty[0]&3&&e(6,i=_.length==s.length-kt&&_.length>0)},[_,s,y,f,b,v,i,d,I,T,U,O,L,M,w,F,a,A,z,R,J,tt,B,S,C,c,o,j,m,G]}class Xt extends Ot{constructor(t){super(),It(this,t,Lt,At,Dt,{data:16},null,[-1,-1])}}export{Xt as component,qt as universal};