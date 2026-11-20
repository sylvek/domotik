'use strict';

/* ---------- helpers ---------- */

const SVG = 'http://www.w3.org/2000/svg';
const SERIES = ['--s1', '--s2', '--s3', '--s4', '--s5', '--s6', '--s7', '--s8'];
const color = (i) => `var(${SERIES[i % SERIES.length]})`;

function h(tag, attrs, ...children) {
  const el = document.createElement(tag);
  for (const [k, v] of Object.entries(attrs || {})) {
    if (k === 'style') el.style.cssText = v; else el.setAttribute(k, v);
  }
  for (const child of children.flat()) {
    if (child != null) el.append(child);
  }
  return el;
}

function s(tag, attrs, text) {
  const el = document.createElementNS(SVG, tag);
  for (const [k, v] of Object.entries(attrs || {})) el.setAttribute(k, v);
  if (text != null) el.textContent = text;
  return el;
}

const round = (v, digits) => v.toLocaleString('fr-FR', { minimumFractionDigits: digits, maximumFractionDigits: digits });

const fmt = {
  temp: (v) => round(v, 1) + ' °C',
  watt: (v) => Math.abs(v) >= 1000 ? round(v / 1000, 2) + ' kW' : round(v, 0) + ' W',
  wattHour: (v) => Math.abs(v) >= 1e6 ? round(v / 1e6, 2) + ' MWh' : Math.abs(v) >= 1000 ? round(v / 1000, 1) + ' kWh' : round(v, 0) + ' Wh',
  euro: (v) => round(v, 2) + ' €',
  liter: (v) => round(v, 0) + ' L',
  percent: (v) => round(v, 1) + ' %',
  count: (v) => round(v, 0),
};

// axis variants: same units, less digits
const axis = {
  temp: (v) => round(v, 0) + ' °C',
  watt: (v) => Math.abs(v) >= 1000 ? round(v / 1000, 1) + ' kW' : round(v, 0) + ' W',
  wattHour: (v) => Math.abs(v) >= 1e6 ? round(v / 1e6, 1) + ' MWh' : Math.abs(v) >= 1000 ? round(v / 1000, 0) + ' kWh' : round(v, 0) + ' Wh',
  euro: (v) => round(v, 0) + ' €',
  liter: (v) => round(v, 0) + ' L',
  percent: (v) => round(v, 0) + ' %',
  count: (v) => round(v, 0),
};

const pad = (n) => String(n).padStart(2, '0');

function niceTicks(min, max, count) {
  const span = max - min || 1;
  const raw = span / count;
  const mag = Math.pow(10, Math.floor(Math.log10(raw)));
  const step = [1, 2, 2.5, 5, 10].map((m) => m * mag).find((m) => m >= raw);
  const ticks = [];
  for (let v = Math.ceil(min / step) * step; v <= max + step * 1e-9; v += step) ticks.push(Math.abs(v) < step * 1e-9 ? 0 : v);
  return ticks;
}

// time ticks & labels. `utc` is used for daily data, stamped at midnight UTC.
function timeTicks(from, to, width, utc) {
  const span = to - from;
  const max = Math.max(2, Math.floor(width / 70));
  const get = (d, what) => d[(utc ? 'getUTC' : 'get') + what]();
  const ticks = [];

  if (span <= 3 * 86400) {
    const step = [1, 2, 3, 6, 12].find((x) => span / (x * 3600) <= max) || 24;
    const d = new Date(from * 1000);
    d.setMinutes(0, 0, 0);
    for (; d / 1000 <= to; d.setHours(d.getHours() + 1)) {
      if (d / 1000 >= from && d.getHours() % step === 0) ticks.push({ t: d / 1000, label: `${pad(d.getHours())}:${pad(d.getMinutes())}` });
    }
  } else if (span <= 200 * 86400) {
    const step = [1, 2, 7, 14].find((x) => span / (x * 86400) <= max) || 30;
    const first = Math.ceil(from / 86400);
    for (let day = first; day * 86400 <= to; day++) {
      if ((day - first) % step === 0) {
        const d = new Date(day * 86400000);
        ticks.push({ t: day * 86400, label: `${pad(d.getUTCDate())}/${pad(d.getUTCMonth() + 1)}` });
      }
    }
  } else {
    const months = span / (30.44 * 86400);
    const step = [1, 2, 3, 6, 12].find((x) => months / x <= max) || 24;
    const d = new Date(from * 1000);
    let y = get(d, 'FullYear'), m = get(d, 'Month') + 1;
    for (;;) {
      const t = Date.UTC(y, m, 1) / 1000;
      if (t > to) break;
      if (m % step === 0) ticks.push({ t, label: step >= 12 ? String(y) : `${y}-${pad(m + 1)}` });
      if (++m > 11) { m = 0; y++; }
    }
  }
  return ticks;
}

