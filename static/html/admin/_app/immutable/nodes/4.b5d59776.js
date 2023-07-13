import{s as St,n as st,r as pt,f as jt}from"../chunks/scheduler.b95eede2.js";import{S as Bt,i as xt,g as u,s as E,h as r,y as F,c as g,j as N,f as c,k,A as G,a as M,z as n,B as K,C as Ut,F as kt,m as at,n as ut,e as Et,D as At,H as Pt,E as Ht,o as it,G as Lt}from"../chunks/index.e4e63e43.js";import{e as _t,i as dt}from"../chunks/navigation.732a7b1b.js";const Rt=async({fetch:l})=>{let t=[];const e=await l("/api/admin/user-columns",{method:"GET",headers:{"Content-Type":"application/json"},credentials:"include"});return e.ok&&(t=await e.json()),{columns:t}},Wt=Object.freeze(Object.defineProperty({__proto__:null,load:Rt},Symbol.toStringTag,{value:"Module"}));function gt(l,t,e){const i=l.slice();return i[30]=t[e],i[32]=e,i}function Dt(l,t,e){const i=l.slice();return i[33]=t[e][0],i[34]=t[e][1],i}function Nt(l){let t,e,i,s,a,m,b,_,y,p,h,v,x,O,U,j,Y,w,P,L="Cancel",A,B,Z="Save",J,$,q=_t(Object.entries(l[7])),I=[];for(let C=0;C<q.length;C+=1)I[C]=Ot(Dt(l,q,C));return{c(){t=u("tr"),e=u("td"),i=E(),s=u("td"),a=u("input"),m=E(),b=u("td"),_=u("input"),y=E(),p=u("td"),h=u("select");for(let C=0;C<I.length;C+=1)I[C].c();v=E(),x=u("td"),O=u("input"),U=E(),j=u("td"),Y=E(),w=u("td"),P=u("button"),P.textContent=L,A=E(),B=u("button"),B.textContent=Z,this.h()},l(C){t=r(C,"TR",{id:!0});var d=N(t);e=r(d,"TD",{}),N(e).forEach(c),i=g(d),s=r(d,"TD",{});var o=N(s);a=r(o,"INPUT",{type:!0,name:!0,placeholder:!0}),o.forEach(c),m=g(d),b=r(d,"TD",{});var f=N(b);_=r(f,"INPUT",{type:!0,name:!0,placeholder:!0}),f.forEach(c),y=g(d),p=r(d,"TD",{});var tt=N(p);h=r(tt,"SELECT",{name:!0});var D=N(h);for(let H=0;H<I.length;H+=1)I[H].l(D);D.forEach(c),tt.forEach(c),v=g(d),x=r(d,"TD",{});var S=N(x);O=r(S,"INPUT",{type:!0,name:!0,placeholder:!0}),S.forEach(c),U=g(d),j=r(d,"TD",{}),N(j).forEach(c),Y=g(d),w=r(d,"TD",{});var R=N(w);P=r(R,"BUTTON",{type:!0,"data-svelte-h":!0}),F(P)!=="svelte-1k872l3"&&(P.textContent=L),A=g(R),B=r(R,"BUTTON",{type:!0,"data-svelte-h":!0}),F(B)!=="svelte-gfswwj"&&(B.textContent=Z),R.forEach(c),d.forEach(c),this.h()},h(){k(a,"type","text"),k(a,"name","display-name"),k(a,"placeholder","Display name"),k(_,"type","text"),k(_,"name","column-code"),k(_,"placeholder","Column code"),k(h,"name","column-type"),l[4]["column-type"]===void 0&&jt(()=>l[20].call(h)),k(O,"type","text"),k(O,"name","column-name"),k(O,"placeholder","Column name"),k(P,"type","button"),k(B,"type","button"),k(t,"id","add-column")},m(C,d){M(C,t,d),n(t,e),n(t,i),n(t,s),n(s,a),G(a,l[4]["display-name"]),n(t,m),n(t,b),n(b,_),G(_,l[4]["column-code"]),n(t,y),n(t,p),n(p,h);for(let o=0;o<I.length;o+=1)I[o]&&I[o].m(h,null);kt(h,l[4]["column-type"],!0),n(t,v),n(t,x),n(x,O),G(O,l[4]["column-name"]),n(t,U),n(t,j),n(t,Y),n(t,w),n(w,P),n(w,A),n(w,B),J||($=[K(a,"input",l[18]),K(_,"input",l[19]),K(h,"change",l[20]),K(O,"input",l[21]),K(P,"click",l[9]),K(B,"click",l[10])],J=!0)},p(C,d){if(d[0]&144&&a.value!==C[4]["display-name"]&&G(a,C[4]["display-name"]),d[0]&144&&_.value!==C[4]["column-code"]&&G(_,C[4]["column-code"]),d[0]&128){q=_t(Object.entries(C[7]));let o;for(o=0;o<q.length;o+=1){const f=Dt(C,q,o);I[o]?I[o].p(f,d):(I[o]=Ot(f),I[o].c(),I[o].m(h,null))}for(;o<I.length;o+=1)I[o].d(1);I.length=q.length}d[0]&144&&kt(h,C[4]["column-type"]),d[0]&144&&O.value!==C[4]["column-name"]&&G(O,C[4]["column-name"])},d(C){C&&c(t),Ut(I,C),J=!1,pt($)}}}function Ot(l){let t,e=l[34]+"",i;return{c(){t=u("option"),i=at(e),this.h()},l(s){t=r(s,"OPTION",{});var a=N(t);i=ut(a,e),a.forEach(c),this.h()},h(){t.__value=l[33],G(t,t.__value)},m(s,a){M(s,t,a),n(t,i)},p:st,d(s){s&&c(t)}}}function qt(l){let t,e,i,s=l[30]["display-name"]+"",a,m,b,_=l[30]["column-code"]+"",y,p,h,v=l[7][l[30]["column-type"]]+"",x,O,U,j=l[30]["column-name"]+"",Y,w,P,L=l[30]["sort-order"]+"",A,B,Z;function J(o,f){return o[32]<=6?Xt:Jt}let q=J(l)(l);function I(o,f){return o[32]<=6?Yt:Gt}let d=I(l)(l);return{c(){t=u("tr"),q.c(),e=E(),i=u("td"),a=at(s),m=E(),b=u("td"),y=at(_),p=E(),h=u("td"),x=at(v),O=E(),U=u("td"),Y=at(j),w=E(),P=u("td"),A=at(L),B=E(),d.c(),Z=E()},l(o){t=r(o,"TR",{});var f=N(t);q.l(f),e=g(f),i=r(f,"TD",{});var tt=N(i);a=ut(tt,s),tt.forEach(c),m=g(f),b=r(f,"TD",{});var D=N(b);y=ut(D,_),D.forEach(c),p=g(f),h=r(f,"TD",{});var S=N(h);x=ut(S,v),S.forEach(c),O=g(f),U=r(f,"TD",{});var R=N(U);Y=ut(R,j),R.forEach(c),w=g(f),P=r(f,"TD",{});var H=N(P);A=ut(H,L),H.forEach(c),B=g(f),d.l(f),Z=g(f),f.forEach(c)},m(o,f){M(o,t,f),q.m(t,null),n(t,e),n(t,i),n(i,a),n(t,m),n(t,b),n(b,y),n(t,p),n(t,h),n(h,x),n(t,O),n(t,U),n(U,Y),n(t,w),n(t,P),n(P,A),n(t,B),d.m(t,null),n(t,Z)},p(o,f){q.p(o,f),f[0]&2&&s!==(s=o[30]["display-name"]+"")&&it(a,s),f[0]&2&&_!==(_=o[30]["column-code"]+"")&&it(y,_),f[0]&2&&v!==(v=o[7][o[30]["column-type"]]+"")&&it(x,v),f[0]&2&&j!==(j=o[30]["column-name"]+"")&&it(Y,j),f[0]&2&&L!==(L=o[30]["sort-order"]+"")&&it(A,L),d.p(o,f)},d(o){o&&c(t),q.d(),d.d()}}}function zt(l){let t,e,i,s,a,m,b,_,y,p,h=l[7][l[30]["column-type"]]+"",v,x,O,U,j,Y,w,P,L,A,B="Cancel",Z,J,$="Save",q,I,C;return{c(){t=u("tr"),e=u("td"),i=E(),s=u("td"),a=u("input"),m=E(),b=u("td"),_=u("input"),y=E(),p=u("td"),v=at(h),x=E(),O=u("td"),U=u("input"),j=E(),Y=u("td"),w=u("input"),P=E(),L=u("td"),A=u("button"),A.textContent=B,Z=E(),J=u("button"),J.textContent=$,q=E(),this.h()},l(d){t=r(d,"TR",{});var o=N(t);e=r(o,"TD",{}),N(e).forEach(c),i=g(o),s=r(o,"TD",{});var f=N(s);a=r(f,"INPUT",{type:!0,name:!0,placeholder:!0}),f.forEach(c),m=g(o),b=r(o,"TD",{});var tt=N(b);_=r(tt,"INPUT",{type:!0,name:!0,placeholder:!0}),tt.forEach(c),y=g(o),p=r(o,"TD",{});var D=N(p);v=ut(D,h),D.forEach(c),x=g(o),O=r(o,"TD",{});var S=N(O);U=r(S,"INPUT",{type:!0,name:!0,placeholder:!0}),S.forEach(c),j=g(o),Y=r(o,"TD",{});var R=N(Y);w=r(R,"INPUT",{type:!0,name:!0}),R.forEach(c),P=g(o),L=r(o,"TD",{});var H=N(L);A=r(H,"BUTTON",{type:!0,"data-svelte-h":!0}),F(A)!=="svelte-5czny3"&&(A.textContent=B),Z=g(H),J=r(H,"BUTTON",{type:!0,"data-svelte-h":!0}),F(J)!=="svelte-wlee15"&&(J.textContent=$),H.forEach(c),q=g(o),o.forEach(c),this.h()},h(){k(a,"type","text"),k(a,"name","display-name"),k(a,"placeholder","Display name"),k(_,"type","text"),k(_,"name","column-code"),k(_,"placeholder","Column code"),k(U,"type","text"),k(U,"name","column-name"),k(U,"placeholder","Column name"),k(w,"type","number"),k(w,"name","sort-order"),k(A,"type","button"),k(J,"type","button")},m(d,o){M(d,t,o),n(t,e),n(t,i),n(t,s),n(s,a),G(a,l[5]["display-name"]),n(t,m),n(t,b),n(b,_),G(_,l[5]["column-code"]),n(t,y),n(t,p),n(p,v),n(t,x),n(t,O),n(O,U),G(U,l[5]["column-name"]),n(t,j),n(t,Y),n(Y,w),G(w,l[5]["sort-order"]),n(t,P),n(t,L),n(L,A),n(L,Z),n(L,J),n(t,q),I||(C=[K(a,"input",l[22]),K(_,"input",l[23]),K(U,"input",l[24]),K(w,"input",l[25]),K(A,"click",l[12]),K(J,"click",l[13])],I=!0)},p(d,o){o[0]&32&&a.value!==d[5]["display-name"]&&G(a,d[5]["display-name"]),o[0]&32&&_.value!==d[5]["column-code"]&&G(_,d[5]["column-code"]),o[0]&2&&h!==(h=d[7][d[30]["column-type"]]+"")&&it(v,h),o[0]&32&&U.value!==d[5]["column-name"]&&G(U,d[5]["column-name"]),o[0]&32&&Pt(w.value)!==d[5]["sort-order"]&&G(w,d[5]["sort-order"])},d(d){d&&c(t),I=!1,pt(C)}}}function Jt(l){let t,e,i,s=!1,a,m,b;return a=Lt(l[27][0]),{c(){t=u("td"),e=u("input"),this.h()},l(_){t=r(_,"TD",{});var y=N(t);e=r(y,"INPUT",{type:!0}),y.forEach(c),this.h()},h(){k(e,"type","checkbox"),e.__value=i=l[30].idx,G(e,e.__value),a.p(e)},m(_,y){M(_,t,y),n(t,e),e.checked=~(l[0]||[]).indexOf(e.__value),m||(b=K(e,"change",l[26]),m=!0)},p(_,y){y[0]&2&&i!==(i=_[30].idx)&&(e.__value=i,G(e,e.__value),s=!0),(s||y[0]&3)&&(e.checked=~(_[0]||[]).indexOf(e.__value))},d(_){_&&c(t),a.r(),m=!1,b()}}}function Xt(l){let t;return{c(){t=u("td")},l(e){t=r(e,"TD",{}),N(t).forEach(c)},m(e,i){M(e,t,i)},p:st,d(e){e&&c(t)}}}function Gt(l){let t,e,i="Edit",s,a,m="Delete",b,_;function y(){return l[28](l[32])}function p(){return l[29](l[32])}return{c(){t=u("td"),e=u("button"),e.textContent=i,s=E(),a=u("button"),a.textContent=m,this.h()},l(h){t=r(h,"TD",{});var v=N(t);e=r(v,"BUTTON",{type:!0,"data-svelte-h":!0}),F(e)!=="svelte-ba0ha4"&&(e.textContent=i),s=g(v),a=r(v,"BUTTON",{type:!0,"data-svelte-h":!0}),F(a)!=="svelte-1glcjy2"&&(a.textContent=m),v.forEach(c),this.h()},h(){k(e,"type","button"),k(a,"type","button")},m(h,v){M(h,t,v),n(t,e),n(t,s),n(t,a),b||(_=[K(e,"click",y),K(a,"click",p)],b=!0)},p(h,v){l=h},d(h){h&&c(t),b=!1,pt(_)}}}function Yt(l){let t;return{c(){t=u("td")},l(e){t=r(e,"TD",{}),N(t).forEach(c)},m(e,i){M(e,t,i)},p:st,d(e){e&&c(t)}}}function It(l){let t;function e(a,m){return a[2]==a[32]?zt:qt}let i=e(l),s=i(l);return{c(){s.c(),t=Et()},l(a){s.l(a),t=Et()},m(a,m){s.m(a,m),M(a,t,m)},p(a,m){i===(i=e(a))&&s?s.p(a,m):(s.d(1),s=i(a),s&&(s.c(),s.m(t.parentNode,t)))},d(a){a&&c(t),s.d(a)}}}function Ft(l){let t,e="User columns",i,s,a="Add column",m,b,_="Delete selected columns",y,p,h,v,x,O,U,j,Y="Display name",w,P,L="Column code",A,B,Z="Column type",J,$,q="Column name",I,C,d="Sort Order",o,f,tt="Control",D,S,R,H,et,nt,mt="Text",lt,vt="Number Integer",ot,yt="Number Real",ft,bt,W=l[3]&&Nt(l),rt=_t(l[1]),Q=[];for(let T=0;T<rt.length;T+=1)Q[T]=It(gt(l,rt,T));return{c(){t=u("h1"),t.textContent=e,i=E(),s=u("button"),s.textContent=a,m=E(),b=u("button"),b.textContent=_,y=E(),p=u("table"),h=u("thead"),v=u("tr"),x=u("td"),O=u("input"),U=E(),j=u("th"),j.textContent=Y,w=E(),P=u("th"),P.textContent=L,A=E(),B=u("th"),B.textContent=Z,J=E(),$=u("th"),$.textContent=q,I=E(),C=u("th"),C.textContent=d,o=E(),f=u("th"),f.textContent=tt,D=E(),W&&W.c(),S=E(),R=u("tbody");for(let T=0;T<Q.length;T+=1)Q[T].c();H=E(),et=u("datalist"),nt=u("option"),nt.textContent=mt,lt=u("option"),lt.textContent=vt,ot=u("option"),ot.textContent=yt,this.h()},l(T){t=r(T,"H1",{"data-svelte-h":!0}),F(t)!=="svelte-1qzb11o"&&(t.textContent=e),i=g(T),s=r(T,"BUTTON",{type:!0,"data-svelte-h":!0}),F(s)!=="svelte-13ld9q5"&&(s.textContent=a),m=g(T),b=r(T,"BUTTON",{type:!0,"data-svelte-h":!0}),F(b)!=="svelte-1kxuovt"&&(b.textContent=_),y=g(T),p=r(T,"TABLE",{id:!0});var z=N(p);h=r(z,"THEAD",{});var X=N(h);v=r(X,"TR",{});var V=N(v);x=r(V,"TD",{});var Ct=N(x);O=r(Ct,"INPUT",{type:!0}),Ct.forEach(c),U=g(V),j=r(V,"TH",{"data-svelte-h":!0}),F(j)!=="svelte-5fbpcz"&&(j.textContent=Y),w=g(V),P=r(V,"TH",{"data-svelte-h":!0}),F(P)!=="svelte-5edxmv"&&(P.textContent=L),A=g(V),B=r(V,"TH",{"data-svelte-h":!0}),F(B)!=="svelte-rfc7ms"&&(B.textContent=Z),J=g(V),$=r(V,"TH",{"data-svelte-h":!0}),F($)!=="svelte-m0r1hp"&&($.textContent=q),I=g(V),C=r(V,"TH",{"data-svelte-h":!0}),F(C)!=="svelte-30jp5e"&&(C.textContent=d),o=g(V),f=r(V,"TH",{"data-svelte-h":!0}),F(f)!=="svelte-10rrrw5"&&(f.textContent=tt),V.forEach(c),X.forEach(c),D=g(z),W&&W.l(z),S=g(z),R=r(z,"TBODY",{id:!0});var Tt=N(R);for(let ht=0;ht<Q.length;ht+=1)Q[ht].l(Tt);Tt.forEach(c),z.forEach(c),H=g(T),et=r(T,"DATALIST",{id:!0});var ct=N(et);nt=r(ct,"OPTION",{"data-svelte-h":!0}),F(nt)!=="svelte-1mqqvia"&&(nt.textContent=mt),lt=r(ct,"OPTION",{"data-svelte-h":!0}),F(lt)!=="svelte-1mkb7wl"&&(lt.textContent=vt),ot=r(ct,"OPTION",{"data-svelte-h":!0}),F(ot)!=="svelte-10956bv"&&(ot.textContent=yt),ct.forEach(c),this.h()},h(){k(s,"type","button"),k(b,"type","button"),k(O,"type","checkbox"),O.checked=l[6],k(R,"id","column-list-body"),k(p,"id","column-list-container"),nt.__value="text",G(nt,nt.__value),lt.__value="number-integer",G(lt,lt.__value),ot.__value="number-real",G(ot,ot.__value),k(et,"id","column-types")},m(T,z){M(T,t,z),M(T,i,z),M(T,s,z),M(T,m,z),M(T,b,z),M(T,y,z),M(T,p,z),n(p,h),n(h,v),n(v,x),n(x,O),n(v,U),n(v,j),n(v,w),n(v,P),n(v,A),n(v,B),n(v,J),n(v,$),n(v,I),n(v,C),n(v,o),n(v,f),n(p,D),W&&W.m(p,null),n(p,S),n(p,R);for(let X=0;X<Q.length;X+=1)Q[X]&&Q[X].m(R,null);M(T,H,z),M(T,et,z),n(et,nt),n(et,lt),n(et,ot),ft||(bt=[K(s,"click",l[17]),K(b,"click",l[15]),K(O,"click",l[8])],ft=!0)},p(T,z){if(z[0]&64&&(O.checked=T[6]),T[3]?W?W.p(T,z):(W=Nt(T),W.c(),W.m(p,S)):W&&(W.d(1),W=null),z[0]&30887){rt=_t(T[1]);let X;for(X=0;X<rt.length;X+=1){const V=gt(T,rt,X);Q[X]?Q[X].p(V,z):(Q[X]=It(V),Q[X].c(),Q[X].m(R,null))}for(;X<Q.length;X+=1)Q[X].d(1);Q.length=rt.length}},i:st,o:st,d(T){T&&(c(t),c(i),c(s),c(m),c(b),c(y),c(p),c(H),c(et)),W&&W.d(),Ut(Q,T),ft=!1,pt(bt)}}}const wt=7;function Mt(l,t,e){let i,s,{data:a}=t,m=[],b=-1,_=!1,y={},p={};const h={text:"Text","number-integer":"Number Integer","number-real":"Number Real"};function v(){s?e(0,m=[]):e(0,m=i.filter(D=>D.idx>wt).map(D=>D.idx))}function x(){e(4,y={}),e(3,_=!1)}async function O(){const S=await fetch("/api/admin/user-columns",{method:"POST",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(y)});S.ok||alert(await S.text()),x(),dt()}function U(D){e(2,b=D),e(5,p={});for(const S in i[D])e(5,p[S]=i[D][S],p)}function j(){e(5,p={}),e(2,b=-1)}async function Y(){const S=await fetch("/api/admin/user-columns",{method:"PUT",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([p])});S.ok||alert(await S.text()),j(),dt()}async function w(D){const S=i[D].idx,H=await fetch("/api/admin/user-columns",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify([{idx:parseInt(S)}])});H.ok||alert(await H.text()),e(0,m=[]),dt()}async function P(){if(m.length==0){alert("Selected nothing");return}const D=[];for(let H=0;H<m.length;H++)D.push({idx:parseInt(m[H])});const R=await fetch("/api/admin/user-columns",{method:"DELETE",headers:{"Content-Type":"application/json"},credentials:"include",body:JSON.stringify(D)});R.ok||alert(await R.text()),e(0,m=[]),dt()}const L=[[]],A=()=>{e(3,_=!0)};function B(){y["display-name"]=this.value,e(4,y),e(7,h)}function Z(){y["column-code"]=this.value,e(4,y),e(7,h)}function J(){y["column-type"]=At(this),e(4,y),e(7,h)}function $(){y["column-name"]=this.value,e(4,y),e(7,h)}function q(){p["display-name"]=this.value,e(5,p)}function I(){p["column-code"]=this.value,e(5,p)}function C(){p["column-name"]=this.value,e(5,p)}function d(){p["sort-order"]=Pt(this.value),e(5,p)}function o(){m=Ht(L[0],this.__value,this.checked),e(0,m)}const f=D=>U(D),tt=D=>w(D);return l.$$set=D=>{"data"in D&&e(16,a=D.data)},l.$$.update=()=>{l.$$.dirty[0]&65536&&e(1,i=a.columns),l.$$.dirty[0]&3&&e(6,s=m.length==i.length-wt&&m.length>0)},[m,i,b,_,y,p,s,h,v,x,O,U,j,Y,w,P,a,A,B,Z,J,$,q,I,C,d,o,L,f,tt]}class Zt extends Bt{constructor(t){super(),xt(this,t,Mt,Ft,St,{data:16},null,[-1,-1])}}export{Zt as component,Wt as universal};
