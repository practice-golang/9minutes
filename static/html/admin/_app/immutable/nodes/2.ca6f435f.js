import{s as g,n as f}from"../chunks/scheduler.e108d1fd.js";import{S as A,i as k,g as h,s as C,h as m,y as x,c as _,k as a,a as o,f as r}from"../chunks/index.c5af3f48.js";function w(v){let n,d="Admin",i,e,c="Login",u,l,p="Logout";return{c(){n=h("h1"),n.textContent=d,i=C(),e=h("a"),e.textContent=c,u=C(),l=h("a"),l.textContent=p,this.h()},l(t){n=m(t,"H1",{"data-svelte-h":!0}),x(n)!=="svelte-1d5mfrr"&&(n.textContent=d),i=_(t),e=m(t,"A",{id:!0,href:!0,rel:!0,"data-svelte-h":!0}),x(e)!=="svelte-whoz57"&&(e.textContent=c),u=_(t),l=m(t,"A",{id:!0,href:!0,rel:!0,"data-svelte-h":!0}),x(l)!=="svelte-1kvwuds"&&(l.textContent=p),this.h()},h(){a(e,"id","login"),a(e,"href","/login"),a(e,"rel","external"),a(l,"id","logout"),a(l,"href","/logout"),a(l,"rel","external")},m(t,s){o(t,n,s),o(t,i,s),o(t,e,s),o(t,u,s),o(t,l,s)},p:f,i:f,o:f,d(t){t&&(r(n),r(i),r(e),r(u),r(l))}}}class S extends A{constructor(n){super(),k(this,n,null,w,g,{})}}export{S as component};