function formatWhen(t, utc, withTime) {
  const d = new Date(t * 1000);
  if (utc) return d.toLocaleDateString('fr-FR', { timeZone: 'UTC', weekday: 'short', day: 'numeric', month: 'short', year: 'numeric' });
  return withTime ? `${pad(d.getHours())}:${pad(d.getMinutes())}` : d.toLocaleDateString('fr-FR');
}

/* ---------- seasons & trends ---------- */

const SEASONS = ['winter', 'spring', 'summer', 'autumn'];
const seasonColor = (season) => `var(--season-${season})`;

// one series by season, so that every season has its color and its statistics
function bySeason(serie) {
  return SEASONS.map((season) => ({
    label: season,
    color: seasonColor(season),
    points: serie.points.filter((p) => p.s === season),
  }));
}

function seasonLegend() {
  return h('div', { class: 'keys' }, SEASONS.map((season) => h('span', {},
    h('span', { class: 'dot', style: `background:${seasonColor(season)}` }), season)));
}

// direction of a trend: an arrow and a signed percentage. `upIs` tells if a rise is good, bad or neutral.
function trendChange(trend, upIs = 'bad') {
  if (!trend || trend.change == null) return null;
  const percent = trend.change * 100;
  const flat = Math.abs(percent) < 0.05;
  const up = percent > 0;
  const status = flat || upIs === 'neutral' ? 'text-muted' : (up === (upIs === 'good') ? 'good' : 'serious');
  return {
    arrow: h('span', { class: 'arrow', style: `color:var(--${status})` }, flat ? '=' : up ? '▲' : '▼'),
    text: `${up && !flat ? '+' : ''}${round(percent, 1)} %`,
  };
}

/* ---------- tooltip ---------- */

const tooltip = document.getElementById('tooltip');

function showTooltip(event, when, rows) {
  tooltip.replaceChildren(
    h('div', { class: 'when' }, when),
    ...rows.map((r) => h('div', { class: 'row' },
      r.color ? h('span', { class: 'dot', style: `background:${r.color}` }) : null,
      h('span', {}, r.label),
      h('b', {}, r.value))));
  tooltip.hidden = false;
  const box = tooltip.getBoundingClientRect();
  let x = event.clientX + 14, y = event.clientY + 14;
  if (x + box.width > innerWidth - 8) x = event.clientX - box.width - 14;
  if (y + box.height > innerHeight - 8) y = event.clientY - box.height - 14;
  tooltip.style.left = Math.max(8, x) + 'px';
  tooltip.style.top = Math.max(8, y) + 'px';
}

function hideTooltip() { tooltip.hidden = true; }

/* ---------- charts ---------- */

const charts = []; // redraw functions, called on resize

function panel(title, subtitle) {
  const body = h('div', {});
  const el = h('section', { class: 'panel' }, h('h2', {}, title), subtitle ? h('p', { class: 'subtitle' }, subtitle) : null, body);
  return { el, body };
}

function stats(points) {
  let min = Infinity, max = -Infinity, sum = 0;
  for (const p of points) { min = Math.min(min, p.v); max = Math.max(max, p.v); sum += p.v; }
  return { min, max, mean: sum / points.length, last: points[points.length - 1].v };
}

