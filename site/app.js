/* Rajaraja site: the link map, the page inspector and the live jq answers.
   Everything is computed in the browser from data/report.js, the unchanged report.json of a real crawl. */
(function () {
  "use strict";

  var REPORT = window.RAJARAJA_REPORT || [];
  var START = "https://learnwebscraping.dev/practice/ecommerce/";
  var SITE = "https://learnwebscraping.dev";
  var SVGNS = "http://www.w3.org/2000/svg";
  var REDUCED = window.matchMedia && window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  var map = document.getElementById("map");
  var inspector = document.getElementById("inspector");
  var jump = document.getElementById("jump");

  /* ---------- model ---------- */

  function uniq(list) {
    var seen = {}, out = [];
    list.forEach(function (x) { if (!seen[x]) { seen[x] = true; out.push(x); } });
    return out;
  }
  function hostOf(u) { try { return new URL(u).host; } catch (e) { return u; } }
  function el(name, attrs, parent) {
    var n = document.createElementNS(SVGNS, name);
    for (var k in attrs) n.setAttribute(k, attrs[k]);
    if (parent) parent.appendChild(n);
    return n;
  }

  var nodes = {};      // id (url) -> node
  var order = [];      // node ids in report order, externals last

  REPORT.forEach(function (p) {
    var u = p.url, role = "product";
    if (u === START) role = "home";
    else if (/\/categories\/$/.test(u)) role = "index";
    else if (/\/categories\/[^/]+\/$/.test(u)) role = "category";
    else if (/\/page\/\d+\/$/.test(u)) role = "listing";
    nodes[u] = { id: u, page: p, role: role, out: [], inn: [], external: false };
    order.push(u);
  });

  var edges = [];
  REPORT.forEach(function (p) {
    uniq(p.outgoing_links).forEach(function (t) {
      if (t === p.url) return;
      if (!nodes[t]) {
        if (t.indexOf(SITE + "/") === 0 || t === SITE) return; // internal page that was not crawled
        nodes[t] = { id: t, page: null, role: "external", out: [], inn: [], external: true };
        order.push(t);
      }
      var e = { from: p.url, to: t, ext: nodes[t].external };
      edges.push(e);
      nodes[p.url].out.push(e);
      nodes[t].inn.push(e);
    });
  });

  // clicks from home, following links in page order
  var discovery = [START];
  var depth = {}; depth[START] = 0;
  for (var qi = 0; qi < discovery.length; qi++) {
    var cur = nodes[discovery[qi]];
    if (!cur.page) continue;
    uniq(cur.page.outgoing_links).forEach(function (t) {
      if (nodes[t] && depth[t] === undefined) { depth[t] = depth[cur.id] + 1; discovery.push(t); }
    });
  }

  // which category each product belongs to: the category page it links to
  var categories = order.filter(function (id) { return nodes[id].role === "category"; });
  order.forEach(function (id) {
    var n = nodes[id];
    if (n.role !== "product") return;
    n.category = n.page.outgoing_links.filter(function (l) { return categories.indexOf(l) >= 0; })[0];
  });

  /* ---------- layout: a temple plan of square courts ---------- */

  var C = 500, R1 = 210, R2 = 378;
  function onSquare(deg, r) {
    var a = deg * Math.PI / 180, c = Math.cos(a), s = Math.sin(a);
    var m = Math.max(Math.abs(c), Math.abs(s));
    return [C + r * c / m, C + r * s / m];
  }

  var index = order.filter(function (id) { return nodes[id].role === "index"; })[0];
  var listing = order.filter(function (id) { return nodes[id].role === "listing"; })[0];
  var slots = [index].concat(categories.slice(0, 3), [listing], categories.slice(3, 6));
  var angleOf = {};
  slots.forEach(function (id, i) {
    if (!id) return;
    var deg = -90 + i * 45;
    angleOf[id] = deg;
    nodes[id].pos = onSquare(deg, R1);
    nodes[id].deg = deg;
  });
  nodes[START].pos = [C, C];

  categories.forEach(function (cat) {
    var members = order.filter(function (id) { return nodes[id].category === cat; });
    members.forEach(function (id, i) {
      var deg = angleOf[cat] + (i - (members.length - 1) / 2) * 7;
      nodes[id].pos = onSquare(deg, R2);
      nodes[id].deg = deg;
    });
  });

  var externals = order.filter(function (id) { return nodes[id].external; });
  externals.forEach(function (id, i) {
    nodes[id].pos = [C + (i - (externals.length - 1) / 2) * 140, 1022];
  });

  /* ---------- draw ---------- */

  function wall(r, gap, cls) {
    var x0 = C - r, x1 = C + r, y0 = C - r, y1 = C + r, g = gap / 2, d;
    if (!gap) d = "M" + x0 + " " + y0 + "H" + x1 + "V" + y1 + "H" + x0 + "Z";
    else d = [
      "M" + (C - g) + " " + y0 + "H" + x0 + "V" + (C - g),
      "M" + x0 + " " + (C + g) + "V" + y1 + "H" + (C - g),
      "M" + (C + g) + " " + y1 + "H" + x1 + "V" + (C + g),
      "M" + x1 + " " + (C - g) + "V" + y0 + "H" + (C + g)
    ].join(" ");
    el("path", { d: d, "class": "court " + (cls || "") }, map);
  }
  wall(118, 0, "court-inner");
  wall(292, 54);
  wall(460, 70);
  var outsideLabel = el("text", { x: C + 86, y: 1028, "class": "outside-label" }, map);
  outsideLabel.textContent = "outside the site, listed but never crawled";

  var edgeLayer = el("g", { "class": "edges" }, map);
  edges.forEach(function (e) {
    var a = nodes[e.from].pos, b = nodes[e.to].pos;
    e.line = el("line", {
      x1: a[0].toFixed(1), y1: a[1].toFixed(1), x2: b[0].toFixed(1), y2: b[1].toFixed(1),
      "class": "edge" + (e.ext ? " ext" : "")
    }, edgeLayer);
  });

  var LABELS = { home: "Home", index: "All categories", listing: "Page 2" };
  function labelFor(n) {
    if (n.external) return hostOf(n.id);
    if (LABELS[n.role]) return LABELS[n.role];
    if (n.role === "category") return n.page.heading;
    return null;
  }
  function nameFor(n) {
    if (n.external) return hostOf(n.id);
    if (n.role === "home") return "Home: " + n.page.heading;
    if (n.role === "listing") return "Page 2 of the catalog";
    return n.page.heading || n.id;
  }

  var nodeLayer = el("g", { "class": "nodes" }, map);
  order.forEach(function (id) {
    var n = nodes[id], p = n.pos;
    var r = n.role === "home" ? 22 : n.external ? 15 : n.role === "product" ? 8 : 12;
    var cls = "node " + (n.role === "home" ? "home" : n.external ? "external" : n.role === "product" ? "product" : "hub");
    var g = el("g", { "class": cls, tabindex: "0", role: "button", "aria-label": nameFor(n) }, nodeLayer);
    el("title", {}, g).textContent = nameFor(n);
    el("circle", { cx: p[0], cy: p[1], r: r + 7, "class": "halo" }, g);
    el("circle", { cx: p[0], cy: p[1], r: r + 10, fill: "transparent" }, g); // larger hit area
    el("circle", { cx: p[0], cy: p[1], r: r, "class": "dot" }, g);

    var text = labelFor(n);
    if (text) {
      var lx = p[0], ly = p[1], anchor = "middle";
      if (n.role === "home") { ly += 50; }
      else if (n.external) { lx -= 0; ly -= 26; }
      else {
        var a = n.deg * Math.PI / 180, cx = Math.cos(a), sy = Math.sin(a);
        lx += cx * 24; ly += sy * 24 + 6;
        if (cx > 0.3) anchor = "start"; else if (cx < -0.3) anchor = "end";
        if (Math.abs(cx) <= 0.3) ly += sy > 0 ? 14 : -4;
      }
      var t = el("text", { x: lx.toFixed(1), y: ly.toFixed(1), "text-anchor": anchor, "class": "label" }, g);
      t.textContent = text;
    }
    n.g = g;
    g.addEventListener("mouseenter", function () { if (!replaying) focusOn(id); });
    g.addEventListener("mouseleave", function () { if (!replaying) focusOn(selected); });
    g.addEventListener("focus", function () { if (!replaying) focusOn(id); });
    g.addEventListener("blur", function () { if (!replaying) focusOn(selected); });
    g.addEventListener("click", function () { select(id); });
    g.addEventListener("keydown", function (ev) {
      if (ev.key === "Enter" || ev.key === " ") { ev.preventDefault(); select(id); }
    });
  });
  // externals sit outside the wall label
  outsideLabel.setAttribute("x", (nodes[externals[0]] ? nodes[externals[0]].pos[0] + 28 : C + 86));

  /* ---------- focus and selection ---------- */

  var selected = "https://learnwebscraping.dev/practice/ecommerce/categories/arcane-staves/";
  if (!nodes[selected]) selected = START;
  var focused = null;

  function focusOn(id) {
    if (focused === id) return;
    focused = id;
    var n = nodes[id];
    map.classList.add("focus");
    edges.forEach(function (e) {
      e.line.classList.toggle("out", e.from === id);
      e.line.classList.toggle("in", e.to === id);
    });
    var outT = {}, inS = {};
    n.out.forEach(function (e) { outT[e.to] = true; });
    n.inn.forEach(function (e) { inS[e.from] = true; });
    var showExt = map.classList.contains("show-ext");
    order.forEach(function (o) {
      var m = nodes[o], isOut = !!outT[o] && (!m.external || showExt), isIn = !!inS[o];
      m.g.classList.toggle("is-out", isOut);
      m.g.classList.toggle("is-in", isIn);
      m.g.classList.toggle("dim", o !== id && !isOut && !isIn);
    });
  }

  function select(id) {
    if (replaying) stopReplay();
    var prev = nodes[selected];
    if (prev) prev.g.classList.remove("selected");
    selected = id;
    nodes[id].g.classList.add("selected");
    focused = null;
    focusOn(id);
    renderInspector(nodes[id]);
    if (jump.value !== id) jump.value = id;
  }

  /* ---------- inspector ---------- */

  function esc(s) { return String(s).replace(/[&<>"]/g, function (c) { return { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" }[c]; }); }
  function highlightJSON(obj) {
    return esc(JSON.stringify(obj, null, 2)).replace(/^(\s*)(&quot;[a-z_]+&quot;)(:)/gm, '$1<span class="k">$2</span>$3');
  }

  function renderInspector(n) {
    var html;
    if (n.external) {
      html =
        '<p class="kind">Another site</p>' +
        "<h2>" + esc(hostOf(n.id)) + "</h2>" +
        '<a class="url" href="' + esc(n.id) + '">' + esc(n.id) + "</a>" +
        "<dl><div><dt>Pages linking here</dt><dd class=\"in\">" + n.inn.length + "</dd></div>" +
        "<div><dt>Times fetched</dt><dd>0</dd></div></dl>" +
        '<p class="explain">Its hostname is not the start URL\'s hostname, so Rajaraja keeps it in every page\'s <code>outgoing_links</code> but never fetches it. Tick "Show links to other sites" to draw those links.</p>';
    } else {
      var p = n.page;
      var kinds = { home: "Home page, where the crawl started", index: "Hub page", category: "Category page", listing: "Catalog page", product: "Product page" };
      var internalOut = n.out.filter(function (e) { return !e.ext; }).length;
      html =
        '<p class="kind">' + kinds[n.role] + "</p>" +
        "<h2>" + esc(p.heading || "No h1 or h2 on this page") + "</h2>" +
        '<a class="url" href="' + esc(p.url) + '">' + esc(p.url) + "</a>" +
        "<dl>" +
        "<div><dt>Clicks from home</dt><dd>" + (depth[n.id] === undefined ? "n/a" : depth[n.id]) + "</dd></div>" +
        "<div><dt>Pages linking here</dt><dd class=\"in\">" + n.inn.length + "</dd></div>" +
        "<div><dt>Pages it links to</dt><dd class=\"out\">" + internalOut + "</dd></div>" +
        "<div><dt>Links on the page</dt><dd>" + p.outgoing_links.length + "</dd></div>" +
        "</dl>" +
        '<p class="entry-label">Its entry in report.json, exactly as written</p>' +
        "<pre>" + highlightJSON(p) + "</pre>";
    }
    inspector.innerHTML = html;
  }

  /* ---------- jump menu and external toggle ---------- */

  function addGroup(label, ids) {
    if (!ids.length) return;
    var og = document.createElement("optgroup");
    og.label = label;
    ids.forEach(function (id) {
      var o = document.createElement("option");
      o.value = id;
      o.textContent = nameFor(nodes[id]);
      og.appendChild(o);
    });
    jump.appendChild(og);
  }
  addGroup("Home and hub pages", [START, index, listing].concat(categories).filter(Boolean));
  categories.forEach(function (cat) {
    addGroup(nodes[cat].page.heading + " products", order.filter(function (id) { return nodes[id].category === cat; }));
  });
  addGroup("Other sites", externals);
  jump.addEventListener("change", function () { select(jump.value); });

  document.getElementById("external").addEventListener("change", function (ev) {
    map.classList.toggle("show-ext", ev.target.checked);
    var f = focused; focused = null; focusOn(f || selected);
  });

  /* ---------- replay the crawl ---------- */

  var replaying = false, timers = [];
  var seen = {};
  function stopReplay() {
    timers.forEach(clearTimeout); timers = [];
    replaying = false;
    map.classList.remove("replaying");
    order.forEach(function (id) { nodes[id].g.classList.remove("just-seen"); });
  }
  function reveal(id) {
    seen[id] = true;
    var g = nodes[id].g;
    g.classList.add("seen", "just-seen");
    timers.push(setTimeout(function () { g.classList.remove("just-seen"); }, 350));
    nodes[id].out.concat(nodes[id].inn).forEach(function (e) {
      if (seen[e.from] && seen[e.to]) e.line.classList.add("seen");
    });
  }
  function replay() {
    stopReplay();
    seen = {};
    order.forEach(function (id) { nodes[id].g.classList.remove("seen"); });
    edges.forEach(function (e) { e.line.classList.remove("seen"); });
    map.classList.remove("focus");
    order.forEach(function (id) { nodes[id].g.classList.remove("dim", "is-in", "is-out"); });
    edges.forEach(function (e) { e.line.classList.remove("in", "out"); });
    focused = null;
    replaying = true;
    map.classList.add("replaying");
    var seq = discovery.slice();
    var total = 2400;
    seq.forEach(function (id, i) {
      timers.push(setTimeout(function () { reveal(id); }, Math.round(i / (seq.length - 1) * total)));
    });
    timers.push(setTimeout(function () { stopReplay(); select(selected); }, total + 500));
  }
  document.getElementById("replay").addEventListener("click", replay);

  /* ---------- jq recipes, answered from the same data ---------- */

  function byCodepoint(a, b) { return a < b ? -1 : a > b ? 1 : 0; }
  var recipes = [
    {
      q: "How many pages were found?",
      cmd: "jq length report.json",
      run: function () { return String(REPORT.length); }
    },
    {
      q: "Which pages have no h1 or h2?",
      cmd: "jq -r '.[] | select(.heading == \"\") | .url' report.json",
      run: function () {
        var r = REPORT.filter(function (p) { return p.heading === ""; }).map(function (p) { return p.url; });
        return r.length ? r.join("\n") : "(no output: every page has a heading)";
      }
    },
    {
      q: "Which URLs are linked from the most pages?",
      cmd: "jq -r '[.[] | .outgoing_links | unique[]] | group_by(.) | map({url: .[0], pages: length}) | sort_by(-.pages) | .[:10][] | \"\\(.pages)\\t\\(.url)\"' report.json",
      run: function () {
        var count = {};
        REPORT.forEach(function (p) { uniq(p.outgoing_links).forEach(function (l) { count[l] = (count[l] || 0) + 1; }); });
        return Object.keys(count).sort(byCodepoint)
          .map(function (u) { return { url: u, pages: count[u] }; })
          .sort(function (a, b) { return b.pages - a.pages; })
          .slice(0, 10)
          .map(function (x) { return x.pages + "\t" + x.url; }).join("\n");
      }
    },
    {
      q: "Which other sites does it link to?",
      cmd: "jq -r --arg site \"" + SITE + "\" '[.[].outgoing_links[] | select(startswith($site) | not)] | unique[]' report.json",
      run: function () {
        var r = [];
        REPORT.forEach(function (p) { p.outgoing_links.forEach(function (l) { if (l.indexOf(SITE) !== 0) r.push(l); }); });
        return uniq(r).sort(byCodepoint).join("\n");
      }
    },
    {
      q: "Which images does the site use?",
      cmd: "jq -r '[.[].image_urls[]] | unique[]' report.json",
      run: function () {
        var r = [];
        REPORT.forEach(function (p) { r = r.concat(p.image_urls); });
        return uniq(r).sort(byCodepoint).join("\n");
      }
    },
    {
      q: "Export the link graph as CSV",
      cmd: "jq -r '.[] | .url as $from | .outgoing_links | unique[] | [$from, .] | @csv' report.json > edges.csv",
      run: function () {
        var rows = [];
        REPORT.forEach(function (p) {
          uniq(p.outgoing_links).sort(byCodepoint).forEach(function (l) { rows.push('"' + p.url + '","' + l + '"'); });
        });
        return rows.slice(0, 8).join("\n") + "\n...\n" + rows.length + " rows in total, ready for Gephi, Cytoscape or a spreadsheet";
      }
    }
  ];

  var tabs = document.getElementById("recipe-tabs");
  var cmdBox = document.getElementById("recipe-cmd");
  var result = document.getElementById("recipe-result");
  var tabButtons = recipes.map(function (r, i) {
    var b = document.createElement("button");
    b.type = "button";
    b.id = "recipe-tab-" + i;
    b.setAttribute("role", "tab");
    b.setAttribute("aria-controls", "recipe-panel");
    b.textContent = r.q;
    b.addEventListener("click", function () { showRecipe(i); });
    b.addEventListener("keydown", function (ev) {
      var k = ev.key, j = null;
      if (k === "ArrowDown" || k === "ArrowRight") j = (i + 1) % recipes.length;
      if (k === "ArrowUp" || k === "ArrowLeft") j = (i - 1 + recipes.length) % recipes.length;
      if (j !== null) { ev.preventDefault(); showRecipe(j); tabButtons[j].focus(); }
    });
    tabs.appendChild(b);
    return b;
  });
  function showRecipe(i) {
    tabButtons.forEach(function (b, j) {
      b.setAttribute("aria-selected", String(i === j));
      b.tabIndex = i === j ? 0 : -1;
    });
    document.getElementById("recipe-panel").setAttribute("aria-labelledby", "recipe-tab-" + i);
    cmdBox.querySelector("code").textContent = recipes[i].cmd;
    cmdBox.setAttribute("data-copy", recipes[i].cmd);
    result.textContent = recipes[i].run();
  }
  showRecipe(2);

  /* ---------- copy buttons and download ---------- */

  function copyText(text, btn) {
    function done() {
      btn.textContent = "Copied";
      btn.classList.add("done");
      setTimeout(function () { btn.textContent = "Copy"; btn.classList.remove("done"); }, 1600);
    }
    if (navigator.clipboard && window.isSecureContext) {
      navigator.clipboard.writeText(text).then(done, function () { fallback(text); done(); });
    } else { fallback(text); done(); }
  }
  function fallback(text) {
    var ta = document.createElement("textarea");
    ta.value = text; ta.setAttribute("readonly", ""); ta.style.position = "fixed"; ta.style.opacity = "0";
    document.body.appendChild(ta); ta.select();
    try { document.execCommand("copy"); } catch (e) { /* nothing else to try */ }
    document.body.removeChild(ta);
  }
  document.addEventListener("click", function (ev) {
    var btn = ev.target.closest && ev.target.closest(".copy");
    if (!btn) return;
    var box = btn.closest(".cmd");
    var text = box.getAttribute("data-copy");
    if (!text) text = box.querySelector("code").textContent.replace(/^\$ /gm, "");
    copyText(text, btn);
  });

  document.getElementById("download").addEventListener("click", function () {
    var blob = new Blob([JSON.stringify(REPORT, null, 2) + "\n"], { type: "application/json" });
    var a = document.createElement("a");
    a.href = URL.createObjectURL(blob);
    a.download = "report.json";
    document.body.appendChild(a); a.click(); a.remove();
    setTimeout(function () { URL.revokeObjectURL(a.href); }, 1000);
  });

  /* ---------- start ---------- */

  select(selected);
  if (!REDUCED) replay();
})();
