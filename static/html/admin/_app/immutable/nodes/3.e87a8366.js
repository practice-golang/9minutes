import{s as Ue,n as He}from"../chunks/scheduler.e108d1fd.js";import{S as we,i as $e,g as a,s as u,m as Se,h as l,y as v,c as d,j as p,f as s,n as Ge,k as t,z as j,a as h,x as e,o as Me}from"../chunks/index.7e31c220.js";const Re=async({})=>({boardName:"boardName",boardCode:"boardCode",boardType:"boardType",boardTable:"boardTable",commentTable:"commentTable",grantRead:"grantRead",grantWrite:"grantWrite",grantComment:"grantComment",grantUpload:"grantUpload",link:"link",page:"page"}),ze=Object.freeze(Object.defineProperty({__proto__:null,load:Re},Symbol.toStringTag,{value:"Module"}));function qe(c){let k,kt="Admin / Board list",W,D,ee="Boards",Dt,z,ae='<label for="search">Search:</label> <input type="text" id="search" onkeyup="pressEnter()" placeholder="Search for..."/> <button type="button" onclick="search()">Search</button>',Nt,I,le="Create",Et,_,at,ne='<tr><th class="svelte-ukyb7v">Index</th> <th class="svelte-ukyb7v">Name</th> <th class="svelte-ukyb7v">Code</th> <th class="svelte-ukyb7v">Type</th> <th class="svelte-ukyb7v">Board table</th> <th class="svelte-ukyb7v">Comment table</th> <th class="svelte-ukyb7v">Grant read</th> <th class="svelte-ukyb7v">Grant write</th> <th class="svelte-ukyb7v">Grant comment</th> <th class="svelte-ukyb7v">Grant upload</th> <th class="svelte-ukyb7v">Control</th></tr>',Gt,R,re='<td class="svelte-ukyb7v"> </td> <td class="svelte-ukyb7v"><input type="text" name="board-name" value="" placeholder="Board name"/></td> <td class="svelte-ukyb7v"><input type="text" name="board-code" onchange="setTableNameByCode(this)" value="" placeholder="Board code"/></td> <td class="svelte-ukyb7v"><input type="text" name="board-type" onchange="restrictDatalist(this)" list="board-types" value="" autocomplete="off" placeholder="Board type"/></td> <td class="svelte-ukyb7v"><input type="text" name="board-table" value="" placeholder="Board table" disabled=""/></td> <td class="svelte-ukyb7v"><input type="text" name="comment-table" value="" placeholder="Comment table" disabled=""/></td> <td class="svelte-ukyb7v"><input type="text" name="grant-read" onchange="restrictDatalist(this)" list="grant-list" value="" autocomplete="off" placeholder="Grant read"/></td> <td class="svelte-ukyb7v"><input type="text" name="grant-write" onchange="restrictDatalist(this)" list="grant-list" value="" autocomplete="off" placeholder="Grant write"/></td> <td class="svelte-ukyb7v"><input type="text" name="grant-comment" onchange="restrictDatalist(this)" list="grant-list" value="" autocomplete="off" placeholder="Grant comment"/></td> <td class="svelte-ukyb7v"><input type="text" name="grant-upload" onchange="restrictDatalist(this)" list="grant-list" value="" autocomplete="off" placeholder="Grant upload"/></td> <td class="svelte-ukyb7v"><button type="button" onclick="closeAdd()">Cancel</button> <button type="button" onclick="addBoard()">Save</button></td>',Mt,P,Y,se='<td class="svelte-ukyb7v">idx</td> <td class="svelte-ukyb7v">board-name</td> <td class="svelte-ukyb7v">board-code</td> <td class="svelte-ukyb7v">board-type</td> <td class="svelte-ukyb7v">board-table</td> <td class="svelte-ukyb7v">comment-table</td> <td class="svelte-ukyb7v">grant-read</td> <td class="svelte-ukyb7v">grant-write</td> <td class="svelte-ukyb7v">grant-comment</td> <td class="svelte-ukyb7v">grant-upload</td> <td class="svelte-ukyb7v"><button type="button" lr-click="moveToListView($index)">View</button> <button type="button" lr-click="openEdit($index)">Edit</button> <button type="button" lr-click="deleteBoard($index)">Delete</button></td>',Ut,o,F,oe='<input type="hidden" name="idx" value="idx" placeholder="Index"/> <span>idx</span>',wt,lt,B,ct,$t,nt,O,bt,Rt,rt,m,ht,qt,st,N,_t,Vt,ot,E,mt,jt,ut,f,ft,Wt,dt,y,yt,zt,it,g,gt,Yt,vt,x,xt,Ft,J,ue='<button type="button" onclick="closeEdit()">Cancel</button> <button type="button" onclick="updateBoard()">Save</button>',It,C,L,de="Admin",A,ie="Manager",H,ve="Regular user",S,pe="Pending user",G,ce="Banned user",M,be="Guest",Pt,q,U,he="Board",w,_e="Gallery",Bt,K,b,Q,me="«",Jt,X,fe="<",Kt,pt,V,Ct=c[0].page+"",Ot,Lt,Qt,$,Tt=c[0].page+"",At,Ht,Xt,Z,ye=">",Zt,tt,ge="»";return{c(){k=a("h1"),k.textContent=kt,W=u(),D=a("h1"),D.textContent=ee,Dt=u(),z=a("div"),z.innerHTML=ae,Nt=u(),I=a("button"),I.textContent=le,Et=u(),_=a("table"),at=a("thead"),at.innerHTML=ne,Gt=u(),R=a("tr"),R.innerHTML=re,Mt=u(),P=a("tbody"),Y=a("tr"),Y.innerHTML=se,Ut=u(),o=a("tr"),F=a("td"),F.innerHTML=oe,wt=u(),lt=a("td"),B=a("input"),$t=u(),nt=a("td"),O=a("input"),Rt=u(),rt=a("td"),m=a("input"),qt=u(),st=a("td"),N=a("input"),Vt=u(),ot=a("td"),E=a("input"),jt=u(),ut=a("td"),f=a("input"),Wt=u(),dt=a("td"),y=a("input"),zt=u(),it=a("td"),g=a("input"),Yt=u(),vt=a("td"),x=a("input"),Ft=u(),J=a("td"),J.innerHTML=ue,It=u(),C=a("datalist"),L=a("option"),L.textContent=de,A=a("option"),A.textContent=ie,H=a("option"),H.textContent=ve,S=a("option"),S.textContent=pe,G=a("option"),G.textContent=ce,M=a("option"),M.textContent=be,Pt=u(),q=a("datalist"),U=a("option"),U.textContent=he,w=a("option"),w.textContent=_e,Bt=u(),K=a("div"),b=a("div"),Q=a("span"),Q.textContent=me,Jt=u(),X=a("span"),X.textContent=fe,Kt=u(),pt=a("b"),V=a("a"),Ot=Se(Ct),Qt=u(),$=a("a"),At=Se(Tt),Xt=u(),Z=a("span"),Z.textContent=ye,Zt=u(),tt=a("span"),tt.textContent=ge,this.h()},l(n){k=l(n,"H1",{"data-svelte-h":!0}),v(k)!=="svelte-197p21a"&&(k.textContent=kt),W=d(n),D=l(n,"H1",{"data-svelte-h":!0}),v(D)!=="svelte-1um8syh"&&(D.textContent=ee),Dt=d(n),z=l(n,"DIV",{"data-svelte-h":!0}),v(z)!=="svelte-1r6baws"&&(z.innerHTML=ae),Nt=d(n),I=l(n,"BUTTON",{type:!0,onclick:!0,"data-svelte-h":!0}),v(I)!=="svelte-m85p0h"&&(I.textContent=le),Et=d(n),_=l(n,"TABLE",{id:!0,class:!0});var r=p(_);at=l(r,"THEAD",{"data-svelte-h":!0}),v(at)!=="svelte-1c459nh"&&(at.innerHTML=ne),Gt=d(r),R=l(r,"TR",{id:!0,class:!0,"data-svelte-h":!0}),v(R)!=="svelte-gb1kvu"&&(R.innerHTML=re),Mt=d(r),P=l(r,"TBODY",{id:!0,"lr-loop":!0});var St=p(P);Y=l(St,"TR",{"lr-if":!0,"data-svelte-h":!0}),v(Y)!=="svelte-17tin2e"&&(Y.innerHTML=se),Ut=d(St),o=l(St,"TR",{"lr-if":!0});var i=p(o);F=l(i,"TD",{class:!0,"data-svelte-h":!0}),v(F)!=="svelte-1mq8n6o"&&(F.innerHTML=oe),wt=d(i),lt=l(i,"TD",{class:!0});var xe=p(lt);B=l(xe,"INPUT",{type:!0,name:!0,placeholder:!0}),xe.forEach(s),$t=d(i),nt=l(i,"TD",{class:!0});var Ce=p(nt);O=l(Ce,"INPUT",{type:!0,name:!0,placeholder:!0}),Ce.forEach(s),Rt=d(i),rt=l(i,"TD",{class:!0});var Te=p(rt);m=l(Te,"INPUT",{type:!0,name:!0,list:!0,autocomplete:!0,placeholder:!0,onchange:!0}),Te.forEach(s),qt=d(i),st=l(i,"TD",{class:!0});var ke=p(st);N=l(ke,"INPUT",{type:!0,name:!0,placeholder:!0}),ke.forEach(s),Vt=d(i),ot=l(i,"TD",{class:!0});var De=p(ot);E=l(De,"INPUT",{type:!0,name:!0,placeholder:!0}),De.forEach(s),jt=d(i),ut=l(i,"TD",{class:!0});var Ne=p(ut);f=l(Ne,"INPUT",{type:!0,name:!0,list:!0,autocomplete:!0,placeholder:!0,onchange:!0}),Ne.forEach(s),Wt=d(i),dt=l(i,"TD",{class:!0});var Ee=p(dt);y=l(Ee,"INPUT",{type:!0,name:!0,list:!0,autocomplete:!0,placeholder:!0,onchange:!0}),Ee.forEach(s),zt=d(i),it=l(i,"TD",{class:!0});var Ie=p(it);g=l(Ie,"INPUT",{type:!0,name:!0,list:!0,autocomplete:!0,placeholder:!0,onchange:!0}),Ie.forEach(s),Yt=d(i),vt=l(i,"TD",{class:!0});var Pe=p(vt);x=l(Pe,"INPUT",{type:!0,name:!0,list:!0,autocomplete:!0,placeholder:!0,onchange:!0}),Pe.forEach(s),Ft=d(i),J=l(i,"TD",{class:!0,"data-svelte-h":!0}),v(J)!=="svelte-4k5ggk"&&(J.innerHTML=ue),i.forEach(s),St.forEach(s),r.forEach(s),It=d(n),C=l(n,"DATALIST",{id:!0});var et=p(C);L=l(et,"OPTION",{"data-svelte-h":!0}),v(L)!=="svelte-17op260"&&(L.textContent=de),A=l(et,"OPTION",{"data-svelte-h":!0}),v(A)!=="svelte-1ob7gsw"&&(A.textContent=ie),H=l(et,"OPTION",{"data-svelte-h":!0}),v(H)!=="svelte-ielivd"&&(H.textContent=ve),S=l(et,"OPTION",{"data-svelte-h":!0}),v(S)!=="svelte-311eip"&&(S.textContent=pe),G=l(et,"OPTION",{"data-svelte-h":!0}),v(G)!=="svelte-1ot027p"&&(G.textContent=ce),M=l(et,"OPTION",{"data-svelte-h":!0}),v(M)!=="svelte-1j47qc6"&&(M.textContent=be),et.forEach(s),Pt=d(n),q=l(n,"DATALIST",{id:!0});var te=p(q);U=l(te,"OPTION",{"data-svelte-h":!0}),v(U)!=="svelte-ueeluy"&&(U.textContent=he),w=l(te,"OPTION",{"data-svelte-h":!0}),v(w)!=="svelte-79feqe"&&(w.textContent=_e),te.forEach(s),Bt=d(n),K=l(n,"DIV",{id:!0});var Be=p(K);b=l(Be,"DIV",{"lr-loop":!0});var T=p(b);Q=l(T,"SPAN",{"lr-if":!0,"data-svelte-h":!0}),v(Q)!=="svelte-2udg8q"&&(Q.textContent=me),Jt=d(T),X=l(T,"SPAN",{"lr-if":!0,"data-svelte-h":!0}),v(X)!=="svelte-1odbese"&&(X.textContent=fe),Kt=d(T),pt=l(T,"B",{"lr-if":!0});var Oe=p(pt);V=l(Oe,"A",{href:!0,rel:!0});var Le=p(V);Ot=Ge(Le,Ct),Le.forEach(s),Oe.forEach(s),Qt=d(T),$=l(T,"A",{"lr-if":!0,href:!0,rel:!0});var Ae=p($);At=Ge(Ae,Tt),Ae.forEach(s),Xt=d(T),Z=l(T,"SPAN",{"lr-if":!0,"data-svelte-h":!0}),v(Z)!=="svelte-8r8q7m"&&(Z.textContent=ye),Zt=d(T),tt=l(T,"SPAN",{"lr-if":!0,"data-svelte-h":!0}),v(tt)!=="svelte-3nt9fp"&&(tt.textContent=ge),T.forEach(s),Be.forEach(s),this.h()},h(){t(I,"type","button"),t(I,"onclick","openAdd()"),t(R,"id","add-board"),t(R,"class","svelte-ukyb7v"),t(Y,"lr-if","boardEditIndex != $index"),t(F,"class","svelte-ukyb7v"),t(B,"type","text"),t(B,"name","board-name"),B.value=ct=c[0].boardName,t(B,"placeholder","Board name"),t(lt,"class","svelte-ukyb7v"),t(O,"type","text"),t(O,"name","board-code"),O.value=bt=c[0].boardCode,t(O,"placeholder","Board code"),t(nt,"class","svelte-ukyb7v"),t(m,"type","text"),t(m,"name","board-type"),t(m,"list","board-types"),m.value=ht=c[0].boardType,t(m,"autocomplete","off"),t(m,"placeholder","Board type"),t(m,"onchange","restrictDatalist(this)"),t(rt,"class","svelte-ukyb7v"),t(N,"type","text"),t(N,"name","board-table"),N.value=_t=c[0].boardTable,t(N,"placeholder","Board table"),N.disabled=!0,t(st,"class","svelte-ukyb7v"),t(E,"type","text"),t(E,"name","comment-table"),E.value=mt=c[0].commentTable,t(E,"placeholder","Comment table"),E.disabled=!0,t(ot,"class","svelte-ukyb7v"),t(f,"type","text"),t(f,"name","grant-read"),t(f,"list","grant-list"),f.value=ft=c[0].grantRead,t(f,"autocomplete","off"),t(f,"placeholder","Grant read"),t(f,"onchange","restrictDatalist(this)"),t(ut,"class","svelte-ukyb7v"),t(y,"type","text"),t(y,"name","grant-write"),t(y,"list","grant-list"),y.value=yt=c[0].grantWrite,t(y,"autocomplete","off"),t(y,"placeholder","Grant write"),t(y,"onchange","restrictDatalist(this)"),t(dt,"class","svelte-ukyb7v"),t(g,"type","text"),t(g,"name","grant-comment"),t(g,"list","grant-list"),g.value=gt=c[0].grantComment,t(g,"autocomplete","off"),t(g,"placeholder","Grant comment"),t(g,"onchange","restrictDatalist(this)"),t(it,"class","svelte-ukyb7v"),t(x,"type","text"),t(x,"name","grant-upload"),t(x,"list","grant-list"),x.value=xt=c[0].grantUpload,t(x,"autocomplete","off"),t(x,"placeholder","Grant upload"),t(x,"onchange","restrictDatalist(this)"),t(vt,"class","svelte-ukyb7v"),t(J,"class","svelte-ukyb7v"),t(o,"lr-if","boardEditIndex == $index"),t(P,"id","boards-list-body"),t(P,"lr-loop","boardsList"),t(_,"id","boards-list-container"),t(_,"class","svelte-ukyb7v"),L.__value="admin",j(L,L.__value),A.__value="manager",j(A,A.__value),H.__value="regular_user",j(H,H.__value),S.__value="pending_user",j(S,S.__value),G.__value="banned_user",j(G,G.__value),M.__value="guest",j(M,M.__value),t(C,"id","grant-list"),U.__value="board",j(U,U.__value),w.__value="gallery",j(w,w.__value),t(q,"id","board-types"),t(Q,"lr-if","$index == 0 && pages[0].page > 1"),t(X,"lr-if","$index == 0 && pages[0].page > 1"),t(V,"href",Lt=c[0].link),t(V,"rel","external"),t(pt,"lr-if","page == boardsData['current-page']"),t($,"lr-if","page != boardsData['current-page']"),t($,"href",Ht=c[0].link),t($,"rel","external"),t(Z,"lr-if","page < boardsData['total-page'] && $index == (pages.length - 1)"),t(tt,"lr-if","page < boardsData['total-page'] && $index == (pages.length - 1)"),t(b,"lr-loop","pages"),t(K,"id","pages-container")},m(n,r){h(n,k,r),h(n,W,r),h(n,D,r),h(n,Dt,r),h(n,z,r),h(n,Nt,r),h(n,I,r),h(n,Et,r),h(n,_,r),e(_,at),e(_,Gt),e(_,R),e(_,Mt),e(_,P),e(P,Y),e(P,Ut),e(P,o),e(o,F),e(o,wt),e(o,lt),e(lt,B),e(o,$t),e(o,nt),e(nt,O),e(o,Rt),e(o,rt),e(rt,m),e(o,qt),e(o,st),e(st,N),e(o,Vt),e(o,ot),e(ot,E),e(o,jt),e(o,ut),e(ut,f),e(o,Wt),e(o,dt),e(dt,y),e(o,zt),e(o,it),e(it,g),e(o,Yt),e(o,vt),e(vt,x),e(o,Ft),e(o,J),h(n,It,r),h(n,C,r),e(C,L),e(C,A),e(C,H),e(C,S),e(C,G),e(C,M),h(n,Pt,r),h(n,q,r),e(q,U),e(q,w),h(n,Bt,r),h(n,K,r),e(K,b),e(b,Q),e(b,Jt),e(b,X),e(b,Kt),e(b,pt),e(pt,V),e(V,Ot),e(b,Qt),e(b,$),e($,At),e(b,Xt),e(b,Z),e(b,Zt),e(b,tt)},p(n,[r]){r&1&&ct!==(ct=n[0].boardName)&&B.value!==ct&&(B.value=ct),r&1&&bt!==(bt=n[0].boardCode)&&O.value!==bt&&(O.value=bt),r&1&&ht!==(ht=n[0].boardType)&&m.value!==ht&&(m.value=ht),r&1&&_t!==(_t=n[0].boardTable)&&N.value!==_t&&(N.value=_t),r&1&&mt!==(mt=n[0].commentTable)&&E.value!==mt&&(E.value=mt),r&1&&ft!==(ft=n[0].grantRead)&&f.value!==ft&&(f.value=ft),r&1&&yt!==(yt=n[0].grantWrite)&&y.value!==yt&&(y.value=yt),r&1&&gt!==(gt=n[0].grantComment)&&g.value!==gt&&(g.value=gt),r&1&&xt!==(xt=n[0].grantUpload)&&x.value!==xt&&(x.value=xt),r&1&&Ct!==(Ct=n[0].page+"")&&Me(Ot,Ct),r&1&&Lt!==(Lt=n[0].link)&&t(V,"href",Lt),r&1&&Tt!==(Tt=n[0].page+"")&&Me(At,Tt),r&1&&Ht!==(Ht=n[0].link)&&t($,"href",Ht)},i:He,o:He,d(n){n&&(s(k),s(W),s(D),s(Dt),s(z),s(Nt),s(I),s(Et),s(_),s(It),s(C),s(Pt),s(q),s(Bt),s(K))}}}function Ve(c,k,kt){let{data:W}=k;return c.$$set=D=>{"data"in D&&kt(0,W=D.data)},[W]}class Ye extends we{constructor(k){super(),$e(this,k,Ve,qe,Ue,{data:0})}}export{Ye as component,ze as universal};