function legend(series, unit) {
  const head = h('tr', {}, h('th', {}, ''), ['max', 'min', 'mean', 'last'].map((c) => h('th', {}, c)));
  const rows = series.map((serie) => {
    const st = stats(serie.points);
    return h('tr', {},
      h('td', {}, h('span', { class: 'key', style: `background:${serie.color}` }), serie.label),
      [st.max, st.min, st.mean, st.last].map((v) => h('td', {}, fmt[unit](v))));
  });
  return h('table', { class: 'legend' }, h('thead', {}, head), h('tbody', {}, rows));
}

/**
 * Chart over time.
 * options: { series: [{label, points: [{t, v}]}], from, to, unit, kind: line|area|dots|bars,
 *            kinds: per series kind, utc, zero: y axis starts at zero, max: y axis maximum, height }
 */
function timeChart(title, options, subtitle) {
  const { el, body } = panel(title, subtitle);
  const series = options.series
    .map((serie, i) => ({ ...serie, color: serie.color || color(i), kind: (options.kinds || [])[i] || options.kind || 'line' }))
    .filter((serie) => serie.points && serie.points.length > 0);

  if (series.length === 0) {
    body.append(h('div', { class: 'empty' }, 'no data'));
    return el;
  }

  const height = options.height || 240;
  const margin = { top: 8, right: 12, bottom: 22, left: 58 };
  const svg = s('svg', { height, role: 'img', 'aria-label': title });
  body.append(svg, legend(series, options.unit));

  const all = series.flatMap((serie) => serie.points);
  const bars = series.some((serie) => serie.kind === 'bars');
  let lo = Math.min(...all.map((p) => p.v)), hi = Math.max(...all.map((p) => p.v));
  if (options.zero || bars) lo = Math.min(0, lo);
  if (options.max != null) hi = Math.max(hi, options.max);
  if (hi === lo) hi = lo + 1;
  const yTicks = niceTicks(lo, hi, Math.max(2, Math.floor(height / 50)));
  const step = yTicks.length > 1 ? yTicks[1] - yTicks[0] : 1;
  lo = Math.min(lo, yTicks[0]);
  if (!(options.zero || bars) || lo < 0) lo = Math.floor(lo / step) * step;
  hi = Math.ceil(hi / step - 1e-9) * step;

  // a daily bar covers its whole day
  const from = options.from, to = options.to;

  function draw() {
    const width = svg.clientWidth || 600;
    svg.setAttribute('viewBox', `0 0 ${width} ${height}`);
    svg.replaceChildren();
    const w = width - margin.left - margin.right, hh = height - margin.top - margin.bottom;
    const x = (t) => margin.left + (t - from) / (to - from) * w;
    const y = (v) => margin.top + (1 - (v - lo) / (hi - lo)) * hh;

    for (let v = lo; v <= hi + step * 1e-9; v += step) {
      svg.append(s('line', { class: 'grid', x1: margin.left, x2: width - margin.right, y1: y(v), y2: y(v) }));
      svg.append(s('text', { x: margin.left - 8, y: y(v) + 4, 'text-anchor': 'end' }, axis[options.unit](v)));
    }
    for (const tick of timeTicks(from, to, w, options.utc)) {
      svg.append(s('text', { x: x(tick.t), y: height - 6, 'text-anchor': 'middle' }, tick.label));
    }

    const days = (to - from) / 86400;
    for (const serie of series) {
      const pts = serie.points;
      if (serie.kind === 'bars') {
        const bw = Math.max(1, Math.min(24, w / days - 2));
        let d = '';
        for (const p of pts) {
          const top = y(p.v), base = y(0), r = Math.min(4, bw / 2, Math.abs(base - top));
          const left = x(p.t) - bw / 2;
          d += `M${left},${base}V${top + r}q0,${-r} ${r},${-r}h${bw - 2 * r}q${r},0 ${r},${r}V${base}Z`;
        }
        svg.append(s('path', { d, fill: serie.color }));
      } else if (serie.kind === 'dots') {
        let d = '';
        for (const p of pts) d += `M${x(p.t)},${y(p.v)}h0`;
        svg.append(s('path', { d, stroke: serie.color, 'stroke-width': 4, 'stroke-linecap': 'round', fill: 'none' }));
      } else {
        // a line is broken when measures are missing for a while
        const gap = median(pts) * 6;
        let d = '', area = '', start = null;
        pts.forEach((p, i) => {
          const broken = i === 0 || p.t - pts[i - 1].t > gap;
          if (broken && start != null) area += `L${x(pts[i - 1].t)},${y(lo)}L${x(start)},${y(lo)}Z`;
          if (broken) start = p.t;
          const cmd = `${broken ? 'M' : 'L'}${x(p.t).toFixed(1)},${y(p.v).toFixed(1)}`;
          d += cmd; area += cmd;
        });
        area += `L${x(pts[pts.length - 1].t)},${y(lo)}L${x(start)},${y(lo)}Z`;
        if (serie.kind === 'area') svg.append(s('path', { d: area, fill: serie.color, 'fill-opacity': 0.12 }));
        svg.append(s('path', { d, stroke: serie.color, 'stroke-width': 2, 'stroke-linejoin': 'round', 'stroke-linecap': 'round', fill: 'none' }));
      }
    }

    // hover layer
    const cross = s('line', { class: 'cross', y1: margin.top, y2: margin.top + hh, visibility: 'hidden' });
    const markers = series.map((serie) => s('circle', { r: 4, fill: serie.color, stroke: 'var(--surface)', 'stroke-width': 2, visibility: 'hidden' }));
    const hit = s('rect', { x: margin.left, y: margin.top, width: w, height: hh, fill: 'transparent' });
    svg.append(cross, ...markers, hit);

    hit.addEventListener('pointermove', (event) => {
      const box = svg.getBoundingClientRect();
      const t = from + (event.clientX - box.left - margin.left) / w * (to - from);
      // the reference time is the nearest point, every series is then read at that time
      let ref = null;
      for (const serie of series) {
        const p = nearest(serie.points, t);
        if (ref == null || Math.abs(p.t - t) < Math.abs(ref - t)) ref = p.t;
      }
      const tolerance = Math.max((to - from) / w * 12, median(series[0].points) * 1.5);
      const rows = [];
      series.forEach((serie, i) => {
        const p = nearest(serie.points, ref);
        if (Math.abs(p.t - ref) > tolerance) { markers[i].setAttribute('visibility', 'hidden'); return; }
        markers[i].setAttribute('cx', x(p.t));
        markers[i].setAttribute('cy', y(p.v));
        markers[i].setAttribute('visibility', 'visible');
        rows.push({ color: serie.color, label: serie.label, value: fmt[options.unit](p.v) });
      });
      cross.setAttribute('x1', x(ref));
      cross.setAttribute('x2', x(ref));
      cross.setAttribute('visibility', 'visible');
      showTooltip(event, formatWhen(ref, options.utc, true), rows);
    });
    hit.addEventListener('pointerleave', () => {
      cross.setAttribute('visibility', 'hidden');
      markers.forEach((m) => m.setAttribute('visibility', 'hidden'));
      hideTooltip();
    });
  }

  charts.push(draw);
  return el;
}

