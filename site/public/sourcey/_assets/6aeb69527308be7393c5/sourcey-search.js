async function e(e) {
	let t = new Worker(new URL(
		/* @vite-ignore */
		new URL("assets/search-worker-Dzc2-hN_.js", import.meta.url).href,
		"" + import.meta.url
	), { type: "module" }), n = 0, r = !1, i = /* @__PURE__ */ new Map();
	function a() {
		if (!r) {
			r = !0, t.terminate();
			for (let e of i.values()) e.reject(/* @__PURE__ */ Error("Search closed"));
			i.clear(), window.removeEventListener("pagehide", a), document.removeEventListener("astro:before-swap", a);
		}
	}
	t.addEventListener("message", (e) => {
		let t = e.data, n = i.get(t.id);
		n && (i.delete(t.id), t.error === void 0 ? n.resolve({
			rows: t.rows,
			more: t.more === !0
		}) : n.reject(Error(t.error)));
	}), t.addEventListener("error", a), t.addEventListener("messageerror", a), window.addEventListener("pagehide", a, { once: !0 }), document.addEventListener("astro:before-swap", a, { once: !0 });
	function o(e) {
		let a = ++n;
		return r ? {
			id: a,
			reply: Promise.reject(/* @__PURE__ */ Error("Search closed"))
		} : {
			id: a,
			reply: new Promise((n, r) => {
				i.set(a, {
					resolve: n,
					reject: r
				}), t.postMessage({
					...e,
					id: a
				});
			})
		};
	}
	let s = 0;
	try {
		let { rows: t } = await o({
			kind: "load",
			urls: e.map((e) => new URL(e, document.baseURI).href)
		}).reply;
		return {
			rows: t,
			query(e) {
				let { id: t, reply: n } = o({
					kind: "query",
					query: e
				});
				return s = t, n;
			},
			more() {
				return o({
					kind: "more",
					of: s
				}).reply;
			},
			dispose: a
		};
	} catch (e) {
		throw a(), e;
	}
}
export { e as loadSearch };