function median(points) {
  if (points.length < 2) return 3600;
  const deltas = [];
  for (let i = 1; i < points.length; i++) deltas.push(points[i].t - points[i - 1].t);
  deltas.sort((a, b) => a - b);
  return deltas[deltas.length >> 1] || 1;
}

function nearest(points, t) {
  let a = 0, b = points.length - 1;
  while (b - a > 1) {
    const m = (a + b) >> 1;
    if (points[m].t < t) a = m; else b = m;
  }
  return Math.abs(points[a].t - t) <= Math.abs(points[b].t - t) ? points[a] : points[b];
}

/** Bar chart by category. options: { buckets: [{label, value}], unit, height } */
function barChart(title, options, subtitle) {
  const { el, body } = panel(title, subtitle);
  const buckets = options.buckets || [];
  if (buckets.length === 0) {
    body.append(h('div', { class: 'empty' }, 'no data'));
    return el;
  }

  const height = options.height || 200;
  const margin = { top: 18, right: 12, bottom: 22, left: 58 };
  const svg = s('svg', { height, role: 'img', 'aria-label': title });
  body.append(svg);
  if (buckets.some((b) => b.season)) body.append(seasonLegend());

  const ticks = niceTicks(0, Math.max(1, ...buckets.map((b) => b.value)), Math.max(2, Math.floor(height / 50)));
  const step = ticks[1] - ticks[0];
  const hi = Math.ceil(Math.max(1, ...buckets.map((b) => b.value)) / step - 1e-9) * step;

  function draw() {
    const width = svg.clientWidth || 600;
    svg.setAttribute('viewBox', `0 0 ${width} ${height}`);
    svg.replaceChildren();
    const w = width - margin.left - margin.right, hh = height - margin.top - margin.bottom;
    const band = w / buckets.length;
    const y = (v) => margin.top + (1 - v / hi) * hh;

    for (let v = 0; v <= hi + step * 1e-9; v += step) {
      svg.append(s('line', { class: 'grid', x1: margin.left, x2: width - margin.right, y1: y(v), y2: y(v) }));
      svg.append(s('text', { x: margin.left - 8, y: y(v) + 4, 'text-anchor': 'end' }, axis[options.unit](v)));
    }

    const bw = Math.max(2, Math.min(24, band - 2));
    const every = Math.ceil(buckets.length / Math.max(1, Math.floor(w / 64)));
    const direct = band >= 56; // values are written on the caps when there is room
    buckets.forEach((b, i) => {
      const cx = margin.left + band * (i + 0.5);
      const top = y(b.value), base = y(0), r = Math.min(4, bw / 2, base - top);
      const left = cx - bw / 2;
      const fill = b.season ? seasonColor(b.season) : color(0);
      svg.append(s('path', { d: `M${left},${base}V${top + r}q0,${-r} ${r},${-r}h${bw - 2 * r}q${r},0 ${r},${r}V${base}Z`, fill }));
      if (direct) svg.append(s('text', { x: cx, y: top - 5, 'text-anchor': 'middle', style: 'fill:var(--text-secondary)' }, fmt[options.unit](b.value)));
      if (i % every === 0) svg.append(s('text', { x: cx, y: height - 6, 'text-anchor': 'middle' }, b.label));
      const hit = s('rect', { x: cx - band / 2, y: margin.top, width: band, height: hh, fill: 'transparent' });
      hit.addEventListener('pointermove', (event) => {
        const rows = [{ color: fill, label: b.season || title, value: fmt[options.unit](b.value) }];
        const change = trendChange(b.trend);
        if (change) {
          rows.push({ label: 'a year before', value: fmt[options.unit](b.trend.previous) });
          rows.push({ label: h('span', {}, change.arrow, ' trend'), value: change.text });
        }
        showTooltip(event, b.label, rows);
      });
      hit.addEventListener('pointerleave', hideTooltip);
      svg.append(hit);
    });
  }

  charts.push(draw);
  return el;
}

/* ---------- tiles ---------- */

// thresholds: [[from, status]…] sorted, the last matching one gives the status of the value.
function tile(label, value, unit, options = {}) {
  const parts = [h('div', { class: 'label' }, label)];
  parts.push(h('div', { class: 'value' }, value, unit ? h('small', {}, unit) : null));
  if (options.max != null) {
    const [, status] = options.thresholds.filter(([from]) => options.raw >= from).pop();
    const ratio = Math.max(0, Math.min(1, options.raw / options.max));
    parts.push(h('div', { class: 'meter' }, h('i', { style: `width:${ratio * 100}%;background:var(--${status})` })));
    parts.push(h('div', { class: 'hint' }, h('span', { class: 'dot', style: `background:var(--${status})` }), options.hint(status)));
  } else if (options.hint) {
    parts.push(h('div', { class: 'hint' }, options.hint));
  }
  return h('div', { class: 'tile' }, parts);
}

// a value compared to the same period, a year before
function trendTile(label, trend, unit, upIs) {
  if (!trend) return tile(label, '—', null, { hint: 'no data' });
  const [value, suffix] = fmt[unit](trend.value).split(/ (?=\S+$)/);
  const change = trendChange(trend, upIs);
  // a percentage of a temperature means nothing: the difference is displayed in degrees
  if (change && unit === 'temp') {
    const delta = trend.value - trend.previous;
    change.text = `${delta > 0 ? '+' : ''}${fmt.temp(delta)}`;
  }
  const hint = change
    ? h('span', {}, change.arrow, ` ${change.text} vs a year before (${fmt[unit](trend.previous)})`)
    : 'no comparable data a year before';
  return tile(label, value, suffix, { hint });
}

function readings(title, list) {
  return h('div', { class: 'tile' },
    h('p', { class: 'group-title' }, title),
    h('div', { class: 'group' }, list.map((r) => h('div', { class: 'reading' },
      h('div', { class: 'label' }, r.label),
      h('div', { class: 'value' }, fmt.temp(r.value))))));
}

const LEVELS = { good: 'low', warning: 'moderate', serious: 'high', critical: 'very high' };

/* ---------- dashboards ---------- */

const dashboards = {
  today(d) {
    const tiles = h('div', { class: 'tiles' });
    const tariff = d.lowTariff == null ? '—' : d.lowTariff ? 'low' : 'high';
    tiles.append(tile('tariff', tariff, null, {
      hint: d.lowTariff == null ? 'no data' : h('span', {}, h('span', { class: 'dot', style: `background:var(--${d.lowTariff ? 'good' : 'warning'})` }), d.lowTariff ? ' heures creuses' : ' heures pleines'),
    }));
    if (d.wattNow != null) {
      tiles.append(tile('now', ...fmt.watt(d.wattNow).split(/ (?=\S+$)/), {
        raw: d.wattNow, max: 10000, hint: (st) => LEVELS[st],
        thresholds: [[-Infinity, 'good'], [500, 'warning'], [3000, 'serious'], [6000, 'critical']],
      }));
    }
    tiles.append(tile('today', round(d.wattConsumedToday / 1000, 1), 'kWh', {
      raw: d.wattConsumedToday, max: 80000, hint: (st) => LEVELS[st],
      thresholds: [[-Infinity, 'good'], [20000, 'warning'], [35000, 'serious'], [50000, 'critical']],
    }));
    tiles.append(tile('cost', round(d.euroSpentToday, 2), '€', {
      raw: d.euroSpentToday, max: 8, hint: (st) => LEVELS[st],
      thresholds: [[-Infinity, 'good'], [2, 'warning'], [4, 'serious'], [6, 'critical']],
    }));
    tiles.append(tile('low tariff ratio', round(d.ratioLowTariffToday * 100, 1), '%', {
      raw: d.ratioLowTariffToday, max: 1,
      hint: (st) => ({ good: 'good', warning: 'fair', serious: 'poor', critical: 'bad' })[st],
      thresholds: [[-Infinity, 'critical'], [0.3, 'serious'], [0.4, 'warning'], [0.5, 'good']],
    }));

    const temps = h('div', { class: 'columns' });
    if (d.inside.length) temps.append(readings('inside', d.inside));
    if (d.outside.length) temps.append(readings('outside', d.outside));

    return [
      tiles,
      temps,
      timeChart('Temperature', { series: d.temperatures, from: d.from, to: d.to, unit: 'temp', height: 320 }),
      timeChart('Electricity consumption', { series: d.power, from: d.from, to: d.to, unit: 'watt', kinds: ['area', 'line'], zero: true, height: 260 }),
    ];
  },

  quarter(d) {
    // daily values are stamped at midnight UTC
    const range = { from: d.from - d.from % 86400 - 43200, to: d.to - d.to % 86400 + 43200, utc: true };
    const period = h('div', { class: 'switch', role: 'group', 'aria-label': 'period' }, d.availableDays.map((days) =>
      h('a', { href: `#quarter-${days}`, ...(days === d.days ? { 'aria-current': 'true' } : {}) }, `${days} days`)));
    return [
      period,
      h('div', { class: 'tiles' },
        trendTile('electricity', d.consumptionTrend, 'wattHour'),
        trendTile('electricity cost', d.costTrend, 'euro'),
        trendTile('water', d.waterTrend, 'liter'),
        trendTile('outside temperature', d.temperatureTrend, 'temp', 'neutral')),
      timeChart('Daily temperature', { ...range, series: d.temperatures, unit: 'temp' }),
      timeChart('Daily electricity consumption', { ...range, series: [d.consumption], unit: 'wattHour', kind: 'area', zero: true }),
      timeChart('Daily water consumption', { ...range, series: [d.water], unit: 'liter', kind: 'bars' }),
      timeChart('Low tariff rate', { ...range, series: [d.lowTariffRate], unit: 'percent', kind: 'bars', max: 100 },
        'share of the electricity consumed during low tariff hours'),
      h('div', { class: 'columns' },
        barChart('Weekly electricity consumption in €', { buckets: d.weeklyCost, unit: 'euro' }),
        barChart('Weekly water consumption in €', { buckets: d.weeklyWaterCost, unit: 'euro' },
          `${round(d.waterEuroByLiter * 1000, 2)} € per m³`)),
    ];
  },

  history(d) {
    const range = { from: d.from, to: d.to, utc: true };
    return [
      h('div', { class: 'tiles' },
        trendTile('electricity since january', d.yearToDateTrend, 'wattHour'),
        trendTile('electricity over the last 12 months', d.twelveMonthsTrend, 'wattHour')),
      timeChart('Daily consumption', { ...range, series: bySeason(d.consumption), unit: 'wattHour', kind: 'dots', zero: true }),
      timeChart('Daily temperature', { ...range, series: bySeason(d.temperature), unit: 'temp', kind: 'dots' }),
      barChart('Monthly consumption', { buckets: d.monthlyConsumption, unit: 'wattHour' },
        'hover a month to compare it with the same month, a year before'),
      barChart('Yearly consumption', { buckets: d.yearlyConsumption, unit: 'wattHour' }),
      h('div', { class: 'columns' },
        barChart('High consumption during cold days', { buckets: d.highConsumptionOnCold, unit: 'count' },
          `days below ${fmt.temp(d.coldDayTemperature)} with more than ${fmt.wattHour(d.highConsumptionWatt)}`),
        barChart('Hot days', { buckets: d.hotDays, unit: 'count' },
          `days above ${fmt.temp(d.hotDayTemperature)}`)),
    ];
  },
};

/* ---------- routing ---------- */

// the application may be served behind a reverse proxy with a path prefix: every url is relative.
const base = location.pathname.replace(/index\.html$/, '').replace(/\/?$/, '/');
const main = document.getElementById('dashboard');
const status = document.getElementById('status');
let timer = null;

async function load() {
  // #quarter-180 displays the quarter dashboard over 180 days
  const [hash, param] = location.hash.slice(1).split('-');
  const name = dashboards[hash] ? hash : 'today';
  const query = name === 'quarter' && param ? `?days=${encodeURIComponent(param)}` : '';
  document.title = `domotik / ${name}`;
  for (const a of document.querySelectorAll('nav a')) {
    if (a.getAttribute('href') === '#' + name) a.setAttribute('aria-current', 'page'); else a.removeAttribute('aria-current');
  }

  clearTimeout(timer);
  try {
    const response = await fetch(`${base}api/${name}${query}`);
    if (!response.ok) throw new Error(await response.text());
    const data = await response.json();
    hideTooltip();
    charts.length = 0;
    main.replaceChildren(...dashboards[name](data));
    charts.forEach((draw) => draw());
    const now = new Date();
    status.textContent = `updated at ${pad(now.getHours())}:${pad(now.getMinutes())}`;
  } catch (error) {
    status.textContent = `unable to load: ${error.message}`;
  }
  if (name === 'today') timer = setTimeout(load, 60000);
}

let resizing = null;
addEventListener('resize', () => {
  clearTimeout(resizing);
  resizing = setTimeout(() => charts.forEach((draw) => draw()), 100);
});
addEventListener('hashchange', load);
load();
